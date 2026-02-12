package tui

import "github.com/charmbracelet/lipgloss"

var (
	accentColor = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	subtleColor = lipgloss.AdaptiveColor{Light: "#666666", Dark: "#999999"}
	warnColor   = lipgloss.AdaptiveColor{Light: "#FF6600", Dark: "#FF9933"}

	// Tab styles
	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor).
			Padding(0, 2)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(subtleColor).
				Padding(0, 2)

	tabBarStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(subtleColor)

	// List styles
	groupHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(accentColor).
				MarginTop(1)

	selectedItemStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(accentColor)

	normalItemStyle = lipgloss.NewStyle()

	// Detail panel styles
	detailBorderStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(accentColor).
				Padding(0, 1)

	detailLabelStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(accentColor)

	detailValueStyle = lipgloss.NewStyle()

	// Help bar
	helpStyle = lipgloss.NewStyle().
			Foreground(subtleColor)

	// Status message
	statusStyle = lipgloss.NewStyle().
			Foreground(warnColor)
)
