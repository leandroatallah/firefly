package weapon_test

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
)

// TestMeleeWeapon_IsSwinging verifies the IsSwinging() accessor reflects the
// weapon's internal swing state across its observable lifecycle: before Fire,
// during the active swing, and after the swing completes.
func TestMeleeWeapon_IsSwinging(t *testing.T) {
	type phase struct {
		name      string
		fire      bool
		updates   int
		wantSwing bool
	}

	// newThreeStepComboWeapon has ActiveFrames[1] == 5 for every step.
	// After swingFrame > 5 (i.e. 6 Update() calls post-Fire), swinging flips false.
	tests := []phase{
		{name: "fresh weapon before Fire is not swinging", fire: false, updates: 0, wantSwing: false},
		{name: "immediately after Fire, IsSwinging is true", fire: true, updates: 0, wantSwing: true},
		{name: "mid-active-window still swinging", fire: true, updates: 3, wantSwing: true},
		{name: "at last active frame still swinging", fire: true, updates: 5, wantSwing: true},
		{name: "after swing window ends, not swinging", fire: true, updates: 7, wantSwing: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			owner := newMeleeOwner(100, 100, combat.FactionPlayer, animation.FaceDirectionRight)
			w := newThreeStepComboWeapon(owner)

			if tc.fire {
				w.Fire(owner.x16, owner.y16, owner.face, body.ShootDirectionStraight, 0)
			}
			for i := 0; i < tc.updates; i++ {
				w.Update()
			}

			if got := w.IsSwinging(); got != tc.wantSwing {
				t.Errorf("IsSwinging() = %v, want %v", got, tc.wantSwing)
			}
		})
	}
}
