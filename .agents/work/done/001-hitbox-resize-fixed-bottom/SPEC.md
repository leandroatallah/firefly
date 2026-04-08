# SPEC-001 — Hitbox Resize Anchored to Bottom

**Branch:** `001-hitbox-resize-fixed-bottom`
**Bounded Context:** Physics
**Package:** `internal/engine/physics/body`

## Technical Requirements

Two pure functions added to `internal/engine/physics/body`:

```go
func ResizeFixedBottom(rect image.Rectangle, newHeight int) image.Rectangle
func ResizeFixedTop(rect image.Rectangle, newHeight int) image.Rectangle
```

- No receiver, no side effects, no dependencies.
- `ResizeFixedBottom`: keeps `rect.Max.Y`, sets `rect.Min.Y = rect.Max.Y - newHeight`.
- `ResizeFixedTop`: keeps `rect.Min.Y`, sets `rect.Max.Y = rect.Min.Y + newHeight`.
- If `newHeight <= 0`, clamp to `0` — result is a zero-height rect at the anchor edge.
- Input rect is never mutated (value semantics — `image.Rectangle` is already a value type).

## Pre-conditions

- `rect` is any valid `image.Rectangle` (may be zero-value).
- `newHeight` is any integer.

## Post-conditions

- Returned rect has `Height() == max(newHeight, 0)`.
- `ResizeFixedBottom`: `returned.Max.Y == rect.Max.Y`.
- `ResizeFixedTop`: `returned.Min.Y == rect.Min.Y`.
- Original `rect` is unchanged.

## Integration Points

- No existing contracts are modified.
- Consumed by US-003 (`Ducking` state enter/exit) and US-004 (dash state enter/exit).

## Red Phase — Failing Test Scenario

File: `internal/engine/physics/body/resize_test.go`

Table-driven test `TestResizeFixedBottom` and `TestResizeFixedTop`:

| case | input rect | newHeight | expected Min.Y | expected Max.Y |
|---|---|---|---|---|
| normal shrink (bottom) | `{0,0,32,64}` | 32 | 32 | 64 |
| normal grow (bottom) | `{0,0,32,64}` | 80 | -16 | 64 |
| zero height (bottom) | `{0,0,32,64}` | 0 | 64 | 64 |
| negative height (bottom) | `{0,0,32,64}` | -5 | 64 | 64 |
| normal shrink (top) | `{0,0,32,64}` | 32 | 0 | 32 |
| zero height (top) | `{0,0,32,64}` | 0 | 0 | 0 |

Test must fail (functions do not exist yet) → implement → test passes.
