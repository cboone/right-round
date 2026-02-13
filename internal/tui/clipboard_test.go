package tui

import (
	"encoding/json"
	"testing"

	"github.com/cboone/right-round/internal/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClipboardPayload_RoundTrip(t *testing.T) {
	// Create an entry with extra fields to verify raw JSON preservation
	rawJSON := `{"id":"test/entry","name":"test","type":"spinner","group":"braille","frames":["a","b"],"extra_field":"preserved","source":{"collection":"test","url":"","original_key":"k","license":"MIT","retrieved":"2024-01-01"}}`

	var env data.EntryEnvelope
	err := json.Unmarshal([]byte(rawJSON), &env)
	require.NoError(t, err)

	// The raw message should preserve the extra_field that isn't in our struct
	var roundTripped map[string]interface{}
	err = json.Unmarshal(env.Raw, &roundTripped)
	require.NoError(t, err)
	assert.Contains(t, roundTripped, "extra_field")
	assert.Equal(t, "preserved", roundTripped["extra_field"])
}

func TestClipboardPayload_PrettyPrint(t *testing.T) {
	rawJSON := `{"id":"test/entry","name":"test","type":"spinner"}`
	env := &data.EntryEnvelope{
		Entry: data.Entry{ID: "test/entry", Name: "test", Type: "spinner"},
		Raw:   json.RawMessage(rawJSON),
	}

	// Verify the raw JSON can be pretty-printed
	var pretty json.RawMessage
	err := json.Unmarshal(env.Raw, &pretty)
	require.NoError(t, err)

	formatted, err := json.MarshalIndent(pretty, "", "  ")
	require.NoError(t, err)
	assert.Contains(t, string(formatted), "\n")
	assert.Contains(t, string(formatted), "  ")
}

func TestCopyToClipboard_NilEntry(t *testing.T) {
	cmd := copyToClipboard(nil)
	msg := cmd()
	result, ok := msg.(clipResultMsg)
	require.True(t, ok)
	assert.NoError(t, result.err)
}
