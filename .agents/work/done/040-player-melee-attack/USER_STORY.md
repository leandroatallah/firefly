# US-040 — Player Melee Attack

**Branch:** `040-player-melee-attack`
**Bounded Context:** Engine (`internal/engine/combat/weapon/`) + Game (`internal/game/entity/actors/states/`)

## Story

As a player,
I want to perform a close-range melee attack with the Z key,
so that I can engage enemies without relying on projectiles.

## Context

US-022 established the `Weapon` interface and `ProjectileWeapon`.
US-010/US-026 established `ShootingSkill` as the model for active combat skills.
US-038 established faction-aware damage via the `Damageable` contract.
`PlayerCommands.Melee` (mapped to `KeyZ`) was added as setup work for this story.

Currently the player can only shoot. No close-range combat option exists. This story introduces a swing-hitbox melee attack: a new `MeleeWeapon`, a new `StateMeleeAttack` Actor state, and a per-state hitbox that is active only during a configurable frame window within the attack animation.

## Acceptance Criteria

- **AC1** — `MeleeWeapon` struct in `internal/engine/combat/weapon/` implements `combat.Weapon`. Fields: `damage`, `cooldownFrames`, `hitboxWidth16`, `hitboxHeight16`, `hitboxOffsetX16`, `hitboxOffsetY16`, `activeFrames` (frame window during which the hitbox deals damage).
- **AC2** — New Actor state `StateMeleeAttack` registered in `internal/game/entity/actors/states/`. Transitions back to the originating state (Grounded or Falling) when `IsAnimationFinished()` returns true.
- **AC3** — During `StateMeleeAttack`, the swing hitbox is live only within the configured `active_frames` window. Outside that window the hitbox is inactive (no damage).
- **AC4** — The hitbox applies `damage` to all overlapping `Damageable` actors with a different faction (faction system from US-038). Same-faction actors are never damaged.
- **AC5** — Melee can be triggered from both `GroundedState` and `Falling` state (air melee allowed). It cannot be triggered while dashing or while already in `StateMeleeAttack`.
- **AC6** — `MeleeWeapon` is JSON-configurable under a `"type": "melee"` key:
  ```json
  {
    "type": "melee",
    "damage": 1,
    "cooldown_frames": 20,
    "active_frames": [4, 10],
    "hitbox": {
      "width": 24,
      "height": 16,
      "offset_x": 12,
      "offset_y": 0
    }
  }
  ```
- **AC7** — `weapon.Factory` is extended to instantiate `MeleeWeapon` from configs with `"type": "melee"`.
- **AC8** — Unit tests verify:
  - Hitbox is inactive before and after `active_frames` window.
  - Damage is applied to an overlapping `Damageable` enemy during the active window.
  - Damage is NOT applied to the player (same faction).
  - Cooldown prevents re-triggering during recovery.
  - `StateMeleeAttack` transitions back correctly after `IsAnimationFinished()`.
  - Air melee transitions back to `Falling` (not `Grounded`).

## Proposed Changes

- `internal/engine/combat/weapon/melee.go` — `MeleeWeapon` struct and logic
- `internal/engine/combat/weapon/factory.go` — extend to handle `"type": "melee"`
- `internal/game/entity/actors/states/melee_state.go` — `StateMeleeAttack` state
- Player JSON config — add `melee` weapon section

## Dependencies

- US-038 — Projectile Damage on Hit (Damageable contract + faction system)
