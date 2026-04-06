# SPEC — 023 Projectile Manager

**Branch:** `023-projectile-manager`
**Bounded Context:** `internal/engine/combat/projectile/`

---

## Overview

Extract bullet lifecycle (spawn, update, draw, despawn) from `PhasesScene` into a standalone `projectile.Manager` that lives on `app.AppContext`. Any scene or weapon can then spawn projectiles without coupling to a specific scene.

---

## New Files

| Path | Purpose |
|---|---|
| `internal/engine/combat/projectile/config.go` | `ProjectileConfig` value type |
| `internal/engine/combat/projectile/projectile.go` | internal `projectile` struct |
| `internal/engine/combat/projectile/manager.go` | `Manager` struct — public API |

---

## Types

### `ProjectileConfig`
```go
// internal/engine/combat/projectile/config.go
package projectile

type ProjectileConfig struct {
    Width   int // pixels
    Height  int // pixels
    Damage  int
}
```

### `projectile` (internal)
```go
// internal/engine/combat/projectile/projectile.go
package projectile

import contractsbody "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"

type projectile struct {
    movable  contractsbody.Movable
    body     contractsbody.Collidable
    space    contractsbody.BodiesSpace
    speedX16 int
    speedY16 int
}
```

Behaviour mirrors the existing `gamestates.Bullet`:
- `update()` — advances position by `(speedX16, speedY16)` each frame, resolves collisions, queues removal when out-of-bounds.
- `OnTouch(other)` — queues removal unless `other` is the owner.
- `OnBlock(other)` — queues removal.

### `Manager`
```go
// internal/engine/combat/projectile/manager.go
package projectile

import (
    contractsbody "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
    "github.com/hajimehoshi/ebiten/v2"
)

type Manager struct {
    projectiles []*projectile
    space       contractsbody.BodiesSpace
    counter     int
}

func NewManager(space contractsbody.BodiesSpace) *Manager

// Spawn creates a new projectile and registers its body in the BodiesSpace.
// x16, y16 — fp16 spawn position; vx16, vy16 — fp16 velocity; owner — collision owner.
func (m *Manager) Spawn(cfg ProjectileConfig, x16, y16, vx16, vy16 int, owner interface{})

// Update advances all active projectiles and removes despawned ones.
func (m *Manager) Update()

// Draw renders all active projectiles to screen.
func (m *Manager) Draw(screen *ebiten.Image)

// Clear removes all projectiles and their bodies from the space.
func (m *Manager) Clear()
```

---

## Contract Change — `contracts/projectile/projectile.go` (new)

A new contract package is needed so other packages can depend on the interface rather than the concrete type:

```go
// internal/engine/contracts/projectile/projectile.go
package projectile

import "github.com/hajimehoshi/ebiten/v2"

type Manager interface {
    Spawn(cfg interface{}, x16, y16, vx16, vy16 int, owner interface{})
    Update()
    Draw(screen *ebiten.Image)
    Clear()
}
```

> Design note: `cfg interface{}` avoids a circular import between `contracts/projectile` and `combat/projectile`. The concrete `Manager` satisfies this interface.

---

## `app.AppContext` Change

Add one field to `internal/engine/app/context.go`:

```go
ProjectileManager *combatprojectile.Manager
// import: "github.com/boilerplate/ebiten-template/internal/engine/combat/projectile"
```

---

## `setup.go` Change

In `internal/game/app/setup.go`, after `Space` is created:

```go
appContext.ProjectileManager = projectile.NewManager(appContext.Space)
```

The `ProjectileManager` must be initialised after `Space` because it holds a reference to `BodiesSpace`.

---

## `PhasesScene` Changes

Remove from `PhasesScene`:
- Field `bullets []*gamestates.Bullet`
- Field `bulletImg *ebiten.Image`
- Field `bulletCounter int`
- Method `SpawnBullet(x16, y16, vx16, vy16 int, owner interface{})`
- Bullet update loop in `Update()`
- Bullet draw loop in `Draw()`
- Bullet cleanup loop in `Update()`

Replace with delegation:
- `SpawnBullet` → `ctx.ProjectileManager.Spawn(projectile.ProjectileConfig{Width:2, Height:1}, x16, y16, vx16, vy16, owner)`
- `Update()` → call `ctx.ProjectileManager.Update()` (replaces the three bullet loops)
- `Draw()` → call `ctx.ProjectileManager.Draw(screen)` (replaces the bullet draw loop)
- `OnFinish()` → call `ctx.ProjectileManager.Clear()`

> `PhasesScene` still satisfies `body.Shooter` via the delegating `SpawnBullet` method.

---

## Pre-conditions

- `app.AppContext.Space` (`body.BodiesSpace`) is non-nil before `ProjectileManager` is created.
- `BodiesSpace.GetTilemapDimensionsProvider()` returns a non-nil provider during `Update()` (same requirement as the existing `Bullet`).

## Post-conditions

- `Manager.Spawn()` increases `len(m.projectiles)` by 1 and adds one body to `BodiesSpace`.
- `Manager.Update()` removes projectiles whose bodies are no longer in the space.
- `Manager.Clear()` leaves `len(m.projectiles) == 0`.

---

## Integration Points

| Point | Detail |
|---|---|
| `body.BodiesSpace` | `Spawn` calls `AddBody`; out-of-bounds/collision calls `QueueForRemoval` + `ProcessRemovals` |
| `body.Shooter` (contract) | `PhasesScene.SpawnBullet` delegates to `Manager.Spawn` — interface still satisfied |
| `app.AppContext` | New `ProjectileManager` field; initialised in `setup.go` |
| `gamestates.Bullet` | Deleted after migration; `Manager`'s internal `projectile` replaces it |

---

## Red Phase — Failing Test Scenario

**File:** `internal/engine/combat/projectile/manager_test.go`

**Scenario 1 — Spawn increases count**
```
GIVEN a Manager with a mock BodiesSpace
WHEN  Spawn is called once with a valid ProjectileConfig
THEN  len(manager.projectiles) == 1
AND   mockSpace.AddBody was called once
```

**Scenario 2 — Out-of-bounds despawn**
```
GIVEN a Manager with one active projectile
AND   the mock BodiesSpace returns a tilemap of size 100×100
AND   the projectile position after Update() exceeds the tilemap bounds
WHEN  Update() is called
THEN  the projectile is removed from manager.projectiles
AND   mockSpace.QueueForRemoval was called
```

**Scenario 3 — Clear resets state**
```
GIVEN a Manager with two active projectiles
WHEN  Clear() is called
THEN  len(manager.projectiles) == 0
AND   mockSpace.RemoveBody was called for each projectile body
```

These three scenarios map directly to AC7 and must be **red** before any implementation is written.
