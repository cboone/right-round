package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/cboone/right-round/internal/data"
)

// repeatToCells repeats s enough times to fill the given number of terminal cells,
// padding any remainder with spaces. This correctly handles multi-cell glyphs (e.g. emoji).
func repeatToCells(s string, cells int) string {
	charWidth := lipgloss.Width(s)
	if charWidth < 1 {
		charWidth = 1
	}
	reps := cells / charWidth
	actual := reps * charWidth
	result := strings.Repeat(s, reps)
	if actual < cells {
		result += strings.Repeat(" ", cells-actual)
	}
	return result
}

// renderProgressBar renders a progress bar at the given fill percentage (0.0 to 1.0)
// within the specified width (in terminal cells).
func renderProgressBar(chars *data.BarCharacters, phases []string, pct float64, width int) string {
	if chars == nil && len(phases) > 0 {
		return renderPhaseOnlyBar(phases, pct, width)
	}
	if chars == nil {
		return ""
	}

	startStr := ""
	endStr := ""
	innerWidth := width

	if chars.Start != nil {
		startStr = *chars.Start
		innerWidth -= lipgloss.Width(startStr)
	}
	if chars.End != nil {
		endStr = *chars.End
		innerWidth -= lipgloss.Width(endStr)
	}
	if innerWidth < 1 {
		innerWidth = 1
	}

	fillCells := int(pct * float64(innerWidth))
	if fillCells > innerWidth {
		fillCells = innerWidth
	}
	emptyCells := innerWidth - fillCells

	var b strings.Builder
	b.WriteString(startStr)

	// Head only renders at the advancing boundary (not at 0% or 100%)
	hasHead := chars.Head != nil && fillCells > 0 && fillCells < innerWidth
	headWidth := 0
	if hasHead {
		headWidth = lipgloss.Width(*chars.Head)
		if headWidth < 1 {
			headWidth = 1
		}
		fillCells -= headWidth
		if fillCells < 0 {
			fillCells = 0
		}
	}

	// Fill portion
	if len(phases) > 0 && fillCells < innerWidth && !hasHead {
		// Use phases for sub-character resolution at the boundary
		b.WriteString(repeatToCells(chars.Fill, fillCells))
		totalFill := pct * float64(innerWidth)
		frac := totalFill - float64(fillCells)
		if frac < 0 {
			frac = 0
		} else if frac > 1 {
			frac = 1
		}
		phaseIdx := int(frac * float64(len(phases)))
		if phaseIdx > 0 && phaseIdx < len(phases) && emptyCells > 0 {
			b.WriteString(phases[phaseIdx])
			emptyCells--
		}
	} else {
		b.WriteString(repeatToCells(chars.Fill, fillCells))
	}

	if hasHead {
		b.WriteString(*chars.Head)
	}

	b.WriteString(repeatToCells(chars.Empty, emptyCells))
	b.WriteString(endStr)

	return b.String()
}

// renderPhaseOnlyBar renders a bar using only the phases array.
func renderPhaseOnlyBar(phases []string, pct float64, width int) string {
	if len(phases) == 0 {
		return ""
	}
	lastPhase := phases[len(phases)-1]
	firstPhase := phases[0]

	fullCells := int(pct * float64(width))
	if fullCells > width {
		fullCells = width
	}
	remaining := width - fullCells

	var b strings.Builder
	b.WriteString(strings.Repeat(lastPhase, fullCells))

	// Partial cell at boundary
	if remaining > 0 && fullCells < width {
		fractional := pct*float64(width) - float64(fullCells)
		phaseIdx := int(fractional * float64(len(phases)-1))
		if phaseIdx >= len(phases) {
			phaseIdx = len(phases) - 1
		}
		if phaseIdx > 0 {
			b.WriteString(phases[phaseIdx])
			remaining--
		}
	}

	b.WriteString(strings.Repeat(firstPhase, remaining))
	return b.String()
}

// renderIndeterminate renders an indeterminate pattern within the given width.
func renderIndeterminate(pattern string, width int, offset int) string {
	if pattern == "" || width < 1 {
		return ""
	}

	// Convert pattern to runes and compute per-rune cell widths for consistent indexing
	patRunes := []rune(pattern)
	patCellWidths := make([]int, len(patRunes))
	totalCells := 0
	for i, r := range patRunes {
		w := lipgloss.Width(string(r))
		patCellWidths[i] = w
		totalCells += w
	}

	var b strings.Builder
	pos := offset % (width + totalCells)
	cellCol := 0
	for cellCol < width {
		relCell := cellCol - pos + totalCells
		// Find which rune (if any) covers this cell position
		wrote := false
		if relCell >= 0 && relCell < totalCells {
			cumulative := 0
			for ri, rw := range patCellWidths {
				if relCell >= cumulative && relCell < cumulative+rw {
					// Only emit the rune at its first cell to avoid duplicates
					if relCell == cumulative {
						b.WriteRune(patRunes[ri])
						cellCol += rw
					} else {
						// Mid-glyph cell: skip with space
						b.WriteRune(' ')
						cellCol++
					}
					wrote = true
					break
				}
				cumulative += rw
			}
		}
		if !wrote {
			b.WriteRune(' ')
			cellCol++
		}
	}

	return b.String()
}

// renderStaticBar renders a static progress bar sample at ~40% for list display.
func renderStaticBar(chars *data.BarCharacters, phases []string, width int) string {
	return renderProgressBar(chars, phases, 0.4, width)
}
