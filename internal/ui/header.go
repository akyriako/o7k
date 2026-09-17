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
	const (
		contextWidth = 60
		logoWidth    = 23
	)

	profile := m.renderContext(contextWidth)
	commands := m.renderCommands()
	renderedLogo := logoStyle.Render(logo)

	commandWidth := max(m.width-contextWidth-logoWidth, 1)

	left := lipgloss.NewStyle().Width(contextWidth).Render(profile)
	center := lipgloss.NewStyle().Width(commandWidth).Align(lipgloss.Left).Render(commands)
	right := lipgloss.NewStyle().Width(logoWidth).Align(lipgloss.Right).Render(renderedLogo)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, center, right)
}

func (m Model) renderContext(width int) string {
	const labelWidth = 10

	valueWidth := max(width-labelWidth, 1)

	field := func(label, value string) string {
		value = truncate(value, valueWidth)

		return headerLabelStyle.Width(labelWidth).Render(label) +
			headerValueStyle.Render(value)
	}

	return strings.Join([]string{
		field("Cloud:", m.context.Cloud),
		field("Identity:", m.context.Identity),
		field("Region:", m.context.Region),
		field("Domain:", m.context.Domain),
		field("Project:", m.context.Project),
	}, "\n")
}

func (m Model) renderCommands() string {
	commands := []resource.Command{
		{Key: "q", Description: "Quit"},
		{Key: ":", Description: "Resource"},
		{Key: "r", Description: "Refresh"},
		{Key: "esc", Description: "Dismiss/Return"},
	}

	if m.resource != nil {
		commands = append(commands, m.resource.Commands()...)
	}

	lines := make([]string, 0, len(commands))

	for _, command := range commands {
		key := headerCommandKeyStyle.Width(12).Render("<" + command.Key + ">")
		text := headerCommandTextStyle.Render(command.Description)
		lines = append(lines, key+text)
	}

	return strings.Join(lines, "\n")
}

func truncate(value string, width int) string {
	if lipgloss.Width(value) <= width {
		return value
	}

	if width <= 3 {
		return strings.Repeat(".", width)
	}

	runes := []rune(value)

	for len(runes) > 0 && lipgloss.Width(string(runes)) > width-3 {
		runes = runes[:len(runes)-1]
	}

	return string(runes) + "..."
}
