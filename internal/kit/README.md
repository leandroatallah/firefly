# Kit Module

The `kit` layer provides genre-reusable concrete implementations for platformer/brawler games built on `internal/engine`. It sits between the engine abstractions and the project-specific `internal/game` layer.

**Dependency rule:** `kit` may import `internal/engine/...`; it must never import `internal/game/...`.

## Packages

- `actors/`: Composable character trait structs (`MeleeCharacter`, `ShooterCharacter`, `DeathBehavior`). Embed independently — a brawler character can use both `MeleeCharacter` and `ShooterCharacter`.
- `combat/`: Weapon inventory, projectile lifecycle, melee controller, and faction-gated damage. See [`combat/README.md`](combat/README.md).
  - `inventory/`: Weapon collections and ammo tracking.
  - `weapon/`: `ProjectileWeapon`, `EnemyShooting`, and a JSON weapon factory.
  - `projectile/`: High-performance projectile manager with lifetime, VFX, and damage hooks.
  - `melee/`: `Controller` + `State` for per-actor melee swings (input buffering, combo, hitbox, VFX).
- `skills/`: Physics-linked actor abilities (`JumpSkill`, `DashSkill`, `HorizontalMovementSkill`, `ShootingSkill`) plus a JSON `FromConfig` factory. Engine-level contracts (`Skill`, `ActiveSkill`, `SkillBase`) live in `internal/engine/skill/`.
- `states/`: Genre-reusable `ActorState` implementations (e.g., `MeleeState`). Parameterised on the caller's enum to avoid coupling to a specific game's state vocabulary.
- `ui/speech/`: `speech.Manager` — the dialogue orchestrator (typing flow, audio scheduling, skip behaviour). Implements `contracts/dialogue.Manager`. Speech primitives live in `internal/engine/ui/speech/`.

## Placement rule

A component belongs in `kit` when 80%+ of games in this genre would use it as-is. If it is specific to *this* game's art, levels, or rules, it belongs in `internal/game/`. See [ADR-006](../../docs/adr/ADR-006-engine-game-layer-separation.md).
