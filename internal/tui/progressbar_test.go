package tui

import (
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/cboone/right-round/internal/data"
	"github.com/stretchr/testify/assert"
)

func strPtr(s string) *string { return &s }

func TestRenderProgressBar_RequiredFieldsOnly(t *testing.T) {
	chars := &data.BarCharacters{Fill: "#", Empty: "-"}
	bar := renderProgressBar(chars, nil, 0.5, 10)
	assert.Equal(t, "#####-----", bar)
}

func TestRenderProgressBar_WithStartEnd(t *testing.T) {
	chars := &data.BarCharacters{Fill: "=", Empty: " ", Start: strPtr("["), End: strPtr("]")}
	bar := renderProgressBar(chars, nil, 0.5, 10)
	// inner width = 10 - 1 - 1 = 8, fill=4, empty=4
	assert.Equal(t, "[====    ]", bar)
}

func TestRenderProgressBar_WithHead(t *testing.T) {
	chars := &data.BarCharacters{Fill: "=", Empty: "-", Head: strPtr(">")}

	// At 50%, head should appear
	bar := renderProgressBar(chars, nil, 0.5, 10)
	assert.Contains(t, bar, ">")
	assert.Contains(t, bar, "=")
	assert.Contains(t, bar, "-")
}

func TestRenderProgressBar_HeadOmittedAt0Percent(t *testing.T) {
	chars := &data.BarCharacters{Fill: "=", Empty: "-", Head: strPtr(">")}
	bar := renderProgressBar(chars, nil, 0.0, 10)
	assert.NotContains(t, bar, ">")
}

func TestRenderProgressBar_HeadOmittedAt100Percent(t *testing.T) {
	chars := &data.BarCharacters{Fill: "=", Empty: "-", Head: strPtr(">")}
	bar := renderProgressBar(chars, nil, 1.0, 10)
	assert.NotContains(t, bar, ">")
}

func TestRenderProgressBar_ZeroPercent(t *testing.T) {
	chars := &data.BarCharacters{Fill: "#", Empty: "-"}
	bar := renderProgressBar(chars, nil, 0.0, 10)
	assert.Equal(t, "----------", bar)
}

func TestRenderProgressBar_FullPercent(t *testing.T) {
	chars := &data.BarCharacters{Fill: "#", Empty: "-"}
	bar := renderProgressBar(chars, nil, 1.0, 10)
	assert.Equal(t, "##########", bar)
}

func TestRenderPhaseOnlyBar(t *testing.T) {
	phases := []string{" ", "▏", "▎", "▍", "▌", "▋", "▊", "▉", "█"}

	// At 0%, should be all first phase (space)
	bar0 := renderPhaseOnlyBar(phases, 0.0, 10)
	assert.Equal(t, 10, lipgloss.Width(bar0))

	// At 100%, should be all last phase
	bar100 := renderPhaseOnlyBar(phases, 1.0, 10)
	assert.Equal(t, "██████████", bar100)
}

func TestRenderPhaseOnlyBar_NilCharsWithPhases(t *testing.T) {
	phases := []string{"⡀", "⡄", "⡆", "⡇", "⣇", "⣧", "⣷", "⣿"}
	bar := renderProgressBar(nil, phases, 0.5, 10)
	assert.NotEmpty(t, bar)
	assert.Equal(t, 10, lipgloss.Width(bar))
}

func TestRenderProgressBar_NilCharsNilPhases(t *testing.T) {
	bar := renderProgressBar(nil, nil, 0.5, 10)
	assert.Equal(t, "", bar)
}

func TestRenderIndeterminate(t *testing.T) {
	pattern := "<=>"
	result := renderIndeterminate(pattern, 20, 5)
	assert.Equal(t, 20, lipgloss.Width(result))
}

func TestRenderIndeterminate_EmptyPattern(t *testing.T) {
	result := renderIndeterminate("", 20, 0)
	assert.Equal(t, "", result)
}

func TestRenderStaticBar(t *testing.T) {
	chars := &data.BarCharacters{Fill: "#", Empty: "-"}
	bar := renderStaticBar(chars, nil, 10)
	// 40% of 10 = 4 fill chars
	assert.Equal(t, "####------", bar)
}

func TestRenderProgressBar_MultiCellHead(t *testing.T) {
	chars := &data.BarCharacters{Fill: "=", Empty: "-", Head: strPtr("=>")}
	bar := renderProgressBar(chars, nil, 0.5, 10)
	// Head "=>" is 2 cells; total should still be 10 cells
	assert.Equal(t, 10, lipgloss.Width(bar))
	assert.Contains(t, bar, "=>")
}

func TestRenderProgressBar_UnicodeWidth(t *testing.T) {
	tests := []struct {
		name   string
		chars  *data.BarCharacters
		width  int
		pct    float64
		expect int // expected display width
	}{
		{
			name:   "ASCII chars",
			chars:  &data.BarCharacters{Fill: "=", Empty: " "},
			width:  20,
			pct:    0.5,
			expect: 20,
		},
		{
			name:   "block chars",
			chars:  &data.BarCharacters{Fill: "█", Empty: "░"},
			width:  20,
			pct:    0.5,
			expect: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bar := renderProgressBar(tt.chars, nil, tt.pct, tt.width)
			assert.Equal(t, tt.expect, lipgloss.Width(bar))
		})
	}
}
