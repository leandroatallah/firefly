package skill

import (
	"testing"

	engineskill "github.com/leandroatallah/firefly/internal/engine/physics/skill"
)

func TestStarSkill_Activation(t *testing.T) {
	skill := NewStarSkill()

	skill.RequestActivation()
	skill.HandleInput(nil, nil, nil)

	if !skill.IsActive() {
		t.Error("Skill should be active after activation")
	}
}

func TestStarSkill_UpdateSequence(t *testing.T) {
	skill := NewStarSkill()
	skill.duration = 3
	skill.cooldown = 2

	vfxCount := 0
	skill.OnActive = func() {
		vfxCount++
	}

	skill.RequestActivation()
	skill.HandleInput(nil, nil, nil)

	// Update 1
	skill.Update(nil, nil)
	if vfxCount != 1 {
		t.Errorf("Expected vfxCount 1, got %d", vfxCount)
	}

	// Update 2
	skill.Update(nil, nil)
	if vfxCount != 2 {
		t.Errorf("Expected vfxCount 2, got %d", vfxCount)
	}

	// Update 3 -> Should deactivate
	skill.Update(nil, nil)
	if skill.IsActive() {
		t.Error("Skill should be inactive after duration")
	}
	if skill.state != engineskill.StateCooldown {
		t.Errorf("Expected state Cooldown, got %v", skill.state)
	}

	// Update during cooldown
	skill.Update(nil, nil) // timer 1
	if skill.state != engineskill.StateCooldown {
		t.Errorf("Still should be in cooldown, got %v", skill.state)
	}
	skill.Update(nil, nil) // timer 0 -> Ready
	if skill.state != engineskill.StateReady {
		t.Errorf("Expected state Ready after cooldown, got %v", skill.state)
	}
}

func TestStarSkill_ResetFromDifferentStates(t *testing.T) {
	// From Active
	skill := NewStarSkill()
	skill.RequestActivation()
	skill.HandleInput(nil, nil, nil)
	skill.Reset()
	if skill.state != engineskill.StateReady {
		t.Errorf("Reset from Active failed, got %v", skill.state)
	}

	// From Cooldown
	skill = NewStarSkill()
	skill.duration = 1
	skill.RequestActivation()
	skill.HandleInput(nil, nil, nil)
	skill.Update(nil, nil) // state Cooldown
	if skill.state != engineskill.StateCooldown {
		t.Fatalf("Failed to enter cooldown, got %v", skill.state)
	}
	skill.Reset()
	if skill.state != engineskill.StateReady {
		t.Errorf("Reset from Cooldown failed, got %v", skill.state)
	}
}
