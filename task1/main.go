package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

type config struct {
	port          int
	windowSize    int
	maxResponseMs int
	thirdPartyAPI string
}

type application struct {
	config config
	logger *log.Logger
	window *NumberWindow
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 9876, "API server port")
	flag.IntVar(&cfg.windowSize, "windowsize", 10, "Window size for number storage")
	flag.IntVar(&cfg.maxResponseMs, "maxresponsems", 500, "Maximum response time in milliseconds")
	flag.Parse()

	logger := log.New(os.Stdout, "APP: ", log.LstdFlags|log.Lshortfile)

	app := &application{
		config: cfg,
		logger: logger,
		window: NewNumberWindow(cfg.windowSize),
	}

	router := httprouter.New()
	router.GET("/numbers/:numberId", app.handleNumbers)
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	serverAddr := fmt.Sprintf(":%d", cfg.port)
	logger.Printf("Starting server on port %d...", cfg.port)
	err := http.ListenAndServe(serverAddr, router)
	if err != nil {
		logger.Fatalf("Server error: %v", err)
	}
}
