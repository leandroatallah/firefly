package melee_test

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/combat/melee"
	"github.com/boilerplate/ebiten-template/internal/engine/combat/weapon"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

//nolint:gochecknoglobals
var (
	stepState0 = actors.RegisterState("test_melee_step_0", func(b actors.BaseState) actors.ActorState {
		return &actors.IdleState{BaseState: b}
	})
	stepState1 = actors.RegisterState("test_melee_step_1", func(b actors.BaseState) actors.ActorState {
		return &actors.IdleState{BaseState: b}
	})
	stepState2 = actors.RegisterState("test_melee_step_2", func(b actors.BaseState) actors.ActorState {
		return &actors.IdleState{BaseState: b}
	})
	meleeAttackEnum = actors.RegisterState("test_melee_attack", func(b actors.BaseState) actors.ActorState {
		return &actors.IdleState{BaseState: b}
	})
)

func threeStepWeapon() *weapon.MeleeWeapon {
	steps := []weapon.ComboStep{
		{Damage: 1, StartupFrames: 0, ActiveFrames: [2]int{0, 3}, HitboxW16: 16 * 16, HitboxH16: 16 * 16},
		{Damage: 1, StartupFrames: 0, ActiveFrames: [2]int{0, 3}, HitboxW16: 16 * 16, HitboxH16: 16 * 16},
		{Damage: 1, StartupFrames: 0, ActiveFrames: [2]int{0, 3}, HitboxW16: 16 * 16, HitboxH16: 16 * 16},
	}
	return weapon.NewMeleeWeapon("ctrl_test", 0, 15, steps)
}

func stepStates() []actors.ActorStateEnum {
	return []actors.ActorStateEnum{stepState0, stepState1, stepState2}
}

func newController(w *weapon.MeleeWeapon) *melee.Controller {
	return melee.New(w, nil, meleeAttackEnum, stepStates(), nil)
}

// ---------------------------------------------------------------------------
// ContributeState
// ---------------------------------------------------------------------------

func TestController_ContributeState_NotSwinging_Defers(t *testing.T) {
	w := threeStepWeapon()
	c := newController(w)

	got, ok := c.ContributeState(actors.Idle)
	if ok {
		t.Errorf("not swinging: ContributeState returned ok=true (state=%v), want false", got)
	}
}

func TestController_ContributeState_SwingingStep0_ReturnsStep0State(t *testing.T) {
	w := threeStepWeapon()
	c := newController(w)

	w.Fire(0, 0, animation.FaceDirectionRight, body.ShootDirectionStraight, 0)

	got, ok := c.ContributeState(actors.Idle)
	if !ok {
		t.Fatalf("swinging step 0: ContributeState returned ok=false, want true")
	}
	if got != stepState0 {
		t.Errorf("swinging step 0: state = %v, want %v", got, stepState0)
	}
}

// ---------------------------------------------------------------------------
// HandleInput — edge detection and buffer
// ---------------------------------------------------------------------------

func TestController_HandleInput_EdgeDetect_FiresOnce(t *testing.T) {
	w := threeStepWeapon()
	c := newController(w)

	// First frame: held=true from idle → edge → want fire
	entered := c.HandleInput(true, false, false, true, false)
	if !entered {
		t.Error("frame 1 (first press): HandleInput returned false, want true")
	}

	// Second frame: held=true again → held-prev=true → no edge → no fire
	entered = c.HandleInput(true, false, false, true, false)
	if entered {
		t.Error("frame 2 (still held): HandleInput returned true (should not re-fire while held)")
	}

	// Third frame: released then pressed again
	_ = c.HandleInput(false, false, false, true, false)
	// Weapon might be in cooldown after step, so reset for a clean test
}

func TestController_HandleInput_Buffer_FiredWhileSwinging(t *testing.T) {
	const animFrames = 8
	w := threeStepWeapon()
	c := melee.New(w, nil, meleeAttackEnum, stepStates(), func(_ int) int { return animFrames })

	// HandleInput edge-detects the press and sets animWait; State.OnStart would normally
	// call w.Fire() — simulate that directly so IsSwinging() is true.
	entered := c.HandleInput(true, false, false, true, false)
	if !entered {
		t.Fatal("precondition: first press must trigger entry")
	}
	w.Fire(0, 0, animation.FaceDirectionRight, body.ShootDirectionStraight, 0)

	// While swinging: re-press should buffer.
	_ = c.HandleInput(false, false, false, true, false) // release
	c.HandleInput(true, false, false, true, false)      // re-press → buffered

	// Drain swing so weapon stops but combo window is open.
	for w.IsSwinging() {
		w.Update()
	}

	// Drain animWait — capture the return so we detect the frame when the
	// buffer fires (which happens the tick animWait hits 0).
	for i := 0; i < animFrames+2; i++ {
		entered = c.HandleInput(false, false, false, true, false)
		if entered {
			break
		}
	}
	if !entered {
		t.Error("buffered press should fire once animWait clears and combo window is open")
	}
}

func TestController_HandleInput_DashInterrupt_ResetsCombo(t *testing.T) {
	const animFrames = 8
	w := threeStepWeapon()
	c := melee.New(w, nil, meleeAttackEnum, stepStates(), func(_ int) int { return animFrames })

	// Fire step 0 (HandleInput + weapon Fire to simulate State.OnStart).
	c.HandleInput(true, false, false, true, false)
	w.Fire(0, 0, animation.FaceDirectionRight, body.ShootDirectionStraight, 0)

	// Run swing to completion so combo window opens.
	for w.IsSwinging() {
		w.Update()
	}
	if w.ComboWindowRemaining() == 0 {
		t.Fatal("precondition: combo window must be open after step-0 swing")
	}

	// Advance combo to step 1.
	if !w.AdvanceCombo() {
		t.Fatal("precondition: AdvanceCombo must succeed while window is open")
	}
	if w.StepIndex() != 1 {
		t.Fatalf("precondition: StepIndex = %d, want 1", w.StepIndex())
	}

	// Inject a combo window for the dash-interrupt check (AdvanceCombo zeroes window).
	// Reopen by firing step 1 and letting it swing out.
	w.Fire(0, 0, animation.FaceDirectionRight, body.ShootDirectionStraight, 0)
	for w.IsSwinging() {
		w.Update()
	}

	// Dash interrupt during combo window.
	c.HandleInput(false, true, false, true, false)

	if w.StepIndex() != 0 {
		t.Errorf("after dash interrupt: StepIndex = %d, want 0", w.StepIndex())
	}
}

func TestController_OnInterrupt_ResetsAllState(t *testing.T) {
	w := threeStepWeapon()
	c := newController(w)

	// Fire to get non-zero animWait
	c.HandleInput(true, false, false, true, false)

	c.OnInterrupt()

	if c.IsBlockingMovement() {
		t.Error("after OnInterrupt: IsBlockingMovement() = true, want false")
	}
}

// ---------------------------------------------------------------------------
// IsBlockingMovement
// ---------------------------------------------------------------------------

func TestController_IsBlockingMovement_FalseWhenIdle(t *testing.T) {
	w := threeStepWeapon()
	c := newController(w)

	if c.IsBlockingMovement() {
		t.Error("fresh controller: IsBlockingMovement() = true, want false")
	}
}

func TestController_IsBlockingMovement_TrueWhileSwinging(t *testing.T) {
	w := threeStepWeapon()
	c := newController(w)

	w.Fire(0, 0, animation.FaceDirectionRight, body.ShootDirectionStraight, 0)

	if !c.IsBlockingMovement() {
		t.Error("while swinging: IsBlockingMovement() = false, want true")
	}
}

// ---------------------------------------------------------------------------
// StepCount
// ---------------------------------------------------------------------------

func TestController_StepCount(t *testing.T) {
	w := threeStepWeapon()
	c := newController(w)

	if got := c.StepCount(); got != 3 {
		t.Errorf("StepCount() = %d, want 3", got)
	}
}
