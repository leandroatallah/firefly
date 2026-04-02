# GREEN PHASE REPORT — SPEC 011

**Date:** 2026-04-02  
**Spec:** Refactor Shooting to Engine Skill  
**Phase:** GREEN (Implementation)

---

## Implementation Summary

Successfully implemented `ShootingSkill` as an `ActiveSkill` in the engine layer with minimal code to pass all 4 test scenarios.

---

## Files Created

### 1. `internal/engine/physics/skill/skill_shooting.go`
**Purpose:** Core `ShootingSkill` implementation

**Key Components:**
- `ShootingSkill` struct with `SkillBase` embedding
- `NewShootingSkill()` constructor
- `HandleInput()` - spawns bullets when ready
- `Update()` - manages cooldown timer
- `ActivationKey()` - returns `ebiten.KeyX`

**Design Decisions:**
- `HandleInput()` being called signals shoot intent (no internal ebiten key check for testability)
- Uses `StateActive` for cooldown period (matches `IsActive()` semantics)
- Instant bullet spawn (no duration phase like DashSkill)
- State flow: Ready → Active (cooldown) → Ready

### 2. `internal/engine/physics/skill/offset_toggler.go`
**Purpose:** Y-offset alternation helper (moved from game layer)

**Implementation:**
```go
type OffsetToggler struct {
	current int
}

func NewOffsetToggler(offset int) *OffsetToggler {
	return &OffsetToggler{current: -offset}
}

func (o *OffsetToggler) Next() int {
	o.current = -o.current
	return o.current
}
```

---

## Test Results

```
=== RUN   TestShootingSkill_CooldownGating
--- PASS: TestShootingSkill_CooldownGating (0.00s)
=== RUN   TestShootingSkill_AlternatingYOffset
--- PASS: TestShootingSkill_AlternatingYOffset (0.00s)
=== RUN   TestShootingSkill_StateTransitions
--- PASS: TestShootingSkill_StateTransitions (0.00s)
=== RUN   TestShootingSkill_NoSpawnWhenNotReady
--- PASS: TestShootingSkill_NoSpawnWhenNotReady (0.00s)
PASS
ok      github.com/boilerplate/ebiten-template/internal/engine/physics/skill   0.966s
```

✅ **All 4 tests pass**

---

## Implementation Details

### Constructor Signature
```go
func NewShootingSkill(
    shooter body.Shooter,
    cooldownFrames int,
    spawnOffsetX16 int,
    bulletSpeedX16 int,
    yOffset int,
) *ShootingSkill
```

### HandleInput Logic
1. Check if state is `StateReady`
2. Get actor position and facing direction
3. Calculate spawn position with X-offset (flipped for left-facing)
4. Get alternating Y-offset from `OffsetToggler`
5. Spawn bullet via `shooter.SpawnBullet()`
6. Transition to `StateActive` with cooldown timer

### Update Logic
1. If state is `StateActive`, decrement timer
2. When timer reaches 0, transition to `StateReady`

### State Semantics
- **StateReady**: Can shoot
- **StateActive**: Cooldown period (IsActive() returns true)
- No `StateCooldown` used (StateActive serves this purpose)

---

## Code Characteristics

✅ **Minimal implementation** - Only code needed to pass tests  
✅ **No unused fields** - Removed `shootHeld` after simplification  
✅ **Testable design** - `HandleInput()` doesn't check ebiten keys internally  
✅ **Follows existing patterns** - Matches `DashSkill` structure  
✅ **Reusable components** - `OffsetToggler` moved to engine layer  

---

## Integration Points

### Dependencies
- `body.Shooter` interface (for bullet spawning)
- `body.MovableCollidable` interface (for actor position/direction)
- `SkillBase` struct (for state management)
- `OffsetToggler` helper (for Y-offset alternation)

### Public API
- `NewShootingSkill()` - Constructor
- `HandleInput()` - Trigger shooting attempt
- `Update()` - Manage cooldown
- `IsActive()` - Check if in cooldown (inherited from SkillBase)
- `ActivationKey()` - Returns `ebiten.KeyX`

---

## Next Steps

1. **Integration Testing** - Test with actual game actors
2. **Input Handling** - Wire up ebiten key checks in game layer
3. **State Transitions** - Implement shooting states (IdleShooting, WalkingShooting, etc.)
4. **Migration** - Move existing shooting logic to use new skill
5. **Cleanup** - Remove old `ShootingSkill` from game layer

---

## Files Modified

| Action | Path |
|--------|------|
| CREATE | `internal/engine/physics/skill/skill_shooting.go` |
| CREATE | `internal/engine/physics/skill/offset_toggler.go` |
| UPDATE | `internal/engine/physics/skill/skill_shooting_test.go` (mock methods) |

---

**Status:** ✅ GREEN PHASE COMPLETE — All tests pass with minimal implementation
