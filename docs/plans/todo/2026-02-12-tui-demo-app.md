# Plan: Create TUI Demo App for Progress Indicators

## Context

The `right-round` repository is a data catalog of 433 terminal progress indicators (333 spinners, 100 progress bars) consolidated from 26 open-source collections into a single `progress-indicators.json` file. There is currently no way to browse or preview these indicators. This plan creates a Go TUI application using Cobra and Bubble Tea that lets users navigate, preview, and copy indicators interactively.

## Directory Structure

```
right-round/                              (repo root)
├── progress-indicators.json              (existing, 433 entries)
├── embedded.go                           (NEW - go:embed for the JSON)
├── go.mod                                (NEW - github.com/cboone/right-round)
├── go.sum                                (NEW)
├── Makefile                              (NEW)
├── README.md                             (existing - rewrite)
├── .gitignore                            (NEW - Go patterns)
├── .goreleaser.yml                       (NEW)
├── .github/
│   └── workflows/
│       ├── ci.yml                        (NEW - build, lint, test)
│       └── release.yml                   (NEW - goreleaser on tag)
├── cmd/
│   └── right-round/
│       └── main.go                       (NEW - Cobra entrypoint)
└── internal/
    ├── data/
    │   ├── types.go                      (NEW - Go structs for JSON schema)
    │   └── loader.go                     (NEW - parse, group, index)
    └── tui/
        ├── app.go                        (NEW - top-level Bubble Tea model)
        ├── keys.go                       (NEW - key bindings)
        ├── styles.go                     (NEW - Lip Gloss styles)
        ├── list.go                       (NEW - grouped list panel)
        ├── detail.go                     (NEW - expandable detail panel)
        ├── preview.go                    (NEW - single-ticker animation engine)
        ├── progressbar.go               (NEW - progress bar rendering)
        └── clipboard.go                 (NEW - copy entry as JSON)
```

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/charmbracelet/bubbletea` v1 | TUI framework |
| `github.com/charmbracelet/bubbles` | Viewport, help, key components |
| `github.com/charmbracelet/lipgloss` v1 | Terminal styling |
| `github.com/spf13/cobra` | CLI command structure |
| `github.com/atotto/clipboard` | System clipboard access |
| `github.com/stretchr/testify` | Test assertions and suites |

Using Bubble Tea v1 (stable) rather than v2 (still RC).

## Architecture

### Data Layer (`internal/data/`)

**`types.go`** — Go structs mirroring the JSON schema:
- `Catalog` — top-level: version, generated, description, stats, entry_schema, entries
- `Entry` — id, name, type, group, frames, interval_ms, characters, phases, indeterminate, completion_states, notes, source, also_found_in
- `BarCharacters` — fill, empty, head, start, end (optional fields are pointers)
- `Source` — collection, url, raw_url, references, original_key, license, license_url, copyright, retrieved
- `AlsoFoundIn` — collection, url, original_key, license
- `IntervalMS` as `*int` (nullable), `Notes` as `*string` (nullable), optional fields as pointers where the JSON omits keys
- `EntryEnvelope` wrapper for each entry that stores both parsed `Entry` and original `json.RawMessage` to guarantee lossless clipboard export

**`loader.go`** — Parses embedded JSON, builds indexed structures:
- `GroupedEntries` with `SpinnerGroups []Group` and `ProgressBarGroups []Group`
- Each `Group` has name, type, and sorted entries
- Groups ordered by count descending (matching stats ordering), then group name ascending as deterministic tie-breaker
- Builds a row index from rendered list rows to backing entries for stable cursor behavior with filtering and collapsed groups
- Validates and normalizes data during load:
  - missing spinner `interval_ms` uses default 100ms at render time
  - spinner entries must have at least one frame
  - progress bar entries must include `characters.fill` and may omit start/end/head

**`embedded.go`** (repo root) — Must live at repo root because `go:embed` cannot traverse `..` paths. Contains only:
```go
package rightround

import _ "embed"

//go:embed progress-indicators.json
var ProgressIndicatorsJSON []byte
```

### TUI Layer (`internal/tui/`)

**Layout** — Two modes based on terminal width:
- Wide (>=100 cols): list panel (40%) + detail panel (60%) side by side
- Narrow (<100 cols): detail panel replaces list panel entirely

**Top bar**: Tabs for "Spinners" / "Progress Bars" switching (Tab key)
**Bottom bar**: Context-sensitive help text

**`app.go`** — Top-level Bubble Tea model:
- Orchestrates layout, panels, animation ticker
- Handles `tea.WindowSizeMsg` for responsive layout
- Dispatches a single `animTickMsg` every 16ms (~60 FPS) for all animations
- Routes key messages to the active panel
- When `--type` is set, applies a hard type filter for the session and disables tab switching with a help hint

**`list.go`** — Custom grouped list (not bubbles `list.Model`):
- Grouped display with section headers needs custom rendering
- Flat cursor + scroll offset tracking
- Each row: `> name                ⠋` (selected) or `  name                ⠹` (unselected)
- Group headers: `BRAILLE (54)` styled bold with accent color
- Progress bar rows show a static sample bar at ~40% fill
- Spinner rows show live animated preview

Reason for custom list: bubbles `list.Model` is designed for flat filterable lists. Our grouped headers, live animation previews, and mixed entry types are better served by a custom implementation.

**`detail.go`** — Expandable detail panel using `viewport.Model` for scrolling:
- Entry name, ID, type, group
- Frames count and interval (spinners) or character set display (progress bars)
- Live animated preview (larger)
- Source: collection, license, copyright, URL, original key
- Notes (if present)
- Also found in (if present): list of other collections
- Progress bar entries: full-width rendered bar with phases demo

**`preview.go`** — Single-ticker animation engine:
- One global `tea.Tick` at 16ms instead of per-spinner tickers
- Per-entry accumulator tracks elapsed time since last frame advance
- Only animates entries currently visible on screen (typically 15-20)
- Default interval of 100ms for spinners with null `interval_ms` (207 of 333)
- Frame advancement: while `accumulated >= entryInterval`, advance frame index and subtract interval (carry remainder forward to avoid drift)
- Supports fractional catch-up when terminal redraw lags by more than one interval

**`progressbar.go`** — Renders progress bars from character sets:
- Uses optional pattern `start + fill*n + head + empty*m + end` where start/head/end are rendered only when present
- Handles entries with phases (sub-character resolution at fill boundary)
- Handles entries with indeterminate patterns for dedicated preview mode
- Preview width adapts to available space

**`clipboard.go`** — Copy selected entry as full indented JSON:
- Runs as `tea.Cmd` (async) to avoid blocking UI
- Shows brief status message on success/failure
- Uses entry `json.RawMessage` for round-trip-safe output so copied JSON never drops optional/unknown fields

**`styles.go`** — Lip Gloss styles with `AdaptiveColor` for light/dark terminals:
- Accent color, subtle color, group header style, selected/normal item styles
- Detail panel border, spinner preview column, help bar, tab styles

**`keys.go`** — Key bindings (vi-style + arrows):
- `j/k` or `up/down`: navigate
- `enter/l`: expand detail view
- `esc/h`: collapse back to list
- `tab`: switch between Spinners / Progress Bars
- `c`: copy selected entry to clipboard
- `/`: search/filter (text input to filter entries by name)
- `q/ctrl+c`: quit
- `?`: toggle full help
- `pgup/pgdn`, `home/end`: fast navigation

### CLI Layer (`cmd/right-round/main.go`)

Cobra root command with optional flags:
- `--filter <group>`: start with a specific group
- `--type <spinner|progress_bar>`: lock UI to one type for the session (tab disabled)

### Unicode Width Handling

Use `lipgloss.Width()` (backed by `go-runewidth`) for all display width calculations. Some spinner frames use emoji or wide Unicode characters that occupy 2 terminal cells. Preview column truncates frames exceeding 8 cells wide; full frames shown in detail view.

## Build & CI/CD

**Makefile targets**: `build`, `install`, `clean`, `lint`, `test`, `run`, `tidy`
- Build output to `bin/right-round`
- Version injected via `-ldflags "-X main.version=$(VERSION)"`

**`.goreleaser.yml`**: Builds for linux/darwin/windows on amd64/arm64, CGO disabled, tar.gz archives (zip for Windows), checksums

**CI workflow** (`.github/workflows/ci.yml`): On push/PR to main — Go 1.23, build, test with race detector, golangci-lint

**Release workflow** (`.github/workflows/release.yml`): On tag push `v*` — goreleaser with `GITHUB_TOKEN`

## Implementation Order

1. `go.mod`, `.gitignore` — Initialize Go module
2. `internal/data/types.go` — Struct definitions (foundation for everything)
3. `embedded.go` — Embed the JSON file
4. `internal/data/loader.go` — Parse and index data
5. `internal/tui/styles.go`, `internal/tui/keys.go` — Styles and keybindings (no deps)
6. `internal/tui/preview.go` — Animation ticker
7. `internal/tui/progressbar.go` — Bar rendering
8. `internal/tui/list.go` — Grouped list panel
9. `internal/tui/detail.go` — Detail panel
10. `internal/tui/clipboard.go` — Copy helper
11. `internal/tui/app.go` — Wire everything together
12. `cmd/right-round/main.go` — Cobra CLI entrypoint
13. `Makefile` — Build targets
14. `.goreleaser.yml`, `.github/workflows/` — CI/CD
15. `README.md` — Usage docs, installation, keybindings reference
16. Add comprehensive test coverage across data, rendering, timing, input handling, and CLI flags

## Test Plan (Comprehensive)

### Unit tests

1. `internal/data/types_test.go`
   - JSON unmarshal coverage for optional/null fields (`interval_ms`, `notes`, `phases`, `indeterminate`, `completion_states`)
   - Source fields (`raw_url`, `license_url`, `references`) are preserved
2. `internal/data/loader_test.go`
   - Loads embedded catalog and validates expected counts (433 total, 333 spinners, 100 bars)
   - Group ordering is count-desc then name-asc
   - Entry sorting within group is stable and deterministic
   - Loader rejects malformed spinner entries with empty frames
3. `internal/tui/preview_test.go`
   - Accumulator carry-over correctness (`accum -= interval`, no drift)
   - Catch-up across delayed ticks advances multiple frames correctly
   - Null interval uses 100ms default
   - Visible-only animation gating
4. `internal/tui/progressbar_test.go`
   - Rendering with only required fields (`fill`, `empty`)
   - Rendering with optional `start`, `head`, `end`
   - Phase rendering at boundary percentages
   - Indeterminate pattern rendering path
   - Unicode width handling with wide glyphs
5. `internal/tui/list_test.go`
   - Cursor navigation across group headers and item rows
   - Scroll offset and page navigation behavior
   - Filtering by name preserves correct selection mapping
6. `internal/tui/detail_test.go`
   - Detail viewport content for spinner vs progress bar entries
   - Conditional sections for notes and also-found-in
7. `internal/tui/clipboard_test.go`
   - Clipboard payload round-trips from `json.RawMessage` without field loss
   - Success and failure status messaging
8. `cmd/right-round/main_test.go`
   - `--type` lock behavior and tab-disable state
   - `--filter` initial group selection
   - Invalid flag values return user-facing errors

### Integration tests

1. Bubble Tea update loop tests that simulate key sequences (`tab`, navigation, enter, esc, search, copy)
2. Responsive layout tests for narrow (`<100`) and wide (`>=100`) terminal widths
3. End-to-end smoke test: load catalog, initialize model, run short animation cycle, render both tabs

### Tooling and quality gates

1. `go test ./... -race -coverprofile=coverage.out`
2. Enforce minimum coverage target of 80% overall and 90% for `internal/data` and `internal/tui/preview.go`
3. `golangci-lint run`
4. `goreleaser check`

## Verification

1. `make build` compiles without errors
2. `make test` passes (full unit + integration suite)
3. `make run` launches the TUI:
   - Spinners animate smoothly in the list view
   - Tab switches between Spinners and Progress Bars
   - Arrow/j/k navigation works, cursor highlights correctly
   - Enter expands detail panel with full metadata
   - Esc collapses back to list
   - `c` copies entry JSON to clipboard (verify with paste)
   - `/` opens search, filters entries by name
   - Progress bars render correctly with their character sets
   - Responsive layout adjusts when terminal is resized
4. `make lint` passes
5. `go test ./... -race -cover` meets coverage thresholds
6. `goreleaser check` validates the release config
