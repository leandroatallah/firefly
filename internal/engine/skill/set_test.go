package skill_test

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
	"github.com/boilerplate/ebiten-template/internal/engine/skill"
	"github.com/hajimehoshi/ebiten/v2"
)

// stubSkill is a minimal Skill implementation for testing Set operations.
type stubSkill struct {
	skill.SkillBase
	updateCalled int
}

func (s *stubSkill) Update(body.MovableCollidable, *physicsmovement.PlatformMovementModel) {
	s.updateCalled++
}

// stubActiveSkill is a minimal ActiveSkill implementation for testing Set.Get.
type stubActiveSkill struct {
	stubSkill
	key ebiten.Key
}

func (s *stubActiveSkill) HandleInput(body.MovableCollidable, *physicsmovement.PlatformMovementModel, body.BodiesSpace) {
}

func (s *stubActiveSkill) ActivationKey() ebiten.Key {
	return s.key
}

// TestEngineSkillSetSurface verifies that the skill.Set registry type exports
// the expected methods and behaves correctly.
func TestEngineSkillSetSurface(t *testing.T) {
	t.Run("NewSet_returns_non_nil", func(t *testing.T) {
		s := skill.NewSet()
		if s == nil {
			t.Fatal("NewSet() returned nil")
		}
	})

	t.Run("ActiveCount_empty_set", func(t *testing.T) {
		s := skill.NewSet()
		if count := s.ActiveCount(); count != 0 {
			t.Errorf("ActiveCount() on empty set = %d, want 0", count)
		}
	})

	t.Run("ActiveCount_reflects_active_skills", func(t *testing.T) {
		s := skill.NewSet()

		stub1 := &stubSkill{}
		stub1.SetState(skill.StateActive)
		s.Add(stub1)

		stub2 := &stubSkill{}
		stub2.SetState(skill.StateReady)
		s.Add(stub2)

		if count := s.ActiveCount(); count != 1 {
			t.Errorf("ActiveCount() = %d, want 1 (only stub1 is active)", count)
		}

		stub2.SetState(skill.StateActive)
		if count := s.ActiveCount(); count != 2 {
			t.Errorf("ActiveCount() after activating stub2 = %d, want 2", count)
		}
	})

	t.Run("Update_invokes_all_skills", func(t *testing.T) {
		s := skill.NewSet()

		stub1 := &stubSkill{}
		stub2 := &stubSkill{}
		s.Add(stub1)
		s.Add(stub2)

		s.Update(nil, nil)

		if stub1.updateCalled != 1 {
			t.Errorf("stub1.updateCalled = %d, want 1", stub1.updateCalled)
		}
		if stub2.updateCalled != 1 {
			t.Errorf("stub2.updateCalled = %d, want 1", stub2.updateCalled)
		}
	})

	t.Run("Get_round_trips_ActiveSkill", func(t *testing.T) {
		s := skill.NewSet()

		stub := &stubActiveSkill{key: ebiten.KeySpace}
		s.Add(stub)

		retrieved, ok := s.Get(ebiten.KeySpace)
		if !ok {
			t.Fatal("Get(KeySpace) returned ok=false, want true")
		}
		if retrieved != stub {
			t.Error("Get(KeySpace) returned different skill than added")
		}

		_, ok = s.Get(ebiten.KeyEnter)
		if ok {
			t.Error("Get(KeyEnter) returned ok=true for unregistered key, want false")
		}
	})

	t.Run("All_returns_registered_skills", func(t *testing.T) {
		s := skill.NewSet()

		stub1 := &stubSkill{}
		stub2 := &stubSkill{}
		s.Add(stub1)
		s.Add(stub2)

		all := s.All()
		if len(all) != 2 {
			t.Fatalf("All() returned %d skills, want 2", len(all))
		}

		if all[0] != stub1 || all[1] != stub2 {
			t.Error("All() did not return skills in insertion order")
		}
	})
}
