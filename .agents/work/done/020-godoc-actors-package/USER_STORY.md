# US-020 — Godoc for Actors Package

**Branch:** `020-godoc-actors-package`
**Bounded Context:** Entity
**Workflow:** Vibe coding (documentation only, no logic changes)

## Story

As a developer, I want `internal/engine/entity/actors/` to have package-level godoc and struct field comments, so that the state machine architecture is understandable without reading the full source.

## Context

The `actors` package is the most complex in the engine — it contains the `Character` struct (15+ fields), the state registry, and all state implementations. None of these have package-level documentation or field comments.

Files to document:

| File | Issue |
|---|---|
| `actor_state.go` | No package-level doc |
| `character.go` | Complex struct, minimal comments |
| `ducking_state.go` | No godoc |
| `shooting_states.go` | 4 state types, zero docs |
| `state_registry.go` | No package doc |
| `actor_manager.go` | Inconsistent comments |

## Acceptance Criteria

- **AC1:** `actor_state.go` has a package-level godoc comment explaining the registry-based state pattern.
- **AC2:** `Character` struct has a type-level godoc comment and inline comments on all exported and key unexported fields.
- **AC3:** All state types (`DuckingState`, shooting states, etc.) have a one-line godoc comment.
- **AC4:** No logic, signatures, or existing behaviour is changed — comments only.
- **AC5:** `go build ./...` and `go test ./...` pass without errors after changes.
