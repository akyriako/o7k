package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

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

func repeat(s string, count int) string {
	return strings.Repeat(s, max(count, 0))
}

func colorizeJSON(data []byte) string {
	lines := strings.Split(string(data), "\n")

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		indent := line[:len(line)-len(strings.TrimLeft(line, " "))]

		if strings.HasPrefix(trimmed, "\"") {
			if colon := strings.Index(trimmed, "\":"); colon >= 0 {
				key := trimmed[:colon+1]
				value := trimmed[colon+1:]

				lines[i] = indent +
					jsonKeyStyle.Render(key) +
					colorizeJSONValue(value)

				continue
			}
		}

		lines[i] = indent + colorizeJSONValue(trimmed)
	}

	return strings.Join(lines, "\n")
}

func colorizeJSONValue(value string) string {
	leading := value[:len(value)-len(strings.TrimLeft(value, " "))]
	trimmed := strings.TrimSpace(value)

	if trimmed == "" {
		return value
	}

	suffix := ""
	raw := trimmed

	if before, ok := strings.CutSuffix(raw, ","); ok {
		raw = before
		suffix = ","
	}

	switch {
	case strings.HasPrefix(raw, "\""):
		raw = jsonStringStyle.Render(raw)

	case raw == "true" || raw == "false":
		raw = jsonBoolStyle.Render(raw)

	case raw == "null":
		raw = jsonNullStyle.Render(raw)

	case raw != "{" && raw != "}" && raw != "[" && raw != "]":
		raw = jsonNumberStyle.Render(raw)
	}

	return leading + raw + suffix
}

func horizontalSlice(content string, offset, width int) string {
	if width <= 0 {
		return ""
	}

	lines := strings.Split(content, "\n")

	for i, line := range lines {
		lines[i] = ansi.Cut(line, offset, offset+width)
	}

	return strings.Join(lines, "\n")
}

func (m Model) maxTableXOffset() int {
	visibleWidth := max(m.width-2, 1)
	tableWidth := lipgloss.Width(m.table.View())

	return max(tableWidth-visibleWidth, 0)
}
