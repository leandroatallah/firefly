package actors_test

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
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
