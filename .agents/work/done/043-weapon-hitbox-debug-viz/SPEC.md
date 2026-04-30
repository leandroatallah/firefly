# Technical Specification — 043-weapon-hitbox-debug-viz

**Branch:** `043-weapon-hitbox-debug-viz`

**Bounded Contexts touched:**
- `internal/engine/render/camera/` (Render — debug overlay primitives)
- `internal/engine/combat/projectile/` (Combat — projectile manager Draw path)
- `internal/engine/combat/weapon/` and `internal/engine/combat/melee/` (Combat — melee hitbox exposure)
- `internal/game/scenes/phases/` (Game Logic — Phase Scene Draw orchestration)

## 1. Goal

Extend the existing `--collision-box` debug flag (already drawing Actor and obstacle collision boxes) to also render:
1. Each active projectile's collision box in the existing non-obstructive (green) style.
2. The active melee hitbox rectangle in a distinct orange style, only during the active-frame window.

The change is purely visual; it must not alter projectile/melee logic, hit detection, damage application, or Actor state transitions.

## 2. Background — Current Render Path

`PhasesScene.Draw` (in `internal/game/scenes/phases/scene.go`) iterates `space.Bodies()` and, for `PlatformerActorEntity` / `items.Item` / `body.Obstacle`, calls `s.Camera().DrawCollisionBox(screen, sb)` when `config.Get().CollisionBox` is true. Projectile bodies are managed by a separate `ProjectileManager` and rendered via `DrawWithOffset(screen, camX, camY)` — they are intentionally not in `space.Bodies()` rendering path of the scene loop, so the existing flag never triggers their collision-box draw.

`MeleeWeapon` (in `internal/engine/combat/weapon/melee.go`) computes a private `hitboxRect()` in pixel space when `IsHitboxActive()` is true. The rect is currently never exposed beyond `ApplyHitbox`.

## 3. Technical Requirements

### 3.1 Engine — Camera debug overlay (orange hitbox primitive)

In `internal/engine/render/camera/camera.go`, add a new method on `*Controller`:

```go
// DrawHitboxRect renders an orange debug rectangle in world space.
// Used by the Phase Scene to visualize an active melee hitbox under the
// --collision-box flag. Caller is responsible for gating on the flag.
func (c *Controller) DrawHitboxRect(screen *ebiten.Image, rect image.Rectangle)
```

Visual contract: outer dark-orange border + inner orange fill, mirroring the two-pass outer/inner pattern of `DrawCollisionBox` so it visually rhymes with existing boxes but is colour-distinct.

- Outer pass: `ColorScale.Scale(0.66, 0.33, 0, 1)` (dark orange).
- Inner pass (only when `rect.Dx() > 2 && rect.Dy() > 2`): `ColorScale.Scale(1, 0.5, 0, 1)` (orange).
- A degenerate rect (`rect.Dx() <= 0 || rect.Dy() <= 0`) is a no-op (no draw calls, no panic).
- Uses the same `collisionBoxImage` 1x1 white texture; transforms via `c.Draw(...)` so camera offset is applied identically to existing collision boxes.

`internal/game/render/camera/camera.go` (game-layer wrapper) gets a thin delegating method:

```go
func (c *Controller) DrawHitboxRect(screen *ebiten.Image, rect image.Rectangle) {
    c.base.DrawHitboxRect(screen, rect)
}
```

### 3.2 Engine — Expose active melee hitbox rect

In `internal/engine/combat/weapon/melee.go`, promote `hitboxRect()` to an exported method and add an active-aware accessor:

```go
// HitboxRect returns the current swing's hitbox rectangle in pixel space
// for the currently selected combo step. Caller should gate on
// IsHitboxActive() if they want frame-accurate visibility.
func (w *MeleeWeapon) HitboxRect() image.Rectangle

// ActiveHitboxRect returns (rect, true) only while IsHitboxActive() is true,
// otherwise (zero, false). This is the surface intended for debug rendering.
func (w *MeleeWeapon) ActiveHitboxRect() (image.Rectangle, bool)
```

`HitboxRect()` is the existing `hitboxRect()` body, exported. `ActiveHitboxRect()` short-circuits when not in the active-frame window.

Pre-conditions:
- Returns `(image.Rectangle{}, false)` when `!IsHitboxActive()` (no swing, in startup, or outside `ActiveFrames` window of the current step).

Post-conditions:
- `ActiveHitboxRect()` returns the same rectangle that `ApplyHitbox` would `space.Query` against during the same frame — they must be byte-identical for any given (step, faceDir, originX16, originY16).
- Calling these methods has no side effects: cooldown, swing frame, hit map are not touched.

### 3.3 Engine — Combat melee state surface

In `internal/engine/combat/melee/state.go`, the `weaponIface` already declares `IsHitboxActive()`. No new surface needed in `weaponIface` (the Phase Scene queries the `*MeleeWeapon` concrete type via the player accessor below).

### 3.4 Engine — Projectile Manager debug draw

In `internal/engine/combat/projectile/manager.go`, add:

```go
// DrawCollisionBoxesWithOffset renders each active projectile's collision box
// using the given camera-space draw helper. The helper is invoked once per
// active projectile body.
func (m *Manager) DrawCollisionBoxesWithOffset(draw func(b body.Collidable))
```

Rationale: The Manager already owns the slice of live `*projectile` instances. It exposes them to the caller through a callback so the camera-rendering decision (which colour scheme, where to draw) stays in the rendering layer — this avoids importing the camera package into combat.

Pre-conditions:
- `draw` is non-nil. Nil callback is a no-op (defensive guard).
- `m.projectiles` may be empty — callback is invoked zero times.

Post-conditions:
- Callback is invoked exactly once per active projectile (those still present in `m.projectiles` after the latest `Update`).
- The `body.Collidable` passed to the callback is the projectile's collision body (`p.body`), already registered in the physics space with a non-empty `CollisionPosition()`.
- No mutation of projectile state (position, lifetime, hit state) occurs.

### 3.5 Game — Player accessor for melee weapon

In `internal/game/entity/actors/player/climber.go`, add a read-only accessor:

```go
// MeleeController returns the per-actor melee Controller, or nil if none is installed.
func (p *ClimberPlayer) MeleeController() *meleeengine.Controller { return p.melee }
```

And in `internal/engine/combat/melee/controller.go`, add a read-only weapon accessor:

```go
// Weapon returns the underlying MeleeWeapon. Used by the Phase Scene to query
// the active hitbox rect for debug rendering.
func (c *Controller) Weapon() *weapon.MeleeWeapon { return c.weapon }
```

Both are pure getters — no state change, no allocation.

### 3.6 Game — Phase Scene wiring (the integration point)

In `internal/game/scenes/phases/scene.go` `Draw`, after the existing body loop and after `ProjectileManager.DrawWithOffset(...)`, add a single new debug block guarded by `config.Get().CollisionBox`:

```go
if config.Get().CollisionBox {
    // AC-1: projectile collision boxes (green/non-obstructive style)
    if pm := s.AppContext().ProjectileManager; pm != nil {
        pm.DrawCollisionBoxesWithOffset(func(b body.Collidable) {
            s.Camera().DrawCollisionBox(screen, b)
        })
    }

    // AC-2/AC-3: active melee hitbox (orange), frame-accurate
    if s.hasPlayer && s.player != nil {
        if cp, ok := s.player.(*player.ClimberPlayer); ok {
            if mc := cp.MeleeController(); mc != nil {
                if rect, active := mc.Weapon().ActiveHitboxRect(); active {
                    s.Camera().DrawHitboxRect(screen, rect)
                }
            }
        }
    }
}
```

(The `*player.ClimberPlayer` type assertion is acceptable here because Phase Scene already imports the concrete game-layer player package; the accessor is added in 3.5.)

## 4. State Machine Changes

None. Melee state machine (`internal/engine/combat/melee/state.go`) is unchanged. The debug overlay only reads `IsHitboxActive()` and `HitboxRect()`; it does not call `Update`, `OnStart`, `OnFinish`, `Fire`, or `ApplyHitbox`.

## 5. Pre/Post-Conditions Summary

| Behaviour | Pre-condition | Post-condition |
|---|---|---|
| Projectile collision box | `--collision-box=true` AND `len(ProjectileManager.projectiles) > 0` | One green outer+inner box per active projectile, every frame, drawn through camera. |
| Melee hitbox (active) | `--collision-box=true` AND `MeleeWeapon.IsHitboxActive() == true` | One orange outer+inner box at `MeleeWeapon.HitboxRect()`, drawn through camera, every frame the predicate holds. |
| Melee hitbox (inactive) | `--collision-box=true` AND `MeleeWeapon.IsHitboxActive() == false` | Zero hitbox draw calls. |
| Flag off | `--collision-box=false` | Zero new draw calls (projectile boxes, melee hitbox). |
| Game logic invariant | (any) | `MeleeWeapon` swing frame, cooldown, combo step, `hitThisSwing`, `Manager.projectiles` slice, projectile positions/lifetimes are byte-identical with and without the flag. |

## 6. Integration Points

| Producer | Consumer | Surface |
|---|---|---|
| `combat/projectile.Manager` | `phases.PhasesScene.Draw` | new `DrawCollisionBoxesWithOffset(draw func(body.Collidable))` |
| `combat/weapon.MeleeWeapon` | `phases.PhasesScene.Draw` (via `melee.Controller.Weapon()`) | new `HitboxRect()`, `ActiveHitboxRect()` |
| `combat/melee.Controller` | `player.ClimberPlayer.MeleeController()` → Phase Scene | new `Weapon()` getter |
| `game/entity/actors/player.ClimberPlayer` | `phases.PhasesScene.Draw` | new `MeleeController()` getter |
| `engine/render/camera.Controller` | `game/render/camera.Controller` → Phase Scene | new `DrawHitboxRect(screen, rect)` |

No new contracts required in `internal/engine/contracts/`. All additions are concrete-type method extensions on existing engine structs. The debug path stays a one-way read from combat → render; no new interfaces, no new mocks.

## 7. Out of Scope (mirrored from User Story)

- No new CLI flags. `--collision-box` is the sole gate.
- Enemy melee hitboxes — only the player's `MeleeController` is queried.
- Performance optimisations for large projectile counts.
- Visual polish beyond the colour contract.

## 8. Red Phase — Failing Test Scenarios

The TDD Specialist must produce the following failing tests **before** any production code in this story is written.

### 8.1 `internal/engine/combat/weapon/melee_hitbox_rect_test.go` — `TestMeleeWeapon_ActiveHitboxRect`

Table-driven across:

| Case | Setup | Expected |
|---|---|---|
| not swinging | freshly constructed weapon, `IsHitboxActive()==false` | `ActiveHitboxRect()` returns `(image.Rectangle{}, false)` |
| in startup | step with `StartupFrames=3`, `Fire` called, `Update` ticked once | returns `(zero, false)` |
| active window, faceRight | swing at `ActiveFrames[0]`, `faceDir=Right`, known `originX16/Y16`, `HitboxOffsetX16/HitboxW16/HitboxH16` | returns `(expectedRect, true)` where rect equals manual computation in pixel space |
| active window, faceLeft | mirrored offset | rect mirrored across origin (centre shifts to `originX16 - offsetX16`) |
| past active window | tick swing past `ActiveFrames[1]` | returns `(zero, false)` |
| `HitboxRect()` parity | during active window | `HitboxRect()` and `ActiveHitboxRect()` rect agree byte-for-byte; both equal what `ApplyHitbox` queried (assertable via the existing query-recording mock used by `TestMeleeWeapon_ApplyHitbox_FactionGating`) |

The test must NOT touch `Update`/`Fire` paths beyond the minimum needed to drive frame state; uses existing test fixtures from `melee_test.go`.

### 8.2 `internal/engine/combat/projectile/manager_debug_test.go` — `TestManager_DrawCollisionBoxesWithOffset`

Table-driven:

| Case | Setup | Expected |
|---|---|---|
| no projectiles | fresh `Manager` with mock `BodiesSpace` | callback invoked 0 times |
| single projectile | `Spawn(...)` once | callback invoked exactly 1 time, with a `body.Collidable` whose `ID()` matches the spawned body's ID |
| multiple projectiles | spawn 3 | callback invoked 3 times; collected IDs equal the set of spawned IDs |
| nil callback | spawn 1, pass `nil` | no panic, no side effects |
| no mutation | spawn 1, snapshot (`len(projectiles)`, position via `GetPositionMin`), call helper | post-call snapshot is identical |

### 8.3 `internal/engine/render/camera/camera_debug_hitbox_test.go` — `TestController_DrawHitboxRect`

Table-driven on rect dimensions:

| Case | Rect | Expected |
|---|---|---|
| degenerate (zero) | `image.Rect(0,0,0,0)` | no panic; underlying `kamera.Draw` not invoked (assert via spying on a wrapped `*ebiten.Image` target — count non-zero pixels remains 0) |
| 1x1 | `image.Rect(10,10,11,11)` | exactly one outer pass drawn (no inner because `Dx<=2`) |
| 4x4 | `image.Rect(10,10,14,14)` | both outer and inner passes drawn; sampled centre pixel has orange channel dominant; sampled border pixel has dark-orange channel dominant |
| negative size | `image.Rect(20,20,10,10)` | no panic; no draw |

These pixel-level assertions follow the existing pattern in `internal/engine/render/camera/camera_test.go` (`ebiten.NewImage` headlessly + `At(x,y)` colour read).

### 8.4 `internal/game/scenes/phases/scene_collision_debug_test.go` — `TestPhasesScene_Draw_CollisionBoxFlag`

Acceptance-level integration test, table-driven on the `CollisionBox` config:

| Case | Setup | Expected |
|---|---|---|
| AC-1: flag on, projectiles active | `config.Set(... CollisionBox=true ...)`, spawn 2 projectiles via `ProjectileManager`, run `Draw` on a headless image | `Manager.DrawCollisionBoxesWithOffset` is invoked (verified via test seam: a thin counter wrapper around the real Manager OR by sampling the screen for green pixels at the projectile positions) |
| AC-2: flag on, melee hitbox active | force `MeleeWeapon` into mid-active-frame state, run `Draw` | screen contains orange pixels overlapping `MeleeWeapon.HitboxRect()` bounds |
| AC-3: flag on, melee hitbox inactive | weapon constructed but never fired, run `Draw` | screen contains zero orange pixels |
| AC-4: flag off, projectiles active and swing active | `CollisionBox=false`, spawn 2 projectiles, force active swing, run `Draw` | screen contains zero green-debug pixels at projectile positions and zero orange pixels at hitbox region |
| AC-5: no logic side-effects | run two parallel scenarios identical except for `CollisionBox` toggle, advance N frames | post-frame snapshots of `MeleeWeapon.swingFrame`, `MeleeWeapon.currentCooldown`, `Manager.projectiles[i].body.GetPosition16()` for every i, and `Manager.projectiles[i].currentLifetime` are identical between the two scenarios |

Implementation notes for the TDD Specialist:
- Reuse the existing test scaffolding pattern in `internal/game/scenes/phases/mocks_test.go`.
- For AC-5, freezing test for game-logic invariance requires running `Update` (not just `Draw`); the assertion compares state after a fixed number of ticks (e.g., 60 frames during which a swing starts, peaks, and ends, and several projectiles fly).
- Pixel sampling uses `ebiten.NewImage(w, h)` + `screen.At(x, y).RGBA()` — the colour channel dominance check (orange: R>>B, R>G; green: G>R, G>B) is sufficient to prove the colour contract without coupling to exact RGBA values.

## 9. Constitution Compliance Notes

- No global mutable state introduced; the only existing global (`collisionBoxImage` in `camera.go`) is reused.
- Fixed-point: `MeleeWeapon.HitboxRect()` returns pixel-space `image.Rectangle` already (existing `hitboxRect()` does the `/16` conversion); we do not invent new fp16 math.
- No `_ = variable` patterns; no new mocks at non-boundary points (Manager helper takes a function, not an interface, but if a boundary mock is preferred the TDD Specialist may instead introduce a `body.Collidable`-emitting iterator interface — at the Specialist's discretion, equally compliant).
- Tests are headless (no `ebiten.RunGame`), table-driven, deterministic, no `time.Sleep`.
- Coverage delta: positive across `internal/engine/render/camera/`, `internal/engine/combat/weapon/`, `internal/engine/combat/projectile/`, and `internal/game/scenes/phases/`.

## 10. Key Design Decisions

1. **Callback over interface for projectile draw.** A `func(body.Collidable)` callback keeps the combat package free of any rendering concept and avoids an extra contract; the rendering decision (colour, camera) stays in the Phase Scene where it belongs. If the TDD Specialist prefers to mock at a boundary, they may introduce a minimal `ProjectileDebugSource` interface — both shapes satisfy the spec.
2. **Two methods on `MeleeWeapon` (`HitboxRect` + `ActiveHitboxRect`)** rather than one. `HitboxRect()` is testable independent of swing state, and `ActiveHitboxRect()` is the safe surface for the renderer. This separates "what would the rect be?" from "should I draw it now?".
3. **No new contract.** Adding a new contract for "thing that has a debug box" would over-generalise a single-feature debug path. Concrete-type extension is sufficient and reversible.
