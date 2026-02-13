package data

import (
	"encoding/json"
	"fmt"
	"sort"
)

// DefaultIntervalMS is the fallback interval for spinners with null interval_ms.
const DefaultIntervalMS = 100

// Group holds entries that share a type and group name.
type Group struct {
	Name    string
	Type    string
	Entries []EntryEnvelope
}

// GroupedEntries organizes the catalog into spinner and progress bar groups.
type GroupedEntries struct {
	SpinnerGroups     []Group
	ProgressBarGroups []Group
	AllEntries        []EntryEnvelope
}

// LoadCatalog parses the embedded JSON bytes and returns indexed, grouped entries.
func LoadCatalog(jsonBytes []byte) (*GroupedEntries, error) {
	var catalog Catalog
	if err := json.Unmarshal(jsonBytes, &catalog); err != nil {
		return nil, fmt.Errorf("parsing catalog: %w", err)
	}

	if err := validateEntries(catalog.Entries); err != nil {
		return nil, err
	}

	spinnerMap := make(map[string][]EntryEnvelope)
	barMap := make(map[string][]EntryEnvelope)

	for _, env := range catalog.Entries {
		switch env.Entry.Type {
		case "spinner":
			spinnerMap[env.Entry.Group] = append(spinnerMap[env.Entry.Group], env)
		case "progress_bar":
			barMap[env.Entry.Group] = append(barMap[env.Entry.Group], env)
		}
	}

	result := &GroupedEntries{
		SpinnerGroups:     buildGroups(spinnerMap, "spinner"),
		ProgressBarGroups: buildGroups(barMap, "progress_bar"),
		AllEntries:        catalog.Entries,
	}
	return result, nil
}

func validateEntries(entries []EntryEnvelope) error {
	for _, env := range entries {
		e := env.Entry
		switch e.Type {
		case "spinner":
			if len(e.Frames) == 0 {
				return fmt.Errorf("spinner %q has no frames", e.ID)
			}
		case "progress_bar":
			hasChars := e.Characters != nil && e.Characters.Fill != ""
			hasPhases := len(e.Phases) > 0
			if !hasChars && !hasPhases {
				return fmt.Errorf("progress bar %q has neither characters nor phases", e.ID)
			}
		}
	}
	return nil
}

// buildGroups converts a map of group name to entries into a sorted slice of Groups.
// Groups are ordered by count descending, then name ascending as a tie-breaker.
func buildGroups(m map[string][]EntryEnvelope, entryType string) []Group {
	groups := make([]Group, 0, len(m))
	for name, entries := range m {
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Entry.Name < entries[j].Entry.Name
		})
		groups = append(groups, Group{
			Name:    name,
			Type:    entryType,
			Entries: entries,
		})
	}
	sort.Slice(groups, func(i, j int) bool {
		if len(groups[i].Entries) != len(groups[j].Entries) {
			return len(groups[i].Entries) > len(groups[j].Entries)
		}
		return groups[i].Name < groups[j].Name
	})
	return groups
}
