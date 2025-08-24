package main

import (
	"log"

	"github.com/monoguard/api/internal/app"
)

func main() {
	// Create and start application
	application, err := app.New()
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	// Start the application
	if err := application.Start(); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}