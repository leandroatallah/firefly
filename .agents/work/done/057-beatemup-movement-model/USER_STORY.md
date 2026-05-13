# USER STORY — 057-beatemup-movement-model

**Branch:** `057-beatemup-movement-model`
**Bounded Context:** Physics (`internal/engine/physics/movement/`)

---

## Story

As an engine developer,
I want a `BeatEmUpMovementModel` that provides 8-directional physics with no Y-axis gravity, obstacle collision, and tilemap-edge clamping,
so that beat-em-up actors move freely on the ground plane without falling, and are stopped by level-designer-placed obstacle tiles.

---

## Background

`PlatformMovementModel` applies gravity; no existing model suits beat-em-up ground-plane movement. `TopDownMovementModel` is the closest analog (diagonal normalization, speed cap, friction, `ApplyValidPosition`). This model is a stripped version: passive input (acceleration set externally by `EightDirectionalMovementSkill`), no Y gravity, no jump, no `onGround` flag, no `maxFallSpeed`.

**Depends on:** Independent of story 056; may be implemented in parallel.

---

## Constraints (resolved in grilling session 2026-05-09)

- **Y axis is ground-plane depth (`y16`)** — NOT altitude. "Walking up" the playfield moves the actor deeper on the ground plane. No `z16` field is introduced; vocabulary in comments uses "groundY" / "ground-plane depth".
- **No playfield boundary args at construction.** `minY`/`maxY` constructor params are eliminated. Walkable-strip bounds are enforced by Tiled obstacle tiles placed by level designers; the existing collision system handles vertical blocking via `ApplyValidPosition(vy16, false, space)`.
- **Existing `clampToPlayArea` reused as-is** — it clamps to engine/tilemap edges only, not to a configurable playfield strip. No new camera or bounds infrastructure needed here.
- **"No gravity" means no Y-axis gravity.** Altitude-axis gravity (for future jump skill) must not be precluded — leave a clean integration point (do not hardcode altitude=0 forever). Spec engineer decides how.
- **Passive model** — no embedded `InputHandler`. Acceleration is set by the skill before `Update` runs.

---

## Acceptance Criteria

- AC-1: `BeatEmUpMovementModel` lives in `internal/engine/physics/movement/movement_model_beatemup.go` and satisfies the same `MovementModel` contract used by `TopDownMovementModel` and `PlatformMovementModel`.
- AC-2: No Y-axis gravity is applied in `Update` — the body does not accumulate downward Y velocity when idle.
- AC-3: Diagonal movement is normalized via the existing `smoothDiagonalMovement` helper (or equivalent) so diagonal speed equals single-axis speed.
- AC-4: Velocity magnitude is capped using the same speed-cap logic as `TopDownMovementModel` (respecting `body.MaxSpeed()` and config `SpeedMultiplier`).
- AC-5: Friction is applied via `reduceVelocity` on both axes each frame.
- AC-6: `ApplyValidPosition` is called for both X and Y axes against the `BodiesSpace`; obstacle tiles placed in the tilemap block movement on both axes.
- AC-7: `clampToPlayArea(body, space)` is called (or equivalent) to clamp the body to tilemap/engine edges — no additional bound arguments needed.
- AC-8: `body.Freeze() == true` causes `Update` to return early (same pattern as `TopDownMovementModel`).
- AC-9: Table-driven unit tests cover: no Y gravity when idle, diagonal speed normalization, X obstacle collision respected, Y obstacle collision respected, friction applied after movement, freeze guard.

---

## Behavioral Edge Cases

- Zero velocity with no input: friction pass is a no-op (velocity stays zero, no oscillation).
- Body with `Freeze() == true`: `Update` returns early.
- Future altitude-gravity integration point: model must not write a hardcoded altitude value; altitude field is left untouched (spec engineer to prescribe exact approach).
