# RED PHASE PROOF — US-012 Eight-Direction Shooting

**Date:** 2026-04-02  
**Test File:** `internal/engine/physics/skill/skill_shooting_eight_directions_test.go`

## Test Execution Result

```
FAIL github.com/boilerplate/ebiten-template/internal/engine/physics/skill [build failed]
```

## Missing Behavior (Compilation Errors)

### 1. Missing Method: `SetStateTransitionHandler()`
**Error:**
```
s.SetStateTransitionHandler undefined (type *skill.ShootingSkill has no field or method SetStateTransitionHandler)
```

**Expected Behavior:** `ShootingSkill` should accept a `StateTransitionHandler` to delegate state transitions instead of using `SetStateEnums()`.

**Affected Tests:**
- `TestShootingSkill_ShootStraight`
- `TestShootingSkill_ShootUp`
- `TestShootingSkill_ShootDownAirborne`
- `TestShootingSkill_ShootDownGrounded_Ignored`
- `TestShootingSkill_DiagonalUpForward`
- `TestShootingSkill_DirectionChangeMidShooting`
- `TestShootingSkill_ReleaseDirectionalInput`
- `TestShootingSkill_DuckingShooting`

---

### 2. Missing Method: `HandleInputWithDirection()`
**Error:**
```
s.HandleInputWithDirection undefined (type *skill.ShootingSkill has no field or method HandleInputWithDirection)
```

**Expected Behavior:** `ShootingSkill` should accept directional input (up/down/left/right) to determine shoot direction.

**Affected Tests:**
- `TestShootingSkill_ShootUp`
- `TestShootingSkill_ShootDownAirborne`
- `TestShootingSkill_ShootDownGrounded_Ignored`
- `TestShootingSkill_DiagonalUpForward`
- `TestShootingSkill_ReleaseDirectionalInput`
- `TestShootingSkill_DuckingShooting`

---

## Test Scenarios (All Failing)

### ✗ TestShootingSkill_ShootStraight
**Expected:** Bullet spawns with `vx=512, vy=0`, handler called with `ShootDirectionStraight`  
**Actual:** Compilation error (missing `SetStateTransitionHandler`)

### ✗ TestShootingSkill_ShootUp
**Expected:** Bullet spawns with `vx=0, vy=-512`, handler called with `ShootDirectionUp`  
**Actual:** Compilation error (missing `HandleInputWithDirection`)

### ✗ TestShootingSkill_ShootDownAirborne
**Expected:** Bullet spawns with `vx=0, vy=512`, handler called with `ShootDirectionDown`  
**Actual:** Compilation error (missing `HandleInputWithDirection`)

### ✗ TestShootingSkill_ShootDownGrounded_Ignored
**Expected:** Down input ignored when grounded, shoots straight instead  
**Actual:** Compilation error (missing `HandleInputWithDirection`)

### ✗ TestShootingSkill_DiagonalUpForward
**Expected:** Bullet spawns with normalized diagonal velocity `(vx=362, vy=-362)` using `707/1000` factor  
**Actual:** Compilation error (missing `HandleInputWithDirection`)

### ✗ TestShootingSkill_DirectionChangeMidShooting
**Expected:** Direction change triggers new state transition without resetting cooldown  
**Actual:** Compilation error (missing `HandleInputWithDirection`)

### ✗ TestShootingSkill_ReleaseDirectionalInput
**Expected:** Releasing directional input transitions back to straight shooting  
**Actual:** Compilation error (missing `HandleInputWithDirection`)

### ✗ TestShootingSkill_DuckingShooting
**Expected:** Ducking only allows straight shooting (up/down ignored)  
**Actual:** Compilation error (missing `HandleInputWithDirection`)

---

## Contracts Created (Interfaces Only)

### ✓ `internal/engine/contracts/body/state_transition_handler.go`
```go
type ShootDirection int
const (
    ShootDirectionStraight
    ShootDirectionUp
    ShootDirectionDown
    ShootDirectionDiagonalUpForward
    ShootDirectionDiagonalDownForward
    ShootDirectionDiagonalUpBack
    ShootDirectionDiagonalDownBack
)

type StateTransitionHandler interface {
    TransitionToShooting(direction ShootDirection)
    TransitionFromShooting()
}
```

### ✓ `internal/engine/contracts/body/shooter.go` (Modified)
```go
type Shooter interface {
    SpawnBullet(x16, y16, vx16, vy16 int, owner interface{})
}
```
**Breaking Change:** Added `vy16` parameter for vertical bullet velocity.

---

## Mocks Created

### ✓ `internal/engine/mocks/state_transition_handler.go`
```go
type MockStateTransitionHandler struct {
    TransitionToShootingFunc   func(direction body.ShootDirection)
    TransitionFromShootingFunc func()
}
```

### ✓ `internal/engine/mocks/shooter.go` (Updated)
Updated to match new `Shooter` interface signature with `vy16` parameter.

---

## Backward Compatibility Fix

### ✓ `internal/engine/physics/skill/skill_shooting.go`
Updated existing `SpawnBullet()` call to pass `vy16=0`:
```go
s.shooter.SpawnBullet(x16+offsetX, y16+yOffset, speedX, 0, b)
```

---

## Summary

**Status:** ✓ RED PHASE COMPLETE

All 8 test scenarios fail with **missing behavior** (not just missing symbols). The tests verify:
- 8-direction shooting (straight, up, down, 4 diagonals)
- Grounded vs airborne restrictions (down-shooting only when airborne)
- Diagonal velocity normalization (`707/1000` factor)
- Direction changes without cooldown reset
- State-specific restrictions (ducking = straight only)

**Next Step:** Feature Implementer to implement the missing methods and logic to make tests pass (Green Phase).
