# US-002 — Last-Pressed-Wins Horizontal Input

**Branch:** `002-last-pressed-wins-input`
**Bounded Context:** Input

## Story

As a player,
I want the most recently pressed horizontal direction to take priority when both left and right are held simultaneously,
So that direction changes feel immediate and responsive.

## Acceptance Criteria

- AC1: When only left is held, horizontal axis returns `-1`.
- AC2: When only right is held, horizontal axis returns `1`.
- AC3: When both are held and right was pressed last, axis returns `1`.
- AC4: When both are held and left was pressed last, axis returns `-1`.
- AC5: Releasing one key while the other is still held immediately returns the held key's direction.
- AC6: The logic is encapsulated and does not use global mutable state — it is injectable/testable.

## Notes

- Lives in `internal/engine/input/`.
- No dependency on physics or Actor state.
- Prerequisite for US-003 (duck/run states rely on clean directional input).
