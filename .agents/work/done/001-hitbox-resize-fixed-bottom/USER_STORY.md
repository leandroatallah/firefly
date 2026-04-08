# US-001 — Hitbox Resize Anchored to Bottom

**Branch:** `001-hitbox-resize-fixed-bottom`
**Bounded Context:** Physics

## Story

As a game developer,
I want to resize an Actor's collision rect while keeping its bottom edge anchored,
So that duck and dash states can shrink the hitbox without visually repositioning the Actor.

## Acceptance Criteria

- AC1: `ResizeFixedBottom(rect image.Rectangle, newHeight int) image.Rectangle` returns a rect with the same bottom edge and the given height.
- AC2: `ResizeFixedTop(rect image.Rectangle, newHeight int) image.Rectangle` returns a rect with the same top edge and the given height.
- AC3: The input rect is never mutated.
- AC4: If `newHeight <= 0`, the result is a zero-height rect anchored at the respective edge.

## Notes

- Pure utility functions, no physics or state dependencies.
- Lives in `internal/engine/physics/body`.
- Prerequisite for US-003 (duck hitbox) and US-004 (tween dash).
