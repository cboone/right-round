# right-round

A TUI for browsing 433 terminal progress indicators (333 spinners, 100 progress bars) from 26 open-source collections.

Navigate, preview live animations, and copy entries as JSON.

## Installation

### From source

```sh
go install github.com/cboone/right-round/cmd/right-round@latest
```

### From release

Download a binary from the [releases page](https://github.com/cboone/right-round/releases).

### Build locally

```sh
git clone https://github.com/cboone/right-round.git
cd right-round
make build
./bin/right-round
```

## Usage

```sh
# Launch the TUI
right-round

# Start locked to spinners only
right-round --type spinner

# Start locked to progress bars only
right-round --type progress_bar

# Start with a specific group selected
right-round --group braille
```

## Keybindings

| Key           | Action                              |
|---------------|-------------------------------------|
| `j` / `down`  | Move cursor down                    |
| `k` / `up`    | Move cursor up                      |
| `pgdn`        | Page down                           |
| `pgup`        | Page up                             |
| `home`        | Go to top                           |
| `end`         | Go to bottom                        |
| `enter` / `l` | Expand detail view (narrow mode)    |
| `esc` / `h`   | Collapse back to list               |
| `tab`         | Switch between Spinners / Progress Bars |
| `/`           | Search/filter by name               |
| `c`           | Copy selected entry as JSON         |
| `?`           | Toggle full help                    |
| `q` / `ctrl+c`| Quit                               |

## Layout

- **Wide terminals** (100+ columns): list panel and detail panel side by side
- **Narrow terminals** (under 100 columns): list and detail panels switch between each other

## Data

All 433 entries are embedded in the binary from `progress-indicators.json`, a consolidated catalog of terminal progress indicators collected from 26 open-source libraries, reference documents, gists, and blog posts.

### Entry types

- **Spinners** (333): animated frame sequences grouped into braille, line, dot, block, geometric, arrow, toggle, bounce, scroll, emoji, novelty, text, and symbol
- **Progress bars** (100): character sets for bar rendering grouped into ascii, block, geometric, decorative, phased, and emoji

## Development

```sh
make build    # Build to bin/right-round
make test     # Run tests with race detector and coverage
make lint     # Run golangci-lint
make run      # Build and run
make tidy     # Run go mod tidy
make clean    # Remove build artifacts
```

## License

See [LICENSE](LICENSE).
