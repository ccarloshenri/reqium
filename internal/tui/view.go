package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	if m.mode == modeRequestForm {
		return appStyle.Render(m.center(m.renderRequestForm()))
	}
	if m.mode == modeEnvForm {
		return appStyle.Render(m.center(m.renderEnvForm()))
	}

	var builder strings.Builder
	builder.WriteString(m.renderHero() + "\n\n")

	if m.status != "" {
		builder.WriteString(m.center(okStyle.Render(m.status)) + "\n\n")
	}
	if m.err != nil {
		builder.WriteString(m.center(badStyle.Render(m.err.Error())) + "\n\n")
	}

	builder.WriteString(m.center(m.renderQuickActions()) + "\n\n")
	builder.WriteString(m.center(m.renderTabs()) + "\n\n")

	var main string
	switch m.activeTab {
	case tabHistory:
		main = m.renderHistory()
	case tabCollections:
		main = m.renderCollections()
	case tabEnvironments:
		main = m.renderEnvironments()
	}

	contentWidth := clamp(m.width-48, 62, 92)
	body := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.renderSidebar(),
		"  ",
		panelStyle.Width(contentWidth).Render(main),
	)
	builder.WriteString(m.center(body))
	builder.WriteString("\n")
	builder.WriteString(m.center(helpStyle.Render("n new request  v variable  1/2/3 panels  r refresh  q quit")))
	builder.WriteString("\n")

	return appStyle.Render(builder.String())
}

func (m model) renderHero() string {
	icon := strings.Join([]string{
		"          +----------------------+        o",
		"          |                      |------o",
		"          |        >_            |------o",
		"          |                      |------o",
		"          +----------------------+        o",
	}, "\n")

	word := strings.Join([]string{
		"  ____   _____   ___    ___   _   _   __  __",
		" |  _ \\ | ____| / _ \\  |_ _| | | | | |  \\/  |",
		" | |_) ||  _|  | | | |  | |  | | | | | |\\/| |",
		" |  _ < | |___ | |_| |  | |  | |_| | | |  | |",
		" |_| \\_\\|_____| \\__\\_\\ |___|  \\___/  |_|  |_|",
	}, "\n")

	hero := lipgloss.JoinVertical(
		lipgloss.Center,
		logoBlue.Render(icon),
		logoViolet.Render(word),
		"",
		welcomeStyle.Render("Welcome to Reqium!"),
		subtitleStyle.Render("Your terminal API workspace. Compose, send, inspect, repeat."),
	)
	return m.center(hero)
}

func (m model) renderQuickActions() string {
	request := actionCard("n", "New request", "Compose and send from inside Reqium.")
	variable := actionCard("v", "Add variable", "Create {{variables}} for environments.")
	history := actionCard("1", "History", "Review and replay recent calls.")
	return lipgloss.JoinHorizontal(lipgloss.Top, request, "  ", variable, "  ", history)
}

func (m model) renderSidebar() string {
	lines := []string{
		titleStyle.Render("Workspace snapshot"),
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
	labels := []string{"1 History", "2 Collections", "3 Environments"}
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
		return emptyState("No requests yet", "Press n to create your first request from the UI.")
	}

	var builder strings.Builder
	for _, entry := range m.history {
		status := okStyle.Render(fmt.Sprintf("%d", entry.StatusCode))
		if entry.Error != "" || entry.StatusCode >= 400 {
			status = badStyle.Render(fmt.Sprintf("%d", entry.StatusCode))
		}
		fmt.Fprintf(&builder, "%s  %s  %s\n%s\n\n", labelStyle.Render(entry.Method), status, entry.URL, mutedStyle.Render(entry.ID))
	}
	return builder.String()
}

func (m model) renderCollections() string {
	if len(m.collections) == 0 {
		return emptyState("No collections yet", "Create one with reqium collection create <name>.")
	}

	var builder strings.Builder
	for _, collection := range m.collections {
		fmt.Fprintf(&builder, "%s  %d requests\n", labelStyle.Render(collection.Name), len(collection.Requests))
		for _, req := range collection.Requests {
			fmt.Fprintf(&builder, "  %s  %s  %s\n", req.Name, req.Method, req.URL)
		}
		builder.WriteString("\n")
	}
	return builder.String()
}

func (m model) renderEnvironments() string {
	if len(m.environments) == 0 {
		return emptyState("No environments yet", "Press v to create a variable and activate an environment.")
	}

	var builder strings.Builder
	for _, env := range m.environments {
		prefix := " "
		if env.Active {
			prefix = "*"
		}
		fmt.Fprintf(&builder, "%s %s\n", prefix, labelStyle.Render(env.Name))
		for key, value := range env.Variables {
			fmt.Fprintf(&builder, "  %s=%s\n", key, value)
		}
		builder.WriteString("\n")
	}
	return builder.String()
}

func (m model) renderRequestForm() string {
	var builder strings.Builder
	builder.WriteString(m.renderHero() + "\n\n")
	builder.WriteString(welcomeStyle.Render("Create a request") + "\n")
	if m.err != nil {
		builder.WriteString(badStyle.Render(m.err.Error()) + "\n")
	}

	builder.WriteString(m.field("Method", m.renderMethodPicker(), m.requestForm.focus == 0) + "\n")
	builder.WriteString(m.field("URL", m.requestForm.url.View(), m.requestForm.focus == 1) + "\n")
	builder.WriteString(m.field("Environment", m.requestForm.env.View(), m.requestForm.focus == 2) + "\n")
	builder.WriteString(m.field("Headers", m.requestForm.headers.View(), m.requestForm.focus == 3) + "\n")
	builder.WriteString(m.field("JSON Body", m.requestForm.body.View(), m.requestForm.focus == 4) + "\n")
	builder.WriteString(helpStyle.Render("tab/down/j next  shift+tab/up/k previous  m method  ctrl+s send  esc dashboard"))
	return panelStyle.Width(clamp(m.width-12, 82, 108)).Render(builder.String())
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
	builder.WriteString(welcomeStyle.Render("Add an environment variable") + "\n")
	if m.err != nil {
		builder.WriteString(badStyle.Render(m.err.Error()) + "\n")
	}
	builder.WriteString(m.field("Environment", m.envForm.env.View(), m.envForm.focus == 0) + "\n")
	builder.WriteString(m.field("Key", m.envForm.key.View(), m.envForm.focus == 1) + "\n")
	builder.WriteString(m.field("Value", m.envForm.value.View(), m.envForm.focus == 2) + "\n")
	builder.WriteString(helpStyle.Render("tab/down/j next  shift+tab/up/k previous  enter/ctrl+s save  esc dashboard"))
	return panelStyle.Width(clamp(m.width-12, 82, 108)).Render(builder.String())
}

func (m model) field(label string, value string, focused bool) string {
	marker := "  "
	style := fieldStyle
	if focused {
		marker = "> "
		style = activeField
	}
	return marker + labelStyle.Render(label) + "\n" + style.Width(clamp(m.width-22, 68, 96)).Render(value)
}

func actionCard(key string, title string, description string) string {
	content := labelStyle.Render("["+key+"] "+title) + "\n" + mutedStyle.Render(description)
	return softPanelStyle.Width(29).Render(content)
}

func emptyState(title string, description string) string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		warnStyle.Render(title),
		mutedStyle.Render(description),
	)
}

func truncate(value string, max int) string {
	value = strings.TrimSpace(value)
	if len(value) <= max {
		return value
	}
	return value[:max] + "..."
}

func (m model) center(value string) string {
	width := m.width - 4
	if width < 80 {
		width = 80
	}
	return lipgloss.PlaceHorizontal(width, lipgloss.Center, value)
}

func clamp(value int, min int, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
