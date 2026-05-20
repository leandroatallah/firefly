# USER STORY — 064-beatemup-footprint-rect

**Branch:** `064-beatemup-footprint-rect`
**Bounded Context:** Kit (`internal/kit/actors/beatemup/`, `internal/engine/data/schemas/`)

---

## Story

As a beat-em-up game developer,
I want actors to declare a per-state `footprint_rect` in their JSON files,
so that movement and world-collision checks use the actor's feet area rather than the full sprite body.

---

## Background

Classical beat-em-ups use a fake-Z ground plane where screen-Y encodes depth. Movement collisions (walls, props, actor-vs-actor on the floor) must be gated by the actor's feet footprint, not the full sprite bounds. Today all collision uses `collision_rect` (the full body), causing:

- An actor stops when their *head* reaches a wall that sits at the top of the screen instead of when their *feet* reach its base.
- Actors at different screen-Y rows (different ground depths) block each other even though they occupy different lanes.

**Related stories:**
- Story 062 (`062-depth-aware-collision`) adds the depth-lane gate inside `internal/engine/physics/space/HasCollision`. That story provides the *gate mechanism*; this story provides the *footprint shape* that beat-em-up actors feed into collision checks. They are independent and can ship in either order.

---

## Acceptance Criteria

- AC-1: `schemas.AssetData` gains an optional `FootprintRect *ShapeRect` field tagged `json:"footprint_rect,omitempty"`; absent field parses without error and leaves the pointer nil.
- AC-2: `internal/kit/actors/beatemup` character exposes a `Footprint() image.Rectangle` method that returns the footprint of the current state's sprite asset, offset by the actor's current world position.
- AC-3: When `footprint_rect` is absent for the current state, `Footprint()` falls back to the existing full collision rect (no regression).
- AC-4: Beat-em-up movement-vs-world collision and actor-vs-actor movement collision use `Footprint()` instead of the full body rect.
- AC-5: Attack hitbox resolution is unchanged — it continues to use the full `collision_rect`.
- AC-6: Platformer and topdown actors are unaffected; `FootprintRect` is only consumed inside `internal/kit/actors/beatemup/`.
- AC-7: Layer rules upheld: `internal/engine/data/schemas/` does not import `internal/kit/` or `internal/game/`; `internal/kit/` does not import `internal/game/`.
- AC-8: Table-driven unit tests cover: state with `footprint_rect` returns correct world-offset rectangle; state without `footprint_rect` falls back to collision rect; facing-left mirroring (if applicable) is consistent with existing collision rect behavior.

---

## Behavioral Edge Cases

- State transition mid-frame: `Footprint()` always reflects the current state at call time; no stale rect from previous state.
- `footprint_rect` with zero width or height: treated as an empty rect; movement collision skips that axis or treats as no-footprint (Spec Engineer decides, must document).
- Jump/fall states may legitimately omit `footprint_rect` to indicate the actor is airborne; `Footprint()` falls back to full collision rect in that case per AC-3 until airborne handling is introduced.
- Actor JSON with multiple collision rects but a single footprint rect: footprint is always a single rect; `CollisionRects` slice is unchanged.
