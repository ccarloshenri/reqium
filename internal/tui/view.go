package tui

import (
	"fmt"
	"strings"
)

func (m model) View() string {
	var builder strings.Builder
	builder.WriteString(titleStyle.Render("Reqium") + "\n")
	builder.WriteString(m.renderTabs() + "\n\n")

	switch m.activeTab {
	case tabHistory:
		builder.WriteString(m.renderHistory())
	case tabCollections:
		builder.WriteString(m.renderCollections())
	case tabEnvironments:
		builder.WriteString(m.renderEnvironments())
	}

	builder.WriteString("\n")
	builder.WriteString(mutedStyle.Render("1 history  2 collections  3 environments  left/right switch  q quit"))
	builder.WriteString("\n")
	return builder.String()
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
		return mutedStyle.Render("No requests yet. Run a request and it will appear here.")
	}

	var builder strings.Builder
	for _, entry := range m.history {
		status := okStyle.Render(fmt.Sprintf("%d", entry.StatusCode))
		if entry.Error != "" || entry.StatusCode >= 400 {
			status = badStyle.Render(fmt.Sprintf("%d", entry.StatusCode))
		}
		fmt.Fprintf(&builder, "%s  %s  %s  %s\n", entry.ID, entry.Method, status, entry.URL)
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
