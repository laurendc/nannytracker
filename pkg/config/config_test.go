package config

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestEnv(t *testing.T) (string, func()) {
	tempDir, err := os.MkdirTemp("", "nannytracker-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create the .nannytracker directory
	dataDir := filepath.Join(tempDir, ".nannytracker")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		t.Fatalf("Failed to create data dir: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

func TestNew(t *testing.T) {
	tempDir, cleanup := setupTestEnv(t)
	defer cleanup()

	// Set environment variables for testing
	os.Setenv("NANNYTRACKER_DATA_DIR", filepath.Join(tempDir, ".nannytracker"))
	os.Setenv("NANNYTRACKER_DATA_FILE", "test_trips.json")
	os.Setenv("NANNYTRACKER_RATE_PER_MILE", "0.655")

	cfg, err := New()
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	// Verify config values
	if cfg.DataDir != filepath.Join(tempDir, ".nannytracker") {
		t.Errorf("Expected DataDir to be %s, got %s", filepath.Join(tempDir, ".nannytracker"), cfg.DataDir)
	}

	if cfg.DataFile != "test_trips.json" {
		t.Errorf("Expected DataFile to be test_trips.json, got %s", cfg.DataFile)
	}

	if cfg.RatePerMile != 0.655 {
		t.Errorf("Expected RatePerMile to be 0.655, got %f", cfg.RatePerMile)
	}
}

func TestDataPath(t *testing.T) {
	tempDir, cleanup := setupTestEnv(t)
	defer cleanup()

	cfg := &Config{
		DataDir:  filepath.Join(tempDir, ".nannytracker"),
		DataFile: "trips.json",
	}

	expectedPath := filepath.Join(tempDir, ".nannytracker", "trips.json")
	if cfg.DataPath() != expectedPath {
		t.Errorf("Expected DataPath to be %s, got %s", expectedPath, cfg.DataPath())
	}
}

func TestDefaultConfig(t *testing.T) {
	// Clear environment variables to test defaults
	os.Unsetenv("NANNYTRACKER_DATA_DIR")
	os.Unsetenv("NANNYTRACKER_DATA_FILE")
	os.Unsetenv("NANNYTRACKER_RATE_PER_MILE")

	cfg, err := New()
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	// Verify default values
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	expectedDataDir := filepath.Join(homeDir, ".nannytracker")
	if cfg.DataDir != expectedDataDir {
		t.Errorf("Expected default DataDir to be %s, got %s", expectedDataDir, cfg.DataDir)
	}

	if cfg.DataFile != "trips.json" {
		t.Errorf("Expected default DataFile to be trips.json, got %s", cfg.DataFile)
	}

	if cfg.RatePerMile != 0.655 {
		t.Errorf("Expected default RatePerMile to be 0.655, got %f", cfg.RatePerMile)
	}
}
