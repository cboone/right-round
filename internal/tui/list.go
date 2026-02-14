package tui

import (
	"sort"
	"strings"

	"charm.land/bubbles/v2/paginator"
	"charm.land/lipgloss/v2"
	"github.com/cboone/right-round/internal/data"
)

type listPaneFocus int

const (
	listPaneGroups listPaneFocus = iota
	listPaneEntries
)

type visibleGroup struct {
	name    string
	entries []data.EntryEnvelope
}

// listModel manages the grouped list panel.
type listModel struct {
	groups        []data.Group
	visibleGroups []visibleGroup
	groupCursor   int
	groupOffset   int
	entryCursor   map[string]int
	entryOffset   map[string]int
	height        int
	width         int
	filter        string
	focusPane     listPaneFocus
	groupPager    paginator.Model
	entryPager    paginator.Model

	anim *animEngine
}

func newListModel(groups []data.Group, anim *animEngine) listModel {
	m := listModel{
		groups:      groups,
		entryCursor: make(map[string]int),
		entryOffset: make(map[string]int),
		focusPane:   listPaneEntries,
		anim:        anim,
	}
	m.groupPager = paginator.New()
	m.groupPager.Type = paginator.Arabic
	m.groupPager.ArabicFormat = "%d/%d"
	m.entryPager = paginator.New()
	m.entryPager.Type = paginator.Arabic
	m.entryPager.ArabicFormat = "%d/%d"
	m.rebuildVisibleGroups("")
	return m
}

func (m *listModel) buildVisibleGroups(groups []data.Group, filter string) []visibleGroup {
	filter = strings.ToLower(filter)
	var visible []visibleGroup
	for i := range groups {
		g := &groups[i]
		var matching []data.EntryEnvelope
		for j := range g.Entries {
			e := g.Entries[j]
			if filter == "" || strings.Contains(strings.ToLower(e.Entry.Name), filter) || strings.Contains(strings.ToLower(e.Entry.ID), filter) {
				matching = append(matching, e)
			}
		}
		if len(matching) == 0 {
			continue
		}
		visible = append(visible, visibleGroup{name: g.Name, entries: matching})
	}

	sort.SliceStable(visible, func(i, j int) bool {
		return strings.ToLower(visible[i].name) < strings.ToLower(visible[j].name)
	})

	return visible
}

func (m *listModel) rebuildVisibleGroups(filter string) {
	oldID := m.selectedID()
	oldGroup := m.selectedGroupName()

	m.filter = filter
	m.visibleGroups = m.buildVisibleGroups(m.groups, filter)

	if len(m.visibleGroups) == 0 {
		m.groupCursor = 0
		m.groupOffset = 0
		return
	}

	m.groupCursor = 0
	if oldGroup != "" {
		for i := range m.visibleGroups {
			if strings.EqualFold(m.visibleGroups[i].name, oldGroup) {
				m.groupCursor = i
				break
			}
		}
	}

	if oldID != "" {
		for gi := range m.visibleGroups {
			for ei := range m.visibleGroups[gi].entries {
				if m.visibleGroups[gi].entries[ei].Entry.ID == oldID {
					m.groupCursor = gi
					m.entryCursor[m.visibleGroups[gi].name] = ei
					break
				}
			}
		}
	}

	m.ensureCurrentGroupState()
	m.ensureGroupVisible()
	m.ensureEntryVisible()
}

func (m *listModel) setFilter(filter string) {
	m.rebuildVisibleGroups(filter)
}

func (m *listModel) setGroups(groups []data.Group) {
	m.groups = groups
	m.rebuildVisibleGroups(m.filter)
}

func (m *listModel) selectGroupByName(name string) bool {
	for i := range m.visibleGroups {
		if strings.EqualFold(m.visibleGroups[i].name, name) {
			m.groupCursor = i
			m.ensureCurrentGroupState()
			m.ensureGroupVisible()
			m.ensureEntryVisible()
			return true
		}
	}
	return false
}

func (m *listModel) selectedEntry() *data.EntryEnvelope {
	g := m.currentGroup()
	if g == nil || len(g.entries) == 0 {
		return nil
	}
	idx := m.entryCursor[g.name]
	if idx < 0 || idx >= len(g.entries) {
		return nil
	}
	return &g.entries[idx]
}

func (m *listModel) selectedID() string {
	if e := m.selectedEntry(); e != nil {
		return e.Entry.ID
	}
	return ""
}

func (m *listModel) selectedGroupName() string {
	g := m.currentGroup()
	if g == nil {
		return ""
	}
	return g.name
}

func (m *listModel) currentGroup() *visibleGroup {
	if m.groupCursor < 0 || m.groupCursor >= len(m.visibleGroups) {
		return nil
	}
	return &m.visibleGroups[m.groupCursor]
}

func (m *listModel) setFocusPane(p listPaneFocus) {
	m.focusPane = p
}

func (m *listModel) setSize(width, height int) {
	m.width = width
	m.height = height
	m.ensureGroupVisible()
	m.ensureEntryVisible()
	m.updatePagers()
}

func (m *listModel) moveUp() {
	m.moveEntryUp()
}

func (m *listModel) moveDown() {
	m.moveEntryDown()
}

func (m *listModel) pageUp() {
	m.pageEntryUp()
}

func (m *listModel) pageDown() {
	m.pageEntryDown()
}

func (m *listModel) goToTop() {
	m.goEntryTop()
}

func (m *listModel) goToBottom() {
	m.goEntryBottom()
}

func (m *listModel) moveGroupUp() {
	if m.groupCursor > 0 {
		m.groupCursor--
		m.ensureCurrentGroupState()
		m.ensureGroupVisible()
		m.ensureEntryVisible()
	}
}

func (m *listModel) moveGroupDown() {
	if m.groupCursor < len(m.visibleGroups)-1 {
		m.groupCursor++
		m.ensureCurrentGroupState()
		m.ensureGroupVisible()
		m.ensureEntryVisible()
	}
}

func (m *listModel) pageGroupUp() {
	m.groupCursor -= m.scrollRows()
	if m.groupCursor < 0 {
		m.groupCursor = 0
	}
	m.ensureCurrentGroupState()
	m.ensureGroupVisible()
	m.ensureEntryVisible()
}

func (m *listModel) pageGroupDown() {
	m.groupCursor += m.scrollRows()
	if m.groupCursor >= len(m.visibleGroups) {
		m.groupCursor = len(m.visibleGroups) - 1
	}
	m.ensureCurrentGroupState()
	m.ensureGroupVisible()
	m.ensureEntryVisible()
}

func (m *listModel) goGroupTop() {
	m.groupCursor = 0
	m.ensureCurrentGroupState()
	m.ensureGroupVisible()
	m.ensureEntryVisible()
}

func (m *listModel) goGroupBottom() {
	m.groupCursor = len(m.visibleGroups) - 1
	m.ensureCurrentGroupState()
	m.ensureGroupVisible()
	m.ensureEntryVisible()
}

func (m *listModel) moveEntryUp() {
	g := m.currentGroup()
	if g == nil {
		return
	}
	name := g.name
	if m.entryCursor[name] > 0 {
		m.entryCursor[name]--
		m.ensureEntryVisible()
	}
}

func (m *listModel) moveEntryDown() {
	g := m.currentGroup()
	if g == nil {
		return
	}
	name := g.name
	if m.entryCursor[name] < len(g.entries)-1 {
		m.entryCursor[name]++
		m.ensureEntryVisible()
	}
}

func (m *listModel) pageEntryUp() {
	g := m.currentGroup()
	if g == nil {
		return
	}
	name := g.name
	m.entryCursor[name] -= m.scrollRows()
	if m.entryCursor[name] < 0 {
		m.entryCursor[name] = 0
	}
	m.ensureEntryVisible()
}

func (m *listModel) pageEntryDown() {
	g := m.currentGroup()
	if g == nil {
		return
	}
	name := g.name
	m.entryCursor[name] += m.scrollRows()
	if m.entryCursor[name] >= len(g.entries) {
		m.entryCursor[name] = len(g.entries) - 1
	}
	m.ensureEntryVisible()
}

func (m *listModel) goEntryTop() {
	g := m.currentGroup()
	if g == nil {
		return
	}
	m.entryCursor[g.name] = 0
	m.ensureEntryVisible()
}

func (m *listModel) goEntryBottom() {
	g := m.currentGroup()
	if g == nil {
		return
	}
	m.entryCursor[g.name] = len(g.entries) - 1
	m.ensureEntryVisible()
}

func (m *listModel) ensureCurrentGroupState() {
	g := m.currentGroup()
	if g == nil {
		return
	}
	if len(g.entries) == 0 {
		m.entryCursor[g.name] = 0
		m.entryOffset[g.name] = 0
		return
	}
	c := m.entryCursor[g.name]
	if c < 0 {
		m.entryCursor[g.name] = 0
	}
	if c >= len(g.entries) {
		m.entryCursor[g.name] = len(g.entries) - 1
	}
	o := m.entryOffset[g.name]
	if o < 0 {
		m.entryOffset[g.name] = 0
	}
}

func (m *listModel) ensureGroupVisible() {
	rows := m.visibleRows()
	if rows < 1 {
		rows = 1
	}
	if m.groupCursor < m.groupOffset {
		m.groupOffset = m.groupCursor
	}
	if m.groupCursor >= m.groupOffset+rows {
		m.groupOffset = m.groupCursor - rows + 1
	}
	if m.groupOffset < 0 {
		m.groupOffset = 0
	}
	m.updatePagers()
}

func (m *listModel) ensureEntryVisible() {
	g := m.currentGroup()
	if g == nil {
		m.updatePagers()
		return
	}
	rows := m.visibleRows()
	if rows < 1 {
		rows = 1
	}
	name := g.name
	if m.entryCursor[name] < m.entryOffset[name] {
		m.entryOffset[name] = m.entryCursor[name]
	}
	if m.entryCursor[name] >= m.entryOffset[name]+rows {
		m.entryOffset[name] = m.entryCursor[name] - rows + 1
	}
	if m.entryOffset[name] < 0 {
		m.entryOffset[name] = 0
	}
	m.updatePagers()
}

// visibleEntryIDs returns the IDs of entries currently visible for animation.
func (m *listModel) visibleEntryIDs() []string {
	var ids []string
	g := m.currentGroup()
	if g == nil {
		return ids
	}
	start := m.entryOffset[g.name]
	rows := m.visibleRows()
	if rows < 1 {
		rows = 1
	}
	end := start + rows
	if end > len(g.entries) {
		end = len(g.entries)
	}
	for i := start; i < end; i++ {
		ids = append(ids, g.entries[i].Entry.ID)
	}
	return ids
}

func (m *listModel) view() string {
	if len(m.visibleGroups) == 0 {
		return helpStyle.Render("  No matches")
	}
	groupWidth, entryWidth := m.columnWidths()
	rows := m.visibleRows()

	groupsTitle := listPaneTitleStyle.Render("Categories")
	groupsHeader := lipgloss.NewStyle().Width(groupWidth).MaxWidth(groupWidth).Render(groupsTitle)

	entriesTitle := listPaneTitleStyle.Render("Entries")
	entriesHeader := lipgloss.NewStyle().Width(entryWidth).MaxWidth(entryWidth).Render(entriesTitle)
	if rows < 1 {
		return groupsHeader + listDividerStyle.Render("|") + entriesHeader
	}

	var groupLines []string
	for i := 0; i < rows; i++ {
		idx := m.groupOffset + i
		if idx >= len(m.visibleGroups) {
			groupLines = append(groupLines, strings.Repeat(" ", groupWidth))
			continue
		}
		g := m.visibleGroups[idx]
		label := truncateWithEllipsis(g.name, groupWidth-2)
		style := normalItemStyle
		prefix := "  "
		if idx == m.groupCursor {
			style = selectedItemStyle
			if m.focusPane == listPaneGroups {
				prefix = "> "
			} else {
				prefix = "* "
			}
		}
		line := prefix + style.Render(label)
		line = lipgloss.NewStyle().Width(groupWidth).MaxWidth(groupWidth).Render(line)
		groupLines = append(groupLines, line)
	}

	previewColWidth := 8
	nameWidth := entryWidth - previewColWidth - 3
	if nameWidth < 1 {
		nameWidth = 1
	}

	var entryLines []string
	g := m.currentGroup()
	if g == nil {
		for i := 0; i < rows; i++ {
			entryLines = append(entryLines, strings.Repeat(" ", entryWidth))
		}
	} else {
		start := m.entryOffset[g.name]
		end := start + rows
		if end > len(g.entries) {
			end = len(g.entries)
		}
		for i := 0; i < rows; i++ {
			idx := start + i
			if idx >= end {
				entryLines = append(entryLines, strings.Repeat(" ", entryWidth))
				continue
			}
			entry := g.entries[idx]
			name := truncateWithEllipsis(entry.Entry.Name, nameWidth)

			var preview string
			if entry.Entry.Type == "spinner" {
				frame := m.anim.currentFrame(entry.Entry.ID, entry.Entry.Frames)
				preview = truncateToWidth(frame, previewColWidth)
			} else if entry.Entry.Indeterminate != nil && *entry.Entry.Indeterminate != "" {
				preview = renderIndeterminate(*entry.Entry.Indeterminate, previewColWidth, m.anim.currentOffset(entry.Entry.ID))
			} else {
				preview = renderProgressBar(entry.Entry.Characters, entry.Entry.Phases, m.anim.currentProgressPct(entry.Entry.ID), previewColWidth)
			}

			selected := idx == m.entryCursor[g.name]
			prefix := "  "
			style := normalItemStyle
			if selected {
				style = selectedItemStyle
				if m.focusPane == listPaneEntries {
					prefix = "> "
				} else {
					prefix = "* "
				}
			}

			nameStr := style.Render(name)
			padding := entryWidth - lipgloss.Width(prefix) - lipgloss.Width(nameStr) - lipgloss.Width(preview)
			if padding < 1 {
				padding = 1
			}

			line := prefix + nameStr + strings.Repeat(" ", padding) + preview
			line = lipgloss.NewStyle().Width(entryWidth).MaxWidth(entryWidth).Render(line)
			entryLines = append(entryLines, line)
		}
	}

	var b strings.Builder
	sep := listDividerStyle.Render("|")
	b.WriteString(groupsHeader + sep + entriesHeader)
	for i := 0; i < rows; i++ {
		b.WriteString("\n")
		b.WriteString(groupLines[i] + sep + entryLines[i])
	}

	return b.String()
}

func (m *listModel) visibleRows() int {
	rows := m.height - 1
	if rows < 0 {
		rows = 0
	}
	return rows
}

func (m *listModel) scrollRows() int {
	rows := m.visibleRows()
	if rows < 1 {
		rows = 1
	}
	return rows
}

func (m *listModel) updatePagers() {
	rows := m.scrollRows()
	m.groupPager.PerPage = rows
	m.groupPager.SetTotalPages(len(m.visibleGroups))
	if rows > 0 {
		m.groupPager.Page = m.groupOffset / rows
	}

	g := m.currentGroup()
	entryTotal := 0
	entryOffset := 0
	if g != nil {
		entryTotal = len(g.entries)
		entryOffset = m.entryOffset[g.name]
	}
	m.entryPager.PerPage = rows
	m.entryPager.SetTotalPages(entryTotal)
	if rows > 0 {
		m.entryPager.Page = entryOffset / rows
	}
}

func (m *listModel) columnWidths() (int, int) {
	if m.width <= 2 {
		return 1, 1
	}
	groupWidth := m.width * 36 / 100
	if groupWidth < 18 {
		groupWidth = 18
	}
	if groupWidth > m.width-20 {
		groupWidth = m.width - 20
	}
	if groupWidth < 1 {
		groupWidth = 1
	}
	entryWidth := m.width - groupWidth - 1
	if entryWidth < 1 {
		entryWidth = 1
	}
	return groupWidth, entryWidth
}

func (m *listModel) isGroupColumn(x int) bool {
	groupWidth, _ := m.columnWidths()
	return x >= 0 && x < groupWidth
}

func (m *listModel) clickGroupRow(y int) bool {
	idx := m.groupOffset + y
	if idx < 0 || idx >= len(m.visibleGroups) {
		return false
	}
	m.groupCursor = idx
	m.ensureCurrentGroupState()
	m.ensureGroupVisible()
	m.ensureEntryVisible()
	return true
}

func (m *listModel) clickEntryRow(y int) bool {
	g := m.currentGroup()
	if g == nil {
		return false
	}
	idx := m.entryOffset[g.name] + y
	if idx < 0 || idx >= len(g.entries) {
		return false
	}
	m.entryCursor[g.name] = idx
	m.ensureEntryVisible()
	return true
}

func (m *listModel) scrollGroup(delta int) {
	steps := delta
	if steps < 0 {
		steps = -steps
	}
	if steps > 3 {
		steps = 3
	}
	for i := 0; i < steps; i++ {
		if delta > 0 {
			m.moveGroupDown()
		} else if delta < 0 {
			m.moveGroupUp()
		}
	}
}

func (m *listModel) scrollEntry(delta int) {
	steps := delta
	if steps < 0 {
		steps = -steps
	}
	if steps > 3 {
		steps = 3
	}
	for i := 0; i < steps; i++ {
		if delta > 0 {
			m.moveEntryDown()
		} else if delta < 0 {
			m.moveEntryUp()
		}
	}
}

func truncateWithEllipsis(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}
	if lipgloss.Width(s) <= maxWidth {
		return s
	}
	if maxWidth <= 3 {
		return strings.Repeat(".", maxWidth)
	}

	targetWidth := maxWidth - 3
	truncated := truncateToWidth(s, targetWidth)
	return truncated + "..."
}

func truncateToWidth(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}
	if lipgloss.Width(s) <= maxWidth {
		return s
	}

	var b strings.Builder
	used := 0
	for _, r := range s {
		rw := lipgloss.Width(string(r))
		if used+rw > maxWidth {
			break
		}
		b.WriteRune(r)
		used += rw
	}
	return b.String()
}
