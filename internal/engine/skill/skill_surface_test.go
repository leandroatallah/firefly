package skill_test

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/skill"
)

// TestEngineSkillPackageSurface verifies that the engine skill package exports
// the expected types, interfaces, and that SkillBase satisfies the Skill interface.
func TestEngineSkillPackageSurface(t *testing.T) {
	t.Run("SkillBase_satisfies_Skill", func(t *testing.T) {
		var _ skill.Skill = (*skill.SkillBase)(nil)
	})

	t.Run("SkillState_constants_exist", func(t *testing.T) {
		var _ = skill.StateReady
		var _ = skill.StateActive
		var _ = skill.StateCooldown
	})

	t.Run("SkillBase_IsActive_behavior", func(t *testing.T) {
		base := &skill.SkillBase{}
		if base.IsActive() {
			t.Error("new SkillBase should not be active")
		}

		base.SetState(skill.StateActive)
		if !base.IsActive() {
			t.Error("SkillBase with StateActive should be active")
		}

		base.SetState(skill.StateReady)
		if base.IsActive() {
			t.Error("SkillBase with StateReady should not be active")
		}
	})

	t.Run("SkillBase_accessor_methods", func(t *testing.T) {
		base := &skill.SkillBase{}

		base.SetDuration(100)
		if base.Duration() != 100 {
			t.Errorf("Duration() = %d, want 100", base.Duration())
		}

		base.SetCooldown(50)
		if base.Cooldown() != 50 {
			t.Errorf("Cooldown() = %d, want 50", base.Cooldown())
		}

		base.SetSpeed(200)
		if base.Speed() != 200 {
			t.Errorf("Speed() = %d, want 200", base.Speed())
		}

		base.SetTimer(10)
		if base.Timer() != 10 {
			t.Errorf("Timer() = %d, want 10", base.Timer())
		}

		base.IncTimer()
		if base.Timer() != 11 {
			t.Errorf("After IncTimer(), Timer() = %d, want 11", base.Timer())
		}
	})
}
