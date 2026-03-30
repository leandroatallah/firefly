package actors_test

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
	"github.com/boilerplate/ebiten-template/internal/engine/render/sprites"
)

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
func (m *mockSkill) Update(actor body.MovableCollidable, model *physicsmovement.PlatformMovementModel) {}
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
