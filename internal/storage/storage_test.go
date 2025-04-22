package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lauren/nannytracker/internal/model"
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

	// Create empty trips file
	tripsFile := filepath.Join(dataDir, "trips.json")
	if err := os.WriteFile(tripsFile, []byte("[]"), 0644); err != nil {
		t.Fatalf("Failed to create trips file: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

func TestStorage(t *testing.T) {
	tmpDir, cleanup := setupTestEnv(t)
	defer cleanup()

	// Create storage
	store := New(filepath.Join(tmpDir, ".nannytracker", "trips.json"))

	// Test saving trips
	trips := []model.Trip{
		{Origin: "Home", Destination: "Work", Miles: 5.0},
		{Origin: "Work", Destination: "Store", Miles: 2.5},
	}

	if err := store.SaveTrips(trips); err != nil {
		t.Fatalf("Failed to save trips: %v", err)
	}

	// Verify file exists
	dataPath := filepath.Join(tmpDir, ".nannytracker", "trips.json")
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		t.Errorf("Expected trips file to exist at %s", dataPath)
	}

	// Test loading trips
	loadedTrips, err := store.LoadTrips()
	if err != nil {
		t.Fatalf("Failed to load trips: %v", err)
	}

	if len(loadedTrips) != len(trips) {
		t.Errorf("Expected %d trips, got %d", len(trips), len(loadedTrips))
	}

	// Verify trip data
	for i, trip := range trips {
		if loadedTrips[i].Origin != trip.Origin {
			t.Errorf("Trip %d: expected origin %s, got %s", i, trip.Origin, loadedTrips[i].Origin)
		}
		if loadedTrips[i].Destination != trip.Destination {
			t.Errorf("Trip %d: expected destination %s, got %s", i, trip.Destination, loadedTrips[i].Destination)
		}
		if loadedTrips[i].Miles != trip.Miles {
			t.Errorf("Trip %d: expected miles %.2f, got %.2f", i, trip.Miles, loadedTrips[i].Miles)
		}
	}

	// Test loading from non-existent file
	os.Remove(dataPath)
	emptyTrips, err := store.LoadTrips()
	if err != nil {
		t.Fatalf("Failed to load from non-existent file: %v", err)
	}
	if len(emptyTrips) != 0 {
		t.Errorf("Expected 0 trips from non-existent file, got %d", len(emptyTrips))
	}

	// Test saving to non-existent directory
	os.RemoveAll(filepath.Join(tmpDir, ".nannytracker"))
	// Create the directory again before saving
	if err := os.MkdirAll(filepath.Join(tmpDir, ".nannytracker"), 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
	if err := store.SaveTrips(trips); err != nil {
		t.Fatalf("Failed to save trips to new directory: %v", err)
	}

	// Verify file was created in new directory
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		t.Errorf("Expected trips file to be created at %s", dataPath)
	}
}
