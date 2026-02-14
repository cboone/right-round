package tui

import (
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cboone/right-round/internal/data"
)

const wideThreshold = 100
const topGapHeight = 1
const tabBarHeight = 1
const contentGapHeight = 1
const helpGapHeight = 1

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
	filterBox.SetWidth(32)
	filterBox.ShowSuggestions = true

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

	m := Model{
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
	m.refreshFilterSuggestions()
	return m
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
			m.updateLayout()
		}

		return m, tick()

	case clipResultMsg:
		if msg.err != nil {
			m.statusMsg = fmt.Sprintf("Copy failed: %v", msg.err)
		} else {
			m.statusMsg = "Copied to clipboard!"
		}
		m.statusExpiry = time.Now().Add(2 * time.Second)
		m.updateLayout()
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case tea.MouseMsg:
		return m.handleMouse(msg)
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	defer m.updateLayout()

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

	case matchKey(msg, keys.Tab):
		if m.typeLock == "" {
			if m.tab == tabSpinners {
				m.tab = tabProgressBars
				m.list.setGroups(m.grouped.ProgressBarGroups)
			} else {
				m.tab = tabSpinners
				m.list.setGroups(m.grouped.SpinnerGroups)
			}
			m.refreshFilterSuggestions()
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
	contentHeight := m.height - topGapHeight - tabBarHeight - contentGapHeight - helpGapHeight - m.bottomBarHeight()
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

func (m Model) bottomBarHeight() int {
	if m.width <= 0 {
		return 1
	}
	if m.filtering || m.statusMsg != "" {
		return 1
	}
	helpModel := m.help
	helpModel.ShowAll = m.showFullHelp
	helpModel.SetWidth(m.width)
	h := lipgloss.Height(helpModel.View(m.currentHelpKeyMap()))
	if h < 1 {
		return 1
	}
	return h
}

func (m *Model) syncListFocus() {
	if m.focus == focusGroups {
		m.list.setFocusPane(listPaneGroups)
	} else {
		m.list.setFocusPane(listPaneEntries)
	}
}

func (m Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	defer m.updateLayout()
	mouse := msg.Mouse()
	x := mouse.X
	y := mouse.Y

	if y >= topGapHeight && y < topGapHeight+tabBarHeight && isMouseClick(msg) && m.typeLock == "" {
		if tab, ok := m.tabAtX(x); ok {
			if tab == tabSpinners && m.tab != tabSpinners {
				m.tab = tabSpinners
				m.list.setGroups(m.grouped.SpinnerGroups)
			} else if tab == tabProgressBars && m.tab != tabProgressBars {
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

	contentY := y - topGapHeight - tabBarHeight - contentGapHeight
	if contentY < 0 || contentY >= m.list.height {
		return m, nil
	}

	if m.width >= wideThreshold {
		if x < m.list.width {
			localX := x
			if delta := mouseWheelDelta(msg); delta != 0 {
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

		if delta := mouseWheelDelta(msg); delta != 0 {
			if delta > 0 {
				m.detail.viewport.ScrollDown(2)
			} else {
				m.detail.viewport.ScrollUp(2)
			}
		} else {
			m.detail.viewport, _ = m.detail.viewport.Update(msg)
		}
		if isMouseClick(msg) {
			m.focus = focusDetail
		}
		m.syncListFocus()
		return m, nil
	}

	if m.focus == focusDetail {
		if delta := mouseWheelDelta(msg); delta != 0 {
			if delta > 0 {
				m.detail.viewport.ScrollDown(2)
			} else {
				m.detail.viewport.ScrollUp(2)
			}
		} else {
			m.detail.viewport, _ = m.detail.viewport.Update(msg)
		}
		return m, nil
	}

	if delta := mouseWheelDelta(msg); delta != 0 {
		if m.list.isGroupColumn(x) {
			m.list.scrollGroup(delta)
			m.focus = focusGroups
		} else {
			m.list.scrollEntry(delta)
			m.focus = focusEntries
		}
	}
	if isMouseClick(msg) {
		if m.list.isGroupColumn(x) {
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
func (m Model) View() tea.View {
	v := tea.NewView(m.render())
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

func (m Model) render() string {
	// Recompute layout from the exact footer we're about to render.
	bottom := ""
	if m.filtering {
		bottom = filterPromptStyle.Render("Filter ") + m.filterBox.View()
	} else if m.statusMsg != "" {
		bottom = statusStyle.Render(m.statusMsg)
	} else {
		helpModel := m.help
		helpModel.ShowAll = m.showFullHelp
		helpModel.SetWidth(m.width)
		bottom = helpModel.View(m.currentHelpKeyMap())
	}

	bottomHeight := lipgloss.Height(bottom)
	if bottomHeight < 1 {
		bottomHeight = 1
	}
	contentHeight := m.height - topGapHeight - tabBarHeight - contentGapHeight - helpGapHeight - bottomHeight
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

	if m.typeLock != "" {
		lockMeta := helpStyle.Render(" lock")
		if usedWidth+lipgloss.Width(lockMeta) <= availableWidth {
			tabLine += lockMeta
		}
	}
	b.WriteString("\n")
	b.WriteString(tabBarStyle.Width(m.width).Render(" " + tabLine))
	b.WriteString("\n\n")

	// Content area
	if m.width >= wideThreshold {
		listView := m.list.view()
		detailView := m.detail.view()
		content := lipgloss.JoinHorizontal(lipgloss.Top, listView, detailView)
		b.WriteString(content)
	} else if m.focus == focusDetail {
		b.WriteString(m.detail.view())
	} else {
		b.WriteString(m.list.view())
	}

	b.WriteString("\n\n")

	// Status / filter bar
	b.WriteString(bottom)

	return b.String()
}

func mouseWheelDelta(msg tea.MouseMsg) int {
	m := msg.Mouse()
	if m.Button == tea.MouseWheelDown {
		return 1
	}
	if m.Button == tea.MouseWheelUp {
		return -1
	}
	return 0
}

func isMouseClick(msg tea.MouseMsg) bool {
	switch msg.(type) {
	case tea.MouseClickMsg, tea.MouseReleaseMsg:
		return true
	default:
		return false
	}
}

func (m *Model) refreshFilterSuggestions() {
	seen := make(map[string]struct{})
	suggestions := make([]string, 0, 96)

	appendSuggestion := func(v string) {
		v = strings.TrimSpace(v)
		if v == "" {
			return
		}
		k := strings.ToLower(v)
		if _, ok := seen[k]; ok {
			return
		}
		seen[k] = struct{}{}
		suggestions = append(suggestions, v)
	}

	for i := range m.list.groups {
		g := m.list.groups[i]
		appendSuggestion(g.Name)
		for j := range g.Entries {
			appendSuggestion(g.Entries[j].Entry.Name)
			appendSuggestion(g.Entries[j].Entry.ID)
			if len(suggestions) >= 96 {
				break
			}
		}
		if len(suggestions) >= 96 {
			break
		}
	}

	m.filterBox.SetSuggestions(suggestions)
}

func (m Model) tabAtX(x int) (activeTab, bool) {
	spinnerWidth := lipgloss.Width(inactiveTabStyle.Render("Spinners"))
	barWidth := lipgloss.Width(inactiveTabStyle.Render("Progress Bars"))
	if x >= 0 && x < spinnerWidth {
		return tabSpinners, true
	}
	if x >= spinnerWidth && x < spinnerWidth+barWidth {
		return tabProgressBars, true
	}
	return tabSpinners, false
}

type contextualHelpKeyMap struct {
	short []key.Binding
	full  [][]key.Binding
}

func (k contextualHelpKeyMap) ShortHelp() []key.Binding {
	return k.short
}

func (k contextualHelpKeyMap) FullHelp() [][]key.Binding {
	return k.full
}

func (m Model) currentHelpKeyMap() contextualHelpKeyMap {
	nav := []key.Binding{keys.Up, keys.Down, keys.PageUp, keys.PageDown, keys.Home, keys.End}
	manage := []key.Binding{keys.Search, keys.Help, keys.Quit}

	if m.typeLock == "" {
		manage = append(manage, keys.Tab)
	}

	switch m.focus {
	case focusGroups:
		short := []key.Binding{keys.Up, keys.Right, keys.PrevGroup, keys.Search, keys.Quit}
		full := [][]key.Binding{
			nav,
			{keys.Left, keys.Right, keys.PrevGroup, keys.NextGroup},
			manage,
		}
		return contextualHelpKeyMap{short: short, full: full}

	case focusDetail:
		short := []key.Binding{keys.Up, keys.Back, keys.Copy, keys.Search, keys.Quit}
		full := [][]key.Binding{
			nav,
			{keys.Back, keys.Left, keys.Copy},
			manage,
		}
		return contextualHelpKeyMap{short: short, full: full}

	default:
		entryFocus := []key.Binding{keys.Left, keys.Right, keys.Enter, keys.Copy}
		if m.width >= wideThreshold {
			entryFocus = []key.Binding{keys.Left, keys.Right, keys.Copy}
		}
		short := []key.Binding{keys.Up, keys.Left, keys.Right, keys.Copy, keys.Search, keys.Quit}
		if m.width < wideThreshold {
			short = []key.Binding{keys.Up, keys.Left, keys.Enter, keys.Copy, keys.Search, keys.Quit}
		}
		full := [][]key.Binding{
			nav,
			entryFocus,
			manage,
		}
		return contextualHelpKeyMap{short: short, full: full}
	}
}
