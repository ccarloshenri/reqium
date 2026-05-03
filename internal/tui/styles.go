package tui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	tabStyle   = lipgloss.NewStyle().Padding(0, 1)
	activeTab  = tabStyle.Bold(true).Foreground(lipgloss.Color("229")).Background(lipgloss.Color("62"))
	mutedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	badStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
	okStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("84"))
)
