package melee_test

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/data/schemas"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/kit/combat/melee"
)

// ---------------------------------------------------------------------------
// recordingWeapon — test double that satisfies melee.weaponIface and records
// SetActiveFramesOverride + Fire ordering for assertions.
// ---------------------------------------------------------------------------

type recordingWeapon struct {
	// recorded
	lastOverride          *[2]int
	overrideCallCount     int
	fireCalled            bool
	overrideSetBeforeFire bool

	// configurable
	stepIndex int
}

// combat.Weapon surface
func (r *recordingWeapon) ID() string { return "recording_melee" }
func (r *recordingWeapon) Fire(_, _ int, _ animation.FacingDirectionEnum, _ body.ShootDirection, _ int) {
	if r.overrideCallCount > 0 {
		r.overrideSetBeforeFire = true
	}
	r.fireCalled = true
}
func (r *recordingWeapon) CanFire() bool          { return true }
func (r *recordingWeapon) Update()                {}
func (r *recordingWeapon) Cooldown() int          { return 0 }
func (r *recordingWeapon) SetCooldown(_ int)      {}
func (r *recordingWeapon) SetOwner(_ interface{}) {}

// weaponIface extra surface (from internal/kit/combat/melee/state.go)
func (r *recordingWeapon) IsHitboxActive() bool           { return false }
func (r *recordingWeapon) IsSwinging() bool               { return false }
func (r *recordingWeapon) IsInStartup() bool              { return false }
func (r *recordingWeapon) ApplyHitbox(_ body.BodiesSpace) {}
func (r *recordingWeapon) StepIndex() int                 { return r.stepIndex }
func (r *recordingWeapon) ComboWindowRemaining() int      { return 0 }
func (r *recordingWeapon) ResetCombo()                    {}

// SetActiveFramesOverride is the NEW method this story introduces; the test
// will fail to compile until the production weaponIface declares it.
func (r *recordingWeapon) SetActiveFramesOverride(o *[2]int) {
	if o != nil {
		// Copy to a fresh array so callers can't mutate our recording from afar.
		cp := *o
		r.lastOverride = &cp
	} else {
		r.lastOverride = nil
	}
	r.overrideCallCount++
}

// ---------------------------------------------------------------------------
// T-M1: OnStart installs override from AssetData
// ---------------------------------------------------------------------------

func TestState_OnStart_InstallsOverrideFromAssetData(t *testing.T) {
	tests := []struct {
		name              string
		assetHitboxFrames *schemas.HitboxFrameRange
		wantOverride      *[2]int
	}{
		{
			name:              "present",
			assetHitboxFrames: &schemas.HitboxFrameRange{Start: 2, End: 4},
			wantOverride:      &[2]int{2, 4},
		},
		{
			name:              "absent",
			assetHitboxFrames: nil,
			wantOverride:      nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := &recordingWeapon{stepIndex: 0}
			owner := &mockOwner{faceDir: animation.FaceDirectionRight}
			owner.SetPosition16(100, 200)
			space := &mockSpace{}

			st := melee.NewState(owner, space, rec, nil, meleeAttackEnum, actors.Idle, actors.Falling)
			st.SetAnimationFrames(10)

			assets := map[string]schemas.AssetData{
				"melee_attack_step_0": {HitboxFrames: tc.assetHitboxFrames},
			}
			stepStateName := func(_ int) string { return "melee_attack_step_0" }
			st.SetHitboxFrameResolver(assets, stepStateName)

			st.OnStart(0)

			if rec.overrideCallCount != 1 {
				t.Fatalf("overrideCallCount = %d, want 1", rec.overrideCallCount)
			}
			if !rec.overrideSetBeforeFire {
				t.Errorf("overrideSetBeforeFire = false, want true (override must be set before Fire)")
			}
			if !rec.fireCalled {
				t.Errorf("fireCalled = false, want true")
			}

			gotNil := rec.lastOverride == nil
			wantNil := tc.wantOverride == nil
			if gotNil != wantNil {
				t.Fatalf("lastOverride nil? = %v, wantNil = %v (lastOverride=%v)", gotNil, wantNil, rec.lastOverride)
			}
			if !gotNil {
				if *rec.lastOverride != *tc.wantOverride {
					t.Errorf("lastOverride = %v, want %v", *rec.lastOverride, *tc.wantOverride)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// T-M2: OnStart clears override when no resolver was installed
// ---------------------------------------------------------------------------

func TestState_OnStart_ClearsOverrideWhenResolverMissing(t *testing.T) {
	rec := &recordingWeapon{stepIndex: 0}
	owner := &mockOwner{faceDir: animation.FaceDirectionRight}
	owner.SetPosition16(100, 200)
	space := &mockSpace{}

	st := melee.NewState(owner, space, rec, nil, meleeAttackEnum, actors.Idle, actors.Falling)
	st.SetAnimationFrames(10)
	// NOTE: SetHitboxFrameResolver is intentionally NOT called.

	st.OnStart(0)

	if rec.overrideCallCount != 1 {
		t.Errorf("overrideCallCount = %d, want 1 (override must be set explicitly to nil each OnStart)", rec.overrideCallCount)
	}
	if rec.lastOverride != nil {
		t.Errorf("lastOverride = %v, want nil", rec.lastOverride)
	}
}

// ---------------------------------------------------------------------------
// T-M3: per-step independent override
// ---------------------------------------------------------------------------

func TestState_OnStart_PerStepIndependentOverride(t *testing.T) {
	rec := &recordingWeapon{stepIndex: 0}
	owner := &mockOwner{faceDir: animation.FaceDirectionRight}
	owner.SetPosition16(100, 200)
	space := &mockSpace{}

	st := melee.NewState(owner, space, rec, nil, meleeAttackEnum, actors.Idle, actors.Falling)
	st.SetAnimationFrames(10)

	assets := map[string]schemas.AssetData{
		"melee_attack_step_0": {HitboxFrames: &schemas.HitboxFrameRange{Start: 1, End: 2}},
		"melee_attack_step_1": {HitboxFrames: nil},
	}
	stepStateName := func(i int) string {
		switch i {
		case 0:
			return "melee_attack_step_0"
		case 1:
			return "melee_attack_step_1"
		default:
			return ""
		}
	}
	st.SetHitboxFrameResolver(assets, stepStateName)

	// Step 0 — has HitboxFrames {1,2}
	rec.stepIndex = 0
	st.OnStart(0)
	if rec.lastOverride == nil {
		t.Fatalf("step 0: lastOverride = nil, want &[2]int{1,2}")
	}
	override1 := *rec.lastOverride
	if override1 != [2]int{1, 2} {
		t.Errorf("step 0: lastOverride = %v, want [1 2]", override1)
	}

	// Step 1 — no HitboxFrames → must clear (nil)
	rec.stepIndex = 1
	st.OnStart(0)
	if rec.lastOverride != nil {
		t.Errorf("step 1: lastOverride = %v, want nil (step 1 has no HitboxFrames)", *rec.lastOverride)
	}
}

// ---------------------------------------------------------------------------
// T-M4: Ducking early-return must not touch the override
// ---------------------------------------------------------------------------

func TestState_OnStart_Ducking_DoesNotTouchOverride(t *testing.T) {
	rec := &recordingWeapon{stepIndex: 0}
	owner := &mockOwner{faceDir: animation.FaceDirectionRight, ducking: true}
	owner.SetPosition16(100, 200)
	space := &mockSpace{}

	st := melee.NewState(owner, space, rec, nil, meleeAttackEnum, actors.Idle, actors.Falling)
	st.SetAnimationFrames(10)

	assets := map[string]schemas.AssetData{
		"melee_attack_step_0": {HitboxFrames: &schemas.HitboxFrameRange{Start: 2, End: 4}},
	}
	stepStateName := func(_ int) string { return "melee_attack_step_0" }
	st.SetHitboxFrameResolver(assets, stepStateName)

	st.OnStart(0)

	if rec.overrideCallCount != 0 {
		t.Errorf("ducking: overrideCallCount = %d, want 0 (early-return path must not touch override)", rec.overrideCallCount)
	}
	if rec.fireCalled {
		t.Errorf("ducking: fireCalled = true, want false (early-return path must not Fire)")
	}
}
