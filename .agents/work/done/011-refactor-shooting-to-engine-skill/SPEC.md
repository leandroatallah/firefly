# SPEC 011 — Refactor Shooting to Explicit Actor States

**Branch:** `011-refactor-shooting-to-engine-skill`
**Bounded Context:** Entity (`internal/engine/entity/actors/`)
**Story:** `USER_STORY.md`

---

## 1. Design Decision: Explicit States vs. Skill Modifiers

**Chosen approach:** Shooting states are **explicit actor states**, not skill modifiers.

**Rationale:**
- In Cuphead, "idle shooting" is visually distinct from "idle" — it's a different state, not idle with an overlay
- Each shooting variant has its own sprite sheet, animation timing, and potentially different hitboxes
- State machine clarity: "What is the actor doing?" should have a single, explicit answer
- Matches existing architecture: `DashState` is a separate state, not a movement modifier

**Consequence:** We register shooting state variants (e.g., `IdleShooting`, `WalkingShooting`) as first-class actor states.

---

## 2. Technical Requirements

### 2.1 Register Shooting State Variants

**File:** `internal/engine/entity/actors/actor_state.go`

Add new state enums:

```go
var (
    // Existing states
    Idle    ActorStateEnum
    Walking ActorStateEnum
    Jumping ActorStateEnum
    Falling ActorStateEnum
    // ... etc.
    
    // NEW: Shooting state variants
    IdleShooting    ActorStateEnum
    WalkingShooting ActorStateEnum
    JumpingShooting ActorStateEnum
    FallingShooting ActorStateEnum
    // Add more as needed (ducking, landing, etc.)
)

func init() {
    // Existing registrations
    Idle = RegisterState("idle", func(b BaseState) ActorState { return &IdleState{BaseState: b} })
    Walking = RegisterState("walk", func(b BaseState) ActorState { return &WalkState{BaseState: b} })
    // ...
    
    // NEW: Shooting state registrations
    IdleShooting = RegisterState("idle_shoot", func(b BaseState) ActorState { 
        return &IdleShootingState{BaseState: b} 
    })
    WalkingShooting = RegisterState("walk_shoot", func(b BaseState) ActorState { 
        return &WalkingShootingState{BaseState: b} 
    })
    JumpingShooting = RegisterState("jump_shoot", func(b BaseState) ActorState { 
        return &JumpingShootingState{BaseState: b} 
    })
    FallingShooting = RegisterState("fall_shoot", func(b BaseState) ActorState { 
        return &FallingShootingState{BaseState: b} 
    })
}
```

### 2.2 Implement Shooting State Types

**New file:** `internal/engine/entity/actors/shooting_states.go`

```go
package actors

// IdleShootingState is the idle state while shooting.
type IdleShootingState struct {
    BaseState
}

// WalkingShootingState is the walking state while shooting.
type WalkingShootingState struct {
    BaseState
}

// JumpingShootingState is the jumping state while shooting.
type JumpingShootingState struct {
    BaseState
}

// FallingShootingState is the falling state while shooting.
type FallingShootingState struct {
    BaseState
}

// All shooting states inherit default behavior from BaseState.
// No custom logic needed unless specific animation timing is required.
```

### 2.3 Move `ShootingSkill` to Engine Layer

**Current location:** `internal/game/entity/actors/states/shooting_skill.go`  
**New location:** `internal/engine/physics/skill/skill_shooting.go`

The refactored `ShootingSkill` implements `ActiveSkill` and **triggers state transitions** instead of spawning bullets directly:

```go
package skill

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
    "github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
    physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
)

type ShootingSkill struct {
    SkillBase
    
    shooter        body.Shooter
    toggler        *OffsetToggler
    spawnOffsetX16 int
    bulletSpeedX16 int
    activationKey  ebiten.Key
    
    shootHeld      bool // tracks if shoot button is currently held
}

func NewShootingSkill(shooter body.Shooter, cooldownFrames int, spawnOffsetX16, bulletSpeedX16, yOffset int) *ShootingSkill {
    return &ShootingSkill{
        SkillBase: SkillBase{
            state:    StateReady,
            cooldown: cooldownFrames,
        },
        shooter:        shooter,
        toggler:        NewOffsetToggler(yOffset),
        spawnOffsetX16: spawnOffsetX16,
        bulletSpeedX16: bulletSpeedX16,
        activationKey:  ebiten.KeyX,
    }
}

func (s *ShootingSkill) ActivationKey() ebiten.Key {
    return s.activationKey
}

// HandleInput checks shoot button and triggers state transitions.
func (s *ShootingSkill) HandleInput(body body.MovableCollidable, model *physicsmovement.PlatformMovementModel, space body.BodiesSpace) {
    wasHeld := s.shootHeld
    s.shootHeld = ebiten.IsKeyPressed(s.activationKey)
    
    // Transition to shooting state when button pressed
    if s.shootHeld && !wasHeld {
        s.transitionToShootingState(body)
    }
    
    // Transition back to base state when button released
    if !s.shootHeld && wasHeld {
        s.transitionToBaseState(body)
    }
    
    // Spawn bullets while in shooting state and cooldown allows
    if s.shootHeld && s.state == StateReady {
        s.spawnBullet(body)
        s.state = StateCooldown
        s.timer = s.cooldown
    }
}

// Update manages cooldown timer.
func (s *ShootingSkill) Update(body body.MovableCollidable, model *physicsmovement.PlatformMovementModel) {
    s.SkillBase.Update(body, model)
    
    if s.state == StateCooldown {
        s.timer--
        if s.timer <= 0 {
            s.state = StateReady
        }
    }
}

func (s *ShootingSkill) transitionToShootingState(body body.MovableCollidable) {
    actor, ok := body.(actors.Stateful)
    if !ok {
        return
    }
    
    currentState := actor.State()
    var newState actors.ActorStateEnum
    
    switch currentState {
    case actors.Idle:
        newState = actors.IdleShooting
    case actors.Walking:
        newState = actors.WalkingShooting
    case actors.Jumping:
        newState = actors.JumpingShooting
    case actors.Falling:
        newState = actors.FallingShooting
    default:
        return // No shooting variant for this state
    }
    
    state, err := actor.NewState(newState)
    if err == nil {
        actor.SetState(state)
    }
}

func (s *ShootingSkill) transitionToBaseState(body body.MovableCollidable) {
    actor, ok := body.(actors.Stateful)
    if !ok {
        return
    }
    
    currentState := actor.State()
    var newState actors.ActorStateEnum
    
    switch currentState {
    case actors.IdleShooting:
        newState = actors.Idle
    case actors.WalkingShooting:
        newState = actors.Walking
    case actors.JumpingShooting:
        newState = actors.Jumping
    case actors.FallingShooting:
        newState = actors.Falling
    default:
        return // Not in a shooting state
    }
    
    state, err := actor.NewState(newState)
    if err == nil {
        actor.SetState(state)
    }
}

func (s *ShootingSkill) spawnBullet(body body.MovableCollidable) {
    x16, y16 := body.GetPosition16()
    
    if body.FaceDirection() == animation.FaceDirectionRight {
        x16 += s.spawnOffsetX16
    } else {
        x16 -= s.spawnOffsetX16
    }
    
    y16 += s.toggler.Next()
    
    speedX16 := s.bulletSpeedX16
    if body.FaceDirection() == animation.FaceDirectionLeft {
        speedX16 = -speedX16
    }
    
    s.shooter.SpawnBullet(x16, y16, speedX16, body)
}
```

### 2.4 Move `OffsetToggler` to Engine

**Current location:** `internal/game/entity/actors/states/offset_toggler.go`  
**New location:** `internal/engine/physics/skill/offset_toggler.go`

Update package from `gamestates` to `skill`. No logic changes.

### 2.5 Move `Bullet` to Engine

**Current location:** `internal/game/entity/actors/states/bullet.go`  
**New location:** `internal/engine/entity/projectiles/bullet.go` (new package)

Update package from `gamestates` to `projectiles`. No logic changes.

### 2.6 Update `Character.handleState()` for Shooting States

**File:** `internal/engine/entity/actors/character.go`

Add shooting state transitions to the state machine logic:

```go
func (c *Character) handleState() {
    // ... existing logic ...
    
    // Shooting states follow the same transitions as their base states
    switch state {
    case IdleShooting:
        if c.IsWalking() {
            setNewState(WalkingShooting)
        } else if c.IsFalling() {
            setNewState(FallingShooting)
        }
    case WalkingShooting:
        if !c.IsWalking() {
            setNewState(IdleShooting)
        } else if c.IsFalling() {
            setNewState(FallingShooting)
        }
    case JumpingShooting:
        if c.state.IsAnimationFinished() {
            setNewState(IdleShooting)
        }
    case FallingShooting:
        if !c.IsFalling() {
            setNewState(Landing) // Landing interrupts shooting
        }
    }
}
```

### 2.7 Remove Shooting Logic from `GroundedState`

**File:** `internal/game/entity/actors/states/grounded_state.go`

Remove shooting-specific code:

```go
// REMOVE these lines:
if g.deps.Shooting != nil && input.ShootHeld() {
    g.deps.Shooting.Update(g.deps.Body)
}
```

Remove from `GroundedDeps`:

```go
type GroundedDeps struct {
    Input GroundedInput
    Body  contractsbody.Movable
    // Shooting removed
}
```

Remove from `GroundedInput`:

```go
type GroundedInput interface {
    HorizontalInput() int
    DuckHeld() bool
    HasCeilingClearance() bool
    JumpPressed() bool
    DashPressed() bool
    AimLockHeld() bool
    // ShootHeld() removed
}
```

---

## 3. Pre-conditions

- `ShootingSkill` exists in `internal/game/entity/actors/states/` (US-010 ✅).
- `OffsetToggler` exists in `internal/game/entity/actors/states/` (US-008 ✅).
- `Bullet` exists in `internal/game/entity/actors/states/` (US-010 ✅).
- `ActiveSkill` interface exists in `internal/engine/physics/skill/skill.go` ✅.
- `body.Shooter` contract exists in `internal/engine/contracts/body/shooter.go` ✅.
- `Character.skills` field exists and is updated in `Character.Update()` ✅.
- Actor state registration system exists in `actor_state.go` ✅.

## 4. Post-conditions

- Shooting state variants registered: `IdleShooting`, `WalkingShooting`, `JumpingShooting`, `FallingShooting`.
- `ShootingSkill` implements `ActiveSkill` and triggers state transitions.
- `ShootingSkill` lives in `internal/engine/physics/skill/`.
- `OffsetToggler` lives in `internal/engine/physics/skill/`.
- `Bullet` lives in `internal/engine/entity/projectiles/`.
- `GroundedState` no longer contains shooting-specific logic.
- `GroundedInput` no longer has `ShootHeld()` method.
- Sprite system can map shooting states to distinct sprite sheets (e.g., "idle_shoot.png").
- All existing shooting tests pass with no behavioral changes.
- Code coverage ≥74.6% (no regression).

---

## 5. Integration Points

| Point | Detail |
|---|---|
| `internal/engine/entity/actors/actor_state.go` | Register shooting state variants |
| `internal/engine/entity/actors/shooting_states.go` | New file; implement shooting state types |
| `internal/engine/entity/actors/character.go` | Update `handleState()` for shooting state transitions |
| `internal/engine/physics/skill/skill_shooting.go` | New file; implements `ActiveSkill`, triggers state transitions |
| `internal/engine/physics/skill/offset_toggler.go` | Moved from `internal/game/entity/actors/states/` |
| `internal/engine/entity/projectiles/bullet.go` | Moved from `internal/game/entity/actors/states/` |
| `internal/game/entity/actors/states/grounded_state.go` | Remove shooting logic |
| `internal/game/entity/actors/states/grounded_input.go` | Remove `ShootHeld()` |
| `internal/game/entity/actors/states/mocks_test.go` | Remove `ShootHeldFunc` from `MockInputSource` |
| `internal/engine/physics/skill/skill_shooting_test.go` | Moved from `shooting_skill_test.go` |

---

## 6. Red Phase — Failing Test Scenarios

File: `internal/engine/physics/skill/skill_shooting_test.go` (moved from game layer)

### Test 1 — State transition on shoot press

```
GIVEN Character in Idle state with ShootingSkill registered
WHEN HandleInput is called with shoot key pressed
THEN Character transitions to IdleShooting state
```

### Test 2 — State transition on shoot release

```
GIVEN Character in IdleShooting state
WHEN HandleInput is called with shoot key released
THEN Character transitions back to Idle state
```

### Test 3 — Cooldown gating (continuous hold)

```
GIVEN ShootingSkill with cooldown=3 frames
WHEN HandleInput is called 4 consecutive frames with shoot key held
THEN SpawnBullet is called on frame 1 and frame 4 (total=2), not on frames 2 or 3
```

### Test 4 — Alternating Y-offset over ≥4 shots

```
GIVEN ShootingSkill with cooldown=0, yOffset=4
WHEN HandleInput is called 4 times (cooldown=0 so each call spawns)
THEN the y16 argument to SpawnBullet alternates: +4, -4, +4, -4
```

### Test 5 — Movement transitions preserve shooting state

```
GIVEN Character in IdleShooting state
WHEN horizontal input applied (triggers walking)
THEN Character transitions to WalkingShooting state (not Walking)
```

### Test 6 — Shooting interrupted by landing

```
GIVEN Character in FallingShooting state
WHEN Character lands (OnGround becomes true)
THEN Character transitions to Landing state (shooting interrupted)
```

---

## 7. Files to Create / Modify / Delete

| Action | Path |
|---|---|
| Modify | `internal/engine/entity/actors/actor_state.go` (register shooting states) |
| Create | `internal/engine/entity/actors/shooting_states.go` |
| Modify | `internal/engine/entity/actors/character.go` (update `handleState()`) |
| Create | `internal/engine/physics/skill/skill_shooting.go` |
| Create | `internal/engine/physics/skill/skill_shooting_test.go` |
| Move | `internal/game/entity/actors/states/offset_toggler.go` → `internal/engine/physics/skill/offset_toggler.go` |
| Move | `internal/game/entity/actors/states/offset_toggler_test.go` → `internal/engine/physics/skill/offset_toggler_test.go` |
| Move | `internal/game/entity/actors/states/bullet.go` → `internal/engine/entity/projectiles/bullet.go` |
| Modify | `internal/game/entity/actors/states/grounded_state.go` (remove shooting logic) |
| Modify | `internal/game/entity/actors/states/grounded_input.go` (remove `ShootHeld()`) |
| Modify | `internal/game/entity/actors/states/mocks_test.go` (remove `ShootHeldFunc`) |
| Delete | `internal/game/entity/actors/states/shooting_skill.go` |
| Delete | `internal/game/entity/actors/states/shooting_skill_test.go` |

---

## 8. Migration Strategy

### Phase 1 (RED): Register States & Move Files
1. Register shooting state enums in `actor_state.go`
2. Create `shooting_states.go` with empty state implementations
3. Move `offset_toggler.go`, `bullet.go` to engine layer
4. Update imports across codebase
5. Tests fail: missing `ShootingSkill` implementation

### Phase 2 (GREEN): Implement ShootingSkill
1. Create `skill_shooting.go` with `ActiveSkill` implementation
2. Implement state transition logic (`transitionToShootingState`, `transitionToBaseState`)
3. Implement bullet spawning with cooldown
4. All tests pass

### Phase 3 (REFACTOR): Clean Up Game Layer
1. Remove shooting logic from `GroundedState`
2. Remove `ShootHeld()` from `GroundedInput`
3. Update `Character.handleState()` for shooting state transitions
4. Verify no regressions

---

## 9. Behavioral Invariants (Must Not Change)

- Bullets spawn at the same rate (cooldown enforcement)
- Y-offset alternates identically to US-010
- Bullet collision and removal behavior unchanged
- Shooting is suppressed while dashing (dash state takes priority)

---

## 10. Sprite Mapping

Shooting states map to distinct sprite sheets:

| State | Sprite Key | Example File |
|---|---|---|
| `IdleShooting` | `"idle_shoot"` | `idle_shoot.png` |
| `WalkingShooting` | `"walk_shoot"` | `walk_shoot.png` |
| `JumpingShooting` | `"jump_shoot"` | `jump_shoot.png` |
| `FallingShooting` | `"fall_shoot"` | `fall_shoot.png` |

The sprite system already supports this via `SpriteAssets.AddSprite(state, path, loop)`.

---

## 11. Future Extensions

This design supports future shooting variants:
- `IdleShootingUp` (shooting straight up while idle)
- `WalkingShootingDiagonal` (shooting diagonally while walking)
- `DuckingShooting` (shooting while ducking)

Each variant is a new state with its own sprite sheet and transition rules.

---

## 3. Pre-conditions

- `ShootingSkill` exists in `internal/game/entity/actors/states/` (US-010 ✅).
- `OffsetToggler` exists in `internal/game/entity/actors/states/` (US-008 ✅).
- `Bullet` exists in `internal/game/entity/actors/states/` (US-010 ✅).
- `ActiveSkill` interface exists in `internal/engine/physics/skill/skill.go` ✅.
- `body.Shooter` contract exists in `internal/engine/contracts/body/shooter.go` ✅.
- `Character.skills` field exists and is updated in `Character.Update()` ✅.
- Actor state registration system exists in `actor_state.go` ✅.

## 4. Post-conditions

- Shooting state variants registered: `IdleShooting`, `WalkingShooting`, `JumpingShooting`, `FallingShooting`.
- `ShootingSkill` implements `ActiveSkill` and triggers state transitions.
- `ShootingSkill` lives in `internal/engine/physics/skill/`.
- `OffsetToggler` lives in `internal/engine/physics/skill/`.
- `Bullet` lives in `internal/engine/entity/projectiles/`.
- `GroundedState` no longer contains shooting-specific logic.
- `GroundedInput` no longer has `ShootHeld()` method.
- Sprite system can map shooting states to distinct sprite sheets (e.g., "idle_shoot.png").
- All existing shooting tests pass with no behavioral changes.
- Code coverage ≥74.6% (no regression).

---

## 4. Integration Points

| Point | Detail |
|---|---|
| `internal/engine/physics/skill/skill_shooting.go` | New file; implements `ActiveSkill` |
| `internal/engine/physics/skill/offset_toggler.go` | Moved from `internal/game/entity/actors/states/` |
| `internal/engine/entity/projectiles/bullet.go` | Moved from `internal/game/entity/actors/states/` |
| `internal/game/entity/actors/states/grounded_state.go` | Remove shooting logic |
| `internal/game/entity/actors/states/grounded_input.go` | Remove `ShootHeld()` |
| `internal/game/entity/actors/states/mocks_test.go` | Remove `ShootHeldFunc` from `MockInputSource` |
| `internal/engine/physics/skill/skill_shooting_test.go` | Moved from `shooting_skill_test.go` |

---

## 5. Red Phase — Failing Test Scenarios

File: `internal/engine/physics/skill/skill_shooting_test.go` (moved from game layer)

### Test 1 — Cooldown gating (continuous hold)

```
GIVEN ShootingSkill with cooldown=3 frames
WHEN HandleInput is called 4 consecutive frames with shoot key held
THEN SpawnBullet is called on frame 1 and frame 4 (total=2), not on frames 2 or 3
```

### Test 2 — Alternating Y-offset over ≥4 shots

```
GIVEN ShootingSkill with cooldown=0, yOffset=4
WHEN HandleInput is called 4 times (cooldown=0 so each call spawns)
THEN the y16 argument to SpawnBullet alternates: +4, -4, +4, -4
```

### Test 3 — State transitions (Ready → Cooldown → Ready)

```
GIVEN ShootingSkill with cooldown=2
WHEN HandleInput called (spawns), then Update called 2 times
THEN state transitions: StateReady → StateCooldown → StateReady
```

### Test 4 — No spawn when state is not Ready

```
GIVEN ShootingSkill in StateCooldown
WHEN HandleInput is called with shoot key held
THEN SpawnBullet is NOT called
```

---

## 6. Files to Create / Modify / Delete

| Action | Path |
|---|---|
| Create | `internal/engine/physics/skill/skill_shooting.go` |
| Create | `internal/engine/physics/skill/skill_shooting_test.go` |
| Move | `internal/game/entity/actors/states/offset_toggler.go` → `internal/engine/physics/skill/offset_toggler.go` |
| Move | `internal/game/entity/actors/states/offset_toggler_test.go` → `internal/engine/physics/skill/offset_toggler_test.go` |
| Move | `internal/game/entity/actors/states/bullet.go` → `internal/engine/entity/projectiles/bullet.go` |
| Modify | `internal/game/entity/actors/states/grounded_state.go` (remove shooting logic) |
| Modify | `internal/game/entity/actors/states/grounded_input.go` (remove `ShootHeld()`) |
| Modify | `internal/game/entity/actors/states/mocks_test.go` (remove `ShootHeldFunc`) |
| Delete | `internal/game/entity/actors/states/shooting_skill.go` |
| Delete | `internal/game/entity/actors/states/shooting_skill_test.go` |

---

## 7. Migration Strategy

1. **Phase 1 (RED):** Move files and update imports; tests will fail due to missing `ActiveSkill` implementation.
2. **Phase 2 (GREEN):** Implement `HandleInput()`, `Update()`, `ActivationKey()` in `ShootingSkill`; all tests pass.
3. **Phase 3 (REFACTOR):** Remove shooting logic from `GroundedState`; verify no regressions.

---

## 8. Behavioral Invariants (Must Not Change)

- Bullets spawn at the same rate (cooldown enforcement).
- Y-offset alternates identically to US-010.
- Shooting is suppressed while dashing (handled by skill priority or state machine).
- Bullet collision and removal behavior unchanged.

---

## 9. Notes

- `HandleInput()` uses `ebiten.IsKeyPressed()` (continuous hold) instead of `inpututil.IsKeyJustPressed()` (single press).
- The `ActivationKey()` method returns the shoot key, but the actual input check is in `HandleInput()`.
- `OffsetToggler` and `Bullet` are now reusable engine components.
- This refactor does NOT introduce a full skill manager; that can be a future story if needed.
