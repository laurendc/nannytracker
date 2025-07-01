package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// LoadEnv loads the .env file from the project root directory
// This function can be called from any subdirectory within the project
func LoadEnv() {
	// Get the project root directory by looking for go.mod
	projectRoot, err := findProjectRoot()
	if err != nil {
		log.Printf("Warning: Could not determine project root: %v", err)
		// Fall back to current directory
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: Error loading .env file: %v", err)
		}
		return
	}

	// Load .env file from project root
	envPath := filepath.Join(projectRoot, ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Warning: Error loading .env file from %s: %v", envPath, err)
	}
}

// findProjectRoot walks up the directory tree to find the project root
// (where go.mod is located)
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		// Check if go.mod exists in current directory
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// We've reached the filesystem root
			return "", filepath.ErrBadPattern
		}
		dir = parent
	}
}
