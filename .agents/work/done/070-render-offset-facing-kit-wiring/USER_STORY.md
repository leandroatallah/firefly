# USER STORY — 070-render-offset-facing-kit-wiring

**Branch:** `070-render-offset-facing-kit-wiring`
**Bounded Context:** Kit (`internal/kit/actors/`, `internal/engine/data/schemas/`, `internal/engine/entity/actors/builder/`)

---

## Story

As a beat-em-up game developer,
I want per-facing-direction render offsets in actor JSON and `ApplyRenderOffsets` wired into the platformer kit constructor,
so that asymmetric sprite art can be nudged correctly for each facing direction, and platformer actors receive the same offset support without silent no-ops.

---

## Background

Story 068 shipped `render_offset {x, y}` in actor JSON. Two deferred items are now in scope (see `068-actor-json-sprite-render-offset/NOTES.md`):

1. The single `x` value is applied identically regardless of facing direction. `cody-melee-0.png`-style art where the fist extends to one side needs a different X nudge when facing left vs right. Adding an optional `x_flipped` field to `SpriteOffset` enables this: when the actor faces left, use `x_flipped` (if set); otherwise fall back to `x`. `y` is unchanged for both facings.

2. `ApplyRenderOffsets` is only called from the beatemup kit constructor. Platformer actor JSONs that declare `render_offset` silently no-op today. The platformer constructor (`internal/kit/actors/platformer/platformer.go`) must call `ApplyRenderOffsets` so the schema is uniformly honoured across kit genres. The shooter kit (`internal/kit/actors/shooter_character.go`) is a trait helper, not a constructor with `SpriteData`; it does not need changes.

---

## Acceptance Criteria

- AC-1: `SpriteOffset` in `internal/engine/data/schemas/json.go` gains an optional `XFlipped *int` field tagged `json:"x_flipped,omitempty"`; absent field parses without error and leaves the pointer nil.
- AC-2: When `XFlipped` is nil, `SpriteOffset` is backward-compatible — the resolved X for all facings equals the existing `X` field (zero regression for all current JSONs).
- AC-3: `Character.UpdateImageOptions()` in `internal/engine/entity/actors/character.go` uses `x_flipped` (when non-nil) as the X translation when `fDirection == FaceDirectionLeft`; uses `X` in all other cases. `Y` is unchanged regardless of facing.
- AC-4: `Character.SetRenderOffset` is extended to accept both `dx` (right-facing) and `dxFlipped *int` (optional left-facing override); or an alternative mechanism documented in `SPEC.md` that achieves the same runtime lookup without breaking existing callers.
- AC-5: `builder.ApplyRenderOffsets` propagates the `XFlipped` pointer from `SpriteOffset` into the character's render-offset store (or equivalent mechanism from AC-4).
- AC-6: `NewPlatformerCharacter` in `internal/kit/actors/platformer/platformer.go` calls `builder.ApplyRenderOffsets(c, spriteData, stateMap)` at the same point the beatemup constructor does — after `actors.NewCharacter` and before `return`.
- AC-7: Existing platformer actor JSON files that omit `render_offset` continue to render identically (zero regression).
- AC-8: Layer rules upheld: `internal/engine/data/schemas/` does not import `internal/kit/` or `internal/game/`; the offset resolution logic lives in `internal/engine/entity/actors/character.go`.
- AC-9: Table-driven unit tests cover: facing-right uses `X` even when `XFlipped` is set; facing-left uses `XFlipped` when set; facing-left falls back to `X` when `XFlipped` is nil; `Y` is identical for both facings.
- AC-10: Table-driven unit tests for `builder.ApplyRenderOffsets` cover: asset with `XFlipped` set registers the flipped value correctly; asset without `XFlipped` registers with nil flipped (falls back to `X` at draw time).
- AC-11: A unit test for `NewPlatformerCharacter` verifies that a `SpriteData` with `render_offset {x:-2}` results in `RenderOffset(Idle)` returning `ok=true` on the returned `PlatformerCharacter`.

---

## Behavioral Edge Cases

- `XFlipped` set to `0` is a valid explicit override (zero-shift when facing left); must not be treated the same as absent.
- Facing direction is determined at draw time (`UpdateImageOptions`) from `fDirection` (acceleration + `SetFaceDirection`); the offset must re-resolve on every call — no caching between frames.
- Assets that declare only `x` (no `x_flipped`) must behave exactly as story 068: same X for both facings, no regression.
- Non-platformer genres (shooter trait, melee archetypes) are unaffected by AC-6; only the platformer kit constructor is wired.
- Large or negative `XFlipped` values are not clamped; art authors are responsible for sensible values.
- State transition mid-frame: the offset is re-read each `UpdateImageOptions` call; no stale value from the previous state.
