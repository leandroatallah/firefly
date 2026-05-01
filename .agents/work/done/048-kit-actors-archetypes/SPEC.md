# SPEC — 048-kit-actors-archetypes

## Overview

Introduce `internal/kit/actors/` as the first actor-level package in the kit layer. Three deliverables:

1. **Relocate** `PlayerDeathBehavior` from `internal/game/entity/actors/methods/` → `internal/kit/actors/`.
2. **Create** `ShooterCharacter` — a reusable trait that holds an `EnemyShooter` and provides update logic.
3. **Create** `MeleeCharacter` — a reusable trait that holds a melee `Controller` and provides accessor methods.

Game-layer concrete types (`BatEnemy`, `WolfEnemy`, `ClimberPlayer`) are updated to embed the new kit traits. No behaviour changes. Traits are independently composable: a future enemy can embed both `ShooterCharacter` and `MeleeCharacter` without requiring a new kit type.

---

## Dependency Rule (non-negotiable)

```
engine  →  (nothing from kit or game)
kit     →  engine only
game    →  engine + kit
```

CI enforcement (must remain green):
```bash
go list -deps ./internal/engine/... | grep -E 'internal/(kit|game)' && exit 1
go list -deps ./internal/kit/...    | grep    'internal/game'        && exit 1
```

---

## Package Layout

```
internal/kit/actors/
├── doc.go                  # package declaration + dependency note
├── death_behavior.go       # PlayerDeathBehavior (moved from game layer)
├── shooter_character.go    # ShooterCharacter trait
└── melee_character.go      # MeleeCharacter trait
```

The original file `internal/game/entity/actors/methods/death_behavior.go` is deleted. The `internal/game/entity/actors/methods/` package may be removed entirely if it becomes empty.

---

## Type Contracts

### `doc.go`

```go
// Package kitactors provides genre-reusable character trait components
// for platformer games built on the Firefly engine.
//
// Traits are independently composable: a concrete game character can embed
// any combination (e.g., ShooterCharacter + MeleeCharacter for a brawler).
//
// Dependency rule (enforced by CI):
//   - kitactors MAY import internal/engine/...
//   - kitactors MUST NOT import internal/game/...
package kitactors
```

---

### `PlayerDeathBehavior`

**File:** `internal/kit/actors/death_behavior.go`  
**Package:** `kitactors`

```go
type PlayerDeathBehavior struct {
    player platformer.PlatformerActorEntity
}

func NewPlayerDeathBehavior(p platformer.PlatformerActorEntity) *PlayerDeathBehavior
func (tm *PlayerDeathBehavior) OnDie()
```

**Behaviour:**
- `OnDie()` calls `player.SetHealth(0)`.
- Identical to the current implementation; only the package path changes.

**Imports (allowed):**
- `internal/engine/entity/actors/platformer`

---

### `ShooterCharacter`

**File:** `internal/kit/actors/shooter_character.go`  
**Package:** `kitactors`

```go
type ShooterCharacter struct {
    shooter combat.EnemyShooter
}

func NewShooterCharacter(shooter combat.EnemyShooter) *ShooterCharacter
func (s *ShooterCharacter) Shooter() combat.EnemyShooter
func (s *ShooterCharacter) SetShooter(shooter combat.EnemyShooter)
func (s *ShooterCharacter) UpdateShooter()
```

**Behaviour:**
- `NewShooterCharacter(shooter)` stores the shooter (may be nil).
- `Shooter()` returns the `EnemyShooter` (may be nil).
- `SetShooter(shooter)` assigns the shooter field.
- `UpdateShooter()` calls `shooter.Update()` if `shooter != nil`; otherwise no-op.
- No `GetCharacter`, no `Update(space)`, no `SetTarget` — those remain on the concrete game type or are promoted from `PlatformerCharacter`.

**Imports (allowed):**
- `internal/engine/contracts/combat`

**Design rationale:**
- `SetTarget` is **not** included because it conflates two concerns: movement AI (`Character.MovementState().SetTarget`) and combat AI (`shooter.SetTarget`). The game layer's concrete type (e.g., `BatEnemy`) implements `SetTarget` to coordinate both.
- `UpdateShooter()` is a focused helper that only ticks the shooter's cooldown and firing logic. The game layer's `Update(space)` calls `UpdateShooter()` then `Character.Update(space)`.

---

### `MeleeCharacter`

**File:** `internal/kit/actors/melee_character.go`  
**Package:** `kitactors`

```go
type MeleeCharacter struct {
    melee *meleeengine.Controller
}

func NewMeleeCharacter() *MeleeCharacter
func (m *MeleeCharacter) MeleeController() *meleeengine.Controller
func (m *MeleeCharacter) SetMeleeController(c *meleeengine.Controller)
```

**Behaviour:**
- `NewMeleeCharacter()` initializes with `melee = nil`.
- `MeleeController()` returns the melee controller (may be nil).
- `SetMeleeController(c)` assigns the controller field. Called by the game layer after wiring.
- No input handling, no inventory, no state-name string literals, no `Update` logic.

**Imports (allowed):**
- `internal/engine/combat/melee`

**Design rationale:**
- The melee controller's `Tick`, `HandleInput`, and `EnterAttackState` calls remain in the game layer's `Update` method because they require game-specific context (input commands, grounded state, ducking state).
- This trait only owns the field and accessor; the game layer owns the wiring and lifecycle.

---

## Game-Layer Updates

### `BatEnemy` and `WolfEnemy`

Both types replace their `shooter combat.EnemyShooter` field with an embedded `*kitactors.ShooterCharacter`:

```go
// Before
type BatEnemy struct {
    *platformer.PlatformerCharacter
    shooter combat.EnemyShooter
}

func (e *BatEnemy) Update(space body.BodiesSpace) error {
    if e.shooter != nil {
        e.shooter.Update()
    }
    return e.Character.Update(space)
}

// After
type BatEnemy struct {
    *platformer.PlatformerCharacter
    *kitactors.ShooterCharacter
}

func (e *BatEnemy) Update(space body.BodiesSpace) error {
    e.ShooterCharacter.UpdateShooter()
    return e.Character.Update(space)
}
```

- The constructor calls `kitactors.NewShooterCharacter(shooter)` after building the shooter via `builder.ConfigureEnemyWeapon`.
- `Shooter()` is promoted from the embedded trait; the game type no longer re-declares it.
- `SetTarget` remains on the concrete game type (coordinates movement + shooter targeting):
  ```go
  func (e *BatEnemy) SetTarget(target body.MovableCollidable) {
      e.Character.MovementState().SetTarget(target)
      if e.ShooterCharacter.Shooter() != nil {
          e.ShooterCharacter.Shooter().SetTarget(target)
      }
  }
  ```
- `OnTouch`, `OnDie`, and `IsEnemy` remain on the concrete game type (game-specific logic).

### `ClimberPlayer`

Replaces the `melee *meleeengine.Controller` field with an embedded `*kitactors.MeleeCharacter`:

```go
// Before
type ClimberPlayer struct {
    *platformer.PlatformerCharacter
    melee *meleeengine.Controller
    *gameplayermethods.PlayerDeathBehavior
    ...
}

func (p *ClimberPlayer) MeleeController() *meleeengine.Controller { return p.melee }

// After
type ClimberPlayer struct {
    *platformer.PlatformerCharacter
    *kitactors.MeleeCharacter
    *kitactors.PlayerDeathBehavior
    ...
}
```

- The constructor calls `kitactors.NewMeleeCharacter()` during initialization.
- `MeleeController()` is promoted from the embedded trait.
- `SetMelee` calls `p.MeleeCharacter.SetMeleeController(controller)` after wiring.
- `PlayerDeathBehavior` import path changes from `internal/game/entity/actors/methods` to `internal/kit/actors`.
- All game-specific methods (`Hurt`, `SetInventory`, `Inventory`, `SetMelee`, `Update`, `OnTouch`, `OnBlock`) remain on `ClimberPlayer`.

---

## Acceptance Criteria (testable)

| # | Criterion | Verification |
|---|-----------|--------------|
| AC-1 | `internal/kit/actors/` package exists with `doc.go` | `go build ./internal/kit/actors/` |
| AC-2 | `PlayerDeathBehavior.OnDie()` sets player health to 0 | Unit test |
| AC-3 | `ShooterCharacter.UpdateShooter()` calls `shooter.Update()` if non-nil | Unit test with mock shooter |
| AC-4 | `ShooterCharacter.UpdateShooter()` with nil shooter does not panic | Unit test |
| AC-5 | `ShooterCharacter.Shooter()` returns the shooter set via constructor or `SetShooter` | Unit test |
| AC-6 | `MeleeCharacter.MeleeController()` returns nil before `SetMeleeController` | Unit test |
| AC-7 | `MeleeCharacter.MeleeController()` returns controller after `SetMeleeController` | Unit test |
| AC-8 | `BatEnemy.Update()` calls `UpdateShooter()` then `Character.Update(space)` | Integration test (existing test suite) |
| AC-9 | `WolfEnemy.Update()` calls `UpdateShooter()` then `Character.Update(space)` | Integration test (existing test suite) |
| AC-10 | `ClimberPlayer.MeleeController()` promoted from embedded trait | Integration test (existing test suite) |
| AC-11 | `go list -deps ./internal/engine/...` contains no `internal/kit` or `internal/game` paths | CI check |
| AC-12 | `go list -deps ./internal/kit/...` contains no `internal/game` paths | CI check |
| AC-13 | `go test ./...` passes with no regressions | Full test run |
| AC-14 | `internal/kit/actors/` coverage ≥ 80% | Coverage report |
| AC-15 | No raw game-state string literals in `internal/kit/actors/` | Code review / grep |

---

## Test File Plan

**`internal/kit/actors/death_behavior_test.go`**
- Table-driven: `OnDie` sets health to 0 on a mock `PlatformerActorEntity`.

**`internal/kit/actors/shooter_character_test.go`**
- `UpdateShooter` with non-nil shooter: verify `shooter.Update()` called.
- `UpdateShooter` with nil shooter: no panic.
- `Shooter` / `SetShooter` round-trip.

**`internal/kit/actors/melee_character_test.go`**
- `MeleeController` returns nil before set.
- `SetMeleeController` / `MeleeController` round-trip.

**Mock placement:** Package-local mocks in `internal/kit/actors/mocks_test.go` (single package, not shared).

---

## Out of Scope

- Moving `internal/engine/combat/` to `internal/kit/combat/` (Story 049).
- Moving `internal/engine/weapon/` to `internal/kit/weapon/` (Story 050).
- Skills system migration (Story 051).
- UI split (Story 052).
- NPC base types.
- New gameplay features.
- Demonstrating multi-trait composition (e.g., shooter + melee brawler) — possible with this design but not required for this story.
