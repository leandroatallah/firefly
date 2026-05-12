# USER STORY — 059-thin-game-phase-scenes

**Branch:** `059-thin-game-phase-scenes`

**Bounded Context:** Kit / Scene

## Story

As a developer building a game on the Firefly boilerplate,
I want the full Update/Draw phase-scene loop to live in the kit for both platformer and beat-em-up genres,
so that `internal/game/scenes/phases/` shrinks to ~30-line factory files that only inject game-specific types.

## Acceptance Criteria

- AC-1: Kit `PlatformerPhaseScene.Update` owns pause, sequence, goal, fall-death, VFX, screen-flipper, camera tick, body iteration with VFX-on-death, projectile update, and trigger resolution; game layer does not reimplement any of these.
- AC-2: Kit `BeatemupPhaseScene.Update` owns the same loop as AC-1 minus fall-death and screen-flipper; game layer does not reimplement any of these.
- AC-3: Kit `PlatformerPhaseScene.Draw` owns sorted-body draw, projectile draw, collision-box debug hook, screen flash, VFX draw, vignette, and pause overlay; game layer does not reimplement any of these.
- AC-4: Kit `BeatemupPhaseScene.Draw` owns the same draw pipeline as AC-3 minus vignette; uses `draworder.SortByGroundYAltitude` for depth sort; game layer does not reimplement any of these.
- AC-5: Each game-layer scene (`gameplatformerphase`, `gamebeatemupphase`) is a single `New*Scene(ctx)` constructor of ≤40 lines that builds the kit scene via an `Options` struct; no genre logic lives in the game layer.
- AC-6: `Options` struct carries `PlayerFactory func(*app.AppContext) (T, error)`, `ItemMap`, `EnemyMap`, `NpcMap`, `DebugDrawHook func(*ebiten.Image)`, `RebootSceneType`, and `MenuSceneType`; all fields are optional except `PlayerFactory`.
- AC-7: A generic `BuildPlayer[T actors.ActorEntity](p T, deps PlayerDeps) (T, error)` function in `internal/kit/actors/builder/` applies skills (always) and optionally wires `Inventory` and `MeleeWeapon` when non-nil; replaces `createPlayer` in both game-layer player files.
- AC-8: `internal/game/render/camera/` wrapper package is deleted; kit calls `camera.SetVerticalOnlyUpward(true)` directly on the engine camera for platformer scenes.
- AC-9: `platformer.PlatformerActorEntity` type switches exist only inside `internal/kit/scenes/phases/platformer/`; beat-em-up actor type switches exist only inside `internal/kit/scenes/phases/beatemup/`.
- AC-10: `SetDebugDrawHook(func(*ebiten.Image))` on kit `PlatformerPhaseScene` is wired from the game layer; no `*gameplayer.ClimberPlayer` cast appears in kit code.
- AC-11: `GoalType`, `ReactEndpointType`, `SequenceGoalType`, `NoGoalType`, and `ReachEndpointGoal` are defined in the engine `phases` package; game and kit import from there.
- AC-12: Tests exercising kit-owned behaviors move to `internal/kit/scenes/phases/{platformer,beatemup}/scene_test.go`; game-layer tests retain only factory/wiring assertions; layer import rules (kit must not import game) continue to pass.

## Behavioral Edge Cases

- Phase with no player: kit scene handles `PlayerFactory` returning nil without panicking; camera and draw loop run in degraded mode.
- Simultaneous death triggers: fall-death and state-machine `Dead` both route through a single guard; sequence fires at most once.
- Sequence-gated pause: pause key suppressed during active sequence; both kit scenes inherit this without re-implementing it.
- No-inventory boilerplate user: `BuildPlayer` with nil `Inventory` applies only skills; no nil-dereference.
- No-weapon boilerplate user: `BuildPlayer` with nil `MeleeWeapon` applies only skills; no nil-dereference.
- `DebugDrawHook` nil: kit `PlatformerPhaseScene.Draw` skips the hook call when nil; no panic.
