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

	headerLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFA500")).
				Bold(true)

	headerValueStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF"))

	headerCommandKeyStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00BFFF")).
				Bold(true)

	headerCommandTextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF"))

	logoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500")).
			Bold(true)
)
