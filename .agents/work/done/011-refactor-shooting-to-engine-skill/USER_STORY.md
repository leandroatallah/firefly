# US-011 — Refactor Shooting to Engine Skill

**Branch:** `011-refactor-shooting-to-engine-skill`
**Bounded Context:** Physics (`internal/engine/physics/skill/`)

## Story

As a developer,
I want shooting to be represented as explicit actor states (IdleShooting, WalkingShooting, etc.),
so that each visual state has its own sprite sheet, animation timing, and clear state machine transitions.

## Context

Currently `ShootingSkill` lives in `internal/game/entity/actors/states/` and is tightly coupled to `GroundedState`. It doesn't follow the engine's state-based architecture.

In Cuphead, shooting changes the visual state:
- Idle → Idle Shooting (different sprite, not just an overlay)
- Walking → Walking Shooting
- Jumping → Jumping Shooting
- Each has distinct animation timing and sprite sheets

This refactor:
1. Registers shooting state variants as first-class actor states
2. Moves `ShootingSkill` to the engine layer as an `ActiveSkill` that triggers state transitions
3. Enables the sprite system to map shooting states to distinct sprite sheets
4. Supports future directional variants (shooting up, diagonal, etc.)

## Acceptance Criteria

- **AC1** — Shooting state variants are registered in `actor_state.go`: `IdleShooting`, `WalkingShooting`, `JumpingShooting`, `FallingShooting`.
- **AC2** — `ShootingSkill` is moved from `internal/game/entity/actors/states/` to `internal/engine/physics/skill/skill_shooting.go`.
- **AC3** — `ShootingSkill` implements the `ActiveSkill` interface: `HandleInput()`, `Update()`, `IsActive()`, `ActivationKey()`.
- **AC4** — `HandleInput()` triggers state transitions: pressing shoot → transition to shooting state variant; releasing shoot → transition back to base state.
- **AC5** — Bullet spawning occurs in `HandleInput()` when cooldown allows, with alternating Y-offset via `OffsetToggler`.
- **AC6** — `Character.handleState()` includes transition logic for shooting states (e.g., IdleShooting → WalkingShooting when moving).
- **AC7** — `OffsetToggler` is moved to `internal/engine/physics/skill/offset_toggler.go`.
- **AC8** — `Bullet` entity is moved to `internal/engine/entity/projectiles/bullet.go`.
- **AC9** — `GroundedState` no longer contains shooting-specific logic; `ShootHeld()` removed from `GroundedInput`.
- **AC10** — Sprite system can map shooting states to distinct sprite sheets (e.g., `"idle_shoot"` → `idle_shoot.png`).
- **AC11** — All existing shooting tests pass after refactor; no behavioral changes to bullet spawning or cooldown.
- **AC12** — Code coverage remains ≥74.6% (no regression).

## Behavioral Edge Cases

- Holding shoot button continuously must spawn bullets at the configured fire rate (cooldown enforcement).
- Releasing and re-pressing shoot within cooldown must not spawn extra bullets.
- Shooting while transitioning between grounded sub-states (idle ↔ walking ↔ aim-lock) must not reset cooldown.
- Shooting must be suppressed while dashing (handled by skill priority or state machine).

## Notes

- This is a **refactor**, not a feature addition. Bullet spawning behavior must remain identical to US-010.
- `OffsetToggler` is moved from `internal/game/entity/actors/states/` to `internal/engine/physics/skill/`.
- The `Shooter` contract (`internal/engine/contracts/body/shooter.go`) remains unchanged.
- Shooting states are **explicit actor states**, not skill modifiers. This matches Cuphead's architecture where "idle shooting" is a distinct visual state.
- Each shooting state variant can have its own sprite sheet, animation timing, and hitboxes.
- Future directional variants (IdleShootingUp, WalkingShootingDiagonal, etc.) can be added as new states.
- The `ShootingSkill` triggers state transitions but doesn't replace the state machine — it works alongside it.
- Sprite mapping: `IdleShooting` → `"idle_shoot"` → `idle_shoot.png` (handled by existing sprite system).

## Success Criteria

- All tests pass (no regressions).
- Shooting states are registered and functional.
- `ShootingSkill` follows the `ActiveSkill` pattern.
- Sprite system can map shooting states to distinct sprite sheets.
- Code coverage remains ≥74.6% (no delta loss).
- `internal/game/entity/actors/states/` no longer contains shooting-specific logic.
