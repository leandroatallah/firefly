---
name: coverage-analyzer
description: Runs go test coverage analysis and produces prioritized list of low-coverage files/functions
tools: Bash, Write
---


# Coverage Analyzer

## Purpose

Runs `go test -coverprofile` on target packages, parses output, and produces a prioritized list of files/functions with low or zero coverage.

## Responsibilities

- Execute coverage analysis on specified packages or entire `internal/engine` module
- Parse `go tool cover -func` output to extract per-function coverage percentages
- Identify packages and functions with <80% coverage
- Prioritize results based on project goals (entity/actors, scenes/phases, sequences)
- Generate structured report for downstream agents

## Inputs

- Target package paths (e.g., `./internal/engine/entity/actors`, `./internal/game/scenes/phases`)
- Coverage threshold (default: 80%)

## Outputs

- JSON or structured text report containing:
  - Package name
  - File path
  - Function name
  - Current coverage percentage
  - Priority level (critical/high/medium/low)

## Commands

```bash
go test ./internal/engine/... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

## Integration

Feeds report to **Gap Detector** for detailed analysis.
