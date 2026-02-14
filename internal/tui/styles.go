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

	listDividerStyle = lipgloss.NewStyle().
				Foreground(subtleColor)

	listPaneTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(accentColor)

	listPaneMetaStyle = lipgloss.NewStyle().
				Foreground(subtleColor)

	// Detail panel styles
	detailBorderStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(accentColor).
				Padding(0, 1)

	detailLabelStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(accentColor)

	// Help bar
	helpStyle = lipgloss.NewStyle().
			Foreground(subtleColor)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(accentColor)

	helpDescStyle = lipgloss.NewStyle().
			Foreground(subtleColor)

	filterPromptStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(accentColor)

	// Status message
	statusStyle = lipgloss.NewStyle().
			Foreground(warnColor)
)
