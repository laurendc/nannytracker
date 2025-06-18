package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	tui "github.com/laurendc/nannytracker/internal/tui"
	"github.com/laurendc/nannytracker/pkg/config"
	"github.com/laurendc/nannytracker/pkg/core/maps"
	"github.com/laurendc/nannytracker/pkg/core/storage"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Load configuration
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize storage
	store := storage.New(cfg.DataPath())

	// Initialize Google Maps client
	mapsClient, err := maps.NewClient()
	if err != nil {
		log.Fatalf("Failed to initialize Google Maps client: %v", err)
	}

	// Initialize UI with Google Maps client
	model, err := tui.NewWithClient(store, cfg.RatePerMile, mapsClient)
	if err != nil {
		log.Fatalf("Failed to initialize UI: %v", err)
	}

	// Start the application
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
