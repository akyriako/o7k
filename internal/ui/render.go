package ui

import (
	"fmt"
	"strings"

	"github.com/akyriako/o7k/internal/resource"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
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
		field("Region:", m.context.Region),
		field("Domain:", m.context.Domain),
		field("Project:", m.context.Project),
		field("Version:", m.version),
	}, "\n")
}

func (m Model) renderHeaderCommands() string {
	commands := []resource.Command{
		{Key: "ctrl+x", Description: "Quit"},
		{Key: ":", Description: "Resource"},
		{Key: "r", Description: "Refresh"},
		{Key: "c", Description: "Copy"},
		{Key: "esc", Description: "Dismiss/Return"},
	}

	if m.resource != nil {
		commands = append(commands, m.resource.Commands()...)
	}

	const commandsPerColumn = 5
	const commandColumnWidth = 28

	columns := make([]string, 0, (len(commands)+commandsPerColumn-1)/commandsPerColumn)

	for start := 0; start < len(commands); start += commandsPerColumn {
		end := min(start+commandsPerColumn, len(commands))
		lines := make([]string, 0, end-start)

		for _, command := range commands[start:end] {
			key := headerCommandKeyStyle.Width(10).Render("<" + command.Key + ">")
			text := headerCommandTextStyle.Render(command.Description)
			lines = append(lines, key+text)
		}

		columns = append(
			columns,
			lipgloss.NewStyle().
				Width(commandColumnWidth).
				Render(strings.Join(lines, "\n")),
		)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, columns...)
}

func (m Model) renderTable() string {
	title := ""

	if m.resource != nil {
		title = fmt.Sprintf(
			" %s[%d] ",
			m.resource.Kind(),
			m.itemCount,
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

	body := horizontalSlice(
		m.table.View(),
		m.tableXOffset,
		innerWidth,
	)

	if !m.loading && m.itemCount == 0 {
		lines := strings.Split(body, "\n")

		if len(lines) > 1 {
			header := lines[0]
			contentHeight := len(lines) - 1

			emptyMessage := "No resources found"
			if !m.loaded {
				emptyMessage = "Failed to load resources"
			}

			emptyContent := lipgloss.Place(
				innerWidth,
				contentHeight,
				lipgloss.Center,
				lipgloss.Center,
				emptyMessage,
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
	title := " details "

	if m.resource != nil {
		title = fmt.Sprintf(" %s details ", m.resource.Kind())
	}

	if m.resource != nil && m.detailID != "" {
		title = fmt.Sprintf(" %s %s ", m.resource.Kind(), m.detailID)
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

func overlayCenter(background, foreground string) string {
	bg := strings.Split(background, "\n")
	fg := strings.Split(foreground, "\n")

	bgWidth := 0
	for _, line := range bg {
		bgWidth = max(bgWidth, ansi.StringWidth(line))
	}

	fgWidth := 0
	for _, line := range fg {
		fgWidth = max(fgWidth, ansi.StringWidth(line))
	}

	x := max((bgWidth-fgWidth)/2, 0)
	y := max((len(bg)-len(fg))/2, 0)

	for i, fgLine := range fg {
		bgIndex := y + i
		if bgIndex >= len(bg) {
			break
		}

		width := ansi.StringWidth(fgLine)

		left := ansi.Cut(bg[bgIndex], 0, x)
		right := ansi.Cut(bg[bgIndex], x+width, bgWidth)

		bg[bgIndex] = left + fgLine + right
	}

	return strings.Join(bg, "\n")
}

//func (m Model) renderErrorModal() string {
//	content := errorTitleStyle.Render("Error") +
//		"\n\n" +
//		m.err.Error() +
//		"\n\n" +
//		lipgloss.NewStyle().
//			Faint(true).
//			Render("Press Esc or Enter to dismiss")
//
//	return errorModalStyle.Render(content)
//}

func (m Model) renderErrorModal() string {
	hint := lipgloss.NewStyle().
		Faint(true).
		Render("Press ") +
		modalButtonStyle.Render(" Esc ") +
		lipgloss.NewStyle().
			Faint(true).
			Render(" or ") +
		modalButtonStyle.Render(" Enter ") +
		lipgloss.NewStyle().
			Faint(true).
			Render(" to dismiss")

	content := errorTitleStyle.Render("Error") +
		"\n\n" +
		m.err.Error() +
		"\n\n" +
		hint

	return errorModalStyle.Render(content)
}
