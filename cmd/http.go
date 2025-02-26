package cmd

import (
	"log"
	"os"

	"github.com/saygik/go-url-shortener/internal/app"
)

func Execute() {

	a, err := app.NewApp()
	if err != nil {
		log.Fatalf("Config error: %s", err)
		os.Exit(1)
	}
	err = a.StartHTTPServer()
	if err != nil {
		log.Fatalf("Cannot start server: %s", err)
		os.Exit(2)
	}
}
