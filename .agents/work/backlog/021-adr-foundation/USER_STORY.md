# 021 — Architecture Decision Records Foundation

**Branch:** `021-adr-foundation`
**Bounded Context:** Cross-cutting
**Workflow:** Vibe coding (documentation only, no code changes)

## Story

As a developer, I want a `docs/adr/` directory with ADRs for the five key architectural decisions in this engine, so that future contributors understand the reasoning behind non-obvious design choices.

## Context

The engine makes several non-obvious architectural choices that are not explained anywhere in the codebase. New contributors (and AI agents) repeatedly need to infer intent from code. Capturing these decisions as ADRs prevents repeated re-litigation of settled choices.

## ADRs to Create

| ADR | Decision |
|---|---|
| `ADR-001-fp16-fixed-point-arithmetic.md` | Why positions use x16 fixed-point instead of float64 |
| `ADR-002-registry-based-state-pattern.md` | Why actor states use a global registry with `init()` registration |
| `ADR-003-goroutine-audio-looping.md` | Why audio looping uses goroutines instead of Ebitengine's built-in loop |
| `ADR-004-space-body-model-physics.md` | Why physics is split into Space / Body / MovementModel layers |
| `ADR-005-composite-grounded-sub-state.md` | Why the grounded state uses a sub-state machine instead of flat states |

## ADR Format

Each file follows this structure:

```
# ADR-NNN — Title

## Status
Accepted

## Context
Why this decision needed to be made.

## Decision
What was decided.

## Consequences
Trade-offs and implications.
```

## Acceptance Criteria

- **AC1:** `docs/adr/` directory exists with all five ADR files.
- **AC2:** Each ADR accurately reflects the actual implementation (verified against source code).
- **AC3:** Each ADR is concise — context, decision, and consequences fit on one screen.
- **AC4:** No source code is modified.
