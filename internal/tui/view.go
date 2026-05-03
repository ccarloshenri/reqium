package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	if m.mode == modeRequestForm {
		return appStyle.Render(m.renderRequestForm())
	}
	if m.mode == modeEnvForm {
		return appStyle.Render(m.renderEnvForm())
	}

	var builder strings.Builder
	builder.WriteString(m.renderHero() + "\n\n")

	if m.status != "" {
		builder.WriteString(okStyle.Render(m.status) + "\n\n")
	}
	if m.err != nil {
		builder.WriteString(badStyle.Render(m.err.Error()) + "\n\n")
	}

	builder.WriteString(m.renderQuickActions() + "\n\n")
	builder.WriteString(m.renderTabs() + "\n\n")

	var main string
	switch m.activeTab {
	case tabHistory:
		main = m.renderHistory()
	case tabCollections:
		main = m.renderCollections()
	case tabEnvironments:
		main = m.renderEnvironments()
	}

	sidebar := m.renderSidebar()
	builder.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, sidebar, "  ", panelStyle.Width(86).Render(main)))
	builder.WriteString("\n")
	builder.WriteString(helpStyle.Render("n new request  v variable  1/2/3 switch panels  r refresh  q quit"))
	builder.WriteString("\n")
	return appStyle.Render(builder.String())
}

func (m model) renderHero() string {
	logo := logoStyle.Render(strings.Join([]string{
		"╭────────────╮     ─────○",
		"│   >_       ├──   ─────○",
		"╰────────────╯     ─────○",
	}, "\n"))
	name := wordmark.Render("reqium")
	copy := mutedStyle.Render("Terminal API workspace")
	return lipgloss.JoinHorizontal(lipgloss.Center, logo, "   ", lipgloss.JoinVertical(lipgloss.Left, name, copy))
}

func (m model) renderQuickActions() string {
	request := panelStyle.Width(38).Render(labelStyle.Render("New Request") + "\n" + mutedStyle.Render("Press n to compose and send from here."))
	variable := panelStyle.Width(38).Render(labelStyle.Render("Environment Variable") + "\n" + mutedStyle.Render("Press v to add {{variables}} quickly."))
	return lipgloss.JoinHorizontal(lipgloss.Top, request, "  ", variable)
}

func (m model) renderSidebar() string {
	lines := []string{
		titleStyle.Render("Workspace"),
		fmt.Sprintf("Active env: %s", warnStyle.Render(m.activeEnvironmentName())),
		fmt.Sprintf("History: %d", len(m.history)),
		fmt.Sprintf("Collections: %d", len(m.collections)),
		fmt.Sprintf("Environments: %d", len(m.environments)),
	}
	if m.response != "" {
		lines = append(lines, "", labelStyle.Render("Last response"), truncate(m.response, 220))
	}
	return sidebarStyle.Render(strings.Join(lines, "\n"))
}

func (m model) renderTabs() string {
	labels := []string{"History", "Collections", "Environments"}
	parts := make([]string, len(labels))
	for i, label := range labels {
		if tab(i) == m.activeTab {
			parts[i] = activeTab.Render(label)
		} else {
			parts[i] = tabStyle.Render(label)
		}
	}
	return strings.Join(parts, " ")
}

func (m model) renderHistory() string {
	if len(m.history) == 0 {
		return mutedStyle.Render("No requests yet. Press n to create one.")
	}

	var builder strings.Builder
	for _, entry := range m.history {
		status := okStyle.Render(fmt.Sprintf("%d", entry.StatusCode))
		if entry.Error != "" || entry.StatusCode >= 400 {
			status = badStyle.Render(fmt.Sprintf("%d", entry.StatusCode))
		}
		fmt.Fprintf(&builder, "%s  %s  %s\n%s\n", labelStyle.Render(entry.Method), status, entry.URL, mutedStyle.Render(entry.ID))
	}
	return builder.String()
}

func (m model) renderCollections() string {
	if len(m.collections) == 0 {
		return mutedStyle.Render("No collections yet. Create one with reqium collection create <name>.")
	}

	var builder strings.Builder
	for _, collection := range m.collections {
		fmt.Fprintf(&builder, "%s  %d requests\n", collection.Name, len(collection.Requests))
		for _, req := range collection.Requests {
			fmt.Fprintf(&builder, "  %s  %s  %s\n", req.Name, req.Method, req.URL)
		}
	}
	return builder.String()
}

func (m model) renderRequestForm() string {
	var builder strings.Builder
	builder.WriteString(m.renderHero() + "\n\n")
	builder.WriteString(titleStyle.Render("Create Request") + "\n")
	if m.err != nil {
		builder.WriteString(badStyle.Render(m.err.Error()) + "\n")
	}

	builder.WriteString(field("Method", m.renderMethodPicker()) + "\n")
	builder.WriteString(field("URL", m.requestForm.url.View()) + "\n")
	builder.WriteString(field("Environment", m.requestForm.env.View()) + "\n")
	builder.WriteString(field("Headers", m.requestForm.headers.View()) + "\n")
	builder.WriteString(field("JSON Body", m.requestForm.body.View()) + "\n")
	builder.WriteString(helpStyle.Render("tab next field  m/left/right method  ctrl+s send  esc dashboard"))
	return panelStyle.Width(92).Render(builder.String())
}

func (m model) renderMethodPicker() string {
	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	parts := make([]string, len(methods))
	for i, method := range methods {
		style := tabStyle
		if i == m.requestForm.methodIndex {
			style = activeTab
		}
		if m.requestForm.focus == 0 && i == m.requestForm.methodIndex {
			style = activeTab.Border(lipgloss.NormalBorder(), false, false, true, false)
		}
		parts[i] = style.Render(method)
	}
	return strings.Join(parts, " ")
}

func (m model) renderEnvForm() string {
	var builder strings.Builder
	builder.WriteString(m.renderHero() + "\n\n")
	builder.WriteString(titleStyle.Render("Add Environment Variable") + "\n")
	if m.err != nil {
		builder.WriteString(badStyle.Render(m.err.Error()) + "\n")
	}
	builder.WriteString(field("Environment", m.envForm.env.View()) + "\n")
	builder.WriteString(field("Key", m.envForm.key.View()) + "\n")
	builder.WriteString(field("Value", m.envForm.value.View()) + "\n")
	builder.WriteString(helpStyle.Render("tab next field  enter/ctrl+s save  esc dashboard"))
	return panelStyle.Width(92).Render(builder.String())
}

func field(label string, value string) string {
	return labelStyle.Render(label) + "\n" + value + "\n"
}

func truncate(value string, max int) string {
	value = strings.TrimSpace(value)
	if len(value) <= max {
		return value
	}
	return value[:max] + "..."
}

func (m model) renderEnvironments() string {
	if len(m.environments) == 0 {
		return mutedStyle.Render("No environments yet. Create one with reqium env create <name>.")
	}

	var builder strings.Builder
	for _, env := range m.environments {
		prefix := " "
		if env.Active {
			prefix = "*"
		}
		fmt.Fprintf(&builder, "%s %s\n", prefix, env.Name)
		for key, value := range env.Variables {
			fmt.Fprintf(&builder, "  %s=%s\n", key, value)
		}
	}
	return builder.String()
}
