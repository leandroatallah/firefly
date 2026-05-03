package kitskills_test

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/skill"
	kitskills "github.com/boilerplate/ebiten-template/internal/kit/skills"
)

// TestKitSkillsPackageSurface verifies that the kit skills package exports
// the expected concrete skill types and that they satisfy the skill contracts.
func TestKitSkillsPackageSurface(t *testing.T) {
	t.Run("HorizontalMovementSkill_satisfies_Skill", func(t *testing.T) {
		var _ skill.Skill = (*kitskills.HorizontalMovementSkill)(nil)
	})

	t.Run("JumpSkill_satisfies_ActiveSkill", func(t *testing.T) {
		var _ skill.ActiveSkill = (*kitskills.JumpSkill)(nil)
	})

	t.Run("DashSkill_satisfies_ActiveSkill", func(t *testing.T) {
		var _ skill.ActiveSkill = (*kitskills.DashSkill)(nil)
	})

	t.Run("ShootingSkill_satisfies_ActiveSkill", func(t *testing.T) {
		var _ skill.ActiveSkill = (*kitskills.ShootingSkill)(nil)
	})

	t.Run("OffsetToggler_type_exists", func(t *testing.T) {
		var _ kitskills.OffsetToggler
	})

	t.Run("FromConfig_returns_empty_slice_for_nil", func(t *testing.T) {
		result := kitskills.FromConfig(nil, kitskills.SkillDeps{})
		if result == nil {
			t.Error("FromConfig should return non-nil slice for nil config")
		}
		if len(result) != 0 {
			t.Errorf("FromConfig(nil, ...) should return empty slice, got %d skills", len(result))
		}
	})
}
