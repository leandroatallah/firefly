package actors_test

import (
	"image"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/tilemaplayer"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
	"github.com/boilerplate/ebiten-template/internal/engine/render/sprites"
	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	// Set up default config for tests
	config.Set(&config.AppConfig{
		Physics: config.PhysicsConfig{
			DownwardGravity: 4,
			UpwardGravity:   2,
		},
	})
}

func TestCharacter_Lifecycle(t *testing.T) {
	// Setup mock sprite map
	img := ebiten.NewImage(1, 1)
	sMap := sprites.SpriteMap{
		actors.Idle:    &sprites.Sprite{Image: img},
		actors.Walking: &sprites.Sprite{Image: img},
		actors.Hurted:  &sprites.Sprite{Image: img},
	}
	rect := bodyphysics.NewRect(0, 0, 16, 16)

	c := actors.NewCharacter(sMap, rect)
	c.SetID("hero")

	// Test ID was set correctly
	if c.ID() != "hero" {
		t.Errorf("expected ID 'hero', got %q", c.ID())
	}

	// Test Initial State
	if c.State() != actors.Idle {
		t.Errorf("expected initial state Idle, got %v", c.State())
	}

	// Test State Transition
	walkState, _ := c.NewState(actors.Walking)
	c.SetState(walkState)
	if c.State() != actors.Walking {
		t.Errorf("expected state Walking, got %v", c.State())
	}

	// Test Health & Hurt
	c.SetMaxHealth(100)
	c.SetHealth(100)
	c.Hurt(20)
	if c.Health() != 80 {
		t.Errorf("expected 80 health, got %d", c.Health())
	}
	if c.State() != actors.Hurted {
		t.Error("should have transitioned to Hurted state")
	}
	if !c.Invulnerable() {
		t.Error("should be invulnerable after taking damage")
	}

	// Test Movement Blockers
	if c.IsMovementBlocked() {
		t.Error("movement should not be blocked initially")
	}
	c.BlockMovement()
	if !c.IsMovementBlocked() {
		t.Error("movement should be blocked")
	}
	c.UnblockMovement()
	if c.IsMovementBlocked() {
		t.Error("movement should be unblocked")
	}
}

func TestCharacter_StateTransitionOverride(t *testing.T) {
	img := ebiten.NewImage(1, 1)
	sMap := sprites.SpriteMap{actors.Idle: &sprites.Sprite{Image: img}}
	rect := bodyphysics.NewRect(0, 0, 16, 16)
	c := actors.NewCharacter(sMap, rect)
	c.SetMaxHealth(100)
	c.SetHealth(100)

	overrideCalled := false
	c.SetStateTransitionHandler(func(char *actors.Character) bool {
		overrideCalled = true
		return true // skip default logic
	})

	// handleState is unexported, but we can trigger it via Update if Freeze is false
	c.Update(nil)
	if !overrideCalled {
		t.Error("StateTransitionHandler was not called")
	}
}

type mockSkill struct{}

func (m *mockSkill) Update(actor body.MovableCollidable, model *physicsmovement.PlatformMovementModel) {
}
func (m *mockSkill) IsActive() bool { return false }

func TestCharacter_Skills(t *testing.T) {
	img := ebiten.NewImage(1, 1)
	sMap := sprites.SpriteMap{actors.Idle: &sprites.Sprite{Image: img}}
	rect := bodyphysics.NewRect(0, 0, 16, 16)
	c := actors.NewCharacter(sMap, rect)

	s := &mockSkill{}
	c.AddSkill(s)
	c.RemoveSkill(s)
	c.ClearSkills()
}

// testCharacterWithState creates a Character with all required states and sets initial state
func testCharacterWithState(initialState actors.ActorStateEnum) *actors.Character {
	img := ebiten.NewImage(1, 1)
	sMap := sprites.SpriteMap{
		actors.Idle:    &sprites.Sprite{Image: img},
		actors.Walking: &sprites.Sprite{Image: img},
		actors.Jumping: &sprites.Sprite{Image: img},
		actors.Falling: &sprites.Sprite{Image: img},
		actors.Landing: &sprites.Sprite{Image: img},
		actors.Hurted:  &sprites.Sprite{Image: img},
		actors.Dying:   &sprites.Sprite{Image: img},
		actors.Dead:    &sprites.Sprite{Image: img},
		actors.Exiting: &sprites.Sprite{Image: img},
	}
	rect := bodyphysics.NewRect(0, 0, 16, 16)
	c := actors.NewCharacter(sMap, rect)
	c.SetMaxHealth(100)
	c.SetHealth(100)

	if initialState != actors.Idle {
		state, _ := c.NewState(initialState)
		c.SetState(state)
	}
	return c
}

func TestCharacter_handleState_DyingToDead(t *testing.T) {
	c := testCharacterWithState(actors.Dying)
	c.Update(nil)
	if c.State() != actors.Dead {
		t.Errorf("expected state Dead, got %v", c.State())
	}
}

func TestCharacter_handleState_HurtedToIdle(t *testing.T) {
	c := testCharacterWithState(actors.Hurted)
	c.Update(nil)
	if c.State() != actors.Idle {
		t.Errorf("expected state Idle, got %v", c.State())
	}
}

func TestCharacter_handleState_JumpingToIdle(t *testing.T) {
	c := testCharacterWithState(actors.Jumping)
	c.Update(nil)
	if c.State() != actors.Idle {
		t.Errorf("expected state Idle, got %v", c.State())
	}
}

func TestCharacter_handleState_LandingToIdle(t *testing.T) {
	c := testCharacterWithState(actors.Landing)
	c.Update(nil)
	if c.State() != actors.Idle {
		t.Errorf("expected state Idle, got %v", c.State())
	}
}

func TestCharacter_handleState_FallingToLanding(t *testing.T) {
	c := testCharacterWithState(actors.Falling)
	if c.State() != actors.Falling {
		t.Fatalf("initial state should be Falling, got %v", c.State())
	}
	c.Update(nil)
	if c.State() != actors.Landing {
		t.Errorf("expected state Landing, got %v", c.State())
	}
}

func TestCharacter_handleState_HealthZeroForceDying(t *testing.T) {
	c := testCharacterWithState(actors.Idle)
	c.SetHealth(0)
	c.Update(nil)
	if c.State() != actors.Dying {
		t.Errorf("expected state Dying when health <= 0, got %v", c.State())
	}
}

func TestCharacter_handleState_StateTransitionHandlerOverride(t *testing.T) {
	c := testCharacterWithState(actors.Idle)
	c.SetHealth(0)

	handlerCalled := false
	c.SetStateTransitionHandler(func(char *actors.Character) bool {
		handlerCalled = true
		return true
	})

	c.Update(nil)
	if !handlerCalled {
		t.Error("StateTransitionHandler was not called")
	}
	if c.State() != actors.Idle {
		t.Errorf("expected state Idle (handler skipped default logic), got %v", c.State())
	}
}

func TestCharacter_handleState_InvulnerabilityTimerDecrement(t *testing.T) {
	c := testCharacterWithState(actors.Idle)
	c.SetMaxHealth(100)
	c.SetHealth(100)

	c.Hurt(10)
	if !c.Invulnerable() {
		t.Error("should be invulnerable after Hurt")
	}

	for i := 0; i < 150; i++ {
		c.Update(nil)
		if !c.Invulnerable() {
			break
		}
	}

	if c.Invulnerable() {
		t.Error("should not be invulnerable after timer reached 0")
	}
}

func TestCharacter_handleState_EarlyExitStates(t *testing.T) {
	tests := []struct {
		name  string
		state actors.ActorStateEnum
	}{
		{"Exiting state", actors.Exiting},
		{"Dead state", actors.Dead},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := testCharacterWithState(tt.state)
			originalState := c.State()
			c.Update(nil)
			if c.State() != originalState {
				t.Errorf("state should not change from %v, got %v", originalState, c.State())
			}
		})
	}
}

func TestCharacter_handleState_AllTransitions(t *testing.T) {
	tests := []struct {
		name          string
		initialState  actors.ActorStateEnum
		expectedState actors.ActorStateEnum
		setup         func(*actors.Character)
	}{
		{
			name:          "Dying with animation finished transitions to Dead",
			initialState:  actors.Dying,
			expectedState: actors.Dead,
			setup:         func(c *actors.Character) {},
		},
		{
			name:          "Hurted with animation finished transitions to Idle",
			initialState:  actors.Hurted,
			expectedState: actors.Idle,
			setup:         func(c *actors.Character) {},
		},
		{
			name:          "Jumping with animation finished transitions to Idle",
			initialState:  actors.Jumping,
			expectedState: actors.Idle,
			setup:         func(c *actors.Character) {},
		},
		{
			name:          "Landing with animation finished transitions to Idle",
			initialState:  actors.Landing,
			expectedState: actors.Idle,
			setup:         func(c *actors.Character) {},
		},
		{
			name:          "Falling transitions to Landing",
			initialState:  actors.Falling,
			expectedState: actors.Landing,
			setup:         func(c *actors.Character) {},
		},
		{
			name:          "Health zero forces Dying from Idle",
			initialState:  actors.Idle,
			expectedState: actors.Dying,
			setup: func(c *actors.Character) {
				c.SetHealth(0)
			},
		},
		{
			name:          "Health zero forces Dying from Walking",
			initialState:  actors.Walking,
			expectedState: actors.Dying,
			setup: func(c *actors.Character) {
				c.SetHealth(0)
			},
		},
		{
			name:          "Exiting state does not transition",
			initialState:  actors.Exiting,
			expectedState: actors.Exiting,
			setup:         func(c *actors.Character) {},
		},
		{
			name:          "Dead state does not transition",
			initialState:  actors.Dead,
			expectedState: actors.Dead,
			setup:         func(c *actors.Character) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := testCharacterWithState(tt.initialState)
			tt.setup(c)
			c.Update(nil)
			if c.State() != tt.expectedState {
				t.Errorf("expected state %v, got %v", tt.expectedState, c.State())
			}
		})
	}
}

// stubStateContributor is a test-only StateContributor.
type stubStateContributor struct {
	target actors.ActorStateEnum
	active bool
	called bool
}

func (s *stubStateContributor) ContributeState(_ actors.ActorStateEnum) (actors.ActorStateEnum, bool) {
	s.called = true
	return s.target, s.active
}

func TestCharacter_handleState_StateContributorWins(t *testing.T) {
	c := testCharacterWithState(actors.Idle)
	contrib := &stubStateContributor{target: actors.Walking, active: true}
	c.AddStateContributor(contrib)
	c.Update(nil)
	if !contrib.called {
		t.Error("contributor was not called")
	}
	if c.State() != actors.Walking {
		t.Errorf("expected Walking from contributor, got %v", c.State())
	}
}

func TestCharacter_handleState_StateContributorDefers(t *testing.T) {
	c := testCharacterWithState(actors.Idle)
	contrib := &stubStateContributor{target: actors.Walking, active: false}
	c.AddStateContributor(contrib)
	c.Update(nil)
	if !contrib.called {
		t.Error("contributor was not called")
	}
	// Contributor deferred; default logic keeps Idle
	if c.State() != actors.Idle {
		t.Errorf("expected Idle when contributor defers, got %v", c.State())
	}
}

func TestCharacter_handleState_StateTransitionHandlerBeatsContributor(t *testing.T) {
	c := testCharacterWithState(actors.Idle)
	contrib := &stubStateContributor{target: actors.Walking, active: true}
	c.AddStateContributor(contrib)
	// StateTransitionHandler returns true → short-circuits before contributors
	c.SetStateTransitionHandler(func(_ *actors.Character) bool { return true })
	c.Update(nil)
	if contrib.called {
		t.Error("contributor should not be called when StateTransitionHandler returns true")
	}
	if c.State() != actors.Idle {
		t.Errorf("expected Idle (handler won), got %v", c.State())
	}
}

func TestCharacter_handleState_ContributorSkippedDuringHurted(t *testing.T) {
	c := testCharacterWithState(actors.Idle)
	c.SetMaxHealth(10)
	c.SetHealth(10)
	c.Hurt(1) // transitions to Hurted
	contrib := &stubStateContributor{target: actors.Walking, active: true}
	c.AddStateContributor(contrib)
	c.Update(nil)
	if contrib.called {
		t.Error("contributor must not be called while Hurted animation plays")
	}
}

type stubSpace struct {
	body.BodiesSpace
}

func (s *stubSpace) GetTilemapDimensionsProvider() tilemaplayer.TilemapDimensionsProvider {
	return nil
}
func (s *stubSpace) Query(_ image.Rectangle) []body.Collidable {
	return nil
}
func (s *stubSpace) ResolveCollisions(_ body.Collidable) (bool, bool) {
	return false, false
}

func TestCharacter_handleState_JumpFlicker(t *testing.T) {
	// Setup mock sprite map
	img := ebiten.NewImage(1, 1)
	sMap := sprites.SpriteMap{
		actors.Idle:    &sprites.Sprite{Image: img},
		actors.Jumping: &sprites.Sprite{Image: img},
		actors.Falling: &sprites.Sprite{Image: img},
	}
	rect := bodyphysics.NewRect(0, 0, 16, 16)
	c := actors.NewCharacter(sMap, rect)
	c.SetMaxSpeed(100)
	c.SetMaxHealth(100)
	c.SetHealth(100)

	// Setup movement model and ensure it's airborne
	model := physicsmovement.NewPlatformMovementModel(nil)
	model.SetOnGround(false)
	c.SetMovementModel(model)

	// Set a negative vertical velocity (going up)
	c.SetVelocity(0, -100)

	// Set initial state to Jumping
	state, _ := c.NewState(actors.Jumping)
	c.SetState(state)

	// Update the character. Since IsAnimationFinished will be true (1x1 image),
	// but onGround is false, it should stay in Jumping state.
	c.Update(&stubSpace{})

	if c.State() != actors.Jumping {
		t.Errorf("expected state to stay Jumping while vy < 0 and airborne, but got %v", c.State())
	}
}

func TestCharacter_handleState_AirPeakStaysJumping(t *testing.T) {
	// Setup default config for tests
	config.Set(&config.AppConfig{
		ScreenWidth:  256,
		ScreenHeight: 240,
		Physics: config.PhysicsConfig{
			SpeedMultiplier: 1.0,
			DownwardGravity: 5,
			UpwardGravity:   4,
			MaxFallSpeed:    100,
		},
	})

	// Setup mock sprite map
	img := ebiten.NewImage(1, 1)
	sMap := sprites.SpriteMap{
		actors.Idle:    &sprites.Sprite{Image: img},
		actors.Jumping: &sprites.Sprite{Image: img},
		actors.Falling: &sprites.Sprite{Image: img},
	}
	rect := bodyphysics.NewRect(0, 0, 16, 16)
	c := actors.NewCharacter(sMap, rect)
	c.SetMaxSpeed(100)
	c.SetMaxHealth(100)
	c.SetHealth(100)

	// Setup movement model and ensure it's airborne
	model := physicsmovement.NewPlatformMovementModel(nil)
	model.SetOnGround(false)
	c.SetMovementModel(model)

	// Set a velocity that will become 0 after gravity (peak)
	// vy is -4, UpwardGravity is 4, so after UpdateMovement it will be 0.
	c.SetVelocity(0, -4)

	// Set initial state to Jumping
	state, _ := c.NewState(actors.Jumping)
	c.SetState(state)

	// Update. After UpdateMovement, vy will be 0.
	// 0 < DownwardGravity(5), so it should stay Jumping.
	c.Update(&stubSpace{})

	if c.State() != actors.Jumping {
		_, curVy := c.Velocity()
		t.Errorf("expected state to stay Jumping at air peak (vy=%d), but got %v", curVy, c.State())
	}

	// Now set vy to DownwardGravity
	c.SetVelocity(0, 5)
	c.Update(&stubSpace{})

	if c.State() != actors.Falling {
		t.Errorf("expected state to transition to Falling when vy >= threshold, but got %v", c.State())
	}
}
