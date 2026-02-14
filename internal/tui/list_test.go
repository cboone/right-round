package tui

import (
	"fmt"
	"testing"

	"github.com/cboone/right-round/internal/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeTestGroups() []data.Group {
	return []data.Group{
		{
			Name: "beta",
			Type: "spinner",
			Entries: []data.EntryEnvelope{
				{Entry: data.Entry{ID: "a/1", Name: "alpha one", Type: "spinner", Frames: []string{"a"}}},
				{Entry: data.Entry{ID: "a/2", Name: "alpha two", Type: "spinner", Frames: []string{"b"}}},
			},
		},
		{
			Name: "alpha",
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
	m.setSize(60, 20)

	require.NotEmpty(t, m.visibleGroups)
	assert.Equal(t, "b/1", m.selectedID())
	assert.Equal(t, "alpha", m.selectedGroupName())

	// Move down in current group
	m.moveGroupDown()
	assert.Equal(t, "a/1", m.selectedID())
	assert.Equal(t, "beta", m.selectedGroupName())

	m.moveEntryDown()
	assert.Equal(t, "a/2", m.selectedID())

	// Previous group's entry cursor is preserved
	m.moveGroupUp()
	assert.Equal(t, "b/1", m.selectedID())
}

func TestListModel_GoToTopBottom(t *testing.T) {
	anim := newAnimEngine()
	m := newListModel(makeTestGroups(), anim)
	m.setSize(60, 20)

	m.goGroupBottom()
	assert.Equal(t, "a/1", m.selectedID())

	m.goGroupTop()
	assert.Equal(t, "b/1", m.selectedID())
}

func TestListModel_Filter(t *testing.T) {
	anim := newAnimEngine()
	m := newListModel(makeTestGroups(), anim)
	m.setSize(60, 20)

	// Filter by "beta"
	m.setFilter("beta")
	require.NotNil(t, m.selectedEntry())
	assert.Equal(t, "b/1", m.selectedID())

	assert.Len(t, m.visibleGroups, 1)
	assert.Equal(t, "alpha", m.visibleGroups[0].name)

	// Clear filter
	m.setFilter("")
	assert.Len(t, m.visibleGroups, 2)
}

func TestListModel_FilterByID(t *testing.T) {
	anim := newAnimEngine()
	m := newListModel(makeTestGroups(), anim)
	m.setSize(60, 20)

	m.setFilter("a/2")
	require.NotNil(t, m.selectedEntry())
	assert.Equal(t, "a/2", m.selectedID())
}

func TestListModel_FilterPreservesSelection(t *testing.T) {
	anim := newAnimEngine()
	m := newListModel(makeTestGroups(), anim)
	m.setSize(60, 20)

	m.moveGroupDown()
	m.moveEntryDown()
	assert.Equal(t, "a/2", m.selectedID())

	// Filter that still includes the selected entry
	m.setFilter("alpha")
	assert.Equal(t, "a/2", m.selectedID())
}

func TestListModel_DefaultSortIsAlphabetical(t *testing.T) {
	anim := newAnimEngine()
	m := newListModel(makeTestGroups(), anim)
	m.setSize(60, 20)

	require.Len(t, m.visibleGroups, 2)
	assert.Equal(t, "alpha", m.visibleGroups[0].name)
	assert.Equal(t, "beta", m.visibleGroups[1].name)
}

func TestListModel_FilterNoMatches(t *testing.T) {
	anim := newAnimEngine()
	m := newListModel(makeTestGroups(), anim)
	m.setSize(60, 20)

	m.setFilter("zzz_no_match")
	assert.Empty(t, m.visibleGroups)
	assert.Nil(t, m.selectedEntry())

	// View should show "No matches"
	view := m.view()
	assert.Contains(t, view, "No matches")
}

func TestListModel_SetGroups(t *testing.T) {
	anim := newAnimEngine()
	m := newListModel(makeTestGroups(), anim)
	m.setSize(60, 20)

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
	m.setSize(60, 20)

	ids := m.visibleEntryIDs()
	assert.Contains(t, ids, "b/1")
	assert.NotContains(t, ids, "a/1")
}

func TestListModel_GroupScrollOffset(t *testing.T) {
	groups := make([]data.Group, 0, 12)
	for i := 0; i < 12; i++ {
		groups = append(groups, data.Group{
			Name:    fmt.Sprintf("group-%02d", i),
			Type:    "spinner",
			Entries: []data.EntryEnvelope{{Entry: data.Entry{ID: fmt.Sprintf("g/%d", i), Name: "item", Type: "spinner", Frames: []string{"x"}}}},
		})
	}

	anim := newAnimEngine()
	m := newListModel(groups, anim)
	m.setSize(60, 3)

	m.goGroupBottom()
	assert.Greater(t, m.groupOffset, 0)
}

func TestListModel_EntryScrollOffsetIndependentFromGroups(t *testing.T) {
	groups := []data.Group{
		{
			Name: "alpha",
			Type: "spinner",
			Entries: []data.EntryEnvelope{
				{Entry: data.Entry{ID: "a/1", Name: "a1", Type: "spinner", Frames: []string{"x"}}},
				{Entry: data.Entry{ID: "a/2", Name: "a2", Type: "spinner", Frames: []string{"x"}}},
				{Entry: data.Entry{ID: "a/3", Name: "a3", Type: "spinner", Frames: []string{"x"}}},
				{Entry: data.Entry{ID: "a/4", Name: "a4", Type: "spinner", Frames: []string{"x"}}},
			},
		},
		{
			Name:    "beta",
			Type:    "spinner",
			Entries: []data.EntryEnvelope{{Entry: data.Entry{ID: "b/1", Name: "b1", Type: "spinner", Frames: []string{"x"}}}},
		},
	}

	anim := newAnimEngine()
	m := newListModel(groups, anim)
	m.setSize(60, 2)

	m.goEntryBottom()
	assert.Greater(t, m.entryOffset["alpha"], 0)

	groupOffsetBefore := m.groupOffset
	m.moveEntryUp()
	assert.Equal(t, groupOffsetBefore, m.groupOffset)
}

func TestListModel_PageNavigation(t *testing.T) {
	anim := newAnimEngine()
	m := newListModel(makeTestGroups(), anim)
	m.setSize(60, 2)
	m.moveGroupDown()

	m.pageEntryDown()
	assert.NotEqual(t, 0, m.entryCursor["beta"])

	m.pageEntryUp()
	assert.Equal(t, "a/1", m.selectedID())
}

func TestListModel_View(t *testing.T) {
	anim := newAnimEngine()
	m := newListModel(makeTestGroups(), anim)
	m.setSize(60, 20)

	view := m.view()
	assert.Contains(t, view, "Categories")
	assert.Contains(t, view, "Entries")
	assert.Contains(t, view, "alpha")
	assert.Contains(t, view, "beta")
	assert.Contains(t, view, "beta one")
}
