package ui

import (
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

	// Test origin input
	uiModel.TextInput.SetValue("123 Main St")
	var updatedModel tea.Model
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

	if uiModel.Trips[0].Origin != "123 Main St" || uiModel.Trips[0].Destination != "456 Oak Ave" {
		t.Errorf("Trip data doesn't match input. Got origin: %s, destination: %s",
			uiModel.Trips[0].Origin, uiModel.Trips[0].Destination)
	}
}

func TestAddingTrip(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Test adding a trip
	trip := model.Trip{
		Origin:      "Home",
		Destination: "Work",
		Miles:       5.0,
	}

	uiModel.AddTrip(trip)

	if len(uiModel.Trips) != 1 {
		t.Errorf("Expected 1 trip after adding, got %d", len(uiModel.Trips))
	}

	if uiModel.Trips[0].Origin != "Home" || uiModel.Trips[0].Destination != "Work" {
		t.Errorf("Added trip data doesn't match. Got origin: %s, destination: %s",
			uiModel.Trips[0].Origin, uiModel.Trips[0].Destination)
	}
}

func TestUIStateTransitions(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Test initial state
	if uiModel.Mode != "origin" {
		t.Errorf("Expected initial mode to be 'origin', got '%s'", uiModel.Mode)
	}

	// Test transition to destination mode
	uiModel.TextInput.SetValue("123 Main St")
	updatedModel, _ := uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	if uiModel.Mode != "destination" {
		t.Errorf("Expected mode to be 'destination' after origin input, got '%s'", uiModel.Mode)
	}

	// Test transition back to origin mode after trip completion
	uiModel.TextInput.SetValue("456 Oak Ave")
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	if uiModel.Mode != "origin" {
		t.Errorf("Expected mode to be 'origin' after trip completion, got '%s'", uiModel.Mode)
	}
}
