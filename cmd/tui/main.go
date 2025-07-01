package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	tui "github.com/laurendc/nannytracker/internal/tui"
	"github.com/laurendc/nannytracker/pkg/config"
	"github.com/laurendc/nannytracker/pkg/core/maps"
	"github.com/laurendc/nannytracker/pkg/core/storage"
	"github.com/laurendc/nannytracker/pkg/version"
)

func main() {
	// Parse command line flags
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&showVersion, "v", false, "Show version information")
	flag.Parse()

	// Show version if requested
	if showVersion {
		fmt.Println(version.FullString())
		os.Exit(0)
	}

	// Load .env file from project root
	config.LoadEnv()

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
