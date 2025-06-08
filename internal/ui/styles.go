package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Main styles
	mainStyle = lipgloss.NewStyle().
			Margin(1, 2)

	// Title styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7B2CBF")).
			Padding(0, 1).
			MarginBottom(1)

	// List styles
	listStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#666666")).
			Padding(1, 2)
		// Width(StandardWidth)

	noteStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 2)

	selectedNoteStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#000000")).
				Background(lipgloss.Color("#9D4EDD")).
				Padding(0, 2)

	// Preview styles
	previewStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#666666")).
			Padding(1, 2).
		// Width(StandardWidth).
		Height(6)

	previewTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#666666")).
				Bold(true).
				MarginTop(1)

	// Input styles
	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#666666")).
			Padding(0, 1)
		// Width(StandardWidth)

	textareaStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#666666")).
			Padding(0, 1).
			Width(StandardWidth).
		// Width( - StandardTextInputPadding).
		Height(StandardHeight)

	// Utility styles
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000"))

	tagStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5af78e")).
			Italic(true)
)
