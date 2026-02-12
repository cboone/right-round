# CLI Spinner Collection Summary

A collection of text/CLI-based spinner character sets gathered from across the
open source ecosystem. These are the animated loading indicators used in
terminal applications, made of character sequences that cycle through frames to
create the illusion of motion.

Collected on 2026-02-12. No deduplication has been performed.

## Files

### cli-spinners.json
- **Source:** [sindresorhus/cli-spinners](https://github.com/sindresorhus/cli-spinners)
- **License:** MIT
- **Spinners:** ~95 named spinners with interval (ms) and frames
- **Notes:** The canonical collection. Used by ora (Node.js), halo (Python),
  yaspin (Python), rich (Python), and many other libraries across languages.
  The single most widely referenced spinner data source.

### briandowns-spinner.json
- **Source:** [briandowns/spinner](https://github.com/briandowns/spinner)
- **License:** Apache-2.0
- **Spinners:** 93 indexed character sets (0-90 plus 2 dynamic clock sets)
- **Notes:** Popular Go package. Includes unique sets like katakana characters,
  card suits, circled numbers, fractions, and alphabet cycling. Also ported to
  C as [briandowns/libspinner](https://github.com/briandowns/libspinner).

### antofthy-spinners.txt
- **Source:** [Anthony Thyssen's Spinners.txt](https://antofthy.gitlab.io/info/ascii/Spinners.txt)
- **License:** No explicit license (reference/catalog)
- **Spinners:** 80+ character sequences across multiple categories
- **Notes:** The most comprehensive text-based reference. Many of the longer
  Braille clock and cycle sequences were figured out by the author and are
  original. Covers ASCII, Unicode, Braille (6-dot and 8-dot), emoji, half-line
  glyphs, box graphics, and progress bar characters. Distinguishes between
  patrol (bounce), clock (ticking hand), and cycle (passing hand) animation
  types. Last updated Dec 2025.

### charmbracelet-bubbles.json
- **Source:** [charmbracelet/bubbles](https://github.com/charmbracelet/bubbles)
- **License:** MIT
- **Spinners:** 12 built-in types
- **Notes:** Spinner component for the Bubble Tea TUI framework (Go). Also used
  by the [gum](https://github.com/charmbracelet/gum) CLI tool via `gum spin`.

### tty-spinner-formats.json
- **Source:** [piotrmurach/tty-spinner](https://github.com/piotrmurach/tty-spinner)
- **License:** MIT
- **Spinners:** 44 predefined formats
- **Notes:** Part of the TTY Ruby toolkit. Includes formats like flip, toss,
  lighthouse, burger, dance (fish), shark, and pong.

### bash-loading-animations.json
- **Source:** [Silejonu/bash_loading_animations](https://github.com/Silejonu/bash_loading_animations)
- **License:** GPL-3.0
- **Spinners:** 38 animations (8 ASCII, 30 UTF-8)
- **Notes:** Ready-to-use for Bash scripts. Includes emoji-based spinners like
  blink, sick, monkey, and camera. Intervals specified in seconds.

### swelljoe-spinner.json
- **Source:** [swelljoe/spinner](https://github.com/swelljoe/spinner)
- **License:** Not explicitly stated
- **Spinners:** 34 definitions (5 ASCII, 24 Unicode, 5 wide multi-character)
- **Notes:** POSIX shell compatible. Includes unique sets like hippie
  (peace/love symbols), pointing hands, card suits, and rotating arrows. Wide
  spinners include ASCII progress, propeller, and snake animations.

### throbber-widgets-tui.json
- **Source:** [arkbig/throbber-widgets-tui](https://github.com/arkbig/throbber-widgets-tui)
- **License:** MIT
- **Spinners:** 21 symbol sets
- **Notes:** Rust TUI widget for ratatui. Each set has full (completion) and
  empty (idle) marker characters in addition to animation frames. Includes
  unique Ogham script, Canadian Aboriginal Syllabics, and parenthesis-based
  spinners.

### asika32764-gist.json
- **Source:** [asika32764 GitHub Gist](https://gist.github.com/asika32764/19956edcc5e893b2cbe3768e91590cf1)
- **License:** No license stated
- **Spinners:** 17 sets
- **Notes:** Ruby implementation. Notable for circle-tail spinners using
  Canadian Aboriginal Syllabics characters (b/d/p/q with tails).

### ora-spinners-extended.json
- **Source:** [RubenVerg/ora_spinners](https://github.com/RubenVerg/ora_spinners)
- **License:** ISC
- **Spinners:** 5 sliding-dot variants
- **Notes:** Extends the cli-spinners standard with succeed/fail/warn/info
  terminal state symbols using filled/empty circles.

### willcarhart-blog.json
- **Source:** [willcarh.art blog](https://willcarh.art/blog/how-to-write-better-bash-spinners)
- **License:** Blog post (no explicit code license)
- **Spinners:** 3 sets
- **Notes:** Tutorial on Bash spinners with a bouncing bullet animation.

### odino-blog.json
- **Source:** [odino.org blog](https://odino.org/command-line-spinners-the-amazing-tale-of-modern-typewriters-and-digital-movies/)
- **License:** Blog post (no explicit code license)
- **Spinners:** 2 sets
- **Notes:** Notable for the quadrant block spinner using half-block characters.

### rosetta-code.json
- **Source:** [Rosetta Code](https://rosettacode.org/wiki/Spinning_rod_animation/Text)
- **License:** GNU FDL 1.2
- **Spinners:** 11 sets
- **Notes:** From the wiki task page with implementations in dozens of
  languages. Includes unique bird flying (CJK bracket) and plant/flower
  spinners.

## Ecosystem Overview

The `sindresorhus/cli-spinners` JSON file is the de facto standard data source,
referenced by libraries in:

| Language | Libraries using cli-spinners data |
|----------|----------------------------------|
| Node.js  | ora, @topcli/spinner, spinnies, cli-spinner |
| Python   | halo, yaspin, rich, py-spinners |
| Go       | yacspin (also draws from briandowns/spinner) |
| Ruby     | tty-spinner (has its own formats too) |
| Rust     | indicatif (references it for custom tick_chars) |
| Swift    | CLISpinner |
| R        | cli (r-lib) |
| Elixir   | cli_spinners |
| Crystal  | spinner-frames.cr |

## Character Categories

Spinners in this collection use characters from these Unicode blocks:

- **ASCII:** Classic `-\|/`, dots, progress bars
- **Box Drawing:** `┤┘┴└├┌┬┐`, bold variants `┫┛┻┗┣┏┳┓`
- **Block Elements:** `▁▂▃▄▅▆▇█`, `▏▎▍▌▋▊▉`, `░▒▓`, `▖▘▝▗▌▀▐▄`
- **Braille Patterns:** 6-dot (`⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏`) and 8-dot variants
- **Geometric Shapes:** `◢◣◤◥`, `◰◳◲◱`, `◴◷◶◵`, `◐◓◑◒`
- **Arrows:** `←↖↑↗→↘↓↙`, `⇐⇑⇒⇓`, `▹▸`
- **Stars/Symbols:** `✶✸✹✺`, `☱☲☴`, `♠♣♥♦`
- **Emoji:** Clocks, moons, earths, faces, hands, weather
- **Miscellaneous:** Katakana, Ogham, Canadian Aboriginal Syllabics, CJK punctuation
