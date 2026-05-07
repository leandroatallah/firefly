# Technical Specification — 053 Altitude Engine Foundation

**Branch:** `053-altitude-engine-foundation`
**Bounded Context:** Physics (`internal/engine/physics/`), Scene rendering (`internal/engine/render/draworder/` new), Kit (`internal/kit/actors/beatemup/` new)
**Pipeline Phase:** Spec (Architect → **Spec** → Mock Generator → TDD → Feature → Gatekeeper)

---

## 1. Goals & Non-Goals

### Goals

- Extend the `Body` and `Movable` engine contracts with an `Altitude` axis using fp16 fixed-point storage.
- Implement `altitude16` storage and accessors in `physics/body.Body` and `vAltitude16` / `accAltitude` in `physics/body.MovableBody`.
- Map `Altitude` to screen Y inside `Body.Position()` so that `Position().Min.Y == groundY - altitude`.
- Add an engine-layer Z-sort helper that returns the scene's collidables in stable ascending Y order, and wire it into the existing scene render path (`internal/game/scenes/phases/scene.go`).
- Scaffold an empty `internal/kit/actors/beatemup/` package containing `doc.go` and `beatemup_character.go`.

### Non-Goals

- No gravity, jump, or ground-detection logic on the altitude axis (Phase 2).
- No 8-way floor movement, depth-aware collision, shadow rendering, or beat-em-up demo content (Phase 2/3).
- No relocation of the existing `internal/engine/entity/actors/platformer/` package (story 054).
- No new `Y`-based ordering inside `BodiesSpace.Bodies()` itself — keep the existing ID sort to avoid breaking collision determinism. Y-sort is a render-only concern.

---

## 2. File-by-File Changes

### 2.1 `internal/engine/contracts/body/body.go`

Add four methods to the `Body` interface:

```go
// Altitude returns the body's altitude above the ground in pixels.
Altitude() int
// SetAltitude sets the body's altitude in pixels.
SetAltitude(alt int)
// Altitude16 returns the body's altitude as a fixed-point fp16 value.
Altitude16() int
// SetAltitude16 sets the body's altitude as a fixed-point fp16 value.
SetAltitude16(alt16 int)
```

Add four methods to the `Movable` interface (right after the existing `Acceleration`/`SetAcceleration` block):

```go
// VAltitude16 returns the velocity on the altitude axis (fp16).
VAltitude16() int
// SetVAltitude16 sets the velocity on the altitude axis (fp16).
SetVAltitude16(v16 int)
// AccelerationAltitude returns the acceleration on the altitude axis.
AccelerationAltitude() int
// SetAccelerationAltitude sets the acceleration on the altitude axis.
SetAccelerationAltitude(acc int)
```

### 2.2 `internal/engine/physics/body/body.go`

- Add field: `altitude16 int` to the `Body` struct (zero-valued by default; no constructor change).
- Add accessor methods:
  ```go
  func (b *Body) Altitude() int        { return fp16.From16(b.altitude16) }
  func (b *Body) Altitude16() int      { return b.altitude16 }
  func (b *Body) SetAltitude(alt int)  { b.altitude16 = fp16.To16(alt) }
  func (b *Body) SetAltitude16(a int)  { b.altitude16 = a }
  ```
- **Modify `Position()`**: subtract altitude from the screen Y. Final form:
  ```go
  func (b *Body) Position() image.Rectangle {
      minX := fp16.From16(b.x16)
      groundY := fp16.From16(b.y16)
      alt := fp16.From16(b.altitude16) // 0 by default → no behaviour change
      minY := groundY - alt
      maxX := minX + b.shape.Width()
      maxY := minY + b.shape.Height()
      return image.Rect(minX, minY, maxX, maxY)
  }
  ```
  Equivalence: when `altitude16 == 0`, `alt == 0`, so `minY == groundY` — bit-identical to the prior implementation (AC-5).

### 2.3 `internal/engine/physics/body/body_movable.go`

- Add fields to `MovableBody`:
  ```go
  vAltitude16  int
  accAltitude  int
  ```
- Add accessors (place beside existing `Velocity`/`Acceleration`):
  ```go
  func (b *MovableBody) VAltitude16() int             { return b.vAltitude16 }
  func (b *MovableBody) SetVAltitude16(v16 int)       { b.vAltitude16 = v16 }
  func (b *MovableBody) AccelerationAltitude() int    { return b.accAltitude }
  func (b *MovableBody) SetAccelerationAltitude(a int){ b.accAltitude = a }
  ```
- No change to existing `Velocity`, `SetVelocity`, `Acceleration`, `SetAcceleration`, or any movement helper. Altitude motion integration is Phase 2.

### 2.4 `internal/engine/render/draworder/` (new package)

New file: `internal/engine/render/draworder/draworder.go`

Purpose: Provide a stable Z-sort helper at the engine layer that the scene rendering path (game layer) can call. Engine-only (does not import kit or game).

```go
// Package draworder provides stable depth-sorting helpers for scene rendering.
//
// In a 2.5D world, draw order should follow ascending ground Y (lower Y drawn
// first, higher Y drawn on top) so closer entities visually overlap farther
// ones. The sort is stable to avoid flicker between entities sharing the same Y.
package draworder

import (
    "sort"

    "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
)

// SortByGroundY returns a new slice containing the input bodies sorted in
// ascending ground-Y order. Sort is stable: equal-Y entities preserve their
// input order. The input slice is NOT mutated.
//
// Ground Y is taken from each body's fp16 position via GetPosition16, so the
// rendering order is independent of altitude (entities jumping up still draw
// in the same depth slot).
func SortByGroundY(in []body.Collidable) []body.Collidable {
    out := make([]body.Collidable, len(in))
    copy(out, in)
    sort.SliceStable(out, func(i, j int) bool {
        _, yi16 := out[i].GetPosition16()
        _, yj16 := out[j].GetPosition16()
        return yi16 < yj16
    })
    return out
}
```

Notes:
- We sort by `y16` (fp16 ground Y) rather than `Position().Min.Y` so that altitude does **not** influence draw order — a jumping entity stays in the same lane.
- We return a fresh copy because `BodiesSpace.Bodies()` documents its result as a non-mutable cache.

Test file: `internal/engine/render/draworder/draworder_test.go` will cover the sort cases (see §4 AC-7).

### 2.5 `internal/game/scenes/phases/scene.go`

- Import the new package: `"github.com/boilerplate/ebiten-template/internal/engine/render/draworder"`.
- In `PhasesScene.Draw` (line 472), replace:
  ```go
  for _, b := range space.Bodies() {
  ```
  with:
  ```go
  for _, b := range draworder.SortByGroundY(space.Bodies()) {
  ```
- Do **not** modify any other iteration of `Bodies()` (e.g., `body_counter.go`, `platform_jump.go`, `commands_actor.go`, `scene.go:408`) — those are non-render iterations and must stay deterministic on ID order.

### 2.6 `internal/kit/actors/beatemup/` (new package, scaffold only)

New file: `internal/kit/actors/beatemup/doc.go`

```go
// Package beatemup provides genre-reusable character traits for beat-em-up
// (2.5D) games built on the Firefly engine.
//
// Phase 1 (this package) is a scaffold: it registers the package and a
// placeholder type so the new altitude axis can be exercised by future
// stories (gravity, jump, ground detection) without touching engine layout.
//
// Dependency rule (enforced by CI):
//   - beatemup MAY import internal/engine/...
//   - beatemup MUST NOT import internal/game/...
package beatemup
```

New file: `internal/kit/actors/beatemup/beatemup_character.go`

```go
package beatemup

import kitactors "github.com/boilerplate/ebiten-template/internal/kit/actors"

// BeatEmUpCharacter is a placeholder trait for beat-em-up actors. It embeds
// the existing kit MeleeCharacter so future Phase 2 stories can layer on
// altitude movement, gravity, and ground-detection logic without restructuring
// the type. No behaviour is implemented in Phase 1.
type BeatEmUpCharacter struct {
    *kitactors.MeleeCharacter
}

// NewBeatEmUpCharacter constructs an empty BeatEmUpCharacter scaffold.
func NewBeatEmUpCharacter() *BeatEmUpCharacter {
    return &BeatEmUpCharacter{
        MeleeCharacter: kitactors.NewMeleeCharacter(),
    }
}
```

Rationale for embedding `MeleeCharacter`: beat-em-up characters are conceptually melee-first; we reuse the existing kit trait so future Phase 2 work composes cleanly. The embedding is a structural placeholder — no methods are overridden.

---

## 3. Pre-conditions / Post-conditions

### Pre-conditions
- `internal/engine/utils/fp16` provides `To16` / `From16`.
- `internal/engine/physics/body.Body` is the canonical implementation of `body.Body`.
- `internal/game/scenes/phases/scene.go` is the only render path that iterates over `space.Bodies()` for drawing.

### Post-conditions
- `body.Body` interface has 4 new methods (`Altitude`, `SetAltitude`, `Altitude16`, `SetAltitude16`).
- `body.Movable` interface has 4 new methods (`VAltitude16`, `SetVAltitude16`, `AccelerationAltitude`, `SetAccelerationAltitude`).
- `physics/body.Body` has an `altitude16 int` field; `Position().Min.Y == GetPosition16Y/65536 - altitude16/65536` exactly.
- `physics/body.MovableBody` has `vAltitude16 int` and `accAltitude int` fields.
- `internal/engine/render/draworder` exists and exposes `SortByGroundY([]body.Collidable) []body.Collidable`.
- `PhasesScene.Draw` iterates over `draworder.SortByGroundY(space.Bodies())` for entity draws.
- `internal/kit/actors/beatemup/` compiles cleanly.
- `go test ./...` passes; coverage delta on `internal/engine/physics/body/` is non-negative.

---

## 4. Interface Contract Diff

### `body.Body` — before / after

**Before** (10 methods): `Ownable`, `ID`, `SetID`, `Position`, `SetPosition`, `SetPosition16`, `SetSize`, `Scale`, `SetScale`, `GetPosition16`, `GetPositionMin`, `GetShape`.

**After** — same plus:
- `Altitude() int`
- `SetAltitude(alt int)`
- `Altitude16() int`
- `SetAltitude16(alt16 int)`

### `body.Movable` — before / after

**Before**: `Body` + 30+ movement & state methods (see file).

**After** — same plus:
- `VAltitude16() int`
- `SetVAltitude16(v16 int)`
- `AccelerationAltitude() int`
- `SetAccelerationAltitude(acc int)`

### Mocks impact

Any mock implementing `body.Body` or `body.Movable` must gain the new methods. Likely files (to be confirmed by Mock Generator):
- `internal/engine/mocks/` — search for files that implement `Body`/`Movable`.
- Any package-local `mocks_test.go` that declares a struct conforming to `body.Body` or `body.Movable`.

The new methods on mocks should be trivial getters/setters backed by an `int` field, mirroring the production implementation.

---

## 5. Z-Sort Implementation Notes

- **Location**: engine layer (`internal/engine/render/draworder/`) so kit and game layers can both use it without owning sort logic.
- **Algorithm**: `sort.SliceStable` keyed on `y16` (fp16 ground Y, retrieved via `GetPosition16`). Stability is essential per AC-7 to avoid flicker for equal-Y bodies.
- **Why fp16 key, not pixel Y**: avoids precision loss for sub-pixel positions and keeps the comparison cheap.
- **Why ground Y, not screen Y (`Position().Min.Y`)**: a jumping entity must remain in the same depth slot relative to other entities. Using screen Y would cause a jumping body to "leap forward" in draw order.
- **Why a copy, not in-place**: `BodiesSpace.Bodies()` documents its slice as cache-shared and must not be mutated.
- **Why not modify `BodiesSpace.Bodies()` itself**: that slice is also consumed by collision/physics paths which depend on deterministic ID order. Render is the only consumer that needs Y-sort.

---

## 6. Test Strategy (per Acceptance Criterion)

Tests are written by the **TDD Specialist** (Red Phase) before implementation. Test files and skeletons:

### AC-1 — Body contract exposes Altitude
- Compile-time assertion in `internal/engine/physics/body/body_test.go`:
  ```go
  var _ body.Body = (*Body)(nil)
  ```
  ensures the contract is satisfied. Already present? If not, add it.
- Table-driven test `TestBody_AltitudeAccessors` covering Set/Get round-trips for pixel and fp16 values: `0`, `+50`, `-25`, `1`.

### AC-2 — Movable contract exposes Altitude velocity/acceleration
- Compile-time assertion in `body_movable_test.go`: `var _ body.Movable = (*MovableBody)(nil)`.
- `TestMovableBody_AltitudeDynamics`: round-trip Set/Get for `VAltitude16`, `AccelerationAltitude` with positive, negative, zero, and large values.

### AC-3 — Physics Body stores altitude16
- `TestBody_Altitude16_StoredDirectly`: `SetAltitude16(123456)` then `Altitude16() == 123456`.
- `TestBody_SetAltitude_UsesFp16`: `SetAltitude(50)` then `Altitude16() == fp16.To16(50)` and `Altitude() == 50`.

### AC-4 — Position() maps Altitude to screen Y
- Table-driven `TestBody_Position_AltitudeMapsToScreenY`:
  | groundY | altitude | expected Min.Y |
  |---|---|---|
  | 200 | 50 | 150 |
  | 100 | 0 | 100 |
  | 100 | 100 | 0 |
  | 50 | 75 | -25 (negative permitted; renderer concern) |
  | 200 | -10 | 210 (negative altitude allowed; clamping is Phase 2) |

### AC-5 — Zero Altitude preserves 2D behaviour
- `TestBody_Position_ZeroAltitude_IsBitIdentical`: instantiate body at `(x, y) = (10, 200)` with default altitude; assert `Position() == image.Rect(10, 200, 10+w, 200+h)`. Run before and after invoking altitude setters with `0`.
- Implicit coverage: existing `body_test.go` and other consumer tests must continue to pass unchanged.

### AC-6 — MovableBody satisfies new Movable methods
- Covered by AC-2 tests on the concrete `MovableBody`. Also assert default zero-value: a freshly-constructed `NewMovableBody(...)` returns `0` for `VAltitude16()` and `AccelerationAltitude()`.

### AC-7 — Z-sort stable by Y
- New file `internal/engine/render/draworder/draworder_test.go`:
  - `TestSortByGroundY_AscendingOrder`: input `{y=300, y=100, y=200}` → output `{100, 200, 300}`.
  - `TestSortByGroundY_StableForEqualY`: input `{idA y=100, idB y=100, idC y=50}` → output preserves `{idA, idB}` relative order after `idC`.
  - `TestSortByGroundY_AltitudeIgnored`: two bodies at the same `y16=200` but altitudes `0` and `100` keep their input order (jumping does not change depth).
  - `TestSortByGroundY_DoesNotMutateInput`: snapshot input slice; assert unchanged after call.
  - `TestSortByGroundY_EmptyAndSingle`: empty slice → empty slice; one element → identity.
- Mocks needed: a small `fakeCollidable` in the test file with `GetPosition16` and `ID` (sufficient surface). Use stub or shared `internal/engine/mocks` as appropriate (Mock Generator decision).
- No new test on `PhasesScene.Draw` itself (it would require GPU). The wiring change in `scene.go` is exercised by existing scene tests continuing to pass.

### AC-8 — Beat-em-up scaffold compiles
- New file `internal/kit/actors/beatemup/beatemup_character_test.go`:
  - `TestNewBeatEmUpCharacter_NotNil`: constructor returns non-nil and embedded `MeleeCharacter` is initialized.
- Compile gate: `go build ./internal/kit/actors/beatemup/...` succeeds.

### AC-9 — Full suite green; coverage delta non-negative
- `go test ./...` runs clean.
- Compare `go test -cover ./internal/engine/physics/body/...` before/after; delta ≥ 0 (Workflow Gatekeeper verifies).

---

## 7. Red Phase Scenario (failing test description for TDD Specialist)

The TDD Specialist should land **all** of the tests above as the Red Phase. The single most representative failing test that demonstrates the new behaviour end-to-end is:

```go
// internal/engine/physics/body/body_test.go
func TestBody_Position_AltitudeMapsToScreenY(t *testing.T) {
    rect := body.NewRect(0, 0, 16, 16)
    b := body.NewBody(rect)
    b.SetPosition(0, 200)
    b.SetAltitude(50)

    pos := b.Position()
    if pos.Min.Y != 150 {
        t.Fatalf("Position().Min.Y = %d; want 150 (groundY=200 - altitude=50)", pos.Min.Y)
    }
}
```

This will not compile until `SetAltitude` exists on `*Body`, and will not pass until `Position()` subtracts altitude from ground Y. Together with the contract-conformance assertion (`var _ body.Body = (*Body)(nil)`), it transitively forces the contract additions, the field, and the accessor implementations.

The companion failing test for Z-sort:

```go
// internal/engine/render/draworder/draworder_test.go
func TestSortByGroundY_AscendingStable(t *testing.T) {
    in := []body.Collidable{
        newFakeCollidable("a", 0, 300),
        newFakeCollidable("b", 0, 100),
        newFakeCollidable("c", 0, 100),
    }
    out := draworder.SortByGroundY(in)
    gotIDs := []string{out[0].ID(), out[1].ID(), out[2].ID()}
    want := []string{"b", "c", "a"} // stable: b before c at equal Y
    if !reflect.DeepEqual(gotIDs, want) {
        t.Fatalf("got %v; want %v", gotIDs, want)
    }
}
```

This will not compile until the `draworder` package and `SortByGroundY` exist.

---

## 8. Integration Points

- **Physics ↔ Engine contracts**: the `physics/body` package implements the augmented `body.Body` and `body.Movable` interfaces. No other physics behaviour changes.
- **Render ↔ Game scene**: `PhasesScene.Draw` becomes the first consumer of `draworder.SortByGroundY`. No other scene rendering path changes.
- **Mocks**: any mock implementing `body.Body` or `body.Movable` must be regenerated/updated by the Mock Generator step. Search scope: `internal/engine/mocks/` and any package-local `mocks_test.go`.
- **Kit beatemup ↔ Kit actors**: `beatemup.BeatEmUpCharacter` embeds `kitactors.MeleeCharacter` only; no other kit/engine surfaces are touched.

---

## 9. Architectural Constraints Honored

- **Layer rules**: new code only in `engine` (contracts, physics impl, draworder helper) and `kit` (beatemup scaffold). No new game-layer code apart from the one-line wiring change in `PhasesScene.Draw`. No engine→kit or engine→game imports introduced.
- **Fixed-point**: altitude is stored as `altitude16 int`; pixel accessors use `fp16.From16` / `fp16.To16` exclusively (no float arithmetic).
- **No global mutable state**: all new fields are instance-level on `Body` / `MovableBody`.
- **Tests deterministic**: Z-sort stability is mandatory (AC-7); all new tests are pure-Go, no GPU, no `time.Sleep`.
- **Mocks at boundaries only**: the new `draworder` test uses a tiny `fakeCollidable` implementing the minimal `body.Collidable` surface required.

---

## 10. Open Questions / Notes for Downstream Agents

- **Mock Generator**: confirm whether `internal/engine/mocks/` already has a Body/Movable mock that needs the four new methods each. If yes, extend; if no, no action.
- **TDD Specialist**: prefer table-driven tests for accessor round-trips; one focused behavioural test per AC; one negative-Y / negative-altitude row each to lock the "no clamping in Phase 1" contract.
- **Feature Implementer**: implement in the order Body → MovableBody → draworder → scene wiring → beatemup scaffold. Each step compiles independently.
- **Workflow Gatekeeper**: verify coverage delta on `internal/engine/physics/body/` is ≥ 0 and that no existing test was modified to mask a regression.
