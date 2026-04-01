# SPEC 008 — Alternating Bullet Spawn Offset (OffsetToggler)

**Branch:** `008-bullet-spawn-offset-toggler`
**Bounded Context:** Entity / Game Logic
**Package:** `internal/game/` (alongside shooting skill or player actor)

## Technical Requirements

```go
type OffsetToggler struct {
    offset  int
    current int
}

func NewOffsetToggler(offset int) *OffsetToggler
func (o *OffsetToggler) Next() int
```

- `NewOffsetToggler(n)`: stores `offset = n`, initialises `current = n` (first call returns `+n`).
- `Next()`: flips sign of `current` (`current = -current`), returns new `current`.
  - Sequence: `+n`, `-n`, `+n`, `-n`, …
- One instance per shooting Actor/skill — never shared globally.
- Applied to bullet spawn Y-position only; Actor position is untouched.

## Pre-conditions

- `offset` may be any integer (including 0 or negative — sign-flip still works).

## Post-conditions

- Consecutive `Next()` calls strictly alternate sign.
- `|Next()| == offset` for all calls (when offset != 0).
- Actor position is never modified.

## Integration Points

- No existing contracts are modified.
- Owned by the shooting skill/state (created in `NewShootingSkill` or equivalent).
- Offset value passed to bullet factory at spawn time.
- Depends on SPEC-006 (shooting is a grounded sub-state behaviour) for context, but `OffsetToggler` itself has zero dependencies.

## Red Phase — Failing Test Scenario

File: `internal/game/offset_toggler_test.go` (or alongside the shooting skill)

`TestOffsetTogglerSequence`:

| call # | offset | expected Next() |
|---|---|---|
| 1 | 4 | +4 |
| 2 | 4 | -4 |
| 3 | 4 | +4 |
| 4 | 4 | -4 |

`TestOffsetTogglerZero`:
- `NewOffsetToggler(0)` → all `Next()` calls return `0`.

`TestOffsetTogglerNegativeInit`:
- `NewOffsetToggler(-3)` → sequence: `-3`, `+3`, `-3`, `+3`.

Test must fail (type does not exist yet) → implement → test passes.
