package tui

import (
	"fmt"
	"strings"

	"github.com/cboone/right-round/internal/data"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

type detailModel struct {
	viewport viewport.Model
	entry    *data.EntryEnvelope
	anim     *animEngine
	width    int
	height   int
}

func newDetailModel(anim *animEngine) detailModel {
	vp := viewport.New(0, 0)
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 2
	return detailModel{
		viewport: vp,
		anim:     anim,
	}
}

func (m *detailModel) setEntry(entry *data.EntryEnvelope) {
	// Only reset scroll position when the entry actually changes
	oldID := ""
	newID := ""
	if m.entry != nil {
		oldID = m.entry.Entry.ID
	}
	if entry != nil {
		newID = entry.Entry.ID
	}
	m.entry = entry
	if newID != oldID {
		m.viewport.GotoTop()
	}
}

func (m *detailModel) setSize(width, height int) {
	m.width = width
	m.height = height
	vpWidth := width - 2 // border padding
	vpHeight := height - 2
	if vpWidth < 0 {
		vpWidth = 0
	}
	if vpHeight < 0 {
		vpHeight = 0
	}
	m.viewport.Width = vpWidth
	m.viewport.Height = vpHeight
}

func (m *detailModel) updateContent() {
	if m.entry == nil {
		m.viewport.SetContent(helpStyle.Render("No entry selected"))
		return
	}

	e := &m.entry.Entry
	contentWidth := m.width - 4
	if contentWidth < 20 {
		contentWidth = 20
	}

	var b strings.Builder
	truncateForPrefix := func(prefix string, s string) string {
		maxWidth := contentWidth - lipgloss.Width(prefix)
		if maxWidth < 1 {
			maxWidth = 1
		}
		return truncateWithEllipsis(s, maxWidth)
	}
	writePrefixed := func(prefix string, s string) {
		b.WriteString(prefix + truncateForPrefix(prefix, s) + "\n")
	}
	section := func(title string) {
		b.WriteString("\n" + detailLabelStyle.Render(title) + "\n")
	}

	b.WriteString(detailLabelStyle.Render(e.Name) + "\n")
	b.WriteString(helpStyle.Render(e.ID) + "\n")

	section("Preview")
	if e.Type == "spinner" {
		frame := m.anim.currentFrame(e.ID, e.Frames)
		writePrefixed("  live: ", frame)
		interval := data.DefaultIntervalMS
		if e.IntervalMS != nil {
			interval = *e.IntervalMS
		}
		b.WriteString(fmt.Sprintf("  frames: %d  interval: %dms\n", len(e.Frames), interval))
		limit := len(e.Frames)
		if limit > 8 {
			limit = 8
		}
		writePrefixed("  sample: ", strings.Join(e.Frames[:limit], " "))
	} else {
		barWidth := contentWidth - 8
		if barWidth > 40 {
			barWidth = 40
		}
		if barWidth < 8 {
			barWidth = 8
		}
		livePct := m.anim.currentProgressPct(e.ID)
		live := renderProgressBar(e.Characters, e.Phases, livePct, barWidth)
		b.WriteString(fmt.Sprintf("  live %3.0f%% %s\n", livePct*100, live))
		for _, pct := range []float64{0.0, 0.5, 1.0} {
			bar := renderProgressBar(e.Characters, e.Phases, pct, barWidth)
			b.WriteString(fmt.Sprintf("  %3.0f%% %s\n", pct*100, bar))
		}

		if e.Indeterminate != nil && *e.Indeterminate != "" {
			preview := renderIndeterminate(*e.Indeterminate, barWidth, m.anim.currentOffset(e.ID))
			writePrefixed("  indeterminate: ", preview)
		}
	}

	section("Essentials")
	writePrefixed("  type: ", e.Type)
	writePrefixed("  group: ", e.Group)
	writePrefixed("  source: ", e.Source.Collection)
	writePrefixed("  license: ", e.Source.License)

	section("Rendering")
	if e.Type == "progress_bar" {
		if e.Characters != nil {
			line := fmt.Sprintf("  chars: fill=%q empty=%q", e.Characters.Fill, e.Characters.Empty)
			if e.Characters.Head != nil {
				line += fmt.Sprintf(" head=%q", *e.Characters.Head)
			}
			if e.Characters.Start != nil {
				line += fmt.Sprintf(" start=%q", *e.Characters.Start)
			}
			if e.Characters.End != nil {
				line += fmt.Sprintf(" end=%q", *e.Characters.End)
			}
			writePrefixed("", line)
		} else {
			b.WriteString("  chars: none\n")
		}
		if len(e.Phases) > 0 {
			writePrefixed("  phases: ", strings.Join(e.Phases, " "))
		}
		if e.Indeterminate != nil && *e.Indeterminate != "" {
			writePrefixed("  pattern: ", *e.Indeterminate)
		}
	} else {
		interval := data.DefaultIntervalMS
		if e.IntervalMS != nil {
			interval = *e.IntervalMS
		}
		b.WriteString(fmt.Sprintf("  interval_ms: %d\n", interval))
	}

	section("Source")
	writePrefixed("  key: ", e.Source.OriginalKey)
	if e.Source.URL != "" {
		writePrefixed("  url: ", e.Source.URL)
	}

	if e.Notes != nil {
		section("Notes")
		writePrefixed("  ", *e.Notes)
	}

	m.viewport.SetContent(b.String())
}

func (m *detailModel) view() string {
	m.updateContent()
	w := m.width - 2
	h := m.height - 2
	if w < 0 {
		w = 0
	}
	if h < 0 {
		h = 0
	}
	return detailBorderStyle.
		Width(w).
		Height(h).
		Render(m.viewport.View())
}
