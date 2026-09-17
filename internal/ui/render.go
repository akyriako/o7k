package ui

import (
	"fmt"
	"strings"
	"time"

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

	profile := m.renderHeaderContext(contextWidth)
	commands := m.renderHeaderCommands()
	renderedLogo := logoStyle.Render(logo)

	commandWidth := max(m.width-contextWidth-logoWidth, 1)

	left := lipgloss.NewStyle().Width(contextWidth).Render(profile)
	center := lipgloss.NewStyle().Width(commandWidth).Align(lipgloss.Left).Render(commands)
	right := lipgloss.NewStyle().Width(logoWidth).Align(lipgloss.Right).Render(renderedLogo)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, center, right)
}

func (m Model) renderHeaderContext(width int) string {
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

func (m Model) renderHeaderCommands() string {
	commands := []resource.Command{
		{Key: "q", Description: "Quit"},
		{Key: ":", Description: "Resource"},
		{Key: "r", Description: "Refresh"},
		{Key: "c", Description: "Copy"},
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

func (m Model) renderTable() string {
	title := ""

	if m.resource != nil {
		title = fmt.Sprintf(
			" %s[%d] | last update: %s ",
			m.resource.Kind(),
			m.itemCount,
			time.Now().Format("02.01.2006 15:04:05"),
		)
	}

	innerWidth := max(m.width-2, 1)

	titleWidth := lipgloss.Width(title)
	remaining := max(innerWidth-titleWidth, 0)

	left := remaining / 2
	right := remaining - left

	borderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00BFFF"))

	top := borderStyle.Render(
		"┌" +
			repeat("─", left) +
			title +
			repeat("─", right) +
			"┐",
	)

	bottom := borderStyle.Render(
		"└" +
			repeat("─", innerWidth) +
			"┘",
	)

	body := lipgloss.NewStyle().
		Width(innerWidth).
		Render(m.table.View())

	if m.loaded && !m.loading && m.itemCount == 0 {
		lines := strings.Split(body, "\n")

		if len(lines) > 1 {
			header := lines[0]
			contentHeight := len(lines) - 1

			emptyContent := lipgloss.Place(
				innerWidth,
				contentHeight,
				lipgloss.Center,
				lipgloss.Center,
				"No resources found",
			)

			body = header + "\n" + emptyContent
		}
	}

	bodyLines := strings.Split(body, "\n")

	for i, line := range bodyLines {
		lineWidth := lipgloss.Width(line)

		if lineWidth < innerWidth {
			line += strings.Repeat(" ", innerWidth-lineWidth)
		}

		bodyLines[i] =
			borderStyle.Render("│") +
				line +
				borderStyle.Render("│")
	}

	return top +
		"\n" +
		strings.Join(bodyLines, "\n") +
		"\n" +
		bottom
}

func (m Model) renderDetails() string {
	title := " server details "

	if m.detailID != "" {
		title = " server " + m.detailID + " "
	}

	innerWidth := max(m.width-2, 1)

	titleWidth := lipgloss.Width(title)
	remaining := max(innerWidth-titleWidth, 0)

	left := remaining / 2
	right := remaining - left

	borderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00BFFF"))

	top := borderStyle.Render(
		"┌" +
			repeat("─", left) +
			title +
			repeat("─", right) +
			"┐",
	)

	bottom := borderStyle.Render(
		"└" +
			repeat("─", innerWidth) +
			"┘",
	)

	body := m.detail.View()
	bodyLines := strings.Split(body, "\n")

	for i, line := range bodyLines {
		lineWidth := lipgloss.Width(line)

		if lineWidth < innerWidth {
			line += strings.Repeat(" ", innerWidth-lineWidth)
		}

		bodyLines[i] =
			borderStyle.Render("│") +
				line +
				borderStyle.Render("│")
	}

	return top +
		"\n" +
		strings.Join(bodyLines, "\n") +
		"\n" +
		bottom
}
