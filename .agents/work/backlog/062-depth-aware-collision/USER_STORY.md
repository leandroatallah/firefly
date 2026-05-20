# USER STORY — 062-depth-aware-collision

**Branch:** `062-depth-aware-collision`
**Bounded Context:** Physics (`internal/engine/physics/space/`)

---

## Story

As an engine developer,
I want `HasCollision` in `internal/engine/physics/space/space.go` to gate 2.5D collisions on a depth-lane check,
so that airborne or depth-separated entities do not register collisions with bodies they cannot physically reach.

---

## Background

Current `HasCollision` checks only 2D bounding-box (X, screen-Y) overlap. In a 2.5D beat-em-up the screen-Y of an airborne entity is `Y - Altitude`, so two actors at different ground depths can appear to overlap on screen without being in the same lane. A collision should only fire when both the 2D bbox overlap AND `abs(a.GroundY - b.GroundY) <= LaneWidth`.

`GroundY` is the body's Y coordinate (depth), available via `GetPosition16()` (second return). `LaneWidth` is a new constant or config value; the Spec Engineer decides the source (config vs. hardcoded constant). The check must be opt-in: bodies that do not expose a depth-lane tag (i.e., plain 2D bodies) must use the existing path with no regression.

**Depends on:** 061 (altitude grounding) is logically prior but not a hard compile dependency. These can be implemented concurrently; 062 should be sequenced after 061 in the roadmap.

---

## Acceptance Criteria

- AC-1: `HasCollision` returns `false` when the 2D bboxes overlap but `abs(a.GroundY - b.GroundY) > LaneWidth` for bodies that opt into depth-lane checking.
- AC-2: `HasCollision` returns `true` when 2D bboxes overlap AND `abs(a.GroundY - b.GroundY) <= LaneWidth`.
- AC-3: Bodies that do not opt into depth-lane checking use the existing overlap-only path — no regression for 2D platformer bodies.
- AC-4: `LaneWidth` is configurable (config struct or named constant); its value must not be hardcoded as a magic number inline.
- AC-5: The opt-in mechanism (interface assertion, tag field, or other) is decided by the Spec Engineer and documented in SPEC.md; this story does not prescribe the mechanism.
- AC-6: Layer rules upheld: `internal/engine/physics/space/` does not import `internal/kit/` or `internal/game/`.
- AC-7: Table-driven unit tests cover: same-lane overlap collides, different-lane overlap does not collide, 2D-only bodies (no depth opt-in) still collide on bbox overlap, same-lane no-bbox-overlap does not collide.

---

## Behavioral Edge Cases

- Two airborne entities at identical `GroundY` but different altitudes: bbox check governs (they share a lane; whether screen-Y bboxes overlap is the determining factor).
- `LaneWidth` of 0: only entities at exactly the same `GroundY` collide — valid edge case, must not panic.
- Mixed scene (2D platformer bodies + 2.5D bodies): each pair evaluated independently via its own opt-in status.
