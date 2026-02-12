package data

import (
	"encoding/json"
	"testing"

	rightround "github.com/cboone/right-round"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadCatalog_EmbeddedData(t *testing.T) {
	jsonBytes := rightround.EmbeddedCatalogJSON()
	require.NotEmpty(t, jsonBytes)

	grouped, err := LoadCatalog(jsonBytes)
	require.NoError(t, err)

	// Validate total counts
	totalSpinners := 0
	for _, g := range grouped.SpinnerGroups {
		totalSpinners += len(g.Entries)
	}
	assert.Equal(t, 333, totalSpinners)

	totalBars := 0
	for _, g := range grouped.ProgressBarGroups {
		totalBars += len(g.Entries)
	}
	assert.Equal(t, 100, totalBars)

	assert.Equal(t, 433, len(grouped.AllEntries))
}

func TestLoadCatalog_GroupOrdering(t *testing.T) {
	jsonBytes := rightround.EmbeddedCatalogJSON()
	grouped, err := LoadCatalog(jsonBytes)
	require.NoError(t, err)

	// Groups should be ordered by count descending
	for i := 1; i < len(grouped.SpinnerGroups); i++ {
		prev := len(grouped.SpinnerGroups[i-1].Entries)
		curr := len(grouped.SpinnerGroups[i].Entries)
		if prev == curr {
			// Tie-breaker: name ascending
			assert.LessOrEqual(t, grouped.SpinnerGroups[i-1].Name, grouped.SpinnerGroups[i].Name,
				"groups with equal count should be sorted by name")
		} else {
			assert.GreaterOrEqual(t, prev, curr,
				"groups should be sorted by count descending")
		}
	}
}

func TestLoadCatalog_EntrySortingWithinGroup(t *testing.T) {
	jsonBytes := rightround.EmbeddedCatalogJSON()
	grouped, err := LoadCatalog(jsonBytes)
	require.NoError(t, err)

	for _, g := range grouped.SpinnerGroups {
		for i := 1; i < len(g.Entries); i++ {
			assert.LessOrEqual(t, g.Entries[i-1].Entry.Name, g.Entries[i].Entry.Name,
				"entries in group %q should be sorted by name", g.Name)
		}
	}
}

func TestLoadCatalog_RejectsEmptyFrames(t *testing.T) {
	catalog := makeCatalog([]Entry{
		{ID: "bad/spinner", Name: "bad", Type: "spinner", Group: "test", Frames: []string{}},
	})

	_, err := LoadCatalog(catalog)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no frames")
}

func TestLoadCatalog_AcceptsPhaseOnlyBar(t *testing.T) {
	catalog := makeCatalog([]Entry{
		{ID: "test/phase-bar", Name: "phase bar", Type: "progress_bar", Group: "phased", Phases: []string{"a", "b", "c"}},
	})

	grouped, err := LoadCatalog(catalog)
	require.NoError(t, err)
	assert.Len(t, grouped.ProgressBarGroups, 1)
	assert.Len(t, grouped.ProgressBarGroups[0].Entries, 1)
}

func TestLoadCatalog_RejectsBarWithNeitherCharsNorPhases(t *testing.T) {
	catalog := makeCatalog([]Entry{
		{ID: "bad/bar", Name: "bad bar", Type: "progress_bar", Group: "test"},
	})

	_, err := LoadCatalog(catalog)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "neither characters nor phases")
}

func TestLoadCatalog_MalformedJSON(t *testing.T) {
	_, err := LoadCatalog([]byte(`{invalid json`))
	assert.Error(t, err)
}

func TestLoadCatalog_Deterministic(t *testing.T) {
	jsonBytes := rightround.EmbeddedCatalogJSON()

	grouped1, err := LoadCatalog(jsonBytes)
	require.NoError(t, err)

	grouped2, err := LoadCatalog(jsonBytes)
	require.NoError(t, err)

	// Same group ordering
	require.Equal(t, len(grouped1.SpinnerGroups), len(grouped2.SpinnerGroups))
	for i := range grouped1.SpinnerGroups {
		assert.Equal(t, grouped1.SpinnerGroups[i].Name, grouped2.SpinnerGroups[i].Name)
		require.Equal(t, len(grouped1.SpinnerGroups[i].Entries), len(grouped2.SpinnerGroups[i].Entries))
		for j := range grouped1.SpinnerGroups[i].Entries {
			assert.Equal(t, grouped1.SpinnerGroups[i].Entries[j].Entry.ID, grouped2.SpinnerGroups[i].Entries[j].Entry.ID)
		}
	}
}

// makeCatalog creates a minimal JSON catalog from entries for testing.
func makeCatalog(entries []Entry) []byte {
	type testCatalog struct {
		Version string  `json:"version"`
		Entries []Entry `json:"entries"`
	}
	c := testCatalog{
		Version: "1.0.0",
		Entries: entries,
	}
	b, _ := json.Marshal(c)
	return b
}
