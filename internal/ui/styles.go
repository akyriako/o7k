package ui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var (
	resourceTagStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("255")).
				Background(lipgloss.Color("#ED1944")).
				Bold(true)

	navigationTagStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("0")).
				Background(lipgloss.Color("255"))

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

	loadingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#00BFFF")).
			Bold(true)
)

var (
	jsonKeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#87CEFA"))

	jsonStringStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#98C379"))

	jsonNumberStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E5C07B"))

	jsonBoolStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C678DD"))

	jsonNullStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7F848E"))
)

func tableStyles() table.Styles {
	styles := table.DefaultStyles()

	styles.Selected = styles.Selected.
		Foreground(lipgloss.Color("255")).
		Background(lipgloss.Color("#ED1944")).
		Bold(true)

	return styles
}
