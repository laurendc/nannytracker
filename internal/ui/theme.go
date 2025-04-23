package ui

import "github.com/charmbracelet/lipgloss"

// Theme defines the color scheme and styles for the UI
var (
	// Colors
	primaryColor   = lipgloss.Color("#FF5F87")
	secondaryColor = lipgloss.Color("#00FF9F")
	accentColor    = lipgloss.Color("#FFB86C")
	errorColor     = lipgloss.Color("#FF5555")
	successColor   = lipgloss.Color("#50FA7B")
	textColor      = lipgloss.Color("#FFFFFF")
	mutedColor     = lipgloss.Color("#626262")

	// Styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor)

	modeStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true)

	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accentColor).
			Padding(0, 1)

	tripStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Padding(0, 1)

	totalStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(successColor)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	containerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2)
)
