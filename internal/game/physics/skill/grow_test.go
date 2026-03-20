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
	// Reduce duration and cooldown for testing
	skill.duration = 2
	skill.cooldown = 2

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

	// Note: Size is NOT restored immediately in the actual implementation.
	// The deactivate() method sets the Shrinking state, and size restoration
	// happens when the Shrinking animation completes via climberStateTransitionLogic.
	// This test verifies the skill state machine, not the animation-driven size restoration.

	// Update to expire cooldown (need 5 updates since cooldown=5)
	for i := 0; i < 5; i++ {
		skill.Update(player, nil)
	}

	// Verify skill returns to Ready state after cooldown
	if skill.state != engineskill.StateReady {
		t.Errorf("Skill should be ready after cooldown, got %v", skill.state)
	}
}

// TestGrowSkill_ActivateTwice verifies that the grow skill can be activated,
// goes through cooldown, and can be activated again.
func TestGrowSkill_ActivateTwice(t *testing.T) {
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
	// Reduce duration and cooldown for testing
	skill.duration = 2
	skill.cooldown = 2

	// First activation
	skill.RequestActivation()
	skill.HandleInput(player, nil, sp)

	if !skill.IsActive() {
		t.Fatal("First activation: skill should be active")
	}
	if player.GetShape().Width() != 20 {
		t.Errorf("First activation: width should be 20, got %d", player.GetShape().Width())
	}

	// Expire active duration
	skill.Update(player, nil)
	skill.Update(player, nil) // Now in cooldown with timer=2

	if skill.state != engineskill.StateCooldown {
		t.Fatalf("Expected Cooldown state, got %v", skill.state)
	}

	// Simulate size restoration (normally done by Shrinking animation completion)
	// In the actual game, climberStateTransitionLogic restores size when animation finishes
	player.SetSize(10, 10)
	player.SetScale(1.0)

	// Expire cooldown
	skill.Update(player, nil)
	skill.Update(player, nil) // Now ready

	if skill.state != engineskill.StateReady {
		t.Fatalf("Expected Ready state after cooldown, got %v", skill.state)
	}

	// Second activation - this is the key test!
	skill.RequestActivation()
	skill.HandleInput(player, nil, sp)

	if !skill.IsActive() {
		t.Error("Second activation: skill should be active again")
	}
	// Size should double from restored size (10 -> 20), not from previous grown size
	if player.GetShape().Width() != 20 {
		t.Errorf("Second activation: width should be 20, got %d", player.GetShape().Width())
	}
}

// TestGrowSkill_CanActivateDuringCooldown verifies that collecting a power-up
// during cooldown resets the cooldown and activates the skill immediately.
func TestGrowSkill_CanActivateDuringCooldown(t *testing.T) {
	sp := space.NewSpace()

	rect := bodyphysics.NewRect(0, 0, 10, 10)
	obs := bodyphysics.NewObstacleRect(rect)
	player := &mockScalablePlayer{
		ObstacleRect: obs,
		scale:        1.0,
	}
	player.SetID("player")
	sp.AddBody(player)

	skill := NewGrowSkill()
	skill.duration = 2
	skill.cooldown = 5 // Longer cooldown for testing

	// First activation
	skill.RequestActivation()
	skill.HandleInput(player, nil, sp)

	// Expire active duration
	skill.Update(player, nil)
	skill.Update(player, nil) // Now in cooldown with timer=5

	// Try to activate during cooldown - should now work!
	skill.RequestActivation()
	skill.HandleInput(player, nil, sp)

	// Skill should be active again, not in cooldown
	if skill.state != engineskill.StateActive {
		t.Errorf("Expected skill to be Active, got %v", skill.state)
	}
	// Width should double again (10 -> 20, but since size wasn't restored, it becomes 40)
	// In actual gameplay, the Shrinking animation would restore size before this happens
	if player.GetShape().Width() != 40 {
		t.Errorf("Width should be 40 (doubled from 20), got %d", player.GetShape().Width())
	}
}

// TestGrowSkill_RespawnItem verifies that the power-up item is respawned when the skill deactivates.
func TestGrowSkill_RespawnItem(t *testing.T) {
	sp := space.NewSpace()

	rect := bodyphysics.NewRect(0, 0, 10, 10)
	obs := bodyphysics.NewObstacleRect(rect)
	player := &mockScalablePlayer{
		ObstacleRect: obs,
		scale:        1.0,
	}
	player.SetID("player")
	sp.AddBody(player)

	// Create a mock item that tracks removal state
	mockItem := &mockScalablePlayer{
		ObstacleRect: obs,
		scale:        1.0,
	}
	mockItem.SetID("item")
	mockItem.SetPosition(100, 200)

	skill := NewGrowSkill()
	skill.duration = 2
	skill.cooldown = 2

	// Register item for respawn
	skill.RequestActivationWithItem(mockItem, sp)
	skill.HandleInput(player, nil, sp)

	// Verify item was tracked
	if skill.itemToRespawn == nil {
		t.Fatal("Expected item to be tracked for respawn")
	}

	// Expire active duration to trigger deactivate and respawn
	skill.Update(player, nil)
	skill.Update(player, nil) // Deactivates and should respawn

	// Verify item was cleared (respawned)
	if skill.itemToRespawn != nil {
		t.Error("Expected item to be cleared after respawn")
	}
}
