package simpletime

import (
	"context"
	"fmt"
	"time"
)

type Module struct{}

func (m *Module) Name() string {
	return "time"
}

func (m *Module) Run(ctx context.Context, inbox chan string, outbox chan string) error {
	enabled := true
	for {
		select {
		case msg := <-inbox:
			fmt.Println(msg)
			enabled = false
			outbox <- `<div id="time">wow</div>`
		case <-time.After(1 * time.Second):
			if enabled {
				outbox <- fmt.Sprintf(
					`<div id="time">%s</div>`,
					time.Now().Format("15:04:05"),
				)
			}
		case <-ctx.Done():
			return nil
		}
	}
}
