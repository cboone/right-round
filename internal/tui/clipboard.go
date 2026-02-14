package tui

import (
	"encoding/json"

	tea "charm.land/bubbletea/v2"
	"github.com/atotto/clipboard"
	"github.com/cboone/right-round/internal/data"
)

type clipResultMsg struct {
	err error
}

// copyToClipboard returns a tea.Cmd that copies the entry's raw JSON to the system clipboard.
func copyToClipboard(entry *data.EntryEnvelope) tea.Cmd {
	return func() tea.Msg {
		if entry == nil {
			return clipResultMsg{err: nil}
		}

		// Pretty-print the raw JSON for readability
		var pretty json.RawMessage
		if err := json.Unmarshal(entry.Raw, &pretty); err != nil {
			return clipResultMsg{err: err}
		}
		formatted, err := json.MarshalIndent(pretty, "", "  ")
		if err != nil {
			return clipResultMsg{err: err}
		}

		err = clipboard.WriteAll(string(formatted))
		return clipResultMsg{err: err}
	}
}
