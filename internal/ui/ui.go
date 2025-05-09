package ui

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lauren/nannytracker/internal/maps"
	"github.com/lauren/nannytracker/internal/model"
	"github.com/lauren/nannytracker/internal/storage"
)

// Model represents the UI state
type Model struct {
	TextInput         textinput.Model
	Trips             []model.Trip
	RecurringTrips    []model.RecurringTrip
	CurrentTrip       model.Trip
	CurrentRecurring  model.RecurringTrip
	CurrentExpense    model.Expense
	Mode              string // "date", "origin", "destination", "type", "edit", "delete", "delete_confirm", "expense_date", "expense_amount", "expense_description", "expense_edit", "expense_delete_confirm", "search", "recurring_date", "recurring_weekday", "recurring_end_date", "convert_to_recurring"
	Err               error
	Storage           storage.Storage
	RatePerMile       float64
	MapsClient        maps.DistanceCalculator
	Data              *model.StorageData
	EditIndex         int    // Index of trip being edited
	SelectedTrip      int    // Index of selected trip for operations
	SelectedRecurring int    // Index of selected recurring trip for operations
	SelectedExpense   int    // Index of selected expense for operations
	SearchQuery       string // Current search query
	SearchMode        bool   // Whether we're in search mode
}

// New creates a new UI model with a mock maps client (for backward compatibility)
func New(storage storage.Storage, ratePerMile float64) (*Model, error) {
	data, err := storage.LoadData()
	if err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
	}

	// Calculate weekly summaries after loading data
	model.CalculateAndUpdateWeeklySummaries(data, ratePerMile)

	ti := textinput.New()
	ti.Placeholder = "Enter date (YYYY-MM-DD)..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50

	return &Model{
		TextInput:   ti,
		Trips:       data.Trips,
		CurrentTrip: model.Trip{},
		Mode:        "date",
		Storage:     storage,
		RatePerMile: ratePerMile,
		MapsClient:  maps.NewMockClient(),
		Data:        data,
	}, nil
}

// NewWithClient creates a new UI model with a provided maps client (useful for testing)
func NewWithClient(storage storage.Storage, ratePerMile float64, mapsClient maps.DistanceCalculator) (*Model, error) {
	data, err := storage.LoadData()
	if err != nil {
		// Initialize empty data if loading fails
		data = &model.StorageData{
			Trips:           make([]model.Trip, 0),
			WeeklySummaries: make([]model.WeeklySummary, 0),
		}
	}

	ti := textinput.New()
	ti.Placeholder = "Enter date (YYYY-MM-DD)..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50

	return &Model{
		TextInput:   ti,
		Trips:       data.Trips,
		CurrentTrip: model.Trip{},
		Mode:        "date",
		Storage:     storage,
		RatePerMile: ratePerMile,
		MapsClient:  mapsClient,
		Data:        data,
		EditIndex:   -1,
	}, nil
}

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages and updates the model accordingly
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyCtrlE:
			// Enter edit mode
			if m.Mode == "date" {
				if len(m.Trips) > 0 && m.SelectedTrip < len(m.Trips) {
					m.CurrentTrip = m.Trips[m.SelectedTrip]
					m.EditIndex = m.SelectedTrip
					m.Mode = "edit"
					m.TextInput.SetValue(m.CurrentTrip.Date)
					m.TextInput.Placeholder = "Enter date (YYYY-MM-DD)..."
				}
			}
			return m, cmd
		case tea.KeyCtrlF:
			if m.SearchMode {
				m.SearchMode = false
				m.SearchQuery = ""
				m.TextInput.Reset()
				m.Mode = "date"
				m.TextInput.Placeholder = "Enter date (YYYY-MM-DD)..."
			} else {
				m.SearchMode = true
				m.Mode = "search"
				m.TextInput.Reset()
				m.TextInput.Placeholder = "Enter search term..."
			}
			return m, cmd
		case tea.KeyCtrlR:
			// Toggle recurring trip mode or convert selected trip to recurring
			if m.Mode == "date" {
				if m.SelectedTrip >= 0 && m.SelectedTrip < len(m.Trips) {
					// Convert selected trip to recurring
					m.CurrentTrip = m.Trips[m.SelectedTrip]
					m.CurrentRecurring = model.RecurringTrip{
						Origin:      m.CurrentTrip.Origin,
						Destination: m.CurrentTrip.Destination,
						Miles:       m.CurrentTrip.Miles,
						StartDate:   m.CurrentTrip.Date,
						Type:        m.CurrentTrip.Type,
					}
					m.Mode = "convert_to_recurring"
					m.TextInput.Placeholder = "Enter weekday (0-6, where 0 is Sunday)..."
				} else {
					// Start new recurring trip
					m.Mode = "recurring_date"
					m.TextInput.Placeholder = "Enter start date (YYYY-MM-DD)..."
				}
			} else if strings.HasPrefix(m.Mode, "recurring_") || m.Mode == "convert_to_recurring" {
				m.Mode = "date"
				m.TextInput.Placeholder = "Enter date (YYYY-MM-DD)..."
			}
			return m, cmd
		case tea.KeyEnter:
			if m.Mode == "convert_to_recurring" {
				weekday, err := strconv.Atoi(m.TextInput.Value())
				if err != nil || weekday < 0 || weekday > 6 {
					m.Err = fmt.Errorf("invalid weekday: must be between 0 and 6")
					return m, cmd
				}
				m.CurrentRecurring.Weekday = weekday

				// Set end date to end of current month
				var now time.Time
				if m.Data.ReferenceDate != "" {
					now, err = time.Parse("2006-01-02", m.Data.ReferenceDate)
					if err != nil {
						m.Err = err
						return m, cmd
					}
				} else {
					now = time.Now()
				}
				endOfMonth := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location())
				m.CurrentRecurring.EndDate = endOfMonth.Format("2006-01-02")

				// Validate the recurring trip
				if err := m.CurrentRecurring.Validate(); err != nil {
					m.Err = fmt.Errorf("invalid recurring trip: %w", err)
					return m, cmd
				}

				// Delete the original trip first
				if err := m.Data.DeleteTrip(m.SelectedTrip); err != nil {
					m.Err = err
					return m, cmd
				}
				m.Trips = m.Data.Trips

				// Add the recurring trip
				if err := m.Data.AddRecurringTrip(m.CurrentRecurring); err != nil {
					m.Err = err
					return m, cmd
				}
				m.RecurringTrips = m.Data.RecurringTrips

				// Generate trips from recurring trips
				if err := m.Data.GenerateTripsFromRecurring(); err != nil {
					m.Err = err
					return m, cmd
				}
				m.Trips = m.Data.Trips

				// Update weekly summaries
				model.CalculateAndUpdateWeeklySummaries(m.Data, m.RatePerMile)
				if err := m.Storage.SaveData(m.Data); err != nil {
					m.Err = err
					return m, cmd
				}

				// Reset state
				m.CurrentRecurring = model.RecurringTrip{}
				m.Mode = "date"
				m.TextInput.Reset()
				m.TextInput.Placeholder = "Enter date (YYYY-MM-DD)..."
			} else if m.Mode == "search" {
				m.SearchQuery = m.TextInput.Value()
			} else if m.Mode == "date" {
				if m.TextInput.Value() == "" {
					return m, cmd
				}
				// Create a temporary trip to validate the date
				tempTrip := model.Trip{
					Date:        m.TextInput.Value(),
					Origin:      "temp",   // Dummy value for validation
					Destination: "temp",   // Dummy value for validation
					Miles:       1.0,      // Dummy value for validation
					Type:        "single", // Dummy value for validation
				}
				if err := tempTrip.Validate(); err != nil {
					m.Err = err
					return m, cmd
				}
				m.CurrentTrip.Date = m.TextInput.Value()
				m.TextInput.Reset()
				m.Mode = "origin"
				m.TextInput.Placeholder = "Enter origin location..."
			} else if m.Mode == "recurring_date" {
				if m.TextInput.Value() == "" {
					return m, cmd
				}
				// Create a temporary recurring trip to validate the date
				tempTrip := model.RecurringTrip{
					StartDate:   m.TextInput.Value(),
					Origin:      "temp",   // Dummy value for validation
					Destination: "temp",   // Dummy value for validation
					Miles:       1.0,      // Dummy value for validation
					Type:        "single", // Dummy value for validation
					Weekday:     0,        // Dummy value for validation
				}
				if err := tempTrip.Validate(); err != nil {
					m.Err = err
					return m, cmd
				}
				m.CurrentRecurring.StartDate = m.TextInput.Value()
				m.TextInput.Reset()
				m.Mode = "recurring_weekday"
				m.TextInput.Placeholder = "Enter weekday (0-6, where 0 is Sunday)..."
			} else if m.Mode == "recurring_weekday" {
				weekday, err := strconv.Atoi(m.TextInput.Value())
				if err != nil || weekday < 0 || weekday > 6 {
					m.Err = fmt.Errorf("invalid weekday: must be between 0 and 6")
					return m, cmd
				}
				m.CurrentRecurring.Weekday = weekday
				m.TextInput.Reset()
				m.Mode = "origin"
				m.TextInput.Placeholder = "Enter origin location..."
			} else if m.Mode == "recurring_end_date" {
				if m.TextInput.Value() != "" {
					// Create a temporary recurring trip to validate the end date
					tempTrip := model.RecurringTrip{
						StartDate:   m.CurrentRecurring.StartDate,
						EndDate:     m.TextInput.Value(),
						Origin:      "temp",   // Dummy value for validation
						Destination: "temp",   // Dummy value for validation
						Miles:       1.0,      // Dummy value for validation
						Type:        "single", // Dummy value for validation
						Weekday:     0,        // Dummy value for validation
					}
					if err := tempTrip.Validate(); err != nil {
						m.Err = err
						return m, cmd
					}
					m.CurrentRecurring.EndDate = m.TextInput.Value()
				}
				m.TextInput.Reset()
				m.Mode = "origin"
				m.TextInput.Placeholder = "Enter origin location..."
			} else if m.Mode == "origin" {
				if m.TextInput.Value() == "" && m.EditIndex >= 0 {
					// Keep existing value if no new input
					m.TextInput.Reset()
					m.Mode = "destination"
					m.TextInput.Placeholder = "Enter destination location..."
				} else {
					if strings.HasPrefix(m.Mode, "recurring_") {
						m.CurrentRecurring.Origin = m.TextInput.Value()
					} else {
						m.CurrentTrip.Origin = m.TextInput.Value()
					}
					m.TextInput.Reset()
					m.Mode = "destination"
					m.TextInput.Placeholder = "Enter destination location..."
				}
			} else if m.Mode == "destination" {
				if m.TextInput.Value() == "" && m.EditIndex >= 0 {
					// Keep existing value if no new input
					m.TextInput.Reset()
					m.Mode = "type"
					m.TextInput.Placeholder = "Enter trip type (single/round)..."
				} else {
					if strings.HasPrefix(m.Mode, "recurring_") {
						m.CurrentRecurring.Destination = m.TextInput.Value()
						// Calculate miles using maps client
						miles, err := m.MapsClient.CalculateDistance(context.Background(), m.CurrentRecurring.Origin, m.CurrentRecurring.Destination)
						if err != nil {
							m.Err = fmt.Errorf("failed to calculate distance: %w", err)
							return m, cmd
						}
						m.CurrentRecurring.Miles = miles
					} else {
						m.CurrentTrip.Destination = m.TextInput.Value()
						// Calculate miles using maps client
						miles, err := m.MapsClient.CalculateDistance(context.Background(), m.CurrentTrip.Origin, m.CurrentTrip.Destination)
						if err != nil {
							m.Err = fmt.Errorf("failed to calculate distance: %w", err)
							return m, cmd
						}
						m.CurrentTrip.Miles = miles
					}
					m.TextInput.Reset()
					m.Mode = "type"
					m.TextInput.Placeholder = "Enter trip type (single/round)..."
				}
			} else if m.Mode == "type" {
				if m.TextInput.Value() == "" && m.EditIndex >= 0 {
					// Keep existing value if no new input
					tripType := m.CurrentTrip.Type
					if tripType != "single" && tripType != "round" {
						m.Err = fmt.Errorf("invalid trip type: %s. Must be 'single' or 'round'", tripType)
						return m, cmd
					}
				} else {
					tripType := strings.ToLower(m.TextInput.Value())
					if tripType != "single" && tripType != "round" {
						m.Err = fmt.Errorf("invalid trip type: %s. Must be 'single' or 'round'", tripType)
						return m, cmd
					}
					if strings.HasPrefix(m.Mode, "recurring_") {
						m.CurrentRecurring.Type = tripType
					} else {
						m.CurrentTrip.Type = tripType
					}
				}

				if strings.HasPrefix(m.Mode, "recurring_") {
					// Validate the recurring trip before saving
					if err := m.CurrentRecurring.Validate(); err != nil {
						m.Err = fmt.Errorf("invalid recurring trip: %w", err)
						return m, cmd
					}

					if m.EditIndex >= 0 {
						// Update existing recurring trip
						if err := m.Data.EditRecurringTrip(m.EditIndex, m.CurrentRecurring); err != nil {
							m.Err = err
							return m, cmd
						}
						m.RecurringTrips[m.EditIndex] = m.CurrentRecurring
					} else {
						// Add new recurring trip
						newTrip := m.CurrentRecurring // Create a copy to avoid reference issues
						m.Data.RecurringTrips = append(m.Data.RecurringTrips, newTrip)
						m.RecurringTrips = m.Data.RecurringTrips
					}

					// Generate trips from recurring trips
					if err := m.Data.GenerateTripsFromRecurring(); err != nil {
						m.Err = err
						return m, cmd
					}

					// Update the UI state with the generated trips
					m.Trips = m.Data.Trips

					// Update weekly summaries
					model.CalculateAndUpdateWeeklySummaries(m.Data, m.RatePerMile)
					if err := m.Storage.SaveData(m.Data); err != nil {
						m.Err = err
						return m, cmd
					}

					// Reset state
					m.EditIndex = -1
					m.CurrentRecurring = model.RecurringTrip{}
					m.Mode = "date"
					m.TextInput.Reset()
					m.TextInput.Placeholder = "Enter date (YYYY-MM-DD)..."
				} else {
					// Validate the trip before saving
					if err := m.CurrentTrip.Validate(); err != nil {
						m.Err = fmt.Errorf("invalid trip: %w", err)
						return m, cmd
					}

					if m.EditIndex >= 0 {
						// Update existing trip
						if err := m.Data.EditTrip(m.EditIndex, m.CurrentTrip); err != nil {
							m.Err = err
							return m, cmd
						}
						m.Trips[m.EditIndex] = m.CurrentTrip
					} else {
						// Add new trip
						newTrip := m.CurrentTrip // Create a copy to avoid reference issues
						m.Data.Trips = append(m.Data.Trips, newTrip)
						m.Trips = m.Data.Trips
					}

					model.CalculateAndUpdateWeeklySummaries(m.Data, m.RatePerMile)
					if err := m.Storage.SaveData(m.Data); err != nil {
						m.Err = err
						return m, cmd
					}

					// Reset state
					m.EditIndex = -1
					m.CurrentTrip = model.Trip{}
					m.Mode = "date"
					m.TextInput.Reset()
					m.TextInput.Placeholder = "Enter date (YYYY-MM-DD)..."
				}
			} else if m.Mode == "edit" {
				if m.TextInput.Value() != "" {
					// Create a temporary trip to validate the date
					tempTrip := model.Trip{
						Date:        m.TextInput.Value(),
						Origin:      "temp",   // Dummy value for validation
						Destination: "temp",   // Dummy value for validation
						Miles:       1.0,      // Dummy value for validation
						Type:        "single", // Dummy value for validation
					}
					if err := tempTrip.Validate(); err != nil {
						m.Err = err
						return m, cmd
					}
					m.CurrentTrip.Date = m.TextInput.Value()
				}
				m.TextInput.Reset()
				m.Mode = "origin"
				m.TextInput.Placeholder = "Enter origin location..."
			} else if m.Mode == "delete_confirm" {
				if strings.ToLower(m.TextInput.Value()) == "yes" {
					if err := m.Data.DeleteTrip(m.SelectedTrip); err != nil {
						m.Err = err
						return m, cmd
					}
					m.Trips = m.Data.Trips
					model.CalculateAndUpdateWeeklySummaries(m.Data, m.RatePerMile)
					if err := m.Storage.SaveData(m.Data); err != nil {
						m.Err = err
						return m, cmd
					}
				}
				m.Mode = "date"
				m.TextInput.Reset()
				m.TextInput.Placeholder = "Enter date (YYYY-MM-DD)..."
			} else if m.Mode == "expense_date" {
				if m.TextInput.Value() == "" {
					return m, cmd
				}
				// Create a temporary expense to validate the date
				tempExpense := model.Expense{
					Date:        m.TextInput.Value(),
					Amount:      1.0,    // Dummy value for validation
					Description: "temp", // Dummy value for validation
				}
				if err := tempExpense.Validate(); err != nil {
					m.Err = err
					return m, cmd
				}
				m.CurrentExpense.Date = m.TextInput.Value()
				m.TextInput.Reset()
				m.Mode = "expense_amount"
				m.TextInput.Placeholder = "Enter expense amount..."
			} else if m.Mode == "expense_amount" {
				if m.TextInput.Value() == "" && m.EditIndex >= 0 {
					// Keep existing value if no new input
					m.TextInput.Reset()
					m.Mode = "expense_description"
					m.TextInput.Placeholder = "Enter expense description..."
				} else {
					amount, err := strconv.ParseFloat(m.TextInput.Value(), 64)
					if err != nil {
						m.Err = fmt.Errorf("invalid amount: %w", err)
						return m, cmd
					}
					if amount <= 0 {
						m.Err = fmt.Errorf("amount must be greater than 0")
						return m, cmd
					}
					m.CurrentExpense.Amount = amount
					m.TextInput.Reset()
					m.Mode = "expense_description"
					m.TextInput.Placeholder = "Enter expense description..."
				}
			} else if m.Mode == "expense_description" {
				if m.TextInput.Value() == "" {
					m.Err = fmt.Errorf("description cannot be empty")
					return m, cmd
				}
				m.CurrentExpense.Description = m.TextInput.Value()

				// Validate the expense before saving
				if err := m.CurrentExpense.Validate(); err != nil {
					m.Err = fmt.Errorf("invalid expense: %w", err)
					return m, cmd
				}

				// Add new expense
				if err := m.Data.AddExpense(m.CurrentExpense); err != nil {
					m.Err = err
					return m, cmd
				}

				model.CalculateAndUpdateWeeklySummaries(m.Data, m.RatePerMile)
				if err := m.Storage.SaveData(m.Data); err != nil {
					m.Err = err
					return m, cmd
				}

				// Reset state
				m.CurrentExpense = model.Expense{}
				m.Mode = "date"
				m.TextInput.Reset()
				m.TextInput.Placeholder = "Enter date (YYYY-MM-DD)..."
			}
			return m, cmd
		case tea.KeyCtrlX:
			// Enter expense mode
			m.Mode = "expense_date"
			m.TextInput.Reset()
			m.TextInput.Placeholder = "Enter expense date (YYYY-MM-DD)..."
			return m, cmd
		case tea.KeyCtrlD:
			// Enter delete confirmation mode
			if m.Mode == "date" {
				m.Mode = "delete_confirm"
				m.TextInput.Reset()
				m.TextInput.Placeholder = "Type 'yes' to confirm deletion..."
			}
			return m, cmd
		case tea.KeyUp:
			// Navigate up in the trips or expenses list
			if m.Mode == "date" {
				if m.SelectedExpense >= 0 && len(m.Data.Expenses) > 0 {
					// Navigate expenses
					m.SelectedExpense = (m.SelectedExpense - 1 + len(m.Data.Expenses)) % len(m.Data.Expenses)
					m.SelectedTrip = -1
				} else if len(m.Trips) > 0 {
					// Navigate trips
					m.SelectedTrip = (m.SelectedTrip - 1 + len(m.Trips)) % len(m.Trips)
					m.SelectedExpense = -1
				}
			}
			return m, cmd
		case tea.KeyDown:
			// Navigate down in the trips or expenses list
			if m.Mode == "date" {
				if m.SelectedExpense >= 0 && len(m.Data.Expenses) > 0 {
					// Navigate expenses
					m.SelectedExpense = (m.SelectedExpense + 1) % len(m.Data.Expenses)
					m.SelectedTrip = -1
				} else if len(m.Trips) > 0 {
					// Navigate trips
					m.SelectedTrip = (m.SelectedTrip + 1) % len(m.Trips)
					m.SelectedExpense = -1
				}
			}
			return m, cmd
		case tea.KeyTab:
			// Switch between trips and expenses
			if m.Mode == "date" {
				if m.SelectedExpense == -1 && len(m.Data.Expenses) > 0 {
					m.SelectedExpense = 0
					m.SelectedTrip = -1
				} else if m.SelectedTrip == -1 && len(m.Trips) > 0 {
					m.SelectedTrip = 0
					m.SelectedExpense = -1
				} else if m.SelectedExpense >= 0 {
					// If we're already in expenses, switch to trips
					m.SelectedTrip = 0
					m.SelectedExpense = -1
				} else if m.SelectedTrip >= 0 {
					// If we're already in trips, switch to expenses
					m.SelectedExpense = 0
					m.SelectedTrip = -1
				}
			}
			return m, cmd
		}

		// Handle search input
		if m.Mode == "search" {
			m.SearchQuery = m.TextInput.Value()
		}
	}

	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

// filterBySearch filters trips based on the search query
func (m *Model) filterBySearch() []model.Trip {
	if m.SearchQuery == "" {
		return m.Trips
	}

	query := strings.ToLower(m.SearchQuery)
	var filteredTrips []model.Trip

	// Filter trips
	for _, trip := range m.Trips {
		if strings.Contains(strings.ToLower(trip.Origin), query) ||
			strings.Contains(strings.ToLower(trip.Destination), query) ||
			strings.Contains(strings.ToLower(trip.Date), query) ||
			strings.Contains(strings.ToLower(trip.Type), query) {
			filteredTrips = append(filteredTrips, trip)
		}
	}

	return filteredTrips
}

// View renders the UI
func (m *Model) View() string {
	var s strings.Builder

	// Title style
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFA500")) // Orange
	title := titleStyle.Render("Nanny Tracker")
	underline := titleStyle.Render(strings.Repeat("─", len("Nanny Tracker")))
	s.WriteString(title + "\n" + underline + "\n\n")

	// Create styles
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00FF00")).
		Padding(0, 1)

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")).
		Padding(0, 1)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFF00")).
		Bold(true)

	editingStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FFFF")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))

	// Show error if any
	if m.Err != nil {
		s.WriteString(errorStyle.Render(m.Err.Error()) + "\n\n")
		m.Err = nil
	}

	// Show mode and input field
	s.WriteString(headerStyle.Render(fmt.Sprintf("Mode: %s", m.Mode)) + "\n")
	s.WriteString(m.TextInput.View() + "\n\n")

	// Show weekly summaries
	if len(m.Data.WeeklySummaries) > 0 {
		s.WriteString(headerStyle.Render("Weekly Summaries:") + "\n")
		for _, summary := range m.Data.WeeklySummaries {
			s.WriteString(fmt.Sprintf("Week of %s to %s:\n", summary.WeekStart, summary.WeekEnd))
			s.WriteString(normalStyle.SetString(fmt.Sprintf("    Total Miles:          %.2f\n"+
				"    Total Mileage Amount: $%.2f\n"+
				"    Total Expenses:       $%.2f\n",
				summary.TotalMiles,
				summary.TotalAmount,
				summary.TotalExpenses)).String())
			s.WriteString("\n") // Ensure a blank line after each summary
		}
	}

	// Get trips to display (filtered or all)
	displayTrips := m.Trips
	if m.SearchMode {
		displayTrips = m.filterBySearch()
	}

	// Show recurring trips
	if len(m.RecurringTrips) > 0 {
		s.WriteString(headerStyle.Render("Recurring Trips:") + "\n")
		for i, trip := range m.RecurringTrips {
			weekday := time.Weekday(trip.Weekday).String()
			displayMiles := trip.Miles
			if trip.Type == "round" {
				displayMiles *= 2
			}
			tripLine := fmt.Sprintf("%s → %s (%.2f miles) [%s] - Every %s",
				trip.Origin, trip.Destination, displayMiles, trip.Type, weekday)

			// Apply appropriate style based on selection/edit state
			if m.EditIndex == i {
				tripLine = editingStyle.Render("> " + tripLine)
			} else if m.SelectedRecurring == i {
				tripLine = selectedStyle.Render("* " + tripLine)
			} else {
				tripLine = normalStyle.Render("  " + tripLine)
			}
			s.WriteString(tripLine + "\n")
		}
		s.WriteString("\n")
	}

	// Show trips
	if len(displayTrips) > 0 {
		s.WriteString(headerStyle.Render("Trips:") + "\n")
		for i, trip := range displayTrips {
			displayMiles := trip.Miles
			if trip.Type == "round" {
				displayMiles *= 2
			}
			tripLine := fmt.Sprintf("%s → %s (%.2f miles) [%s]", trip.Origin, trip.Destination, displayMiles, trip.Type)

			// Apply appropriate style based on selection/edit state
			if m.EditIndex == i {
				tripLine = editingStyle.Render("> " + tripLine)
			} else if m.SelectedTrip == i {
				tripLine = selectedStyle.Render("* " + tripLine)
			} else {
				tripLine = normalStyle.Render("  " + tripLine)
			}
			s.WriteString(tripLine + "\n")
		}
		s.WriteString("\n")
	}

	// Show expenses
	if len(m.Data.Expenses) > 0 {
		s.WriteString(headerStyle.Render("Expenses:") + "\n")
		for i, expense := range m.Data.Expenses {
			expenseLine := fmt.Sprintf("%s: $%.2f - %s", expense.Date, expense.Amount, expense.Description)
			if m.SelectedExpense == i {
				expenseLine = selectedStyle.Render("* " + expenseLine)
			} else {
				expenseLine = normalStyle.Render("  " + expenseLine)
			}
			s.WriteString(expenseLine + "\n")
		}
		s.WriteString("\n")
	}

	// Show help text
	s.WriteString("Controls:\n")
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		SetString("↑/↓: Navigate trips and expenses\n" +
			"Tab: Switch between trips and expenses\n" +
			"Ctrl+E: Edit selected trip or expense\n" +
			"Ctrl+D: Delete selected trip or expense\n" +
			"Ctrl+F: Toggle search mode\n" +
			"Ctrl+X: Add expense\n" +
			"Ctrl+R: Toggle recurring trip mode or convert selected trip to recurring\n")
	s.WriteString(helpStyle.String())

	return s.String()
}

// CalculateTotalMiles calculates the total miles for a list of trips
func (m *Model) CalculateTotalMiles(trips []model.Trip) float64 {
	return model.CalculateTotalMiles(trips)
}

// CalculateReimbursement calculates the total reimbursement for a list of trips
func (m *Model) CalculateReimbursement(trips []model.Trip, ratePerMile float64) float64 {
	return model.CalculateReimbursement(trips, ratePerMile)
}

// AddTrip adds a new trip to the model's trips list and updates weekly summaries
func (m *Model) AddTrip(trip model.Trip) {
	m.Trips = append(m.Trips, trip)
	m.Data.Trips = m.Trips
	model.CalculateAndUpdateWeeklySummaries(m.Data, m.RatePerMile)
	if err := m.Storage.SaveData(m.Data); err != nil {
		m.Err = err
	}
}
