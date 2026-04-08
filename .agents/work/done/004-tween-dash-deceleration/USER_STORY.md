# US-004 — Tween-Based Dash Deceleration

**Branch:** `004-tween-dash-deceleration`
**Bounded Context:** Physics / Entity

## Story

As a player,
I want the dash to start at full speed and smoothly decelerate to a stop,
So that the dash feels snappy and impactful rather than abruptly cut off.

## Acceptance Criteria

- AC1: On dash activation, horizontal velocity is driven by an `InOutSine` tween from `DashSpeed` to `0` over a fixed duration (e.g. 300ms / 18 frames at 60fps).
- AC2: Vertical velocity is zeroed and gravity is suspended for the duration of the dash.
- AC3: The Actor's collision rect is resized to duck height (`ResizeFixedBottom`, US-001) for the full dash duration.
- AC4: When the tween completes, the Actor transitions to `Falling` if airborne, or `Idle` if grounded.
- AC5: Dash is blocked if a solid Body is within `DashBlockDistance` pixels in the facing direction.
- AC6: One air dash is allowed per jump; the allowance resets on landing.
- AC7: A cooldown prevents re-triggering dash immediately after it ends.

## Edge Cases

- Dash into a wall: velocity is zeroed by collision, tween is abandoned, Actor transitions to `Falling`/`Idle`.
- Dash triggered while already dashing: ignored.

## Notes

- Depends on US-001 (hitbox resize during dash).
- The tween is owned by the dash skill/state, not the movement model.
- Lives in `internal/engine/physics/skill/` or `internal/game/entity/actors/states/`.
