package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/cboone/right-round/internal/data"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

const wideThreshold = 100
const tabBarHeight = 2

// focus tracks which panel has focus.
type focus int

const (
	focusGroups focus = iota
	focusEntries
	focusDetail
)

// activeTab tracks the current tab.
type activeTab int

const (
	tabSpinners activeTab = iota
	tabProgressBars
)

// Model is the top-level Bubble Tea model.
type Model struct {
	grouped    *data.GroupedEntries
	entryIndex map[string]*data.Entry

	list   listModel
	detail detailModel
	anim   *animEngine

	focus    focus
	tab      activeTab
	typeLock string // non-empty locks to one type

	width  int
	height int

	filtering   bool
	filterInput string
	filterBox   textinput.Model

	help help.Model

	optionsOpen   bool
	optionsForm   *huh.Form
	optionsFilter string
	optionsSort   string
	optionsDetail string
	optionsHelp   bool
	optionsType   string

	showFullHelp bool
	statusMsg    string
	statusExpiry time.Time

	lastTick time.Time
}

// New creates a new Model from grouped entries with optional type lock and initial group.
func New(grouped *data.GroupedEntries, typeLock string, initialGroup string) Model {
	anim := newAnimEngine()

	// Build entry index
	idx := make(map[string]*data.Entry)
	for i := range grouped.AllEntries {
		e := &grouped.AllEntries[i].Entry
		idx[e.ID] = e
	}

	tab := tabSpinners
	groups := grouped.SpinnerGroups
	if typeLock == "progress_bar" {
		tab = tabProgressBars
		groups = grouped.ProgressBarGroups
	}

	list := newListModel(groups, anim)
	filterBox := textinput.New()
	filterBox.Prompt = ""
	filterBox.Placeholder = "name or id"
	filterBox.CharLimit = 80
	filterBox.Width = 32

	helpModel := help.New()
	helpModel.Styles.ShortKey = helpKeyStyle
	helpModel.Styles.FullKey = helpKeyStyle
	helpModel.Styles.ShortDesc = helpDescStyle
	helpModel.Styles.FullDesc = helpDescStyle
	helpModel.Styles.ShortSeparator = helpDescStyle
	helpModel.Styles.FullSeparator = helpDescStyle
	helpModel.Ellipsis = " ..."

	// If initial group specified, move cursor to it
	if initialGroup != "" {
		list.selectGroupByName(initialGroup)
	}

	detail := newDetailModel(anim)

	return Model{
		grouped:    grouped,
		entryIndex: idx,
		list:       list,
		detail:     detail,
		anim:       anim,
		focus:      focusEntries,
		tab:        tab,
		typeLock:   typeLock,
		filterBox:  filterBox,
		help:       helpModel,
		lastTick:   time.Now(),
	}
}

// Init starts the animation ticker.
func (m Model) Init() tea.Cmd {
	return tick()
}

// Update handles messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateLayout()
		if m.optionsOpen && m.optionsForm != nil {
			m.optionsForm.WithWidth(m.optionsFormWidth()).WithHeight(m.optionsFormHeight())
		}
		if m.width >= wideThreshold {
			m.detail.setEntry(m.list.selectedEntry())
		}
		return m, nil

	case animTickMsg:
		now := time.Time(msg)
		elapsed := now.Sub(m.lastTick)
		m.lastTick = now

		// Collect visible entry IDs
		visibleIDs := m.list.visibleEntryIDs()
		if selected := m.list.selectedEntry(); selected != nil {
			alreadyVisible := false
			for _, id := range visibleIDs {
				if id == selected.Entry.ID {
					alreadyVisible = true
					break
				}
			}
			if !alreadyVisible {
				visibleIDs = append(visibleIDs, selected.Entry.ID)
			}
		}
		m.anim.advance(elapsed, visibleIDs, m.entryIndex)

		// Clear expired status message
		if m.statusMsg != "" && now.After(m.statusExpiry) {
			m.statusMsg = ""
		}

		return m, tick()

	case clipResultMsg:
		if msg.err != nil {
			m.statusMsg = fmt.Sprintf("Copy failed: %v", msg.err)
		} else {
			m.statusMsg = "Copied to clipboard!"
		}
		m.statusExpiry = time.Now().Add(2 * time.Second)
		return m, nil

	case tea.KeyMsg:
		if m.optionsOpen {
			return m.updateOptions(msg)
		}
		return m.handleKey(msg)

	case tea.MouseMsg:
		if m.optionsOpen {
			return m.updateOptions(msg)
		}
		return m.handleMouse(msg)
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Always allow ctrl+c to quit, even during filtering
	if msg.String() == "ctrl+c" {
		return m, tea.Quit
	}

	// Handle filtering mode
	if m.filtering {
		switch msg.String() {
		case "esc":
			m.filtering = false
			m.filterBox.Blur()
			m.filterInput = ""
			m.filterBox.SetValue("")
			m.list.setFilter("")
			return m, nil
		case "enter":
			m.filtering = false
			m.filterBox.Blur()
			m.filterInput = m.filterBox.Value()
			return m, nil
		default:
			var cmd tea.Cmd
			m.filterBox, cmd = m.filterBox.Update(msg)
			m.filterInput = m.filterBox.Value()
			m.list.setFilter(m.filterInput)
			return m, cmd
		}
	}

	switch {
	case matchKey(msg, keys.Quit):
		return m, tea.Quit

	case matchKey(msg, keys.Up):
		switch m.focus {
		case focusGroups:
			m.list.moveGroupUp()
		case focusEntries:
			m.list.moveEntryUp()
		default:
			m.detail.viewport.ScrollUp(1)
		}

	case matchKey(msg, keys.Down):
		switch m.focus {
		case focusGroups:
			m.list.moveGroupDown()
		case focusEntries:
			m.list.moveEntryDown()
		default:
			m.detail.viewport.ScrollDown(1)
		}

	case matchKey(msg, keys.PageUp):
		switch m.focus {
		case focusGroups:
			m.list.pageGroupUp()
		case focusEntries:
			m.list.pageEntryUp()
		default:
			m.detail.viewport.HalfPageUp()
		}

	case matchKey(msg, keys.PageDown):
		switch m.focus {
		case focusGroups:
			m.list.pageGroupDown()
		case focusEntries:
			m.list.pageEntryDown()
		default:
			m.detail.viewport.HalfPageDown()
		}

	case matchKey(msg, keys.Home):
		switch m.focus {
		case focusGroups:
			m.list.goGroupTop()
		case focusEntries:
			m.list.goEntryTop()
		default:
			m.detail.viewport.GotoTop()
		}

	case matchKey(msg, keys.End):
		switch m.focus {
		case focusGroups:
			m.list.goGroupBottom()
		case focusEntries:
			m.list.goEntryBottom()
		default:
			m.detail.viewport.GotoBottom()
		}

	case matchKey(msg, keys.Enter):
		if m.focus == focusGroups {
			m.focus = focusEntries
		} else if m.focus == focusEntries {
			if m.width < wideThreshold {
				m.focus = focusDetail
				m.detail.setEntry(m.list.selectedEntry())
			}
		}

	case matchKey(msg, keys.Back):
		if m.focus == focusDetail {
			m.focus = focusEntries
		} else if m.focus == focusEntries {
			m.focus = focusGroups
		}

	case matchKey(msg, keys.Left):
		if m.focus == focusEntries {
			m.focus = focusGroups
		} else if m.focus == focusDetail {
			m.focus = focusEntries
		}

	case matchKey(msg, keys.Right):
		if m.focus == focusGroups {
			m.focus = focusEntries
		} else if m.focus == focusEntries && m.width >= wideThreshold {
			m.focus = focusDetail
		}

	case matchKey(msg, keys.PrevGroup):
		m.list.moveGroupUp()
		m.focus = focusGroups

	case matchKey(msg, keys.NextGroup):
		m.list.moveGroupDown()
		m.focus = focusGroups

	case matchKey(msg, keys.Sort):
		m.list.toggleGroupSort()

	case matchKey(msg, keys.Verbose):
		m.detail.toggleVerbose()

	case matchKey(msg, keys.Options):
		cmd := m.openOptions()
		return m, cmd

	case matchKey(msg, keys.Tab):
		if m.typeLock == "" {
			if m.tab == tabSpinners {
				m.tab = tabProgressBars
				m.list.setGroups(m.grouped.ProgressBarGroups)
			} else {
				m.tab = tabSpinners
				m.list.setGroups(m.grouped.SpinnerGroups)
			}
			m.focus = focusEntries
		}

	case matchKey(msg, keys.Search):
		m.filtering = true
		m.filterBox.SetValue(m.filterInput)
		m.filterBox.CursorEnd()
		return m, m.filterBox.Focus()

	case matchKey(msg, keys.Copy):
		if entry := m.list.selectedEntry(); entry != nil {
			return m, copyToClipboard(entry)
		}

	case matchKey(msg, keys.Help):
		m.showFullHelp = !m.showFullHelp
	}

	// Update detail panel with current selection in wide mode
	if m.width >= wideThreshold {
		m.detail.setEntry(m.list.selectedEntry())
	}
	m.syncListFocus()

	return m, nil
}

func matchKey(msg tea.KeyMsg, binding key.Binding) bool {
	for _, k := range binding.Keys() {
		if msg.String() == k {
			return true
		}
	}
	return false
}

func (m *Model) updateLayout() {
	// Reserve space for tab bar and help bar.
	helpHeight := 1
	if m.showFullHelp {
		helpHeight = 2
	}
	contentHeight := m.height - tabBarHeight - helpHeight
	if contentHeight < 1 {
		contentHeight = 1
	}

	if m.width >= wideThreshold {
		listWidth := m.width * 55 / 100
		detailWidth := m.width - listWidth
		m.list.setSize(listWidth, contentHeight)
		m.detail.setSize(detailWidth, contentHeight)
	} else {
		m.list.setSize(m.width, contentHeight)
		m.detail.setSize(m.width, contentHeight)
	}
}

func (m *Model) syncListFocus() {
	if m.focus == focusGroups {
		m.list.setFocusPane(listPaneGroups)
	} else {
		m.list.setFocusPane(listPaneEntries)
	}
}

func (m *Model) openOptions() tea.Cmd {
	m.optionsFilter = m.filterInput
	m.optionsSort = m.list.groupSortLabel()
	if m.detail.verbose {
		m.optionsDetail = "verbose"
	} else {
		m.optionsDetail = "concise"
	}
	m.optionsHelp = m.showFullHelp
	if m.tab == tabProgressBars {
		m.optionsType = "progress bars"
	} else {
		m.optionsType = "spinners"
	}

	fields := []huh.Field{
		huh.NewInput().
			Title("Filter").
			Description("Search by name or id").
			Placeholder("braille, ascii, dots, loader").
			Value(&m.optionsFilter),
		huh.NewSelect[string]().
			Title("Group order").
			Options(
				huh.NewOption("Alphabetical", "alpha"),
				huh.NewOption("By group size", "size"),
			).
			Value(&m.optionsSort),
		huh.NewSelect[string]().
			Title("Detail panel").
			Options(
				huh.NewOption("Concise", "concise"),
				huh.NewOption("Verbose", "verbose"),
			).
			Value(&m.optionsDetail),
		huh.NewConfirm().
			Title("Show expanded key help").
			Affirmative("yes").
			Negative("no").
			Value(&m.optionsHelp),
	}

	if m.typeLock == "" {
		fields = append(fields,
			huh.NewSelect[string]().
				Title("Indicator type").
				Options(
					huh.NewOption("Spinners", "spinners"),
					huh.NewOption("Progress Bars", "progress bars"),
				).
				Value(&m.optionsType),
		)
	}

	m.optionsForm = huh.NewForm(
		huh.NewGroup(fields...).
			Title("Display and Navigation"),
	).WithTheme(huh.ThemeCharm()).
		WithWidth(m.optionsFormWidth()).
		WithHeight(m.optionsFormHeight())

	m.optionsOpen = true
	return m.optionsForm.Init()
}

func (m *Model) updateOptions(msg tea.Msg) (tea.Model, tea.Cmd) {
	if km, ok := msg.(tea.KeyMsg); ok && km.String() == "ctrl+c" {
		return m, tea.Quit
	}
	if m.optionsForm == nil {
		m.optionsOpen = false
		return m, nil
	}

	updated, cmd := m.optionsForm.Update(msg)
	if f, ok := updated.(*huh.Form); ok {
		m.optionsForm = f
	}

	switch m.optionsForm.State {
	case huh.StateCompleted:
		m.applyOptionsFromForm()
		m.optionsOpen = false
		m.optionsForm = nil
		m.statusMsg = "Options applied"
		m.statusExpiry = time.Now().Add(2 * time.Second)
		return m, nil
	case huh.StateAborted:
		m.optionsOpen = false
		m.optionsForm = nil
		return m, nil
	default:
		return m, cmd
	}
}

func (m *Model) applyOptionsFromForm() {
	m.filterInput = strings.TrimSpace(m.optionsFilter)
	m.filterBox.SetValue(m.filterInput)
	m.list.setFilter(m.filterInput)

	if m.optionsSort == "size" {
		m.list.setGroupSort(groupSortBySize)
	} else {
		m.list.setGroupSort(groupSortAlphabetical)
	}

	m.detail.setVerbose(m.optionsDetail == "verbose")
	m.showFullHelp = m.optionsHelp

	if m.typeLock == "" {
		if m.optionsType == "progress bars" {
			m.tab = tabProgressBars
			m.list.setGroups(m.grouped.ProgressBarGroups)
		} else {
			m.tab = tabSpinners
			m.list.setGroups(m.grouped.SpinnerGroups)
		}
	}

	if m.width >= wideThreshold {
		m.detail.setEntry(m.list.selectedEntry())
	}
	m.syncListFocus()
}

func (m Model) optionsFormWidth() int {
	if m.width <= 0 {
		return 64
	}
	w := m.width - 10
	if w > 76 {
		w = 76
	}
	if w < 36 {
		w = 36
	}
	return w
}

func (m Model) optionsFormHeight() int {
	if m.height <= 0 {
		return 16
	}
	h := m.height - 8
	if h > 20 {
		h = 20
	}
	if h < 10 {
		h = 10
	}
	return h
}

func (m Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if msg.Y < tabBarHeight && isMouseClick(msg) && m.typeLock == "" {
		if msg.X < m.width/2 {
			if m.tab != tabSpinners {
				m.tab = tabSpinners
				m.list.setGroups(m.grouped.SpinnerGroups)
			}
		} else {
			if m.tab != tabProgressBars {
				m.tab = tabProgressBars
				m.list.setGroups(m.grouped.ProgressBarGroups)
			}
		}
		m.focus = focusEntries
		m.syncListFocus()
		if m.width >= wideThreshold {
			m.detail.setEntry(m.list.selectedEntry())
		}
		return m, nil
	}

	contentY := msg.Y - tabBarHeight
	if contentY < 0 || contentY >= m.list.height {
		return m, nil
	}

	if m.width >= wideThreshold {
		if msg.X < m.list.width {
			localX := msg.X
			if isMouseWheel(msg) {
				delta := 0
				if msg.Button == tea.MouseButtonWheelDown {
					delta = 1
				} else if msg.Button == tea.MouseButtonWheelUp {
					delta = -1
				}
				if delta != 0 {
					if m.list.isGroupColumn(localX) {
						m.list.scrollGroup(delta)
						m.focus = focusGroups
					} else {
						m.list.scrollEntry(delta)
						m.focus = focusEntries
					}
				}
			} else if isMouseClick(msg) {
				if m.list.isGroupColumn(localX) {
					if m.list.clickGroupRow(contentY) {
						m.focus = focusGroups
					}
				} else {
					if m.list.clickEntryRow(contentY) {
						m.focus = focusEntries
					}
				}
			}
			m.syncListFocus()
			m.detail.setEntry(m.list.selectedEntry())
			return m, nil
		}

		if isMouseWheel(msg) {
			if msg.Button == tea.MouseButtonWheelDown {
				m.detail.viewport.ScrollDown(2)
			} else if msg.Button == tea.MouseButtonWheelUp {
				m.detail.viewport.ScrollUp(2)
			}
		}
		if isMouseClick(msg) {
			m.focus = focusDetail
		}
		m.syncListFocus()
		return m, nil
	}

	if m.focus == focusDetail {
		if isMouseWheel(msg) {
			if msg.Button == tea.MouseButtonWheelDown {
				m.detail.viewport.ScrollDown(2)
			} else if msg.Button == tea.MouseButtonWheelUp {
				m.detail.viewport.ScrollUp(2)
			}
		}
		return m, nil
	}

	if isMouseWheel(msg) {
		delta := 0
		if msg.Button == tea.MouseButtonWheelDown {
			delta = 1
		} else if msg.Button == tea.MouseButtonWheelUp {
			delta = -1
		}
		if m.list.isGroupColumn(msg.X) {
			m.list.scrollGroup(delta)
			m.focus = focusGroups
		} else {
			m.list.scrollEntry(delta)
			m.focus = focusEntries
		}
	}
	if isMouseClick(msg) {
		if m.list.isGroupColumn(msg.X) {
			if m.list.clickGroupRow(contentY) {
				m.focus = focusGroups
			}
		} else {
			if m.list.clickEntryRow(contentY) {
				m.focus = focusEntries
			}
		}
	}
	m.syncListFocus()
	return m, nil
}

// View renders the full TUI.
func (m Model) View() string {
	var b strings.Builder

	// Tab bar
	var tabLine string
	spinnerLabel := "Spinners"
	barLabel := "Progress Bars"
	if m.tab == tabSpinners {
		tabLine = activeTabStyle.Render(spinnerLabel) + inactiveTabStyle.Render(barLabel)
	} else {
		tabLine = inactiveTabStyle.Render(spinnerLabel) + activeTabStyle.Render(barLabel)
	}

	availableWidth := m.width - 2
	if availableWidth < 0 {
		availableWidth = 0
	}
	usedWidth := lipgloss.Width(tabLine)

	groupMeta := helpStyle.Render(" g:" + m.list.groupSortLabel())
	if usedWidth+lipgloss.Width(groupMeta) <= availableWidth {
		tabLine += groupMeta
		usedWidth += lipgloss.Width(groupMeta)
	}

	detailMode := "conc"
	if m.detail.verbose {
		detailMode = "verb"
	}
	detailMeta := helpStyle.Render(" d:" + detailMode)
	if usedWidth+lipgloss.Width(detailMeta) <= availableWidth {
		tabLine += detailMeta
		usedWidth += lipgloss.Width(detailMeta)
	}

	if m.typeLock != "" {
		lockMeta := helpStyle.Render(" lock")
		if usedWidth+lipgloss.Width(lockMeta) <= availableWidth {
			tabLine += lockMeta
		}
	}
	b.WriteString(tabBarStyle.Width(m.width).Render(tabLine))

	// Content area
	if m.optionsOpen && m.optionsForm != nil {
		modal := optionsModalStyle.
			Width(m.optionsFormWidth()).
			Render(m.optionsForm.View())
		b.WriteString("\n")
		b.WriteString(lipgloss.Place(m.width, m.list.height, lipgloss.Center, lipgloss.Center, modal))
	} else if m.width >= wideThreshold {
		listView := m.list.view()
		detailView := m.detail.view()
		content := lipgloss.JoinHorizontal(lipgloss.Top, listView, detailView)
		b.WriteString("\n")
		b.WriteString(content)
	} else if m.focus == focusDetail {
		b.WriteString("\n")
		b.WriteString(m.detail.view())
	} else {
		b.WriteString("\n")
		b.WriteString(m.list.view())
	}

	b.WriteString("\n")

	// Status / filter bar
	if m.filtering {
		b.WriteString(filterPromptStyle.Render("Filter ") + m.filterBox.View())
	} else if m.statusMsg != "" {
		b.WriteString(statusStyle.Render(m.statusMsg))
	} else if m.optionsOpen {
		b.WriteString(helpStyle.Render("Options: enter apply  esc cancel"))
	} else {
		helpModel := m.help
		helpModel.ShowAll = m.showFullHelp
		helpModel.Width = m.width
		b.WriteString(helpModel.View(keys))
	}

	return b.String()
}

func isMouseWheel(msg tea.MouseMsg) bool {
	return msg.Button == tea.MouseButtonWheelUp || msg.Button == tea.MouseButtonWheelDown
}

func isMouseClick(msg tea.MouseMsg) bool {
	return msg.Action == tea.MouseActionPress || msg.Action == tea.MouseActionRelease
}
