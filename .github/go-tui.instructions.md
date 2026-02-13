---
applyTo: "internal/tui/**/*.go"
---

- **Page navigation header resolution**: In the TUI list panel, `pageUp` uses `moveToNextEntry(1)` and `pageDown` uses `moveToNextEntry(-1)` intentionally. When paging lands on a group header, the cursor resolves to the nearest entry within the landed-on group, not in the movement direction. This is a deliberate UX convention validated by `TestListModel_PageNavigation`.
- **Progress bar head width handling**: The `renderBar` function in `progressbar.go` already accounts for multi-cell head glyphs using `lipgloss.Width(*chars.Head)` and `fillCells -= headWidth`. Do not flag this as assuming a 1-cell glyph; the width calculation is correct and tested in `TestRenderBar_UnicodeWidthHandling`.
