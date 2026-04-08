# US-025 — Weapon Inventory System

**Branch:** `025-weapon-inventory`
**Bounded Context:** Engine (`internal/engine/combat/inventory/`)

## Story

As a game developer using this boilerplate,
I want a weapon inventory that holds multiple weapons and supports switching between them,
so that I can implement Megaman-style weapon selection or Metal Slug weapon pickups.

## Context

With the weapon system (US-022) in place, characters need a way to hold multiple weapons and switch between them. The inventory also manages per-weapon ammo, enabling energy-limited weapons like Megaman's special weapons.

## Acceptance Criteria

- **AC1** — `Inventory` struct with `AddWeapon()`, `ActiveWeapon()`, `SwitchNext()`, `SwitchPrev()`, `SwitchTo(index int)`.
- **AC2** — `Inventory` tracks ammo per weapon ID; `HasAmmo()` and `ConsumeAmmo()` methods provided.
- **AC3** — Unlimited ammo is represented as `-1`; `HasAmmo()` returns `true` for `-1`.
- **AC4** — `SwitchNext()` / `SwitchPrev()` wrap around (last → first, first → last).
- **AC5** — `ActiveWeapon()` returns `nil` when inventory is empty (no panic).
- **AC6** — Unit tests cover: add/switch/wrap-around, ammo consumption, empty inventory guard.

## Notes

- Input handling for weapon switching is wired in US-026 (Enhanced Shooting Skill).
- Package path: `internal/engine/combat/inventory/`.
- Ammo map key is weapon `ID()` string.
