package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cboone/right-round/internal/data"
	"github.com/charmbracelet/bubbles/viewport"
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

	var b strings.Builder

	// Name and ID
	b.WriteString(detailLabelStyle.Render(e.Name) + "\n")
	b.WriteString(helpStyle.Render(e.ID) + "\n\n")

	// Type and Group
	b.WriteString(detailLabelStyle.Render("Type: ") + e.Type + "  ")
	b.WriteString(detailLabelStyle.Render("Group: ") + e.Group + "\n\n")

	// Live preview
	if e.Type == "spinner" {
		frame := m.anim.currentFrame(e.ID, e.Frames)
		b.WriteString(detailLabelStyle.Render("Preview: ") + frame + "\n")
		b.WriteString(detailLabelStyle.Render("Frames: ") + fmt.Sprintf("%d", len(e.Frames)) + "\n")
		interval := data.DefaultIntervalMS
		if e.IntervalMS != nil {
			interval = *e.IntervalMS
		}
		b.WriteString(detailLabelStyle.Render("Interval: ") + fmt.Sprintf("%dms", interval) + "\n")

		// Show all frames
		b.WriteString(detailLabelStyle.Render("All frames: ") + strings.Join(e.Frames, " ") + "\n")
	} else {
		// Progress bar preview
		barWidth := contentWidth
		if barWidth > 50 {
			barWidth = 50
		}
		b.WriteString(detailLabelStyle.Render("Preview:") + "\n")
		for _, pct := range []float64{0.0, 0.25, 0.5, 0.75, 1.0} {
			bar := renderProgressBar(e.Characters, e.Phases, pct, barWidth)
			b.WriteString(fmt.Sprintf("  %3.0f%% %s\n", pct*100, bar))
		}

		if e.Indeterminate != nil && *e.Indeterminate != "" {
			preview := renderIndeterminate(*e.Indeterminate, barWidth, m.anim.currentOffset(e.ID))
			b.WriteString("\n" + detailLabelStyle.Render("Indeterminate preview: ") + preview + "\n")
		}

		if e.Characters != nil {
			b.WriteString("\n" + detailLabelStyle.Render("Characters:") + "\n")
			b.WriteString(fmt.Sprintf("  fill: %q  empty: %q", e.Characters.Fill, e.Characters.Empty))
			if e.Characters.Head != nil {
				b.WriteString(fmt.Sprintf("  head: %q", *e.Characters.Head))
			}
			if e.Characters.Start != nil {
				b.WriteString(fmt.Sprintf("  start: %q", *e.Characters.Start))
			}
			if e.Characters.End != nil {
				b.WriteString(fmt.Sprintf("  end: %q", *e.Characters.End))
			}
			b.WriteString("\n")
		}

		if len(e.Phases) > 0 {
			b.WriteString(detailLabelStyle.Render("Phases: ") + strings.Join(e.Phases, " ") + "\n")
		}

		if e.Indeterminate != nil {
			b.WriteString(detailLabelStyle.Render("Indeterminate pattern: ") + *e.Indeterminate + "\n")
		}
	}

	if len(e.CompletionStates) > 0 {
		b.WriteString("\n" + detailLabelStyle.Render("Completion states:") + "\n")
		csKeys := make([]string, 0, len(e.CompletionStates))
		for k := range e.CompletionStates {
			csKeys = append(csKeys, k)
		}
		sort.Strings(csKeys)
		for _, k := range csKeys {
			b.WriteString(fmt.Sprintf("  %s: %q\n", k, e.CompletionStates[k]))
		}
	}

	// Source
	b.WriteString("\n" + detailLabelStyle.Render("Source:") + "\n")
	b.WriteString("  " + detailLabelStyle.Render("Collection: ") + e.Source.Collection + "\n")
	if e.Source.URL != "" {
		b.WriteString("  " + detailLabelStyle.Render("URL: ") + e.Source.URL + "\n")
	}
	if len(e.Source.References) > 0 {
		b.WriteString("  " + detailLabelStyle.Render("References:") + "\n")
		for _, ref := range e.Source.References {
			b.WriteString("    " + ref + "\n")
		}
	}
	b.WriteString("  " + detailLabelStyle.Render("Original key: ") + e.Source.OriginalKey + "\n")
	b.WriteString("  " + detailLabelStyle.Render("License: ") + e.Source.License + "\n")
	if e.Source.Copyright != "" {
		b.WriteString("  " + detailLabelStyle.Render("Copyright: ") + e.Source.Copyright + "\n")
	}

	// Notes
	if e.Notes != nil {
		b.WriteString("\n" + detailLabelStyle.Render("Notes: ") + *e.Notes + "\n")
	}

	// Also found in
	if len(e.AlsoFoundIn) > 0 {
		b.WriteString("\n" + detailLabelStyle.Render("Also found in:") + "\n")
		for _, afi := range e.AlsoFoundIn {
			b.WriteString(fmt.Sprintf("  %s (%s) — key: %s\n", afi.Collection, afi.License, afi.OriginalKey))
		}
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
