package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lauren/nannytracker/internal/model"
	"github.com/lauren/nannytracker/internal/storage"
)

// Model represents the UI state
type Model struct {
	TextInput   textinput.Model
	Trips       []model.Trip
	CurrentTrip model.Trip
	Mode        string
	Err         error
	Storage     storage.Storage
	RatePerMile float64
}

// New creates a new UI model
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
				m.CurrentTrip.Miles = 10.0 // Placeholder value
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

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF5F87")).
		Render("Nanny Mileage Tracker")
	s.WriteString(title + "\n\n")

	// Current mode
	modeText := fmt.Sprintf("Current mode: %s", m.Mode)
	s.WriteString(modeText + "\n\n")

	// Input field
	s.WriteString(m.TextInput.View() + "\n\n")

	// Current trip info
	if m.CurrentTrip.Origin != "" {
		s.WriteString(fmt.Sprintf("Origin: %s\n", m.CurrentTrip.Origin))
	}
	if m.CurrentTrip.Destination != "" {
		s.WriteString(fmt.Sprintf("Destination: %s\n", m.CurrentTrip.Destination))
	}

	// Trip history
	if len(m.Trips) > 0 {
		s.WriteString("\nTrip History:\n")
		for i, t := range m.Trips {
			s.WriteString(fmt.Sprintf("%d. %s → %s (%.2f miles)\n", i+1, t.Origin, t.Destination, t.Miles))
		}
	}

	// Total mileage and reimbursement
	if len(m.Trips) > 0 {
		totalMiles := model.CalculateTotalMiles(m.Trips)
		totalReimbursement := model.CalculateReimbursement(m.Trips, m.RatePerMile)
		s.WriteString(fmt.Sprintf("\nTotal Miles: %.2f\n", totalMiles))
		s.WriteString(fmt.Sprintf("Total Reimbursement: $%.2f\n", totalReimbursement))
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
		Render("\nPress Ctrl+C to quit")
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

// AddTrip adds a new trip to the model's trips list
func (m *Model) AddTrip(trip model.Trip) {
	m.Trips = append(m.Trips, trip)
}
