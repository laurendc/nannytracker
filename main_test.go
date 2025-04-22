package main

import (
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/lauren/nannytracker/internal/maps"
	"github.com/lauren/nannytracker/internal/model"
	"github.com/lauren/nannytracker/internal/storage"
	"github.com/lauren/nannytracker/internal/ui"
	"github.com/lauren/nannytracker/pkg/config"
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

func TestTripCreation(t *testing.T) {
	tempDir, cleanup := setupTestEnv(t)
	defer cleanup()

	cfg := &config.Config{
		DataDir:     filepath.Join(tempDir, ".nannytracker"),
		DataFile:    "trips.json",
		RatePerMile: 0.655,
	}

	store := storage.New(cfg.DataPath())
	mockClient := maps.NewMockClient()

	uiModel, err := ui.NewWithClient(store, cfg.RatePerMile, mockClient)
	if err != nil {
		t.Fatalf("Failed to create UI model: %v", err)
	}

	// Test origin input
	uiModel.TextInput.SetValue("123 Main St")
	var updatedModel tea.Model
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*ui.Model)

	if uiModel.CurrentTrip.Origin != "123 Main St" {
		t.Errorf("Expected origin to be '123 Main St', got '%s'", uiModel.CurrentTrip.Origin)
	}

	if uiModel.Mode != "destination" {
		t.Errorf("Expected mode to be 'destination', got '%s'", uiModel.Mode)
	}

	// Test destination input - first update the text input
	uiModel.TextInput.SetValue("456 Oak Ave")
	// Then send the enter key
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*ui.Model)

	if len(uiModel.Trips) != 1 {
		t.Errorf("Expected 1 trip, got %d", len(uiModel.Trips))
	}

	if uiModel.Trips[0].Origin != "123 Main St" || uiModel.Trips[0].Destination != "456 Oak Ave" {
		t.Errorf("Trip data doesn't match input. Got origin: %s, destination: %s",
			uiModel.Trips[0].Origin, uiModel.Trips[0].Destination)
	}

	// Verify saved trips
	savedTrips, err := store.LoadTrips()
	if err != nil {
		t.Fatalf("Failed to load trips: %v", err)
	}

	if len(savedTrips) != 1 {
		t.Errorf("Expected 1 saved trip, got %d", len(savedTrips))
	}

	if savedTrips[0].Origin != "123 Main St" || savedTrips[0].Destination != "456 Oak Ave" {
		t.Errorf("Saved trip data doesn't match input")
	}

	// Test total miles calculation
	totalMiles := uiModel.CalculateTotalMiles(uiModel.Trips)
	if totalMiles != 10.0 {
		t.Errorf("Expected total miles to be 10.0, got %.2f", totalMiles)
	}

	// Test reimbursement calculation
	totalReimbursement := uiModel.CalculateReimbursement(uiModel.Trips, cfg.RatePerMile)
	expectedReimbursement := 10.0 * cfg.RatePerMile
	if totalReimbursement != expectedReimbursement {
		t.Errorf("Expected reimbursement to be %.2f, got %.2f", expectedReimbursement, totalReimbursement)
	}
}

func TestTotalMilesCalculation(t *testing.T) {
	tempDir, cleanup := setupTestEnv(t)
	defer cleanup()

	cfg := &config.Config{
		DataDir:     filepath.Join(tempDir, ".nannytracker"),
		DataFile:    "trips.json",
		RatePerMile: 0.70,
	}

	store := storage.New(cfg.DataPath())
	mockClient := maps.NewMockClient()

	uiModel, err := ui.NewWithClient(store, cfg.RatePerMile, mockClient)
	if err != nil {
		t.Fatalf("Failed to create UI model: %v", err)
	}

	// Add multiple trips
	trips := []model.Trip{
		{Origin: "A", Destination: "B", Miles: 10.0},
		{Origin: "C", Destination: "D", Miles: 15.0},
		{Origin: "E", Destination: "F", Miles: 5.0},
	}

	for _, trip := range trips {
		uiModel.AddTrip(trip)
	}

	totalMiles := uiModel.CalculateTotalMiles(uiModel.Trips)
	if totalMiles != 30.0 {
		t.Errorf("Expected total miles to be 30.0, got %.2f", totalMiles)
	}

	totalReimbursement := uiModel.CalculateReimbursement(uiModel.Trips, cfg.RatePerMile)
	expectedReimbursement := 30.0 * cfg.RatePerMile
	if totalReimbursement != expectedReimbursement {
		t.Errorf("Expected total reimbursement to be $%.2f, got $%.2f", expectedReimbursement, totalReimbursement)
	}
}

func TestStorage(t *testing.T) {
	tmpDir, cleanup := setupTestEnv(t)
	defer cleanup()

	// Create config
	cfg, err := config.New()
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	// Create storage
	store := storage.New(cfg.DataPath())

	// Test saving trips
	trips := []model.Trip{
		{Origin: "Home", Destination: "Work", Miles: 5.0},
		{Origin: "Work", Destination: "Store", Miles: 2.5},
	}

	if err := store.SaveTrips(trips); err != nil {
		t.Fatalf("Failed to save trips: %v", err)
	}

	// Verify file exists
	dataPath := cfg.DataPath()
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
	if err := store.SaveTrips(trips); err != nil {
		t.Fatalf("Failed to save trips to new directory: %v", err)
	}

	// Verify file was created in new directory
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		t.Errorf("Expected trips file to be created at %s", dataPath)
	}
}

func TestAddingTrip(t *testing.T) {
	tempDir, cleanup := setupTestEnv(t)
	defer cleanup()

	cfg := &config.Config{
		DataDir:     filepath.Join(tempDir, ".nannytracker"),
		DataFile:    "trips.json",
		RatePerMile: 0.655,
	}

	store := storage.New(cfg.DataPath())
	mockClient := maps.NewMockClient()

	uiModel, err := ui.NewWithClient(store, cfg.RatePerMile, mockClient)
	if err != nil {
		t.Fatalf("Failed to create UI model: %v", err)
	}

	// Test adding a trip
	trip := model.Trip{
		Origin:      "Home",
		Destination: "Work",
		Miles:       10.5,
	}

	uiModel.AddTrip(trip)

	if len(uiModel.Trips) != 1 {
		t.Errorf("Expected 1 trip, got %d", len(uiModel.Trips))
	}

	if uiModel.Trips[0] != trip {
		t.Errorf("Trip data doesn't match. Expected %+v, got %+v", trip, uiModel.Trips[0])
	}
}
