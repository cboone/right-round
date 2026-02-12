package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/cboone/right-round/internal/data"
)

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
	if hasHead {
		fillCells--
	}

	// Fill portion
	if len(phases) > 0 && fillCells < innerWidth && !hasHead {
		// Use phases for sub-character resolution at the boundary
		b.WriteString(strings.Repeat(chars.Fill, fillCells))
		phaseIdx := int(pct*float64(innerWidth)*float64(len(phases))) % len(phases)
		if phaseIdx > 0 && emptyCells > 0 {
			b.WriteString(phases[phaseIdx])
			emptyCells--
		}
	} else {
		b.WriteString(strings.Repeat(chars.Fill, fillCells))
	}

	if hasHead {
		b.WriteString(*chars.Head)
	}

	b.WriteString(strings.Repeat(chars.Empty, emptyCells))
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

	patLen := lipgloss.Width(pattern)
	totalLen := width + patLen
	var b strings.Builder

	pos := offset % (width + patLen)
	for i := 0; i < width; i++ {
		relPos := i - pos + patLen
		if relPos >= 0 && relPos < patLen {
			// Inside the pattern - extract the character at this position
			runeIdx := 0
			for _, r := range pattern {
				if runeIdx == relPos {
					b.WriteRune(r)
					break
				}
				runeIdx++
			}
		} else {
			b.WriteRune(' ')
		}
	}

	_ = totalLen
	return b.String()
}

// renderStaticBar renders a static progress bar sample at ~40% for list display.
func renderStaticBar(chars *data.BarCharacters, phases []string, width int) string {
	return renderProgressBar(chars, phases, 0.4, width)
}
