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
	Mode            string // "date", "origin", "destination", "type", "edit", "delete", "delete_confirm", "expense_date", "expense_amount", "expense_description", "expense_edit", "expense_delete_confirm", "search"
	Err             error
	Storage         storage.Storage
	RatePerMile     float64
	MapsClient      maps.DistanceCalculator
	Data            *model.StorageData
	EditIndex       int    // Index of trip being edited
	SelectedTrip    int    // Index of selected trip for operations
	SelectedExpense int    // Index of selected expense for operations
	SearchQuery     string // Current search query
	SearchMode      bool   // Whether we're in search mode
	CurrentGroup    int    // Index of current time group being viewed
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
			if m.Mode == "search" {
				// Exit search mode
				m.Mode = "date"
				m.SearchMode = false
				m.SearchQuery = ""
				m.TextInput.Reset()
				m.TextInput.Placeholder = "Enter date (YYYY-MM-DD)..."
				return m, cmd
			}
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
		case tea.KeyCtrlF:
			// Enter search mode
			if m.Mode == "date" {
				m.Mode = "search"
				m.SearchMode = true
				m.TextInput.Reset()
				m.TextInput.Placeholder = "Search trips and expenses..."
			}
		case tea.KeyCtrlT:
			// Jump to Today's group
			if m.Mode == "date" {
				groups := m.groupByTimePeriod()
				for i, group := range groups {
					if group.Title == "Today" {
						m.CurrentGroup = i
						break
					}
				}
			}
		case tea.KeyCtrlW:
			// Jump to This Week's group
			if m.Mode == "date" {
				groups := m.groupByTimePeriod()
				for i, group := range groups {
					if group.Title == "This Week" {
						m.CurrentGroup = i
						break
					}
				}
			}
		case tea.KeyCtrlY:
			// Jump to This Month's group
			if m.Mode == "date" {
				groups := m.groupByTimePeriod()
				for i, group := range groups {
					if group.Title == "This Month" {
						m.CurrentGroup = i
						break
					}
				}
			}
		case tea.KeyCtrlO:
			// Jump to Older group
			if m.Mode == "date" {
				groups := m.groupByTimePeriod()
				for i, group := range groups {
					if group.Title == "Older" {
						m.CurrentGroup = i
						break
					}
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

		// Handle search input
		if m.Mode == "search" {
			m.SearchQuery = m.TextInput.Value()
		}
	}

	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

// TimeGroup represents a group of trips or expenses for a specific time period
type TimeGroup struct {
	Title     string
	Trips     []model.Trip
	Expenses  []model.Expense
	StartDate string
	EndDate   string
}

// filterBySearch filters trips and expenses based on the search query
func (m *Model) filterBySearch(trips []model.Trip, expenses []model.Expense) ([]model.Trip, []model.Expense) {
	if m.SearchQuery == "" {
		return trips, expenses
	}

	query := strings.ToLower(m.SearchQuery)
	var filteredTrips []model.Trip
	var filteredExpenses []model.Expense

	// Filter trips
	for _, trip := range trips {
		if strings.Contains(strings.ToLower(trip.Date), query) ||
			strings.Contains(strings.ToLower(trip.Origin), query) ||
			strings.Contains(strings.ToLower(trip.Destination), query) ||
			strings.Contains(strings.ToLower(trip.Type), query) {
			filteredTrips = append(filteredTrips, trip)
		}
	}

	// Filter expenses
	for _, expense := range expenses {
		if strings.Contains(strings.ToLower(expense.Date), query) ||
			strings.Contains(strings.ToLower(expense.Description), query) ||
			strings.Contains(strings.ToLower(fmt.Sprintf("%.2f", expense.Amount)), query) {
			filteredExpenses = append(filteredExpenses, expense)
		}
	}

	return filteredTrips, filteredExpenses
}

// groupByTimePeriod groups trips and expenses by time periods (Today, This Week, This Month, Older)
func (m *Model) groupByTimePeriod() []TimeGroup {
	now := time.Now()
	today := now.Format("2006-01-02")
	weekStart := now.AddDate(0, 0, -int(now.Weekday())).Format("2006-01-02")
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")

	groups := []TimeGroup{
		{Title: "Today", StartDate: today, EndDate: today},
		{Title: "This Week", StartDate: weekStart, EndDate: today},
		{Title: "This Month", StartDate: monthStart, EndDate: today},
		{Title: "Older", StartDate: "0000-01-01", EndDate: monthStart},
	}

	// Get filtered trips and expenses if in search mode
	trips := m.Trips
	expenses := m.Data.Expenses
	if m.SearchMode {
		trips, expenses = m.filterBySearch(trips, expenses)
	}

	// Sort trips and expenses by date
	sort.Slice(trips, func(i, j int) bool {
		return trips[i].Date > trips[j].Date // Most recent first
	})
	sort.Slice(expenses, func(i, j int) bool {
		return expenses[i].Date > expenses[j].Date // Most recent first
	})

	// Group trips
	for _, trip := range trips {
		for i := range groups {
			if trip.Date >= groups[i].StartDate && trip.Date <= groups[i].EndDate {
				groups[i].Trips = append(groups[i].Trips, trip)
				break
			}
		}
	}

	// Group expenses
	for _, expense := range expenses {
		for i := range groups {
			if expense.Date >= groups[i].StartDate && expense.Date <= groups[i].EndDate {
				groups[i].Expenses = append(groups[i].Expenses, expense)
				break
			}
		}
	}

	// Remove empty groups
	var nonEmptyGroups []TimeGroup
	for _, group := range groups {
		if len(group.Trips) > 0 || len(group.Expenses) > 0 {
			nonEmptyGroups = append(nonEmptyGroups, group)
		}
	}

	return nonEmptyGroups
}

// renderTimeGroup renders a time group with its trips and expenses
func (m *Model) renderTimeGroup(group TimeGroup) string {
	var s strings.Builder

	// Group header with summary
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF5F87")).
		Padding(0, 1)

	summaryStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Padding(0, 1)

	// Calculate group totals
	totalMiles := model.CalculateTotalMiles(group.Trips)
	totalReimbursement := model.CalculateReimbursement(group.Trips, m.RatePerMile)
	totalExpenses := model.CalculateTotalExpenses(group.Expenses)

	s.WriteString(headerStyle.Render(group.Title) + "\n")
	s.WriteString(summaryStyle.Render(fmt.Sprintf("Total Miles: %.2f | Mileage Amount: $%.2f | Expenses: $%.2f\n",
		totalMiles, totalReimbursement, totalExpenses)))

	// Render trips
	if len(group.Trips) > 0 {
		s.WriteString("\nTrips:\n")
		for i, t := range group.Trips {
			style := lipgloss.NewStyle()
			if i == m.SelectedTrip {
				style = style.Foreground(lipgloss.Color("#FF5F87"))
			}
			// Calculate display miles based on trip type
			displayMiles := t.Miles
			if t.Type == "round" {
				displayMiles = t.Miles * 2
			}
			tripType := "→"
			if t.Type == "round" {
				tripType = "↔"
			}
			s.WriteString(style.Render(fmt.Sprintf("  %s - %s %s %s (%.2f miles)\n",
				t.Date, t.Origin, tripType, t.Destination, displayMiles)))
		}
	}

	// Render expenses
	if len(group.Expenses) > 0 {
		s.WriteString("\nExpenses:\n")
		for i, e := range group.Expenses {
			style := lipgloss.NewStyle()
			if i == m.SelectedExpense {
				style = style.Foreground(lipgloss.Color("#FF5F87"))
			}
			s.WriteString(style.Render(fmt.Sprintf("  %s - $%.2f - %s\n",
				e.Date, e.Amount, e.Description)))
		}
	}

	s.WriteString("\n" + strings.Repeat("─", 50) + "\n")
	return s.String()
}

func (m *Model) View() string {
	var s strings.Builder

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF5F87")).
		Render("Nanny Tracker")
	s.WriteString(title + "\n\n")

	// Current mode and search status
	modeText := fmt.Sprintf("Current mode: %s", m.Mode)
	if m.SearchMode {
		modeText += " (Search: " + m.SearchQuery + ")"
	}
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

	// Group and display trips and expenses by time period
	groups := m.groupByTimePeriod()
	for _, group := range groups {
		s.WriteString(m.renderTimeGroup(group))
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
		Render("\nPress Ctrl+C to quit | Ctrl+E to edit | Ctrl+D to delete | Ctrl+X for expenses | " +
			"Ctrl+F to search | Ctrl+T/W/Y/O to jump to time periods | ↑/↓ to select | Tab to switch between trips/expenses")
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
