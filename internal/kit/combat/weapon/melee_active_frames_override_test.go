package weapon_test

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
	"github.com/boilerplate/ebiten-template/internal/kit/combat/weapon"
)

// newOverrideTestWeapon constructs a single-step MeleeWeapon with the canonical
// override-test shape: ActiveFrames=[3,5], 0 startup, short cooldown, 24x16
// hitbox forward of the owner. Cooldown is small so the T-W2 test can elapse it
// without manually advancing dozens of frames.
func newOverrideTestWeapon(owner interface{}) *weapon.MeleeWeapon {
	steps := []weapon.ComboStep{{
		Damage:          1,
		ActiveFrames:    [2]int{3, 5},
		HitboxW16:       24 * 16,
		HitboxH16:       16 * 16,
		HitboxOffsetX16: 12 * 16,
		HitboxOffsetY16: 0,
	}}
	w := weapon.NewMeleeWeapon("override_melee", 8 /*cooldown*/, 0 /*comboWindow*/, steps)
	w.SetOwner(owner)
	return w
}

// T-W1: IsHitboxActive honors the override when present and the step window
// otherwise. ComboStep.ActiveFrames is [3,5] for all rows; rows specify
// override and the swing-frame to probe.
func TestMeleeWeapon_IsHitboxActive_RespectsOverride(t *testing.T) {
	type probe struct {
		frame      int
		wantActive bool
	}
	tests := []struct {
		name     string
		override *[2]int
		probes   []probe
	}{
		{
			name:     "nil override, frame in step window",
			override: nil,
			probes:   []probe{{frame: 3, wantActive: true}},
		},
		{
			name:     "nil override, frame outside step window",
			override: nil,
			probes:   []probe{{frame: 6, wantActive: false}},
		},
		{
			name:     "override [1,2] uses override not step",
			override: &[2]int{1, 2},
			probes: []probe{
				{frame: 2, wantActive: true},
				{frame: 3, wantActive: false},
			},
		},
		{
			name:     "override exact Start",
			override: &[2]int{4, 7},
			probes:   []probe{{frame: 4, wantActive: true}},
		},
		{
			name:     "override exact End",
			override: &[2]int{4, 7},
			probes:   []probe{{frame: 7, wantActive: false}},
		},
		{
			name:     "override Start>End never activates",
			override: &[2]int{6, 5},
			probes: []probe{
				{frame: 5, wantActive: false},
				{frame: 6, wantActive: false},
			},
		},
		{
			name:     "override Start==End single-frame",
			override: &[2]int{4, 4},
			probes: []probe{
				{frame: 3, wantActive: false},
				{frame: 4, wantActive: true},
				{frame: 5, wantActive: false},
			},
		},
	}

	for _, tc := range tests {
		for _, p := range tc.probes {
			tc, p := tc, p
			t.Run(tc.name, func(t *testing.T) {
				owner := newMeleeOwner(100, 100, combat.FactionPlayer, animation.FaceDirectionRight)
				w := newOverrideTestWeapon(owner)

				w.SetActiveFramesOverride(tc.override)
				w.Fire(owner.x16, owner.y16, owner.face, body.ShootDirectionStraight, 0)
				for i := 0; i < p.frame; i++ {
					w.Update()
				}

				if got := w.IsHitboxActive(); got != p.wantActive {
					t.Errorf("override=%v frame=%d IsHitboxActive() = %v, want %v",
						tc.override, p.frame, got, p.wantActive)
				}
			})
		}
	}
}

// T-W2: startSwing clears a previously-installed override so the next Fire()
// without a fresh SetActiveFramesOverride uses the step's ActiveFrames again.
func TestMeleeWeapon_StartSwing_ClearsOverride(t *testing.T) {
	owner := newMeleeOwner(100, 100, combat.FactionPlayer, animation.FaceDirectionRight)
	w := newOverrideTestWeapon(owner)

	// Swing 1: install a narrow override [1,2] and fire.
	w.SetActiveFramesOverride(&[2]int{1, 2})
	w.Fire(owner.x16, owner.y16, owner.face, body.ShootDirectionStraight, 0)

	// Drain swing 1 (advance past ActiveFrames[1]=5 so swinging goes false).
	for w.IsSwinging() {
		w.Update()
	}
	// Drain cooldown so CanFire() returns true again.
	for !w.CanFire() {
		w.Update()
	}

	// Swing 2: do NOT call SetActiveFramesOverride. If startSwing cleared the
	// stale override, IsHitboxActive should follow ComboStep.ActiveFrames=[3,5].
	w.Fire(owner.x16, owner.y16, owner.face, body.ShootDirectionStraight, 0)

	// Frame 2: outside step window → must be inactive (would be ACTIVE if the
	// stale [1,2] override survived the swing transition).
	for i := 0; i < 2; i++ {
		w.Update()
	}
	if got := w.IsHitboxActive(); got != false {
		t.Errorf("swing 2 frame 2: IsHitboxActive() = %v, want false (stale override should be cleared)", got)
	}

	// Frame 3: inside step window → must be active.
	w.Update()
	if got := w.IsHitboxActive(); got != true {
		t.Errorf("swing 2 frame 3: IsHitboxActive() = %v, want true (step ActiveFrames=[3,5] should apply)", got)
	}
}

// T-W3: Shared-weapon mid-swing scenario. AC-8 requires the override to be
// cleared at the START of each swing (inside startSwing), NOT at swing end
// inside Update. This guards the case where a MeleeWeapon instance is shared
// across actors and a second actor reaches startSwing while the previous
// swing has not yet drained through Update's end-of-swing path. If clearing
// only happens at swing end, the stale override leaks into the next swing.
//
// To reach startSwing mid-swing through the public API we force cooldown to
// zero (the only blocker on Fire mid-swing); this models any path that
// re-enters startSwing without the prior swing's natural termination.
func TestMeleeWeapon_StartSwing_ClearsOverride_SharedWeaponMidSwing(t *testing.T) {
	owner := newMeleeOwner(100, 100, combat.FactionPlayer, animation.FaceDirectionRight)
	w := newOverrideTestWeapon(owner)

	// Swing 1: install narrow override [1,2] (distinct from step [3,5]).
	w.SetActiveFramesOverride(&[2]int{1, 2})
	w.Fire(owner.x16, owner.y16, owner.face, body.ShootDirectionStraight, 0)

	// Advance to mid-swing (frame 2). Swing 1 is still in flight; Update's
	// end-of-swing override-clear branch has NOT executed.
	w.Update()
	w.Update()

	// Force cooldown to zero so the next Fire() reaches startSwing while the
	// stale override is still installed (simulates a second actor sharing the
	// weapon instance firing before swing 1 fully drained).
	w.SetCooldown(0)

	// Swing 2: a different actor (or re-fire) does NOT install a new override.
	// If startSwing clears the stale override (AC-8), the new swing must
	// follow ComboStep.ActiveFrames=[3,5].
	w.Fire(owner.x16, owner.y16, owner.face, body.ShootDirectionStraight, 0)

	// Swing 2 frame 4: outside stale override [1,2], inside step window [3,5].
	//   - Buggy impl (clears only at swing end): stale [1,2] survives the
	//     re-entry into startSwing, 4 not in [1,2] → IsHitboxActive == false.
	//   - Correct impl (clears in startSwing per AC-8): override nil, step
	//     [3,5] applies, 4 in [3,5] → IsHitboxActive == true.
	for i := 0; i < 4; i++ {
		w.Update()
	}
	if got := w.IsHitboxActive(); got != true {
		t.Errorf("shared-weapon swing 2 frame 4: IsHitboxActive() = %v, want true (startSwing must clear stale override so step ActiveFrames=[3,5] applies)", got)
	}
}
