package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/cboone/right-round/internal/data"
)

// listRow represents a single renderable row in the list panel.
type listRow struct {
	isHeader bool
	header   string    // non-empty for group headers
	count    int       // entry count for group headers
	entry    *data.EntryEnvelope // non-nil for entry rows
}

// listModel manages the grouped list panel.
type listModel struct {
	groups    []data.Group
	rows      []listRow
	cursor    int
	offset    int
	height    int
	width     int
	filter    string
	filtering bool

	anim *animEngine
}

func newListModel(groups []data.Group, anim *animEngine) listModel {
	m := listModel{
		groups: groups,
		anim:   anim,
	}
	m.rows = m.buildRows(groups, "")
	m.moveToNextEntry(1)
	return m
}

func (m *listModel) buildRows(groups []data.Group, filter string) []listRow {
	filter = strings.ToLower(filter)
	var rows []listRow
	for i := range groups {
		g := &groups[i]
		var matching []*data.EntryEnvelope
		for j := range g.Entries {
			e := &g.Entries[j]
			if filter == "" || strings.Contains(strings.ToLower(e.Entry.Name), filter) || strings.Contains(strings.ToLower(e.Entry.ID), filter) {
				matching = append(matching, e)
			}
		}
		if len(matching) == 0 {
			continue
		}
		rows = append(rows, listRow{
			isHeader: true,
			header:   g.Name,
			count:    len(matching),
		})
		for _, e := range matching {
			rows = append(rows, listRow{entry: e})
		}
	}
	return rows
}

func (m *listModel) setFilter(filter string) {
	oldID := m.selectedID()
	m.filter = filter
	m.rows = m.buildRows(m.groups, filter)

	// Try to keep selection on same entry
	if oldID != "" {
		for i, row := range m.rows {
			if !row.isHeader && row.entry != nil && row.entry.Entry.ID == oldID {
				m.cursor = i
				m.ensureVisible()
				return
			}
		}
	}
	// Move to first entry row
	m.cursor = 0
	m.moveToNextEntry(1)
	m.ensureVisible()
}

func (m *listModel) setGroups(groups []data.Group) {
	m.groups = groups
	m.rows = m.buildRows(groups, m.filter)
	m.cursor = 0
	m.offset = 0
	m.moveToNextEntry(1)
}

func (m *listModel) selectedEntry() *data.EntryEnvelope {
	if m.cursor < 0 || m.cursor >= len(m.rows) {
		return nil
	}
	row := m.rows[m.cursor]
	if row.isHeader {
		return nil
	}
	return row.entry
}

func (m *listModel) selectedID() string {
	if e := m.selectedEntry(); e != nil {
		return e.Entry.ID
	}
	return ""
}

func (m *listModel) moveUp() {
	if m.cursor > 0 {
		m.cursor--
		if m.rows[m.cursor].isHeader {
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor++
			}
		}
		m.ensureVisible()
	}
}

func (m *listModel) moveDown() {
	if m.cursor < len(m.rows)-1 {
		m.cursor++
		if m.rows[m.cursor].isHeader {
			if m.cursor < len(m.rows)-1 {
				m.cursor++
			} else {
				m.cursor--
			}
		}
		m.ensureVisible()
	}
}

func (m *listModel) pageUp() {
	m.cursor -= m.height
	if m.cursor < 0 {
		m.cursor = 0
	}
	m.moveToNextEntry(1)
	m.ensureVisible()
}

func (m *listModel) pageDown() {
	m.cursor += m.height
	if m.cursor >= len(m.rows) {
		m.cursor = len(m.rows) - 1
	}
	m.moveToNextEntry(-1)
	m.ensureVisible()
}

func (m *listModel) goToTop() {
	m.cursor = 0
	m.moveToNextEntry(1)
	m.ensureVisible()
}

func (m *listModel) goToBottom() {
	m.cursor = len(m.rows) - 1
	m.moveToNextEntry(-1)
	m.ensureVisible()
}

func (m *listModel) moveToNextEntry(dir int) {
	for m.cursor >= 0 && m.cursor < len(m.rows) && m.rows[m.cursor].isHeader {
		m.cursor += dir
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.cursor >= len(m.rows) {
		m.cursor = len(m.rows) - 1
	}
}

func (m *listModel) ensureVisible() {
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+m.height {
		m.offset = m.cursor - m.height + 1
	}
	if m.offset < 0 {
		m.offset = 0
	}
}

// visibleEntryIDs returns the IDs of entries currently visible for animation.
func (m *listModel) visibleEntryIDs() []string {
	var ids []string
	end := m.offset + m.height
	if end > len(m.rows) {
		end = len(m.rows)
	}
	for i := m.offset; i < end; i++ {
		if !m.rows[i].isHeader && m.rows[i].entry != nil {
			ids = append(ids, m.rows[i].entry.Entry.ID)
		}
	}
	return ids
}

func (m *listModel) view() string {
	if len(m.rows) == 0 {
		return helpStyle.Render("  No matches")
	}

	previewColWidth := 8
	nameWidth := m.width - previewColWidth - 4 // 4 for cursor + padding

	var b strings.Builder
	end := m.offset + m.height
	if end > len(m.rows) {
		end = len(m.rows)
	}

	for i := m.offset; i < end; i++ {
		row := m.rows[i]
		if row.isHeader {
			header := fmt.Sprintf("%s (%d)", strings.ToUpper(row.header), row.count)
			b.WriteString(groupHeaderStyle.Render(header))
			b.WriteString("\n")
			continue
		}

		entry := row.entry
		selected := i == m.cursor
		name := entry.Entry.Name
		if lipgloss.Width(name) > nameWidth {
			name = name[:nameWidth-1] + "..."
		}

		var preview string
		if entry.Entry.Type == "spinner" {
			frame := m.anim.currentFrame(entry.Entry.ID, entry.Entry.Frames)
			if lipgloss.Width(frame) > previewColWidth {
				frame = frame[:previewColWidth]
			}
			preview = frame
		} else {
			preview = renderStaticBar(entry.Entry.Characters, entry.Entry.Phases, previewColWidth)
		}

		cursor := "  "
		style := normalItemStyle
		if selected {
			cursor = "> "
			style = selectedItemStyle
		}

		nameStr := style.Render(name)
		padding := m.width - lipgloss.Width(cursor) - lipgloss.Width(nameStr) - lipgloss.Width(preview)
		if padding < 1 {
			padding = 1
		}

		b.WriteString(cursor + nameStr + strings.Repeat(" ", padding) + preview + "\n")
	}

	return b.String()
}
