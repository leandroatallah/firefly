# US-019 — Godoc for Contracts Interfaces

**Branch:** `019-godoc-contracts-interfaces`
**Bounded Context:** Contracts
**Workflow:** Vibe coding (documentation only, no logic changes)

## Story

As a developer, I want all interfaces in `internal/engine/contracts/` to have godoc comments, so that the engine's public API is self-documenting and IDE tooling surfaces useful descriptions.

## Context

Most interfaces in `internal/engine/contracts/` have no documentation. Only `one_way_platform.go`, `shooter.go`, and `sequences.go` have any comments. The following files need godoc added:

| File | Issue |
|---|---|
| `body/body.go` | 50+ line interface, no comments |
| `navigation/navigation.go` | No interface or method docs |
| `vfx/vfx.go` | 15+ methods, no godoc |
| `animation/animation.go` | No docs |
| `context/context.go` | No docs |
| `scene/freeze.go` | No docs |

## Acceptance Criteria

- **AC1:** Every interface type in the listed files has a godoc comment explaining its purpose.
- **AC2:** Every method in those interfaces has a one-line godoc comment.
- **AC3:** No logic, signatures, or existing behaviour is changed — comments only.
- **AC4:** `go build ./...` and `go test ./...` pass without errors after changes.
