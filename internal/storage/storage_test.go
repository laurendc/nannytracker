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
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "nannytracker-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create the storage file path
	filePath := filepath.Join(tmpDir, "trips.json")

	// Create a new storage instance
	store := New(filePath)

	// Test saving and loading data
	data := &model.StorageData{
		Trips: []model.Trip{
			{Date: "2024-03-20", Origin: "Home", Destination: "Work", Miles: 10.0},
			{Date: "2024-03-21", Origin: "Work", Destination: "Home", Miles: 10.0},
		},
		WeeklySummaries: []model.WeeklySummary{
			{
				WeekStart:   "2024-03-17",
				WeekEnd:     "2024-03-23",
				TotalMiles:  20.0,
				TotalAmount: 13.10,
			},
		},
	}

	// Save the data
	if err := store.SaveData(data); err != nil {
		t.Fatalf("Failed to save data: %v", err)
	}

	// Load the data
	loadedData, err := store.LoadData()
	if err != nil {
		t.Fatalf("Failed to load data: %v", err)
	}

	// Verify the loaded data
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
	os.Remove(filePath)
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
