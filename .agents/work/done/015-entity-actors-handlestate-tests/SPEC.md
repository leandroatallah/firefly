# 015 — Entity Actors handleState Test Coverage

**Branch:** `015-entity-actors-handlestate-tests`
**Bounded Context:** Entity
**Package:** `internal/engine/entity/actors`

## Technical Requirements

### Core Function Under Test
`Character.handleState()` in `internal/engine/entity/actors/character.go` (lines ~200–280)

The function manages state machine transitions for an Actor across 10 states:
- **Idle, Walking, Jumping, Falling, Landing, Hurted, Dying, Dead, Exiting, Ducking**

### State Machine Logic

#### Early-Exit States (no transitions allowed)
- `Exiting`, `Dying`, `Dead` → return immediately

#### Health-Driven Transition
- If `Health() <= 0` and state ≠ `Dying` → transition to `Dying`

#### StateTransitionHandler Override
- If `StateTransitionHandler` is set and returns `true` → skip all default logic

#### Invulnerability Timer
- Decrement `invulnerabilityTimer` each frame
- When timer reaches 0 → call `SetInvulnerability(false)`

#### Standard State Transitions
- `Dying` + animation finished → `Dead`
- `Hurted` + animation finished → `Idle`
- `Landing` + walking → `Walking`; + animation finished → `Idle`
- `Jumping` + animation finished → `Idle`
- `Falling` + not falling → `Landing`
- Not falling + is falling → `Falling`
- Not jumping + is going up → `Jumping`
- `Ducking` + not ducking → `Idle`
- Not ducking + is ducking → `Ducking`
- Not walking + is walking → `Walking`
- Not idle + is idle → `Idle`

### Contracts & Interfaces

**Existing Contracts Used:**
- `body.BodiesSpace` — physics world (passed to `Update()`)
- `animation.FaceDirection` — sprite facing direction

**No new contracts required** — all state transitions use existing `Character` methods:
- `IsAnimationFinished()` — checks if current state animation is done
- `IsWalking()`, `IsFalling()`, `IsGoingUp()`, `IsDucking()`, `IsIdle()` — movement state queries
- `Health()`, `SetInvulnerability()` — health/invulnerability state
- `SetNewStateFatal()` — transition to new state

### Pre-Conditions
- Character is initialized with a valid sprite map and body rect
- Character has a valid initial state (Idle)
- Health and invulnerability state are set up

### Post-Conditions
- State transitions occur only when conditions are met
- Invulnerability timer decrements and resets correctly
- StateTransitionHandler can override default behavior
- Early-exit states prevent further transitions

## Red Phase: Failing Test Scenario

### Test: `TestCharacter_handleState_AllTransitions`

**Scenario:** Table-driven test covering all 15+ conditional branches in `handleState()`.

**Failing Test Description:**
```
Test Case: Dying → Dead transition
  Given: Character in Dying state with animation finished
  When: handleState() is called
  Then: Character should transition to Dead state
  
Test Case: Hurted → Idle transition
  Given: Character in Hurted state with animation finished
  When: handleState() is called
  Then: Character should transition to Idle state

Test Case: StateTransitionHandler override
  Given: Character with StateTransitionHandler returning true
  When: handleState() is called
  Then: Default state logic should be skipped

Test Case: Invulnerability timer decrement
  Given: Character with invulnerabilityTimer = 5
  When: handleState() is called 5 times
  Then: invulnerabilityTimer should reach 0 and SetInvulnerability(false) called

Test Case: Health <= 0 forces Dying
  Given: Character in Idle state with Health = 0
  When: handleState() is called
  Then: Character should transition to Dying state

Test Case: Early-exit states (Exiting, Dying, Dead)
  Given: Character in Exiting/Dying/Dead state
  When: handleState() is called
  Then: No state transitions should occur

Test Case: Falling → Landing transition
  Given: Character in Falling state, IsFalling() returns false
  When: handleState() is called
  Then: Character should transition to Landing state

Test Case: Landing with walking
  Given: Character in Landing state with IsWalking() = true
  When: handleState() is called
  Then: Character should transition to Walking state

Test Case: Jumping → Idle transition
  Given: Character in Jumping state with animation finished
  When: handleState() is called
  Then: Character should transition to Idle state

Test Case: Ducking → Idle transition
  Given: Character in Ducking state with IsDucking() = false
  When: handleState() is called
  Then: Character should transition to Idle state
```

## Integration Points

### Within Entity Bounded Context
- `Character` state machine is the core of actor behavior
- State transitions trigger `OnStateChange` callback (if set)
- Collision shapes refresh on state change via `RefreshCollisions()`

### External Dependencies
- Physics queries: `IsFalling()`, `IsWalking()`, `IsGoingUp()`, `IsDucking()`, `IsIdle()`
- Health system: `Health()`, `SetInvulnerability()`
- Animation system: `IsAnimationFinished()`

## Coverage Goals

- **Current:** `handleState` at 19.5% coverage
- **Target:** `internal/engine/entity/actors` package ≥ 60% coverage
- **Method-level:** All 15+ branches in `handleState()` covered by at least one test

## Test Strategy

1. **Table-driven tests** for all state transitions
2. **Mock movement queries** (IsFalling, IsWalking, etc.) via test helpers
3. **Mock animation state** via `IsAnimationFinished()` override
4. **No `time.Sleep`** — use frame counters and virtual state
5. **Deterministic** — all tests use fixed inputs and expected outputs
6. **No GPU calls** — use headless `ebiten.NewImage()` for sprites

## Key Design Decisions

- Tests will use a **test helper** to create a Character with mocked movement/animation state
- **StateTransitionHandler** will be tested by setting it and verifying it's called and can skip default logic
- **Invulnerability timer** will be tested by directly manipulating the timer and calling `handleState()` multiple times
- **Early-exit states** will be tested by setting state and verifying no transitions occur
- **Health-driven transitions** will be tested by setting health to 0 and verifying Dying state is entered
