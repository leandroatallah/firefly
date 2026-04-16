---
name: go-testing
description: Guidelines for table-driven tests, coverage analysis, and package-specific testing patterns.
---



# Go Testing & Coverage

Write structured table-driven tests with multiple input/output cases.

```go
tests := []struct {
    name  string
    input int
    want  int
}{
    {"Case A", 1, 2},
    {"Case B", 2, 4},
}
```

Use `go test -coverprofile` and `go tool cover` to measure and track coverage per package:

```bash
go test ./internal/engine/[package] -coverprofile=coverage.out
go tool cover -func=coverage.out
```

Prioritize packages with 0% or low coverage. Verify improvement after each test addition.
