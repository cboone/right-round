package data

import "encoding/json"

// Catalog is the top-level structure of the progress-indicators.json file.
type Catalog struct {
	Version     string          `json:"version"`
	Generated   string          `json:"generated"`
	Description string          `json:"description"`
	Stats       json.RawMessage `json:"stats"`
	EntrySchema json.RawMessage `json:"entry_schema"`
	Entries     []EntryEnvelope `json:"entries"`
}

// EntryEnvelope wraps a parsed Entry alongside its original JSON bytes
// to guarantee lossless clipboard export.
type EntryEnvelope struct {
	Entry Entry
	Raw   json.RawMessage
}

// UnmarshalJSON decodes an entry, keeping a copy of the raw bytes.
func (e *EntryEnvelope) UnmarshalJSON(data []byte) error {
	e.Raw = make(json.RawMessage, len(data))
	copy(e.Raw, data)
	return json.Unmarshal(data, &e.Entry)
}

// Entry represents a single progress indicator entry.
type Entry struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	Type             string            `json:"type"`
	Group            string            `json:"group"`
	Frames           []string          `json:"frames,omitempty"`
	IntervalMS       *int              `json:"interval_ms,omitempty"`
	Characters       *BarCharacters    `json:"characters,omitempty"`
	Phases           []string          `json:"phases,omitempty"`
	Indeterminate    *string           `json:"indeterminate,omitempty"`
	CompletionStates map[string]string `json:"completion_states,omitempty"`
	Notes            *string           `json:"notes,omitempty"`
	Source           Source            `json:"source"`
	AlsoFoundIn      []AlsoFoundIn     `json:"also_found_in,omitempty"`
}

// BarCharacters holds the character set for progress bar rendering.
type BarCharacters struct {
	Fill  string  `json:"fill"`
	Empty string  `json:"empty"`
	Head  *string `json:"head,omitempty"`
	Start *string `json:"start,omitempty"`
	End   *string `json:"end,omitempty"`
}

// Source holds attribution information for an entry.
type Source struct {
	Collection string   `json:"collection"`
	URL        string   `json:"url"`
	RawURL     string   `json:"raw_url,omitempty"`
	References []string `json:"references,omitempty"`
	OriginalKey string  `json:"original_key"`
	License    string   `json:"license"`
	LicenseURL string   `json:"license_url,omitempty"`
	Copyright  string   `json:"copyright,omitempty"`
	Retrieved  string   `json:"retrieved"`
}

// AlsoFoundIn records another collection where the same indicator was found.
type AlsoFoundIn struct {
	Collection  string `json:"collection"`
	URL         string `json:"url"`
	OriginalKey string `json:"original_key"`
	License     string `json:"license"`
}
