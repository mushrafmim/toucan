package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"toucan/internal/app"
	"toucan/internal/config"
)

func main() {
	logger := log.New(os.Stdout, "toucan ", log.LstdFlags)
	cfg := config.Load()
	application, err := app.New(cfg.Database, cfg.Storage, cfg.Seed, cfg.Identity, logger)
	if err != nil {
		logger.Fatalf("bootstrap application: %v", err)
	}
	defer func() {
		if closeErr := application.Close(); closeErr != nil {
			logger.Printf("close application: %v", closeErr)
		}
	}()

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           application.Handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.Printf("starting toucan on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("server error: %v", err)
	}
}
