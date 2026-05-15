# USER STORY — 058-wire-beatemup-movement

**Branch:** `058-wire-beatemup-movement`
**Bounded Context:** Kit / Game Logic (`internal/kit/actors/beatemup/`, `internal/game/scenes/phases/beatemup/`)

---

## Story

As a developer assembling the beat-em-up boilerplate,
I want `BeatEmUpCharacter` to use `EightDirectionalMovementSkill` and `BeatEmUpMovementModel`, with camera bounds wired to the tilemap rectangle,
so that the player walks in 8 directions, is blocked by obstacle tiles, and the camera follows within the arena.

---

## Background

Story 055 delivered the beat-em-up phase scene shell (`internal/game/scenes/phases/beatemup/scene.go`). Stories 056 and 057 deliver the skill and model respectively. This story wires them together and extends `kitskills.FromConfig` with a `Mode` discriminator.

**Depends on:** 056 and 057 must be complete.

---

## Constraints (resolved in grilling session 2026-05-09)

- **No playfield bounds injection.** `minY`/`maxY` are eliminated. Walkable-strip bounds come from Tiled obstacle tiles; the scene does not compute or pass bound args to the model.
- **Camera follow** is bounded to the full tilemap rectangle via `Camera().SetBounds(tilemapRect)` and `SetFollowTarget(player)`. No new camera infrastructure; both calls already exist in the engine.
- **Skill registration via `cfg.Movement.Mode`**: extend `kitskills.FromConfig` in `internal/kit/skills/factory.go` with a `Mode` discriminator on `schemas.SkillsConfig.Movement`:
  - `mode: "horizontal"` (default, backward-compatible) → `HorizontalMovementSkill`
  - `mode: "eight_dir"` → `EightDirectionalMovementSkill`
  - Per-actor config selects its own mode. Factory remains genre-agnostic.
- **Model ownership follows platformer precedent**: `BeatEmUpCharacter` owns its `BeatEmUpMovementModel`, created in `NewBeatEmUpCharacter` and accessible via `MovementModel()`. The scene does not inject the model.
- **Cross-phase state preservation** (HP/inventory across actor swaps when genre-switching) is out of scope; flag for a future story.

---

## Acceptance Criteria

- AC-1: `kitskills.FromConfig` accepts `mode: "eight_dir"` and instantiates `EightDirectionalMovementSkill`; `mode: "horizontal"` (or absent) is backward-compatible and still produces `HorizontalMovementSkill`.
- AC-2: `BeatEmUpCharacter` owns a `BeatEmUpMovementModel` created at construction; `MovementModel()` returns it. `PlatformMovementModel` is not referenced in the beatemup actor package.
- AC-3: `BeatEmUpCharacter` registers `EightDirectionalMovementSkill` so `HandleInput` is called each frame during actor update.
- AC-4: The beat-em-up phase scene calls `Camera().SetBounds(tilemapRect)` and `Camera().SetFollowTarget(player)` in `OnStart`; camera stays within the tilemap boundary.
- AC-5: Tilemap collision bodies are created via `s.Tilemap().CreateCollisionBodies(s.PhysicsSpace(), nil)`; obstacle tiles block X and Y movement via the model's `ApplyValidPosition` calls.
- AC-6: In runtime, the player moves in all 8 directions, is stopped by obstacle tiles on both axes, and does not accumulate downward velocity when idle.
- AC-7: Existing platformer phase scene and its tests are unaffected.
- AC-8: Layer rules upheld: `internal/kit/actors/beatemup/` does not import `internal/game/`.
- AC-9: Unit tests for `BeatEmUpCharacter` construction: skill present, model present, no panic on zero-input update frame. Unit test for `kitskills.FromConfig`: `mode: "eight_dir"` returns correct type.

---

## Behavioral Edge Cases

- Immobile flag set during cutscene: `EightDirectionalMovementSkill` zeroes velocity — no drift while immobile.
- Tilemap with no collision bodies: player moves freely on both axes; no panic.
- Camera at tilemap edge: `SetBounds` prevents the camera from scrolling past the tilemap boundary on any side.
- `mode` field absent from config: treated as `"horizontal"` — no regression for existing actors.

---

## Out of Scope

- Cross-phase HP/inventory preservation across genre actor swaps — future story.
- PlayArea object layer for configurable walkable strip — revisit if obstacle-tile approach proves insufficient.
