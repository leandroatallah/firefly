# US-009 — Variable Jump Height

**Branch:** `009-variable-jump-height`

## Story

As a player, I want the jump height to vary based on how long I hold the jump button, so that I have precise control over my movement like in classic platformer games.

## Bounded Context

- **Physics** (`internal/engine/physics/`) — controls `Body` velocity during jump
- **Game Logic** (`internal/game/`) — reads input and drives the jump state

## Background

The `Movable` contract already exposes `TryJump(force int)` and `SetJumpForceMultiplier(multiplier float64)`. Variable jump height is achieved by **cutting the upward velocity** when the player releases the jump button early (the "jump cut" technique), rather than applying a different initial force.

- Full press → `Body` reaches maximum jump apex.
- Short press → upward `vy16` is reduced by a configurable `JumpCutMultiplier` on button release, producing a lower arc.

No new contract methods are required.

## Acceptance Criteria

| # | Criterion |
|---|---|
| AC-1 | When the player presses the jump button and holds it, the `Actor` reaches the full jump apex defined by `TryJump(force)`. |
| AC-2 | When the player releases the jump button while the `Body` is still moving upward (`vy16 < 0`), the upward velocity is multiplied by a configurable `JumpCutMultiplier` (e.g. `0.5`), immediately reducing the apex. |
| AC-3 | `JumpCutMultiplier` is in the range `(0.0, 1.0]`. A value of `1.0` disables the cut (full jump always). |
| AC-4 | If the `Body` is already falling (`vy16 >= 0`) when the button is released, no velocity change is applied. |
| AC-5 | The jump cut is applied at most once per jump (re-pressing the button mid-air does not trigger another cut). |
| AC-6 | The behaviour is deterministic: given the same input frame sequence, the resulting apex is always identical. |

## Behavioral Edge Cases

- **Button released on the same frame as `TryJump`:** the cut is still applied (short-tap = minimum jump).
- **Button held across a scene transition / freeze:** the jump-cut state resets; the next jump starts fresh.
- **`JumpCutMultiplier = 0.0` is invalid** and must be rejected (or clamped to a minimum, e.g. `0.1`) to avoid zeroing velocity entirely.
- **Coyote-time / buffered jumps:** if the project later adds coyote time, the cut logic must remain decoupled from grounded detection — it only cares about `vy16`.
