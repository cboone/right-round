package data

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntryEnvelope_UnmarshalJSON(t *testing.T) {
	raw := `{"id":"test/spinner","name":"test","type":"spinner","group":"braille","frames":["a","b"],"interval_ms":80,"source":{"collection":"test","url":"","original_key":"k","license":"MIT","retrieved":"2024-01-01"}}`

	var env EntryEnvelope
	err := json.Unmarshal([]byte(raw), &env)
	require.NoError(t, err)

	assert.Equal(t, "test/spinner", env.Entry.ID)
	assert.Equal(t, "test", env.Entry.Name)
	assert.Equal(t, "spinner", env.Entry.Type)
	assert.Equal(t, "braille", env.Entry.Group)
	assert.Equal(t, []string{"a", "b"}, env.Entry.Frames)
	assert.NotNil(t, env.Entry.IntervalMS)
	assert.Equal(t, 80, *env.Entry.IntervalMS)

	// Raw bytes should be preserved
	assert.NotEmpty(t, env.Raw)
	assert.JSONEq(t, raw, string(env.Raw))
}

func TestEntry_NullOptionalFields(t *testing.T) {
	raw := `{"id":"test/bar","name":"test bar","type":"progress_bar","group":"ascii","characters":{"fill":"#","empty":" "},"source":{"collection":"test","url":"","original_key":"k","license":"MIT","retrieved":"2024-01-01"}}`

	var env EntryEnvelope
	err := json.Unmarshal([]byte(raw), &env)
	require.NoError(t, err)

	assert.Nil(t, env.Entry.IntervalMS)
	assert.Nil(t, env.Entry.Notes)
	assert.Nil(t, env.Entry.Indeterminate)
	assert.Nil(t, env.Entry.Phases)
	assert.Nil(t, env.Entry.CompletionStates)
	assert.Nil(t, env.Entry.AlsoFoundIn)
	assert.Nil(t, env.Entry.Characters.Head)
	assert.Nil(t, env.Entry.Characters.Start)
	assert.Nil(t, env.Entry.Characters.End)
}

func TestEntry_WithAllOptionalFields(t *testing.T) {
	raw := `{
		"id":"test/full",
		"name":"full entry",
		"type":"spinner",
		"group":"line",
		"frames":["-","\\","|","/"],
		"interval_ms":100,
		"completion_states":{"completed":" "},
		"notes":"A test note",
		"source":{
			"collection":"test-col",
			"url":"https://example.com",
			"raw_url":"https://example.com/raw",
			"references":["https://ref1.com","https://ref2.com"],
			"original_key":"full",
			"license":"MIT",
			"license_url":"https://example.com/license",
			"copyright":"(c) Test",
			"retrieved":"2024-01-01"
		},
		"also_found_in":[
			{"collection":"other","url":"https://other.com","original_key":"x","license":"Apache-2.0"}
		]
	}`

	var env EntryEnvelope
	err := json.Unmarshal([]byte(raw), &env)
	require.NoError(t, err)

	e := env.Entry
	assert.NotNil(t, e.IntervalMS)
	assert.Equal(t, 100, *e.IntervalMS)
	assert.NotNil(t, e.Notes)
	assert.Equal(t, "A test note", *e.Notes)
	assert.Equal(t, map[string]string{"completed": " "}, e.CompletionStates)

	// Source fields
	assert.Equal(t, "https://example.com/raw", e.Source.RawURL)
	assert.Equal(t, "https://example.com/license", e.Source.LicenseURL)
	assert.Equal(t, []string{"https://ref1.com", "https://ref2.com"}, e.Source.References)
	assert.Equal(t, "(c) Test", e.Source.Copyright)

	// Also found in
	require.Len(t, e.AlsoFoundIn, 1)
	assert.Equal(t, "other", e.AlsoFoundIn[0].Collection)
	assert.Equal(t, "Apache-2.0", e.AlsoFoundIn[0].License)
}

func TestBarCharacters_WithOptionalFields(t *testing.T) {
	raw := `{"fill":"=","empty":"-","head":">","start":"[","end":"]"}`

	var chars BarCharacters
	err := json.Unmarshal([]byte(raw), &chars)
	require.NoError(t, err)

	assert.Equal(t, "=", chars.Fill)
	assert.Equal(t, "-", chars.Empty)
	assert.NotNil(t, chars.Head)
	assert.Equal(t, ">", *chars.Head)
	assert.NotNil(t, chars.Start)
	assert.Equal(t, "[", *chars.Start)
	assert.NotNil(t, chars.End)
	assert.Equal(t, "]", *chars.End)
}
