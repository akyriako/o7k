package ui

import "strings"

func (m Model) renderHeader() string {
	lines := []string{
		headerLabelStyle.Render("Profile: ") +
			headerValueStyle.Render(m.profile),

		headerLabelStyle.Render("Region:  ") +
			headerValueStyle.Render(m.region),

		headerLabelStyle.Render("Project: ") +
			headerValueStyle.Render(m.project),

		headerLabelStyle.Render("Domain:  ") +
			headerValueStyle.Render(m.domain),
	}

	return strings.Join(lines, "\n")
}
