package ui

import (
	"strings"

	"github.com/akyriako/o7k/internal/resource"
	"github.com/charmbracelet/lipgloss"
)

const logo = `  ___   _____  _    
 / _ \ |___  || | __
| | | |   / / | |/ /
| |_| |  / /  |   < 
 \___/  /_/   |_|\_\`

func (m Model) renderHeader() string {
	profile := strings.Join([]string{
		headerLabelStyle.Render("Profile: ") +
			headerValueStyle.Render(m.profile),

		headerLabelStyle.Render("Region:  ") +
			headerValueStyle.Render(m.region),

		headerLabelStyle.Render("Project: ") +
			headerValueStyle.Render(m.project),

		headerLabelStyle.Render("Domain:  ") +
			headerValueStyle.Render(m.domain),
	}, "\n")

	commands := m.renderCommands()

	renderedLogo := logoStyle.Render(logo)

	const (
		profileWidth = 32
		logoWidth    = 23
	)

	centerWidth := max(
		m.width-profileWidth-logoWidth,
		1,
	)

	left := lipgloss.NewStyle().
		Width(profileWidth).
		Render(profile)

	center := lipgloss.NewStyle().
		Width(centerWidth).
		Align(lipgloss.Center).
		Render(commands)

	right := lipgloss.NewStyle().
		Width(logoWidth).
		Align(lipgloss.Right).
		Render(renderedLogo)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		left,
		center,
		right,
	)
}

func renderHeaderCommand(key, description string) string {
	return headerCommandKeyStyle.Render(key) +
		" " +
		headerCommandTextStyle.Render(description)
}

func (m Model) renderCommands() string {
	commands := []resource.Command{
		{Key: "<q>", Description: "Quit"},
		{Key: "<:>", Description: "Resource"},
		{Key: "<r>", Description: "Refresh"},
		{Key: "<Esc>", Description: "Dismiss"},
	}

	if m.resource != nil {
		commands = append(commands, m.resource.Commands()...)
	}

	lines := make([]string, 0, len(commands))

	for _, command := range commands {
		lines = append(
			lines,
			renderHeaderCommand(
				command.Key,
				command.Description,
			),
		)
	}

	return strings.Join(lines, "\n")
}
