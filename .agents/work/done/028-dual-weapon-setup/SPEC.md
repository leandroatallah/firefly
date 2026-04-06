# Technical Specification — 028 Dual Weapon Setup

**Branch:** `028-dual-weapon-setup`  
**Bounded Context:** Game (`internal/game/`)

## Overview

Wire two concrete weapons into the climber player's inventory at startup. Both weapon definitions live in the game layer. The engine's `weapon.NewProjectileWeapon`, `inventory.New`, and `projectile.Manager` are already in place — this story only adds the game-specific wiring.

---

## New File

**`internal/game/entity/actors/player/weapons.go`**

```go
package gameplayer

import (
    "github.com/boilerplate/ebiten-template/internal/engine/combat/inventory"
    "github.com/boilerplate/ebiten-template/internal/engine/combat/weapon"
    "github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
)

func newClimberInventory(manager combat.ProjectileManager) *inventory.Inventory {
    inv := inventory.New()
    inv.AddWeapon(weapon.NewProjectileWeapon("light_blaster", 8,  "bullet_small", 393216, manager))
    inv.AddWeapon(weapon.NewProjectileWeapon("heavy_cannon",  30, "bullet_large", 589824, manager))
    return inv
}
```

> Speeds in fp16: `393216 = 6 * 65536`, `589824 = 9 * 65536`.

---

## Changed File

**`internal/game/scenes/phases/player.go`** — pass the inventory into `SkillDeps`:

```go
deps := skill.SkillDeps{
    Inventory:         gameplayer.NewClimberInventory(ctx.ProjectileManager),
    ProjectileManager: ctx.ProjectileManager,
    OnJump:            ...,
    EventManager:      ctx.EventManager,
}
```

> `newClimberInventory` must be exported as `NewClimberInventory` so `player.go` can call it.

---

## Pre-conditions

- US-022 through US-027 complete.
- `ctx.ProjectileManager` is non-nil in `AppContext` (set in `setup.go`).
- `skill.FromConfig` skips the shooting skill when `deps.Inventory == nil`; passing a real inventory enables it.

## Post-conditions

- `ShootingSkill` receives a populated inventory with two weapons.
- Active weapon starts as `light_blaster` (index 0).
- `WeaponNext` switches to `heavy_cannon`; `WeaponPrev` wraps back.
- No engine packages are modified.

---

## Projectile Size Mapping

The `projectileType` string passed to `SpawnProjectile` is already used by `projectile.Manager` to look up a `ProjectileConfig`. A follow-up story can map `"bullet_small"` → `{Width:2, Height:1}` and `"bullet_large"` → `{Width:4, Height:3}` inside the manager. For this story the type string is sufficient — the manager falls back to a default config if the type is unknown.

---

## Red Phase: Failing Test Scenario

**File:** `internal/game/entity/actors/player/weapons_test.go`

```
GIVEN  a mock ProjectileManager
WHEN   NewClimberInventory(mockManager) is called
THEN   inv.ActiveWeapon().ID() == "light_blaster"
AND    inv.ActiveWeapon().CanFire() == true

WHEN   inv.SwitchNext() is called
THEN   inv.ActiveWeapon().ID() == "heavy_cannon"
```

Expected failure: `NewClimberInventory` does not exist yet.
