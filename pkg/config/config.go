package config

import (
	"os"
	"path/filepath"
)

const (
	DefaultRatePerMile = 0.70
	DefaultDataFile    = "trips.json"
)

type Config struct {
	RatePerMile float64
	DataFile    string
	DataDir     string
}

func New() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dataDir := filepath.Join(homeDir, ".nannytracker")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	return &Config{
		RatePerMile: DefaultRatePerMile,
		DataFile:    DefaultDataFile,
		DataDir:     dataDir,
	}, nil
}

func (c *Config) DataPath() string {
	return filepath.Join(c.DataDir, c.DataFile)
}
