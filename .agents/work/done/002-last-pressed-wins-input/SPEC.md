# SPEC-002 — Last-Pressed-Wins Horizontal Input

**Branch:** `002-last-pressed-wins-input`
**Bounded Context:** Input
**Package:** `internal/engine/input/`

## Technical Requirements

New type `HorizontalAxis` in `internal/engine/input/`:

```go
type HorizontalAxis struct { /* unexported fields */ }

func NewHorizontalAxis() *HorizontalAxis
func (h *HorizontalAxis) Press(dir int)   // dir: -1 (left) or +1 (right)
func (h *HorizontalAxis) Release(dir int)
func (h *HorizontalAxis) Value() int      // returns -1, 0, or +1
```

- Internally tracks which keys are held and which was pressed last.
- `Value()` logic:
  - Neither held → `0`
  - Only one held → that direction
  - Both held → direction of the last `Press()` call
- No global state; each instance is independent.
- Caller (input handler) maps raw key events to `Press`/`Release` calls.

## Pre-conditions

- `Press` is called when a directional key goes down.
- `Release` is called when a directional key goes up.

## Post-conditions

- `Value()` always returns one of `{-1, 0, 1}`.
- After `Release` of the last-pressed key while the other is still held, `Value()` returns the still-held direction.

## Integration Points

- No existing contracts are modified.
- Consumed by US-003 (duck/run state input reads).
- No dependency on physics or Actor.

## Red Phase — Failing Test Scenario

File: `internal/engine/input/horizontal_axis_test.go`

Table-driven test `TestHorizontalAxisValue`:

| case | actions | expected Value() |
|---|---|---|
| only left held | Press(-1) | -1 |
| only right held | Press(1) | 1 |
| both held, right last | Press(-1), Press(1) | 1 |
| both held, left last | Press(1), Press(-1) | -1 |
| release last-pressed, other still held | Press(-1), Press(1), Release(1) | -1 |
| release all | Press(1), Release(1) | 0 |

Test must fail (type does not exist yet) → implement → test passes.
