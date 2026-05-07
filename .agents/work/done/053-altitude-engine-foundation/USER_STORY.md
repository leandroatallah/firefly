# User Story 053 — Altitude Engine Foundation

**Branch:** `053-altitude-engine-foundation`
**Bounded Context:** Physics (`internal/engine/physics/`)

---

## Story

As an engine developer, I want the `Body` and `Movable` contracts to carry an `Altitude` axis and the physics `Body` implementation to map it to screen coordinates via `ScreenY = Y - Altitude`, so that future beat-em-up actors can exist in a 2.5D space while all existing 2D levels remain unaffected.

---

## Background

The engine currently models position as a flat 2D plane (X, Y). To support beat-em-up (2.5D) gameplay, a third logical axis — `Altitude` — must be introduced. In 2.5D:

- **Y** represents ground depth (the "lane" on screen).
- **Altitude** represents how high above the ground floor the entity is (jump height).
- **ScreenY** = `Y - Altitude`, so a jumping entity visually rises while its depth (Y) stays fixed.

Existing 2D levels are fully compatible: `Altitude` defaults to `0`, making `ScreenY = Y` exactly as before.

A new kit-layer actor package for beat-em-up characters (`internal/kit/actors/beatemup/`) is introduced in this story as an empty scaffolded package — it registers the new axis but contains no movement logic (that is Phase 2).

---

## Acceptance Criteria

### AC-1: Body contract exposes Altitude

Given the `Body` interface in `internal/engine/contracts/body/body.go`, when a consumer calls altitude accessors, then the interface must expose:

- `Altitude() int` — current altitude in pixels.
- `SetAltitude(alt int)` — sets altitude in pixels.
- `Altitude16() int` — altitude as fixed-point fp16.
- `SetAltitude16(alt16 int)` — sets altitude as fp16.

### AC-2: Movable contract exposes Altitude velocity and acceleration

Given the `Movable` interface, when a consumer calls altitude dynamics accessors, then the interface must expose:

- `VAltitude16() int` — vertical velocity on the altitude axis (fp16).
- `SetVAltitude16(v16 int)` — sets that velocity.
- `AccelerationAltitude() int` — acceleration on the altitude axis.
- `SetAccelerationAltitude(acc int)` — sets that acceleration.

### AC-3: Physics Body implementation stores altitude16

Given `internal/engine/physics/body/body.go`, when the struct is instantiated, then:

- An `altitude16 int` field is present, zero-initialised.
- `Altitude()` returns `fp16.To16(b.altitude16)` (pixel value).
- `Altitude16()` returns `b.altitude16` directly.
- `SetAltitude(alt int)` stores `fp16.From16(alt)`.
- `SetAltitude16(alt16 int)` stores `alt16` directly.

### AC-4: Position() maps Altitude to screen Y

Given a `Body` whose `Y` ground position is `groundY` pixels and whose `Altitude` is `alt` pixels, when `Position()` is called, then `Position().Min.Y` must equal `groundY - alt`.

Concretely: if `groundY = 200` and `alt = 50`, then `Position().Min.Y == 150`.

### AC-5: Zero Altitude preserves existing 2D behaviour

Given a `Body` with `altitude16 == 0` (the default), when `Position()` is called, then `Position().Min.Y` equals the ground Y exactly as before this change — no regression to existing 2D tests.

### AC-6: MovableBody implementation satisfies new Movable methods

Given `internal/engine/physics/body/body_movable.go`, when the struct is used, then:

- `vAltitude16 int` field is present, zero-initialised.
- `VAltitude16()` / `SetVAltitude16()` read and write that field.
- `AccelerationAltitude()` / `SetAccelerationAltitude()` read and write an `accAltitude int` field.

### AC-7: Z-sort (depth sort) by Y in scene rendering

Given a scene that draws multiple entities, when entities have different `Y` (ground depth) values, then entities are drawn in ascending Y order (lower Y drawn first, higher Y drawn on top) so that depth layering appears correct. This sort must be stable to avoid draw-order flicker between entities sharing the same Y.

The sort must apply to the existing scene rendering path and must not break any existing 2D scene tests.

### AC-8: Beat-em-up actor package scaffold exists

Given `internal/kit/actors/beatemup/`, when the package is compiled, then:

- A `doc.go` with package comment exists.
- A placeholder `beatemup_character.go` file exists with a stub `BeatEmUpCharacter` struct embedding the existing kit `Actor` type (or referencing it by interface) — no movement logic is implemented.
- The package compiles cleanly with `go build ./internal/kit/actors/beatemup/...`.

### AC-9: All existing tests continue to pass

Given the full test suite, when `go test ./...` is run after this change, then there are no regressions. Coverage delta on `internal/engine/physics/body/` must be non-negative.

---

## Behavioural Edge Cases

| Scenario | Expected Behaviour |
|---|---|
| `SetAltitude` called with negative value | `Altitude()` returns the negative value; clamping is Phase 2's responsibility |
| `Position()` called when altitude > groundY | `ScreenY` is negative; clipping to viewport is the renderer's concern, not the Body |
| Entities with equal Y during Z-sort | Draw order is stable (insertion order preserved within the same Y bucket) |
| 2D level entity with `Altitude == 0` | `Position()` is bit-identical to current implementation |

---

## Out of Scope (future phases)

- Gravity and jump physics (`vAltitude16` integration, ground detection) — Phase 2.
- 8-way floor movement model — Phase 2.
- Depth-aware collision (`LaneWidth` threshold) — Phase 3.
- Shadow component rendering — Phase 3.
- Tiled level content for the beat-em-up demo — manual, no story needed.

---

## Notes on Platformer Package Relocation

`internal/engine/entity/actors/platformer/` currently lives in the `entity` layer of the engine, not in `kit`. Architecturally, a platformer character is a genre-reusable concrete actor and therefore belongs in `internal/kit/actors/platformer/` alongside the new `beatemup/` package.

**Recommendation:** Create a follow-up story (054) to migrate `internal/engine/entity/actors/platformer/` into `internal/kit/actors/platformer/`. This is kept out of story 053 to keep Phase 1 small and to avoid breaking the existing platformer integration tests during the altitude refactor. The migration is a pure relocation (no behaviour change) and should be straightforward once Phase 1 is stable.
