# USER STORY — 069-depth-lane-body-impl

**Branch:** `069-depth-lane-body-impl`
**Bounded Context:** Engine (`internal/engine/physics/body/`, `internal/engine/physics/space/`, `internal/engine/physics/movement/`)
**Epic:** `beatemup-mechanics`

---

## Story

As a beat-em-up game developer,
I want background obstacles and beat-em-up characters to implement `DepthLaneBody`,
so that an airborne player no longer false-collides with background walls placed at a different depth lane.

---

## Background

Story 062 wired the `DepthLaneBody` gate into `space.HasCollision` (`internal/engine/physics/space/space.go:274`), but **no production body implements the interface**. The gate therefore never fires, and `HasCollision` falls through to pure 2D bbox overlap for every pair.

Combined with the per-frame shape-shift in `movement_model_beatemup.go:120-126` (which lifts airborne shapes into screen space so combat/items work), this means airborne shapes overlap any background obstacle that sits at a higher screen-Y → false collision the moment the player jumps.

`BeatEmUpMovementModel.Update` currently works around the missing gate with a zero-altitude wrap around `ApplyValidPosition` (Block 1, lines 47-55). It only protects player→world checks initiated by the player. Other-body→player checks (enemy AI, projectiles iterating `ResolveCollisions(self)`) bypass the wrap and still false-collide.

Implementing `DepthLaneBody` on the two body types that participate in beat-em-up collisions closes the gap. Once the gate is reliable, the Block 1 zero-out hack can be removed.

---

## Acceptance Criteria

- AC-1: `*body.ObstacleRect` implements `space.DepthLaneBody` with `GroundY()` returning the obstacle's world-Y (pre-altitude) bottom edge, and `LaneHalfWidth()` returning a value derived from its world-Y height so the lane covers the obstacle's full depth extent.
- AC-2: `*beatemupkit.BeatEmUpCharacter` implements `space.DepthLaneBody` with `GroundY()` returning the body's ground Y (pre-altitude, i.e. `y16 >> 16`) and `LaneHalfWidth()` returning `space.DefaultLaneHalfWidth` (8px).
- AC-3: `space.HasCollision(player, obstacle)` returns `false` when their ground-Y depth difference exceeds the larger of the two `LaneHalfWidth` values, regardless of bbox overlap on screen.
- AC-4: `space.HasCollision(player, obstacle)` returns `true` when bboxes overlap AND ground-Y depth difference is within tolerance, regardless of player altitude.
- AC-5: `*PlatformerCharacter` (and any other non-beat-em-up character) does NOT implement `DepthLaneBody`. 2D platformer scenes retain pure bbox behavior (regression-safe).
- AC-6: The zero-altitude wrap in `movement_model_beatemup.go:47-55` (Block 1) is removed. Wall blocking now relies solely on the depth-lane gate.
- AC-7: Block 2 (shape shift on altitude change, `movement_model_beatemup.go:120-126`) remains in place — combat/item hitboxes still need altitude-aware shapes.
- AC-8: Table-driven tests in `internal/engine/physics/space/` cover: same-depth + bbox overlap → collision; different-depth + bbox overlap → no collision; airborne player + same-depth wall → collision (block); airborne player + different-depth wall → no collision (no block).
- AC-9: Existing 2D platformer regression tests pass without modification.
- AC-10: Manual verification: jumping near a background wall placed at a different depth no longer triggers a false block; jumping into a wall at the player's depth still blocks (when shapes overlap on the depth axis).

---

## Behavioral Edge Cases

- **Obstacle with zero height**: `LaneHalfWidth()` must fall back to `DefaultLaneHalfWidth` to avoid an exact-equal match requirement.
- **Player ground-Y exactly at obstacle ground-Y ± tolerance boundary**: inclusive match (`<=` in the gate, already correct in `space.go:293`).
- **Owner indirection**: `ResolveCollisions` passes parent bodies (CodyPlayer, ObstacleRect) to `HasCollision`, not child collision shapes. Interface lives on the parent — confirm `CodyPlayer` inherits via embedded `BeatEmUpCharacter` without re-declaration.
- **Airborne obstacles** (none today, but future-proof): `GroundY()` must remain pre-altitude even if an obstacle ever sets altitude.

---

## Out of Scope

- Story 064 (`beatemup-footprint-rect`) — independent; footprint JSON schema is orthogonal to the depth-lane gate.
- A depth-lane debug visualizer.
- Variable per-state lane widths on the character.

---

## Migration Note

After this story lands, update `.agents/work/epics/beatemup-mechanics/PLAN_airborne-collision-split.md` to reflect: Option B chosen, Block 1 removed, Option C retired.
