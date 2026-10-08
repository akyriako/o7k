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

const (
	minTerminalWidth  = 60
	minTerminalHeight = 15
	compactWidth      = 120
	logoWidth         = 23
)

func (m Model) renderHeader() string {
	if m.width < minTerminalWidth {
		return ""
	}

	const (
		commandsPerColumn  = 5
		commandColumnWidth = 28
		contextWidth       = 48
	)

	commandCount := 5
	if m.resource != nil {
		commandCount += len(m.resource.Commands())
	}

	columnCount := (commandCount + commandsPerColumn - 1) / commandsPerColumn
	preferredCommandsWidth := columnCount * commandColumnWidth

	// Compact mode: commands only, without context or logo.
	if m.width < compactWidth {
		return m.renderHeaderCommands(m.width)
	}

	profile := m.renderHeaderContext()
	renderedLogo := logoStyle.Render(logo)

	// Preserve the existing context breakpoint.
	leftWidth := contextWidth

	// The logo is independently anchored to the right edge.
	// Hide it if it would reduce the commands below their
	// preferred width.
	showLogo := m.width-leftWidth-preferredCommandsWidth >= logoWidth

	rightWidth := 0
	if showLogo {
		rightWidth = logoWidth
	}

	centerWidth := max(m.width-leftWidth-rightWidth, 0)

	left := lipgloss.NewStyle().
		Width(leftWidth).
		Render(profile)

	center := m.renderHeaderCommands(centerWidth)

	if !showLogo {
		return lipgloss.JoinHorizontal(lipgloss.Top, left, center)
	}

	// Fill all available space before the logo, so the logo
	// remains at the terminal's right edge.
	center = lipgloss.NewStyle().
		Width(centerWidth).
		Render(center)

	right := lipgloss.NewStyle().
		Width(rightWidth).
		Align(lipgloss.Right).
		Render(renderedLogo)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, center, right)
}

func (m Model) renderHeaderContext() string {
	const labelWidth = 10
	const contextWidth = 48

	valueWidth := contextWidth - labelWidth

	field := func(label, value string) string {
		value = truncate(value, valueWidth)

		return headerLabelStyle.Width(labelWidth).Render(label) +
			headerValueStyle.Width(valueWidth).Render(value)
	}

	return strings.Join([]string{
		field("Cloud:", m.context.Cloud),
		field("Region:", m.context.Region),
		field("Domain:", m.context.Domain),
		field("Project:", m.context.Project),
		field("Version:", m.version),
	}, "\n")
}

func (m Model) renderHeaderCommands(availableWidth int) string {
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

	const (
		commandsPerColumn  = 5
		commandColumnWidth = 30
		keyWidth           = 10
	)

	columnCount := (len(commands) + commandsPerColumn - 1) / commandsPerColumn

	if availableWidth <= 0 {
		return strings.Repeat("\n", commandsPerColumn-1)
	}

	widths := make([]int, columnCount)
	baseWidth := availableWidth / columnCount
	remainder := availableWidth % columnCount

	for i := range widths {
		widths[i] = baseWidth
		if i < remainder {
			widths[i]++
		}

		widths[i] = min(widths[i], commandColumnWidth)
	}

	columns := make([]string, 0, columnCount)

	for column := range columnCount {
		start := column * commandsPerColumn
		end := min(start+commandsPerColumn, len(commands))
		width := widths[column]

		lines := make([]string, 0, commandsPerColumn)

		for _, command := range commands[start:end] {
			description := command.Description
			if command.Default {
				description += " (Default)"
			}

			// Preserve the original keybinding spacing.
			key := headerCommandKeyStyle.
				Width(keyWidth).
				Render("<" + command.Key + ">")

			text := headerCommandTextStyle.Render(description)

			// Truncate the description, not the spacing between
			// the keybinding and the description.
			if width <= keyWidth {
				lines = append(lines, ansi.Cut(key, 0, width))
				continue
			}

			text = ansi.Truncate(text, max(width-keyWidth-1, 0), "")
			lines = append(lines, key+text+" ")
		}

		for len(lines) < commandsPerColumn {
			lines = append(lines, "")
		}

		columns = append(columns, strings.Join(lines, "\n"))
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

		if m.width < compactWidth && m.context.Cloud != "" {
			title = fmt.Sprintf(" (%s)%s", m.context.Cloud, title)
		}
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
