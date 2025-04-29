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

	// Create empty trips file with proper StorageData structure
	tripsFile := filepath.Join(dataDir, "trips.json")
	emptyData := `{"trips":[],"weekly_summaries":[]}`
	if err := os.WriteFile(tripsFile, []byte(emptyData), 0644); err != nil {
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

	// Test date input
	uiModel.TextInput.SetValue("2024-03-20")
	var updatedModel tea.Model
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*ui.Model)

	if uiModel.CurrentTrip.Date != "2024-03-20" {
		t.Errorf("Expected date to be '2024-03-20', got '%s'", uiModel.CurrentTrip.Date)
	}

	if uiModel.Mode != "origin" {
		t.Errorf("Expected mode to be 'origin', got '%s'", uiModel.Mode)
	}

	// Test origin input
	uiModel.TextInput.SetValue("123 Main St")
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*ui.Model)

	if uiModel.CurrentTrip.Origin != "123 Main St" {
		t.Errorf("Expected origin to be '123 Main St', got '%s'", uiModel.CurrentTrip.Origin)
	}

	if uiModel.Mode != "destination" {
		t.Errorf("Expected mode to be 'destination', got '%s'", uiModel.Mode)
	}

	// Test destination input
	uiModel.TextInput.SetValue("456 Oak Ave")
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*ui.Model)

	if len(uiModel.Trips) != 1 {
		t.Errorf("Expected 1 trip, got %d", len(uiModel.Trips))
	}

	if uiModel.Trips[0].Date != "2024-03-20" || uiModel.Trips[0].Origin != "123 Main St" || uiModel.Trips[0].Destination != "456 Oak Ave" {
		t.Errorf("Trip data doesn't match input. Got date: %s, origin: %s, destination: %s",
			uiModel.Trips[0].Date, uiModel.Trips[0].Origin, uiModel.Trips[0].Destination)
	}

	// Verify saved trips
	savedData, err := store.LoadData()
	if err != nil {
		t.Fatalf("Failed to load data: %v", err)
	}

	if len(savedData.Trips) != 1 {
		t.Errorf("Expected 1 saved trip, got %d", len(savedData.Trips))
	}

	if savedData.Trips[0].Date != "2024-03-20" || savedData.Trips[0].Origin != "123 Main St" || savedData.Trips[0].Destination != "456 Oak Ave" {
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
		{Date: "2024-03-20", Origin: "A", Destination: "B", Miles: 10.0},
		{Date: "2024-03-21", Origin: "C", Destination: "D", Miles: 15.0},
		{Date: "2024-03-22", Origin: "E", Destination: "F", Miles: 5.0},
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
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "nannytracker-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create config
	cfg := &config.Config{
		DataDir:     filepath.Join(tmpDir, ".nannytracker"),
		DataFile:    "trips.json",
		RatePerMile: 0.655,
	}

	// Ensure data directory exists
	if err := os.MkdirAll(cfg.DataDir, 0755); err != nil {
		t.Fatalf("Failed to create data directory: %v", err)
	}

	// Create storage
	store := storage.New(cfg.DataPath())

	// Test saving data
	data := &model.StorageData{
		Trips: []model.Trip{
			{Date: "2024-03-20", Origin: "Home", Destination: "Work", Miles: 5.0},
			{Date: "2024-03-21", Origin: "Work", Destination: "Store", Miles: 2.5},
		},
		WeeklySummaries: []model.WeeklySummary{
			{
				WeekStart:   "2024-03-17",
				WeekEnd:     "2024-03-23",
				TotalMiles:  7.5,
				TotalAmount: 4.91,
			},
		},
	}

	if err := store.SaveData(data); err != nil {
		t.Fatalf("Failed to save data: %v", err)
	}

	// Verify file exists
	dataPath := cfg.DataPath()
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		t.Errorf("Expected data file to exist at %s", dataPath)
	}

	// Test loading data
	loadedData, err := store.LoadData()
	if err != nil {
		t.Fatalf("Failed to load data: %v", err)
	}

	if len(loadedData.Trips) != len(data.Trips) {
		t.Errorf("Expected %d trips, got %d", len(data.Trips), len(loadedData.Trips))
	}

	if len(loadedData.WeeklySummaries) != len(data.WeeklySummaries) {
		t.Errorf("Expected %d weekly summaries, got %d", len(data.WeeklySummaries), len(loadedData.WeeklySummaries))
	}

	// Verify trip data
	for i, trip := range data.Trips {
		if loadedData.Trips[i].Date != trip.Date {
			t.Errorf("Trip %d: expected date %s, got %s", i, trip.Date, loadedData.Trips[i].Date)
		}
		if loadedData.Trips[i].Origin != trip.Origin {
			t.Errorf("Trip %d: expected origin %s, got %s", i, trip.Origin, loadedData.Trips[i].Origin)
		}
		if loadedData.Trips[i].Destination != trip.Destination {
			t.Errorf("Trip %d: expected destination %s, got %s", i, trip.Destination, loadedData.Trips[i].Destination)
		}
		if loadedData.Trips[i].Miles != trip.Miles {
			t.Errorf("Trip %d: expected miles %.2f, got %.2f", i, trip.Miles, loadedData.Trips[i].Miles)
		}
	}

	// Verify weekly summary data
	for i, summary := range data.WeeklySummaries {
		if loadedData.WeeklySummaries[i].WeekStart != summary.WeekStart {
			t.Errorf("Summary %d: expected week start %s, got %s", i, summary.WeekStart, loadedData.WeeklySummaries[i].WeekStart)
		}
		if loadedData.WeeklySummaries[i].WeekEnd != summary.WeekEnd {
			t.Errorf("Summary %d: expected week end %s, got %s", i, summary.WeekEnd, loadedData.WeeklySummaries[i].WeekEnd)
		}
		if loadedData.WeeklySummaries[i].TotalMiles != summary.TotalMiles {
			t.Errorf("Summary %d: expected total miles %.2f, got %.2f", i, summary.TotalMiles, loadedData.WeeklySummaries[i].TotalMiles)
		}
		if loadedData.WeeklySummaries[i].TotalAmount != summary.TotalAmount {
			t.Errorf("Summary %d: expected total amount %.2f, got %.2f", i, summary.TotalAmount, loadedData.WeeklySummaries[i].TotalAmount)
		}
	}

	// Test loading from non-existent file
	os.Remove(dataPath)
	emptyData, err := store.LoadData()
	if err != nil {
		t.Fatalf("Failed to load from non-existent file: %v", err)
	}
	if len(emptyData.Trips) != 0 {
		t.Errorf("Expected 0 trips from non-existent file, got %d", len(emptyData.Trips))
	}
	if len(emptyData.WeeklySummaries) != 0 {
		t.Errorf("Expected 0 weekly summaries from non-existent file, got %d", len(emptyData.WeeklySummaries))
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
		Date:        "2024-03-20",
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

	// Verify weekly summaries were updated
	if len(uiModel.Data.WeeklySummaries) != 1 {
		t.Errorf("Expected 1 weekly summary, got %d", len(uiModel.Data.WeeklySummaries))
	}

	summary := uiModel.Data.WeeklySummaries[0]
	if summary.WeekStart != "2024-03-17" || summary.WeekEnd != "2024-03-23" {
		t.Errorf("Expected week range 2024-03-17 to 2024-03-23, got %s to %s", summary.WeekStart, summary.WeekEnd)
	}
	if summary.TotalMiles != 10.5 {
		t.Errorf("Expected total miles 10.5, got %.2f", summary.TotalMiles)
	}

	// Calculate exact expected amount
	expectedAmount := 10.5 * cfg.RatePerMile
	if summary.TotalAmount != expectedAmount {
		t.Errorf("Expected total amount %.4f, got %.4f", expectedAmount, summary.TotalAmount)
	}
}
