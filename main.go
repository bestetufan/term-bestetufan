package main

import (
	"log"

	"github.com/bestetufan/beste-store/config"
	"github.com/bestetufan/beste-store/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
