# 026 — Enhanced Shooting Skill (Inventory-Aware)

**Branch:** `026-shooting-skill-inventory`
**Bounded Context:** Engine (`internal/engine/physics/skill/`)

## Story

As a game developer using this boilerplate,
I want `ShootingSkill` to fire through the active weapon in an inventory rather than calling a scene shooter directly,
so that weapon switching and multiple weapon types work without scene coupling.

## Context

`ShootingSkill` currently holds a `body.Shooter` reference (a scene interface) and constructs bullet parameters inline. After US-022 (Weapon System), US-023 (Projectile Manager), and US-025 (Inventory), the skill should delegate firing to `inventory.ActiveWeapon()` and weapon switching to the inventory via the existing `input.CommandsReader()` (US-014).

## Acceptance Criteria

- **AC1** — `ShootingSkill` constructor signature changed to `NewShootingSkill(inv *inventory.Inventory)`.
- **AC2** — `HandleInput()` calls `inv.ActiveWeapon().Fire()` when shoot is pressed and `CanFire()` is true.
- **AC3** — `HandleInput()` calls `inv.SwitchNext()` / `inv.SwitchPrev()` based on `input.CommandsReader().WeaponNext` / `WeaponPrev`.
- **AC4** — Old `NewShootingSkill(shooter, cooldown, offset, speed, yOffset)` constructor is removed.
- **AC5** — `body.Shooter` interface and its mock are removed (no longer needed).
- **AC6** — Unit tests cover: fire delegates to active weapon, weapon switch on input, no fire when `CanFire()` is false.

## Notes

- Depends on US-022, US-023, US-025.
- Direction detection logic (8-directional from US-012) is preserved unchanged.
- `ShootingSkill` no longer owns cooldown state — that moves to `Weapon`.
