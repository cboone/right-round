# GitHub Copilot Instructions

## PR Review

- **Go version**: The Go version in `go.mod` and `.github/workflows/ci.yml` are the source of truth. Do not compare against PR description text for version consistency.
- **golangci-lint version**: This project intentionally uses `version: latest` for golangci-lint. The target Go version is pinned separately via `--go` flag. Do not flag `version: latest` as a concern.
- **gofmt compliance**: Run `gofmt -d` before flagging formatting issues. Chained method calls with tab-indented dot-prefixed lines are standard Go formatting for fluent-style APIs like lipgloss.
