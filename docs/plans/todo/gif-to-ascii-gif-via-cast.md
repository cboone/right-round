# Create `bin/gif2gif`: animated GIF to ASCII art GIF via asciinema cast

## Context

The `feature/vinyl` branch has `images/vinyl-loop.gif` (27 frames, 480x382px, frame delays alternating 6/7 centiseconds, loops forever). A prior version of a `gif2gif` script used VHS for rendering but was removed. The new approach replaces VHS with asciinema `.cast` files and `agg`, which is much simpler: extract frames, convert each to ANSI text with `chafa`, assemble into a `.cast` file, and render back to GIF with `agg`.

## Pipeline

```
vinyl-loop.gif
  -> magick -coalesce -> individual PNGs
  -> chafa --format symbols -> ANSI text per frame
  -> python3 assembles -> asciinema v2 .cast file
  -> agg -> vinyl-ascii.gif
```

## Implementation

Create `bin/gif2gif` (~120 lines) following the style of `scripts/check-coverage.sh`.

### Script structure

1. **Header**: `#!/usr/bin/env bash`, `set -euo pipefail`
2. **Hardcoded constants**:
   - Input: `$repo_root/images/vinyl-loop.gif`
   - Output: `$repo_root/images/vinyl-ascii.gif`
   - chafa size: `80x40` (chafa will produce ~80x32 actual output preserving aspect ratio)
   - agg font size: `14`
3. **Helper function**: `require_command` (same pattern as prior version)
4. **Dependency checks**: `chafa`, `magick`, `agg`, `python3`
5. **Temp directory** with `trap 'rm -rf "$work_dir"' EXIT`
6. **Extract frame delays**: `mapfile -t delays < <(magick identify -format '%T\n' ...)`
7. **Extract frames**: `magick "$input_gif" -coalesce "$work_dir/frames/frame-%03d.png"`
8. **chafa loop**: Convert each frame PNG to ANSI text file
   - `--format symbols --colors 256 --size 80x40 --animate off --polite on`
   - `--polite on` avoids cursor hide/show sequences we don't need
9. **Measure actual dimensions** from frame 0's ANSI output (rows via `wc -l`, cols via stripping ANSI escapes with python3)
10. **Build `.cast` file** with a single `python3 -c` invocation that:
    - Writes header: `{"version": 2, "width": cols, "height": rows}`
    - For each frame, reads the ANSI text file, prepends `\x1b[2J\x1b[H` (clear + home) for frames > 0, and writes `[timestamp, "o", data]`
    - Timestamps are cumulative from the GIF's centisecond delays
    - `json.dumps` handles proper encoding of escape characters
11. **Render with agg**: `agg --font-size 14 output.cast vinyl-ascii.gif`
12. **Summary output**: frame count, duration, output path

### Key design decisions

- **`--polite on`** for chafa: produces cleaner ANSI without cursor hide/show noise
- **Single python3 invocation** for cast file assembly: avoids 27 subprocess spawns and eliminates shell-escaping issues with ANSI data
- **`magick -coalesce`** for extraction: essential since the GIF uses partial-frame updates
- **Measure dimensions** rather than hardcoding: robust if chafa size constant changes later
- **Frame 0 at t=0** with no clear; subsequent frames prepend clear+home before content

## Files

- **Create**: `bin/gif2gif` (already allowed by `.gitignore` line `!bin/gif2gif`)
- **Reference**: `scripts/check-coverage.sh` (style conventions)
- **Input**: `images/vinyl-loop.gif`
- **Output**: `images/vinyl-ascii.gif` (generated, already in `.gitignore` via `bin/*` pattern -- actually this is in `images/`, needs to be gitignored or not committed)

## Prerequisites

`agg` must be installed: `brew install agg`

## Verification

1. Install agg if needed: `brew install agg`
2. Run: `bin/gif2gif`
3. Verify `images/vinyl-ascii.gif` is created and is a valid animated GIF
4. Open it: `open images/vinyl-ascii.gif`
5. Check frame count matches: `magick identify images/vinyl-ascii.gif | wc -l` (should be close to 27)
