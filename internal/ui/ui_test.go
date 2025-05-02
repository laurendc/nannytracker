package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lauren/nannytracker/internal/maps"
	"github.com/lauren/nannytracker/internal/model"
	"github.com/lauren/nannytracker/internal/storage"
)

func setupTestUI(t *testing.T) (*Model, func()) {
	// Create a temporary directory for test data
	tempDir, err := os.MkdirTemp("", "nannytracker-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create storage file path
	storageFile := filepath.Join(tempDir, "trips.json")

	store := storage.New(storageFile)
	mockClient := maps.NewMockClient()
	model, err := NewWithClient(store, 0.655, mockClient)
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create UI model: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return model, cleanup
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

	if uiModel.Mode != "type" {
		t.Errorf("Expected mode to be 'type', got '%s'", uiModel.Mode)
	}

	// Test trip type input
	uiModel.TextInput.SetValue("round")
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	// Check for errors
	if uiModel.Err != nil {
		t.Errorf("Unexpected error: %v", uiModel.Err)
	}

	// Verify the trip was created with the correct data
	if len(uiModel.Trips) != 1 {
		t.Errorf("Expected 1 trip, got %d", len(uiModel.Trips))
	}

	trip := uiModel.Trips[0]
	if trip.Date != "2024-03-20" {
		t.Errorf("Expected date to be '2024-03-20', got '%s'", trip.Date)
	}
	if trip.Origin != "123 Main St" {
		t.Errorf("Expected origin to be '123 Main St', got '%s'", trip.Origin)
	}
	if trip.Destination != "456 Oak Ave" {
		t.Errorf("Expected destination to be '456 Oak Ave', got '%s'", trip.Destination)
	}
	if trip.Miles != 10.0 {
		t.Errorf("Expected miles to be 10.0, got %.2f", trip.Miles)
	}
	if trip.Type != "round" {
		t.Errorf("Expected type to be 'round', got '%s'", trip.Type)
	}

	// Verify the trip is valid
	if err := trip.Validate(); err != nil {
		t.Errorf("Trip validation failed: %v", err)
	}
}

func TestInvalidTripType(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Set up a trip with date, origin, and destination
	uiModel.TextInput.SetValue("2024-03-20")
	updatedModel, _ := uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	uiModel.TextInput.SetValue("123 Main St")
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	uiModel.TextInput.SetValue("456 Oak Ave")
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	// Try invalid trip type
	uiModel.TextInput.SetValue("invalid")
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	if uiModel.Err == nil {
		t.Error("Expected error for invalid trip type")
	}
	if !strings.Contains(uiModel.Err.Error(), "invalid trip type") {
		t.Errorf("Expected error about invalid trip type, got: %v", uiModel.Err)
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

	// Test transition to type mode
	uiModel.TextInput.SetValue("456 Oak Ave")
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	if uiModel.Mode != "type" {
		t.Errorf("Expected mode to be 'type' after destination input, got '%s'", uiModel.Mode)
	}

	// Test transition back to date mode after trip completion
	uiModel.TextInput.SetValue("single")
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

	// Add expenses for different weeks
	expenses := []model.Expense{
		{Date: "2024-03-17", Amount: 25.50, Description: "Lunch"},      // Week 1
		{Date: "2024-03-18", Amount: 15.75, Description: "Snacks"},     // Week 1
		{Date: "2024-03-24", Amount: 30.00, Description: "Activities"}, // Week 2
	}

	for _, trip := range trips {
		uiModel.AddTrip(trip)
	}

	for _, expense := range expenses {
		if err := uiModel.Data.AddExpense(expense); err != nil {
			t.Fatalf("Failed to add expense: %v", err)
		}
	}

	model.CalculateAndUpdateWeeklySummaries(uiModel.Data, uiModel.RatePerMile)

	// Get the view
	view := uiModel.View()

	// Check if weekly summaries are displayed with expenses
	expectedSummaries := []string{
		"Week of 2024-03-17 to 2024-03-23:",
		"  Total Miles: 25.00",
		"  Total Mileage Amount: $16.38",
		"  Total Expenses: $41.25",
		"Week of 2024-03-24 to 2024-03-30:",
		"  Total Miles: 45.00",
		"  Total Mileage Amount: $29.48",
		"  Total Expenses: $30.00",
	}

	for _, expected := range expectedSummaries {
		if !strings.Contains(view, expected) {
			t.Errorf("View does not contain expected weekly summary: %s", expected)
		}
	}
}

func TestExpenseEntry(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Test entering an expense
	expense := model.Expense{
		Date:        "2024-03-20",
		Amount:      25.50,
		Description: "Lunch",
	}

	// Simulate entering expense mode
	msg := tea.KeyMsg{Type: tea.KeyCtrlX}
	model, _ := uiModel.Update(msg)
	uiModel = model.(*Model)

	// Enter date
	uiModel.TextInput.SetValue(expense.Date)
	model, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = model.(*Model)

	// Enter amount
	uiModel.TextInput.SetValue(fmt.Sprintf("%.2f", expense.Amount))
	model, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = model.(*Model)

	// Enter description
	uiModel.TextInput.SetValue(expense.Description)
	model, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = model.(*Model)

	// Verify expense was added
	if len(uiModel.Data.Expenses) != 1 {
		t.Errorf("Expected 1 expense, got %d", len(uiModel.Data.Expenses))
	}

	addedExpense := uiModel.Data.Expenses[0]
	if addedExpense.Date != expense.Date {
		t.Errorf("Expected date %s, got %s", expense.Date, addedExpense.Date)
	}
	if addedExpense.Amount != expense.Amount {
		t.Errorf("Expected amount %.2f, got %.2f", expense.Amount, addedExpense.Amount)
	}
	if addedExpense.Description != expense.Description {
		t.Errorf("Expected description %s, got %s", expense.Description, addedExpense.Description)
	}

	// Verify weekly summary was updated
	if len(uiModel.Data.WeeklySummaries) != 1 {
		t.Errorf("Expected 1 weekly summary, got %d", len(uiModel.Data.WeeklySummaries))
	}

	summary := uiModel.Data.WeeklySummaries[0]
	if summary.TotalExpenses != expense.Amount {
		t.Errorf("Expected total expenses %.2f, got %.2f", expense.Amount, summary.TotalExpenses)
	}
}

func TestExpenseValidation(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Enter expense mode first
	msg := tea.KeyMsg{Type: tea.KeyCtrlX}
	model, _ := uiModel.Update(msg)
	uiModel = model.(*Model)

	if uiModel.Mode != "expense_date" {
		t.Errorf("Expected mode to be 'expense_date', got '%s'", uiModel.Mode)
	}

	// Test invalid date
	uiModel.TextInput.SetValue("invalid-date")
	model, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = model.(*Model)

	// Verify we're still in expense_date mode and have an error
	if uiModel.Mode != "expense_date" {
		t.Errorf("Expected to stay in expense_date mode after invalid date, got '%s'", uiModel.Mode)
	}
	if uiModel.Err == nil || !strings.Contains(uiModel.Err.Error(), "date must be in YYYY-MM-DD format") {
		t.Errorf("Expected error about invalid date format, got: %v", uiModel.Err)
	}

	// Test valid date but invalid amount
	uiModel.TextInput.SetValue("2024-03-20")
	model, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = model.(*Model)

	if uiModel.Mode != "expense_amount" {
		t.Errorf("Expected mode to be 'expense_amount', got '%s'", uiModel.Mode)
	}

	uiModel.TextInput.SetValue("invalid-amount")
	model, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = model.(*Model)

	// Verify we're still in expense_amount mode and have an error
	if uiModel.Mode != "expense_amount" {
		t.Errorf("Expected to stay in expense_amount mode after invalid amount, got '%s'", uiModel.Mode)
	}
	if uiModel.Err == nil || !strings.Contains(uiModel.Err.Error(), "invalid amount") {
		t.Errorf("Expected error about invalid amount, got: %v", uiModel.Err)
	}

	// Test negative amount
	uiModel.TextInput.SetValue("-10.00")
	model, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = model.(*Model)

	// Verify we're still in expense_amount mode and have an error
	if uiModel.Mode != "expense_amount" {
		t.Errorf("Expected to stay in expense_amount mode after negative amount, got '%s'", uiModel.Mode)
	}
	if uiModel.Err == nil || !strings.Contains(uiModel.Err.Error(), "amount must be greater than 0") {
		t.Errorf("Expected error about negative amount, got: %v", uiModel.Err)
	}

	// Test valid amount but empty description
	uiModel.TextInput.SetValue("25.50")
	model, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = model.(*Model)

	if uiModel.Mode != "expense_description" {
		t.Errorf("Expected mode to be 'expense_description', got '%s'", uiModel.Mode)
	}

	uiModel.TextInput.SetValue("")
	model, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = model.(*Model)

	// Verify we're still in expense_description mode and have an error
	if uiModel.Mode != "expense_description" {
		t.Errorf("Expected to stay in expense_description mode after empty description, got '%s'", uiModel.Mode)
	}
	if uiModel.Err == nil || !strings.Contains(uiModel.Err.Error(), "description cannot be empty") {
		t.Errorf("Expected error about empty description, got: %v", uiModel.Err)
	}
}

func TestEditTrip(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Add a trip first
	originalTrip := model.Trip{
		Date:        "2024-03-20",
		Origin:      "Home",
		Destination: "Work",
		Miles:       10.5,
		Type:        "single",
	}
	uiModel.AddTrip(originalTrip)

	// Select the trip
	uiModel.SelectedTrip = 0

	// Enter edit mode
	var updatedModel tea.Model
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyCtrlE})
	uiModel = updatedModel.(*Model)

	if uiModel.Mode != "edit" {
		t.Errorf("Expected mode to be 'edit', got '%s'", uiModel.Mode)
	}

	// Verify initial edit state
	if uiModel.EditIndex != 0 {
		t.Errorf("Expected EditIndex to be 0, got %d", uiModel.EditIndex)
	}
	if uiModel.CurrentTrip != originalTrip {
		t.Errorf("Expected CurrentTrip to match original trip")
	}
	if uiModel.TextInput.Value() != originalTrip.Date {
		t.Errorf("Expected TextInput value to be '%s', got '%s'", originalTrip.Date, uiModel.TextInput.Value())
	}

	// Edit the date
	newDate := "2024-03-21"
	uiModel.TextInput.SetValue(newDate)
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	// Edit origin
	uiModel.TextInput.SetValue("Updated Home")
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	// Edit destination
	uiModel.TextInput.SetValue("Updated Work")
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	// Edit trip type
	uiModel.TextInput.SetValue("round")
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	// Verify final state
	if len(uiModel.Trips) != 1 {
		t.Errorf("Expected 1 trip, got %d", len(uiModel.Trips))
	}

	editedTrip := uiModel.Trips[0]
	if editedTrip.Date != newDate {
		t.Errorf("Expected final date to be '%s', got '%s'", newDate, editedTrip.Date)
	}
	if editedTrip.Origin != "Updated Home" {
		t.Errorf("Expected origin to be 'Updated Home', got '%s'", editedTrip.Origin)
	}
	if editedTrip.Destination != "Updated Work" {
		t.Errorf("Expected destination to be 'Updated Work', got '%s'", editedTrip.Destination)
	}
	if editedTrip.Type != "round" {
		t.Errorf("Expected type to be 'round', got '%s'", editedTrip.Type)
	}

	// Verify edit mode was cleared
	if uiModel.Mode != "date" {
		t.Errorf("Expected mode to reset to 'date', got '%s'", uiModel.Mode)
	}
	if uiModel.EditIndex != -1 {
		t.Errorf("Expected EditIndex to reset to -1, got %d", uiModel.EditIndex)
	}
}

func TestDeleteTrip(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Add a trip first
	trip := model.Trip{
		Date:        "2024-03-20",
		Origin:      "Home",
		Destination: "Work",
		Miles:       10.5,
	}
	uiModel.AddTrip(trip)

	// Select the trip
	uiModel.SelectedTrip = 0

	// Enter delete confirmation mode
	var updatedModel tea.Model
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyCtrlD})
	uiModel = updatedModel.(*Model)

	// Verify we're in delete confirmation mode
	if uiModel.Mode != "delete_confirm" {
		t.Errorf("Expected mode to be 'delete_confirm', got '%s'", uiModel.Mode)
	}

	// Test cancellation by entering something other than 'yes'
	uiModel.TextInput.SetValue("no")
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	// Verify trip wasn't deleted and mode was reset
	if len(uiModel.Trips) != 1 {
		t.Errorf("Expected 1 trip after cancellation, got %d", len(uiModel.Trips))
	}
	if uiModel.Mode != "date" {
		t.Errorf("Expected mode to be 'date' after cancellation, got '%s'", uiModel.Mode)
	}

	// Enter delete confirmation mode again
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyCtrlD})
	uiModel = updatedModel.(*Model)

	// Confirm deletion by entering 'yes'
	uiModel.TextInput.SetValue("yes")
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	// Verify trip was deleted and mode was reset
	if len(uiModel.Trips) != 0 {
		t.Errorf("Expected 0 trips after deletion, got %d", len(uiModel.Trips))
	}
	if uiModel.Mode != "date" {
		t.Errorf("Expected mode to be 'date' after deletion, got '%s'", uiModel.Mode)
	}
}

func TestTripSelection(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Add multiple trips
	trips := []model.Trip{
		{Date: "2024-03-20", Origin: "A", Destination: "B", Miles: 10.0},
		{Date: "2024-03-21", Origin: "C", Destination: "D", Miles: 15.0},
		{Date: "2024-03-22", Origin: "E", Destination: "F", Miles: 5.0},
	}

	for _, trip := range trips {
		uiModel.AddTrip(trip)
	}

	// Initialize selection
	uiModel.SelectedTrip = -1

	// Test selection movement
	var updatedModel tea.Model

	// Move down
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyDown})
	uiModel = updatedModel.(*Model)
	if uiModel.SelectedTrip != 0 {
		t.Errorf("Expected selected trip to be 0, got %d", uiModel.SelectedTrip)
	}

	// Move down again
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyDown})
	uiModel = updatedModel.(*Model)
	if uiModel.SelectedTrip != 1 {
		t.Errorf("Expected selected trip to be 1, got %d", uiModel.SelectedTrip)
	}

	// Move up
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyUp})
	uiModel = updatedModel.(*Model)
	if uiModel.SelectedTrip != 0 {
		t.Errorf("Expected selected trip to be 0, got %d", uiModel.SelectedTrip)
	}
}

func TestEditTripWithType(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Add a trip first
	originalTrip := model.Trip{
		Date:        "2024-03-20",
		Origin:      "Home",
		Destination: "Work",
		Miles:       10.5,
		Type:        "single",
	}
	uiModel.AddTrip(originalTrip)

	// Select the trip
	uiModel.SelectedTrip = 0

	// Enter edit mode
	var updatedModel tea.Model
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyCtrlE})
	uiModel = updatedModel.(*Model)

	if uiModel.Mode != "edit" {
		t.Errorf("Expected mode to be 'edit', got '%s'", uiModel.Mode)
	}

	// Verify initial edit state
	if uiModel.EditIndex != 0 {
		t.Errorf("Expected EditIndex to be 0, got %d", uiModel.EditIndex)
	}
	if uiModel.CurrentTrip != originalTrip {
		t.Errorf("Expected CurrentTrip to match original trip")
	}
	if uiModel.TextInput.Value() != originalTrip.Date {
		t.Errorf("Expected TextInput value to be '%s', got '%s'", originalTrip.Date, uiModel.TextInput.Value())
	}

	// Edit the date
	newDate := "2024-03-21"
	uiModel.TextInput.SetValue(newDate)
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	// Edit origin
	uiModel.TextInput.SetValue("Updated Home")
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	// Edit destination
	uiModel.TextInput.SetValue("Updated Work")
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	// Edit trip type
	uiModel.TextInput.SetValue("round")
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	// Verify final state
	if len(uiModel.Trips) != 1 {
		t.Errorf("Expected 1 trip, got %d", len(uiModel.Trips))
	}

	editedTrip := uiModel.Trips[0]
	if editedTrip.Date != newDate {
		t.Errorf("Expected final date to be '%s', got '%s'", newDate, editedTrip.Date)
	}
	if editedTrip.Origin != "Updated Home" {
		t.Errorf("Expected origin to be 'Updated Home', got '%s'", editedTrip.Origin)
	}
	if editedTrip.Destination != "Updated Work" {
		t.Errorf("Expected destination to be 'Updated Work', got '%s'", editedTrip.Destination)
	}
	if editedTrip.Type != "round" {
		t.Errorf("Expected type to be 'round', got '%s'", editedTrip.Type)
	}

	// Verify edit mode was cleared
	if uiModel.Mode != "date" {
		t.Errorf("Expected mode to reset to 'date', got '%s'", uiModel.Mode)
	}
	if uiModel.EditIndex != -1 {
		t.Errorf("Expected EditIndex to reset to -1, got %d", uiModel.EditIndex)
	}
}

func TestTripHistoryDisplay(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Add trips with different types
	trips := []model.Trip{
		{Date: "2024-03-20", Origin: "Home", Destination: "Work", Miles: 10.0, Type: "single"},
		{Date: "2024-03-21", Origin: "Work", Destination: "Store", Miles: 5.0, Type: "round"},
		{Date: "2024-03-22", Origin: "Home", Destination: "Gym", Miles: 3.0, Type: "single"},
	}

	for _, trip := range trips {
		uiModel.AddTrip(trip)
	}

	// Get the view
	view := uiModel.View()

	// Check if trip history shows correct miles
	expectedTrips := []string{
		"1. 2024-03-20 - Home → Work (10.00 miles, single)",
		"2. 2024-03-21 - Work → Store (10.00 miles, round)", // Should show doubled miles
		"3. 2024-03-22 - Home → Gym (3.00 miles, single)",
	}

	for _, expected := range expectedTrips {
		if !strings.Contains(view, expected) {
			t.Errorf("View does not contain expected trip: %s", expected)
		}
	}

	// Verify total miles calculation
	totalMiles := uiModel.CalculateTotalMiles(uiModel.Trips)
	expectedTotal := 10.0 + (5.0 * 2) + 3.0 // single + round + single
	if totalMiles != expectedTotal {
		t.Errorf("Expected total miles to be %.2f, got %.2f", expectedTotal, totalMiles)
	}
}

func TestEditTripWithSkippedFields(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Add a trip first
	originalTrip := model.Trip{
		Date:        "2024-03-20",
		Origin:      "Home",
		Destination: "Work",
		Miles:       10.5,
		Type:        "single",
	}
	uiModel.AddTrip(originalTrip)

	// Select the trip
	uiModel.SelectedTrip = 0

	// Enter edit mode
	var updatedModel tea.Model
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyCtrlE})
	uiModel = updatedModel.(*Model)

	// Skip date field (press Enter without input)
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	// Skip origin field (press Enter without input)
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	// Skip destination field (press Enter without input)
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	// Skip trip type field (press Enter without input)
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	// Verify the trip remains unchanged
	if len(uiModel.Trips) != 1 {
		t.Errorf("Expected 1 trip, got %d", len(uiModel.Trips))
	}

	editedTrip := uiModel.Trips[0]
	if editedTrip.Date != originalTrip.Date {
		t.Errorf("Expected date to remain '%s', got '%s'", originalTrip.Date, editedTrip.Date)
	}
	if editedTrip.Origin != originalTrip.Origin {
		t.Errorf("Expected origin to remain '%s', got '%s'", originalTrip.Origin, editedTrip.Origin)
	}
	if editedTrip.Destination != originalTrip.Destination {
		t.Errorf("Expected destination to remain '%s', got '%s'", originalTrip.Destination, editedTrip.Destination)
	}
	if editedTrip.Type != originalTrip.Type {
		t.Errorf("Expected type to remain '%s', got '%s'", originalTrip.Type, editedTrip.Type)
	}
}

func TestEditTripWithMixedChanges(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Add a trip first
	originalTrip := model.Trip{
		Date:        "2024-03-20",
		Origin:      "Home",
		Destination: "Work",
		Miles:       10.5,
		Type:        "single",
	}
	uiModel.AddTrip(originalTrip)

	// Select the trip
	uiModel.SelectedTrip = 0

	// Enter edit mode
	var updatedModel tea.Model
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyCtrlE})
	uiModel = updatedModel.(*Model)

	// Change date
	uiModel.TextInput.SetValue("2024-03-21")
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	// Skip origin field
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	// Change destination
	uiModel.TextInput.SetValue("Updated Work")
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	// Skip trip type field
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	// Verify the trip has mixed changes
	if len(uiModel.Trips) != 1 {
		t.Errorf("Expected 1 trip, got %d", len(uiModel.Trips))
	}

	editedTrip := uiModel.Trips[0]
	if editedTrip.Date != "2024-03-21" {
		t.Errorf("Expected date to be '2024-03-21', got '%s'", editedTrip.Date)
	}
	if editedTrip.Origin != originalTrip.Origin {
		t.Errorf("Expected origin to remain '%s', got '%s'", originalTrip.Origin, editedTrip.Origin)
	}
	if editedTrip.Destination != "Updated Work" {
		t.Errorf("Expected destination to be 'Updated Work', got '%s'", editedTrip.Destination)
	}
	if editedTrip.Type != originalTrip.Type {
		t.Errorf("Expected type to remain '%s', got '%s'", originalTrip.Type, editedTrip.Type)
	}
}

func TestExpenseHistoryDisplay(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Add expenses
	expenses := []model.Expense{
		{Date: "2024-03-20", Amount: 25.50, Description: "Lunch"},
		{Date: "2024-03-21", Amount: 15.75, Description: "Snacks"},
		{Date: "2024-03-22", Amount: 30.00, Description: "Activities"},
	}

	for _, expense := range expenses {
		if err := uiModel.Data.AddExpense(expense); err != nil {
			t.Fatalf("Failed to add expense: %v", err)
		}
	}

	// Get the view
	view := uiModel.View()

	// Check if expense history shows correct format
	expectedExpenses := []string{
		"1. 2024-03-20 - $25.50 - Lunch",
		"2. 2024-03-21 - $15.75 - Snacks",
		"3. 2024-03-22 - $30.00 - Activities",
	}

	for _, expected := range expectedExpenses {
		if !strings.Contains(view, expected) {
			t.Errorf("View does not contain expected expense: %s", expected)
		}
	}

	// Verify total expenses
	totalExpenses := model.CalculateTotalExpenses(uiModel.Data.Expenses)
	expectedTotal := 25.50 + 15.75 + 30.00
	if totalExpenses != expectedTotal {
		t.Errorf("Expected total expenses to be %.2f, got %.2f", expectedTotal, totalExpenses)
	}
}
