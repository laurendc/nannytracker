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
	model "github.com/laurendc/nannytracker/pkg/core"
	"github.com/laurendc/nannytracker/pkg/core/maps"
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
	Width             int                  // Terminal width in characters
	// Phase 2: Help System
	HelpVisible bool // Whether help overlay is visible
	HelpLevel   int  // Help level: 1=Quick, 2=Detailed, 3=Advanced
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

	// Handle window size updates
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
	}

	// Update text input
	m.TextInput, cmd = m.TextInput.Update(msg)
	cmds = append(cmds, cmd)

	// Handle key messages
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			if m.HelpVisible {
				// Close help overlay
				m.HelpVisible = false
				return m, cmd
			}
			return m, tea.Quit
		case tea.KeyF1:
			m.HelpVisible = true
			m.HelpLevel = 1
			return m, cmd
		case tea.KeyF2:
			m.HelpVisible = true
			m.HelpLevel = 2
			return m, cmd
		case tea.KeyF3:
			m.HelpVisible = true
			m.HelpLevel = 3
			return m, cmd
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
				// For templates, we need to navigate through the sorted display order
				// Create a sorted copy to determine the next/previous template
				displayTemplates := make([]model.TripTemplate, len(m.TripTemplates))
				copy(displayTemplates, m.TripTemplates)
				sort.Slice(displayTemplates, func(i, j int) bool {
					return strings.ToLower(displayTemplates[i].Name) < strings.ToLower(displayTemplates[j].Name)
				})

				// Find current selection in sorted order
				currentSortedIndex := -1
				for i, displayTemplate := range displayTemplates {
					if m.SelectedTemplate >= 0 && m.SelectedTemplate < len(m.TripTemplates) {
						originalTemplate := m.TripTemplates[m.SelectedTemplate]
						if displayTemplate.Name == originalTemplate.Name &&
							displayTemplate.Origin == originalTemplate.Origin &&
							displayTemplate.Destination == originalTemplate.Destination &&
							displayTemplate.TripType == originalTemplate.TripType {
							currentSortedIndex = i
							break
						}
					}
				}

				// Navigate in sorted order
				if currentSortedIndex <= 0 {
					// Go to last template
					if len(displayTemplates) > 0 {
						lastTemplate := displayTemplates[len(displayTemplates)-1]
						// Find original index of last template
						for i, originalTemplate := range m.TripTemplates {
							if lastTemplate.Name == originalTemplate.Name &&
								lastTemplate.Origin == originalTemplate.Origin &&
								lastTemplate.Destination == originalTemplate.Destination &&
								lastTemplate.TripType == originalTemplate.TripType {
								m.SelectedTemplate = i
								break
							}
						}
					}
				} else {
					// Go to previous template
					prevTemplate := displayTemplates[currentSortedIndex-1]
					// Find original index of previous template
					for i, originalTemplate := range m.TripTemplates {
						if prevTemplate.Name == originalTemplate.Name &&
							prevTemplate.Origin == originalTemplate.Origin &&
							prevTemplate.Destination == originalTemplate.Destination &&
							prevTemplate.TripType == originalTemplate.TripType {
							m.SelectedTemplate = i
							break
						}
					}
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
				// For templates, we need to navigate through the sorted display order
				// Create a sorted copy to determine the next/previous template
				displayTemplates := make([]model.TripTemplate, len(m.TripTemplates))
				copy(displayTemplates, m.TripTemplates)
				sort.Slice(displayTemplates, func(i, j int) bool {
					return strings.ToLower(displayTemplates[i].Name) < strings.ToLower(displayTemplates[j].Name)
				})

				// Find current selection in sorted order
				currentSortedIndex := -1
				for i, displayTemplate := range displayTemplates {
					if m.SelectedTemplate >= 0 && m.SelectedTemplate < len(m.TripTemplates) {
						originalTemplate := m.TripTemplates[m.SelectedTemplate]
						if displayTemplate.Name == originalTemplate.Name &&
							displayTemplate.Origin == originalTemplate.Origin &&
							displayTemplate.Destination == originalTemplate.Destination &&
							displayTemplate.TripType == originalTemplate.TripType {
							currentSortedIndex = i
							break
						}
					}
				}

				// Navigate in sorted order
				if currentSortedIndex >= len(displayTemplates)-1 {
					// Go to first template
					if len(displayTemplates) > 0 {
						firstTemplate := displayTemplates[0]
						// Find original index of first template
						for i, originalTemplate := range m.TripTemplates {
							if firstTemplate.Name == originalTemplate.Name &&
								firstTemplate.Origin == originalTemplate.Origin &&
								firstTemplate.Destination == originalTemplate.Destination &&
								firstTemplate.TripType == originalTemplate.TripType {
								m.SelectedTemplate = i
								break
							}
						}
					}
				} else {
					// Go to next template
					nextTemplate := displayTemplates[currentSortedIndex+1]
					// Find original index of next template
					for i, originalTemplate := range m.TripTemplates {
						if nextTemplate.Name == originalTemplate.Name &&
							nextTemplate.Origin == originalTemplate.Origin &&
							nextTemplate.Destination == originalTemplate.Destination &&
							nextTemplate.TripType == originalTemplate.TripType {
							m.SelectedTemplate = i
							break
						}
					}
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
		case tea.KeyRunes:
			// Handle single key presses like "U" for template usage
			if len(msg.Runes) == 1 {
				switch msg.Runes[0] {
				case 'u', 'U':
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

	// Show status bar
	s.WriteString(m.renderStatusBar() + "\n")

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

			// Create a mapping from sorted display index to original index
			displayToOriginal := make(map[int]int)
			for i, displayTemplate := range displayTemplates {
				for j, originalTemplate := range m.TripTemplates {
					if displayTemplate.Name == originalTemplate.Name &&
						displayTemplate.Origin == originalTemplate.Origin &&
						displayTemplate.Destination == originalTemplate.Destination &&
						displayTemplate.TripType == originalTemplate.TripType {
						displayToOriginal[i] = j
						break
					}
				}
			}

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

				// Use the mapping to find the original index for selection highlighting
				originalIndex := displayToOriginal[i]

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

	// Show context-aware controls
	s.WriteString(m.renderContextualControls())

	// Show help overlay if visible
	if m.HelpVisible {
		s.WriteString("\n")
		s.WriteString(m.renderHelpOverlay())
	}

	return s.String()
}

// getHelpContent returns help content based on current tab and help level
func (m *Model) getHelpContent() string {
	var content strings.Builder

	// Create styles for help content
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00FF00"))

	sectionStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFF00"))

	shortcutStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FFFF"))

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))

	tipStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888"))

	// Title based on help level
	switch m.HelpLevel {
	case 1:
		content.WriteString(titleStyle.Render("Quick Help [F1]") + "\n\n")
	case 2:
		content.WriteString(titleStyle.Render("Detailed Help [F2]") + "\n\n")
	case 3:
		content.WriteString(titleStyle.Render("Advanced Help [F3]") + "\n\n")
	}

	// Universal navigation (always shown)
	content.WriteString(sectionStyle.Render("NAVIGATION") + "\n")
	content.WriteString(shortcutStyle.Render("↑/↓") + " " + descStyle.Render("Navigate items") + "\n")
	content.WriteString(shortcutStyle.Render("[Tab]") + " " + descStyle.Render("Switch tabs") + "\n")
	content.WriteString(shortcutStyle.Render("←/→") + " " + descStyle.Render("Navigate pages") + "\n")
	content.WriteString(shortcutStyle.Render("[Enter]") + " " + descStyle.Render("Select item") + "\n")
	content.WriteString(shortcutStyle.Render("[Esc]") + " " + descStyle.Render("Cancel/Close") + "\n")

	if m.HelpLevel >= 2 {
		content.WriteString(shortcutStyle.Render("[Home]") + " " + descStyle.Render("First item") + "\n")
		content.WriteString(shortcutStyle.Render("[End]") + " " + descStyle.Render("Last item") + "\n")
	}

	content.WriteString("\n")

	// Context-specific actions based on active tab
	switch m.ActiveTab {
	case TabWeeklySummaries:
		content.WriteString(sectionStyle.Render("WEEKLY SUMMARIES") + "\n")
		content.WriteString(shortcutStyle.Render("←/→") + " " + descStyle.Render("Switch weeks") + "\n")
		if m.HelpLevel >= 2 {
			content.WriteString(shortcutStyle.Render("[W]") + " " + descStyle.Render("Jump to current week") + "\n")
			content.WriteString(shortcutStyle.Render("[M]") + " " + descStyle.Render("Jump to current month") + "\n")
		}
		if m.HelpLevel >= 3 {
			content.WriteString(shortcutStyle.Render("[P]") + " " + descStyle.Render("Print summary") + "\n")
			content.WriteString(shortcutStyle.Render("[E]") + " " + descStyle.Render("Export week") + "\n")
		}

	case TabTrips:
		content.WriteString(sectionStyle.Render("TRIPS") + "\n")
		content.WriteString(shortcutStyle.Render("[Ctrl+E]") + " " + descStyle.Render("Edit trip") + "\n")
		content.WriteString(shortcutStyle.Render("[Ctrl+F]") + " " + descStyle.Render("Search trips") + "\n")
		content.WriteString(shortcutStyle.Render("[Ctrl+T]") + " " + descStyle.Render("Use template") + "\n")
		content.WriteString(shortcutStyle.Render("[Ctrl+X]") + " " + descStyle.Render("Add expense") + "\n")
		content.WriteString(shortcutStyle.Render("[Ctrl+R]") + " " + descStyle.Render("Add recurring trip") + "\n")
		content.WriteString(shortcutStyle.Render("[Ctrl+D]") + " " + descStyle.Render("Delete trip") + "\n")

		if m.HelpLevel >= 2 {
			content.WriteString("\n" + sectionStyle.Render("TRIP TIPS") + "\n")
			content.WriteString(tipStyle.Render("• Use templates for common routes") + "\n")
			content.WriteString(tipStyle.Render("• Search works on origin, destination, date, type") + "\n")
			content.WriteString(tipStyle.Render("• Round trips automatically double mileage") + "\n")
		}

		if m.HelpLevel >= 3 {
			content.WriteString("\n" + sectionStyle.Render("ADVANCED TRIP FEATURES") + "\n")
			content.WriteString(shortcutStyle.Render("[Ctrl+Shift+E]") + " " + descStyle.Render("Bulk edit trips") + "\n")
			content.WriteString(shortcutStyle.Render("[Ctrl+Shift+D]") + " " + descStyle.Render("Bulk delete trips") + "\n")
			content.WriteString(tipStyle.Render("• Chain templates: Create → Save → Use") + "\n")
			content.WriteString(tipStyle.Render("• Recurring trips generate automatically") + "\n")
		}

	case TabExpenses:
		content.WriteString(sectionStyle.Render("EXPENSES") + "\n")
		content.WriteString(shortcutStyle.Render("[Ctrl+E]") + " " + descStyle.Render("Edit expense") + "\n")
		content.WriteString(shortcutStyle.Render("[Ctrl+F]") + " " + descStyle.Render("Filter expenses") + "\n")
		content.WriteString(shortcutStyle.Render("[Ctrl+X]") + " " + descStyle.Render("Add expense") + "\n")
		content.WriteString(shortcutStyle.Render("[Ctrl+D]") + " " + descStyle.Render("Delete expense") + "\n")

		if m.HelpLevel >= 2 {
			content.WriteString("\n" + sectionStyle.Render("EXPENSE TIPS") + "\n")
			content.WriteString(tipStyle.Render("• Expenses are sorted by date (newest first)") + "\n")
			content.WriteString(tipStyle.Render("• Use clear descriptions for easy tracking") + "\n")
		}

		if m.HelpLevel >= 3 {
			content.WriteString("\n" + sectionStyle.Render("ADVANCED EXPENSE FEATURES") + "\n")
			content.WriteString(shortcutStyle.Render("[Ctrl+Shift+X]") + " " + descStyle.Render("Bulk import expenses") + "\n")
			content.WriteString(shortcutStyle.Render("[Ctrl+Shift+C]") + " " + descStyle.Render("Categorize expenses") + "\n")
		}

	case TabTemplates:
		content.WriteString(sectionStyle.Render("TEMPLATES") + "\n")
		content.WriteString(shortcutStyle.Render("[Ctrl+E]") + " " + descStyle.Render("Edit template") + "\n")
		content.WriteString(shortcutStyle.Render("[Ctrl+F]") + " " + descStyle.Render("Search templates") + "\n")
		content.WriteString(shortcutStyle.Render("[Ctrl+T]") + " " + descStyle.Render("Create template") + "\n")
		content.WriteString(shortcutStyle.Render("[U]") + " " + descStyle.Render("Use template") + "\n")
		content.WriteString(shortcutStyle.Render("[Ctrl+D]") + " " + descStyle.Render("Delete template") + "\n")

		if m.HelpLevel >= 2 {
			content.WriteString("\n" + sectionStyle.Render("TEMPLATE TIPS") + "\n")
			content.WriteString(tipStyle.Render("• Templates speed up common trip entry") + "\n")
			content.WriteString(tipStyle.Render("• Use descriptive names for easy identification") + "\n")
			content.WriteString(tipStyle.Render("• Templates can include notes for context") + "\n")
		}

		if m.HelpLevel >= 3 {
			content.WriteString("\n" + sectionStyle.Render("ADVANCED TEMPLATE FEATURES") + "\n")
			content.WriteString(shortcutStyle.Render("[Ctrl+Shift+T]") + " " + descStyle.Render("Template management") + "\n")
			content.WriteString(shortcutStyle.Render("[Ctrl+Shift+I]") + " " + descStyle.Render("Import templates") + "\n")
			content.WriteString(tipStyle.Render("• Templates support recurring trip patterns") + "\n")
		}
	}

	content.WriteString("\n")

	// Help navigation
	content.WriteString(sectionStyle.Render("HELP NAVIGATION") + "\n")
	content.WriteString(shortcutStyle.Render("[F1]") + " " + descStyle.Render("Quick Help (essentials)") + "\n")
	content.WriteString(shortcutStyle.Render("[F2]") + " " + descStyle.Render("Detailed Help (complete)") + "\n")
	content.WriteString(shortcutStyle.Render("[F3]") + " " + descStyle.Render("Advanced Help (power user)") + "\n")
	content.WriteString(shortcutStyle.Render("[Esc]") + " " + descStyle.Render("Close help") + "\n")

	// Advanced features section (F3 only)
	if m.HelpLevel >= 3 {
		content.WriteString("\n" + sectionStyle.Render("KEYBOARD COMBINATIONS") + "\n")
		content.WriteString(tipStyle.Render("• [Ctrl+F] + [Enter] = Quick search") + "\n")
		content.WriteString(tipStyle.Render("• [Ctrl+E] + [Tab] = Edit next field") + "\n")
		content.WriteString(tipStyle.Render("• [Ctrl+X] + [Ctrl+R] = Add expense to trip") + "\n")

		content.WriteString("\n" + sectionStyle.Render("DEVELOPER TOOLS") + "\n")
		content.WriteString(shortcutStyle.Render("[Ctrl+Shift+L]") + " " + descStyle.Render("Show logs") + "\n")
		content.WriteString(shortcutStyle.Render("[Ctrl+Shift+D]") + " " + descStyle.Render("Debug mode") + "\n")
		content.WriteString(shortcutStyle.Render("[Ctrl+Shift+S]") + " " + descStyle.Render("Save backup") + "\n")

		content.WriteString("\n" + sectionStyle.Render("DATA EXPORT") + "\n")
		content.WriteString(shortcutStyle.Render("[Ctrl+Shift+J]") + " " + descStyle.Render("Export JSON") + "\n")
		content.WriteString(shortcutStyle.Render("[Ctrl+Shift+C]") + " " + descStyle.Render("Export CSV") + "\n")
		content.WriteString(shortcutStyle.Render("[Ctrl+Shift+P]") + " " + descStyle.Render("Export PDF") + "\n")
	}

	return content.String()
}

// renderHelpOverlay renders the help overlay modal
func (m *Model) renderHelpOverlay() string {
	if !m.HelpVisible {
		return ""
	}

	// Create overlay styles
	overlayStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00FF00")).
		Padding(1, 2).
		Margin(1, 2).
		Background(lipgloss.Color("#1a1a1a")).
		Foreground(lipgloss.Color("#FFFFFF"))

	// Get help content
	helpContent := m.getHelpContent()

	// Calculate overlay dimensions based on terminal width
	overlayWidth := 80
	if m.Width > 0 && m.Width < 100 {
		overlayWidth = m.Width - 10 // Leave margin
	}

	// Apply width constraint
	overlayStyle = overlayStyle.Width(overlayWidth)

	return overlayStyle.Render(helpContent)
}

// renderStatusBar renders the status bar with current mode and context information
func (m *Model) renderStatusBar() string {
	// Create status bar styles
	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#333333")).
		Padding(0, 1)

	// Get current tab name
	tabNames := []string{"📊 Weekly Summary", "🚗 Trips", "💰 Expenses", "📋 Templates"}
	currentTab := tabNames[m.ActiveTab]

	// Build status information
	statusInfo := fmt.Sprintf("%s | Mode: %s", currentTab, m.Mode)

	// Add context-specific information
	switch m.ActiveTab {
	case TabTrips:
		if m.SearchMode {
			statusInfo += fmt.Sprintf(" | Search: \"%s\"", m.SearchQuery)
		}
		if len(m.Trips) > 0 {
			statusInfo += fmt.Sprintf(" | %d trips", len(m.Trips))
		}
	case TabExpenses:
		if len(m.Data.Expenses) > 0 {
			statusInfo += fmt.Sprintf(" | %d expenses", len(m.Data.Expenses))
		}
	case TabTemplates:
		if len(m.TripTemplates) > 0 {
			statusInfo += fmt.Sprintf(" | %d templates", len(m.TripTemplates))
		}
	}

	// Add pagination info if applicable
	if m.ActiveTab == TabTrips || m.ActiveTab == TabExpenses || m.ActiveTab == TabTemplates {
		var totalItems int
		switch m.ActiveTab {
		case TabTrips:
			totalItems = len(m.Trips)
		case TabExpenses:
			totalItems = len(m.Data.Expenses)
		case TabTemplates:
			totalItems = len(m.TripTemplates)
		}

		if totalItems > m.PageSize {
			totalPages := (totalItems + m.PageSize - 1) / m.PageSize
			statusInfo += fmt.Sprintf(" | Page %d/%d", m.CurrentPage+1, totalPages)
		}
	}

	return statusStyle.Render(statusInfo)
}

// renderContextualControls renders context-aware controls based on current tab and mode
func (m *Model) renderContextualControls() string {
	var s strings.Builder

	// Create control styles
	navigationStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")). // Green
		Bold(true)

	actionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFF00")). // Yellow
		Bold(true)

	destructiveStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")). // Red
		Bold(true)

	quickAddStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FFFF")). // Cyan
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888"))

	// NAVIGATION (always shown)
	s.WriteString(navigationStyle.Render("NAVIGATION:  ↑/↓ Navigate  [Tab] Switch  ←/→ Pages  [Enter] Select  [Esc] Cancel") + "\n")

	// Blank line for separation
	s.WriteString("\n")

	// ACTIONS (context-specific)
	switch m.ActiveTab {
	case TabWeeklySummaries:
		s.WriteString(actionStyle.Render("ACTIONS:     ←/→ Switch weeks") + "\n")
	case TabTrips:
		s.WriteString(actionStyle.Render("ACTIONS:     [Ctrl+E] Edit  [Ctrl+F] Search  [Ctrl+T] Template") + "\n")
	case TabExpenses:
		s.WriteString(actionStyle.Render("ACTIONS:     [Ctrl+E] Edit  [Ctrl+F] Filter") + "\n")
	case TabTemplates:
		s.WriteString(actionStyle.Render("ACTIONS:     [Ctrl+E] Edit  [Ctrl+F] Search") + "\n")
	}

	// QUICK ADD (context-specific)
	switch m.ActiveTab {
	case TabTrips:
		s.WriteString(quickAddStyle.Render("QUICK ADD:   [Ctrl+X] Expense  [Ctrl+R] Recurring") + "\n")
	case TabExpenses:
		s.WriteString(quickAddStyle.Render("QUICK ADD:   [Ctrl+X] Add expense") + "\n")
	case TabTemplates:
		s.WriteString(quickAddStyle.Render("QUICK ADD:   [Ctrl+T] New template  [U] Use template") + "\n")
	}

	// DELETE (context-specific)
	switch m.ActiveTab {
	case TabTrips, TabExpenses, TabTemplates:
		s.WriteString(destructiveStyle.Render("DELETE:      [Ctrl+D] Delete selected") + "\n")
	}

	// Show status information
	if m.SearchMode {
		s.WriteString(normalStyle.Render("🔍 Search mode active") + "\n")
	}

	// Add help hints
	s.WriteString("\n")
	s.WriteString(normalStyle.Render("HELP:        [F1] Quick  [F2] Detailed  [F3] Advanced") + "\n")

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
