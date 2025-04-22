package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lauren/nannytracker/internal/storage"
	"github.com/lauren/nannytracker/internal/ui"
	"github.com/lauren/nannytracker/pkg/config"
)

func main() {
	// Load configuration
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize storage
	store := storage.New(cfg.DataPath())

	// Initialize UI
	model, err := ui.New(store, cfg.RatePerMile)
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
