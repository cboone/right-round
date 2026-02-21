package tui

import (
	"charm.land/lipgloss/v2"
)

var (
	// Warmer palette inspired by Huh's Base16 theme.
	primaryColor     = lipgloss.Color("9")
	interactiveColor = lipgloss.Color("11")
	selectedColor    = lipgloss.Color("13")
	subtleColor      = lipgloss.Color("8")
	warnColor        = lipgloss.Color("1")

	// Tab styles
	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Padding(0, 2)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(subtleColor).
				Padding(0, 2)

	tabBarStyle = lipgloss.NewStyle().
			Foreground(subtleColor)

	// List styles
	groupHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(primaryColor).
				MarginTop(1)

	selectedItemStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(selectedColor)

	normalItemStyle = lipgloss.NewStyle()

	listPaneTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(primaryColor)

	listPaneMetaStyle = lipgloss.NewStyle().
				Foreground(subtleColor)

	// Detail panel styles
	detailBorderStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(subtleColor).
				Padding(0, 1)

	detailLabelStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(primaryColor)

	// Help bar
	helpStyle = lipgloss.NewStyle().
			Foreground(subtleColor)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(interactiveColor)

	helpDescStyle = lipgloss.NewStyle().
			Foreground(subtleColor)

	filterPromptStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(interactiveColor)

	// Status message
	statusStyle = lipgloss.NewStyle().
			Foreground(warnColor)
)
