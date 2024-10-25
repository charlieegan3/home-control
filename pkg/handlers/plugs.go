package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"

	"github.com/charlieegan3/home-control/pkg/config"
	"github.com/charlieegan3/home-control/pkg/tasmota/plugs"
)

func BuildPlugsHandler(
	opts *Options,
	cfg *config.Config,
) (func(http.ResponseWriter, *http.Request), error) {

	tmpl, err := template.ParseFS(templates, "templates/plugs.html", "templates/base.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %s", err)
	}

	plugIndex := make(map[string]config.Plug)
	for _, plug := range cfg.Plugs {
		plugIndex[plug.ID] = plug
	}

	return func(w http.ResponseWriter, r *http.Request) {

		err := tmpl.ExecuteTemplate(w, "base", struct {
			Opts   *Options
			Plugs  map[string]config.Plug
			Groups []config.Group
		}{Opts: opts, Plugs: plugIndex, Groups: cfg.Groups})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to render template: %s", err), http.StatusInternalServerError)
			return
		}

	}, nil
}

type webSocketMessage struct {
	ID string
}

func BuildPlugsWebsocketHandler(ctx context.Context, cfg *config.Config) func(ws *websocket.Conn) {
	plugIndex := make(map[string]config.Plug)
	for _, plug := range cfg.Plugs {
		plugIndex[plug.ID] = plug
	}

	return func(ws *websocket.Conn) {
		ctx, cancel := context.WithCancel(ctx)

		statusUpdates := make(chan plugs.Status, 10)

		plugActivePower := make(map[string]int)

		go func() {
			for {
				var msg string
				err := websocket.Message.Receive(ws, &msg)
				if err != nil && err == io.EOF {
					log.Println("EOF received, closing connection")
					cancel()
					return
				}
				if err != nil {
					log.Println("Error receiving message: ", err)
					continue
				}
				if msg == "" {
					continue
				}

				var parsedMessage webSocketMessage
				err = json.Unmarshal([]byte(msg), &parsedMessage)
				if err != nil {
					log.Println("Error parsing message: ", err)
					continue
				}

				plug, ok := plugIndex[parsedMessage.ID]
				if !ok {
					log.Println("Error finding plug: ", parsedMessage.ID)
					continue
				}

				req, err := http.NewRequest(
					"GET",
					fmt.Sprintf("http://%s/?m=1&o=1", plug.Hostname),
					nil,
				)
				if err != nil {
					log.Println("Error creating request: ", err)
					continue
				}

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					log.Println("Error sending request: ", err)
					continue
				}

				if resp.StatusCode != http.StatusOK {
					log.Println("Error sending request, status code: ", resp.StatusCode)
					continue
				}

				bodyBs, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Println("Error reading response body: ", err)
					continue
				}

				status, err := plugs.ParseStatus(string(bodyBs))
				if err != nil {
					log.Println("Error parsing response body: ", err)
					continue
				}

				status.ID = parsedMessage.ID

				statusUpdates <- *status
			}
		}()

		go func() {
			for {
				select {
				case status := <-statusUpdates:
					plug, ok := plugIndex[status.ID]
					if !ok {
						log.Println("Error finding plug: ", status.ID)
						break
					}

					if status.PowerOn {
						plugActivePower[plug.ID] = status.ActivePower
					} else {
						plugActivePower[plug.ID] = 0
					}

					html := statusToHTML(plug, status)
					err := websocket.Message.Send(ws, html)
					if err != nil {
						log.Println("Error sending message: ", err)
					}

					totalWatts := 0.0
					for _, watts := range plugActivePower {
						totalWatts += float64(watts)
					}
					ppkw := 14.77
					year := totalWatts / 1000 * 24 * 365 * ppkw / 100
					month := year / 12
					day := year / 365
					err = websocket.Message.Send(ws,
						fmt.Sprintf(
							`<span id="total"><strong>%.fW</strong> (year: £%.f, month: £%.f, day: £%.1f)</span>`,
							totalWatts,
							year,
							month,
							day,
						),
					)
					if err != nil {
						log.Println("Error sending totalWattage message: ", err)
					}

					for _, group := range cfg.Groups {
						totalWatts := 0
						for _, plugID := range group.Plugs {
							totalWatts += plugActivePower[plugID]
						}
						err = websocket.Message.Send(ws,
							fmt.Sprintf(
								`<span id="group-%s"><strong>%dW</strong></span>`,
								group.ID,
								totalWatts,
							),
						)
						if err != nil {
							log.Println("Error sending group totalWattage message: ", err)
						}
					}

				case <-ctx.Done():
					return
				}
			}
		}()

		go func() {
			for {
				select {
				case <-time.After(1 * time.Second):
					for _, plug := range cfg.Plugs {
						req, err := http.NewRequest(
							"GET",
							fmt.Sprintf("http://%s/?m=1", plug.Hostname),
							nil,
						)
						if err != nil {
							log.Println("Error creating request: ", err)
							continue
						}

						client := &http.Client{}
						resp, err := client.Do(req)
						if err != nil {
							log.Println("Error sending request: ", err)
							continue
						}

						if resp.StatusCode != http.StatusOK {
							log.Println("Error sending request, status code: ", resp.StatusCode)
							continue
						}

						bodyBs, err := io.ReadAll(resp.Body)
						if err != nil {
							log.Println("Error reading response body: ", err)
							continue
						}

						status, err := plugs.ParseStatus(string(bodyBs))
						if err != nil {
							log.Println("Error parsing response body: ", err)
							continue
						}

						status.ID = plug.ID

						statusUpdates <- *status
					}
				case <-ctx.Done():
					return
				}
			}
		}()

		<-ctx.Done()
		ws.Close()
	}
}

func statusToHTML(plug config.Plug, status plugs.Status) string {
	if status.PowerOn {
		return fmt.Sprintf(
			`<span id="%s" class="green white"><strong>ON</strong> %dW</span>`,
			plug.ID,
			status.ActivePower,
		)
	}

	return fmt.Sprintf(`<span id="%s" class=""><strong>OFF</strong></span>`, plug.ID)
}
