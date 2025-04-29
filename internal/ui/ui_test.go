package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lauren/nannytracker/internal/maps"
	"github.com/lauren/nannytracker/internal/model"
	"github.com/lauren/nannytracker/internal/storage"
)

func setupTestUI(t *testing.T) (*Model, func()) {
	store := storage.New("testdata/trips.json")
	mockClient := maps.NewMockClient()
	model, err := NewWithClient(store, 0.655, mockClient)
	if err != nil {
		t.Fatalf("Failed to create UI model: %v", err)
	}
	return model, func() {}
}

func TestTripCreation(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Test date input
	uiModel.TextInput.SetValue("2024-03-20")
	var updatedModel tea.Model
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	if uiModel.CurrentTrip.Date != "2024-03-20" {
		t.Errorf("Expected date to be '2024-03-20', got '%s'", uiModel.CurrentTrip.Date)
	}

	if uiModel.Mode != "origin" {
		t.Errorf("Expected mode to be 'origin', got '%s'", uiModel.Mode)
	}

	// Test origin input
	uiModel.TextInput.SetValue("123 Main St")
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	if uiModel.CurrentTrip.Origin != "123 Main St" {
		t.Errorf("Expected origin to be '123 Main St', got '%s'", uiModel.CurrentTrip.Origin)
	}

	if uiModel.Mode != "destination" {
		t.Errorf("Expected mode to be 'destination', got '%s'", uiModel.Mode)
	}

	// Test destination input
	uiModel.TextInput.SetValue("456 Oak Ave")
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	if len(uiModel.Trips) != 1 {
		t.Errorf("Expected 1 trip, got %d", len(uiModel.Trips))
	}

	if uiModel.Trips[0].Date != "2024-03-20" || uiModel.Trips[0].Origin != "123 Main St" || uiModel.Trips[0].Destination != "456 Oak Ave" {
		t.Errorf("Trip data doesn't match input. Got date: %s, origin: %s, destination: %s",
			uiModel.Trips[0].Date, uiModel.Trips[0].Origin, uiModel.Trips[0].Destination)
	}
}

func TestAddingTrip(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Test adding a trip
	trip := model.Trip{
		Date:        "2024-03-20",
		Origin:      "Home",
		Destination: "Work",
		Miles:       5.0,
	}

	uiModel.AddTrip(trip)

	if len(uiModel.Trips) != 1 {
		t.Errorf("Expected 1 trip after adding, got %d", len(uiModel.Trips))
	}

	if uiModel.Trips[0].Date != "2024-03-20" || uiModel.Trips[0].Origin != "Home" || uiModel.Trips[0].Destination != "Work" {
		t.Errorf("Added trip data doesn't match. Got date: %s, origin: %s, destination: %s",
			uiModel.Trips[0].Date, uiModel.Trips[0].Origin, uiModel.Trips[0].Destination)
	}
}

func TestUIStateTransitions(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Test initial state
	if uiModel.Mode != "date" {
		t.Errorf("Expected initial mode to be 'date', got '%s'", uiModel.Mode)
	}

	// Test transition to origin mode
	uiModel.TextInput.SetValue("2024-03-20")
	updatedModel, _ := uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	if uiModel.Mode != "origin" {
		t.Errorf("Expected mode to be 'origin' after date input, got '%s'", uiModel.Mode)
	}

	// Test transition to destination mode
	uiModel.TextInput.SetValue("123 Main St")
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	if uiModel.Mode != "destination" {
		t.Errorf("Expected mode to be 'destination' after origin input, got '%s'", uiModel.Mode)
	}

	// Test transition back to date mode after trip completion
	uiModel.TextInput.SetValue("456 Oak Ave")
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	if uiModel.Mode != "date" {
		t.Errorf("Expected mode to be 'date' after trip completion, got '%s'", uiModel.Mode)
	}
}

func TestWeeklySummaryDisplay(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Add trips for different weeks
	trips := []model.Trip{
		{Date: "2024-03-17", Origin: "Home", Destination: "Work", Miles: 10.0}, // Week 1
		{Date: "2024-03-18", Origin: "Work", Destination: "Home", Miles: 15.0}, // Week 1
		{Date: "2024-03-24", Origin: "Home", Destination: "Work", Miles: 20.0}, // Week 2
		{Date: "2024-03-25", Origin: "Work", Destination: "Home", Miles: 25.0}, // Week 2
	}

	for _, trip := range trips {
		uiModel.AddTrip(trip)
	}

	// Get the view
	view := uiModel.View()

	// Check if weekly summaries are displayed
	expectedSummaries := []string{
		"Week of 2024-03-17 to 2024-03-23:",
		"  Total Miles: 25.00",
		"  Amount Owed: $16.38",
		"Week of 2024-03-24 to 2024-03-30:",
		"  Total Miles: 45.00",
		"  Amount Owed: $29.48",
	}

	for _, expected := range expectedSummaries {
		if !strings.Contains(view, expected) {
			t.Errorf("View does not contain expected weekly summary: %s", expected)
		}
	}
}
