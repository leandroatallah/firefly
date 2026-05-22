# USER STORY — 063-shadow-component

**Branch:** `063-shadow-component`
**Bounded Context:** Kit (`internal/kit/`)

---

## Story

As a game developer,
I want airborne beat-em-up entities to render a simple oval shadow at their ground position `(X, GroundY)`,
so that players have a clear visual anchor for where an entity will land.

---

## Background

When an entity is airborne (`Altitude > 0`) its sprite is drawn at `ScreenY = Y - Altitude`, which lifts it visually off the floor. Without a shadow the player cannot judge landing position. A `Shadow` component renders a fixed-size semi-transparent oval at `(X, Y)` (the ground-plane Y, not the screen-projected Y) before the actor sprite is drawn each frame.

The shadow should be drawn in the beat-em-up phase scene's draw loop, behind all entity sprites. The Spec Engineer decides whether `Shadow` is a standalone `Drawable` registered in the scene or a method on the actor. The component lives in `internal/kit/` and must not import `internal/game/`.

**Depends on:** 061 (altitude grounding) — shadow is only meaningful once altitude is live. 062 has no dependency relationship with this story.

---

## Acceptance Criteria

- AC-1: A `Shadow` type (or equivalent) exists in `internal/kit/` and renders a semi-transparent oval at the body's ground position `(X, GroundY)`.
- AC-2: The shadow is drawn each frame when `body.Altitude() > 0`; when `Altitude == 0` no shadow is drawn (entity is on the ground).
- AC-3: The oval is scaled relative to altitude: it shrinks as altitude increases, providing a depth cue (exact scale formula is Spec Engineer's decision).
- AC-4: Shadow alpha is fixed (e.g. 50% opacity); the value is a named constant, not a magic number.
- AC-5: The shadow is drawn behind all entity sprites in the beat-em-up phase scene draw order.
- AC-6: Layer rules upheld: shadow implementation does not import `internal/game/`.
- AC-7: Existing platformer and 2D scenes are unaffected (shadow is only registered in the beat-em-up scene).
- AC-8: Unit tests (headless, no GPU) cover: shadow drawn when airborne, no shadow at zero altitude, oval bounds computed correctly for a given altitude.

---

## Behavioral Edge Cases

- `Altitude` very large (entity launched high): shadow shrinks to minimum size rather than disappearing entirely — minimum size is a named constant.
- `Altitude` exactly 0 on the same frame entity lands: shadow is not drawn (landing frame is ground-level).
- Scene with no entities airborne: shadow draw pass is a no-op with no allocations.
- Shadow drawn with `ebiten.NewImage` in tests: no GPU calls; use headless image creation pattern from constitution.
