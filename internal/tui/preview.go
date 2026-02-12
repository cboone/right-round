package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cboone/right-round/internal/data"
)

const tickInterval = 16 * time.Millisecond // ~60 FPS

type animTickMsg time.Time

// animState tracks per-entry animation state.
type animState struct {
	frameIndex  int
	accumulator time.Duration
}

// animEngine manages a single global ticker and per-entry animation accumulators.
type animEngine struct {
	states map[string]*animState // keyed by entry ID
}

func newAnimEngine() *animEngine {
	return &animEngine{
		states: make(map[string]*animState),
	}
}

// tick returns a Cmd that produces the next animTickMsg.
func tick() tea.Cmd {
	return tea.Tick(tickInterval, func(t time.Time) tea.Msg {
		return animTickMsg(t)
	})
}

// advance updates all visible entries by the given elapsed duration.
// visibleIDs lists entry IDs currently on screen.
func (a *animEngine) advance(elapsed time.Duration, visibleIDs []string, entries map[string]*data.Entry) {
	for _, id := range visibleIDs {
		entry, ok := entries[id]
		if !ok || entry.Type != "spinner" || len(entry.Frames) == 0 {
			continue
		}

		state := a.getOrCreate(id)
		state.accumulator += elapsed

		interval := entryInterval(entry)
		for state.accumulator >= interval {
			state.frameIndex = (state.frameIndex + 1) % len(entry.Frames)
			state.accumulator -= interval
		}
	}
}

// currentFrame returns the current frame for a spinner entry.
func (a *animEngine) currentFrame(id string, frames []string) string {
	if len(frames) == 0 {
		return ""
	}
	state := a.getOrCreate(id)
	return frames[state.frameIndex%len(frames)]
}

func (a *animEngine) getOrCreate(id string) *animState {
	if s, ok := a.states[id]; ok {
		return s
	}
	s := &animState{}
	a.states[id] = s
	return s
}

func entryInterval(e *data.Entry) time.Duration {
	if e.IntervalMS != nil {
		return time.Duration(*e.IntervalMS) * time.Millisecond
	}
	return time.Duration(data.DefaultIntervalMS) * time.Millisecond
}
