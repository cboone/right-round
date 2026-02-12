package tui

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cboone/right-round/internal/data"
	"github.com/stretchr/testify/assert"
)

func makeTestGroupedEntries() *data.GroupedEntries {
	spinnerGroups := []data.Group{
		{
			Name: "braille",
			Type: "spinner",
			Entries: []data.EntryEnvelope{
				{Entry: data.Entry{ID: "s/1", Name: "spinner one", Type: "spinner", Group: "braille", Frames: []string{"a", "b"}}},
				{Entry: data.Entry{ID: "s/2", Name: "spinner two", Type: "spinner", Group: "braille", Frames: []string{"c", "d"}}},
			},
		},
	}
	barGroups := []data.Group{
		{
			Name: "ascii",
			Type: "progress_bar",
			Entries: []data.EntryEnvelope{
				{Entry: data.Entry{ID: "b/1", Name: "bar one", Type: "progress_bar", Group: "ascii",
					Characters: &data.BarCharacters{Fill: "#", Empty: "-"}}},
			},
		},
	}

	var all []data.EntryEnvelope
	for _, g := range spinnerGroups {
		all = append(all, g.Entries...)
	}
	for _, g := range barGroups {
		all = append(all, g.Entries...)
	}

	return &data.GroupedEntries{
		SpinnerGroups:     spinnerGroups,
		ProgressBarGroups: barGroups,
		AllEntries:        all,
	}
}

func TestNew_DefaultsToSpinners(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "", "")
	assert.Equal(t, tabSpinners, m.tab)
	assert.Equal(t, "s/1", m.list.selectedID())
}

func TestNew_TypeLockProgressBar(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "progress_bar", "")
	assert.Equal(t, tabProgressBars, m.tab)
	assert.Equal(t, "b/1", m.list.selectedID())
}

func TestNew_InitialGroup(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "", "braille")
	assert.Equal(t, "s/1", m.list.selectedID())
}

func TestModel_TabSwitching(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "", "")
	m.width = 120
	m.height = 40
	m.updateLayout()

	assert.Equal(t, tabSpinners, m.tab)

	// Press tab
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(Model)
	assert.Equal(t, tabProgressBars, m.tab)
	assert.Equal(t, "b/1", m.list.selectedID())

	// Tab again
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(Model)
	assert.Equal(t, tabSpinners, m.tab)
}

func TestModel_TabLocked(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "spinner", "")
	m.width = 120
	m.height = 40
	m.updateLayout()

	// Tab should not switch when locked
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(Model)
	assert.Equal(t, tabSpinners, m.tab)
}

func TestModel_Navigation(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "", "")
	m.width = 120
	m.height = 40
	m.updateLayout()

	assert.Equal(t, "s/1", m.list.selectedID())

	// j to move down
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(Model)
	assert.Equal(t, "s/2", m.list.selectedID())

	// k to move up
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = updated.(Model)
	assert.Equal(t, "s/1", m.list.selectedID())
}

func TestModel_SearchMode(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "", "")
	m.width = 120
	m.height = 40
	m.updateLayout()

	// Enter search mode with /
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m = updated.(Model)
	assert.True(t, m.filtering)

	// Type "two"
	for _, ch := range "two" {
		updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		m = updated.(Model)
	}
	assert.Equal(t, "two", m.filterInput)
	assert.Equal(t, "s/2", m.list.selectedID())

	// Esc exits search and clears filter
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = updated.(Model)
	assert.False(t, m.filtering)
	assert.Empty(t, m.filterInput)
}

func TestModel_SearchEnterConfirms(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "", "")
	m.width = 120
	m.height = 40
	m.updateLayout()

	// Enter search mode
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m = updated.(Model)

	// Type and confirm with enter
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	m = updated.(Model)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(Model)
	assert.False(t, m.filtering)
	// Filter should still be active (not cleared)
}

func TestModel_NarrowLayoutEnterExpands(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "", "")
	m.width = 80 // narrow
	m.height = 40
	m.updateLayout()

	assert.Equal(t, focusList, m.focus)

	// Enter expands detail
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(Model)
	assert.Equal(t, focusDetail, m.focus)

	// Esc goes back
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = updated.(Model)
	assert.Equal(t, focusList, m.focus)
}

func TestModel_WideLayoutNoFocusSwitch(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "", "")
	m.width = 120 // wide
	m.height = 40
	m.updateLayout()

	// In wide mode, enter shouldn't switch focus
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(Model)
	assert.Equal(t, focusList, m.focus)
}

func TestModel_HelpToggle(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "", "")
	m.width = 120
	m.height = 40
	m.updateLayout()

	assert.False(t, m.showFullHelp)

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	m = updated.(Model)
	assert.True(t, m.showFullHelp)

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	m = updated.(Model)
	assert.False(t, m.showFullHelp)
}

func TestModel_AnimTickMsg(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "", "")
	m.width = 120
	m.height = 40
	m.updateLayout()

	now := time.Now()
	m.lastTick = now.Add(-100 * time.Millisecond)

	updated, cmd := m.Update(animTickMsg(now))
	m = updated.(Model)

	// Should produce another tick command
	assert.NotNil(t, cmd)
}

func TestModel_WindowSizeMsg(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "", "")

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 150, Height: 50})
	m = updated.(Model)

	assert.Equal(t, 150, m.width)
	assert.Equal(t, 50, m.height)
}

func TestModel_QuitKey(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "", "")
	m.width = 120
	m.height = 40
	m.updateLayout()

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	assert.NotNil(t, cmd)
}

func TestModel_View(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "", "")
	m.width = 120
	m.height = 40
	m.updateLayout()

	view := m.View()
	assert.Contains(t, view, "Spinners")
	assert.Contains(t, view, "Progress Bars")
}

func TestModel_ViewNarrow(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "", "")
	m.width = 80
	m.height = 40
	m.updateLayout()

	view := m.View()
	assert.Contains(t, view, "spinner one")
}

func TestModel_StatusMessage(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "", "")
	m.width = 120
	m.height = 40
	m.updateLayout()

	m.statusMsg = "Copied to clipboard!"
	m.statusExpiry = time.Now().Add(2 * time.Second)

	view := m.View()
	assert.Contains(t, view, "Copied to clipboard!")
}

func TestModel_Init(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "", "")
	cmd := m.Init()
	assert.NotNil(t, cmd)
}
