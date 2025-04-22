package config

import (
	"os"
	"path/filepath"
)

const (
	DefaultRatePerMile = 0.655
	DefaultDataFile    = "trips.json"
)

type Config struct {
	RatePerMile float64
	DataFile    string
	DataDir     string
}

func New() (*Config, error) {
	// Check environment variables first
	dataDir := os.Getenv("NANNYTRACKER_DATA_DIR")
	dataFile := os.Getenv("NANNYTRACKER_DATA_FILE")
	ratePerMile := DefaultRatePerMile

	// If no environment variables are set, use defaults
	if dataDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		dataDir = filepath.Join(homeDir, ".nannytracker")
	}

	if dataFile == "" {
		dataFile = DefaultDataFile
	}

	// Create the data directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	return &Config{
		RatePerMile: ratePerMile,
		DataFile:    dataFile,
		DataDir:     dataDir,
	}, nil
}

func (c *Config) DataPath() string {
	return filepath.Join(c.DataDir, c.DataFile)
}
