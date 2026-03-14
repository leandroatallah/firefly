package skill

import (
	"testing"

	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
	engineskill "github.com/leandroatallah/firefly/internal/engine/physics/skill"
	"github.com/leandroatallah/firefly/internal/engine/physics/space"
)

func TestFreezeSkill_Activate(t *testing.T) {
	// Setup
	sp := space.NewSpace()
	
	// Create Player
	player := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 10, 10))
	player.SetID("player")
	
	// Create Enemy 1 (Movable)
	enemy1 := bodyphysics.NewObstacleRect(bodyphysics.NewRect(20, 0, 10, 10))
	enemy1.SetID("enemy1")
	sp.AddBody(enemy1)
	
	// Create Enemy 2 (Movable)
	enemy2 := bodyphysics.NewObstacleRect(bodyphysics.NewRect(40, 0, 10, 10))
	enemy2.SetID("enemy2")
	sp.AddBody(enemy2)

	// Create Skill
	skill := NewFreezeSkill()
	
	// Request Activation
	skill.RequestActivation()
	
	// Handle Input (Triggers Activation)
	skill.HandleInput(player, nil, sp)
	
	// Verify State
	if !skill.IsActive() {
		t.Error("Skill should be active")
	}
	
	// Verify Enemies Frozen
	if !enemy1.Freeze() {
		t.Error("Enemy 1 should be frozen")
	}
	if !enemy2.Freeze() {
		t.Error("Enemy 2 should be frozen")
	}
	
	// Verify Player Not Frozen (Player wasn't in space, but if it was?)
	// Let's add player to space to be sure
	sp.AddBody(player)
	// Reset and try again
	enemy1.SetFreeze(false)
	enemy2.SetFreeze(false)
	skill = NewFreezeSkill()
	skill.RequestActivation()
	skill.HandleInput(player, nil, sp)
	
	if player.Freeze() {
		t.Error("Player should NOT be frozen")
	}
	if !enemy1.Freeze() {
		t.Error("Enemy 1 should be frozen")
	}
}

func TestFreezeSkill_Deactivate(t *testing.T) {
	sp := space.NewSpace()
	player := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 10, 10))
	player.SetID("player")
	enemy := bodyphysics.NewObstacleRect(bodyphysics.NewRect(20, 0, 10, 10))
	enemy.SetID("enemy")
	sp.AddBody(enemy)
	
	skill := NewFreezeSkill()
	skill.duration = 2 // Short duration
	skill.RequestActivation()
	skill.HandleInput(player, nil, sp)
	
	if !enemy.Freeze() {
		t.Fatal("Enemy should be frozen")
	}
	
	// Update to expire timer
	skill.Update(player, nil) // timer becomes 1
	skill.Update(player, nil) // timer becomes 0 -> deactivate
	
	if enemy.Freeze() {
		t.Error("Enemy should be unfrozen after duration")
	}
	
	if skill.state != engineskill.StateCooldown {
		t.Errorf("Skill should be in cooldown, got %v", skill.state)
	}
}

func TestFreezeSkill_AlreadyFrozen(t *testing.T) {
	sp := space.NewSpace()
	player := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 10, 10))
	player.SetID("player")
	enemy := bodyphysics.NewObstacleRect(bodyphysics.NewRect(20, 0, 10, 10))
	enemy.SetID("enemy")
	sp.AddBody(enemy)
	
	// Pre-freeze enemy (e.g. by another mechanic)
	enemy.SetFreeze(true)
	
	skill := NewFreezeSkill()
	skill.RequestActivation()
	skill.HandleInput(player, nil, sp)
	
	// Skill should NOT have added enemy to its tracking list because it was already frozen
	// We verify this by deactivating skill and ensuring enemy STAYS frozen
	
	// We can't access `deactivate` directly as it's unexported in `freeze.go`?
	// Wait, `deactivate` IS unexported. I can only trigger it via Update.
	skill.duration = 1
	skill.Update(player, nil) // Trigger deactivation
	
	if !enemy.Freeze() {
		t.Error("Enemy should remain frozen because it was frozen before skill activation")
	}
}
