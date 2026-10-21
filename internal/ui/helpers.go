package ui

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
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

func (m *Model) suggestions() []string {
	var suggestions []string

	for _, r := range m.registry.Resources() {
		suggestions = append(suggestions, r.Kind())
		suggestions = append(suggestions, r.Aliases()...)
	}

	return suggestions
}

func (m *Model) setError(err error) {
	m.err = err

	if err == nil {
		m.errorViewport.SetContent("")
		m.errorViewport.GotoTop()
		return
	}

	m.autoRefreshPaused = true
	m.autoRefreshGeneration++
	m.autoRefreshScheduled = false

	m.resizeErrorViewport()
}

func (m Model) buildCacheKey(kind string) string {
	c := m.context
	identity := strings.Join([]string{
		c.Cloud,
		c.Region,
		c.Domain,
		c.Project,
	}, "\x00")

	hash := sha256.Sum256([]byte(identity))
	return fmt.Sprintf("%s:%s", hex.EncodeToString(hash[:8]), kind)
}

var ansiSGR = regexp.MustCompile("\x1b\\[[0-9;]*m")

func preserveBackground(content string, background lipgloss.Color) string {
	style := lipgloss.NewStyle().Background(background)
	prefix := ansiSGR.FindString(style.Render("x"))

	if prefix == "" {
		return content
	}

	return prefix + ansiSGR.ReplaceAllStringFunc(content, func(sequence string) string {
		return sequence + prefix
	})
}

//func (m *Model) wrapDetailContent() {
//	content := m.detailContent
//
//	if m.detailContentIsJson {
//		m.detailViewport.SetContent(colorizeJSON([]byte(content)))
//
//		if !m.detailContentWrap {
//			return
//		}
//	}
//
//	wrapped := lipgloss.NewStyle().
//		Width(max(m.detailViewport.Width, 1)).
//		Render(content)
//
//	if m.detailContentIsJson {
//		m.detailViewport.SetContent(colorizeJSON([]byte(wrapped)))
//	}
//
//	m.detailViewport.SetContent(wrapped)
//}

func (m *Model) wrapDetailContent() {
	content := m.detailContent

	if !m.detailContentWrap {
		if m.detailContentIsJson {
			content = colorizeJSON([]byte(content))
		}
		m.detailViewport.SetContent(content)
		return
	}

	wrapped := lipgloss.NewStyle().
		Width(max(m.detailViewport.Width, 1)).
		Render(content)

	if m.detailContentIsJson {
		wrapped = colorizeJSON([]byte(wrapped))
	}

	m.detailViewport.SetContent(wrapped)
}
