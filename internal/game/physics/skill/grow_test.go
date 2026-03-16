package skill

import (
	"testing"

	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
	engineskill "github.com/leandroatallah/firefly/internal/engine/physics/skill"
	"github.com/leandroatallah/firefly/internal/engine/physics/space"
)

type mockScalablePlayer struct {
	*bodyphysics.ObstacleRect
	scale float64
}

func (m *mockScalablePlayer) SetScale(s float64) {
	m.scale = s
}

func (m *mockScalablePlayer) RefreshCollisions() {}

func TestGrowSkill_ActivateDeactivate(t *testing.T) {
	// Setup
	sp := space.NewSpace()

	// Create Player with mock scalable capability
	rect := bodyphysics.NewRect(0, 0, 10, 10)
	obs := bodyphysics.NewObstacleRect(rect)
	player := &mockScalablePlayer{
		ObstacleRect: obs,
		scale:        1.0,
	}
	player.SetID("player")
	sp.AddBody(player)

	// Create Skill
	skill := NewGrowSkill()
	// Reduce duration for testing
	skill.duration = 2

	// Verify Initial State
	if player.GetShape().Width() != 10 {
		t.Errorf("Initial width should be 10, got %d", player.GetShape().Width())
	}
	if player.scale != 1.0 {
		t.Errorf("Initial scale should be 1.0, got %f", player.scale)
	}

	// Request Activation
	skill.RequestActivation()
	skill.HandleInput(player, nil, sp)

	// Verify Activated State
	if !skill.IsActive() {
		t.Error("Skill should be active")
	}

	// Verify Size Doubled
	if player.GetShape().Width() != 20 {
		t.Errorf("Width should be 20, got %d", player.GetShape().Width())
	}
	if player.GetShape().Height() != 20 {
		t.Errorf("Height should be 20, got %d", player.GetShape().Height())
	}

	// Verify Scale Doubled
	if player.scale != 2.0 {
		t.Errorf("Scale should be 2.0, got %f", player.scale)
	}

	// Update to expire timer (Deactivate)
	skill.Update(player, nil) // timer 1
	skill.Update(player, nil) // timer 0 -> deactivate

	// Verify Deactivated State
	if skill.state != engineskill.StateCooldown {
		t.Errorf("Skill should be in cooldown, got %v", skill.state)
	}

	// Verify Size Restored
	if player.GetShape().Width() != 10 {
		t.Errorf("Width should be restored to 10, got %d", player.GetShape().Width())
	}
	if player.GetShape().Height() != 10 {
		t.Errorf("Height should be restored to 10, got %d", player.GetShape().Height())
	}

	// Verify Scale Restored
	if player.scale != 1.0 {
		t.Errorf("Scale should be restored to 1.0, got %f", player.scale)
	}
}
