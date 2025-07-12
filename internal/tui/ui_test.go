package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	model "github.com/laurendc/nannytracker/pkg/core"
	"github.com/laurendc/nannytracker/pkg/core/maps"
	"github.com/laurendc/nannytracker/pkg/core/storage"
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
		{Date: "2024-03-17", Origin: "Home", Destination: "Work", Miles: 10.0, Type: "single"}, // Week 1
		{Date: "2024-03-18", Origin: "Work", Destination: "Home", Miles: 15.0, Type: "round"},  // Week 1
		{Date: "2024-03-24", Origin: "Home", Destination: "Work", Miles: 20.0, Type: "single"}, // Week 2
		{Date: "2024-03-25", Origin: "Work", Destination: "Home", Miles: 25.0, Type: "round"},  // Week 2
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

	// Check the most recent week (default selected)
	view := uiModel.View()

	// Verify the extra newline between Trips and Expenses sections
	tripsSection := " Trips:"
	expensesSection := "\n Expenses:"
	tripsIndex := strings.Index(view, tripsSection)
	expensesIndex := strings.Index(view, expensesSection)
	if tripsIndex == -1 || expensesIndex == -1 {
		t.Fatalf("Could not find Trips or Expenses sections in view.\nView output:\n%s", view)
	}
	expensesIndex += 1 // skip the leading '\n' for correct slicing
	tripsSectionEnd := tripsIndex + len(tripsSection)
	if tripsSectionEnd > expensesIndex {
		t.Fatalf("Trips section end (%d) is after Expenses section start (%d).\nFull view output:\n%s", tripsSectionEnd, expensesIndex, view)
	}
	gapSubstring := view[tripsSectionEnd:expensesIndex]
	newlineCount := strings.Count(gapSubstring, "\n")
	if newlineCount < 2 {
		t.Errorf("Expected at least 2 newlines between Trips and Expenses sections, got %d", newlineCount)
	}

	expectedSummaries := []string{
		"Week of 2024-03-24 to 2024-03-30 (Week 1 of 2):",
		"    Total Miles:          70.00",
		"    Total Mileage Amount: $45.85",
		"    Total Expenses:       $30.00",
	}
	for _, expected := range expectedSummaries {
		if !strings.Contains(view, expected) {
			t.Errorf("View does not contain expected weekly summary: %s", expected)
		}
	}
	expectedTrips := []string{
		"2024-03-25: Work → Home (50.00 miles) [round]",
		"2024-03-24: Home → Work (20.00 miles) [single]",
	}
	for i, expected := range expectedTrips {
		if !strings.Contains(view, expected) {
			t.Errorf("View does not contain expected trip: %s", expected)
		}
		if i > 0 {
			prevTrip := expectedTrips[i-1]
			prevIndex := strings.Index(view, prevTrip)
			currentIndex := strings.Index(view, expected)
			if prevIndex > currentIndex {
				t.Errorf("Trips not in descending order: %s appears before %s", expected, prevTrip)
			}
		}
	}
	expectedExpenses := []string{
		"2024-03-24: $30.00 - Activities",
	}
	for _, expected := range expectedExpenses {
		if !strings.Contains(view, expected) {
			t.Errorf("View does not contain expected expense: %s", expected)
		}
	}

	// Move to the second week (older week)
	updatedModel, _ := uiModel.Update(tea.KeyMsg{Type: tea.KeyRight})
	uiModel = updatedModel.(*Model)
	view = uiModel.View()
	expectedSummaries = []string{
		"Week of 2024-03-17 to 2024-03-23 (Week 2 of 2):",
		"    Total Miles:          40.00",
		"    Total Mileage Amount: $26.20",
		"    Total Expenses:       $41.25",
	}
	for _, expected := range expectedSummaries {
		if !strings.Contains(view, expected) {
			t.Errorf("View does not contain expected weekly summary: %s", expected)
		}
	}
	expectedTrips = []string{
		"2024-03-18: Work → Home (30.00 miles) [round]",
		"2024-03-17: Home → Work (10.00 miles) [single]",
	}
	for i, expected := range expectedTrips {
		if !strings.Contains(view, expected) {
			t.Errorf("View does not contain expected trip: %s", expected)
		}
		if i > 0 {
			prevTrip := expectedTrips[i-1]
			prevIndex := strings.Index(view, prevTrip)
			currentIndex := strings.Index(view, expected)
			if prevIndex > currentIndex {
				t.Errorf("Trips not in descending order: %s appears before %s", expected, prevTrip)
			}
		}
	}
	expectedExpenses = []string{
		"2024-03-18: $15.75 - Snacks",
		"2024-03-17: $25.50 - Lunch",
	}
	for i, expected := range expectedExpenses {
		if !strings.Contains(view, expected) {
			t.Errorf("View does not contain expected expense: %s", expected)
		}
		if i > 0 {
			prevExpense := expectedExpenses[i-1]
			prevIndex := strings.Index(view, prevExpense)
			currentIndex := strings.Index(view, expected)
			if prevIndex > currentIndex {
				t.Errorf("Expenses not in descending order: %s appears before %s", expected, prevExpense)
			}
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
	uiModel.ActiveTab = TabTrips // Ensure we are on the Trips tab

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

	// Set active tab to Trips (required for delete to work)
	uiModel.ActiveTab = TabTrips

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

func TestTripHistoryDisplay(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Add test trips
	trip1 := model.Trip{
		Date:        "2024-03-20",
		Origin:      "Home",
		Destination: "Work",
		Miles:       10.0,
		Type:        "single",
	}
	trip2 := model.Trip{
		Date:        "2024-03-21",
		Origin:      "Work",
		Destination: "Store",
		Miles:       10.0,
		Type:        "single",
	}
	trip3 := model.Trip{
		Date:        "2024-03-22",
		Origin:      "Home",
		Destination: "Gym",
		Miles:       3.0,
		Type:        "single",
	}
	uiModel.AddTrip(trip1)
	uiModel.AddTrip(trip2)
	uiModel.AddTrip(trip3)

	// Add more trips to trigger pagination
	validDates := []string{
		"2024-03-23", "2024-03-24", "2024-03-25", "2024-03-26", "2024-03-27", "2024-03-28", "2024-03-29", "2024-03-30", "2024-03-31",
		"2024-04-01", "2024-04-02", "2024-04-03", "2024-04-04", "2024-04-05", "2024-04-06",
	}
	for _, date := range validDates {
		trip := model.Trip{
			Date:        date,
			Origin:      "Home",
			Destination: "Work",
			Miles:       5.0,
			Type:        "single",
		}
		uiModel.AddTrip(trip)
	}

	// Set active tab to Trips
	uiModel.ActiveTab = TabTrips

	// Navigate to page 2 (where the original 3 trips will appear)
	uiModel.CurrentPage = 1
	view := uiModel.View()
	expectedTrips := []string{
		"2024-03-22: Home → Gym (3.00 miles) [single]",
		"2024-03-21: Work → Store (10.00 miles) [single]",
		"2024-03-20: Home → Work (10.00 miles) [single]",
	}
	for _, trip := range expectedTrips {
		if !strings.Contains(view, trip) {
			t.Errorf("View does not contain expected trip: %s", trip)
		}
	}

	// Check if pagination info is displayed
	if !strings.Contains(view, "Page 2 of") {
		t.Error("View does not contain second page information")
	}

	// Test going back to first page
	updatedModel, _ := uiModel.Update(tea.KeyMsg{Type: tea.KeyLeft})
	uiModel = updatedModel.(*Model)
	if uiModel.CurrentPage != 0 {
		t.Errorf("Expected current page to be 0, got %d", uiModel.CurrentPage)
	}
}

func TestExpenseHistoryDisplay(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Add test expenses
	expense1 := model.Expense{
		Date:        "2024-03-20",
		Amount:      25.50,
		Description: "Lunch",
	}
	expense2 := model.Expense{
		Date:        "2024-03-21",
		Amount:      15.75,
		Description: "Snacks",
	}
	expense3 := model.Expense{
		Date:        "2024-03-22",
		Amount:      30.00,
		Description: "Activities",
	}
	if err := uiModel.Data.AddExpense(expense1); err != nil {
		t.Fatalf("Failed to add expense1: %v", err)
	}
	if err := uiModel.Data.AddExpense(expense2); err != nil {
		t.Fatalf("Failed to add expense2: %v", err)
	}
	if err := uiModel.Data.AddExpense(expense3); err != nil {
		t.Fatalf("Failed to add expense3: %v", err)
	}

	// Set active tab to Expenses
	uiModel.ActiveTab = TabExpenses

	// Check if expenses are displayed
	view := uiModel.View()
	expectedExpenses := []string{
		"$25.50 - Lunch",
		"$15.75 - Snacks",
		"$30.00 - Activities",
	}
	for _, expense := range expectedExpenses {
		if !strings.Contains(view, expense) {
			t.Errorf("View does not contain expected expense: %s", expense)
		}
	}
}

func TestTimeBasedGrouping(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Add test trips
	trip1 := model.Trip{
		Date:        "2024-03-20",
		Origin:      "Home",
		Destination: "Work",
		Miles:       10.0,
		Type:        "single",
	}
	trip2 := model.Trip{
		Date:        "2024-03-20",
		Origin:      "Work",
		Destination: "Home",
		Miles:       10.0,
		Type:        "single",
	}
	trip3 := model.Trip{
		Date:        "2024-03-21",
		Origin:      "Home",
		Destination: "Store",
		Miles:       5.0,
		Type:        "single",
	}
	trip4 := model.Trip{
		Date:        "2024-03-22",
		Origin:      "Home",
		Destination: "Gym",
		Miles:       3.0,
		Type:        "single",
	}
	uiModel.AddTrip(trip1)
	uiModel.AddTrip(trip2)
	uiModel.AddTrip(trip3)
	uiModel.AddTrip(trip4)

	// Set active tab to Trips
	uiModel.ActiveTab = TabTrips

	// Check if trips are displayed
	view := uiModel.View()
	expectedTrips := []string{
		"Home → Work (10.00 miles)",
		"Work → Home (10.00 miles)",
		"Home → Store (5.00 miles)",
		"Home → Gym (3.00 miles)",
	}
	for _, trip := range expectedTrips {
		if !strings.Contains(view, trip) {
			t.Errorf("View does not contain expected trip: %s", trip)
		}
	}
}

func TestTimeGroupNavigation(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Add test trips
	trip1 := model.Trip{
		Date:        "2024-03-20",
		Origin:      "Home",
		Destination: "Work",
		Miles:       10.0,
		Type:        "single",
	}
	trip2 := model.Trip{
		Date:        "2024-03-21",
		Origin:      "Work",
		Destination: "Store",
		Miles:       10.0,
		Type:        "single",
	}
	trip3 := model.Trip{
		Date:        "2024-03-22",
		Origin:      "Home",
		Destination: "Gym",
		Miles:       3.0,
		Type:        "single",
	}
	uiModel.AddTrip(trip1)
	uiModel.AddTrip(trip2)
	uiModel.AddTrip(trip3)

	// Set active tab to Trips
	uiModel.ActiveTab = TabTrips

	// Check if trips are displayed
	view := uiModel.View()
	expectedTrips := []string{
		"Home → Work (10.00 miles)",
		"Work → Store (10.00 miles)",
		"Home → Gym (3.00 miles)",
	}
	for _, trip := range expectedTrips {
		if !strings.Contains(view, trip) {
			t.Errorf("View does not contain expected trip: %s", trip)
		}
	}
}

func TestSearchFunctionality(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Add test trips
	trip1 := model.Trip{
		Date:        "2024-03-20",
		Origin:      "Home",
		Destination: "Work",
		Miles:       10.0,
		Type:        "single",
	}
	trip2 := model.Trip{
		Date:        "2024-03-21",
		Origin:      "Work",
		Destination: "Store",
		Miles:       10.0,
		Type:        "single",
	}
	trip3 := model.Trip{
		Date:        "2024-03-22",
		Origin:      "Home",
		Destination: "Gym",
		Miles:       3.0,
		Type:        "single",
	}
	uiModel.AddTrip(trip1)
	uiModel.AddTrip(trip2)
	uiModel.AddTrip(trip3)

	// Add more trips to trigger pagination for the initial search
	validDates := []string{
		"2024-03-23", "2024-03-24", "2024-03-25", "2024-03-26", "2024-03-27", "2024-03-28", "2024-03-29", "2024-03-30", "2024-03-31",
		"2024-04-01", "2024-04-02", "2024-04-03", "2024-04-04", "2024-04-05", "2024-04-06",
	}
	for _, date := range validDates {
		trip := model.Trip{
			Date:        date,
			Origin:      "Work",
			Destination: "Home",
			Miles:       5.0,
			Type:        "single",
		}
		uiModel.AddTrip(trip)
	}

	// Set active tab to Trips
	uiModel.ActiveTab = TabTrips

	// Enter search mode
	updatedModel, _ := uiModel.Update(tea.KeyMsg{Type: tea.KeyCtrlF})
	uiModel = updatedModel.(*Model)

	// Search for "Work"
	uiModel.TextInput.SetValue("Work")
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	// Navigate to the last page (page 4) by pressing the right arrow key multiple times
	for i := 0; i < 3; i++ {
		updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyRight})
		uiModel = updatedModel.(*Model)
	}
	view := uiModel.View()
	expectedTrips := []string{
		"2024-03-21: Work → Store (10.00 miles) [single]",
		"2024-03-20: Home → Work (10.00 miles) [single]",
	}
	for _, trip := range expectedTrips {
		if !strings.Contains(view, trip) {
			t.Errorf("View does not contain trip: %s", trip)
		}
	}

	// Check if pagination info is displayed
	if !strings.Contains(view, "Page 2 of 2") {
		t.Error("View does not contain pagination information")
	}
}

func TestExpenseNavigation(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Add test expenses
	expense1 := model.Expense{
		Date:        "2024-03-20",
		Amount:      25.50,
		Description: "Lunch",
	}
	expense2 := model.Expense{
		Date:        "2024-03-21",
		Amount:      15.75,
		Description: "Snacks",
	}
	expense3 := model.Expense{
		Date:        "2024-03-22",
		Amount:      30.00,
		Description: "Activities",
	}
	if err := uiModel.Data.AddExpense(expense1); err != nil {
		t.Fatalf("Failed to add expense1: %v", err)
	}
	if err := uiModel.Data.AddExpense(expense2); err != nil {
		t.Fatalf("Failed to add expense2: %v", err)
	}
	if err := uiModel.Data.AddExpense(expense3); err != nil {
		t.Fatalf("Failed to add expense3: %v", err)
	}

	// Set active tab to Expenses
	uiModel.ActiveTab = TabExpenses

	// Test initial selection
	updatedModel, _ := uiModel.Update(tea.KeyMsg{Type: tea.KeyDown})
	uiModel = updatedModel.(*Model)
	if uiModel.SelectedExpense != 0 {
		t.Errorf("Expected to start at expense 0 after Tab, got %d", uiModel.SelectedExpense)
	}
	if uiModel.SelectedTrip != -1 {
		t.Errorf("Expected trip selection to be -1 after Tab, got %d", uiModel.SelectedTrip)
	}

	// Test moving selection down
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyDown})
	uiModel = updatedModel.(*Model)
	if uiModel.SelectedExpense != 1 {
		t.Errorf("Expected selected expense to be 1, got %d", uiModel.SelectedExpense)
	}

	// Test moving selection down again
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyDown})
	uiModel = updatedModel.(*Model)
	if uiModel.SelectedExpense != 2 {
		t.Errorf("Expected selected expense to be 2, got %d", uiModel.SelectedExpense)
	}

	// Test moving selection up
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyUp})
	uiModel = updatedModel.(*Model)
	if uiModel.SelectedExpense != 1 {
		t.Errorf("Expected selected expense to be 1, got %d", uiModel.SelectedExpense)
	}

	// Test wrap-around when moving down
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyDown})
	uiModel = updatedModel.(*Model)
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyDown})
	uiModel = updatedModel.(*Model)
	if uiModel.SelectedExpense != 0 {
		t.Errorf("Expected selected expense to wrap to 0, got %d", uiModel.SelectedExpense)
	}

	// Test wrap-around when moving up
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyUp})
	uiModel = updatedModel.(*Model)
	if uiModel.SelectedExpense != 2 {
		t.Errorf("Expected selected expense to wrap to 2, got %d", uiModel.SelectedExpense)
	}
}

func TestConvertTripToRecurring(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Set reference date for testing
	uiModel.Data.ReferenceDate = "2024-03-20"

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

	// Set active tab to Trips (required for convert to recurring to work)
	uiModel.ActiveTab = TabTrips

	// Enter convert mode
	var updatedModel tea.Model
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyCtrlR})
	uiModel = updatedModel.(*Model)

	// Verify we're in convert mode
	if uiModel.Mode != "convert_to_recurring" {
		t.Errorf("Expected mode to be 'convert_to_recurring', got '%s'", uiModel.Mode)
	}

	// Verify the trip data was copied to the recurring trip
	if uiModel.CurrentRecurring.Origin != originalTrip.Origin {
		t.Errorf("Expected origin to be '%s', got '%s'", originalTrip.Origin, uiModel.CurrentRecurring.Origin)
	}
	if uiModel.CurrentRecurring.Destination != originalTrip.Destination {
		t.Errorf("Expected destination to be '%s', got '%s'", originalTrip.Destination, uiModel.CurrentRecurring.Destination)
	}
	if uiModel.CurrentRecurring.Miles != originalTrip.Miles {
		t.Errorf("Expected miles to be %.2f, got %.2f", originalTrip.Miles, uiModel.CurrentRecurring.Miles)
	}
	if uiModel.CurrentRecurring.StartDate != originalTrip.Date {
		t.Errorf("Expected start date to be '%s', got '%s'", originalTrip.Date, uiModel.CurrentRecurring.StartDate)
	}
	if uiModel.CurrentRecurring.Type != originalTrip.Type {
		t.Errorf("Expected type to be '%s', got '%s'", originalTrip.Type, uiModel.CurrentRecurring.Type)
	}

	// Test invalid weekday
	uiModel.TextInput.SetValue("7") // Invalid weekday
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	if uiModel.Err == nil {
		t.Error("Expected error for invalid weekday")
	}
	if !strings.Contains(uiModel.Err.Error(), "invalid weekday") {
		t.Errorf("Expected error about invalid weekday, got: %v", uiModel.Err)
	}

	// Set end date to end of March 2024 before setting the weekday
	uiModel.CurrentRecurring.EndDate = "2024-03-31"

	// Test valid weekday (Wednesday is 3)
	uiModel.TextInput.SetValue("3")
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	// Verify the recurring trip was added
	if len(uiModel.RecurringTrips) != 1 {
		t.Errorf("Expected 1 recurring trip, got %d", len(uiModel.RecurringTrips))
	}

	recurringTrip := uiModel.RecurringTrips[0]
	if recurringTrip.Origin != originalTrip.Origin {
		t.Errorf("Expected origin to be '%s', got '%s'", originalTrip.Origin, recurringTrip.Origin)
	}
	if recurringTrip.Destination != originalTrip.Destination {
		t.Errorf("Expected destination to be '%s', got '%s'", originalTrip.Destination, recurringTrip.Destination)
	}
	if recurringTrip.Miles != originalTrip.Miles {
		t.Errorf("Expected miles to be %.2f, got %.2f", originalTrip.Miles, recurringTrip.Miles)
	}
	if recurringTrip.StartDate != originalTrip.Date {
		t.Errorf("Expected start date to be '%s', got '%s'", originalTrip.Date, recurringTrip.StartDate)
	}
	if recurringTrip.Type != originalTrip.Type {
		t.Errorf("Expected type to be '%s', got '%s'", originalTrip.Type, recurringTrip.Type)
	}
	if recurringTrip.Weekday != 3 {
		t.Errorf("Expected weekday to be 3, got %d", recurringTrip.Weekday)
	}
	if recurringTrip.EndDate != "2024-03-31" {
		t.Errorf("Expected end date to be '2024-03-31', got '%s'", recurringTrip.EndDate)
	}

	// Verify mode was reset
	if uiModel.Mode != "date" {
		t.Errorf("Expected mode to reset to 'date', got '%s'", uiModel.Mode)
	}

	// Verify trips were generated (should be 2 Wednesdays: March 20 and March 27)
	if len(uiModel.Trips) != 2 {
		t.Errorf("Expected 2 trips to be generated, got %d", len(uiModel.Trips))
	}

	// Verify the generated trips are on the correct dates
	dates := make(map[string]bool)
	for _, trip := range uiModel.Trips {
		dates[trip.Date] = true
	}
	if !dates["2024-03-20"] {
		t.Error("Expected trip on March 20, 2024")
	}
	if !dates["2024-03-27"] {
		t.Error("Expected trip on March 27, 2024")
	}
}

func TestConvertTripToRecurringWithNoSelection(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Set active tab to Trips (required for convert to recurring to work)
	uiModel.ActiveTab = TabTrips

	// Try to convert without selecting a trip
	var updatedModel tea.Model
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyCtrlR})
	uiModel = updatedModel.(*Model)

	// Verify we're in recurring trip creation mode
	if uiModel.Mode != "recurring_date" {
		t.Errorf("Expected mode to be 'recurring_date', got '%s'", uiModel.Mode)
	}

	// Verify no recurring trip was created
	if len(uiModel.RecurringTrips) != 0 {
		t.Errorf("Expected 0 recurring trips, got %d", len(uiModel.RecurringTrips))
	}
}

func TestTabNavigation(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Test initial tab state
	if uiModel.ActiveTab != TabWeeklySummaries {
		t.Errorf("Expected initial tab to be Weekly Summaries (0), got %d", uiModel.ActiveTab)
	}

	// Test forward tab navigation
	updatedModel, _ := uiModel.Update(tea.KeyMsg{Type: tea.KeyTab})
	uiModel = updatedModel.(*Model)
	if uiModel.ActiveTab != TabTrips {
		t.Errorf("Expected tab to be Trips (1) after Tab, got %d", uiModel.ActiveTab)
	}

	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyTab})
	uiModel = updatedModel.(*Model)
	if uiModel.ActiveTab != TabExpenses {
		t.Errorf("Expected tab to be Expenses (2) after Tab, got %d", uiModel.ActiveTab)
	}

	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyTab})
	uiModel = updatedModel.(*Model)
	if uiModel.ActiveTab != TabTemplates {
		t.Errorf("Expected tab to be Templates (3) after Tab, got %d", uiModel.ActiveTab)
	}

	// Test wrap-around
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyTab})
	uiModel = updatedModel.(*Model)
	if uiModel.ActiveTab != TabWeeklySummaries {
		t.Errorf("Expected tab to wrap around to Weekly Summaries (0), got %d", uiModel.ActiveTab)
	}

	// Test reverse tab navigation
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	uiModel = updatedModel.(*Model)
	if uiModel.ActiveTab != TabTemplates {
		t.Errorf("Expected tab to be Templates (3) after Shift+Tab, got %d", uiModel.ActiveTab)
	}

	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	uiModel = updatedModel.(*Model)
	if uiModel.ActiveTab != TabExpenses {
		t.Errorf("Expected tab to be Expenses (2) after Shift+Tab, got %d", uiModel.ActiveTab)
	}

	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	uiModel = updatedModel.(*Model)
	if uiModel.ActiveTab != TabTrips {
		t.Errorf("Expected tab to be Trips (1) after Shift+Tab, got %d", uiModel.ActiveTab)
	}

	// Test wrap-around in reverse
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	uiModel = updatedModel.(*Model)
	if uiModel.ActiveTab != TabWeeklySummaries {
		t.Errorf("Expected tab to wrap around to Weekly Summaries (0) in reverse, got %d", uiModel.ActiveTab)
	}
}

func TestTabContentNavigation(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Add some test data
	trip := model.Trip{
		Date:        "2024-03-20",
		Origin:      "Home",
		Destination: "Work",
		Miles:       5.0,
		Type:        "single",
	}
	uiModel.AddTrip(trip)

	expense := model.Expense{
		Date:        "2024-03-20",
		Amount:      10.0,
		Description: "Test expense",
	}
	if err := uiModel.Data.AddExpense(expense); err != nil {
		t.Fatalf("Failed to add expense: %v", err)
	}

	// Test navigation in Trips tab
	uiModel.ActiveTab = TabTrips
	updatedModel, _ := uiModel.Update(tea.KeyMsg{Type: tea.KeyDown})
	uiModel = updatedModel.(*Model)
	if uiModel.SelectedTrip != 0 {
		t.Errorf("Expected selected trip to be 0, got %d", uiModel.SelectedTrip)
	}
	if uiModel.SelectedExpense != -1 {
		t.Errorf("Expected selected expense to be -1, got %d", uiModel.SelectedExpense)
	}

	// Test navigation in Expenses tab
	uiModel.ActiveTab = TabExpenses
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyDown})
	uiModel = updatedModel.(*Model)
	if uiModel.SelectedExpense != 0 {
		t.Errorf("Expected selected expense to be 0, got %d", uiModel.SelectedExpense)
	}
	if uiModel.SelectedTrip != -1 {
		t.Errorf("Expected selected trip to be -1, got %d", uiModel.SelectedTrip)
	}

	// Test navigation in Weekly Summaries tab (should not affect selections)
	uiModel.ActiveTab = TabWeeklySummaries
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyDown})
	uiModel = updatedModel.(*Model)
	if uiModel.SelectedExpense != 0 {
		t.Errorf("Expected selected expense to remain 0, got %d", uiModel.SelectedExpense)
	}
	if uiModel.SelectedTrip != -1 {
		t.Errorf("Expected selected trip to remain -1, got %d", uiModel.SelectedTrip)
	}
}

func TestTabContentDisplay(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Add test data
	trip := model.Trip{
		Date:        "2024-03-20",
		Origin:      "Home",
		Destination: "Work",
		Miles:       5.0,
		Type:        "single",
	}
	uiModel.AddTrip(trip)

	expense := model.Expense{
		Date:        "2024-03-20",
		Amount:      10.0,
		Description: "Test expense",
	}
	if err := uiModel.Data.AddExpense(expense); err != nil {
		t.Fatalf("Failed to add expense: %v", err)
	}

	// Test Weekly Summaries tab content
	uiModel.ActiveTab = TabWeeklySummaries
	view := uiModel.View()
	if !strings.Contains(view, "Weekly Summaries") {
		t.Error("Weekly Summaries tab should display weekly summaries content")
	}

	// Test Trips tab content
	uiModel.ActiveTab = TabTrips
	view = uiModel.View()
	if !strings.Contains(view, "Regular Trips") {
		t.Error("Trips tab should display trips content")
	}

	// Test Expenses tab content
	uiModel.ActiveTab = TabExpenses
	view = uiModel.View()
	if !strings.Contains(view, "Test expense") {
		t.Error("Expenses tab should display expenses content")
	}
}

func TestTripSelection(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Add test trips
	trip1 := model.Trip{
		Date:        "2024-03-20",
		Origin:      "Home",
		Destination: "Work",
		Miles:       5.0,
		Type:        "single",
	}
	trip2 := model.Trip{
		Date:        "2024-03-21",
		Origin:      "Work",
		Destination: "Store",
		Miles:       3.0,
		Type:        "single",
	}
	uiModel.AddTrip(trip1)
	uiModel.AddTrip(trip2)

	// Set active tab to Trips
	uiModel.ActiveTab = TabTrips

	// Test initial selection
	updatedModel, _ := uiModel.Update(tea.KeyMsg{Type: tea.KeyDown})
	uiModel = updatedModel.(*Model)
	if uiModel.SelectedTrip != 0 {
		t.Errorf("Expected selected trip to be 0, got %d", uiModel.SelectedTrip)
	}

	// Test moving selection down
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyDown})
	uiModel = updatedModel.(*Model)
	if uiModel.SelectedTrip != 1 {
		t.Errorf("Expected selected trip to be 1, got %d", uiModel.SelectedTrip)
	}

	// Test moving selection up
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyUp})
	uiModel = updatedModel.(*Model)
	if uiModel.SelectedTrip != 0 {
		t.Errorf("Expected selected trip to be 0, got %d", uiModel.SelectedTrip)
	}

	// Test page navigation
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyRight})
	uiModel = updatedModel.(*Model)
	if uiModel.CurrentPage != 0 {
		t.Errorf("Expected current page to be 0 (no change), got %d", uiModel.CurrentPage)
	}

	// Add more trips to test pagination
	for i := 0; i < 15; i++ {
		trip := model.Trip{
			Date:        fmt.Sprintf("2024-03-%02d", i+23),
			Origin:      "Home",
			Destination: "Work",
			Miles:       5.0,
			Type:        "single",
		}
		uiModel.AddTrip(trip)
	}

	// Test page navigation with multiple pages
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyRight})
	uiModel = updatedModel.(*Model)
	if uiModel.CurrentPage != 1 {
		t.Errorf("Expected current page to be 1, got %d", uiModel.CurrentPage)
	}

	// Test going back to first page
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyLeft})
	uiModel = updatedModel.(*Model)
	if uiModel.CurrentPage != 0 {
		t.Errorf("Expected current page to be 0, got %d", uiModel.CurrentPage)
	}
}

func TestExpenseDisplay(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Add test expenses in random order
	expenses := []model.Expense{
		{Date: "2024-03-20", Amount: 25.50, Description: "Lunch"},
		{Date: "2024-03-22", Amount: 30.00, Description: "Activities"},
		{Date: "2024-03-21", Amount: 15.75, Description: "Snacks"},
	}

	for _, expense := range expenses {
		if err := uiModel.Data.AddExpense(expense); err != nil {
			t.Fatalf("Failed to add expense: %v", err)
		}
	}

	// Set active tab to Expenses
	uiModel.ActiveTab = TabExpenses

	// Get the view
	view := uiModel.View()

	// Check if expenses are displayed in descending order
	expectedExpenses := []string{
		"2024-03-22: $30.00 - Activities",
		"2024-03-21: $15.75 - Snacks",
		"2024-03-20: $25.50 - Lunch",
	}

	for i, expected := range expectedExpenses {
		if !strings.Contains(view, expected) {
			t.Errorf("View does not contain expected expense: %s", expected)
		}
		// Verify order
		if i > 0 {
			prevExpense := expectedExpenses[i-1]
			prevIndex := strings.Index(view, prevExpense)
			currentIndex := strings.Index(view, expected)
			if prevIndex > currentIndex {
				t.Errorf("Expenses not in descending order: %s appears before %s", expected, prevExpense)
			}
		}
	}
}

func TestExpensePagination(t *testing.T) {
	m := &Model{
		Data: &model.StorageData{
			Expenses: []model.Expense{
				{Date: "2024-03-27", Amount: 10.00, Description: "Expense 1"},
				{Date: "2024-03-26", Amount: 20.00, Description: "Expense 2"},
				{Date: "2024-03-25", Amount: 30.00, Description: "Expense 3"},
				{Date: "2024-03-24", Amount: 40.00, Description: "Expense 4"},
				{Date: "2024-03-23", Amount: 50.00, Description: "Expense 5"},
				{Date: "2024-03-22", Amount: 60.00, Description: "Expense 6"},
				{Date: "2024-03-21", Amount: 70.00, Description: "Expense 7"},
				{Date: "2024-03-20", Amount: 80.00, Description: "Expense 8"},
				{Date: "2024-03-19", Amount: 90.00, Description: "Expense 9"},
				{Date: "2024-03-18", Amount: 100.00, Description: "Expense 10"},
				{Date: "2024-03-17", Amount: 110.00, Description: "Expense 11"},
			},
		},
		PageSize:    5,
		CurrentPage: 0,
		ActiveTab:   TabExpenses,
	}

	// Test first page
	view := m.View()
	if !strings.Contains(view, "Page 1 of 3 (Showing 1-5 of 11 expenses)") {
		t.Errorf("View does not contain pagination information")
	}

	// Test second page
	m.CurrentPage = 1
	view = m.View()
	if !strings.Contains(view, "Page 2 of 3 (Showing 6-10 of 11 expenses)") {
		t.Errorf("View does not contain second page information")
	}

	// Test third page
	m.CurrentPage = 2
	view = m.View()
	if !strings.Contains(view, "Page 3 of 3 (Showing 11-11 of 11 expenses)") {
		t.Errorf("View does not contain third page information")
	}

	// Test expense sorting
	if !strings.Contains(view, "2024-03-17: $110.00 - Expense 11") {
		t.Errorf("Expenses are not sorted by date in descending order")
	}
}

func TestWeeklySummarySorting(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Add trips for different weeks
	trips := []model.Trip{
		{Date: "2024-03-17", Origin: "Home", Destination: "Work", Miles: 10.0, Type: "single"}, // Week 1
		{Date: "2024-03-18", Origin: "Work", Destination: "Home", Miles: 15.0, Type: "round"},  // Week 1
		{Date: "2024-03-24", Origin: "Home", Destination: "Work", Miles: 20.0, Type: "single"}, // Week 2
		{Date: "2024-03-25", Origin: "Work", Destination: "Home", Miles: 25.0, Type: "round"},  // Week 2
	}

	// Add expenses for different weeks
	expenses := []model.Expense{
		{Date: "2024-03-17", Amount: 25.50, Description: "Lunch"},      // Week 1
		{Date: "2024-03-18", Amount: 15.75, Description: "Snacks"},     // Week 1
		{Date: "2024-03-24", Amount: 30.00, Description: "Activities"}, // Week 2
		{Date: "2024-03-25", Amount: 35.25, Description: "Materials"},  // Week 2
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

	// Set active tab to Weekly Summaries
	uiModel.ActiveTab = TabWeeklySummaries

	// Test first week's sorting (most recent week)
	view := uiModel.View()
	expectedTrips := []string{
		"2024-03-25: Work → Home (50.00 miles) [round]",
		"2024-03-24: Home → Work (20.00 miles) [single]",
	}

	for i, expected := range expectedTrips {
		if !strings.Contains(view, expected) {
			t.Errorf("View does not contain expected trip: %s", expected)
		}

		if i > 0 {
			prevTrip := expectedTrips[i-1]
			prevIndex := strings.Index(view, prevTrip)
			currentIndex := strings.Index(view, expected)
			if prevIndex > currentIndex {
				t.Errorf("Trips not in descending order: %s appears before %s", expected, prevTrip)
			}
		}
	}

	expectedExpenses := []string{
		"2024-03-25: $35.25 - Materials",
		"2024-03-24: $30.00 - Activities",
	}

	for i, expected := range expectedExpenses {
		if !strings.Contains(view, expected) {
			t.Errorf("View does not contain expected expense: %s", expected)
		}

		if i > 0 {
			prevExpense := expectedExpenses[i-1]
			prevIndex := strings.Index(view, prevExpense)
			currentIndex := strings.Index(view, expected)
			if prevIndex > currentIndex {
				t.Errorf("Expenses not in descending order: %s appears before %s", expected, prevExpense)
			}
		}
	}

	// Test second week's sorting (older week)
	updatedModel, _ := uiModel.Update(tea.KeyMsg{Type: tea.KeyRight})
	uiModel = updatedModel.(*Model)
	view = uiModel.View()

	expectedTrips = []string{
		"2024-03-18: Work → Home (30.00 miles) [round]",
		"2024-03-17: Home → Work (10.00 miles) [single]",
	}

	for i, expected := range expectedTrips {
		if !strings.Contains(view, expected) {
			t.Errorf("View does not contain expected trip: %s", expected)
		}

		if i > 0 {
			prevTrip := expectedTrips[i-1]
			prevIndex := strings.Index(view, prevTrip)
			currentIndex := strings.Index(view, expected)
			if prevIndex > currentIndex {
				t.Errorf("Trips not in descending order: %s appears before %s", expected, prevTrip)
			}
		}
	}

	expectedExpenses = []string{
		"2024-03-18: $15.75 - Snacks",
		"2024-03-17: $25.50 - Lunch",
	}

	for i, expected := range expectedExpenses {
		if !strings.Contains(view, expected) {
			t.Errorf("View does not contain expected expense: %s", expected)
		}

		if i > 0 {
			prevExpense := expectedExpenses[i-1]
			prevIndex := strings.Index(view, prevExpense)
			currentIndex := strings.Index(view, expected)
			if prevIndex > currentIndex {
				t.Errorf("Expenses not in descending order: %s appears before %s", expected, prevExpense)
			}
		}
	}
}

func TestTemplateCreation(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Enter template creation mode
	msg := tea.KeyMsg{Type: tea.KeyCtrlT}
	model, _ := uiModel.Update(msg)
	uiModel = model.(*Model)

	if uiModel.Mode != "template_name" {
		t.Errorf("Expected mode to be 'template_name', got '%s'", uiModel.Mode)
	}

	// Enter template name
	uiModel.TextInput.SetValue("Work Commute")
	model, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = model.(*Model)

	if uiModel.Mode != "template_origin" {
		t.Errorf("Expected mode to be 'template_origin', got '%s'", uiModel.Mode)
	}

	// Enter origin
	uiModel.TextInput.SetValue("123 Home St")
	model, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = model.(*Model)

	if uiModel.Mode != "template_destination" {
		t.Errorf("Expected mode to be 'template_destination', got '%s'", uiModel.Mode)
	}

	// Enter destination
	uiModel.TextInput.SetValue("456 Work Ave")
	model, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = model.(*Model)

	if uiModel.Mode != "template_type" {
		t.Errorf("Expected mode to be 'template_type', got '%s'", uiModel.Mode)
	}

	// Enter trip type
	uiModel.TextInput.SetValue("round")
	model, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = model.(*Model)

	if uiModel.Mode != "template_notes" {
		t.Errorf("Expected mode to be 'template_notes', got '%s'", uiModel.Mode)
	}

	// Enter notes
	uiModel.TextInput.SetValue("Regular work commute")
	model, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = model.(*Model)

	// Verify template was created
	if len(uiModel.TripTemplates) != 1 {
		t.Errorf("Expected 1 template, got %d", len(uiModel.TripTemplates))
	}

	template := uiModel.TripTemplates[0]
	if template.Name != "Work Commute" {
		t.Errorf("Expected name 'Work Commute', got '%s'", template.Name)
	}
	if template.Origin != "123 Home St" {
		t.Errorf("Expected origin '123 Home St', got '%s'", template.Origin)
	}
	if template.Destination != "456 Work Ave" {
		t.Errorf("Expected destination '456 Work Ave', got '%s'", template.Destination)
	}
	if template.TripType != "round" {
		t.Errorf("Expected type 'round', got '%s'", template.TripType)
	}
	if template.Notes != "Regular work commute" {
		t.Errorf("Expected notes 'Regular work commute', got '%s'", template.Notes)
	}
}

func TestTemplateEditing(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Add a template first
	template := model.TripTemplate{
		Name:        "Work Commute",
		Origin:      "123 Home St",
		Destination: "456 Work Ave",
		TripType:    "single",
		Notes:       "Regular work commute",
	}
	uiModel.TripTemplates = append(uiModel.TripTemplates, template)
	uiModel.Data.TripTemplates = uiModel.TripTemplates

	// Select the template
	uiModel.SelectedTemplate = 0

	// Set active tab to Templates (required for edit to work)
	uiModel.ActiveTab = TabTemplates

	// Enter edit mode
	msg := tea.KeyMsg{Type: tea.KeyCtrlE}
	model, _ := uiModel.Update(msg)
	uiModel = model.(*Model)

	if uiModel.Mode != "template_edit" {
		t.Errorf("Expected mode to be 'template_edit', got '%s'", uiModel.Mode)
	}

	// Edit name
	uiModel.TextInput.SetValue("Updated Work Commute")
	model, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = model.(*Model)

	// Skip origin (keep existing)
	model, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = model.(*Model)

	// Skip destination (keep existing)
	model, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = model.(*Model)

	// Edit trip type
	uiModel.TextInput.SetValue("round")
	model, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = model.(*Model)

	// Skip notes (keep existing)
	model, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = model.(*Model)

	// Verify template was updated
	if len(uiModel.TripTemplates) != 1 {
		t.Errorf("Expected 1 template, got %d", len(uiModel.TripTemplates))
	}

	updatedTemplate := uiModel.TripTemplates[0]
	if updatedTemplate.Name != "Updated Work Commute" {
		t.Errorf("Expected name 'Updated Work Commute', got '%s'", updatedTemplate.Name)
	}
	if updatedTemplate.Origin != "123 Home St" {
		t.Errorf("Expected origin '123 Home St', got '%s'", updatedTemplate.Origin)
	}
	if updatedTemplate.Destination != "456 Work Ave" {
		t.Errorf("Expected destination '456 Work Ave', got '%s'", updatedTemplate.Destination)
	}
	if updatedTemplate.TripType != "round" {
		t.Errorf("Expected type 'round', got '%s'", updatedTemplate.TripType)
	}
	if updatedTemplate.Notes != "Regular work commute" {
		t.Errorf("Expected notes 'Regular work commute', got '%s'", updatedTemplate.Notes)
	}
}

func TestTemplateDeletion(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Add a template first
	template := model.TripTemplate{
		Name:        "Work Commute",
		Origin:      "123 Home St",
		Destination: "456 Work Ave",
		TripType:    "single",
		Notes:       "Regular work commute",
	}
	uiModel.TripTemplates = append(uiModel.TripTemplates, template)
	uiModel.Data.TripTemplates = uiModel.TripTemplates

	// Select the template
	uiModel.SelectedTemplate = 0

	// Set active tab to Templates (required for delete to work)
	uiModel.ActiveTab = TabTemplates

	// Enter delete confirmation mode
	msg := tea.KeyMsg{Type: tea.KeyCtrlD}
	model, _ := uiModel.Update(msg)
	uiModel = model.(*Model)

	if uiModel.Mode != "template_delete_confirm" {
		t.Errorf("Expected mode to be 'template_delete_confirm', got '%s'", uiModel.Mode)
	}

	// Test cancellation by entering something other than 'yes'
	uiModel.TextInput.SetValue("no")
	model, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = model.(*Model)

	// Verify template wasn't deleted and mode was reset
	if len(uiModel.TripTemplates) != 1 {
		t.Errorf("Expected 1 template after cancellation, got %d", len(uiModel.TripTemplates))
	}
	if uiModel.Mode != "date" {
		t.Errorf("Expected mode to be 'date' after cancellation, got '%s'", uiModel.Mode)
	}

	// Enter delete confirmation mode again
	msg = tea.KeyMsg{Type: tea.KeyCtrlD}
	model, _ = uiModel.Update(msg)
	uiModel = model.(*Model)

	// Confirm deletion by entering 'yes'
	uiModel.TextInput.SetValue("yes")
	model, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = model.(*Model)

	// Verify template was deleted and mode was reset
	if len(uiModel.TripTemplates) != 0 {
		t.Errorf("Expected 0 templates after deletion, got %d", len(uiModel.TripTemplates))
	}
	if uiModel.Mode != "date" {
		t.Errorf("Expected mode to be 'date' after deletion, got '%s'", uiModel.Mode)
	}
}

func TestTabNavigationWithTemplates(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Add a test template
	template := model.TripTemplate{
		Name:        "Work Commute",
		Origin:      "123 Home St",
		Destination: "456 Work Ave",
		TripType:    "single",
		Notes:       "Regular work commute",
	}
	uiModel.TripTemplates = append(uiModel.TripTemplates, template)
	uiModel.Data.TripTemplates = uiModel.TripTemplates

	// Test initial tab state
	if uiModel.ActiveTab != TabWeeklySummaries {
		t.Errorf("Expected initial tab to be Weekly Summaries (0), got %d", uiModel.ActiveTab)
	}

	// Test forward tab navigation including Templates
	updatedModel, _ := uiModel.Update(tea.KeyMsg{Type: tea.KeyTab})
	uiModel = updatedModel.(*Model)
	if uiModel.ActiveTab != TabTrips {
		t.Errorf("Expected tab to be Trips (1) after Tab, got %d", uiModel.ActiveTab)
	}

	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyTab})
	uiModel = updatedModel.(*Model)
	if uiModel.ActiveTab != TabExpenses {
		t.Errorf("Expected tab to be Expenses (2) after Tab, got %d", uiModel.ActiveTab)
	}

	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyTab})
	uiModel = updatedModel.(*Model)
	if uiModel.ActiveTab != TabTemplates {
		t.Errorf("Expected tab to be Templates (3) after Tab, got %d", uiModel.ActiveTab)
	}

	// Test wrap-around
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyTab})
	uiModel = updatedModel.(*Model)
	if uiModel.ActiveTab != TabWeeklySummaries {
		t.Errorf("Expected tab to wrap around to Weekly Summaries (0), got %d", uiModel.ActiveTab)
	}

	// Test reverse tab navigation including Templates
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	uiModel = updatedModel.(*Model)
	if uiModel.ActiveTab != TabTemplates {
		t.Errorf("Expected tab to be Templates (3) after Shift+Tab, got %d", uiModel.ActiveTab)
	}

	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	uiModel = updatedModel.(*Model)
	if uiModel.ActiveTab != TabExpenses {
		t.Errorf("Expected tab to be Expenses (2) after Shift+Tab, got %d", uiModel.ActiveTab)
	}

	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	uiModel = updatedModel.(*Model)
	if uiModel.ActiveTab != TabTrips {
		t.Errorf("Expected tab to be Trips (1) after Shift+Tab, got %d", uiModel.ActiveTab)
	}

	// Test wrap-around in reverse
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	uiModel = updatedModel.(*Model)
	if uiModel.ActiveTab != TabWeeklySummaries {
		t.Errorf("Expected tab to wrap around to Weekly Summaries (0) in reverse, got %d", uiModel.ActiveTab)
	}
}

func TestTemplateTabContent(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Add test templates
	template1 := model.TripTemplate{
		Name:        "Work Commute",
		Origin:      "123 Home St",
		Destination: "456 Work Ave",
		TripType:    "single",
		Notes:       "Regular work commute",
	}
	template2 := model.TripTemplate{
		Name:        "Grocery Run",
		Origin:      "123 Home St",
		Destination: "789 Market St",
		TripType:    "round",
		Notes:       "Weekly grocery shopping",
	}
	uiModel.TripTemplates = append(uiModel.TripTemplates, template1, template2)
	uiModel.Data.TripTemplates = uiModel.TripTemplates

	// Set active tab to Templates
	uiModel.ActiveTab = TabTemplates

	// Test template selection - now navigates in sorted order (Grocery Run first, then Work Commute)
	updatedModel, _ := uiModel.Update(tea.KeyMsg{Type: tea.KeyDown})
	uiModel = updatedModel.(*Model)

	// Find which template is selected (should be Grocery Run alphabetically)
	selectedTemplate := uiModel.TripTemplates[uiModel.SelectedTemplate]
	if selectedTemplate.Name != "Grocery Run" {
		t.Errorf("Expected first selected template to be 'Grocery Run' (alphabetically first), got '%s'", selectedTemplate.Name)
	}

	// Test moving selection down (should go to Work Commute)
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyDown})
	uiModel = updatedModel.(*Model)
	selectedTemplate = uiModel.TripTemplates[uiModel.SelectedTemplate]
	if selectedTemplate.Name != "Work Commute" {
		t.Errorf("Expected second selected template to be 'Work Commute' (alphabetically second), got '%s'", selectedTemplate.Name)
	}

	// Test moving selection up (should go back to Grocery Run)
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyUp})
	uiModel = updatedModel.(*Model)
	selectedTemplate = uiModel.TripTemplates[uiModel.SelectedTemplate]
	if selectedTemplate.Name != "Grocery Run" {
		t.Errorf("Expected up arrow to go back to 'Grocery Run', got '%s'", selectedTemplate.Name)
	}

	// Test template content display
	view := uiModel.View()
	expectedContent := []string{
		"Work Commute: 123 Home St → 456 Work Ave [single] - Regular work commute",
		"Grocery Run: 123 Home St → 789 Market St [round] - Weekly grocery shopping",
	}
	for _, content := range expectedContent {
		if !strings.Contains(view, content) {
			t.Errorf("View does not contain expected template content: %s", content)
		}
	}
}

func TestTemplateToTripConversion(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Add a test template
	template := model.TripTemplate{
		Name:        "Work Commute",
		Origin:      "123 Home St",
		Destination: "456 Work Ave",
		TripType:    "single",
		Notes:       "Regular work commute",
	}
	uiModel.TripTemplates = append(uiModel.TripTemplates, template)
	uiModel.Data.TripTemplates = uiModel.TripTemplates

	// Select the template
	uiModel.SelectedTemplate = 0

	// Set active tab to Templates (required for template-to-trip conversion)
	uiModel.ActiveTab = TabTemplates

	// Trigger template-to-trip conversion with Ctrl+U
	msg := tea.KeyMsg{Type: tea.KeyCtrlU}
	model, _ := uiModel.Update(msg)
	uiModel = model.(*Model)

	// Verify we switched to Trips tab
	if uiModel.ActiveTab != TabTrips {
		t.Errorf("Expected active tab to be Trips, got %d", uiModel.ActiveTab)
	}

	// Verify we're in date input mode
	if uiModel.Mode != "date" {
		t.Errorf("Expected mode to be 'date', got '%s'", uiModel.Mode)
	}

	// Verify the current trip has template data
	if uiModel.CurrentTrip.Origin != template.Origin {
		t.Errorf("Expected origin to be '%s', got '%s'", template.Origin, uiModel.CurrentTrip.Origin)
	}
	if uiModel.CurrentTrip.Destination != template.Destination {
		t.Errorf("Expected destination to be '%s', got '%s'", template.Destination, uiModel.CurrentTrip.Destination)
	}
	if uiModel.CurrentTrip.Type != template.TripType {
		t.Errorf("Expected type to be '%s', got '%s'", template.TripType, uiModel.CurrentTrip.Type)
	}
	if uiModel.CurrentTrip.Miles != 0 {
		t.Errorf("Expected initial miles to be 0, got %.2f", uiModel.CurrentTrip.Miles)
	}

	// Enter date
	uiModel.TextInput.SetValue("2024-03-20")
	model, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = model.(*Model)

	// Enter origin (accept pre-filled value)
	model, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = model.(*Model)

	// Enter destination (accept pre-filled value)
	model, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = model.(*Model)

	// Enter type (accept pre-filled value)
	model, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = model.(*Model)

	// Verify the trip was created with calculated miles
	if len(uiModel.Trips) != 1 {
		t.Errorf("Expected 1 trip, got %d", len(uiModel.Trips))
	}

	trip := uiModel.Trips[0]
	if trip.Date != "2024-03-20" {
		t.Errorf("Expected date to be '2024-03-20', got '%s'", trip.Date)
	}
	if trip.Origin != template.Origin {
		t.Errorf("Expected origin to be '%s', got '%s'", template.Origin, trip.Origin)
	}
	if trip.Destination != template.Destination {
		t.Errorf("Expected destination to be '%s', got '%s'", template.Destination, trip.Destination)
	}
	if trip.Type != template.TripType {
		t.Errorf("Expected type to be '%s', got '%s'", template.TripType, trip.Type)
	}
	if trip.Miles != 10.0 { // Mock client returns 10.0 miles
		t.Errorf("Expected miles to be 10.0, got %.2f", trip.Miles)
	}
}

func TestTemplateSorting(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Add test templates in unsorted order
	templates := []model.TripTemplate{
		{
			Name:        "Zebra Trip",
			Origin:      "123 Home St",
			Destination: "456 Work Ave",
			TripType:    "single",
			Notes:       "Last alphabetically",
		},
		{
			Name:        "apple trip",
			Origin:      "123 Home St",
			Destination: "789 Market St",
			TripType:    "round",
			Notes:       "First alphabetically (case-insensitive)",
		},
		{
			Name:        "Banana Trip",
			Origin:      "123 Home St",
			Destination: "101 Park Ave",
			TripType:    "single",
			Notes:       "Middle alphabetically",
		},
	}
	uiModel.TripTemplates = templates
	uiModel.Data.TripTemplates = templates

	// Set active tab to Templates
	uiModel.ActiveTab = TabTemplates

	// Get the view
	view := uiModel.View()

	// Verify templates are displayed in alphabetical order
	expectedOrder := []string{
		"apple trip: 123 Home St → 789 Market St [round] - First alphabetically (case-insensitive)",
		"Banana Trip: 123 Home St → 101 Park Ave [single] - Middle alphabetically",
		"Zebra Trip: 123 Home St → 456 Work Ave [single] - Last alphabetically",
	}

	// Find the template section in the view
	templateSection := strings.Split(view, "Trip Templates:")[1]
	templateSection = strings.Split(templateSection, "Controls:")[0]

	// Verify each template appears in the correct order
	for i, expected := range expectedOrder {
		if !strings.Contains(templateSection, expected) {
			t.Errorf("Expected template '%s' not found in view", expected)
		}
		if i > 0 {
			// Verify order by checking that each template appears after the previous one
			prevTemplate := expectedOrder[i-1]
			prevIndex := strings.Index(templateSection, prevTemplate)
			currIndex := strings.Index(templateSection, expected)
			if prevIndex > currIndex {
				t.Errorf("Template '%s' appears before '%s' in view, but should be after",
					prevTemplate, expected)
			}
		}
	}
}

func TestWeeklySummaryUpdateAfterTripEdit(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Add a trip
	trip := model.Trip{
		Date:        "2024-03-20",
		Origin:      "Home",
		Destination: "Work",
		Miles:       10.0,
		Type:        "round",
	}
	uiModel.AddTrip(trip)

	// Switch to Trips tab
	uiModel.ActiveTab = TabTrips
	uiModel.SelectedTrip = 0

	// Enter edit mode
	updatedModel, _ := uiModel.Update(tea.KeyMsg{Type: tea.KeyCtrlE})
	uiModel = updatedModel.(*Model)

	// Follow the edit sequence
	// 1. Edit date (keep existing)
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	// 2. Edit origin (keep existing)
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	// 3. Edit destination (keep existing)
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	// 4. Edit trip type
	uiModel.TextInput.SetValue("single")
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	uiModel = updatedModel.(*Model)

	// Switch back to Weekly Summaries tab
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyTab})
	uiModel = updatedModel.(*Model)

	// Ensure SelectedWeek is valid
	uiModel.SelectedWeek = 0
	// Ensure the correct tab is active
	uiModel.ActiveTab = TabWeeklySummaries

	// Verify the weekly summary shows the updated trip type
	view := uiModel.View()
	expectedTrip := "2024-03-20: Home → Work (10.00 miles) [single]"
	if !strings.Contains(view, expectedTrip) {
		t.Errorf("Weekly summary does not show updated trip type. Expected to find: %s", expectedTrip)
		t.Logf("View output:\n%s", view)
		t.Logf("WeeklySummaries: %+v", uiModel.Data.WeeklySummaries)
		t.Logf("SelectedWeek: %d, WeeklySummariesLen: %d", uiModel.SelectedWeek, len(uiModel.Data.WeeklySummaries))
	}

	// Verify the total miles are updated (should be 10 instead of 20 for single trip)
	expectedMiles := "Total Miles:          10.00"
	if !strings.Contains(view, expectedMiles) {
		t.Errorf("Weekly summary does not show updated total miles. Expected to find: %s", expectedMiles)
		t.Logf("View output:\n%s", view)
		t.Logf("WeeklySummaries: %+v", uiModel.Data.WeeklySummaries)
		t.Logf("SelectedWeek: %d, WeeklySummariesLen: %d", uiModel.SelectedWeek, len(uiModel.Data.WeeklySummaries))
	}
}

func TestTemplateNavigationWithSorting(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Add test templates in non-alphabetical order
	template1 := model.TripTemplate{
		Name:        "Zebra Trip",
		Origin:      "Home",
		Destination: "Zoo",
		TripType:    "single",
		Notes:       "Visit zebras",
	}
	template2 := model.TripTemplate{
		Name:        "Apple Trip",
		Origin:      "Home",
		Destination: "Store",
		TripType:    "round",
		Notes:       "Buy apples",
	}
	template3 := model.TripTemplate{
		Name:        "Banana Trip",
		Origin:      "Home",
		Destination: "Market",
		TripType:    "single",
		Notes:       "Buy bananas",
	}

	// Add templates to the model
	if err := uiModel.Data.AddTripTemplate(template1); err != nil {
		t.Fatalf("Failed to add template1: %v", err)
	}
	if err := uiModel.Data.AddTripTemplate(template2); err != nil {
		t.Fatalf("Failed to add template2: %v", err)
	}
	if err := uiModel.Data.AddTripTemplate(template3); err != nil {
		t.Fatalf("Failed to add template3: %v", err)
	}

	// Update the model's template list
	uiModel.TripTemplates = uiModel.Data.TripTemplates

	// Set active tab to Templates
	uiModel.ActiveTab = TabTemplates

	// Test initial selection (should start with no selection)
	if uiModel.SelectedTemplate != -1 {
		t.Errorf("Expected initial template selection to be -1, got %d", uiModel.SelectedTemplate)
	}

	// Test first down arrow (should select first template in sorted order - Apple Trip)
	updatedModel, _ := uiModel.Update(tea.KeyMsg{Type: tea.KeyDown})
	uiModel = updatedModel.(*Model)

	// Verify that Apple Trip is selected (first alphabetically)
	if uiModel.SelectedTemplate == -1 {
		t.Error("Expected a template to be selected after first down arrow")
	} else {
		selectedTemplate := uiModel.TripTemplates[uiModel.SelectedTemplate]
		if selectedTemplate.Name != "Apple Trip" {
			t.Errorf("Expected Apple Trip to be selected first (alphabetically), got %s", selectedTemplate.Name)
		}
	}

	// Test second down arrow (should select Banana Trip)
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyDown})
	uiModel = updatedModel.(*Model)

	selectedTemplate := uiModel.TripTemplates[uiModel.SelectedTemplate]
	if selectedTemplate.Name != "Banana Trip" {
		t.Errorf("Expected Banana Trip to be selected second (alphabetically), got %s", selectedTemplate.Name)
	}

	// Test third down arrow (should select Zebra Trip)
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyDown})
	uiModel = updatedModel.(*Model)

	selectedTemplate = uiModel.TripTemplates[uiModel.SelectedTemplate]
	if selectedTemplate.Name != "Zebra Trip" {
		t.Errorf("Expected Zebra Trip to be selected third (alphabetically), got %s", selectedTemplate.Name)
	}

	// Test fourth down arrow (should wrap around to Apple Trip)
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyDown})
	uiModel = updatedModel.(*Model)

	selectedTemplate = uiModel.TripTemplates[uiModel.SelectedTemplate]
	if selectedTemplate.Name != "Apple Trip" {
		t.Errorf("Expected wrap-around to Apple Trip, got %s", selectedTemplate.Name)
	}

	// Test up arrow (should go back to Zebra Trip)
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyUp})
	uiModel = updatedModel.(*Model)

	selectedTemplate = uiModel.TripTemplates[uiModel.SelectedTemplate]
	if selectedTemplate.Name != "Zebra Trip" {
		t.Errorf("Expected up arrow to go to Zebra Trip, got %s", selectedTemplate.Name)
	}

	// Test up arrow again (should go to Banana Trip)
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyUp})
	uiModel = updatedModel.(*Model)

	selectedTemplate = uiModel.TripTemplates[uiModel.SelectedTemplate]
	if selectedTemplate.Name != "Banana Trip" {
		t.Errorf("Expected up arrow to go to Banana Trip, got %s", selectedTemplate.Name)
	}

	// Test up arrow again (should go to Apple Trip)
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyUp})
	uiModel = updatedModel.(*Model)

	selectedTemplate = uiModel.TripTemplates[uiModel.SelectedTemplate]
	if selectedTemplate.Name != "Apple Trip" {
		t.Errorf("Expected up arrow to go to Apple Trip, got %s", selectedTemplate.Name)
	}

	// Test up arrow again (should wrap around to Zebra Trip)
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyUp})
	uiModel = updatedModel.(*Model)

	selectedTemplate = uiModel.TripTemplates[uiModel.SelectedTemplate]
	if selectedTemplate.Name != "Zebra Trip" {
		t.Errorf("Expected wrap-around to Zebra Trip, got %s", selectedTemplate.Name)
	}
}

func TestRenderStatusBar(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Test Weekly Summary tab
	uiModel.ActiveTab = TabWeeklySummaries
	uiModel.Mode = "date"
	status := uiModel.renderStatusBar()
	if !strings.Contains(status, "📊 Weekly Summary") {
		t.Errorf("Expected status to contain '📊 Weekly Summary', got: %s", status)
	}

	// Test Trips tab with search
	uiModel.ActiveTab = TabTrips
	uiModel.SearchMode = true
	uiModel.SearchQuery = "grocery"
	status = uiModel.renderStatusBar()
	if !strings.Contains(status, "🚗 Trips") {
		t.Errorf("Expected status to contain '🚗 Trips', got: %s", status)
	}
	if !strings.Contains(status, "Search: \"grocery\"") {
		t.Errorf("Expected status to contain search query, got: %s", status)
	}
}

func TestRenderContextualControls(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Test Weekly Summary controls
	uiModel.ActiveTab = TabWeeklySummaries
	controls := uiModel.renderContextualControls()
	if !strings.Contains(controls, "NAVIGATION:") {
		t.Errorf("Expected controls to contain navigation section, got: %s", controls)
	}
	if !strings.Contains(controls, "ACTIONS:") {
		t.Errorf("Expected controls to contain actions section, got: %s", controls)
	}

	// Test Trips controls
	uiModel.ActiveTab = TabTrips
	controls = uiModel.renderContextualControls()
	if !strings.Contains(controls, "DELETE:") {
		t.Errorf("Expected controls to contain delete section, got: %s", controls)
	}
	if !strings.Contains(controls, "QUICK ADD:") {
		t.Errorf("Expected controls to contain quick add section, got: %s", controls)
	}

	// Test that each control section is on its own line
	uiModel.ActiveTab = TabTrips
	controls = uiModel.renderContextualControls()
	lines := strings.Split(controls, "\n")

	// Should have: NAVIGATION, blank line, ACTIONS, QUICK ADD, DELETE, and potentially empty lines
	var navigationFound, actionsFound, quickAddFound, deleteFound bool
	var blankLineAfterNavigation bool

	for i, line := range lines {
		if strings.Contains(line, "NAVIGATION:") {
			navigationFound = true
			// Check if next line is blank (for separation)
			if i+1 < len(lines) && strings.TrimSpace(lines[i+1]) == "" {
				blankLineAfterNavigation = true
			}
		}
		if strings.Contains(line, "ACTIONS:") {
			actionsFound = true
		}
		if strings.Contains(line, "QUICK ADD:") {
			quickAddFound = true
		}
		if strings.Contains(line, "DELETE:") {
			deleteFound = true
		}
	}

	if !navigationFound {
		t.Error("Expected to find NAVIGATION section")
	}
	if !actionsFound {
		t.Error("Expected to find ACTIONS section")
	}
	if !quickAddFound {
		t.Error("Expected to find QUICK ADD section")
	}
	if !deleteFound {
		t.Error("Expected to find DELETE section")
	}
	if !blankLineAfterNavigation {
		t.Error("Expected blank line after NAVIGATION section")
	}
}

// TestHelpSystemFunctionality tests the Phase 2 help system
func TestHelpSystemFunctionality(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Test initial state (help not visible)
	if uiModel.HelpVisible {
		t.Error("Help should not be visible initially")
	}
	if uiModel.HelpLevel != 0 {
		t.Errorf("Expected help level to be 0, got %d", uiModel.HelpLevel)
	}

	// Test F1 key (Quick Help)
	updatedModel, _ := uiModel.Update(tea.KeyMsg{Type: tea.KeyF1})
	uiModel = updatedModel.(*Model)
	if !uiModel.HelpVisible {
		t.Error("Help should be visible after F1")
	}
	if uiModel.HelpLevel != 1 {
		t.Errorf("Expected help level to be 1, got %d", uiModel.HelpLevel)
	}

	// Test F2 key (Detailed Help)
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyF2})
	uiModel = updatedModel.(*Model)
	if !uiModel.HelpVisible {
		t.Error("Help should be visible after F2")
	}
	if uiModel.HelpLevel != 2 {
		t.Errorf("Expected help level to be 2, got %d", uiModel.HelpLevel)
	}

	// Test F3 key (Advanced Help)
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyF3})
	uiModel = updatedModel.(*Model)
	if !uiModel.HelpVisible {
		t.Error("Help should be visible after F3")
	}
	if uiModel.HelpLevel != 3 {
		t.Errorf("Expected help level to be 3, got %d", uiModel.HelpLevel)
	}

	// Test Esc key to close help
	updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEsc})
	uiModel = updatedModel.(*Model)
	if uiModel.HelpVisible {
		t.Error("Help should not be visible after Esc")
	}
}

func TestHelpContentGeneration(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Test help content for different tabs and levels
	tabs := []int{TabWeeklySummaries, TabTrips, TabExpenses, TabTemplates}
	levels := []int{1, 2, 3}

	for _, tab := range tabs {
		for _, level := range levels {
			uiModel.ActiveTab = tab
			uiModel.HelpLevel = level

			content := uiModel.getHelpContent()

			// Verify content is not empty
			if content == "" {
				t.Errorf("Help content should not be empty for tab %d, level %d", tab, level)
			}

			// Verify level-specific content
			switch level {
			case 1:
				if !strings.Contains(content, "Quick Help [F1]") {
					t.Errorf("Quick help should contain F1 indicator for tab %d", tab)
				}
			case 2:
				if !strings.Contains(content, "Detailed Help [F2]") {
					t.Errorf("Detailed help should contain F2 indicator for tab %d", tab)
				}
			case 3:
				if !strings.Contains(content, "Advanced Help [F3]") {
					t.Errorf("Advanced help should contain F3 indicator for tab %d", tab)
				}
				if !strings.Contains(content, "DEVELOPER TOOLS") {
					t.Errorf("Advanced help should contain developer tools for tab %d", tab)
				}
			}

			// Verify common navigation is always present
			if !strings.Contains(content, "NAVIGATION") {
				t.Errorf("Help content should contain navigation section for tab %d, level %d", tab, level)
			}
		}
	}
}

func TestHelpOverlayRendering(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Test help overlay when not visible
	overlay := uiModel.renderHelpOverlay()
	if overlay != "" {
		t.Error("Help overlay should be empty when not visible")
	}

	// Test help overlay when visible
	uiModel.HelpVisible = true
	uiModel.HelpLevel = 1
	overlay = uiModel.renderHelpOverlay()
	if overlay == "" {
		t.Error("Help overlay should not be empty when visible")
	}
}

func TestHelpSystemIntegration(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Test help system integration with main view
	view := uiModel.View()
	if !strings.Contains(view, "HELP:") {
		t.Error("Main view should contain help hints")
	}
	if !strings.Contains(view, "[F1] Quick") {
		t.Error("Main view should contain F1 help hint")
	}
	if !strings.Contains(view, "[F2] Detailed") {
		t.Error("Main view should contain F2 help hint")
	}
	if !strings.Contains(view, "[F3] Advanced") {
		t.Error("Main view should contain F3 help hint")
	}

	// Test help overlay appears in view when visible
	uiModel.HelpVisible = true
	uiModel.HelpLevel = 2
	view = uiModel.View()
	if !strings.Contains(view, "Detailed Help [F2]") {
		t.Error("View should contain help overlay when help is visible")
	}
}

func TestHelpContentContextAwareness(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Test context-specific help content
	testCases := []struct {
		tab             int
		expectedContent string
		level           int
	}{
		{TabWeeklySummaries, "WEEKLY SUMMARIES", 1},
		{TabTrips, "TRIPS", 1},
		{TabExpenses, "EXPENSES", 1},
		{TabTemplates, "TEMPLATES", 1},
		{TabTrips, "TRIP TIPS", 2},
		{TabExpenses, "EXPENSE TIPS", 2},
		{TabTemplates, "TEMPLATE TIPS", 2},
		{TabTrips, "ADVANCED TRIP FEATURES", 3},
		{TabExpenses, "ADVANCED EXPENSE FEATURES", 3},
		{TabTemplates, "ADVANCED TEMPLATE FEATURES", 3},
	}

	for _, tc := range testCases {
		uiModel.ActiveTab = tc.tab
		uiModel.HelpLevel = tc.level
		content := uiModel.getHelpContent()

		if !strings.Contains(content, tc.expectedContent) {
			t.Errorf("Help content for tab %d, level %d should contain '%s'", tc.tab, tc.level, tc.expectedContent)
		}
	}
}

func TestHelpSystemKeyHandling(t *testing.T) {
	uiModel, cleanup := setupTestUI(t)
	defer cleanup()

	// Test that help keys work in different application states
	modes := []string{"date", "origin", "destination", "type", "search"}

	for _, mode := range modes {
		uiModel.Mode = mode
		uiModel.HelpVisible = false

		// Test F1 key
		updatedModel, _ := uiModel.Update(tea.KeyMsg{Type: tea.KeyF1})
		uiModel = updatedModel.(*Model)
		if !uiModel.HelpVisible {
			t.Errorf("F1 should work in mode %s", mode)
		}

		// Test Esc to close
		updatedModel, _ = uiModel.Update(tea.KeyMsg{Type: tea.KeyEsc})
		uiModel = updatedModel.(*Model)
		if uiModel.HelpVisible {
			t.Errorf("Esc should close help in mode %s", mode)
		}
	}
}
