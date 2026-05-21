# NOTES — 062-depth-aware-collision

Human-readable companion to SPEC.md. Captures rationale, risks, and decisions deferred by the story.

## Investigation findings

- `HasCollision` lives at `internal/engine/physics/space/space.go:223`. It is a free function (not a method), called from `Space.ResolveCollisions`. The current body iterates `collisionRects(a)` × `collisionRects(b)` and returns on first overlap. Augmentation is straightforward: same prologue + bbox loop, new gate appended.
- `body.Collidable` already exposes `GetPosition16()` (story suggestion for `GroundY` source) but its semantics are screen-Y after altitude. Using it directly inside `space/` would either (a) require `space/` to know about altitude (`b.Altitude()`) to recover ground-Y, or (b) silently consume the wrong axis. Both are worse than an explicit opt-in method.
- Story 061 (`done/061-altitude-jump-ground-detection`) introduced `VAltitude16`, `Altitude`, `Altitude16` on `Body`. The bbox `Position()` is already drawn at `Y - Altitude` (screen-Y). So bbox overlap reflects what the player SEES. Lane gating is the correct way to filter "looks-overlapping-but-can't-touch".
- No existing usages of any `LaneWidth` constant in the codebase (grep clean). Greenfield for naming.

## Why interface assertion (over tag field or config)

- **Tag field on `Body`** would require modifying `internal/engine/contracts/body/body.go` (a heavy contract change touching every implementer), and would force ALL bodies to declare an opt-in flag even when irrelevant. Rejected.
- **Config-level pair check** (e.g., a `ScenePolicy` injected into Space) would require `space/` to take a new dependency and would couple `space/` to scene config. Rejected.
- **Local interface assertion** is the smallest surface: downstream `kit`/`game` types add two methods and become opt-in automatically. Discovery is via Go's structural typing — zero pressure on the contracts layer. Chosen.

## Why per-body `LaneHalfWidth()` (over a constant or scene config)

- A package-level constant in `space/` would mean every 2.5D body has the same depth tolerance — fine for a single beat-em-up but brittle if multiple genres coexist (the engine is genre-agnostic by design per the constitution).
- A scene-config injection would force `space/` to import a config struct (layer concern).
- Per-body method matches the principle that physics tuning lives with the entity. The story permits "field on a config struct, per-body method, or scene-level config" — per-body is the most decoupled.
- `DefaultLaneHalfWidth = 8` is exported as documentation for downstream implementers; `space/` itself never reads it.

## Why `GroundY()` on the interface (not `GetPosition16()`)

- The story says "GroundY is available via GetPosition16() (second return)", but `GetPosition16()` returns the body's drawn Y (screen-Y). After altitude was added in 061, screen-Y diverged from ground-Y. Calling it "GroundY" inside `space/` would mislabel the value when an airborne body is checked.
- Making `GroundY()` an explicit method on `DepthLaneBody` lets downstream beat-em-up bodies return `y + altitude` (i.e., re-add the altitude offset to recover ground depth) without exposing that arithmetic to `space/`.

## Risks

- **Risk:** Downstream `kit`/`game` bodies forget to implement `DepthLaneBody` and silently fall back to the 2D path.
  **Mitigation:** Story 063+ (or a follow-up) should add an `IsDepthAware()` smoke test on the player/enemy beat-em-up bodies. Out of scope here.
- **Risk:** Per-body `LaneHalfWidth` values diverge between player and enemy, causing surprising "max wins" behavior.
  **Mitigation:** Documented in SPEC §4. Recommend downstream uses a single shared constant (e.g., from the beat-em-up scene module) when constructing bodies.
- **Risk:** `GroundY()` implementation drifts from the body's actual depth (e.g., not updated when the body moves).
  **Mitigation:** Out of scope for `space/`; this is a downstream contract. A future story should formalize a `BeatEmUpActor` interface in `kit/` that wires `GroundY()` directly to `y + altitude`.

## Out of scope

- Modifying `body.Collidable` or any contract under `internal/engine/contracts/`.
- Mock generation: no new contracts means no `internal/engine/mocks/` work.
- Downstream wiring of `DepthLaneBody` onto concrete `kit`/`game` bodies (separate story).
- Scene-level depth-lane configuration (e.g., per-phase lane width override).
- Visual debug overlay for lane bounds.

## Pipeline note

The Spec Engineer has decided NOT to introduce new contracts. Mock Generator can be skipped. Sequence: TDD Specialist → Feature Implementer → Gatekeeper.
