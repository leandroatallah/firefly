# US-008 — Alternating Bullet Spawn Offset (OffsetToggler)

**Branch:** `008-bullet-spawn-offset-toggler`
**Bounded Context:** Entity / Game Logic

## Story

As a game designer,
I want each bullet fired to alternate its spawn Y-offset between +N and -N pixels,
So that rapid fire produces a double-barrel visual effect without additional state management.

## Acceptance Criteria

- AC1: An `OffsetToggler` type holds a fixed offset value and flips its sign on each call to `Next() int`.
- AC2: Consecutive calls to `Next()` alternate: `+N`, `-N`, `+N`, `-N`, …
- AC3: The toggler is owned by the shooting Actor/skill, not shared globally.
- AC4: The offset is applied to the bullet's spawn position, not to the Actor's position.
- AC5: `OffsetToggler` is covered by a unit test asserting the alternating sequence over at least 4 calls.

## Notes

- Depends on US-006 (shooting is a grounded sub-state behaviour).
- Tiny type; lives alongside the shooting skill or player actor in `internal/game/`.
