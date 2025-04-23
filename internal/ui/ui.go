package ui

import (
	"context"
	"fmt"
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
	TextInput   textinput.Model
	Trips       []model.Trip
	CurrentTrip model.Trip
	Mode        string // "origin", "destination", or "date"
	Err         error
	Storage     storage.Storage
	RatePerMile float64
	MapsClient  maps.DistanceCalculator
}

// New creates a new UI model with a mock maps client (for backward compatibility)
func New(storage storage.Storage, ratePerMile float64) (*Model, error) {
	trips, err := storage.LoadTrips()
	if err != nil {
		return nil, fmt.Errorf("failed to load trips: %w", err)
	}

	ti := textinput.New()
	ti.Placeholder = "Enter address..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50

	return &Model{
		TextInput:   ti,
		Trips:       trips,
		CurrentTrip: model.Trip{},
		Mode:        "origin",
		Storage:     storage,
		RatePerMile: ratePerMile,
		MapsClient:  maps.NewMockClient(), // Use mock client by default
	}, nil
}

// NewWithClient creates a new UI model with a provided maps client (useful for testing)
func NewWithClient(storage storage.Storage, ratePerMile float64, mapsClient maps.DistanceCalculator) (*Model, error) {
	trips, err := storage.LoadTrips()
	if err != nil {
		return nil, fmt.Errorf("failed to load trips: %w", err)
	}

	ti := textinput.New()
	ti.Placeholder = "Enter address..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50

	return &Model{
		TextInput:   ti,
		Trips:       trips,
		CurrentTrip: model.Trip{},
		Mode:        "origin",
		Storage:     storage,
		RatePerMile: ratePerMile,
		MapsClient:  mapsClient,
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
			if m.Mode == "origin" {
				m.CurrentTrip.Origin = m.TextInput.Value()
				m.TextInput.Reset()
				m.Mode = "destination"
				m.TextInput.Placeholder = "Enter destination address..."
			} else if m.Mode == "destination" {
				m.CurrentTrip.Destination = m.TextInput.Value()
				m.TextInput.Reset()
				m.Mode = "date"
				m.TextInput.Placeholder = "Enter date (YYYY-MM-DD)..."
			} else if m.Mode == "date" {
				m.CurrentTrip.Date = m.TextInput.Value()

				// Calculate distance using Google Maps API
				distance, err := m.MapsClient.CalculateDistance(context.Background(), m.CurrentTrip.Origin, m.CurrentTrip.Destination)
				if err != nil {
					m.Err = fmt.Errorf("failed to calculate distance: %w", err)
					return m, cmd
				}
				m.CurrentTrip.Miles = distance

				m.Trips = append(m.Trips, m.CurrentTrip)
				if err := m.Storage.SaveTrips(m.Trips); err != nil {
					m.Err = err
				}
				m.CurrentTrip = model.Trip{}
				m.TextInput.Reset()
				m.Mode = "origin"
				m.TextInput.Placeholder = "Enter origin address..."
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	var s strings.Builder

	// Title with custom style
	title := titleStyle.Render("Nanny Mileage Tracker")
	s.WriteString(title + "\n\n")

	// Mode indicator with custom style
	modeText := modeStyle.Render(fmt.Sprintf("Current mode: %s", m.Mode))
	s.WriteString(modeText + "\n\n")

	// Input field with custom style
	input := inputStyle.Render(m.TextInput.View())
	s.WriteString(input + "\n\n")

	// Current trip info
	if m.CurrentTrip.Origin != "" {
		s.WriteString(tripStyle.Render(fmt.Sprintf("Origin: %s\n", m.CurrentTrip.Origin)))
	}
	if m.CurrentTrip.Destination != "" {
		s.WriteString(tripStyle.Render(fmt.Sprintf("Destination: %s\n", m.CurrentTrip.Destination)))
	}
	if m.CurrentTrip.Date != "" {
		s.WriteString(tripStyle.Render(fmt.Sprintf("Date: %s\n", m.CurrentTrip.Date)))
	}

	// Trip history
	if len(m.Trips) > 0 {
		s.WriteString("\n" + titleStyle.Render("Trip History") + "\n")

		// Create a list container
		listContainer := lipgloss.NewStyle().
			PaddingLeft(2).
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(accentColor)

		var tripList strings.Builder
		for i, t := range m.Trips {
			// Format each trip with consistent spacing
			trip := fmt.Sprintf("%d. %s → %s (%.2f miles) - %s",
				i+1, t.Origin, t.Destination, t.Miles, t.Date)
			tripList.WriteString(tripStyle.Render(trip) + "\n")
		}

		// Add the formatted list to the main view
		s.WriteString(listContainer.Render(tripList.String()))
	}

	// Total mileage and reimbursement
	if len(m.Trips) > 0 {
		totalMiles := model.CalculateTotalMiles(m.Trips)
		totalReimbursement := model.CalculateReimbursement(m.Trips, m.RatePerMile)

		// Add extra spacing before totals section
		s.WriteString("\n\n" + titleStyle.Render("Totals") + "\n")

		// Create a totals container with similar styling to trip list
		totalsContainer := lipgloss.NewStyle().
			PaddingLeft(2).
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(accentColor)

		var totalsList strings.Builder
		// Format totals to match trip history style with consistent spacing
		totalsList.WriteString(fmt.Sprintf("1. Total Miles: %.2f miles\n", totalMiles))
		totalsList.WriteString(fmt.Sprintf("2. Total Reimbursement: $%.2f", totalReimbursement))
		s.WriteString(totalsContainer.Render(tripStyle.Render(totalsList.String())))
	}

	// Error message
	if m.Err != nil {
		errorMsg := errorStyle.Render(fmt.Sprintf("\n\nError: %v", m.Err))
		s.WriteString(errorMsg)
	}

	// Help text
	help := helpStyle.Render("\n\nPress Ctrl+C to quit")
	s.WriteString(help)

	// Wrap everything in a container
	return containerStyle.Render(s.String())
}

// CalculateTotalMiles calculates the total miles for a list of trips
func (m *Model) CalculateTotalMiles(trips []model.Trip) float64 {
	return model.CalculateTotalMiles(trips)
}

// CalculateReimbursement calculates the total reimbursement for a list of trips
func (m *Model) CalculateReimbursement(trips []model.Trip, ratePerMile float64) float64 {
	return model.CalculateReimbursement(trips, ratePerMile)
}

// AddTrip adds a new trip to the model's trips list
func (m *Model) AddTrip(trip model.Trip) {
	m.Trips = append(m.Trips, trip)
}
