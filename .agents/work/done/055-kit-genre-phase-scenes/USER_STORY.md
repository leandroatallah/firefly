# USER STORY — 055-kit-genre-phase-scenes

**Branch:** `055-kit-genre-phase-scenes`

**Bounded Context:** Kit / Scene

---

## Story

As a developer building a game on the Firefly boilerplate,
I want genre-specific Phase Scene implementations (`PlatformerPhaseScene`, `BeatemupPhaseScene`) to live in `internal/kit/scenes/phases/`,
so that I can assemble a new game by composing a ready-made kit scene instead of copying and modifying `internal/game/` logic that is presently hard-coupled to the platformer genre.

---

## Background

`internal/game/scenes/phases/scene.go` — the `PhasesScene` type — is the runtime Phase Scene for the example game. It violates the three-layer architecture rule (see Constitution §Architecture) because it contains reusable, genre-specific boilerplate that belongs in the kit layer, not the game layer. Concrete violations on `main` as of story creation:

| Location | Coupling point |
|---|---|
| Field `player platformer.PlatformerActorEntity` (line 51) | Typed to platformer kit interface |
| Field `screenFlipper *scene.ScreenFlipper` (line 68) | Room-flip mechanic — platformer-only |
| Method `checkPlayerFallDeath` (line 264) | Gravity/fall-out-of-camera death — platformer-only |
| Field `gameCamera *gamecamera.Controller` (line 74) | Wraps engine camera with vertical-only-upward constraint — platformer constraint |
| Body-iteration type switch `case platformer.PlatformerActorEntity` (line 413) | Forces platformer interface on all actor iteration |
| Draw type switch `case platformer.PlatformerActorEntity` (line 476) | Cross-genre type switch in draw path |
| `cp.MeleeController()` cast to `*gameplayer.ClimberPlayer` (line 518) | Game-specific player type leaking into scene Draw |

A sibling actor package `internal/kit/actors/beatemup/` now exists alongside `internal/kit/actors/platformer/`. Beat-em-up scenes require an altitude axis, depth-sorted draw (`draworder.SortByGroundY` on the Y+altitude composite), an arena-scrolling camera (no vertical-only-upward lock), and no fall-death.

The project is a boilerplate, not a shipped game. Per the Constitution, the engine must remain genre-agnostic, the kit must contain opinionated genre implementations, and `internal/game/` must hold only thin, game-specific assembly.

---

## Constitution References

- **§Architecture**: "three-layer — engine, kit, game; dependency rules: engine must not import kit or game; kit may import engine, must not import game; game may import both engine and kit."
- **§Bounded Contexts**: Kit maps to `internal/kit/`; Scene maps to `internal/engine/scene/`.
- **§Non-Negotiable Standards / Production Code**: "No global mutable state. Inject dependencies via interfaces (`internal/engine/contracts/`)."

---

## Acceptance Criteria

### AC-1 — Engine retains only genre-agnostic phase-scene base

The engine package `internal/engine/scene/` (and any sub-package such as `internal/engine/scene/phases/`) must contain:
- A base phase-scene composition type (e.g., an embeddable struct or interface) that provides: pause lifecycle, goal system integration, sequence-player wiring, vignette, body-iteration scaffolding (item removal, obstacle skip, projectile update), and a common `Update`/`Draw` skeleton via hooks or template-method pattern.
- No import of `internal/kit/`, `internal/game/`, or any genre-specific interface such as `platformer.PlatformerActorEntity`.

### AC-2 — Kit gains `PlatformerPhaseScene`

A new type `PlatformerPhaseScene` in `internal/kit/scenes/phases/platformer/` must:
- Embed the engine base phase-scene (AC-1).
- Hold `player platformer.PlatformerActorEntity` (typed to the kit platformer interface, not a game type).
- Wire `scene.ScreenFlipper` for room-flipping on player push.
- Apply a vertical-only-upward camera constraint (the logic currently in `internal/game/render/camera`) — either by moving the constraint into a kit camera decorator at `internal/kit/render/camera/` or by configuring the engine camera with a constraint option.
- Implement `checkPlayerFallDeath`: if the player's top Y exceeds camera bottom, trigger the death sequence.
- Body iteration and Draw must use a `platformer.PlatformerActorEntity` type switch scoped exclusively to this scene; no reference to beat-em-up types.

### AC-3 — Kit gains `BeatemupPhaseScene`

A new type `BeatemupPhaseScene` in `internal/kit/scenes/phases/beatemup/` must:
- Embed the engine base phase-scene (AC-1).
- Hold an actor typed to the beat-em-up kit interface (a `BeatEmUpActorEntity` interface to be defined or extended from `internal/kit/actors/beatemup/`).
- Apply an arena-scrolling camera with no vertical-only-upward lock.
- Sort draw calls by ground-Y + altitude composite (depth sort appropriate for 2.5D).
- Not implement `checkPlayerFallDeath` — fall-death is not a beat-em-up concept.
- Body iteration and Draw must not reference `platformer.PlatformerActorEntity`.

### AC-4 — `internal/game/scenes/phases/` shrinks to thin wiring

After the refactor, `internal/game/scenes/phases/` must contain only game-specific assembly:
- Wiring `PlatformerPhaseScene` (or `BeatemupPhaseScene`) with the concrete game player type (`ClimberPlayer`), game-specific item map, enemy map, NPC map, and scene-navigation identifiers.
- No genre logic (no `checkPlayerFallDeath`, no `ScreenFlipper` instantiation, no camera constraint logic) duplicated from the kit layer.
- The file may be removed entirely if the genre kit scenes are self-sufficient for the boilerplate examples, in which case `internal/game/app/` wires the kit scene directly.

### AC-5 — Cross-genre type switches eliminated

No file outside `internal/kit/scenes/phases/platformer/` may contain a type switch or type assertion against `platformer.PlatformerActorEntity`. No file outside `internal/kit/scenes/phases/beatemup/` may contain a type switch against a beat-em-up actor type. Engine body iteration (AC-1) must use only engine-level contracts (`body.Body`, `body.Obstacle`, `items.Item`).

### AC-6 — Melee-hitbox debug draw relocated via extension hook

The `*gameplayer.ClimberPlayer` cast in `Draw` (currently line 518 of `scene.go`) that draws the active melee hitbox is a game-specific debug concern. It must not appear in any kit scene. The kit `PlatformerPhaseScene` must expose a `SetDebugDrawHook(func(screen *ebiten.Image))` (or equivalent extension point) that `internal/game/` calls to inject game-specific debug rendering (e.g., melee hitbox overlay) without the kit importing any game type.

### AC-7 — Existing tests pass; new kit-level tests added

- All tests in `internal/engine/` and `internal/game/` must pass after the refactor with no regressions.
- `internal/kit/scenes/phases/platformer/` must have unit tests covering:
  - `checkPlayerFallDeath` fires death sequence when player top Y exceeds camera bottom.
  - `checkPlayerFallDeath` does not fire when death sequence is already active.
  - `ScreenFlipper` callbacks set/unset player immobility on flip start/finish.
- `internal/kit/scenes/phases/beatemup/` must have unit tests covering:
  - Body-iteration draw order is sorted by ground-Y + altitude (no fall-death path exists).
  - Actor removal on `Dead` state does not panic when actor has no altitude component.
- Tests must follow table-driven patterns (Constitution §Tests) and must not call `ebiten.RunGame` or use `time.Sleep`.

### AC-8 — Layer import rules enforced

CI (or `go vet`/build tags) must continue to enforce:
- `internal/kit/scenes/phases/platformer/` does not import `internal/game/`.
- `internal/kit/scenes/phases/beatemup/` does not import `internal/game/`.
- `internal/engine/scene/` does not import `internal/kit/`.

---

## Behavioral Edge Cases

1. **Phase with no player** — both kit scenes must handle `hasPlayer = false` gracefully (fixed camera, no flipper, no fall-death check) without panicking.
2. **Simultaneous death trigger sources** — fall-death and state-machine death (`Dying`/`Dead` state) must both route through a single `startDeathSequence` guard so the sequence fires at most once per phase.
3. **Sequence-gated pause** — when a phase uses `SequenceGoalType`, the pause key must be suppressed; both kit scenes inherit this from the engine base (AC-1) without re-implementing it.
4. **ScreenFlipper during death** — if a flip is in progress when death triggers, the flip must be cancelled or completed before navigation, not left in an intermediate state.
5. **Camera constraint direction** — the vertical-only-upward constraint means the camera Y can increase (player moves down) but never snaps upward past the player; this must be tested as an isolated unit on the camera decorator, not only through the full scene.

---

## Out of Scope

- Redesigning the goal system (`phases.Goal`, `ReachEndpointGoal`, `SequenceGoal`) — goals are engine primitives and remain unchanged.
- Rewriting engine camera internals (`internal/engine/render/camera`) — the constraint lives in a kit-layer decorator, not the engine camera itself.
- Implementing beat-em-up gameplay mechanics (altitude gravity, jump, ground detection) — those are covered by story `053-altitude-engine-foundation` and future stories; `BeatemupPhaseScene` only scaffolds the scene shell with correct draw order and camera.
- Migrating audio, i18n, or VFX systems.
- Changing the scene navigation identifiers in `internal/game/scenes/types/`.

---

## Open Questions / Risks

1. **Where does `gameplayer.ClimberPlayer` melee-hitbox debug draw belong?**
   Currently at line 518 of `scene.go` inside `Draw`, gated by `config.Get().CollisionBox`. This is a game-specific debug concern (it casts to a concrete game type). Proposed resolution: the kit `PlatformerPhaseScene` exposes a `SetDebugDrawHook(func(*ebiten.Image))` called at the end of `Draw` after all kit-managed rendering. `internal/game/` registers a closure that captures `*gameplayer.ClimberPlayer` and draws the hitbox. This keeps the kit free of game imports while preserving the debug feature (AC-6).

2. **Vertical-only-upward camera constraint location** — `internal/game/render/camera/camera.go` is a thin delegating wrapper; the actual vertical-only-upward constraint logic does not appear in this file (it may be in the engine camera's follow logic). The Spec Engineer must locate where this constraint is implemented before deciding whether it moves to `internal/kit/render/camera/` or is expressed as a configuration option on the engine camera.

3. **`BeatEmUpActorEntity` interface** — `internal/kit/actors/beatemup/` currently exposes only `BeatEmUpCharacter` (a concrete struct embedding `MeleeCharacter`). An interface (`BeatEmUpActorEntity`) analogous to `PlatformerActorEntity` does not yet exist. The Spec Engineer should define this interface in `internal/kit/actors/beatemup/` before speccing `BeatemupPhaseScene`.

4. **Altitude draw sort** — `draworder.SortByGroundY` currently sorts by body ground-Y alone. For a beat-em-up, depth order is ground-Y plus altitude. Whether this requires a new `draworder.SortByGroundYAltitude` function in the engine or a local sort in `BeatemupPhaseScene` should be decided in the spec.

5. **Risk: test coverage regression** — `internal/game/scenes/phases` currently has 18.4% coverage (AGENTS.md §Priorities). Migrating logic to the kit layer will reduce the lines counted under `internal/game/`, potentially lowering or eliminating coverage in that sub-path. The Gatekeeper must confirm the 80%+ coverage goal is met across the new kit paths, not just the residual game path.
