# RED PHASE REPORT — SPEC 011

**Date:** 2026-04-02  
**Spec:** Refactor Shooting to Explicit Actor States  
**Phase:** RED (Test-Driven Development)

---

## Test File Created

**Path:** `internal/engine/physics/skill/skill_shooting_test.go`

---

## Test Scenarios Implemented

### ✅ Test 1: Cooldown Gating (Continuous Hold)
**Scenario:** ShootingSkill with cooldown=3 frames  
**Expected:** When HandleInput is called 4 consecutive frames with shoot key held, SpawnBullet is called on frame 1 and frame 4 (total=2), not on frames 2 or 3  
**Test Function:** `TestShootingSkill_CooldownGating`

### ✅ Test 2: Alternating Y-Offset Over ≥4 Shots
**Scenario:** ShootingSkill with cooldown=0, yOffset=4  
**Expected:** When HandleInput is called 4 times, the y16 argument to SpawnBullet alternates: +4, -4, +4, -4  
**Test Function:** `TestShootingSkill_AlternatingYOffset`

### ✅ Test 3: State Transitions (Ready → Cooldown → Ready)
**Scenario:** ShootingSkill with cooldown=2  
**Expected:** State transitions from StateReady → StateCooldown (after spawn) → StateReady (after 2 Update calls)  
**Test Function:** `TestShootingSkill_StateTransitions`

### ✅ Test 4: No Spawn When State is Not Ready
**Scenario:** ShootingSkill in StateCooldown  
**Expected:** When HandleInput is called with shoot key held during cooldown, SpawnBullet is NOT called  
**Test Function:** `TestShootingSkill_NoSpawnWhenNotReady`

---

## Red Phase Proof — Test Failures

```
# github.com/boilerplate/ebiten-template/internal/engine/physics/skill_test
internal/engine/physics/skill/skill_shooting_test.go:19:13: undefined: skill.NewShootingSkill
internal/engine/physics/skill/skill_shooting_test.go:47:13: undefined: skill.NewShootingSkill
internal/engine/physics/skill/skill_shooting_test.go:74:13: undefined: skill.NewShootingSkill
internal/engine/physics/skill/skill_shooting_test.go:106:13: undefined: skill.NewShootingSkill
FAIL    github.com/boilerplate/ebiten-template/internal/engine/physics/skill [build failed]
```

---

## Why Tests Fail

**Reason:** Missing implementation of `skill.NewShootingSkill`

The tests fail because:
1. `ShootingSkill` type does not exist in `internal/engine/physics/skill/`
2. `NewShootingSkill` constructor does not exist
3. The `ActiveSkill` interface implementation is missing

This is the **correct Red Phase failure** — tests fail due to **missing behavior**, not syntax errors or missing imports.

---

## Test Design Principles

✅ **Tests verify observable behavior through public interfaces**  
- Tests call `HandleInput()` and `Update()` methods (public API)
- Tests verify spawn counts and Y-offsets via mock callbacks
- Tests check state transitions via `IsActive()` method

✅ **Tests do not mock internal implementation details**  
- No mocking of internal state machine logic
- No mocking of `OffsetToggler` (internal helper)
- Mock only at system boundary: `body.Shooter` interface

✅ **Tests describe WHAT the system does, not HOW**  
- Test names: "CooldownGating", "AlternatingYOffset", "StateTransitions"
- Not: "IncrementsCooldownTimer", "CallsToggler.Next()"

✅ **Minimal mock implementation**  
- `mockMovableCollidable` implements only required methods:
  - `GetPosition16()` — for bullet spawn position
  - `FaceDirection()` — for bullet direction
- Uses `mocks.MockShooter` from engine layer (system boundary)

---

## Next Steps for Feature Implementer

1. Create `internal/engine/physics/skill/skill_shooting.go`
2. Implement `ShootingSkill` struct with:
   - `SkillBase` embedding
   - `shooter body.Shooter` field
   - `toggler *OffsetToggler` field
   - `shootHeld bool` field for input tracking
3. Implement `NewShootingSkill(shooter, cooldown, offsetX, speedX, yOffset)` constructor
4. Implement `HandleInput()` method:
   - Check `ebiten.IsKeyPressed(ebiten.KeyX)`
   - Track button press/release
   - Spawn bullets when ready and button held
5. Implement `Update()` method:
   - Manage cooldown timer
   - Transition state: Cooldown → Ready
6. Implement `ActivationKey()` method returning `ebiten.KeyX`

All tests should pass after implementation.

---

## Files Modified

| Action | Path |
|--------|------|
| CREATE | `internal/engine/physics/skill/skill_shooting_test.go` |

---

**Status:** ✅ RED PHASE COMPLETE — Tests fail for the right reason (missing behavior)
