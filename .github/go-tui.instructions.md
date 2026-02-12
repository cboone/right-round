---
applyTo: "internal/tui/**/*.go"
---

- **Page navigation header resolution**: In the TUI list panel, `pageUp` uses `moveToNextEntry(1)` and `pageDown` uses `moveToNextEntry(-1)` intentionally. When paging lands on a group header, the cursor resolves to the nearest entry within the landed-on group, not in the movement direction. This is a deliberate UX convention validated by `TestListModel_PageNavigation`.
