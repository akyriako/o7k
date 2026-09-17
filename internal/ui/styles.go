package ui

import "github.com/charmbracelet/lipgloss"

var (
	resourceTagStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("255")).
				Background(lipgloss.Color("#ED1944")).
				Bold(true)

	commandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#87CEFA"))

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#FFA500")).
			Align(lipgloss.Center)

	tableContainerStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("#00BFFF"))
)
