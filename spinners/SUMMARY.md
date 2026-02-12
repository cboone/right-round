# CLI Progress Indicator Collection Summary

A collection of text/CLI-based progress indicators gathered from across the open
source ecosystem. Covers spinners (animated loading indicators), progress bars,
sparklines, and related character sets used in terminal applications.

Collected on 2026-02-12. No deduplication has been performed on the per-source files under `spinners/`; a deduplicated unified list is available in `progress-indicators.json`.

## Spinner Collections

### cli-spinners.json

- **Source:** [sindresorhus/cli-spinners](https://github.com/sindresorhus/cli-spinners)
- **License:** MIT
- **Items:** ~95 named spinners with interval (ms) and frames
- **Notes:** The canonical collection. Used by ora (Node.js), halo (Python),
  yaspin (Python), rich (Python), and many other libraries across languages.
  The single most widely referenced spinner data source.

### briandowns-spinner.json

- **Source:** [briandowns/spinner](https://github.com/briandowns/spinner)
- **License:** Apache-2.0
- **Items:** 93 indexed character sets (0-90 plus 2 dynamic clock sets)
- **Notes:** Popular Go package. Includes unique sets like katakana characters,
  card suits, circled numbers, fractions, and alphabet cycling.

### antofthy-spinners.txt

- **Source:** [Anthony Thyssen's Spinners.txt](https://antofthy.gitlab.io/info/ascii/Spinners.txt)
- **License:** No explicit license (reference/catalog)
- **Items:** 80+ character sequences across multiple categories
- **Notes:** The most comprehensive text-based reference. Many of the longer
  Braille clock and cycle sequences were originally figured out by the author.
  Covers ASCII, Unicode, Braille (6-dot and 8-dot), emoji, half-line glyphs,
  box graphics, and progress bar characters. Last updated Dec 2025.

### charmbracelet-bubbles.json

- **Source:** [charmbracelet/bubbles](https://github.com/charmbracelet/bubbles)
- **License:** MIT
- **Items:** 12 built-in spinner types
- **Notes:** For the Bubble Tea TUI framework (Go). Also used by
  [gum](https://github.com/charmbracelet/gum) via `gum spin`.

### tty-spinner-formats.json

- **Source:** [piotrmurach/tty-spinner](https://github.com/piotrmurach/tty-spinner)
- **License:** MIT
- **Items:** 44 predefined spinner formats
- **Notes:** Part of the TTY Ruby toolkit. Includes flip, toss, lighthouse,
  burger, dance (fish), shark, and pong.

### bash-loading-animations.json

- **Source:** [Silejonu/bash_loading_animations](https://github.com/Silejonu/bash_loading_animations)
- **License:** GPL-3.0
- **Items:** 38 animations (8 ASCII, 30 UTF-8)
- **Notes:** Ready-to-use for Bash scripts. Intervals specified in seconds.

### swelljoe-spinner.json

- **Source:** [swelljoe/spinner](https://github.com/swelljoe/spinner)
- **License:** Not explicitly stated
- **Items:** 34 definitions (5 ASCII, 24 Unicode, 5 wide multi-character)
- **Notes:** POSIX shell compatible. Includes hippie, hands, card suits, and
  rotating arrows. Wide spinners include progress, propeller, and snake.

### throbber-widgets-tui.json

- **Source:** [arkbig/throbber-widgets-tui](https://github.com/arkbig/throbber-widgets-tui)
- **License:** MIT
- **Items:** 21 symbol sets
- **Notes:** Rust TUI widget for ratatui. Includes Ogham script, Canadian
  Aboriginal Syllabics, and parenthesis-based spinners.

### asika32764-gist.json

- **Source:** [asika32764 GitHub Gist](https://gist.github.com/asika32764/19956edcc5e893b2cbe3768e91590cf1)
- **License:** No license stated
- **Items:** 17 sets
- **Notes:** Ruby implementation with circle-tail spinners.

### ora-spinners-extended.json

- **Source:** [RubenVerg/ora_spinners](https://github.com/RubenVerg/ora_spinners)
- **License:** ISC
- **Items:** 5 sliding-dot variants
- **Notes:** Extends cli-spinners with succeed/fail/warn/info terminal states.

### rosetta-code.json

- **Source:** [Rosetta Code](https://rosettacode.org/wiki/Spinning_rod_animation/Text)
- **License:** GNU FDL 1.2
- **Items:** 11 sets
- **Notes:** Unique bird flying (CJK bracket) and plant/flower spinners.

### willcarhart-blog.json

- **Source:** [willcarh.art blog](https://willcarh.art/blog/how-to-write-better-bash-spinners)
- **License:** Blog post (no explicit code license)
- **Items:** 3 sets

### odino-blog.json

- **Source:** [odino.org blog](https://odino.org/command-line-spinners-the-amazing-tale-of-modern-typewriters-and-digital-movies/)
- **License:** Blog post (no explicit code license)
- **Items:** 2 sets (including quadrant block spinner)

## Progress Bar Character Collections

### unicode-progress-bar-characters.json

- **Source:** Compiled from multiple references
- **References:** [changaco.oy.lc](https://changaco.oy.lc/unicode-progress-bars/),
  [ascii.bar](https://www.ascii.bar/),
  [mike42.me](https://mike42.me/blog/2018-06-make-better-cli-progress-bars-with-unicode-block-characters),
  [naut.ca](https://www.naut.ca/blog/2024/12/26/making-a-unicode-progress-bar/),
  [rougier gist](https://gist.github.com/rougier/c0d31f5cbdaac27b876c),
  [dernocua](https://dernocua.github.io/notes/unicode-graphics.html),
  [Cygra/interesting-unicode-symbols](https://github.com/Cygra/interesting-unicode-symbols/issues/1),
  [Rosetta Code sparklines](https://rosettacode.org/wiki/Sparkline_in_unicode)
- **License:** Unicode characters are public domain
- **Items:** 9 categories of characters
- **Notes:** The comprehensive Unicode character catalog. Covers horizontal
  fractional blocks (8 sub-steps), vertical blocks (sparklines), shade blocks,
  quadrant blocks, Braille progress fills, geometric filled/empty pairs, circle
  fill progressions, halfwidth CJK forms, ASCII styles, and bracket delimiters.

### tty-progressbar-formats.json

- **Source:** [piotrmurach/tty-progressbar](https://github.com/piotrmurach/tty-progressbar)
- **License:** MIT
- **Items:** 22 preconfigured bar formats
- **Notes:** Each defines complete, incomplete, and unknown (indeterminate)
  characters. Creative styles: heart, star, diamond, crate, burger, chevron,
  bracket, tread, track, wave.

### verigak-progress.json

- **Source:** [verigak/progress](https://github.com/verigak/progress)
- **License:** ISC
- **Items:** 7 bar types + 5 spinner types
- **Notes:** Python library. IncrementalBar uses 9-phase fractional blocks,
  PixelBar uses 8-phase Braille, ShadyBar uses 5-phase shading.
  FillingSquaresBar and FillingCirclesBar use geometric pairs.

### cli-progress-presets.json

- **Source:** [npkgz/cli-progress](https://github.com/npkgz/cli-progress)
- **License:** MIT
- **Items:** 4 built-in presets
- **Notes:** Node.js library. legacy (`=`/`-`), shades_classic (`█`/`░`),
  shades_grey (same with ANSI grey), rect (`■`/` `).

### tqdm-bar-characters.json

- **Source:** [tqdm/tqdm](https://github.com/tqdm/tqdm)
- **License:** MIT AND MPL-2.0
- **Items:** 2 built-in character sets + custom support
- **Notes:** The most popular Python progress bar (44M+ weekly pip installs).
  Bar.ASCII uses `" 123456789#"`, Bar.UTF uses fractional blocks. Custom
  strings supported via the `ascii` parameter.

### schollz-progressbar.json

- **Source:** [schollz/progressbar](https://github.com/schollz/progressbar)
- **License:** MIT
- **Items:** 3 predefined themes + examples
- **Notes:** Go library. Theme struct: Saucer, SaucerHead, SaucerPadding,
  BarStart, BarEnd. Includes a PAC-MAN example.

### cheggaaa-pb.json

- **Source:** [cheggaaa/pb](https://github.com/cheggaaa/pb)
- **License:** BSD-3-Clause
- **Items:** v1 FORMAT string + v3 template function
- **Notes:** Go library. v1 uses `[=>-]` format string. v3 uses Go templates
  with a `bar` function taking 5 character arguments. Supports animated
  current-position via cycle function.

### fancybar-python.json

- **Source:** [jenca-adam/fancybar](https://github.com/jenca-adam/fancybar)
- **License:** MIT
- **Items:** 6 named bar types
- **Notes:** Python library. Includes triangle (`▶`/`▷`), gradient (left half
  block with truecolor), and standard ASCII types.

### ruby-progressbar.json

- **Source:** [jfelchner/ruby-progressbar](https://github.com/jfelchner/ruby-progressbar)
- **License:** MIT
- **Items:** Default characters + PAC-MAN example
- **Notes:** Ruby library. Famous PAC-MAN bar uses Canadian Syllabics Carrier
  KHEE and Halfwidth Katakana Middle Dot.

### indicatif-rust.json

- **Source:** [console-rs/indicatif](https://github.com/console-rs/indicatif)
- **License:** MIT
- **Items:** Common progress_chars and tick_chars examples
- **Notes:** Rust library. No hardcoded themes; template-based system where
  users compose styles from character strings. progress_chars takes 3+
  characters for filled/current/empty.

### social-media-progress-bars.json

- **Source:** Compiled from [copy-paste.net](https://copy-paste.net/en/loading-bar.php),
  [emojicombos.com](https://emojicombos.com/ascii-progress-bar),
  [emojimagic.io](https://emojimagic.io/progress-bar),
  [textgenerator.net](https://textgenerator.net/patterns/loading-symbol), and others
- **License:** Public domain Unicode characters
- **Items:** 13 style categories
- **Notes:** Static text progress bars for social media (Discord, TikTok,
  Instagram). Covers block, shade, circle, square, parallelogram, emoji
  colored squares, and decorative star/heart/diamond styles.

### osc-9-4-terminal-progress.json

- **Source:** [ConEmu](https://conemu.github.io/en/AnsiEscapeCodes.html),
  [Windows Terminal](https://learn.microsoft.com/en-us/windows/terminal/tutorials/progress-bar-sequences),
  [rockorager.dev](https://rockorager.dev/misc/osc-9-4-progress-bars/)
- **License:** Protocol specification
- **Items:** 5 progress states
- **Notes:** Not character-based, but an escape sequence protocol for native
  terminal progress rendering. Supported by ConEmu, Windows Terminal, iTerm2,
  Ghostty, Konsole, and mintty.

## Ecosystem Overview

### Spinner Data Sources

The `sindresorhus/cli-spinners` JSON file is the de facto standard, referenced
by libraries in:

| Language | Libraries                                       |
| -------- | ----------------------------------------------- |
| Node.js  | ora, @topcli/spinner, spinnies, cli-spinner     |
| Python   | halo, yaspin, rich, py-spinners                 |
| Go       | yacspin (also draws from briandowns/spinner)    |
| Ruby     | tty-spinner (has its own formats too)           |
| Rust     | indicatif (references it for custom tick_chars) |
| Swift    | CLISpinner                                      |
| R        | cli (r-lib)                                     |
| Elixir   | cli_spinners                                    |
| Crystal  | spinner-frames.cr                               |

### Progress Bar Character Conventions

Most progress bar libraries define these character slots:

| Slot                     | Description       | Common defaults         |
| ------------------------ | ----------------- | ----------------------- |
| Fill/Saucer/Complete     | Filled portion    | `=`, `#`, `█`, `▰`, `■` |
| Head/Tip/Current         | Leading edge      | `>`, `█`, `▸`, none     |
| Empty/Padding/Incomplete | Remaining portion | `-`, ` `, `░`, `▱`, `□` |
| Start/Prefix             | Left delimiter    | `[`, `\|`, `▐`          |
| End/Suffix               | Right delimiter   | `]`, `\|`, `▌`          |

### Sub-Character Resolution

For smoother progress bars, several libraries support fractional fills:

| Technique              | Characters      | Resolution    |
| ---------------------- | --------------- | ------------- |
| Left fractional blocks | `▏▎▍▌▋▊▉█`      | 8 steps/char  |
| Shade progression      | `░▒▓█`          | 4 steps/char  |
| Braille fill           | `⡀⡄⡆⡇⣇⣧⣷⣿`      | 8 steps/char  |
| tqdm ASCII digits      | `" 123456789#"` | 10 steps/char |

## Character Categories

This collection uses characters from these Unicode blocks:

- **ASCII:** Classic `-\|/`, `=#-.>`, progress bars
- **Box Drawing:** `┤┘┴└├┌┬┐`, bold variants `┫┛┻┗┣┏┳┓`
- **Block Elements:** Horizontal `▏▎▍▌▋▊▉█`, vertical `▁▂▃▄▅▆▇█`,
  shades `░▒▓`, quadrants `▖▘▝▗▌▀▐▄▛▜▙▟`
- **Braille Patterns:** 6-dot spinners, 8-dot spinners, progress fills,
  sparklines
- **Geometric Shapes:** Triangles `◢◣◤◥`, squares `◰◳◲◱`, circles
  `◴◷◶◵`, `◐◓◑◒`, filled/empty pairs `●○`, `■□`, `▰▱`, `▶▷`,
  `★☆`, `♥♡`, `♦♢`
- **Arrows:** `←↖↑↗→↘↓↙`, `⇐⇑⇒⇓`, `▹▸`, `➞➟➠➡`
- **Miscellaneous Symbols:** `✶✸✹✺`, `☱☲☴`, `♠♣♥♦`, `⦿⦾`, `◉◎`
- **Emoji:** Clocks, moons, earths, faces, hands, weather, colored squares
- **CJK/Halfwidth:** `￭` (U+FFED), `･` (U+FF65), `｢｣`
- **Scripts:** Katakana, Ogham, Canadian Aboriginal Syllabics, CJK punctuation
- **Mathematical:** `⊶⊷`, `≡`, `⋅∙`, `⎺⎻⎼⎽`

## File Count

- **25 collection files** (13 spinner, 11 progress bar, 1 protocol spec)
- **1 summary document** (this file)
