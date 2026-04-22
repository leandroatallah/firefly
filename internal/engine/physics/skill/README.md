# Physics Skills

The `physics/skill` package encapsulates discrete player abilities that interact with the platformer movement model. Skills are **composable units**: an actor holds a `[]Skill`, updates them each frame, and optionally forwards input to those implementing `ActiveSkill`.

## Interfaces

- `Skill` — `Update(actor, model)` + `IsActive() bool`. Required for every skill.
- `ActiveSkill` — extends `Skill` with `HandleInput(actor, model, space)` and `ActivationKey()`. Implemented by skills that react to a key press.
- `SkillBase` — a small embeddable struct providing the `StateReady / StateActive / StateCooldown` state machine and a frame timer.

## Built-in Skills

| Skill | File | Responsibility |
|---|---|---|
| `HorizontalMovementSkill` | `skill_platform_move.go` | Left/right movement with acceleration, inertia, and last-pressed-wins axis via `input.HorizontalAxis`. |
| `JumpSkill` | `skill_platform_jump.go` | Coyote time, jump buffering, and variable jump height via `JumpCutMultiplier` (clamped to `(0.1, 1.0]`). |
| `DashSkill` | `skill_dash.go` | Tween-based dash using `InOutSine`. One air dash per jump; suspends gravity; rect swap handled by `StateCollisionManager`. |
| `ShootingSkill` | `skill_shooting.go` | Consumes `Shoot`/`WeaponNext`/`WeaponPrev` commands, detects fire direction (up/down/diagonal/straight) from input + grounded/ducking state, and fires via `combat.Inventory`. Reports `IsActive()` while shoot is held. |

`skill_shooting_eight_directions.go` holds the direction-table lookup used by `ShootingSkill.detectShootDirection`.

## JSON-Driven Factory

`factory.go` exposes `FromConfig(cfg *schemas.SkillsConfig, deps SkillDeps) []Skill`. Each skill's config can enable/disable it and override tunables:

```go
skills := skill.FromConfig(cfg.Skills, skill.SkillDeps{
    Inventory:         playerInventory,
    ProjectileManager: projectileManager,
    OnJump:            func(b interface{}) { vfxManager.SpawnJumpPuff(...) },
    EventManager:      eventManager,
})
```

- `ShootingSkill` is skipped when `deps.Inventory` is nil.
- Disabled skills (`enabled: false`) are filtered out entirely.
- Defaults are applied inline (see `NewDashSkill`, `NewJumpSkill`) and overridden only when the config provides a positive value.

## State Integration

Skills **do not** implement `actors.StateContributor` directly; the game layer wraps them in small adapters that translate `skill.IsActive()` into the right game-specific actor state. See [ADR-008](../../../../docs/adr/ADR-008-state-contributor-pattern.md) and `internal/game/entity/actors/player/state_contributors.go`.

## Helpers

- `OffsetToggler` (`offset_toggler.go`) — alternates a spawn Y-offset between `-N` and `+N` on each `Next()` call. Used by the player's double-barrel visual effect so consecutive bullets spawn from alternating barrels.

## Testing

Skills accept injectable sources (`input.CommandsReader`, `physicsmovement.PlatformMovementModel`, `body.MovableCollidable`) so they can be driven deterministically from tests without real input or physics. See `skill_platform_jump_test.go` and `skill_shooting_test.go` for table-driven patterns.
