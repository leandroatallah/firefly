# SPEC — 055-kit-genre-phase-scenes

**Branch:** `055-kit-genre-phase-scenes`
**Bounded Context:** Kit / Scene
**Author:** Spec Engineer
**Status:** Active

---

## 1. Overview

Split `internal/game/scenes/phases/scene.go` into three layers per Constitution §Architecture:

1. **Engine** — `internal/engine/scene/phases/`: a genre-agnostic phase-scene base (`PhaseSceneBase`) providing pause, sequence, vignette, goal, body-iteration scaffolding, and template-method hooks.
2. **Kit** — `internal/kit/scenes/phases/platformer/PlatformerPhaseScene` and `internal/kit/scenes/phases/beatemup/BeatemupPhaseScene`: opinionated, genre-specific assemblies that embed the engine base and override hooks.
3. **Game** — `internal/game/scenes/phases/`: thin wiring that composes a kit scene with concrete game player/items/enemies/NPC factories and a debug-draw closure.

The story is a refactor with no behaviour change for the platformer (existing tests must continue to pass) plus a new `BeatemupPhaseScene` shell.

---

## 2. Pre-Spec Findings

### 2.1 Vertical-only-upward camera constraint location

`internal/engine/render/camera/camera.go` already exposes:

```go
type Controller struct {
    ...
    VerticalOnlyUpward bool
    lastTargetY        float64
    initialized        bool
}
```

The constraint logic is implemented inside `Controller.Update()` (lines 145–153). `internal/game/render/camera/Controller` is a **pure delegating wrapper** with no constraint logic of its own — confirming Risk #2 in the User Story. The constraint is therefore an *engine-level* configuration knob, not a game-layer behaviour.

**Decision:** the kit `PlatformerPhaseScene` will set `engineCamera.VerticalOnlyUpward = true` directly via a new accessor `Controller.SetVerticalOnlyUpward(bool)` on the engine camera. The whole `internal/game/render/camera/` wrapper package becomes redundant after the refactor and is deleted in this story.

### 2.2 Beat-em-up actor interface

`internal/kit/actors/beatemup/beatemup_character.go` exposes only the concrete `BeatEmUpCharacter` struct. A `BeatEmUpActorEntity` interface analogous to `PlatformerActorEntity` does not exist. **Decision:** introduce one in this story (see §4.2).

### 2.3 Altitude-aware draw sort

`internal/engine/contracts/body/body.go` already exposes `Altitude16()` on the relevant body interface (line 213). `draworder.SortByGroundY` only sorts by Y. **Decision:** introduce `draworder.SortByGroundYAltitude` in the engine, sorting by `(y16 + altitude16)` ascending; falls back to y16 if a body does not implement an `Altitudable` interface.

### 2.4 Game-specific debug draw

The active-melee-hitbox debug draw at line 518 of `scene.go` casts to `*gameplayer.ClimberPlayer`. **Decision:** the engine base exposes a `DebugDrawHook func(*ebiten.Image)` field invoked at the end of `Draw()`. The game wires it via a setter; the kit never imports any game type.

---

## 3. Engine Layer — `internal/engine/scene/phases/`

### 3.1 New type: `PhaseSceneBase`

File: `internal/engine/scene/phases/phase_scene_base.go`

Embeds `*scene.TilemapScene`. Owns all genre-agnostic phase-scene state: pause screen, pause menu, sequence player, goal, vignette, completion/death triggers, body counter, screen-flash counter, and lifecycle flags (`hasPlayer`, `allowPause`, `reachedEndpoint`).

```go
package phases

type PhaseSceneBase struct {
    *scene.TilemapScene

    count int

    // Goal & completion
    hasPlayer         bool
    reachedEndpoint   bool
    goal              Goal
    completionTrigger utils.DelayTrigger
    deathTrigger      utils.DelayTrigger
    deathActive       bool

    // UI / fx
    showDrawScreenFlash int
    vignette            *enginevfx.Vignette
    pauseScreen         *pause.PauseScreen
    pauseMenu           *menu.Menu
    allowPause          bool

    // Sequencing
    sequencePlayer sequencestypes.Player

    // Debug hook (game-injected)
    debugDrawHook func(*ebiten.Image)

    // Template hooks (set by subclasses; nil-safe)
    onUpdateHook        func() (handled bool, err error)
    onPreActorIterHook  func()
    onPostActorIterHook func()
    onDrawWorldHook     func(screen *ebiten.Image)
    onDrawOverlayHook   func(screen *ebiten.Image)
    actorBodyHandler    func(b body.Body, space physics.Space) (handled bool, err error)
    actorDrawHandler    func(screen *ebiten.Image, b body.Collidable) (handled bool)
    deathStarter        func()  // genre-specific death sequence; nil if genre has none
}
```

#### 3.1.1 Methods provided by base

- `OnStart()`: builds pause screen, pause menu, sequence player, vignette, body counter, goal init.
- `Update() error`: pause → sequence → vfx → `onUpdateHook` (genre overrides for screen-flip / fall-death / state-machine death) → triggers → goal check → camera update → body-iteration loop (using `actorBodyHandler` for genre-typed bodies, then engine-only `items.Item` and `body.Obstacle` cases) → projectile update → trigger collisions → `space.ProcessRemovals()`.
- `Draw(screen)`: tilemap → body draw loop using `draworder.SortByGroundY` (override-able sort key via `sortBodiesFunc func([]body.Collidable) []body.Collidable`) and `actorDrawHandler` for the genre case → projectiles → collision-box debug → `onDrawWorldHook` → vignette → pause draw → `debugDrawHook(screen)`.
- `OnFinish()`: clears projectiles, unblocks player movement (via `ActorManager.GetPlayer()`), unregisters primary actor.
- `SetDebugDrawHook(func(*ebiten.Image))`, `TriggerScreenFlash()`, `EnableVignetteDarkness(float64)`, `DisableVignetteDarkness()`.
- `SetGoal(Goal)`, `StartDeathSequence()` (template — calls `deathStarter` if set).
- Helpers: `freezeAllActors()`, `canPause()`, `refreshPauseMenuLabels()`, `drawPause()`.

#### 3.1.2 Hooks contract

- `actorBodyHandler` — called inside the body-update loop for every `body.Body` that is **not** an `items.Item` or `body.Obstacle`. Returns `handled=true` if the body was processed (continue loop); `handled=false` lets the base ignore the body. This is the **only** place a genre type-switch may occur outside its own kit package.
- `actorDrawHandler` — analogous, for the draw loop.
- `onUpdateHook` — runs after sequence/vfx update, before triggers; returns `handled=true` to short-circuit (e.g., during a screen-flip).
- `deathStarter` — invoked when the base detects a generic death-trigger condition (none in the engine; platformer wires fall-death and state-machine death detection through `onUpdateHook`).

### 3.2 Goal types kept in engine

Move `phases.Goal`, `SequenceGoal`, `NoGoal` (already in `internal/engine/scene/phases/goals.go`); also move `ReachEndpointGoal`, `ReactEndpointType`, `SequenceGoalType`, `NoGoalType` from `internal/game/scenes/phases/goals.go` and `goal_type.go` into the engine `phases` package — they are genre-agnostic.

### 3.3 New: `draworder.SortByGroundYAltitude`

File: `internal/engine/render/draworder/draworder.go` — add:

```go
type Altitudable interface { Altitude16() int }

func SortByGroundYAltitude(in []body.Collidable) []body.Collidable
```

Sort key per body `b`:
- `_, y16 := b.GetPosition16()`
- `alt16 := 0; if a, ok := b.(Altitudable); ok { alt16 = a.Altitude16() }`
- depth = `y16 - alt16` (higher altitude → drawn earlier / behind, since less ground-Y; final spec: `y16 - alt16` means a body lifted off the ground sorts as if its ground projection is higher on screen).

### 3.4 Engine camera setter

File: `internal/engine/render/camera/camera.go` — add:

```go
func (c *Controller) SetVerticalOnlyUpward(v bool) { c.VerticalOnlyUpward = v }
```

The struct field is already exported but a setter keeps the kit code idiomatic and lets us hide the field in a future pass.

### 3.5 No new contract interfaces required for engine base

`PhaseSceneBase` operates on `body.Body`, `body.Collidable`, `body.Obstacle`, `items.Item`, all already in `internal/engine/contracts/`. No new contract files in `internal/engine/contracts/` are introduced by the engine layer of this story.

---

## 4. Kit Layer

### 4.1 `internal/kit/scenes/phases/platformer/`

#### 4.1.1 Type

```go
package platformerphasescene

type PlatformerPhaseScene struct {
    *phases.PhaseSceneBase

    player        platformerkit.PlatformerActorEntity
    screenFlipper *scene.ScreenFlipper

    deathStateChecker func(actors.ActorStateEnum) bool // returns true for Dying/Dead in the host project's state map
    dyingState        actors.ActorStateEnum            // injected: gamestates.Dying value
    deadState         actors.ActorStateEnum            // injected: gamestates.Dead value

    playerFactory func(ctx *app.AppContext) (platformerkit.PlatformerActorEntity, error)
}

func NewPlatformerPhaseScene(ctx *app.AppContext, opts Options) *PlatformerPhaseScene
```

`Options` carries: `PlayerFactory`, `DyingState`, `DeadState`, `ItemMap`, `EnemyMap`, `NpcMap`, `RebootSceneType navigation.SceneType`, and an optional `DebugDrawHook`.

#### 4.1.2 Behaviour

- `OnStart`: chains to `PhaseSceneBase.OnStart`. If `Tilemap().HasPlayerStartPosition()`:
  - calls `opts.PlayerFactory(ctx)`, registers as primary, adds to space.
  - sets camera config to `CameraModeFollow`, sets `Camera().SetFollowTarget(player)`, calls `Camera().SetVerticalOnlyUpward(true)`.
  - constructs `screen_flipper` with `OnFlipStart`/`OnFlipFinish` toggling `player.SetImmobile`.
- Wires `actorBodyHandler` to a closure that type-switches on `platformerkit.PlatformerActorEntity` and emits death-explosion VFX + remove on `Dead` state, else calls `b.Update(space)`.
- Wires `actorDrawHandler` likewise (image options + collision box).
- Wires `onUpdateHook` to call `checkPlayerFallDeath()` and the state-machine death check (`player.State() == dyingState || deadState`).
- Provides `checkPlayerFallDeath()`: identical algorithm to current `scene.go` lines 264–287, but reads camera bottom from the engine `Camera()` directly.
- Provides `startDeathSequence()`: VFX explosion at player position, `player.GetCharacter().SetNewStateFatal(dyingState)`, `player.SetImmobile(true)`, enable `deathTrigger` (1 s), then on trigger fire the base navigates to `opts.RebootSceneType`.

#### 4.1.3 Vignette

Override `onDrawOverlayHook` only if a vignette darkness is enabled — but vignette draw stays in the base because it requires only `body.Body` (the player). Pass the player to base via `SetVignetteTarget(body.Body)` (added on base) so the base can draw vignette without genre coupling.

### 4.2 `internal/kit/actors/beatemup/`

Add interface:

```go
package beatemup

type BeatEmUpActorEntity interface {
    actors.ActorEntity
    context.ContextProvider
    Altitude16() int       // depth-sort component
    SetAltitude16(int)
}
```

`BeatEmUpCharacter` already satisfies `actors.ActorEntity` via the embedded `MeleeCharacter`. The interface only requires altitude accessors which are already on the body contract; verify and expose via embedding. If the underlying `MeleeCharacter` body does not surface `Altitude16()` on the entity, add a forwarding method on `BeatEmUpCharacter`.

### 4.3 `internal/kit/scenes/phases/beatemup/`

#### 4.3.1 Type

```go
package beatemupphasescene

type BeatemupPhaseScene struct {
    *phases.PhaseSceneBase

    player        beatemupkit.BeatEmUpActorEntity
    playerFactory func(ctx *app.AppContext) (beatemupkit.BeatEmUpActorEntity, error)
}

func NewBeatemupPhaseScene(ctx *app.AppContext, opts Options) *BeatemupPhaseScene
```

#### 4.3.2 Behaviour

- `OnStart`: chains base. If has-player: factory creates player, sets `CameraModeFollow`, `SetFollowTarget(player)`, **does not** call `SetVerticalOnlyUpward(true)` (arena-scrolling, no upward lock). Optional `Camera().SetBounds(...)` if the tilemap exposes arena bounds.
- No `screen_flipper`. No `checkPlayerFallDeath`.
- Sets base `sortBodiesFunc = draworder.SortByGroundYAltitude`.
- Wires `actorBodyHandler` to a closure that type-switches on `beatemupkit.BeatEmUpActorEntity` (Dead-state removal works without panic when the actor implements only `actors.ActorEntity` — altitude is read by the draw sort, not by removal logic).
- Wires `actorDrawHandler` to the BeatEmUp type case.
- `onUpdateHook` checks state-machine death only (no fall-death).
- No `screen_flipper` cancellation logic needed (Behavioral Edge Case 4 trivially holds: no flipper exists).

---

## 5. Game Layer — `internal/game/scenes/phases/`

After the refactor, this directory is reduced to:

- `scene.go`: a thin file constructing a `platformerphasescene.PlatformerPhaseScene` with:
  - `PlayerFactory = func(ctx) { return createPlayer(ctx, gameentitytypes.ClimberPlayerType) }`
  - `DyingState = gamestates.Dying`, `DeadState = gamestates.Dead`
  - `ItemMap = gameitems.InitItemMap(ctx)`, `EnemyMap = gameenemies.InitEnemyMap(ctx)`, `NpcMap = gamenpcs.InitNpcMap(ctx)`
  - `RebootSceneType = scenestypes.ScenePhaseReboot`
  - `DebugDrawHook` — closure that finds the player via `ctx.ActorManager.GetPlayer()`, type-asserts to `*gameplayer.ClimberPlayer`, draws active melee hitbox via `Camera().DrawHitboxRect`.
- `events.go`, `player.go`, `sequences.go`, `body_counter.go`: kept if they wire game-level events; `body_counter.go` is moved into the engine base. `goal_type.go`/`goals.go` move into engine `phases` (§3.2).
- `internal/game/render/camera/`: **deleted**; callers inside game switch to `engineCamera.Controller` directly.

`internal/game/app/setup.go` and `phases_list.go` are updated to construct the kit scene — already in-flight per `git status` showing local modifications. The Spec Engineer assumes those local edits will be reconciled by the Feature Implementer.

---

## 6. Pre-conditions / Post-conditions

### Pre-conditions
- `internal/engine/scene/TilemapScene` is unchanged.
- `platformer.PlatformerActorEntity` and `actors.ActorEntity` are stable.
- `engine/render/camera.Controller.VerticalOnlyUpward` semantics already validated by `camera_test.go`.

### Post-conditions
- No file under `internal/engine/` imports `internal/kit/` or `internal/game/`.
- No file under `internal/kit/scenes/phases/platformer/` imports `internal/game/`.
- No file under `internal/kit/scenes/phases/beatemup/` imports `internal/game/`.
- No file outside `internal/kit/scenes/phases/platformer/` contains `case platformer.PlatformerActorEntity` or `.(platformer.PlatformerActorEntity)`.
- No file outside `internal/kit/scenes/phases/beatemup/` references a beat-em-up actor type in a switch/assertion.
- No file under `internal/kit/` references `*gameplayer.ClimberPlayer`.
- All existing `internal/engine/...` and `internal/game/...` tests pass.
- New tests added per §8 pass.
- `internal/game/render/camera/` package is removed.

---

## 7. Integration Points (within Bounded Context)

- **Scene** (`internal/engine/scene/`): `PhaseSceneBase` embeds `TilemapScene`; uses `pause`, `transition`, `ScreenFlipper` (consumed by kit platformer scene only).
- **Sequences** (`internal/engine/sequences/`): unchanged; the base owns the `sequencePlayer`.
- **Physics / Body** (`internal/engine/physics/`, `internal/engine/contracts/body/`): base iterates space via `body.Body` / `body.Collidable`; kit overrides genre type-switch via hooks.
- **Render / Camera** (`internal/engine/render/camera/`): new `SetVerticalOnlyUpward` setter; existing `VerticalOnlyUpward` field continues to drive `Update()`.
- **Render / Draworder**: new `SortByGroundYAltitude` consumed by `BeatemupPhaseScene` only.
- **App / AppContext**: unchanged.

---

## 8. Red Phase — Failing Test Scenarios

The TDD Specialist must produce the following failing tests **before** any production code is written. All tests are table-driven where applicable, headless (`ebiten.NewImage`), no `time.Sleep`, no `ebiten.RunGame`.

### 8.1 `internal/kit/scenes/phases/platformer/scene_test.go`

#### T-P1: `TestPlatformerPhaseScene_CheckPlayerFallDeath_FiresWhenBelowCamera`
- Build a `PlatformerPhaseScene` with a stub player (mock `PlatformerActorEntity` returning a position whose top Y exceeds the camera bottom).
- Camera positioned at center (100,100), screen height 200 → bottom at Y=200. Set player top Y=250.
- Call `checkPlayerFallDeath()`.
- Assert: `deathActive == true` AND mock player received `GetCharacter().SetNewStateFatal(Dying)` AND `SetImmobile(true)`.

#### T-P2: `TestPlatformerPhaseScene_CheckPlayerFallDeath_NoOpWhenDeathActive`
- Same setup but pre-set `deathActive = true`.
- Call `checkPlayerFallDeath()` twice.
- Assert: mock player's `SetNewStateFatal` was called **0 times**.

#### T-P3: `TestPlatformerPhaseScene_ScreenFlipperCallbacksToggleImmobility`
- Construct scene with a stub player and tilemap with two rooms.
- Trigger `OnFlipStart` callback; assert mock player received `SetImmobile(true)`.
- Trigger `OnFlipFinish`; assert mock player received `SetImmobile(false)` exactly once.

#### T-P4: `TestPlatformerPhaseScene_NoPlayer_DoesNotPanic`
- Tilemap reports `HasPlayerStartPosition() == false`.
- Call `OnStart`, `Update`, `Draw(headlessImage)`.
- Assert: no panic; `screenFlipper == nil`; camera is in fixed mode.

#### T-P5: `TestPlatformerPhaseScene_DebugDrawHookInvoked`
- Set `DebugDrawHook` to a closure that increments a counter.
- Call `Draw(headlessImage)`.
- Assert: counter incremented once.

### 8.2 `internal/kit/scenes/phases/beatemup/scene_test.go`

#### T-B1: `TestBeatemupPhaseScene_DrawOrderSortsByGroundYPlusAltitude`
- Construct the scene; populate space with three stub bodies implementing `BeatEmUpActorEntity` with positions/altitudes:
  - A: y16=10·fp16, alt16=0       → effective 10
  - B: y16= 5·fp16, alt16=8·fp16  → effective −3
  - C: y16=20·fp16, alt16=5·fp16  → effective 15
- Capture draw order via spy `actorDrawHandler`.
- Assert order = [B, A, C] (ascending by `y16 - alt16`).

#### T-B2: `TestBeatemupPhaseScene_RemovesDeadActor_NoAltitudePanic`
- Stub actor implements only `actors.ActorEntity` + `BeatEmUpActorEntity` minimal subset; sets state to `Dead`.
- Run one `Update()` tick.
- Assert: `space.RemoveBody(actor)` was called; no panic.

#### T-B3: `TestBeatemupPhaseScene_NoFallDeathPath`
- Place player far below the camera bottom.
- Run `Update()`.
- Assert: `deathActive == false`; player did not receive `SetNewStateFatal`.

#### T-B4: `TestBeatemupPhaseScene_CameraNotVerticalOnlyUpward`
- After `OnStart`, assert `engineCamera.VerticalOnlyUpward == false`.

### 8.3 `internal/engine/render/draworder/draworder_test.go` (extension)

#### T-D1: `TestSortByGroundYAltitude` — table-driven over 4 inputs covering:
- All non-altitude bodies → identical to `SortByGroundY`.
- Mixed altitudable + non-altitudable bodies → non-altitudable treated as alt=0.
- Stable order on equal effective depth.

### 8.4 `internal/engine/render/camera/camera_test.go` (extension)

#### T-C1: `TestController_SetVerticalOnlyUpward` — toggling the setter is read back via `VerticalOnlyUpward` field; running `Update()` with a target moving downward does not advance camera Y when the flag is on, advances when off (table-driven on/off).

### 8.5 Regression
- All current tests in `internal/game/scenes/phases/scene_test.go` and `scene_collision_debug_test.go` must continue to pass after migration. If a test referenced removed types, it must be ported (not deleted) to `internal/kit/scenes/phases/platformer/`.

---

## 9. Out of Scope (re-stated)

- Goal-system redesign.
- Engine camera internal rewrite (only a setter is added).
- Beat-em-up gameplay (jump/altitude gravity/ground detection — story 053 territory).
- Audio, i18n, VFX subsystem changes.
- Scene navigation identifier renames.

---

## 10. Risks & Mitigations

| Risk | Mitigation |
|---|---|
| Coverage regression in `internal/game/scenes/phases/` (story risk #5) | New kit packages own the migrated tests; Gatekeeper measures aggregate `internal/kit/...` + `internal/engine/...` + `internal/game/...` coverage. |
| Hidden coupling from `events.go`/`sequences.go`/`player.go` to `gameplayer.ClimberPlayer` | Feature Implementer reads each file before migration; any game-typed reference stays in `internal/game/`, only generic logic is moved. |
| Engine base hooks become a "god struct" | Hooks are plain function fields, nil-safe, with single-responsibility names; no dynamic dispatch beyond the four hooks listed. |
| Beat-em-up scene tests need a stubbed `BeatEmUpActorEntity` | Mock Generator produces a shared mock in `internal/engine/mocks/` only if the interface is added to `internal/engine/contracts/`; otherwise package-local `mocks_test.go` per Constitution §Tests. **Decision:** keep `BeatEmUpActorEntity` in `internal/kit/actors/beatemup/` (kit interface, not engine contract) and use a package-local mock. |

---

## 11. Mock & Contract Inventory

- **No new files in `internal/engine/contracts/`** — the engine base operates entirely on existing contracts.
- **New kit interfaces** (no shared mocks generated; package-local mocks only):
  - `internal/kit/actors/beatemup/BeatEmUpActorEntity`
- **Existing kit interface reused**:
  - `internal/kit/actors/platformer/PlatformerActorEntity`

Therefore the **Mock Generator stage may be skipped** — mocks live in `mocks_test.go` of each new test package. If the developer prefers shared mocks, they can be added to `internal/engine/mocks/` later without spec changes.

---

## 12. Acceptance-Criteria Traceability

| AC | Section |
|---|---|
| AC-1 (engine base) | §3 |
| AC-2 (PlatformerPhaseScene) | §4.1 |
| AC-3 (BeatemupPhaseScene) | §4.2, §4.3 |
| AC-4 (game shrinks) | §5 |
| AC-5 (no cross-genre type switches) | §3.1.2 hooks; §6 post-conditions |
| AC-6 (debug draw hook) | §3.1.1 `SetDebugDrawHook`; §5 wiring |
| AC-7 (tests) | §8 |
| AC-8 (layer import rules) | §6 post-conditions |
