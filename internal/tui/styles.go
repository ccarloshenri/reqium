package tui

import "github.com/charmbracelet/lipgloss"

var (
	appStyle     = lipgloss.NewStyle().Padding(1, 2)
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	logoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("45")).Bold(true)
	wordmark     = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true)
	panelStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62")).Padding(1, 2)
	sidebarStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("238")).Padding(1, 2).Width(34)
	tabStyle     = lipgloss.NewStyle().Padding(0, 1)
	activeTab    = tabStyle.Bold(true).Foreground(lipgloss.Color("229")).Background(lipgloss.Color("62"))
	labelStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("117")).Bold(true)
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	mutedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	badStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
	okStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("84"))
	warnStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
)
