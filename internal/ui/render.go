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

	// Compact mode: commands only, without context or logo.
	if m.width < compactWidth {
		return m.renderHeaderCommands(m.width)
	}

	context, contextWidth := m.renderHeaderContext()
	renderedLogo := logoStyle.Render(logo)

	leftWidth := contextWidth
	rightWidth := logoWidth

	// The logo is independently anchored to the right edge.
	// Hide the context if it would reduce the commands below their
	// preferred width.
	showContext := m.showHeaderContext()

	if !showContext {
		leftWidth = 0
	}

	centerWidth := max(m.width-leftWidth-rightWidth, 0)

	left := lipgloss.NewStyle().
		Width(leftWidth).
		Render(context)

	center := m.renderHeaderCommands(centerWidth)
	// Fill all available space before the logo, so the logo
	// remains at the terminal's right edge.
	center = lipgloss.NewStyle().
		Width(centerWidth).
		Render(center)

	right := lipgloss.NewStyle().
		Width(rightWidth).
		Align(lipgloss.Right).
		Render(renderedLogo)

	if !showContext {
		return lipgloss.JoinHorizontal(lipgloss.Top, center, right)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, left, center, right)
}

func (m Model) showHeaderContext() bool {
	if m.width < compactWidth {
		return false
	}

	const (
		commandsPerColumn  = 5
		commandColumnWidth = 28
	)

	commandCount := 5
	if m.resource != nil {
		commandCount += len(m.resource.Commands())
	}

	columnCount := (commandCount + commandsPerColumn - 1) / commandsPerColumn
	preferredCommandsWidth := columnCount * commandColumnWidth

	_, contextWidth := m.renderHeaderContext()

	return m.width-logoWidth-preferredCommandsWidth >= contextWidth
}

func (m Model) renderHeaderContext() (string, int) {
	const labelWidth = 10
	const contextRightMargin = 4
	var contextWidth = 0

	ctxFields := []struct {
		label string
		value string
	}{
		{"Cloud", m.context.Cloud},
		{"Region", m.context.Region},
		{"Domain", m.context.Domain},
		{"Project", m.context.Project},
		{"Version", m.version},
	}

	for _, f := range ctxFields {
		width := labelWidth + len(f.value)
		contextWidth = max(contextWidth, width)
	}

	valueWidth := contextWidth - labelWidth

	field := func(label, value string) string {
		value = truncate(value, valueWidth)

		return headerLabelStyle.Width(labelWidth).Render(label) +
			headerValueStyle.Width(valueWidth).Render(value)
	}

	fields := make([]string, 0, len(ctxFields))
	for _, f := range ctxFields {
		fields = append(fields, field(fmt.Sprintf("%s:", f.label), f.value))
	}

	return strings.Join(fields, "\n"), contextWidth + contextRightMargin
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

		if !m.showHeaderContext() && m.context.Cloud != "" {
			title = fmt.Sprintf(" (%s)%s", m.context.Cloud, title)
		}
	}

	innerWidth := max(m.width-2, 1)
	title = ansi.Truncate(title, innerWidth, "…")
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

	body := m.detailViewport.View()
	bodyLines := strings.Split(body, "\n")

	scrollbar := m.renderDetailScrollbar()
	var scrollbarLines []string

	if scrollbar != "" {
		scrollbarLines = strings.Split(scrollbar, "\n")
	}

	for i, line := range bodyLines {
		lineWidth := lipgloss.Width(line)

		if lineWidth < innerWidth {
			line += strings.Repeat(" ", innerWidth-lineWidth)
		}

		rightBorder := borderStyle.Render("│")

		if i < len(scrollbarLines) {
			rightBorder = scrollbarLines[i]
		}

		bodyLines[i] =
			borderStyle.Render("│") +
				line +
				rightBorder
	}

	return top +
		"\n" +
		strings.Join(bodyLines, "\n") +
		"\n" +
		bottom
}

func (m Model) renderDetailScrollbar() string {
	height := m.detailViewport.Height

	if height <= 0 || m.detailViewport.TotalLineCount() <= height {
		return ""
	}

	totalLines := m.detailViewport.TotalLineCount()
	maxOffset := max(totalLines-height, 1)

	thumbHeight := max(height*height/totalLines, 1)
	thumbHeight = min(thumbHeight, height)

	thumbTravel := height - thumbHeight
	thumbStart := m.detailViewport.YOffset * thumbTravel / maxOffset

	lines := make([]string, height)

	for i := range lines {
		if i >= thumbStart && i < thumbStart+thumbHeight {
			lines[i] = detailScrollbarThumbStyle.Render("█")
		} else {
			lines[i] = detailScrollbarTrackStyle.Render("│")
		}
	}

	return strings.Join(lines, "\n")
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
		m.errorViewport.View() +
		"\n\n" +
		hint

	modal := errorModalStyle.Render(content)

	scrollbar := m.renderErrorScrollbar()
	if scrollbar == "" {
		return modal
	}

	lines := strings.Split(modal, "\n")
	scrollbarLines := strings.Split(scrollbar, "\n")

	// One border row, one padding row, title row,
	// one blank row, then the viewport.
	const viewportStartRow = 4

	for i, scrollbarLine := range scrollbarLines {
		row := viewportStartRow + i

		if row >= len(lines)-1 {
			break
		}

		line := lines[row]
		width := ansi.StringWidth(line)

		// Preserve everything except the final border cell.
		lines[row] = ansi.Cut(line, 0, width-1) + scrollbarLine
	}

	return strings.Join(lines, "\n")
}

func (m Model) renderErrorScrollbar() string {
	height := m.errorViewport.Height

	if height <= 0 || m.errorViewport.TotalLineCount() <= height {
		return ""
	}

	totalLines := m.errorViewport.TotalLineCount()
	maxOffset := max(totalLines-height, 1)

	thumbHeight := max(height*height/totalLines, 1)
	thumbHeight = min(thumbHeight, height)

	thumbTravel := height - thumbHeight
	thumbStart := m.errorViewport.YOffset * thumbTravel / maxOffset

	lines := make([]string, height)

	for i := range lines {
		if i >= thumbStart && i < thumbStart+thumbHeight {
			lines[i] = errorScrollbarThumbStyle.Render("█")
		} else {
			lines[i] = errorScrollbarTrackStyle.Render("│")
		}
	}

	return strings.Join(lines, "\n")
}

func (m *Model) resizeErrorViewport() {
	const (
		modalWidth      = 60
		horizontalTrim  = 4
		scrollbarWidth  = 2
		verticalTrim    = 8
		maxVisibleLines = 12
	)

	contentWidth := modalWidth - horizontalTrim - scrollbarWidth
	m.errorViewport.Width = contentWidth

	if m.err == nil {
		m.errorViewport.Height = 1
		return
	}

	wrapped := lipgloss.NewStyle().
		Width(contentWidth).
		Render(m.err.Error())

	m.errorViewport.SetContent(wrapped)

	contentHeight := lipgloss.Height(wrapped)
	availableHeight := max(m.height-verticalTrim, 1)

	m.errorViewport.Height = max(
		min(contentHeight, availableHeight, maxVisibleLines),
		1,
	)

	m.errorViewport.SetYOffset(
		min(
			m.errorViewport.YOffset,
			max(contentHeight-m.errorViewport.Height, 0),
		),
	)
}
