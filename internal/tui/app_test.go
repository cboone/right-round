package tui

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	rightround "github.com/cboone/right-round"
	"github.com/cboone/right-round/internal/data"
	"github.com/stretchr/testify/assert"
)

func keyRune(r rune) tea.KeyPressMsg {
	return tea.KeyPressMsg{Code: r, Text: string(r)}
}

func keyCode(code rune) tea.KeyPressMsg {
	return tea.KeyPressMsg{Code: code}
}

func keyCtrlRune(r rune) tea.KeyPressMsg {
	return tea.KeyPressMsg{Code: r, Mod: tea.ModCtrl}
}

func mouseClick(x, y int, button tea.MouseButton) tea.MouseClickMsg {
	return tea.MouseClickMsg{X: x, Y: y, Button: button}
}

func mouseRelease(x, y int, button tea.MouseButton) tea.MouseReleaseMsg {
	return tea.MouseReleaseMsg{X: x, Y: y, Button: button}
}

func mouseWheel(x, y int, button tea.MouseButton) tea.MouseWheelMsg {
	return tea.MouseWheelMsg{X: x, Y: y, Button: button}
}

func viewString(m Model) string {
	return fmt.Sprint(m.View().Content)
}

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
	assert.Contains(t, m.filterBox.AvailableSuggestions(), "braille")
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
	updated, _ := m.Update(keyCode(tea.KeyTab))
	m = updated.(Model)
	assert.Equal(t, tabProgressBars, m.tab)
	assert.Equal(t, "b/1", m.list.selectedID())

	// Tab again
	updated, _ = m.Update(keyCode(tea.KeyTab))
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
	updated, _ := m.Update(keyCode(tea.KeyTab))
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
	updated, _ := m.Update(keyRune('j'))
	m = updated.(Model)
	assert.Equal(t, "s/2", m.list.selectedID())

	// k to move up
	updated, _ = m.Update(keyRune('k'))
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
	updated, _ := m.Update(keyRune('/'))
	m = updated.(Model)
	assert.True(t, m.filtering)

	// Type "two"
	for _, ch := range "two" {
		updated, _ = m.Update(keyRune(ch))
		m = updated.(Model)
	}
	assert.Equal(t, "two", m.filterInput)
	assert.Equal(t, "s/2", m.list.selectedID())

	// Esc exits search and clears filter
	updated, _ = m.Update(keyCode(tea.KeyEsc))
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
	updated, _ := m.Update(keyRune('/'))
	m = updated.(Model)

	// Type and confirm with enter
	updated, _ = m.Update(keyRune('t'))
	m = updated.(Model)
	updated, _ = m.Update(keyCode(tea.KeyEnter))
	m = updated.(Model)
	assert.False(t, m.filtering)
	// Filter should still be active (not cleared)
}

func TestModel_CtrlCQuitsWhileFiltering(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "", "")
	m.width = 120
	m.height = 40
	m.updateLayout()

	// Enter search mode with /
	updated, _ := m.Update(keyRune('/'))
	m = updated.(Model)
	assert.True(t, m.filtering)

	// ctrl+c should quit even while filtering
	_, cmd := m.Update(keyCtrlRune('c'))
	assert.NotNil(t, cmd)
}

func TestModel_NarrowLayoutEnterExpands(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "", "")
	m.width = 80 // narrow
	m.height = 40
	m.updateLayout()

	assert.Equal(t, focusEntries, m.focus)

	// Enter expands detail
	updated, _ := m.Update(keyCode(tea.KeyEnter))
	m = updated.(Model)
	assert.Equal(t, focusDetail, m.focus)

	// Esc goes back
	updated, _ = m.Update(keyCode(tea.KeyEsc))
	m = updated.(Model)
	assert.Equal(t, focusEntries, m.focus)
}

func TestModel_WideLayoutNoFocusSwitch(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "", "")
	m.width = 120 // wide
	m.height = 40
	m.updateLayout()

	// In wide mode, enter shouldn't switch focus
	updated, _ := m.Update(keyCode(tea.KeyEnter))
	m = updated.(Model)
	assert.Equal(t, focusEntries, m.focus)
}

func TestModel_HelpToggle(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "", "")
	m.width = 120
	m.height = 40
	m.updateLayout()

	assert.False(t, m.showFullHelp)

	updated, _ := m.Update(keyRune('?'))
	m = updated.(Model)
	assert.True(t, m.showFullHelp)

	updated, _ = m.Update(keyRune('?'))
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

	_, cmd := m.Update(keyRune('q'))
	assert.NotNil(t, cmd)
}

func TestModel_View(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "", "")
	m.width = 120
	m.height = 40
	m.updateLayout()

	view := viewString(m)
	assert.Contains(t, view, "Spinners")
	assert.Contains(t, view, "Progress Bars")
}

func TestModel_ViewNarrow(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "", "")
	m.width = 80
	m.height = 40
	m.updateLayout()

	view := viewString(m)
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

	view := viewString(m)
	assert.Contains(t, view, "Copied to clipboard!")
}

func TestModel_SmallTerminalHeight(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "", "")

	// Simulate very small terminal (height=2 yields contentHeight=2-2-1=-1 unclamped)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 2})
	m = updated.(Model)

	assert.Equal(t, 2, m.height)
	// Layout should not panic and list height should be at least 1
	assert.GreaterOrEqual(t, m.list.height, 1)
}

func TestModel_Init(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "", "")
	cmd := m.Init()
	assert.NotNil(t, cmd)
}

func makeMouseTestGroupedEntries() *data.GroupedEntries {
	spinnerGroups := []data.Group{
		{
			Name:    "alpha",
			Type:    "spinner",
			Entries: []data.EntryEnvelope{{Entry: data.Entry{ID: "s/a", Name: "alpha spin", Type: "spinner", Group: "alpha", Frames: []string{"a", "b"}}}},
		},
		{
			Name:    "beta",
			Type:    "spinner",
			Entries: []data.EntryEnvelope{{Entry: data.Entry{ID: "s/b", Name: "beta spin", Type: "spinner", Group: "beta", Frames: []string{"c", "d"}}}},
		},
	}
	barGroups := []data.Group{
		{
			Name:    "ascii",
			Type:    "progress_bar",
			Entries: []data.EntryEnvelope{{Entry: data.Entry{ID: "b/1", Name: "bar one", Type: "progress_bar", Group: "ascii", Characters: &data.BarCharacters{Fill: "#", Empty: "-"}}}},
		},
	}

	all := append([]data.EntryEnvelope{}, spinnerGroups[0].Entries...)
	all = append(all, spinnerGroups[1].Entries...)
	all = append(all, barGroups[0].Entries...)

	return &data.GroupedEntries{
		SpinnerGroups:     spinnerGroups,
		ProgressBarGroups: barGroups,
		AllEntries:        all,
	}
}

func TestModel_MouseTabSwitching(t *testing.T) {
	grouped := makeMouseTestGroupedEntries()
	m := New(grouped, "", "")
	m.width = 120
	m.height = 40
	m.updateLayout()
	spinnerWidth := lipgloss.Width(inactiveTabStyle.Render("Spinners"))
	barX := spinnerWidth + 1

	updated, _ := m.Update(mouseClick(barX, 0, tea.MouseLeft))
	m = updated.(Model)
	assert.Equal(t, tabProgressBars, m.tab)

	updated, _ = m.Update(mouseClick(1, 0, tea.MouseLeft))
	m = updated.(Model)
	assert.Equal(t, tabSpinners, m.tab)
}

func TestModel_MouseSelectsGroup(t *testing.T) {
	grouped := makeMouseTestGroupedEntries()
	m := New(grouped, "", "")
	m.width = 120
	m.height = 40
	m.updateLayout()

	groupWidth, _ := m.list.columnWidths()
	assert.Equal(t, "alpha", m.list.selectedGroupName())

	updated, _ := m.Update(mouseClick(groupWidth-2, 2, tea.MouseLeft))
	m = updated.(Model)

	assert.Equal(t, "beta", m.list.selectedGroupName())
	assert.Equal(t, focusGroups, m.focus)
}

func TestModel_MouseWheelScrollsEntriesWithoutChangingGroup(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "", "")
	m.width = 120
	m.height = 8
	m.updateLayout()

	_, entryWidth := m.list.columnWidths()
	entryX := m.list.width - entryWidth + 1
	groupBefore := m.list.selectedGroupName()

	updated, _ := m.Update(mouseWheel(entryX, 2, tea.MouseWheelDown))
	m = updated.(Model)

	assert.Equal(t, groupBefore, m.list.selectedGroupName())
}

func TestModel_MouseReleaseTriggersTabSwitch(t *testing.T) {
	grouped := makeMouseTestGroupedEntries()
	m := New(grouped, "", "")
	m.width = 120
	m.height = 40
	m.updateLayout()
	spinnerWidth := lipgloss.Width(inactiveTabStyle.Render("Spinners"))
	barX := spinnerWidth + 1

	updated, _ := m.Update(mouseRelease(barX, 0, tea.MouseLeft))
	m = updated.(Model)
	assert.Equal(t, tabProgressBars, m.tab)

	updated, _ = m.Update(mouseRelease(1, 0, tea.MouseLeft))
	m = updated.(Model)
	assert.Equal(t, tabSpinners, m.tab)
}

func TestModel_MouseWheelDeprecatedTypeScrollsDetail(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "", "")
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 30})
	m = updated.(Model)

	m.focus = focusDetail
	m.detail.viewport.SetContent(strings.Repeat("line\n", 200))
	before := m.detail.viewport.YOffset()

	updated, _ = m.Update(mouseWheel(m.list.width+2, 5, tea.MouseWheelDown))
	m = updated.(Model)

	assert.Greater(t, m.detail.viewport.YOffset(), before)
}

func TestModel_MouseWheelScrollsDetailFromRightPane(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "", "")
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 30})
	m = updated.(Model)

	m.detail.viewport.SetContent(strings.Repeat("line\n", 200))
	before := m.detail.viewport.YOffset()

	updated, _ = m.Update(mouseWheel(m.list.width+2, 6, tea.MouseWheelDown))
	m = updated.(Model)

	assert.Greater(t, m.detail.viewport.YOffset(), before)
}

func TestModel_ViewHeightStableAfterHelpToggle(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "", "")

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 30})
	m = updated.(Model)
	assert.Equal(t, 30, lipgloss.Height(viewString(m)))

	updated, _ = m.Update(keyRune('?'))
	m = updated.(Model)
	assert.Equal(t, 30, lipgloss.Height(viewString(m)))

	updated, _ = m.Update(keyRune('?'))
	m = updated.(Model)
	assert.Equal(t, 30, lipgloss.Height(viewString(m)))
}

func TestModel_ViewHeightStableAcrossMouseClick(t *testing.T) {
	grouped := makeMouseTestGroupedEntries()
	m := New(grouped, "", "")

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 30})
	m = updated.(Model)
	assert.Equal(t, 30, lipgloss.Height(viewString(m)))

	groupWidth, _ := m.list.columnWidths()
	updated, _ = m.Update(mouseClick(groupWidth-2, 3, tea.MouseLeft))
	m = updated.(Model)
	assert.Equal(t, 30, lipgloss.Height(viewString(m)))

	updated, _ = m.Update(mouseRelease(groupWidth-2, 3, tea.MouseLeft))
	m = updated.(Model)
	assert.Equal(t, 30, lipgloss.Height(viewString(m)))
}

func TestModel_ViewHeightStableAcrossRepeatedGroupClicks(t *testing.T) {
	grouped := makeMouseTestGroupedEntries()
	m := New(grouped, "", "")

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 78, Height: 24})
	m = updated.(Model)
	assert.Equal(t, 24, lipgloss.Height(viewString(m)))

	groupWidth, _ := m.list.columnWidths()
	for i := 0; i < 20; i++ {
		updated, _ = m.Update(mouseClick(groupWidth-2, 3, tea.MouseLeft))
		m = updated.(Model)
		assert.Equal(t, 24, lipgloss.Height(viewString(m)))

		updated, _ = m.Update(mouseRelease(groupWidth-2, 3, tea.MouseLeft))
		m = updated.(Model)
		assert.Equal(t, 24, lipgloss.Height(viewString(m)))
	}
}

func TestModel_ViewHeightStableForSymbolGroupAnimation(t *testing.T) {
	grouped, err := data.LoadCatalog(rightround.EmbeddedCatalogJSON())
	assert.NoError(t, err)

	m := New(grouped, "spinner", "symbol")
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 30})
	m = updated.(Model)
	baseline := lipgloss.Height(viewString(m))
	assert.GreaterOrEqual(t, baseline, 1)

	for i := 0; i < 60; i++ {
		next := m.lastTick.Add(16 * time.Millisecond)
		updated, _ = m.Update(animTickMsg(next))
		m = updated.(Model)
		assert.Equal(t, baseline, lipgloss.Height(viewString(m)))
	}
}

func TestModel_ContextualHelpByFocus(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "", "")
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 30})
	m = updated.(Model)

	km := m.currentHelpKeyMap()
	assert.True(t, hasBindingKey(km.ShortHelp(), "/"))
	assert.True(t, hasBindingKey(km.ShortHelp(), "c"))

	m.focus = focusGroups
	km = m.currentHelpKeyMap()
	assert.True(t, hasBindingKey(km.ShortHelp(), "["))
	assert.True(t, hasBindingKey(km.ShortHelp(), "/"))

	m.focus = focusDetail
	km = m.currentHelpKeyMap()
	assert.True(t, hasBindingKey(km.ShortHelp(), "c"))
	assert.True(t, hasBindingKey(km.ShortHelp(), "esc"))
}

func TestModel_ContextualHelpHonorsTypeLock(t *testing.T) {
	grouped := makeTestGroupedEntries()
	m := New(grouped, "spinner", "")
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 30})
	m = updated.(Model)

	km := m.currentHelpKeyMap()
	for _, group := range km.FullHelp() {
		assert.False(t, hasBindingKey(group, "tab"))
	}
}

func hasBindingKey(bindings []key.Binding, target string) bool {
	for _, b := range bindings {
		for _, k := range b.Keys() {
			if k == target {
				return true
			}
		}
	}
	return false
}
