package tui

import (
	"fmt"
	"strings"

	"github.com/cboone/right-round/internal/data"
	"github.com/charmbracelet/lipgloss"
)

// listRow represents a single renderable row in the list panel.
type listRow struct {
	isHeader bool
	header   string              // non-empty for group headers
	count    int                 // entry count for group headers
	entry    *data.EntryEnvelope // non-nil for entry rows
}

// listModel manages the grouped list panel.
type listModel struct {
	groups []data.Group
	rows   []listRow
	cursor int
	offset int
	height int
	width  int
	filter string

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
	if len(m.rows) == 0 {
		m.cursor = 0
		return
	}
	// Skip headers in the given direction
	for m.cursor >= 0 && m.cursor < len(m.rows) && m.rows[m.cursor].isHeader {
		m.cursor += dir
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.cursor >= len(m.rows) {
		m.cursor = len(m.rows) - 1
	}
	// If still on a header after clamping (hit boundary), try opposite direction
	if m.rows[m.cursor].isHeader {
		for m.cursor >= 0 && m.cursor < len(m.rows) && m.rows[m.cursor].isHeader {
			m.cursor -= dir
		}
		if m.cursor < 0 {
			m.cursor = 0
		}
		if m.cursor >= len(m.rows) {
			m.cursor = len(m.rows) - 1
		}
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
	if nameWidth < 1 {
		nameWidth = 1
	}

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
		name = truncateWithEllipsis(name, nameWidth)

		var preview string
		if entry.Entry.Type == "spinner" {
			frame := m.anim.currentFrame(entry.Entry.ID, entry.Entry.Frames)
			frame = truncateToWidth(frame, previewColWidth)
			preview = frame
		} else if entry.Entry.Indeterminate != nil && *entry.Entry.Indeterminate != "" {
			preview = renderIndeterminate(*entry.Entry.Indeterminate, previewColWidth, m.anim.currentOffset(entry.Entry.ID))
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
