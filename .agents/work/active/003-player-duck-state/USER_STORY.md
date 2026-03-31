# USER STORY 003 — Duck State with Reduced Hitbox

**Branch:** `003-player-duck-state`
**Bounded Context:** Entity / Physics

## Story

As a player,
I want to duck by holding a down input,
So that I can avoid high obstacles and fit through low passages.

## Acceptance Criteria

- AC1: While duck input is held, the Actor transitions to a `Ducking` state.
- AC2: On entering `Ducking`, the collision rect is resized to half height using `ResizeFixedBottom` (US-001).
- AC3: Horizontal velocity is set to zero while ducking.
- AC4: On releasing duck input, the Actor returns to `Idle` and the full hitbox is restored.
- AC5: The Actor cannot jump while in `Ducking` state.
- AC6: `Ducking` → `Idle` transition only occurs when the full-height hitbox has clearance above the Actor (no ceiling collision).

## Edge Cases

- Ducking while moving: Actor stops moving and enters `Ducking`.
- Ducking while falling: not applicable — duck is a grounded-only state.

## Notes

- Depends on US-001 (hitbox resize) and US-002 (clean directional input).
- Lives in `internal/game/entity/actors/states/`.
