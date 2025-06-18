package ui

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/laurendc/nannytracker/pkg/core/maps"
	model "github.com/laurendc/nannytracker/pkg/core"
	"github.com/laurendc/nannytracker/pkg/core/storage"
)

// Model represents the UI state
type Model struct {
	TextInput         textinput.Model
	Trips             []model.Trip
	RecurringTrips    []model.RecurringTrip
	CurrentTrip       model.Trip
	CurrentRecurring  model.RecurringTrip
	CurrentExpense    model.Expense
	Mode              string // "date", "origin", "destination", "type", "edit", "delete", "delete_confirm", "expense_date", "expense_amount", "expense_description", "expense_edit", "expense_delete_confirm", "search", "recurring_date", "recurring_weekday", "recurring_end_date", "convert_to_recurring", "template_name", "template_origin", "template_destination", "template_type", "template_notes", "template_edit", "template_delete_confirm"
	Err               error
	Storage           storage.Storage
	RatePerMile       float64
	MapsClient        maps.DistanceCalculator
	Data              *model.StorageData
	EditIndex         int                  // Index of trip being edited
	SelectedTrip      int                  // Index of selected trip for operations
	SelectedRecurring int                  // Index of selected recurring trip for operations
	SelectedExpense   int                  // Index of selected expense for operations
	SearchQuery       string               // Current search query
	SearchMode        bool                 // Whether we're in search mode
	ActiveTab         int                  // Index of the active tab (0: Weekly Summaries, 1: Trips, 2: Expenses, 3: Templates)
	SelectedWeek      int                  // Index of the currently selected week in WeeklySummaries
	PageSize          int                  // Number of items to show per page
	CurrentPage       int                  // Current page number (0-based)
	TripTemplates     []model.TripTemplate // List of saved trip templates
	SelectedTemplate  int                  // Index of selected template for operations
	CurrentTemplate   model.TripTemplate   // Current template being edited
	JustChangedMode   bool                 // Flag to prevent double-processing after mode change
}

const (
	TabWeeklySummaries = iota
	TabTrips
	TabExpenses
	TabTemplates
)

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

	m := &Model{
		TextInput:         ti,
		Trips:             data.Trips,
		CurrentTrip:       model.Trip{},
		Mode:              "date",
		Storage:           storage,
		RatePerMile:       ratePerMile,
		MapsClient:        maps.NewMockClient(),
		Data:              data,
		EditIndex:         -1,
		SelectedTrip:      -1,
		SelectedExpense:   -1,
		SelectedRecurring: -1,
		SelectedTemplate:  -1,
		PageSize:          10, // Default page size
		CurrentPage:       0,  // Start at first page
		TripTemplates:     data.TripTemplates,
	}
	m.SelectedWeek = m.getCurrentWeekIndex()
	return m, nil
}

// NewWithClient creates a new UI model with a provided maps client (useful for testing)
func NewWithClient(storage storage.Storage, ratePerMile float64, mapsClient maps.DistanceCalculator) (*Model, error) {
	data, err := storage.LoadData()
	if err != nil {
		// Initialize empty data if loading fails
		data = &model.StorageData{
			Trips:           make([]model.Trip, 0),
			WeeklySummaries: make([]model.WeeklySummary, 0),
			TripTemplates:   make([]model.TripTemplate, 0),
		}
	}

	ti := textinput.New()
	ti.Placeholder = "Enter date (YYYY-MM-DD)..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50

	m := &Model{
		TextInput:         ti,
		Trips:             data.Trips,
		CurrentTrip:       model.Trip{},
		Mode:              "date",
		Storage:           storage,
		RatePerMile:       ratePerMile,
		MapsClient:        mapsClient,
		Data:              data,
		EditIndex:         -1,
		SelectedTrip:      -1,
		SelectedExpense:   -1,
		SelectedRecurring: -1,
		SelectedTemplate:  -1,
		PageSize:          10, // Default page size
		CurrentPage:       0,  // Start at first page
		TripTemplates:     data.TripTemplates,
	}
	m.SelectedWeek = m.getCurrentWeekIndex()
	return m, nil
}

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages and updates the model accordingly
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	// Update text input
	m.TextInput, cmd = m.TextInput.Update(msg)
	cmds = append(cmds, cmd)

	// Handle key messages
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyCtrlE:
			if m.ActiveTab == TabTrips && m.SelectedTrip >= 0 {
				m.Mode = "edit"
				m.EditIndex = 0
				m.CurrentTrip = m.Trips[m.SelectedTrip]
				m.TextInput.SetValue(m.CurrentTrip.Date)
				m.TextInput.Placeholder = "Enter date (YYYY-MM-DD)..."
			} else if m.ActiveTab == TabTemplates && m.SelectedTemplate >= 0 {
				m.Mode = "template_edit"
				m.EditIndex = m.SelectedTemplate
				m.CurrentTemplate = m.TripTemplates[m.SelectedTemplate]
				m.TextInput.SetValue(m.CurrentTemplate.Name)
				m.TextInput.Placeholder = "Enter template name..."
			}
		case tea.KeyCtrlD:
			if m.ActiveTab == TabTrips && m.SelectedTrip >= 0 {
				m.Mode = "delete_confirm"
				m.TextInput.Reset()
				m.TextInput.Placeholder = "Type 'yes' and press Enter to confirm deletion, or anything else to cancel."
			} else if m.ActiveTab == TabTemplates && m.SelectedTemplate >= 0 {
				m.Mode = "template_delete_confirm"
				m.TextInput.Reset()
				m.TextInput.Placeholder = "Type 'yes' and press Enter to confirm deletion, or anything else to cancel."
			}
		case tea.KeyEnter:
			if m.Mode == "date" {
				if m.TextInput.Value() == "" {
					return m, cmd
				}
				// Validate just the date format
				if err := model.ValidateDate(m.TextInput.Value()); err != nil {
					m.Err = err
					return m, cmd
				}
				m.CurrentTrip.Date = m.TextInput.Value()

				// If CurrentTrip already has origin, destination, and type (i.e., from a template), prompt for origin with pre-filled value
				if m.CurrentTrip.Origin != "" && m.CurrentTrip.Destination != "" && m.CurrentTrip.Type != "" {
					m.TextInput.Reset()
					m.TextInput.SetValue(m.CurrentTrip.Origin)
					m.Mode = "origin"
					m.TextInput.Placeholder = "Enter origin location..."
					return m, cmd
				}

				// Otherwise, continue normal flow
				m.TextInput.Reset()
				m.Mode = "origin"
				m.TextInput.Placeholder = "Enter origin location..."
				return m, cmd
			}
			if m.Mode == "origin" {
				if m.TextInput.Value() == "" {
					return m, cmd
				}
				m.CurrentTrip.Origin = m.TextInput.Value()
				m.TextInput.Reset()
				// If destination is already set (from template), pre-fill it
				if m.CurrentTrip.Destination != "" {
					m.TextInput.SetValue(m.CurrentTrip.Destination)
				} else {
					m.TextInput.SetValue("")
				}
				m.Mode = "destination"
				m.TextInput.Placeholder = "Enter destination location..."
				return m, cmd
			}
			if m.Mode == "destination" {
				if m.TextInput.Value() == "" {
					return m, cmd
				}
				m.CurrentTrip.Destination = m.TextInput.Value()
				m.TextInput.Reset()
				// If type is already set (from template), pre-fill it
				if m.CurrentTrip.Type != "" {
					m.TextInput.SetValue(m.CurrentTrip.Type)
				} else {
					m.TextInput.SetValue("")
				}
				m.Mode = "type"
				m.TextInput.Placeholder = "Enter trip type (single/round)..."
				return m, cmd
			} else if m.Mode == "edit" {
				if m.TextInput.Value() == "" {
					return m, cmd
				}
				m.CurrentTrip.Date = m.TextInput.Value()
				m.TextInput.Reset()
				m.Mode = "edit_origin"
				m.TextInput.Placeholder = "Enter origin location..."
				m.TextInput.SetValue(m.CurrentTrip.Origin)
				m.EditIndex = 1
				return m, cmd
			} else if m.Mode == "edit_origin" {
				if m.TextInput.Value() != "" {
					m.CurrentTrip.Origin = m.TextInput.Value()
				}
				m.TextInput.Reset()
				m.Mode = "edit_destination"
				m.TextInput.Placeholder = "Enter destination location..."
				m.TextInput.SetValue(m.CurrentTrip.Destination)
				m.EditIndex = 2
				return m, cmd
			} else if m.Mode == "edit_destination" {
				if m.TextInput.Value() != "" {
					m.CurrentTrip.Destination = m.TextInput.Value()
				}
				m.TextInput.Reset()
				m.Mode = "edit_type"
				m.TextInput.Placeholder = "Enter trip type (single/round)..."
				m.TextInput.SetValue(m.CurrentTrip.Type)
				m.EditIndex = 3
				return m, cmd
			} else if m.Mode == "edit_type" {
				if m.TextInput.Value() != "" {
					tripType := strings.ToLower(m.TextInput.Value())
					if tripType != "single" && tripType != "round" {
						m.Err = fmt.Errorf("invalid trip type: %s. Must be 'single' or 'round'", tripType)
						return m, cmd
					}
					m.CurrentTrip.Type = tripType
				}
				// Save edited trip
				if err := m.CurrentTrip.Validate(); err != nil {
					m.Err = fmt.Errorf("invalid trip: %w", err)
					return m, cmd
				}
				if m.SelectedTrip >= 0 && m.SelectedTrip < len(m.Trips) {
					m.Trips[m.SelectedTrip] = m.CurrentTrip
					m.Data.Trips = m.Trips
					model.CalculateAndUpdateWeeklySummaries(m.Data, m.RatePerMile)
					if err := m.Storage.SaveData(m.Data); err != nil {
						m.Err = fmt.Errorf("failed to save trip: %w", err)
						return m, cmd
					}
				}
				m.EditIndex = -1
				m.CurrentTrip = model.Trip{}
				m.Mode = "date"
				m.TextInput.Reset()
				m.TextInput.Placeholder = "Enter date (YYYY-MM-DD)..."
				return m, cmd
			} else if m.Mode == "template_edit" {
				// Name input
				if m.TextInput.Value() != "" {
					m.CurrentTemplate.Name = m.TextInput.Value()
				}
				m.TextInput.Reset()
				m.TextInput.SetValue(m.CurrentTemplate.Origin)
				m.TextInput.Placeholder = "Enter origin location..."
				m.Mode = "template_edit_origin"
				return m, cmd
			} else if m.Mode == "template_edit_origin" {
				if m.TextInput.Value() != "" {
					m.CurrentTemplate.Origin = m.TextInput.Value()
				}
				m.TextInput.Reset()
				m.TextInput.SetValue(m.CurrentTemplate.Destination)
				m.TextInput.Placeholder = "Enter destination location..."
				m.Mode = "template_edit_destination"
				return m, cmd
			} else if m.Mode == "template_edit_destination" {
				if m.TextInput.Value() != "" {
					m.CurrentTemplate.Destination = m.TextInput.Value()
				}
				m.TextInput.Reset()
				m.TextInput.SetValue(m.CurrentTemplate.TripType)
				m.TextInput.Placeholder = "Enter trip type (single/round)..."
				m.Mode = "template_edit_type"
				return m, cmd
			} else if m.Mode == "template_edit_type" {
				if m.TextInput.Value() != "" {
					tripType := strings.ToLower(m.TextInput.Value())
					if tripType != "single" && tripType != "round" {
						m.Err = fmt.Errorf("invalid trip type: %s. Must be 'single' or 'round'", tripType)
						return m, cmd
					}
					m.CurrentTemplate.TripType = tripType
				}
				m.TextInput.Reset()
				m.TextInput.SetValue(m.CurrentTemplate.Notes)
				m.TextInput.Placeholder = "Enter notes (optional, press Enter to skip)..."
				m.Mode = "template_edit_notes"
				return m, cmd
			} else if m.Mode == "template_edit_notes" {
				if m.TextInput.Value() != "" {
					m.CurrentTemplate.Notes = m.TextInput.Value()
				}
				// Validate the template before saving
				if err := m.CurrentTemplate.Validate(); err != nil {
					m.Err = fmt.Errorf("invalid template: %w", err)
					return m, cmd
				}
				if m.EditIndex >= 0 {
					if err := m.Data.EditTripTemplate(m.EditIndex, m.CurrentTemplate); err != nil {
						m.Err = err
						return m, cmd
					}
					m.TripTemplates[m.EditIndex] = m.CurrentTemplate
				}
				if err := m.Storage.SaveData(m.Data); err != nil {
					m.Err = err
					return m, cmd
				}
				m.EditIndex = -1
				m.CurrentTemplate = model.TripTemplate{}
				m.Mode = "date"
				m.TextInput.Reset()
				m.TextInput.Placeholder = "Enter date (YYYY-MM-DD)..."
				return m, cmd
			} else if m.Mode == "template_name" {
				if m.TextInput.Value() == "" {
					return m, cmd
				}
				m.CurrentTemplate.Name = m.TextInput.Value()
				m.TextInput.Reset()
				m.Mode = "template_origin"
				m.TextInput.Placeholder = "Enter origin location..."
			} else if m.Mode == "template_origin" {
				if m.TextInput.Value() == "" && m.EditIndex >= 0 {
					// Keep existing value if no new input
					m.TextInput.Reset()
					m.Mode = "template_destination"
					m.TextInput.Placeholder = "Enter destination location..."
				} else {
					m.CurrentTemplate.Origin = m.TextInput.Value()
					m.TextInput.Reset()
					m.Mode = "template_destination"
					m.TextInput.Placeholder = "Enter destination location..."
				}
			} else if m.Mode == "template_destination" {
				if m.TextInput.Value() == "" && m.EditIndex >= 0 {
					// Keep existing value if no new input
					m.TextInput.Reset()
					m.Mode = "template_type"
					m.TextInput.Placeholder = "Enter trip type (single/round)..."
				} else {
					m.CurrentTemplate.Destination = m.TextInput.Value()
					m.TextInput.Reset()
					m.Mode = "template_type"
					m.TextInput.Placeholder = "Enter trip type (single/round)..."
				}
			} else if m.Mode == "template_type" {
				if m.TextInput.Value() == "" && m.EditIndex >= 0 {
					// Keep existing value if no new input
					tripType := m.CurrentTemplate.TripType
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
					m.CurrentTemplate.TripType = tripType
				}
				m.TextInput.Reset()
				m.TextInput.SetValue(m.CurrentTemplate.Notes)
				m.TextInput.Placeholder = "Enter notes (optional, press Enter to skip)..."
				m.Mode = "template_notes"
			} else if m.Mode == "template_notes" {
				m.CurrentTemplate.Notes = m.TextInput.Value()
				// Validate the template before saving
				if err := m.CurrentTemplate.Validate(); err != nil {
					m.Err = fmt.Errorf("invalid template: %w", err)
					return m, cmd
				}

				if m.EditIndex >= 0 {
					// Update existing template
					if err := m.Data.EditTripTemplate(m.EditIndex, m.CurrentTemplate); err != nil {
						m.Err = err
						return m, cmd
					}
					m.TripTemplates[m.EditIndex] = m.CurrentTemplate
				} else {
					// Add new template
					if err := m.Data.AddTripTemplate(m.CurrentTemplate); err != nil {
						m.Err = err
						return m, cmd
					}
					m.TripTemplates = append(m.TripTemplates, m.CurrentTemplate)
					m.Data.TripTemplates = m.TripTemplates
				}

				if err := m.Storage.SaveData(m.Data); err != nil {
					m.Err = err
					return m, cmd
				}

				// Reset state
				m.EditIndex = -1
				m.CurrentTemplate = model.TripTemplate{}
				m.Mode = "date"
				m.TextInput.Reset()
				m.TextInput.Placeholder = "Enter date (YYYY-MM-DD)..."
			} else if m.Mode == "convert_to_recurring" {
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
			} else if m.Mode == "recurring_date" {
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
					// Calculate miles if not already set
					if m.CurrentTrip.Miles == 0 {
						distance, err := m.MapsClient.CalculateDistance(context.Background(), m.CurrentTrip.Origin, m.CurrentTrip.Destination)
						if err != nil {
							m.Err = fmt.Errorf("failed to calculate distance: %w", err)
							return m, cmd
						}
						m.CurrentTrip.Miles = distance
					}

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
				return m, cmd
			} else if m.Mode == "delete_confirm" {
				if m.TextInput.Value() == "yes" {
					if m.SelectedTrip >= 0 && m.SelectedTrip < len(m.Trips) {
						// Remove the trip
						m.Trips = append(m.Trips[:m.SelectedTrip], m.Trips[m.SelectedTrip+1:]...)
						m.Data.Trips = m.Trips
						model.CalculateAndUpdateWeeklySummaries(m.Data, m.RatePerMile)
						if err := m.Storage.SaveData(m.Data); err != nil {
							m.Err = fmt.Errorf("failed to save after deletion: %w", err)
							return m, cmd
						}
						m.SelectedTrip = -1
					}
				}
				m.Mode = "date"
				m.TextInput.Reset()
				m.TextInput.Placeholder = "Enter date (YYYY-MM-DD)..."
				return m, cmd
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
			} else if m.Mode == "template_delete_confirm" {
				if m.SelectedTemplate >= 0 && m.SelectedTemplate < len(m.TripTemplates) {
					if m.TextInput.Value() == "yes" {
						// Remove the template
						m.TripTemplates = append(m.TripTemplates[:m.SelectedTemplate], m.TripTemplates[m.SelectedTemplate+1:]...)
						m.Data.TripTemplates = m.TripTemplates
						if err := m.Storage.SaveData(m.Data); err != nil {
							m.Err = fmt.Errorf("failed to save after deletion: %w", err)
							return m, cmd
						}
						m.SelectedTemplate = -1
					}
				}
				m.Mode = "date"
				m.TextInput.Reset()
				m.TextInput.Placeholder = "Enter date (YYYY-MM-DD)..."
			}
		case tea.KeyCtrlX:
			// Enter expense mode
			m.Mode = "expense_date"
			m.TextInput.Reset()
			m.TextInput.Placeholder = "Enter expense date (YYYY-MM-DD)..."
			return m, cmd
		case tea.KeyCtrlT:
			// Enter template creation mode
			if m.Mode == "date" {
				m.Mode = "template_name"
				m.TextInput.Reset()
				m.TextInput.Placeholder = "Enter template name..."
			}
			return m, cmd
		case tea.KeyUp:
			if m.ActiveTab == TabTrips {
				if len(m.Trips) == 0 {
					return m, cmd
				}
				if m.SelectedTrip <= 0 {
					m.SelectedTrip = len(m.Trips) - 1
				} else {
					m.SelectedTrip--
				}
				m.SelectedExpense = -1
				m.SelectedTemplate = -1
			} else if m.ActiveTab == TabExpenses {
				if len(m.Data.Expenses) == 0 {
					return m, cmd
				}
				if m.SelectedExpense <= 0 {
					m.SelectedExpense = len(m.Data.Expenses) - 1
				} else {
					m.SelectedExpense--
				}
				m.SelectedTrip = -1
				m.SelectedTemplate = -1
			} else if m.ActiveTab == TabTemplates {
				if len(m.TripTemplates) == 0 {
					return m, cmd
				}
				if m.SelectedTemplate <= 0 {
					m.SelectedTemplate = len(m.TripTemplates) - 1
				} else {
					m.SelectedTemplate--
				}
				m.SelectedTrip = -1
				m.SelectedExpense = -1
			}
		case tea.KeyDown:
			if m.ActiveTab == TabTrips {
				if len(m.Trips) == 0 {
					return m, cmd
				}
				if m.SelectedTrip >= len(m.Trips)-1 {
					m.SelectedTrip = 0
				} else {
					m.SelectedTrip++
				}
				m.SelectedExpense = -1
				m.SelectedTemplate = -1
			} else if m.ActiveTab == TabExpenses {
				if len(m.Data.Expenses) == 0 {
					return m, cmd
				}
				if m.SelectedExpense >= len(m.Data.Expenses)-1 {
					m.SelectedExpense = 0
				} else {
					m.SelectedExpense++
				}
				m.SelectedTrip = -1
				m.SelectedTemplate = -1
			} else if m.ActiveTab == TabTemplates {
				if len(m.TripTemplates) == 0 {
					return m, cmd
				}
				if m.SelectedTemplate >= len(m.TripTemplates)-1 {
					m.SelectedTemplate = 0
				} else {
					m.SelectedTemplate++
				}
				m.SelectedTrip = -1
				m.SelectedExpense = -1
			}
		case tea.KeyLeft:
			if m.ActiveTab == TabWeeklySummaries && len(m.Data.WeeklySummaries) > 0 {
				if m.SelectedWeek > 0 {
					m.SelectedWeek--
				}
			} else if m.ActiveTab == TabTrips {
				if m.CurrentPage > 0 {
					m.CurrentPage--
					// Adjust selected trip to stay within the current page
					if m.SelectedTrip >= 0 {
						m.SelectedTrip = 0
					}
				}
			} else if m.ActiveTab == TabExpenses {
				if m.CurrentPage > 0 {
					m.CurrentPage--
					// Adjust selected expense to stay within the current page
					if m.SelectedExpense >= 0 {
						m.SelectedExpense = 0
					}
				}
			} else if m.ActiveTab == TabTemplates {
				if m.CurrentPage > 0 {
					m.CurrentPage--
					// Adjust selected template to stay within the current page
					if m.SelectedTemplate >= 0 {
						m.SelectedTemplate = 0
					}
				}
			}
			return m, cmd
		case tea.KeyRight:
			if m.ActiveTab == TabWeeklySummaries && len(m.Data.WeeklySummaries) > 0 {
				if m.SelectedWeek < len(m.Data.WeeklySummaries)-1 {
					m.SelectedWeek++
				}
			} else if m.ActiveTab == TabTrips {
				displayTrips := m.Trips
				if m.SearchMode {
					displayTrips = m.filterBySearch()
				}
				// Sort trips in descending order (most recent first)
				sort.Slice(displayTrips, func(i, j int) bool {
					return displayTrips[i].Date > displayTrips[j].Date
				})
				if m.CurrentPage < (len(displayTrips)-1)/m.PageSize {
					m.CurrentPage++
					// Adjust selected trip to stay within the current page
					if m.SelectedTrip >= 0 {
						m.SelectedTrip = 0
					}
				}
			} else if m.ActiveTab == TabExpenses {
				if m.CurrentPage < (len(m.Data.Expenses)-1)/m.PageSize {
					m.CurrentPage++
					// Adjust selected expense to stay within the current page
					if m.SelectedExpense >= 0 {
						m.SelectedExpense = 0
					}
				}
			} else if m.ActiveTab == TabTemplates {
				if m.CurrentPage < (len(m.TripTemplates)-1)/m.PageSize {
					m.CurrentPage++
					// Adjust selected template to stay within the current page
					if m.SelectedTemplate >= 0 {
						m.SelectedTemplate = 0
					}
				}
			}
			return m, cmd
		case tea.KeyPgUp:
			// Remove Page Up handler since we're using left arrow
			return m, cmd
		case tea.KeyPgDown:
			// Remove Page Down handler since we're using right arrow
			return m, cmd
		case tea.KeyCtrlR:
			if m.ActiveTab == TabTrips {
				if m.SelectedTrip >= 0 && m.SelectedTrip < len(m.Trips) {
					trip := m.Trips[m.SelectedTrip]
					m.Mode = "convert_to_recurring"
					m.CurrentRecurring = model.RecurringTrip{
						Origin:      trip.Origin,
						Destination: trip.Destination,
						Miles:       trip.Miles,
						StartDate:   trip.Date,
						Type:        trip.Type,
					}
					m.TextInput.Reset()
					m.TextInput.Placeholder = "Enter weekday (0=Sunday, 6=Saturday)..."
					return m, cmd
				} else {
					m.Mode = "recurring_date"
					m.CurrentRecurring = model.RecurringTrip{}
					m.TextInput.Reset()
					m.TextInput.Placeholder = "Enter start date (YYYY-MM-DD)..."
					return m, cmd
				}
			}
		case tea.KeyTab:
			// Cycle forward through tabs: Weekly Summaries -> Trips -> Expenses -> Templates -> Weekly Summaries
			switch m.ActiveTab {
			case TabWeeklySummaries:
				m.ActiveTab = TabTrips
			case TabTrips:
				m.ActiveTab = TabExpenses
			case TabExpenses:
				m.ActiveTab = TabTemplates
			case TabTemplates:
				m.ActiveTab = TabWeeklySummaries
				// Refresh weekly summaries when switching to Weekly Summaries tab
				model.CalculateAndUpdateWeeklySummaries(m.Data, m.RatePerMile)
			}
			// Reset selections when changing tabs
			m.CurrentPage = 0
			if m.ActiveTab == TabWeeklySummaries {
				if len(m.Data.WeeklySummaries) > 0 && (m.SelectedWeek < 0 || m.SelectedWeek >= len(m.Data.WeeklySummaries)) {
					m.SelectedWeek = 0
				} else if len(m.Data.WeeklySummaries) == 0 {
					m.SelectedWeek = -1
				}
			} else {
				m.SelectedWeek = -1
			}
			m.SelectedTrip = -1
			m.SelectedExpense = -1
			m.SelectedTemplate = -1
			return m, cmd
		case tea.KeyShiftTab:
			// Cycle backward through tabs: Weekly Summaries -> Templates -> Expenses -> Trips -> Weekly Summaries
			switch m.ActiveTab {
			case TabWeeklySummaries:
				m.ActiveTab = TabTemplates
			case TabTemplates:
				m.ActiveTab = TabExpenses
			case TabExpenses:
				m.ActiveTab = TabTrips
			case TabTrips:
				m.ActiveTab = TabWeeklySummaries
				// Refresh weekly summaries when switching to Weekly Summaries tab
				model.CalculateAndUpdateWeeklySummaries(m.Data, m.RatePerMile)
			}
			// Reset selections when changing tabs
			m.CurrentPage = 0
			if m.ActiveTab == TabWeeklySummaries {
				m.SelectedWeek = m.getCurrentWeekIndex() // Show week containing today
			} else {
				m.SelectedWeek = -1
			}
			m.SelectedTrip = -1
			m.SelectedExpense = -1
			m.SelectedTemplate = -1
			return m, cmd
		case tea.KeyCtrlU:
			if m.ActiveTab == TabTemplates && m.SelectedTemplate >= 0 {
				// Create a new trip from the selected template
				template := m.TripTemplates[m.SelectedTemplate]
				m.CurrentTrip = model.Trip{
					Origin:      template.Origin,
					Destination: template.Destination,
					Type:        template.TripType,
					Miles:       0, // Will be calculated when the trip is saved
				}
				m.Mode = "date"
				m.TextInput.Reset()
				m.TextInput.Placeholder = "Enter date (YYYY-MM-DD)..."
				// Switch to trips tab
				m.ActiveTab = TabTrips
				m.SelectedTrip = -1
				m.SelectedTemplate = -1
				return m, cmd
			}
		}

		// Handle search input
		if m.Mode == "search" {
			m.SearchQuery = m.TextInput.Value()
		}
	}

	// Only return early for edit modes after handling key events
	if m.Mode == "edit" || m.Mode == "edit_origin" || m.Mode == "edit_destination" || m.Mode == "edit_type" {
		return m, tea.Batch(cmds...)
	}

	return m, tea.Batch(cmds...)
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

	// Tab styles
	tabStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Padding(0, 1)

	activeTabStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")).
		Bold(true).
		Padding(0, 1)

	// Show error if any
	if m.Err != nil {
		s.WriteString(errorStyle.Render(m.Err.Error()) + "\n\n")
		m.Err = nil
	}

	// Show mode and input field
	s.WriteString(headerStyle.Render(fmt.Sprintf("Mode: %s", m.Mode)) + "\n")
	s.WriteString(m.TextInput.View() + "\n\n")

	// Render tabs
	tabs := []string{"Weekly Summaries", "Trips", "Expenses", "Trip Templates"}
	var tabLine strings.Builder
	for i, tab := range tabs {
		if i == m.ActiveTab {
			tabLine.WriteString(activeTabStyle.Render(tab))
		} else {
			tabLine.WriteString(tabStyle.Render(tab))
		}
		if i < len(tabs)-1 {
			tabLine.WriteString(" | ")
		}
	}
	s.WriteString(tabLine.String() + "\n\n")

	// Show content based on active tab
	switch m.ActiveTab {
	case TabWeeklySummaries:
		if m.SelectedWeek >= 0 && m.SelectedWeek < len(m.Data.WeeklySummaries) {
			summary := m.Data.WeeklySummaries[m.SelectedWeek]
			s.WriteString(headerStyle.Render(fmt.Sprintf("Week of %s to %s (Week %d of %d):", summary.WeekStart, summary.WeekEnd, m.SelectedWeek+1, len(m.Data.WeeklySummaries))) + "\n")
			s.WriteString(normalStyle.Render(fmt.Sprintf("    Total Miles:          %.2f", summary.TotalMiles)) + "\n")
			s.WriteString(normalStyle.Render(fmt.Sprintf("    Total Mileage Amount: $%.2f", summary.TotalAmount)) + "\n")
			s.WriteString(normalStyle.Render(fmt.Sprintf("    Total Expenses:       $%.2f", summary.TotalExpenses)) + "\n")
			s.WriteString(normalStyle.Render(" Trips:") + "\n")
			for _, trip := range summary.Trips {
				displayMiles := trip.Miles
				if trip.Type == "round" {
					displayMiles *= 2
				}
				tripLine := fmt.Sprintf(" %s: %s → %s (%.2f miles) [%s]", trip.Date, trip.Origin, trip.Destination, displayMiles, trip.Type)
				s.WriteString(normalStyle.Render(tripLine) + "\n")
			}
			s.WriteString("\n")
			s.WriteString(normalStyle.Render(" Expenses:") + "\n")
			if len(summary.Expenses) > 0 {
				for _, exp := range summary.Expenses {
					s.WriteString(normalStyle.Render(fmt.Sprintf(" %s: $%.2f - %s", exp.Date, exp.Amount, exp.Description)) + "\n")
				}
			} else {
				s.WriteString(normalStyle.Render(" (No expenses available.)") + "\n")
			}
		} else {
			s.WriteString(normalStyle.Render(" (No weekly summary available.)") + "\n")
		}

	case TabTrips:
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

		// Show regular trips with pagination
		if len(displayTrips) > 0 {
			s.WriteString(headerStyle.Render("Regular Trips:") + "\n")

			// Sort trips in descending order (most recent first)
			sort.Slice(displayTrips, func(i, j int) bool {
				return displayTrips[i].Date > displayTrips[j].Date
			})
			startIdx := m.CurrentPage * m.PageSize
			endIdx := startIdx + m.PageSize
			if endIdx > len(displayTrips) {
				endIdx = len(displayTrips)
			}

			// Display trips for current page
			for i := startIdx; i < endIdx; i++ {
				trip := displayTrips[i]
				displayMiles := trip.Miles
				if trip.Type == "round" {
					displayMiles *= 2
				}
				tripLine := fmt.Sprintf("%s: %s → %s (%.2f miles) [%s]",
					trip.Date, trip.Origin, trip.Destination, displayMiles, trip.Type)

				if m.EditIndex == i {
					tripLine = editingStyle.Render("> " + tripLine)
				} else if m.SelectedTrip == i {
					tripLine = selectedStyle.Render("* " + tripLine)
				} else {
					tripLine = normalStyle.Render("  " + tripLine)
				}
				s.WriteString(tripLine + "\n")
			}

			// Show pagination info
			totalPages := (len(displayTrips) + m.PageSize - 1) / m.PageSize
			if totalPages > 1 {
				paginationInfo := fmt.Sprintf("\nPage %d of %d (Showing %d-%d of %d trips)",
					m.CurrentPage+1, totalPages, startIdx+1, endIdx, len(displayTrips))
				s.WriteString(normalStyle.Render(paginationInfo) + "\n")
			}
		} else {
			s.WriteString(normalStyle.Render("No trips available.\n"))
		}

	case TabExpenses:
		if len(m.Data.Expenses) > 0 {
			// Sort expenses by date in descending order
			sort.Slice(m.Data.Expenses, func(i, j int) bool {
				return m.Data.Expenses[i].Date > m.Data.Expenses[j].Date
			})

			// Calculate pagination
			startIdx := m.CurrentPage * m.PageSize
			endIdx := startIdx + m.PageSize
			if endIdx > len(m.Data.Expenses) {
				endIdx = len(m.Data.Expenses)
			}

			// Display expenses for current page
			for i := startIdx; i < endIdx; i++ {
				expense := m.Data.Expenses[i]
				expenseLine := fmt.Sprintf("%s: $%.2f - %s", expense.Date, expense.Amount, expense.Description)
				if m.SelectedExpense == i {
					expenseLine = selectedStyle.Render("* " + expenseLine)
				} else {
					expenseLine = normalStyle.Render("  " + expenseLine)
				}
				s.WriteString(expenseLine + "\n")
			}

			// Show pagination info
			totalPages := (len(m.Data.Expenses) + m.PageSize - 1) / m.PageSize
			if totalPages > 1 {
				paginationInfo := fmt.Sprintf("\nPage %d of %d (Showing %d-%d of %d expenses)",
					m.CurrentPage+1, totalPages, startIdx+1, endIdx, len(m.Data.Expenses))
				s.WriteString(normalStyle.Render(paginationInfo) + "\n")
			}
		} else {
			s.WriteString(normalStyle.Render("No expenses available.\n"))
		}

	case TabTemplates:
		// Show templates with pagination
		if len(m.TripTemplates) > 0 {
			s.WriteString(headerStyle.Render("Trip Templates:") + "\n")

			// Create a copy of templates for sorting
			displayTemplates := make([]model.TripTemplate, len(m.TripTemplates))
			copy(displayTemplates, m.TripTemplates)

			// Sort templates alphabetically by name (case-insensitive)
			sort.Slice(displayTemplates, func(i, j int) bool {
				return strings.ToLower(displayTemplates[i].Name) < strings.ToLower(displayTemplates[j].Name)
			})

			startIdx := m.CurrentPage * m.PageSize
			endIdx := startIdx + m.PageSize
			if endIdx > len(displayTemplates) {
				endIdx = len(displayTemplates)
			}

			// Display sorted templates
			for i := startIdx; i < endIdx; i++ {
				template := displayTemplates[i]
				templateLine := fmt.Sprintf("%s: %s → %s [%s]",
					template.Name, template.Origin, template.Destination, template.TripType)
				if template.Notes != "" {
					templateLine += fmt.Sprintf(" - %s", template.Notes)
				}

				// Find the original index in m.TripTemplates for selection highlighting
				originalIndex := -1
				for j, t := range m.TripTemplates {
					if t.Name == template.Name && t.Origin == template.Origin &&
						t.Destination == template.Destination && t.TripType == template.TripType {
						originalIndex = j
						break
					}
				}

				if m.SelectedTemplate == originalIndex {
					templateLine = selectedStyle.Render("* " + templateLine)
				} else {
					templateLine = normalStyle.Render("  " + templateLine)
				}
				s.WriteString(templateLine + "\n")
			}

			// Show pagination info if there are multiple pages
			if len(m.TripTemplates) > m.PageSize {
				totalPages := (len(m.TripTemplates) + m.PageSize - 1) / m.PageSize
				s.WriteString(fmt.Sprintf("\nPage %d of %d\n", m.CurrentPage+1, totalPages))
			}
		} else {
			s.WriteString(normalStyle.Render("No trip templates available.\n"))
		}

	}

	s.WriteString("\n")

	// Show help text
	s.WriteString("Controls:\n")
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		SetString("↑/↓: Navigate items\n" +
			"Tab/Shift+Tab: Switch tabs\n" +
			"←/→: Switch weeks/pages\n" +
			"Ctrl+E: Edit selected item\n" +
			"Ctrl+D: Delete selected item\n" +
			"Ctrl+F: Toggle search mode\n" +
			"Ctrl+X: Add expense\n" +
			"Ctrl+R: Toggle recurring trip mode or convert selected trip to recurring\n" +
			"Ctrl+T: Create new trip template\n" +
			"Ctrl+U: Use selected template to create a new trip\n")
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

// Helper: find the index of the week containing today
func (m *Model) getCurrentWeekIndex() int {
	today := time.Now().Format("2006-01-02")
	todayTime, err := time.Parse("2006-01-02", today)
	if err != nil {
		return 0
	}
	for i, summary := range m.Data.WeeklySummaries {
		start, err1 := time.Parse("2006-01-02", summary.WeekStart)
		end, err2 := time.Parse("2006-01-02", summary.WeekEnd)
		if err1 != nil || err2 != nil {
			continue
		}
		if !todayTime.Before(start) && !todayTime.After(end) {
			return i
		}
	}
	return 0 // fallback to most recent week
}
