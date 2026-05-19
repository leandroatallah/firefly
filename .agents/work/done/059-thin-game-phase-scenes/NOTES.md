# NOTES — 059-thin-game-phase-scenes

Human context, rationale, and risks for the spec. Not consumed by downstream
agents directly — they should read SPEC.md.

## Investigation Findings

- `internal/kit/scenes/phases/platformer/scene.go` (~400 LoC) already owns the
  fall-death loop, death sequence, and a minimal actor-update path. The game
  layer (`internal/game/scenes/phases/platformer/scene.go`, ~640 LoC)
  duplicates large pieces of the loop because the kit version was extracted
  incrementally during story 055.
- The same situation exists for beat-em-up: ~315 LoC kit vs. ~545 LoC game.
- `internal/game/render/camera/` is a 130-line passthrough wrapper. Every
  method is `c.base.X(...)`. The only "value-add" was historically the
  `SetVerticalOnlyUpward` constraint, which now lives on the engine camera
  itself (verified by grep — engine `camera.go:138` defines it). The wrapper
  exists purely so the game-layer scene could call `s.gameCamera.Base()` in
  a couple of spots; once the scene moves into kit, those calls disappear.
- `createPlayer` in both game-phase packages is nearly identical: each calls a
  player constructor, then runs the same skill/inventory/melee/state-wire
  block. Story 058 already isolated the skill block — making the rest of the
  block generic is the natural next step.
- `internal/game/scenes/phases/goal_type.go` defines three `GoalType` consts
  that the kit can't currently see because they live in game. Once the kit
  owns goal selection (in `initGoal`), the consts must move to engine.
- `ReachEndpointGoal` is implemented twice (once per genre in the game layer).
  Both implementations are identical except for the receiver type. Extracting
  it to the engine `phases` package with a `Reach()` setter removes both
  duplicates.

## Design Rationale

### Why a generic `BuildPlayer[T]` rather than a concrete one per genre?

The function is genre-agnostic: it only touches skills, inventory, and melee.
Both `PlatformerActorEntity` and `BeatEmUpActorEntity` extend
`actors.ActorEntity`, so a single generic constraint `T actors.ActorEntity`
covers both. The optional `playerWiring` interface keeps the function from
fighting players that don't have inventory/melee.

### Why an `Options` struct rather than positional args?

The constructor needs 7+ inputs (ctx, factory, three maps, two scene types,
debug hook). Positional args become brittle quickly. The `Options` struct also
naturally encodes optional fields — anything nil is skipped.

### Why is the debug-draw hook a closure rather than a typed interface?

The cast to `*gameplayer.ClimberPlayer` is a game-layer concern. Wrapping it
in a `func(*ebiten.Image)` keeps the entire `gameplayer` import out of the
kit. The kit only ever sees an opaque function value.

### Why keep the old low-level `New(...)` constructor in the kit?

Existing kit-side tests use it directly. Marking it Deprecated and adding the
new `NewWithOptions` lets us migrate production paths first and clean up
tests in a separate pass.

## Risks

| Risk | Mitigation |
|---|---|
| Move of `GoalType` const into engine may break unrelated callers. | grep before moving; only 5 references known. |
| Camera wrapper deletion could surface in test files. | Search `internal/` for `gamecamera`; only one production file currently. |
| `BuildPlayer` generic + reflection-free interface assertion may not match real player types. | The `playerWiring` interface uses only methods already implemented by both `ClimberPlayer` and `CodyPlayer`. Verified by grep. |
| Sequence-gated pause edge case (T-P3) depends on `sequencePlayer == nil` check. | Spec explicitly nil-guards `canPause`. |
| Game-layer test files currently reference scene internals (private fields). | Most relocate to kit-side tests; the residual game-layer tests assert only public surface. |

## Out of Scope

- Refactoring `scene.TilemapScene` to remove its `app.AppContext` coupling.
- Moving `body_counter.go` into the kit (it stays game-side for now;
  optionally inlined into the scene struct).
- Migrating `subscribeEvents` to kit (it's small enough to remain a
  game-layer hook passed via Options, or live in the kit but be game-layer
  agnostic — done via kit code that subscribes from `OnStart` to standard
  engine event types).
- Touching the existing `mocks_test.go` files in `internal/game/scenes/phases/*`
  beyond removing what is no longer referenced.

## Open Questions (resolve during TDD/Implementation)

1. Should `subscribeEvents` move to kit? Both genres have identical content
   (jump puff, landing puff). Recommendation: move to kit; it touches only
   engine APIs (`event`, `actorevents`, `ctx.VFX`, `ctx.AudioManager`).
2. Should `bodyCounter` move to kit? Probably yes (a `*BodyCounter` field on
   the kit scene). Low risk.
3. The `phase.BlockPlayerMovement` interpretation involves
   `ctx.ActorManager.GetPlayer()` returning a more general type than `Player`.
   Kit-side this stays type-correct because `BlockMovement()` is on
   `actors.ActorEntity`.
