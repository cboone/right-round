package tui

import (
	"encoding/json"
	"testing"

	"github.com/cboone/right-round/internal/data"
	"github.com/stretchr/testify/assert"
)

func makeTestSpinnerEnvelope() *data.EntryEnvelope {
	e := data.Entry{
		ID:     "test/spinner",
		Name:   "test spinner",
		Type:   "spinner",
		Group:  "braille",
		Frames: []string{"a", "b", "c"},
		Source: data.Source{
			Collection:  "test-collection",
			URL:         "https://example.com",
			OriginalKey: "test",
			License:     "MIT",
			Copyright:   "(c) Test",
			Retrieved:   "2024-01-01",
		},
	}
	raw, _ := json.Marshal(e)
	return &data.EntryEnvelope{Entry: e, Raw: raw}
}

func makeTestBarEnvelope() *data.EntryEnvelope {
	notes := "Test notes"
	e := data.Entry{
		ID:         "test/bar",
		Name:       "test bar",
		Type:       "progress_bar",
		Group:      "ascii",
		Characters: &data.BarCharacters{Fill: "=", Empty: "-", Head: strPtr(">")},
		Notes:      &notes,
		Source: data.Source{
			Collection:  "test-collection",
			URL:         "https://example.com",
			OriginalKey: "test",
			License:     "MIT",
			Retrieved:   "2024-01-01",
		},
		AlsoFoundIn: []data.AlsoFoundIn{
			{Collection: "other", URL: "https://other.com", OriginalKey: "k", License: "Apache-2.0"},
		},
	}
	raw, _ := json.Marshal(e)
	return &data.EntryEnvelope{Entry: e, Raw: raw}
}

func TestDetailModel_SpinnerContent(t *testing.T) {
	anim := newAnimEngine()
	m := newDetailModel(anim)
	m.setSize(80, 30)
	m.setEntry(makeTestSpinnerEnvelope())
	m.updateContent()

	content := m.viewport.View()
	assert.Contains(t, content, "test spinner")
	assert.Contains(t, content, "test/spinner")
	assert.Contains(t, content, "spinner")
	assert.Contains(t, content, "braille")
	assert.Contains(t, content, "Frames:")
	assert.Contains(t, content, "MIT")
	assert.Contains(t, content, "test-collection")
}

func TestDetailModel_ProgressBarContent(t *testing.T) {
	anim := newAnimEngine()
	m := newDetailModel(anim)
	m.setSize(80, 30)
	m.setEntry(makeTestBarEnvelope())
	m.updateContent()

	content := m.viewport.View()
	assert.Contains(t, content, "test bar")
	assert.Contains(t, content, "progress_bar")
	assert.Contains(t, content, "Characters:")
	assert.Contains(t, content, "Test notes")
	assert.Contains(t, content, "Also found in:")
	assert.Contains(t, content, "other")
}

func TestDetailModel_NoEntry(t *testing.T) {
	anim := newAnimEngine()
	m := newDetailModel(anim)
	m.setSize(80, 30)
	m.setEntry(nil)
	m.updateContent()

	content := m.viewport.View()
	assert.Contains(t, content, "No entry selected")
}
