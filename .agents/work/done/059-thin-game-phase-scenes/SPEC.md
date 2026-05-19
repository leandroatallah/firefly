# SPEC — 059-thin-game-phase-scenes

**Branch:** `059-thin-game-phase-scenes`
**Bounded Context:** Kit / Scene

## 1. Layer Map

| Concern | Package | New / Modified |
|---|---|---|
| Goal type constants | `internal/engine/scene/phases` | modified |
| Generic player builder | `internal/kit/actors/builder` | new package |
| Kit platformer phase scene | `internal/kit/scenes/phases/platformer` | modified (absorbs full loop) |
| Kit beat-em-up phase scene | `internal/kit/scenes/phases/beatemup` | modified (absorbs full loop) |
| Game platformer phase scene | `internal/game/scenes/phases/platformer` | shrunk to factory |
| Game beat-em-up phase scene | `internal/game/scenes/phases/beatemup` | shrunk to factory |
| Game camera wrapper | `internal/game/render/camera` | **deleted** |

## 2. Engine: Goal Types [AC-11]

File: `internal/engine/scene/phases/goals.go` (extend existing).

```go
// Move from internal/game/scenes/phases/goal_type.go to here.
var (
    ReactEndpointType GoalType = "reach_endpoint"
    SequenceGoalType  GoalType = "sequence"
    NoGoalType        GoalType = "no_goal"
)

// ReachEndpointGoal: completes when a flag is flipped via Reach().
type ReachEndpointGoal struct {
    reached       bool
    OnCompletion_ func() // optional callback (game-layer freeze/audio fade)
}
func (g *ReachEndpointGoal) IsCompleted() bool { return g.reached }
func (g *ReachEndpointGoal) OnCompletion()     { if g.OnCompletion_ != nil { g.OnCompletion_() } }
func (g *ReachEndpointGoal) Reach()            { g.reached = true }
```

Delete `internal/game/scenes/phases/goal_type.go` and `internal/game/scenes/phases/{platformer,beatemup}/goals.go`.
Game packages import `phases.ReactEndpointType` etc.

## 3. Kit: Generic Player Builder [AC-7]

File: `internal/kit/actors/builder/builder.go` (new package).

```go
package kitbuilder

import (
    "github.com/boilerplate/ebiten-template/internal/engine/contracts/vfx"
    "github.com/boilerplate/ebiten-template/internal/engine/data/schemas"
    "github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
    enginebuilder "github.com/boilerplate/ebiten-template/internal/engine/entity/actors/builder"
    kitskills "github.com/boilerplate/ebiten-template/internal/kit/skills"
    "github.com/boilerplate/ebiten-template/internal/kit/combat/weapon"
)

// PlayerDeps configures optional wiring applied by BuildPlayer.
type PlayerDeps struct {
    SkillDeps    kitskills.SkillDeps          // required if SpriteData has skills
    Inventory    interface{}                  // optional; applied via SetInventory if non-nil
    MeleeWeapon  *weapon.MeleeWeapon          // optional; applied via SetMelee if non-nil
    VFXManager   vfx.Manager                  // passed to SetMelee
    SpriteData   *schemas.SpriteData          // required if applying skills; nil ⇒ skip skills
    WireState    func(*actors.Character)      // optional; e.g., WireStateContributors
}

// playerWiring is the optional interface used to inject inventory/melee.
type playerWiring interface {
    SetInventory(interface{})
    SetMelee(w *weapon.MeleeWeapon, vfxMgr vfx.Manager)
    GetCharacter() *actors.Character
}

// BuildPlayer applies skills (when SpriteData non-nil), then optionally injects
// Inventory and MeleeWeapon. Returns p untouched on success.
func BuildPlayer[T actors.ActorEntity](p T, deps PlayerDeps) (T, error) {
    pw, ok := any(p).(playerWiring)
    if !ok {
        return p, nil // p does not opt-in to inventory/melee wiring
    }
    if deps.Inventory != nil {
        pw.SetInventory(deps.Inventory)
    }
    if deps.MeleeWeapon != nil {
        pw.SetMelee(deps.MeleeWeapon, deps.VFXManager)
    }
    if deps.SpriteData != nil {
        skills := kitskills.FromConfig(deps.SpriteData.Skills, deps.SkillDeps)
        if err := enginebuilder.ApplySkills(p, skills); err != nil {
            return p, err
        }
    }
    if deps.WireState != nil {
        deps.WireState(pw.GetCharacter())
    }
    return p, nil
}
```

## 4. Kit Scene Options [AC-5, AC-6]

Files: `internal/kit/scenes/phases/platformer/options.go`, `internal/kit/scenes/phases/beatemup/options.go`.

```go
type Options[T any] struct {
    Ctx             *app.AppContext                            // required
    PlayerFactory   func(*app.AppContext) (T, error)           // required
    ItemMap         items.ItemMap                              // optional
    EnemyMap        enemies.EnemyMap[T]                        // optional (platformer; beatemup uses its own T)
    NpcMap          npcs.NpcMap[T]                             // optional
    DebugDrawHook   func(*ebiten.Image)                        // optional
    RebootSceneType navigation.SceneType                       // required for death routing
    MenuSceneType   navigation.SceneType                       // required for pause→menu
}
```

Validation in `New*Scene(opts)`:
- `Ctx == nil` ⇒ return constructor error (panic acceptable at init).
- `PlayerFactory == nil` ⇒ return constructor error.
- Other fields default to nil/zero; loop must nil-check at use.

## 5. Kit PlatformerPhaseScene [AC-1, AC-3, AC-9, AC-10]

File: `internal/kit/scenes/phases/platformer/scene.go` (expand).

New struct fields (in addition to existing):
```
tilemapScene      *scene.TilemapScene
appCtx            *app.AppContext
goal              phases.Goal
sequencePlayer    sequences.Player
allowPause        bool
pauseScreen       *pause.PauseScreen
pauseMenu         *menu.Menu
screenFlipper     *scene.ScreenFlipper
completionTrigger utils.DelayTrigger
deathTrigger      utils.DelayTrigger
rebootScene       navigation.SceneType
menuScene         navigation.SceneType
playerFactory     func(*app.AppContext) (Player, error)
itemMap           items.ItemMap
enemyMap          enemies.EnemyMap[Player]
npcMap            npcs.NpcMap[Player]
```

Replace existing `New(...)` with the Options-based factory:
```go
func NewWithOptions(opts Options[Player]) (*PlatformerPhaseScene, error)
```
Keep existing low-level `New(cam, space, sw, sh, dyingState, deadState)` for test-only paths and migration; mark Deprecated.

### 5.1 OnStart [AC-1]

Pseudocode:
```
OnStart:
  TilemapScene.OnStart()
  count = 0
  hasPlayer = Tilemap().HasPlayerStartPosition()
  if hasPlayer:
    p, err := playerFactory(ctx)
    on err: log.Fatal
    player = p
    ctx.ActorManager.Register(p); RegisterPrimary(p)
    space.AddBody(p)
    OnDeathStarted = func() { spawnDeathExplosion(p); deathTrigger.Enable(1s) }
    if phase.BlockPlayerMovement: ctx.ActorManager.GetPlayer().BlockMovement()
  initTilemap(itemMap, enemyMap, npcMap)
  bodyCounter.set(space)
  Tilemap().CreateCollisionBodies(space, touchTriggerFor(endpointTrigger))
  if hasPlayer:
    SetCameraConfig(Follow)
    camera.SetFollowing(true)
    camera.SetVerticalOnlyUpward(true)         // [AC-8] direct engine call
    camera.SetFollowTarget(player)
    screenFlipper = scene.NewScreenFlipper(camera, player, Tilemap(), ctx)
    screenFlipper.PlayerPushDistance = Tilemap().Tilewidth / 2
    screenFlipper.FlipStrategy = (dx,dy)→Instant if dy!=0 else Smooth
    screenFlipper.OnFlipStart  = ()→player.SetImmobile(true)
    screenFlipper.OnFlipFinish = ()→player.SetImmobile(false)
    screenFlipper.SnapToCurrentRoom()
  else:
    SetCameraConfig(Fixed)
    pos = Tilemap().GetCameraStartPosition() or (0,0)
    camera.SetPositionTopLeft(pos)
  buildPauseScreen(ctx, menuScene)
  buildSequencePlayer(ctx)   // sets allowPause=phase.GoalType!=SequenceGoalType
  initGoal(phase.GoalType)
```

### 5.2 Update [AC-1]

Pseudocode (exact order matters; mirrors current game-layer scene):
```
Update:
  if pauseScreen and canPause():
    pauseScreen.Update()
    if pauseScreen.IsPaused(): return nil
  if sequencePlayer: sequencePlayer.Update()
  if ctx.VFX: ctx.VFX.Update()
  if screenFlipper:
    screenFlipper.Update()
    if screenFlipper.IsFlipping(): return nil
  if hasPlayer:
    checkPlayerFallDeath()
    if !deathActive and (player.State() in {Dying,Dead}):
      startDeathSequence()
  completionTrigger.Update(); deathTrigger.Update()
  if deathTrigger.Trigger():
    SceneManager.NavigateTo(rebootScene, fader, false)
  if goal!=nil and goal.IsCompleted() and !completionTrigger.IsEnabled():
    goal.OnCompletion()
  if cfg.CamDebug: camera.CamDebug()
  if completionTrigger.Trigger():
    ctx.CompleteCurrentPhase(fader, true)
  camera.Update()
  BaseScene.Update()
  count++
  for b in space.Bodies():
    switch b.(type):
      case platformer.PlatformerActorEntity:           // [AC-9] only inside this package
        if b.State()==Dead: emit VFX death (size 30); space.RemoveBody(b); continue
        b.Update(space)
      case items.Item:
        if b.IsRemoved(): space.RemoveBody(b); continue
        b.Update(space)
      case body.Obstacle: continue
  if ctx.ProjectileManager: ctx.ProjectileManager.Update()
  if hasPlayer: space.ResolveCollisions(player)
  space.ProcessRemovals()
```

### 5.3 Draw [AC-3]

```
Draw(screen):
  screen.Fill(black)
  tilemapImg, _ := Tilemap().Image(screen)
  camera.Draw(tilemapImg, Tilemap().ImageOptions(), screen)
  for b in draworder.SortByGroundY(space.Bodies()):
    case PlatformerActorEntity / items.Item: camera.Draw(sb.Image(), sb.ImageOptions(), screen)
                                              if cfg.CollisionBox: camera.DrawCollisionBox(screen, sb)
    case body.Obstacle:                        if cfg.CollisionBox: camera.DrawCollisionBox(screen, sb)
  if ctx.ProjectileManager: ctx.ProjectileManager.DrawWithOffset(screen, camOffset)
  if cfg.CollisionBox and ctx.ProjectileManager: pm.DrawCollisionBoxesWithOffset(camera.DrawCollisionBox)
  if flashCount>0: screenutil.DrawScreenFlash(screen); flashCount--
  if ctx.VFX: ctx.VFX.Draw(screen, camera)
  if vignette and hasPlayer: vignette.Draw(screen, camera, player)
  if debugDrawHook != nil: debugDrawHook(screen)     // [AC-10] replaces ClimberPlayer cast
  if pauseScreen.IsPaused(): drawPause(screen)
```

### 5.4 Other Methods

```
SetDebugDrawHook(f func(*ebiten.Image))   // [AC-10]
TriggerScreenFlash()                       // existing, keep
EnableVignetteDarkness(radiusPx float64)
DisableVignetteDarkness()
canPause() bool: allowPause && (sequencePlayer==nil || !sequencePlayer.IsPlaying())
endpointTrigger(id string):
  if !hasPlayer: return
  if deathActive: return
  switch id:
    "SPIKE":    startDeathSequence()
    "CUTSCENE": pass
    default:    if goal is *phases.ReachEndpointGoal: goal.Reach()
```

## 6. Kit BeatemupPhaseScene [AC-2, AC-4, AC-9]

File: `internal/kit/scenes/phases/beatemup/scene.go` (expand).

Same fields and OnStart as platformer with these deltas:
- No `screenFlipper`, no `vignette`.
- `camera.SetBounds(tilemapRect)` after initTilemap.
- `camera.SetVerticalOnlyUpward(false)`.

Update loop deltas vs. platformer:
- **No** `checkPlayerFallDeath`.
- **No** `screenFlipper.Update`.
- Type switch is `case beatemupkit.BeatEmUpActorEntity` (the only place this assertion appears).

Draw deltas vs. platformer:
- Use `draworder.SortByGroundYAltitude(space.Bodies())`.
- **No** vignette.
- `debugDrawHook` still supported.

## 7. Game Layer: Thin Factories [AC-5]

### 7.1 `internal/game/scenes/phases/platformer/scene.go` (≤40 lines)

```go
func NewPlatformerPhaseScene(ctx *app.AppContext) *platformerphasescene.PlatformerPhaseScene {
    s, err := platformerphasescene.NewWithOptions(platformerphasescene.Options[platformerphasescene.Player]{
        Ctx:             ctx,
        PlayerFactory:   newClimberPlayer,      // local thin wrapper around BuildPlayer
        ItemMap:         gameitems.InitItemMap(ctx),
        EnemyMap:        gameenemies.InitEnemyMap(ctx),
        NpcMap:          gamenpcs.InitNpcMap(ctx),
        DebugDrawHook:   makeClimberDebugHook(ctx),    // captures *gameplayer.ClimberPlayer
        RebootSceneType: scenestypes.ScenePhaseReboot,
        MenuSceneType:   scenestypes.SceneMenu,
    })
    if err != nil { log.Fatal(err) }
    return s
}
```

Delete from this package: `Update`, `Draw`, `OnFinish`, `Camera`, `BaseCamera`, `initTilemap`, `initGoal`, `freezeAllActors`, `defaultCompletion`, `canPause`, `refreshPauseMenuLabels`, `drawPause`, `endpointTrigger`, `TriggerScreenFlash`, `EnableVignetteDarkness`, `DisableVignetteDarkness`, `body_counter.go`, `goals.go`, `events.go` (subscribeEvents moves into kit OnStart).

Files retained: `scene.go` (factory only), `player.go` (renamed `newClimberPlayer`, uses kitbuilder.BuildPlayer), debug-hook helper for ClimberPlayer melee hitbox.

### 7.2 `internal/game/scenes/phases/beatemup/scene.go` (≤40 lines)

Identical pattern; `DebugDrawHook` may be nil.

### 7.3 Player Factory Helpers (game layer)

```go
// newClimberPlayer wires the climber via kitbuilder.BuildPlayer.
func newClimberPlayer(ctx *app.AppContext) (platformerphasescene.Player, error) {
    p, err := gameplayer.NewClimberPlayer(ctx)
    if err != nil { return nil, err }
    return kitbuilder.BuildPlayer(p, kitbuilder.PlayerDeps{
        SpriteData:  p.GetSpriteData(),
        Inventory:   gameplayer.NewClimberInventory(ctx.ProjectileManager, ctx.VFX),
        MeleeWeapon: gameplayer.NewPlayerMeleeWeapon(),
        VFXManager:  ctx.VFX,
        SkillDeps:   buildSkillDeps(ctx),
        WireState:   func(c *actors.Character) { gameplayer.WireStateContributors(c, p) },
    })
}
```

### 7.4 Debug-Hook Helper (platformer game layer)

```go
func makeClimberDebugHook(ctx *app.AppContext) func(*ebiten.Image) {
    return func(screen *ebiten.Image) {
        if !config.Get().CollisionBox { return }
        // ClimberPlayer melee hitbox draw — uses ctx.ActorManager.GetPlayer().(*gameplayer.ClimberPlayer)
    }
}
```

The kit must not import `gameplayer`; the cast lives only inside this closure.

## 8. Camera Wrapper Deletion [AC-8]

Delete:
- `internal/game/render/camera/camera.go`
- `internal/game/render/camera/camera_test.go`

All call sites in `internal/game/scenes/phases/platformer/scene.go` are removed (the file shrinks to a factory). No other production package imports `gamecamera`.

## 9. Pre-Conditions & Post-Conditions

| Check | One-liner |
|---|---|
| C-1 | Kit `*PlatformerPhaseScene` exposes `Update`, `Draw`, `OnStart`, `OnFinish`, `SetDebugDrawHook`. |
| C-2 | `internal/game/render/camera/` directory no longer exists. |
| C-3 | `grep -r "PlatformerActorEntity" internal/kit/` → only inside `kit/scenes/phases/platformer/` and `kit/actors/platformer/`. |
| C-4 | `grep -r "BeatEmUpActorEntity" internal/kit/` → only inside `kit/scenes/phases/beatemup/` and `kit/actors/beatemup/`. |
| C-5 | `grep -r "ClimberPlayer\|CodyPlayer" internal/kit/` → 0 matches. |
| C-6 | `internal/game/scenes/phases/platformer/scene.go` ≤ 40 non-blank lines. |
| C-7 | `internal/game/scenes/phases/beatemup/scene.go` ≤ 40 non-blank lines. |
| C-8 | `go vet ./...` clean; `go test ./...` green. |
| C-9 | `phases.ReactEndpointType`, `phases.SequenceGoalType`, `phases.NoGoalType`, `phases.ReachEndpointGoal` resolve in engine package. |
| C-10 | `kitbuilder.BuildPlayer[T](p, PlayerDeps{Inventory:nil, MeleeWeapon:nil})` does not panic. |

## 10. Red Phase Tests

### 10.1 Kit Platformer Scene Update Tests
Location: `internal/kit/scenes/phases/platformer/scene_test.go`.

```
T-P1: Update fall-death triggers exactly once
  pre:  scene with player; player.TopY=250; camera.Bottom=200; deathActive=false
  act:  scene.Update(); scene.Update()
  post: deathActive==true; recordedSetNewStateFatal calls == 1; player.SetImmobile(true)==1

T-P2: Update routes Dead actors through space.Remove + VFX death explosion
  pre:  body b with State()=Dead added to space; ctx.VFX recorder
  act:  scene.Update()
  post: vfx.SpawnDeathExplosionCalls == 1; b not in space.Bodies()

T-P3: Update suppressed during active sequence-gated pause
  pre:  phase.GoalType=SequenceGoalType; sequencePlayer.IsPlaying()=true
  act:  press pause; scene.Update()
  post: pauseScreen.IsPaused()==false (canPause returns false)

T-P4: Update skips player loop when screenFlipper.IsFlipping()
  pre:  screenFlipper.flipping=true; spaceBodies has 1 actor
  act:  scene.Update()
  post: actor.UpdateCalls == 0

T-P5: Simultaneous fall + state-Dead routes once
  pre:  player.TopY>camera.Bottom; player.State()=Dying; deathActive=false
  act:  scene.Update()
  post: startDeathSequence invocations == 1; deathTrigger.Enable called once

T-P6: nil PlayerFactory result handled gracefully
  pre:  PlayerFactory returns nil player
  act:  scene.OnStart(); scene.Update(); scene.Draw(img)
  post: no panic; hasPlayer==false; camera in Fixed mode
```

### 10.2 Kit Platformer Scene Draw Tests

```
T-P7: Draw calls debugDrawHook when set
  pre:  scene.SetDebugDrawHook(recorder); cfg.CollisionBox=false
  act:  scene.Draw(screen)
  post: recorder.calls == 1

T-P8: Draw skips debugDrawHook when nil
  pre:  debugDrawHook=nil
  act:  scene.Draw(screen)
  post: no panic

T-P9: Draw routes vignette only when hasPlayer
  pre:  hasPlayer=false; vignette.Enable(10)
  act:  scene.Draw(screen)
  post: vignette.DrawCalls == 0

T-P10: Draw uses SortByGroundY (NOT altitude)
  pre:  bodies with varying GroundY
  act:  scene.Draw(screen) with recording handler
  post: draw order matches draworder.SortByGroundY output
```

### 10.3 Kit Beat-em-up Scene Update Tests
Location: `internal/kit/scenes/phases/beatemup/scene_test.go`.

```
T-B1: Update does NOT perform fall-death check
  pre:  player.TopY=999; camera.Bottom=0
  act:  scene.Update()
  post: deathActive==false; setNewStateFatal calls == 0

T-B2: Update routes Dead BeatEmUpActorEntity through VFX + remove
  pre:  body b implementing BeatEmUpActorEntity; State()=Dead
  act:  scene.Update()
  post: vfx.SpawnDeathExplosionCalls == 1; b not in space

T-B3: Update advances player state-machine death to reboot trigger
  pre:  player.State()=Dying; deathActive=false
  act:  scene.Update(); advance deathTrigger; scene.Update()
  post: SceneManager.LastNavigateTo == rebootScene
```

### 10.4 Kit Beat-em-up Scene Draw Tests

```
T-B4: Draw uses SortByGroundYAltitude
  pre:  bodies with varying Altitude16
  act:  scene.Draw(screen) with recording handler
  post: draw order matches draworder.SortByGroundYAltitude output

T-B5: Draw does NOT call vignette
  pre:  scene built; any state
  act:  scene.Draw(screen)
  post: no vignette field exists (compile-time) OR no draw call recorded
```

### 10.5 BuildPlayer Tests
Location: `internal/kit/actors/builder/builder_test.go`.

```
T-BP1: BuildPlayer with nil Inventory applies only skills
  pre:  mockPlayer satisfies playerWiring; deps.Inventory=nil; deps.MeleeWeapon=nil; deps.SpriteData={Skills:[]}
  act:  BuildPlayer(p, deps)
  post: pw.SetInventoryCalls==0; pw.SetMeleeCalls==0; no error

T-BP2: BuildPlayer with nil MeleeWeapon skips melee wiring
  pre:  deps.Inventory=stubInv; deps.MeleeWeapon=nil
  act:  BuildPlayer(p, deps)
  post: pw.SetInventoryCalls==1; pw.SetMeleeCalls==0

T-BP3: BuildPlayer applies both when both non-nil
  pre:  deps.Inventory=stubInv; deps.MeleeWeapon=stubWpn
  act:  BuildPlayer(p, deps)
  post: pw.SetInventoryCalls==1; pw.SetMeleeCalls==1

T-BP4: BuildPlayer on player without playerWiring is a no-op
  pre:  mockPlayer NOT implementing playerWiring
  act:  BuildPlayer(p, deps{Inventory:x, MeleeWeapon:y})
  post: returns p; no error; no panic
```

### 10.6 Engine Goal Tests
Location: `internal/engine/scene/phases/goals_test.go`.

```
T-G1: ReachEndpointGoal.IsCompleted false before Reach
T-G2: ReachEndpointGoal.IsCompleted true after Reach
T-G3: ReachEndpointGoal.OnCompletion invokes callback when set
T-G4: ReachEndpointGoal.OnCompletion no-op when callback nil
```

### 10.7 Game-Layer Wiring Tests [AC-12]
Location: `internal/game/scenes/phases/platformer/scene_test.go`.

Retain only:
```
T-W1: NewPlatformerPhaseScene returns non-nil and Options.PlayerFactory != nil
T-W2: NewPlatformerPhaseScene DebugDrawHook closure captures ClimberPlayer cast
T-W3: Layer rule: package import graph excludes internal/game/* in kit
       (executed via go list -deps; or moved to a dedicated lint test)
```

## 11. Mock / Contract Inventory

| Mock | Where | Purpose |
|---|---|---|
| `mockPlayer` (impl `Player` + `playerWiring`) | `kit/scenes/phases/{platformer,beatemup}/mocks_test.go` | scene unit tests |
| `mockBodiesSpace` | already exists in kit scene test files | reuse |
| `mockSceneManager` | `internal/engine/mocks/` or kit-local | assert NavigateTo arg |
| `mockVFXManager` | kit-local or shared | record `SpawnDeathExplosion` calls |
| `mockSequencePlayer` | kit-local | toggle `IsPlaying` |
| `recorderDebugDrawHook` | kit-local | counts calls |
| `stubInventory`, `stubMeleeWeapon` | `kit/actors/builder/mocks_test.go` | BuildPlayer tests |

No new contracts are introduced (BuildPlayer uses an unexported `playerWiring` interface defined locally). **Mock Generator can be skipped.**

## 12. Migration Order (for implementer)

1. Move `GoalType` constants + `ReachEndpointGoal` to `internal/engine/scene/phases/goals.go`.
2. Create `internal/kit/actors/builder/` package with `BuildPlayer`.
3. Extend kit `Options` structs and `NewWithOptions` constructors (both genres).
4. Move Update/Draw/OnStart/OnFinish bodies into kit scenes.
5. Shrink game-layer scene files to factories; rewrite game-layer `player.go` to use `kitbuilder.BuildPlayer`.
6. Delete `internal/game/render/camera/` package.
7. Move/rename tests per AC-12.
