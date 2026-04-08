# SPEC-026 — Enhanced Shooting Skill (Inventory-Aware)

**Story ID:** 026  
**Branch:** `026-shooting-skill-inventory`  
**Bounded Context:** Physics (`internal/engine/physics/skill/`)  
**Dependencies:** US-022 (Weapon System), US-023 (Projectile Manager), US-025 (Inventory)

---

## Overview

Refactor `ShootingSkill` to delegate firing to `inventory.Inventory` instead of directly calling `body.Shooter`. This decouples the skill from scene-level shooting logic and enables weapon switching via input commands.

---

## Technical Requirements

### 1. Constructor Signature Change

**Current:**
```go
func NewShootingSkill(shooter body.Shooter, cooldownFrames, spawnOffsetX16, bulletSpeedX16, yOffset int) *ShootingSkill
```

**New:**
```go
func NewShootingSkill(inv *inventory.Inventory) *ShootingSkill
```

**Rationale:** Cooldown, offset, speed, and yOffset are now weapon-specific properties. The skill only needs the inventory reference.

---

### 2. Field Changes in `ShootingSkill`

**Remove:**
- `shooter body.Shooter`
- `spawnOffsetX int`
- `bulletSpeed int`
- `toggler *OffsetToggler`

**Add:**
- `inv *inventory.Inventory`

**Preserve:**
- `SkillBase` (state machine: `StateReady`, `StateActive`)
- `shootHeld bool` (for release detection)
- `handler body.StateTransitionHandler` (animation transitions)
- `lastDirection body.ShootDirection` (8-directional logic from US-012)
- `directionSet bool`

---

### 3. Input Handling Changes

#### `HandleInput(b body.MovableCollidable, model *physicsmovement.PlatformMovementModel, space body.BodiesSpace)`

**Current behavior:**
- Reads `input.CommandsReader().Shoot`
- Calls `HandleInputWithDirection` if shoot is pressed

**New behavior:**
- Read `input.CommandsReader()` once
- If `WeaponNext` → call `inv.SwitchNext()`
- If `WeaponPrev` → call `inv.SwitchPrev()`
- If `Shoot` → call `HandleInputWithDirection`

**Note:** `WeaponNext` and `WeaponPrev` fields must be added to `input.PlayerCommands` (see Pre-conditions).

---

#### `HandleInputWithDirection(...)`

**Current behavior:**
- Detects shoot direction (8-directional)
- Constructs bullet parameters inline
- Calls `s.shooter.SpawnBullet(...)`
- Manages cooldown state

**New behavior:**
- Detect shoot direction (unchanged)
- Get active weapon: `weapon := s.inv.ActiveWeapon()`
- If `weapon == nil` → return (no weapon equipped)
- If `!weapon.CanFire()` → return (weapon on cooldown)
- Calculate spawn position (x16, y16) from body position and face direction
- Call `weapon.Fire(x16, y16, b.FaceDirection(), direction)`
- Trigger animation transition via `s.handler.TransitionToShooting(direction)` if handler is set
- **Do NOT manage cooldown state** — weapon owns cooldown

**Removed logic:**
- `s.state = StateActive`
- `s.timer = s.cooldown`
- Bullet velocity calculation (moved to weapon)
- Spawn offset calculation (moved to weapon)

---

### 4. Update Method Changes

**Current behavior:**
```go
func (s *ShootingSkill) Update(b body.MovableCollidable, model *physicsmovement.PlatformMovementModel) {
	if s.state == StateActive {
		s.timer--
		if s.timer <= 0 {
			s.state = StateReady
		}
	}
	// shoot release detection
}
```

**New behavior:**
```go
func (s *ShootingSkill) Update(b body.MovableCollidable, model *physicsmovement.PlatformMovementModel) {
	wasHeld := s.shootHeld
	s.shootHeld = input.CommandsReader().Shoot

	if !s.shootHeld && wasHeld && s.handler != nil {
		s.handler.TransitionFromShooting()
	}
}
```

**Rationale:** Cooldown is now managed by `Weapon.Update()`, not by the skill. The skill only tracks shoot button release for animation transitions.

---

### 5. Direction Detection (Unchanged)

The 8-directional logic in `detectShootDirection` remains identical:
- Ducking → `ShootDirectionStraight`
- Down + airborne → `ShootDirectionDown` or `ShootDirectionDiagonalDownForward`
- Up → `ShootDirectionUp` or `ShootDirectionDiagonalUpForward`
- Default → `ShootDirectionStraight`

---

### 6. Removed Code

#### `body.Shooter` Interface
**File:** `internal/engine/contracts/body/shooter.go`

**Action:** Delete file.

**Rationale:** No longer needed. Weapons handle firing via `combat.Weapon.Fire()`.

#### Mock Shooter
**File:** Search for `MockShooter` or `FakeShooter` in test files.

**Action:** Remove all mock shooter implementations.

---

## Pre-conditions

1. `internal/engine/combat/inventory/inventory.go` exists and implements `combat.Inventory` interface.
2. `combat.Weapon` interface includes:
   - `Fire(x16, y16 int, faceDir animation.FacingDirectionEnum, direction body.ShootDirection)`
   - `CanFire() bool`
   - `Update()`
3. `input.PlayerCommands` struct must be extended with:
   ```go
   WeaponNext bool
   WeaponPrev bool
   ```
   Mapped to keys (e.g., `Q` and `E`) in `ReadPlayerCommands()`.

---

## Post-conditions

1. `ShootingSkill` no longer depends on `body.Shooter`.
2. Firing delegates to `inventory.ActiveWeapon().Fire()`.
3. Weapon switching works via `input.CommandsReader().WeaponNext` / `WeaponPrev`.
4. Old constructor `NewShootingSkill(shooter, cooldown, ...)` is removed.
5. `body.Shooter` interface and mocks are deleted.
6. All existing tests pass with updated constructor signature.
7. New tests cover:
   - Fire delegates to active weapon
   - No fire when `weapon.CanFire()` returns false
   - No fire when inventory is empty
   - Weapon switching on input

---

## Integration Points

### Within Bounded Context (Physics)
- `internal/engine/physics/skill/skill_shooting.go` (modified)
- `internal/engine/physics/skill/skill_shooting_test.go` (updated)
- `internal/engine/physics/skill/skill_shooting_eight_directions_test.go` (updated)
- `internal/engine/physics/skill/factory.go` (constructor call updated)

### Cross-Context Dependencies
- `internal/engine/combat/inventory/inventory.go` (read-only)
- `internal/engine/contracts/combat/inventory.go` (read-only)
- `internal/engine/contracts/combat/weapon.go` (read-only)
- `internal/engine/input/commands.go` (modified: add `WeaponNext`, `WeaponPrev`)

---

## Red Phase: Failing Test Scenario

**Test File:** `internal/engine/physics/skill/skill_shooting_test.go`

### Test 1: Fire Delegates to Active Weapon
```go
func TestShootingSkill_FireDelegatesToActiveWeapon(t *testing.T) {
	// Given: inventory with a mock weapon
	mockWeapon := &MockWeapon{canFire: true}
	inv := inventory.New()
	inv.AddWeapon(mockWeapon)
	
	skill := NewShootingSkill(inv)
	
	// When: HandleInput is called with shoot pressed
	// Then: mockWeapon.Fire() should be called once
	// Assert: mockWeapon.fireCalled == true
}
```

**Expected failure:** `NewShootingSkill` still requires old signature; `Fire()` not called.

---

### Test 2: No Fire When CanFire Returns False
```go
func TestShootingSkill_NoFireWhenWeaponOnCooldown(t *testing.T) {
	// Given: weapon with CanFire() == false
	mockWeapon := &MockWeapon{canFire: false}
	inv := inventory.New()
	inv.AddWeapon(mockWeapon)
	
	skill := NewShootingSkill(inv)
	
	// When: HandleInput is called with shoot pressed
	// Then: mockWeapon.Fire() should NOT be called
	// Assert: mockWeapon.fireCalled == false
}
```

**Expected failure:** Old code calls `shooter.SpawnBullet` regardless of cooldown.

---

### Test 3: Weapon Switching on Input
```go
func TestShootingSkill_WeaponSwitchingOnInput(t *testing.T) {
	// Given: inventory with 2 weapons
	weapon1 := &MockWeapon{id: "pistol"}
	weapon2 := &MockWeapon{id: "shotgun"}
	inv := inventory.New()
	inv.AddWeapon(weapon1)
	inv.AddWeapon(weapon2)
	
	skill := NewShootingSkill(inv)
	
	// When: HandleInput is called with WeaponNext pressed
	// Then: inv.ActiveWeapon().ID() should be "shotgun"
	
	// When: HandleInput is called with WeaponPrev pressed
	// Then: inv.ActiveWeapon().ID() should be "pistol"
}
```

**Expected failure:** `WeaponNext` / `WeaponPrev` fields don't exist in `PlayerCommands`.

---

### Test 4: No Fire When Inventory Empty
```go
func TestShootingSkill_NoFireWhenInventoryEmpty(t *testing.T) {
	// Given: empty inventory
	inv := inventory.New()
	skill := NewShootingSkill(inv)
	
	// When: HandleInput is called with shoot pressed
	// Then: no panic, no fire
}
```

**Expected failure:** Old code panics or calls shooter with nil weapon.

---

## Design Decisions

1. **Cooldown ownership moved to Weapon:** Each weapon manages its own cooldown state. `ShootingSkill` no longer tracks `StateActive` / `StateReady` for cooldown purposes.

2. **Spawn position calculation remains in skill:** The skill calculates `x16, y16` from body position and face direction, then passes it to `weapon.Fire()`. This keeps body-specific logic (width adjustment, ducking) in the skill layer.

3. **Direction detection unchanged:** The 8-directional logic (US-012) is preserved as-is. Only the firing delegation changes.

4. **Input commands extended:** `WeaponNext` and `WeaponPrev` are added to `input.PlayerCommands` to support weapon switching without introducing a new input system.

5. **No new interfaces:** Uses existing `combat.Inventory` and `combat.Weapon` contracts. No new contracts needed.
