package tui

import (
	"testing"

	"github.com/cboone/right-round/internal/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeTestGroups() []data.Group {
	return []data.Group{
		{
			Name: "alpha",
			Type: "spinner",
			Entries: []data.EntryEnvelope{
				{Entry: data.Entry{ID: "a/1", Name: "alpha one", Type: "spinner", Frames: []string{"a"}}},
				{Entry: data.Entry{ID: "a/2", Name: "alpha two", Type: "spinner", Frames: []string{"b"}}},
			},
		},
		{
			Name: "beta",
			Type: "spinner",
			Entries: []data.EntryEnvelope{
				{Entry: data.Entry{ID: "b/1", Name: "beta one", Type: "spinner", Frames: []string{"c"}}},
			},
		},
	}
}

func TestListModel_CursorNavigation(t *testing.T) {
	anim := newAnimEngine()
	m := newListModel(makeTestGroups(), anim)
	m.width = 60
	m.height = 20

	// Initial position should be on first entry (past header)
	require.NotEmpty(t, m.rows)
	assert.False(t, m.rows[m.cursor].isHeader)
	assert.Equal(t, "a/1", m.selectedID())

	// Move down
	m.moveDown()
	assert.Equal(t, "a/2", m.selectedID())

	// Move down past group boundary (should skip header)
	m.moveDown()
	assert.Equal(t, "b/1", m.selectedID())

	// Move up past group boundary
	m.moveUp()
	assert.Equal(t, "a/2", m.selectedID())
}

func TestListModel_GoToTopBottom(t *testing.T) {
	anim := newAnimEngine()
	m := newListModel(makeTestGroups(), anim)
	m.width = 60
	m.height = 20

	m.goToBottom()
	assert.Equal(t, "b/1", m.selectedID())

	m.goToTop()
	assert.Equal(t, "a/1", m.selectedID())
}

func TestListModel_Filter(t *testing.T) {
	anim := newAnimEngine()
	m := newListModel(makeTestGroups(), anim)
	m.width = 60
	m.height = 20

	// Filter by "beta"
	m.setFilter("beta")
	require.NotNil(t, m.selectedEntry())
	assert.Equal(t, "b/1", m.selectedID())

	// Only beta group should be visible
	entryCount := 0
	for _, row := range m.rows {
		if !row.isHeader {
			entryCount++
		}
	}
	assert.Equal(t, 1, entryCount)

	// Clear filter
	m.setFilter("")
	entryCount = 0
	for _, row := range m.rows {
		if !row.isHeader {
			entryCount++
		}
	}
	assert.Equal(t, 3, entryCount)
}

func TestListModel_FilterByID(t *testing.T) {
	anim := newAnimEngine()
	m := newListModel(makeTestGroups(), anim)
	m.width = 60
	m.height = 20

	m.setFilter("a/2")
	require.NotNil(t, m.selectedEntry())
	assert.Equal(t, "a/2", m.selectedID())
}

func TestListModel_FilterPreservesSelection(t *testing.T) {
	anim := newAnimEngine()
	m := newListModel(makeTestGroups(), anim)
	m.width = 60
	m.height = 20

	m.moveDown()
	assert.Equal(t, "a/2", m.selectedID())

	// Filter that still includes the selected entry
	m.setFilter("alpha")
	assert.Equal(t, "a/2", m.selectedID())
}

func TestListModel_FilterNoMatches(t *testing.T) {
	anim := newAnimEngine()
	m := newListModel(makeTestGroups(), anim)
	m.width = 60
	m.height = 20

	m.setFilter("zzz_no_match")
	assert.Empty(t, m.rows)
	assert.Nil(t, m.selectedEntry())

	// View should show "No matches"
	view := m.view()
	assert.Contains(t, view, "No matches")
}

func TestListModel_SetGroups(t *testing.T) {
	anim := newAnimEngine()
	m := newListModel(makeTestGroups(), anim)
	m.width = 60
	m.height = 20

	newGroups := []data.Group{
		{
			Name: "gamma",
			Type: "progress_bar",
			Entries: []data.EntryEnvelope{
				{Entry: data.Entry{ID: "g/1", Name: "gamma one", Type: "progress_bar",
					Characters: &data.BarCharacters{Fill: "#", Empty: "-"}}},
			},
		},
	}

	m.setGroups(newGroups)
	assert.Equal(t, "g/1", m.selectedID())
}

func TestListModel_VisibleEntryIDs(t *testing.T) {
	anim := newAnimEngine()
	m := newListModel(makeTestGroups(), anim)
	m.width = 60
	m.height = 20

	ids := m.visibleEntryIDs()
	assert.Contains(t, ids, "a/1")
	assert.Contains(t, ids, "a/2")
	assert.Contains(t, ids, "b/1")
}

func TestListModel_ScrollOffset(t *testing.T) {
	anim := newAnimEngine()
	m := newListModel(makeTestGroups(), anim)
	m.width = 60
	m.height = 2 // Very small to force scrolling

	m.goToBottom()
	assert.Greater(t, m.offset, 0)
}

func TestListModel_PageNavigation(t *testing.T) {
	anim := newAnimEngine()
	m := newListModel(makeTestGroups(), anim)
	m.width = 60
	m.height = 2

	m.pageDown()
	// Should move cursor forward by height
	assert.NotEqual(t, 0, m.cursor)

	m.pageUp()
	assert.Equal(t, "a/1", m.selectedID())
}

func TestListModel_View(t *testing.T) {
	anim := newAnimEngine()
	m := newListModel(makeTestGroups(), anim)
	m.width = 60
	m.height = 20

	view := m.view()
	assert.Contains(t, view, "ALPHA")
	assert.Contains(t, view, "BETA")
	assert.Contains(t, view, "alpha one")
}
