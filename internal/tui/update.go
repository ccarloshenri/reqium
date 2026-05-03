package tui

import tea "github.com/charmbracelet/bubbletea"

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "h", "left":
			if m.activeTab == tabHistory {
				m.activeTab = tabEnvironments
			} else {
				m.activeTab--
			}
		case "l", "right", "tab":
			if m.activeTab == tabEnvironments {
				m.activeTab = tabHistory
			} else {
				m.activeTab++
			}
		case "1":
			m.activeTab = tabHistory
		case "2":
			m.activeTab = tabCollections
		case "3":
			m.activeTab = tabEnvironments
		}
	}
	return m, nil
}
