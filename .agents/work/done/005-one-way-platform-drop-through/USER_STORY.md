# US-005 — One-Way Platform Drop-Through

**Branch:** `005-one-way-platform-drop-through`
**Bounded Context:** Physics / Entity

## Story

As a player,
I want to drop through one-way platforms by pressing down + jump,
So that I can descend to lower areas without needing a gap.

## Acceptance Criteria

- AC1: While standing on a one-way platform, pressing down + jump causes the Actor to pass through it.
- AC2: The platform's solidity is temporarily disabled for the Actor for a minimum of 1 frame, long enough to clear the top surface.
- AC3: After the Actor's bottom edge is below the platform's top edge, normal collision resumes.
- AC4: Drop-through does not trigger a jump — vertical velocity is not set to a positive (upward) value.
- AC5: Drop-through is not possible from a solid (two-way) platform.

## Edge Cases

- Holding down without jump: no drop-through, Actor ducks (US-003).
- Pressing jump without down on a one-way platform: normal jump.

## Notes

- Depends on US-003 (duck input gesture is down; must not conflict).
- Lives in `internal/engine/physics/movement/` and the platform body contract.
