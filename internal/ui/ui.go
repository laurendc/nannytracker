package ui

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lauren/nannytracker/internal/maps"
	"github.com/lauren/nannytracker/internal/model"
	"github.com/lauren/nannytracker/internal/storage"
)

// Model represents the UI state
type Model struct {
	TextInput       textinput.Model
	Trips           []model.Trip
	CurrentTrip     model.Trip
	CurrentExpense  model.Expense
	Mode            string // "date", "origin", "destination", "type", "edit", "delete", "delete_confirm", "expense_date", "expense_amount", "expense_description", "expense_edit", "expense_delete_confirm"
	Err             error
	Storage         storage.Storage
	RatePerMile     float64
	MapsClient      maps.DistanceCalculator
	Data            *model.StorageData
	EditIndex       int // Index of trip being edited
	SelectedTrip    int // Index of selected trip for operations
	SelectedExpense int // Index of selected expense for operations
}

// New creates a new UI model with a mock maps client (for backward compatibility)
func New(storage storage.Storage, ratePerMile float64) (*Model, error) {
	data, err := storage.LoadData()
	if err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
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

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.Mode == "date" || m.Mode == "edit" {
				if m.TextInput.Value() == "" && m.EditIndex >= 0 {
					// Keep existing value if no new input
					m.TextInput.Reset()
					m.Mode = "origin"
					m.TextInput.Placeholder = "Edit origin address..."
				} else {
					m.CurrentTrip.Date = m.TextInput.Value()
					m.TextInput.Reset()
					m.Mode = "origin"
					if m.EditIndex >= 0 {
						m.TextInput.Placeholder = "Edit origin address..."
					} else {
						m.TextInput.Placeholder = "Enter origin address..."
					}
				}
			} else if m.Mode == "origin" {
				if m.TextInput.Value() == "" && m.EditIndex >= 0 {
					// Keep existing value if no new input
					m.TextInput.Reset()
					m.Mode = "destination"
					m.TextInput.Placeholder = "Edit destination address..."
				} else {
					m.CurrentTrip.Origin = m.TextInput.Value()
					m.TextInput.Reset()
					m.Mode = "destination"
					if m.EditIndex >= 0 {
						m.TextInput.Placeholder = "Edit destination address..."
					} else {
						m.TextInput.Placeholder = "Enter destination address..."
					}
				}
			} else if m.Mode == "destination" {
				if m.TextInput.Value() == "" && m.EditIndex >= 0 {
					// Keep existing value if no new input
					m.TextInput.Reset()
					m.Mode = "type"
					m.TextInput.Placeholder = "Enter trip type (single/round)..."
				} else {
					m.CurrentTrip.Destination = m.TextInput.Value()

					// Calculate distance using Google Maps API
					distance, err := m.MapsClient.CalculateDistance(context.Background(), m.CurrentTrip.Origin, m.CurrentTrip.Destination)
					if err != nil {
						m.Err = fmt.Errorf("failed to calculate distance: %w", err)
						return m, cmd
					}
					m.CurrentTrip.Miles = distance

					// Move to trip type selection
					m.Mode = "type"
					m.TextInput.Reset()
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
					m.CurrentTrip.Type = tripType
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
				m.CurrentTrip = model.Trip{}
				m.Mode = "date"
				m.EditIndex = -1
				m.TextInput.Reset()
				m.TextInput.Placeholder = "Enter date (YYYY-MM-DD)..."
			} else if m.Mode == "expense_date" {
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
				amount, err := strconv.ParseFloat(m.TextInput.Value(), 64)
				if err != nil {
					m.Err = fmt.Errorf("invalid amount: %w", err)
					return m, cmd
				}
				// Create a temporary expense to validate the amount
				tempExpense := model.Expense{
					Date:        m.CurrentExpense.Date,
					Amount:      amount,
					Description: "temp", // Dummy value for validation
				}
				if err := tempExpense.Validate(); err != nil {
					m.Err = err
					return m, cmd
				}
				m.CurrentExpense.Amount = amount
				m.TextInput.Reset()
				m.Mode = "expense_description"
				m.TextInput.Placeholder = "Enter expense description..."
			} else if m.Mode == "expense_description" {
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
			} else if m.Mode == "expense_edit" {
				if m.TextInput.Value() == "" && m.EditIndex >= 0 {
					// Keep existing value if no new input
					m.TextInput.Reset()
					m.Mode = "expense_amount"
					m.TextInput.Placeholder = "Edit expense amount..."
				} else {
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
					m.TextInput.Placeholder = "Edit expense amount..."
				}
			} else if m.Mode == "expense_amount" && m.EditIndex >= 0 {
				if m.TextInput.Value() == "" {
					// Keep existing value if no new input
					m.TextInput.Reset()
					m.Mode = "expense_description"
					m.TextInput.Placeholder = "Edit expense description..."
				} else {
					amount, err := strconv.ParseFloat(m.TextInput.Value(), 64)
					if err != nil {
						m.Err = fmt.Errorf("invalid amount: %w", err)
						return m, cmd
					}
					// Create a temporary expense to validate the amount
					tempExpense := model.Expense{
						Date:        m.CurrentExpense.Date,
						Amount:      amount,
						Description: "temp", // Dummy value for validation
					}
					if err := tempExpense.Validate(); err != nil {
						m.Err = err
						return m, cmd
					}
					m.CurrentExpense.Amount = amount
					m.TextInput.Reset()
					m.Mode = "expense_description"
					m.TextInput.Placeholder = "Edit expense description..."
				}
			} else if m.Mode == "expense_description" && m.EditIndex >= 0 {
				if m.TextInput.Value() == "" {
					// Keep existing value if no new input
					if m.CurrentExpense.Description == "" {
						m.Err = fmt.Errorf("description cannot be empty")
						return m, cmd
					}
				} else {
					m.CurrentExpense.Description = m.TextInput.Value()
				}

				// Validate the expense before saving
				if err := m.CurrentExpense.Validate(); err != nil {
					m.Err = fmt.Errorf("invalid expense: %w", err)
					return m, cmd
				}

				// Update existing expense
				if err := m.Data.EditExpense(m.EditIndex, m.CurrentExpense); err != nil {
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
				m.EditIndex = -1
				m.TextInput.Reset()
				m.TextInput.Placeholder = "Enter date (YYYY-MM-DD)..."
			} else if m.Mode == "expense_delete_confirm" {
				if strings.ToLower(m.TextInput.Value()) == "yes" {
					// Delete the expense
					if err := m.Data.DeleteExpense(m.SelectedExpense); err != nil {
						m.Err = err
						return m, cmd
					}
					model.CalculateAndUpdateWeeklySummaries(m.Data, m.RatePerMile)
					if err := m.Storage.SaveData(m.Data); err != nil {
						m.Err = err
					}
					if m.SelectedExpense >= len(m.Data.Expenses) {
						m.SelectedExpense = len(m.Data.Expenses) - 1
					}
					// Reset mode and input
					m.Mode = "date"
					m.TextInput.Reset()
					m.TextInput.Placeholder = "Enter date (YYYY-MM-DD)..."
				} else {
					// Cancel deletion
					m.Mode = "date"
					m.TextInput.Reset()
					m.TextInput.Placeholder = "Enter date (YYYY-MM-DD)..."
				}
			} else if m.Mode == "delete_confirm" {
				if strings.ToLower(m.TextInput.Value()) == "yes" {
					// Delete the trip
					if err := m.Data.DeleteTrip(m.SelectedTrip); err != nil {
						m.Err = err
						return m, cmd
					}
					model.CalculateAndUpdateWeeklySummaries(m.Data, m.RatePerMile)
					if err := m.Storage.SaveData(m.Data); err != nil {
						m.Err = err
					}
					m.Trips = m.Data.Trips
					if m.SelectedTrip >= len(m.Trips) {
						m.SelectedTrip = len(m.Trips) - 1
					}
					// Reset mode and input
					m.Mode = "date"
					m.TextInput.Reset()
					m.TextInput.Placeholder = "Enter date (YYYY-MM-DD)..."
				} else {
					// Cancel deletion
					m.Mode = "date"
					m.TextInput.Reset()
					m.TextInput.Placeholder = "Enter date (YYYY-MM-DD)..."
				}
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyCtrlE:
			// Enter edit mode
			if m.Mode == "date" {
				if m.SelectedTrip >= 0 && m.SelectedTrip < len(m.Trips) {
					m.Mode = "edit"
					m.EditIndex = m.SelectedTrip
					m.CurrentTrip = m.Trips[m.SelectedTrip]
					m.TextInput.SetValue(m.CurrentTrip.Date)
					m.TextInput.Placeholder = "Edit date (YYYY-MM-DD)..."
				} else if m.SelectedExpense >= 0 && m.SelectedExpense < len(m.Data.Expenses) {
					m.Mode = "expense_edit"
					m.EditIndex = m.SelectedExpense
					m.CurrentExpense = m.Data.Expenses[m.SelectedExpense]
					m.TextInput.SetValue(m.CurrentExpense.Date)
					m.TextInput.Placeholder = "Edit date (YYYY-MM-DD)..."
				}
			}
		case tea.KeyCtrlX:
			// Enter expense mode
			m.Mode = "expense_date"
			m.TextInput.Reset()
			m.TextInput.Placeholder = "Enter expense date (YYYY-MM-DD)..."
		case tea.KeyCtrlD:
			// Enter delete confirmation mode
			if m.Mode == "date" {
				if m.SelectedTrip >= 0 && m.SelectedTrip < len(m.Trips) {
					m.Mode = "delete_confirm"
					m.TextInput.Reset()
					m.TextInput.Placeholder = "Type 'yes' to confirm deletion..."
				} else if m.SelectedExpense >= 0 && m.SelectedExpense < len(m.Data.Expenses) {
					m.Mode = "expense_delete_confirm"
					m.TextInput.Reset()
					m.TextInput.Placeholder = "Type 'yes' to confirm deletion..."
				}
			}
		case tea.KeyDown:
			// Move selection down
			if len(m.Trips) > 0 {
				if m.SelectedTrip == -1 {
					m.SelectedTrip = 0
					m.SelectedExpense = -1
				} else if m.SelectedTrip < len(m.Trips)-1 {
					m.SelectedTrip++
					m.SelectedExpense = -1
				}
			} else if len(m.Data.Expenses) > 0 {
				if m.SelectedExpense == -1 {
					m.SelectedExpense = 0
					m.SelectedTrip = -1
				} else if m.SelectedExpense < len(m.Data.Expenses)-1 {
					m.SelectedExpense++
					m.SelectedTrip = -1
				}
			}
		case tea.KeyUp:
			// Move selection up
			if m.SelectedTrip > 0 {
				m.SelectedTrip--
				m.SelectedExpense = -1
			} else if m.SelectedExpense > 0 {
				m.SelectedExpense--
				m.SelectedTrip = -1
			}
		case tea.KeyTab:
			// Switch between trips and expenses
			if m.SelectedTrip >= 0 {
				m.SelectedTrip = -1
				if len(m.Data.Expenses) > 0 {
					m.SelectedExpense = 0
				}
			} else {
				m.SelectedExpense = -1
				if len(m.Trips) > 0 {
					m.SelectedTrip = 0
				}
			}
		}
	}

	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	var s strings.Builder

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF5F87")).
		Render("Nanny Tracker")
	s.WriteString(title + "\n\n")

	// Current mode
	modeText := fmt.Sprintf("Current mode: %s", m.Mode)
	s.WriteString(modeText + "\n\n")

	// Input field
	s.WriteString(m.TextInput.View() + "\n\n")

	// Delete confirmation message
	if m.Mode == "delete_confirm" {
		confirmStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Render("WARNING: This will permanently delete the selected trip. Type 'yes' to confirm.")
		s.WriteString(confirmStyle + "\n\n")
	} else if m.Mode == "expense_delete_confirm" {
		confirmStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Render("WARNING: This will permanently delete the selected expense. Type 'yes' to confirm.")
		s.WriteString(confirmStyle + "\n\n")
	}

	// Current trip info
	if m.CurrentTrip.Date != "" {
		s.WriteString(fmt.Sprintf("Date: %s\n", m.CurrentTrip.Date))
	}
	if m.CurrentTrip.Origin != "" {
		s.WriteString(fmt.Sprintf("Origin: %s\n", m.CurrentTrip.Origin))
	}
	if m.CurrentTrip.Destination != "" {
		s.WriteString(fmt.Sprintf("Destination: %s\n", m.CurrentTrip.Destination))
	}
	if m.CurrentTrip.Type != "" {
		s.WriteString(fmt.Sprintf("Type: %s\n", m.CurrentTrip.Type))
	}

	// Current expense info
	if m.Mode == "expense_date" || m.Mode == "expense_amount" || m.Mode == "expense_description" || m.Mode == "expense_edit" {
		s.WriteString("\nEntering Expense:\n")
		if m.CurrentExpense.Date != "" {
			s.WriteString(fmt.Sprintf("Date: %s\n", m.CurrentExpense.Date))
		}
		if m.CurrentExpense.Amount != 0 {
			s.WriteString(fmt.Sprintf("Amount: $%.2f\n", m.CurrentExpense.Amount))
		}
		if m.CurrentExpense.Description != "" {
			s.WriteString(fmt.Sprintf("Description: %s\n", m.CurrentExpense.Description))
		}
	}

	// Weekly summaries
	if len(m.Data.WeeklySummaries) > 0 {
		s.WriteString("\nWeekly Summaries:\n")
		for _, summary := range m.Data.WeeklySummaries {
			s.WriteString(fmt.Sprintf("Week of %s to %s:\n", summary.WeekStart, summary.WeekEnd))
			s.WriteString(fmt.Sprintf("  Total Miles: %.2f\n", summary.TotalMiles))
			s.WriteString(fmt.Sprintf("  Total Mileage Amount: $%.2f\n", summary.TotalAmount))
			s.WriteString(fmt.Sprintf("  Total Expenses: $%.2f\n\n", summary.TotalExpenses))
		}
	}

	// Trip history
	if len(m.Trips) > 0 {
		s.WriteString("\nTrip History:\n")
		for i, t := range m.Trips {
			style := lipgloss.NewStyle()
			if i == m.SelectedTrip {
				style = style.Foreground(lipgloss.Color("#FF5F87"))
			}
			// Calculate display miles based on trip type
			displayMiles := t.Miles
			if t.Type == "round" {
				displayMiles = t.Miles * 2
			}
			s.WriteString(style.Render(fmt.Sprintf("%d. %s - %s → %s (%.2f miles, %s)\n",
				i+1, t.Date, t.Origin, t.Destination, displayMiles, t.Type)))
		}
	}

	// Expense history
	if len(m.Data.Expenses) > 0 {
		s.WriteString("\nExpense History:\n")
		for i, e := range m.Data.Expenses {
			style := lipgloss.NewStyle()
			if i == m.SelectedExpense {
				style = style.Foreground(lipgloss.Color("#FF5F87"))
			}
			s.WriteString(style.Render(fmt.Sprintf("%d. %s - $%.2f - %s\n",
				i+1, e.Date, e.Amount, e.Description)))
		}
	}

	// Total mileage and reimbursement
	if len(m.Trips) > 0 {
		totalMiles := model.CalculateTotalMiles(m.Trips)
		totalReimbursement := model.CalculateReimbursement(m.Trips, m.RatePerMile)
		s.WriteString(fmt.Sprintf("\nTotal Miles: %.2f\n", totalMiles))
		s.WriteString(fmt.Sprintf("Total Mileage Amount: $%.2f\n", totalReimbursement))
	}

	// Total expenses
	if len(m.Data.Expenses) > 0 {
		totalExpenses := model.CalculateTotalExpenses(m.Data.Expenses)
		s.WriteString(fmt.Sprintf("Total Expenses: $%.2f\n", totalExpenses))
	}

	// Error message
	if m.Err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Render(fmt.Sprintf("\nError: %v", m.Err))
		s.WriteString(errorStyle)
	}

	// Help text
	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Render("\nPress Ctrl+C to quit | Ctrl+E to edit | Ctrl+D to delete | Ctrl+X for expenses | ↑/↓ to select | Tab to switch between trips/expenses")
	s.WriteString(help)

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
