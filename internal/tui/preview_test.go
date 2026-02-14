package tui

import (
	"testing"
	"time"

	"github.com/cboone/right-round/internal/data"
	"github.com/stretchr/testify/assert"
)

func intPtr(v int) *int { return &v }

func TestAnimEngine_BasicFrameAdvance(t *testing.T) {
	engine := newAnimEngine()
	entries := map[string]*data.Entry{
		"test": {
			ID:         "test",
			Type:       "spinner",
			Frames:     []string{"a", "b", "c"},
			IntervalMS: intPtr(100),
		},
	}

	// Initially at frame 0
	assert.Equal(t, "a", engine.currentFrame("test", entries["test"].Frames))

	// Advance by 100ms -> frame 1
	engine.advance(100*time.Millisecond, []string{"test"}, entries)
	assert.Equal(t, "b", engine.currentFrame("test", entries["test"].Frames))

	// Advance by 100ms -> frame 2
	engine.advance(100*time.Millisecond, []string{"test"}, entries)
	assert.Equal(t, "c", engine.currentFrame("test", entries["test"].Frames))

	// Advance by 100ms -> wraps to frame 0
	engine.advance(100*time.Millisecond, []string{"test"}, entries)
	assert.Equal(t, "a", engine.currentFrame("test", entries["test"].Frames))
}

func TestAnimEngine_AccumulatorCarryOver(t *testing.T) {
	engine := newAnimEngine()
	entries := map[string]*data.Entry{
		"test": {
			ID:         "test",
			Type:       "spinner",
			Frames:     []string{"a", "b", "c"},
			IntervalMS: intPtr(100),
		},
	}

	// Advance by 150ms -> one frame advance, 50ms carried over
	engine.advance(150*time.Millisecond, []string{"test"}, entries)
	assert.Equal(t, "b", engine.currentFrame("test", entries["test"].Frames))

	// Advance by 60ms -> accumulator now 110ms, another frame advance with 10ms carry
	engine.advance(60*time.Millisecond, []string{"test"}, entries)
	assert.Equal(t, "c", engine.currentFrame("test", entries["test"].Frames))
}

func TestAnimEngine_CatchUpMultipleFrames(t *testing.T) {
	engine := newAnimEngine()
	entries := map[string]*data.Entry{
		"test": {
			ID:         "test",
			Type:       "spinner",
			Frames:     []string{"a", "b", "c", "d"},
			IntervalMS: intPtr(50),
		},
	}

	// Advance by 150ms -> should catch up 3 frames
	engine.advance(150*time.Millisecond, []string{"test"}, entries)
	assert.Equal(t, "d", engine.currentFrame("test", entries["test"].Frames))
}

func TestAnimEngine_NullIntervalUsesDefault(t *testing.T) {
	engine := newAnimEngine()
	entries := map[string]*data.Entry{
		"test": {
			ID:     "test",
			Type:   "spinner",
			Frames: []string{"a", "b"},
			// IntervalMS is nil
		},
	}

	// Default is 100ms. Advance by 100ms -> frame 1
	engine.advance(100*time.Millisecond, []string{"test"}, entries)
	assert.Equal(t, "b", engine.currentFrame("test", entries["test"].Frames))
}

func TestAnimEngine_OnlyAdvancesVisibleEntries(t *testing.T) {
	engine := newAnimEngine()
	entries := map[string]*data.Entry{
		"visible": {
			ID:         "visible",
			Type:       "spinner",
			Frames:     []string{"a", "b"},
			IntervalMS: intPtr(100),
		},
		"hidden": {
			ID:         "hidden",
			Type:       "spinner",
			Frames:     []string{"x", "y"},
			IntervalMS: intPtr(100),
		},
	}

	// Only "visible" is in the visible list
	engine.advance(100*time.Millisecond, []string{"visible"}, entries)
	assert.Equal(t, "b", engine.currentFrame("visible", entries["visible"].Frames))
	assert.Equal(t, "x", engine.currentFrame("hidden", entries["hidden"].Frames))
}

func TestAnimEngine_EmptyFrames(t *testing.T) {
	engine := newAnimEngine()
	assert.Equal(t, "", engine.currentFrame("empty", nil))
	assert.Equal(t, "", engine.currentFrame("empty", []string{}))
}

func TestAnimEngine_NonSpinnerSkipped(t *testing.T) {
	engine := newAnimEngine()
	entries := map[string]*data.Entry{
		"bar": {
			ID:   "bar",
			Type: "progress_bar",
		},
	}

	// Should advance determinate bar animation state
	engine.advance(100*time.Millisecond, []string{"bar"}, entries)
	assert.Greater(t, engine.currentOffset("bar"), 0)
}

func TestAnimEngine_IndeterminateProgressBarAdvance(t *testing.T) {
	engine := newAnimEngine()
	pattern := "<=>"
	entries := map[string]*data.Entry{
		"bar": {
			ID:            "bar",
			Type:          "progress_bar",
			Indeterminate: &pattern,
		},
	}

	assert.Equal(t, 0, engine.currentOffset("bar"))
	engine.advance(80*time.Millisecond, []string{"bar"}, entries)
	assert.Equal(t, 1, engine.currentOffset("bar"))

	engine.advance(160*time.Millisecond, []string{"bar"}, entries)
	assert.Equal(t, 3, engine.currentOffset("bar"))
}

func TestAnimEngine_DeterminateProgressBarAdvancesWithoutPattern(t *testing.T) {
	engine := newAnimEngine()
	entries := map[string]*data.Entry{
		"bar": {
			ID:   "bar",
			Type: "progress_bar",
		},
	}

	engine.advance(500*time.Millisecond, []string{"bar"}, entries)
	assert.Greater(t, engine.currentOffset("bar"), 0)
}

func TestAnimEngine_CurrentProgressPctPingPong(t *testing.T) {
	engine := newAnimEngine()
	engine.states["bar"] = &animState{frameIndex: 0}
	assert.Equal(t, 0.0, engine.currentProgressPct("bar"))

	engine.states["bar"].frameIndex = 100
	assert.Equal(t, 1.0, engine.currentProgressPct("bar"))

	engine.states["bar"].frameIndex = 150
	assert.Equal(t, 0.5, engine.currentProgressPct("bar"))
}

func TestEntryInterval(t *testing.T) {
	tests := []struct {
		name     string
		entry    *data.Entry
		expected time.Duration
	}{
		{
			name:     "explicit interval",
			entry:    &data.Entry{IntervalMS: intPtr(80)},
			expected: 80 * time.Millisecond,
		},
		{
			name:     "nil interval uses default",
			entry:    &data.Entry{},
			expected: 100 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, entryInterval(tt.entry))
		})
	}
}
