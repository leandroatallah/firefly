# SPEC-012 — 8-Direction Shooting (Cuphead-style)

**Branch:** `012-eight-direction-shooting`  
**Story:** `.agents/work/active/USER_STORY_012.md`  
**Bounded Context:** Game Logic (`internal/game/`) + Physics (`internal/engine/physics/skill/`)

---

## Technical Requirements

### 1. Refactor State Transition Architecture

**Problem:** Current `ShootingSkill.SetStateEnums()` injects 8 state enum fields, violating separation of concerns.

**Solution:** Replace with `StateTransitionHandler` interface.

#### New Contract: `internal/engine/contracts/body/state_transition_handler.go`

```go
package body

type ShootDirection int

const (
	ShootDirectionStraight ShootDirection = iota
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

#### Modified: `internal/engine/physics/skill/skill_shooting.go`

- Remove `SetStateEnums()` method
- Remove 8 `ActorStateEnum` fields
- Add `handler StateTransitionHandler` field
- Add `SetStateTransitionHandler(handler StateTransitionHandler)` method
- Replace `transitionToShootingState()` with `handler.TransitionToShooting(direction)`
- Replace `transitionToBaseState()` with `handler.TransitionFromShooting()`

### 2. Directional Input Detection

#### Modified: `internal/engine/physics/skill/skill_shooting.go`

Add method:
```go
func (s *ShootingSkill) detectShootDirection(b body.MovableCollidable, model *physicsmovement.PlatformMovementModel) body.ShootDirection
```

**Logic:**
- Read arrow key state (up/down/left/right)
- Determine if body is grounded (via `model.IsGrounded()`)
- Determine current movement state (idle/walking/jumping/falling)
- Return direction based on:
  - Down input + airborne → `ShootDirectionDown` or diagonal-down variants
  - Down input + grounded → `ShootDirectionStraight` (ignore down, ducking takes priority)
  - Up input → `ShootDirectionUp` or diagonal-up variants
  - Diagonal priority: up+forward > up > forward
  - No directional input → `ShootDirectionStraight`
  - Forward/back relative to `b.FaceDirection()`

**State-specific restrictions:**
- **Idle:** Only straight and up allowed (no horizontal input)
- **Walking:** Only straight and diagonal-up allowed (horizontal input already active; up input = diagonal-up, not straight up; down = ducking)
- **Ducking:** Only straight allowed (down input already active)
- **Jumping/Falling:** All 8 directions allowed

### 3. Bullet Velocity Calculation

#### Modified: `internal/engine/physics/skill/skill_shooting.go`

Add method:
```go
func (s *ShootingSkill) calculateBulletVelocity(direction body.ShootDirection, faceDir animation.FacingDirectionEnum) (vx16, vy16 int)
```

**Velocity mapping (fixed-point):**
- `ShootDirectionStraight`: `(±bulletSpeed, 0)`
- `ShootDirectionUp`: `(0, -bulletSpeed)`
- `ShootDirectionDown`: `(0, +bulletSpeed)`
- `ShootDirectionDiagonalUpForward`: `(±bulletSpeed*707/1000, -bulletSpeed*707/1000)`
- `ShootDirectionDiagonalDownForward`: `(±bulletSpeed*707/1000, +bulletSpeed*707/1000)`
- `ShootDirectionDiagonalUpBack`: `(∓bulletSpeed*707/1000, -bulletSpeed*707/1000)`
- `ShootDirectionDiagonalDownBack`: `(∓bulletSpeed*707/1000, +bulletSpeed*707/1000)`

Sign of `vx16` determined by `faceDir`.

### 4. Bullet Spawn Offset

#### Modified: `internal/engine/physics/skill/skill_shooting.go`

Add method:
```go
func (s *ShootingSkill) calculateSpawnOffset(direction body.ShootDirection, faceDir animation.FacingDirectionEnum) (offsetX16, offsetY16 int)
```

**Offset logic:**
- Straight: `(±spawnOffsetX, toggler.Next())`
- Up: `(0, -spawnOffsetX)` (reuse X offset as vertical offset)
- Down: `(0, +spawnOffsetX)`
- Diagonal: `(±spawnOffsetX*707/1000, ±spawnOffsetX*707/1000)`

### 5. Modified Shooter Contract

#### Modified: `internal/engine/contracts/body/shooter.go`

```go
type Shooter interface {
	SpawnBullet(x16, y16, vx16, vy16 int, owner interface{})
}
```

**Breaking change:** Add `vy16` parameter for vertical velocity.

### 6. Directional State Transitions

#### Modified: `internal/engine/physics/skill/skill_shooting.go`

In `HandleInput()`:
- Detect current direction via `detectShootDirection()`
- Store last direction in `lastDirection ShootDirection` field
- On direction change while shooting: call `handler.TransitionToShooting(newDirection)`
- On shoot button press: call `handler.TransitionToShooting(direction)`
- On shoot button release: call `handler.TransitionFromShooting()`

In `Update()`:
- Continue checking direction even while shooting
- Transition between directional shooting states without resetting cooldown

---

## Pre-conditions

- US-011 complete: `ShootingSkill` exists with horizontal shooting
- `body.Shooter` interface exists
- Actor state machine supports dynamic state transitions
- Input system provides arrow key state

---

## Post-conditions

- `ShootingSkill` uses `StateTransitionHandler` (no `SetStateEnums()`)
- Bullets spawn with 2D velocity (`vx16`, `vy16`)
- 8 directional shooting states registered in game layer
- Down-shooting restricted to airborne states
- Directional transitions work without cooldown reset
- All existing US-011 tests pass (regression-free)

---

## Integration Points

### Engine Layer (`internal/engine/`)
- `contracts/body/state_transition_handler.go` — new interface
- `contracts/body/shooter.go` — modified signature
- `physics/skill/skill_shooting.go` — refactored with directional logic

### Game Layer (`internal/game/`)
- `entity/player/player.go` — implements `StateTransitionHandler`
- `entity/player/state.go` — defines directional shooting state enums:
  - **Idle:** `IdleShooting`, `IdleShootingUp` (2 states)
  - **Walking:** `WalkingShooting`, `WalkingShootingDiagonalUp` (2 states)
  - **Ducking:** `DuckingShooting` (1 state)
  - **Jumping:** `JumpingShooting`, `JumpingShootingUp`, `JumpingShootingDown`, `JumpingShootingDiagonalUp`, `JumpingShootingDiagonalDown` (5 states)
  - **Falling:** `FallingShooting`, `FallingShootingUp`, `FallingShootingDown`, `FallingShootingDiagonalUp`, `FallingShootingDiagonalDown` (5 states)
  - **Total:** 15 directional shooting states
- `entity/player/state_machine.go` — registers directional state transitions
- `entity/bullet/bullet.go` — accepts `vy16` in constructor

---

## Red Phase Scenario

### Test: `TestShootingSkill_EightDirections`

**Setup:**
- Create `ShootingSkill` with mock `Shooter` and `StateTransitionHandler`
- Create mock `MovableCollidable` body
- Create mock `PlatformMovementModel` (grounded vs airborne)

**Scenario 1: Shoot Straight (no directional input)**
1. Press shoot button
2. Assert `handler.TransitionToShooting(ShootDirectionStraight)` called
3. Assert bullet spawned with `(vx16=bulletSpeed, vy16=0)`

**Scenario 2: Shoot Up**
1. Hold up arrow + press shoot
2. Assert `handler.TransitionToShooting(ShootDirectionUp)` called
3. Assert bullet spawned with `(vx16=0, vy16=-bulletSpeed)`

**Scenario 3: Shoot Down (airborne)**
1. Set `model.IsGrounded() = false`
2. Hold down arrow + press shoot
3. Assert `handler.TransitionToShooting(ShootDirectionDown)` called
4. Assert bullet spawned with `(vx16=0, vy16=+bulletSpeed)`

**Scenario 4: Shoot Down (grounded) — ignored**
1. Set `model.IsGrounded() = true`
2. Hold down arrow + press shoot
3. Assert `handler.TransitionToShooting(ShootDirectionStraight)` called (down ignored)

**Scenario 5: Diagonal Up-Forward**
1. Hold up + right arrow + press shoot (facing right)
2. Assert `handler.TransitionToShooting(ShootDirectionDiagonalUpForward)` called
3. Assert bullet spawned with `(vx16=bulletSpeed*707/1000, vy16=-bulletSpeed*707/1000)`

**Scenario 6: Direction Change Mid-Shooting**
1. Press shoot (straight)
2. Assert `handler.TransitionToShooting(ShootDirectionStraight)` called
3. Press up arrow (while still holding shoot)
4. Assert `handler.TransitionToShooting(ShootDirectionUp)` called
5. Assert cooldown NOT reset

**Scenario 7: Release Directional Input**
1. Hold up + shoot
2. Release up arrow (still holding shoot)
3. Assert `handler.TransitionToShooting(ShootDirectionStraight)` called

**Scenario 8: Ducking Shooting (grounded)**
1. Set `model.IsGrounded() = true`
2. Hold down arrow (enter ducking state)
3. Press shoot button
4. Assert `handler.TransitionToShooting(ShootDirectionStraight)` called
5. Assert bullet spawned with `(vx16=±bulletSpeed, vy16=0)`
6. Verify up/down directional input ignored while ducking (only straight allowed)

**Expected Result (Red Phase):**
- All assertions fail (methods/fields don't exist yet)
- Compilation errors on `StateTransitionHandler`, `ShootDirection`, `detectShootDirection()`, `calculateBulletVelocity()`

---

## Notes

- Diagonal normalization: `707/1000 ≈ 0.707 ≈ 1/√2` (integer-friendly)
- `OffsetToggler` still used for straight shots (alternating Y offset)
- Sprite mapping handled by game layer (not in this spec)
- Backward compatibility: existing `SpawnBullet()` calls must be updated to include `vy16=0`
