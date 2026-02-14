package tui

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
)

var (
	accentColor = compat.AdaptiveColor{Light: lipgloss.Color("#874BFD"), Dark: lipgloss.Color("#7D56F4")}
	subtleColor = compat.AdaptiveColor{Light: lipgloss.Color("#666666"), Dark: lipgloss.Color("#999999")}
	warnColor   = compat.AdaptiveColor{Light: lipgloss.Color("#FF6600"), Dark: lipgloss.Color("#FF9933")}

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
