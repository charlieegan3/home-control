package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/websocket"

	"github.com/charlieegan3/home-control/pkg/config"
	"github.com/charlieegan3/home-control/pkg/handlers"
)

func main() {
	var err error

	var port string
	var addr string
	var configPath string

	flag.StringVar(&port, "port", "3000", "Set server port")
	flag.StringVar(&addr, "addr", "127.0.0.1", "Set server address")
	flag.StringVar(&configPath, "config", "config.yaml", "Set config file")
	flag.Parse()

	log.Println("Reading from config file", configPath)

	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	cfg, err := config.Parse(file)
	if err != nil {
		log.Fatalf("error parsing config: %v", err)
	}

	opts := &handlers.Options{}

	if os.Getenv("GO_ENV") == "dev" {
		opts.DevMode = true
	}

	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	go func() {
		<-sigChan
		cancel()
	}()

	mux := http.NewServeMux()

	stylesEtag, cssHandler, err := handlers.BuildCSSHandler(opts)
	if err != nil {
		panic(err)
	}
	mux.HandleFunc("/styles.css", cssHandler)
	opts.EtagScript = stylesEtag

	scriptEtag, jsHandler, err := handlers.BuildJSHandler(opts)
	if err != nil {
		panic(err)
	}
	mux.HandleFunc("/script.js", jsHandler)
	opts.EtagStyles = scriptEtag

	mux.HandleFunc(
		"/fonts/",
		handlers.BuildFontHandler(opts),
	)
	mux.HandleFunc(
		"/static/",
		handlers.BuildStaticHandler(opts),
	)
	mux.HandleFunc("/favicon.ico", handlers.BuildFaviconHandler(opts))

	plugsHandler, err := handlers.BuildPlugsHandler(opts, cfg)
	if err != nil {
		log.Fatalf("failed to build plugs handler: %v", err)
	}
	mux.HandleFunc("/plugs", plugsHandler)
	mux.Handle("/plugs/websocket", websocket.Handler(handlers.BuildPlugsWebsocketHandler(ctx, cfg)))

	indexHandler, err := handlers.BuildIndexHandler(opts)
	if err != nil {
		log.Fatalf("failed to build index handler: %v", err)
	}
	mux.HandleFunc("/", indexHandler)

	serverAddrPort := fmt.Sprintf("%s:%s", addr, port)
	fmt.Fprintln(
		os.Stdout,
		"Starting server on",
		fmt.Sprintf("http://%s", serverAddrPort),
	)
	server := &http.Server{
		Addr:    serverAddrPort,
		Handler: mux,
	}

	err = server.ListenAndServe()
	if err != nil {
		fmt.Fprint(os.Stderr, fmt.Errorf("failed to start server: %s", err))
	}
}
