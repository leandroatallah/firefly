# USER STORY — 068-actor-json-sprite-render-offset

**Branch:** `068-actor-json-sprite-render-offset`
**Bounded Context:** Kit (`internal/kit/actors/`, `internal/engine/data/schemas/`)

---

## Story

As a beat-em-up game developer,
I want to declare an optional per-state `render_offset` in an actor's JSON sprite config,
so that artwork that is not centered the same way as the idle pose can be nudged visually without moving the body rect or collision rects.

---

## Background

The render system positions every animation frame at the same body-rect origin. `cody-melee-0.png` is a 64×32 sprite whose character art sits toward the right half of the canvas, causing a visible leftward jump when the idle-to-melee transition plays. There is no mechanism to compensate without editing art. A per-state pixel offset applied only at draw time — not in physics — solves this without touching collision or movement.

---

## Acceptance Criteria

- AC-1: `schemas.AssetData` gains an optional `RenderOffset *SpriteOffset` field tagged `json:"render_offset,omitempty"`; absent field parses without error and leaves the pointer nil.
- AC-2: `SpriteOffset` is a new struct in `internal/engine/data/schemas/` with fields `X int` (`json:"x"`) and `Y int` (`json:"y"`).
- AC-3: The actor sprite renderer applies `RenderOffset` (in pixels) as an additive draw-position nudge when the field is non-nil for the current state.
- AC-4: When `RenderOffset` is nil the draw position is identical to the current behavior (zero regression).
- AC-5: `RenderOffset` has no effect on `Body` position, velocity, `CollisionRects`, or `FootprintRect` — physics and collision are entirely unaffected.
- AC-6: The offset is applied after all existing draw transformations (facing-direction flip, body-rect anchoring); it is the final screen-space nudge.
- AC-7: Layer rules upheld: `internal/engine/data/schemas/` does not import `internal/kit/` or `internal/game/`; the offset is consumed only inside the kit-level sprite renderer.
- AC-8: No existing actor JSON files require modification; all omit `render_offset` and continue to render identically.
- AC-9: Table-driven unit tests cover: state with `render_offset {x:-4, y:0}` produces a draw call shifted by (-4, 0) relative to the base position; state without `render_offset` produces the same draw position as baseline; facing-left mirroring does not invert the offset unless the Spec Engineer explicitly decides it should (decision must be documented in SPEC.md).

---

## Behavioral Edge Cases

- Facing-left mirroring interaction: whether `X` is mirrored when the actor faces left must be an explicit Spec Engineer decision; story leaves it open.
- Large offsets (e.g., `x: 200`) are allowed by the schema; the render system does not clamp; art authors are responsible for sensible values.
- State transition mid-frame: the offset is re-read from the current state's `AssetData` on every draw call; no stale value from the previous state.
- `RenderOffset {x:0, y:0}` is functionally identical to nil; both must produce zero-delta draw positions.
