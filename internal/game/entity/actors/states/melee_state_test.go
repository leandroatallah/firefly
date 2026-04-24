package gamestates_test

import (
	"image"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/combat/weapon"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	gamestates "github.com/boilerplate/ebiten-template/internal/game/entity/actors/states"
	"github.com/hajimehoshi/ebiten/v2"
)

// ---------------------------------------------------------------------------
// Test doubles (package-local to avoid colliding with other state tests).
// ---------------------------------------------------------------------------

// meleeEnemy is a Damageable + Factioned + Collidable target for
// MeleeAttackState tests (distinct from meleeTarget in the weapon_test package).
type meleeEnemy struct {
	id          string
	pos         image.Rectangle
	faction     combat.Faction
	damageCalls []int
}

func newMeleeEnemy(id string, rect image.Rectangle, f combat.Faction) *meleeEnemy {
	return &meleeEnemy{id: id, pos: rect, faction: f}
}

func (e *meleeEnemy) TakeDamage(amount int)                           { e.damageCalls = append(e.damageCalls, amount) }
func (e *meleeEnemy) Faction() combat.Faction                         { return e.faction }
func (e *meleeEnemy) ID() string                                      { return e.id }
func (e *meleeEnemy) SetID(string)                                    {}
func (e *meleeEnemy) Position() image.Rectangle                       { return e.pos }
func (e *meleeEnemy) SetPosition(int, int)                            {}
func (e *meleeEnemy) SetPosition16(int, int)                          {}
func (e *meleeEnemy) SetSize(int, int)                                {}
func (e *meleeEnemy) Scale() float64                                  { return 1 }
func (e *meleeEnemy) SetScale(float64)                                {}
func (e *meleeEnemy) GetPosition16() (int, int)                       { return e.pos.Min.X * 16, e.pos.Min.Y * 16 }
func (e *meleeEnemy) GetPositionMin() (int, int)                      { return e.pos.Min.X, e.pos.Min.Y }
func (e *meleeEnemy) GetShape() body.Shape                            { return nil }
func (e *meleeEnemy) Owner() interface{}                              { return nil }
func (e *meleeEnemy) SetOwner(interface{})                            {}
func (e *meleeEnemy) LastOwner() interface{}                          { return nil }
func (e *meleeEnemy) GetTouchable() body.Touchable                    { return nil }
func (e *meleeEnemy) DrawCollisionBox(*ebiten.Image, image.Rectangle) {}
func (e *meleeEnemy) CollisionPosition() []image.Rectangle            { return []image.Rectangle{e.pos} }
func (e *meleeEnemy) CollisionShapes() []body.Collidable              { return nil }
func (e *meleeEnemy) IsObstructive() bool                             { return false }
func (e *meleeEnemy) SetIsObstructive(bool)                           {}
func (e *meleeEnemy) AddCollision(...body.Collidable)                 {}
func (e *meleeEnemy) ClearCollisions()                                {}
func (e *meleeEnemy) SetTouchable(body.Touchable)                     {}
func (e *meleeEnemy) OnTouch(body.Collidable)                         {}
func (e *meleeEnemy) OnBlock(body.Collidable)                         {}
func (e *meleeEnemy) ApplyValidPosition(int, bool, body.BodiesSpace) (int, int, bool) {
	return e.pos.Min.X, e.pos.Min.Y, false
}

// owningMockBody wraps mockBody and adds a Faction() getter so the MeleeWeapon
// can gate same-faction targets when used as an owner.
type owningMockBody struct {
	*mockBody
	faction combat.Faction
}

func (o *owningMockBody) Faction() combat.Faction { return o.faction }

func newPlayerOwner(xPx, yPx int, face animation.FacingDirectionEnum) *owningMockBody { //nolint:unparam
	b := &mockBody{
		pos:      image.Rect(xPx, yPx, xPx+16, yPx+24),
		grounded: true,
		faceDir:  face,
	}
	return &owningMockBody{mockBody: b, faction: combat.FactionPlayer}
}

// ---------------------------------------------------------------------------
// §4 RED-3 helpers and tests
// ---------------------------------------------------------------------------

// newTestMeleeWeaponForState builds a single-step combo-capable weapon.
// Matches the pre-US-041 defaults (damage=1, active=[3,5], cooldown=20,
// hitbox 24x16 offset 12,0).
func newTestMeleeWeaponForState(owner interface{}) *weapon.MeleeWeapon {
	steps := []weapon.ComboStep{{
		Damage:          1,
		ActiveFrames:    [2]int{3, 5},
		HitboxW16:       24 * 16,
		HitboxH16:       16 * 16,
		HitboxOffsetX16: 12 * 16,
		HitboxOffsetY16: 0,
	}}
	w := weapon.NewMeleeWeapon("player_melee", 20, 0 /*comboWindowFrames*/, steps)
	w.SetOwner(owner)
	return w
}

// newThreeStepStateMeleeWeapon builds a 3-step weapon for combo-step tests.
func newThreeStepStateMeleeWeapon(owner interface{}) *weapon.MeleeWeapon {
	steps := []weapon.ComboStep{
		{Damage: 1, ActiveFrames: [2]int{3, 5}, HitboxW16: 24 * 16, HitboxH16: 16 * 16, HitboxOffsetX16: 12 * 16, HitboxOffsetY16: 0},
		{Damage: 1, ActiveFrames: [2]int{3, 5}, HitboxW16: 28 * 16, HitboxH16: 16 * 16, HitboxOffsetX16: 14 * 16, HitboxOffsetY16: -4 * 16},
		{Damage: 2, ActiveFrames: [2]int{3, 5}, HitboxW16: 32 * 16, HitboxH16: 20 * 16, HitboxOffsetX16: 16 * 16, HitboxOffsetY16: 0},
	}
	w := weapon.NewMeleeWeapon("player_melee", 0 /*cooldown*/, 15 /*window*/, steps)
	w.SetOwner(owner)
	return w
}

func TestMeleeAttackState_ReturnsToGrounded_WhenAnimationFinishes(t *testing.T) {
	owner := newPlayerOwner(100, 100, animation.FaceDirectionRight)
	sp := &mockSpace{}
	w := newTestMeleeWeaponForState(owner)

	// US-041: ClimberPlayer owns Fire. Tests must fire explicitly before OnStart.
	w.Fire(owner.pos.Min.X*16, owner.pos.Min.Y*16, owner.faceDir, body.ShootDirectionStraight, 0)

	const animFrames = 8
	st := gamestates.NewMeleeAttackState(owner, sp, w, gamestates.StateGrounded)
	st.SetAnimationFrames(animFrames)
	st.OnStart(0)

	for i := 0; i < animFrames-1; i++ {
		if got := st.Update(); got != gamestates.StateMeleeAttack {
			t.Errorf("tick %d Update() = %v, want StateMeleeAttack", i, got)
		}
	}
	if got := st.Update(); got != gamestates.StateGrounded {
		t.Errorf("final tick Update() = %v, want StateGrounded", got)
	}
}

func TestMeleeAttackState_AirMelee_ReturnsToFalling(t *testing.T) {
	owner := newPlayerOwner(100, 100, animation.FaceDirectionRight)
	owner.grounded = false
	sp := &mockSpace{}
	w := newTestMeleeWeaponForState(owner)

	w.Fire(owner.pos.Min.X*16, owner.pos.Min.Y*16, owner.faceDir, body.ShootDirectionStraight, 0)

	const animFrames = 6
	st := gamestates.NewMeleeAttackState(owner, sp, w, actors.Falling)
	st.SetAnimationFrames(animFrames)
	st.OnStart(0)

	var got actors.ActorStateEnum
	for i := 0; i < animFrames; i++ {
		got = st.Update()
	}
	if got != actors.Falling {
		t.Errorf("air-melee final Update() = %v, want actors.Falling (not StateGrounded)", got)
	}
}

// US-041: OnStart no longer calls Fire — the climber owns that call.
// The state's job is to drive the hitbox while the weapon swings.
func TestMeleeAttackState_Update_AppliesHitboxDuringActiveWindow(t *testing.T) {
	owner := newPlayerOwner(100, 100, animation.FaceDirectionRight)
	enemy := newMeleeEnemy("enemy", image.Rect(110, 100, 118, 108), combat.FactionEnemy)

	sp := &mockSpace{queryResult: []body.Collidable{enemy}}
	w := newTestMeleeWeaponForState(owner)

	// Fire explicitly before entering the state (per SPEC §1.4).
	w.Fire(owner.pos.Min.X*16, owner.pos.Min.Y*16, owner.faceDir, body.ShootDirectionStraight, 0)

	const animFrames = 10
	st := gamestates.NewMeleeAttackState(owner, sp, w, gamestates.StateGrounded)
	st.SetAnimationFrames(animFrames)
	st.OnStart(0)

	for i := 0; i < animFrames; i++ {
		st.Update()
	}

	if len(enemy.damageCalls) != 1 {
		t.Errorf("TakeDamage called %d times over full swing, want exactly 1 (single-hit per swing)", len(enemy.damageCalls))
	}
}

func TestGroundedState_MeleePressed_TransitionsToMeleeAttack(t *testing.T) {
	input := &MockInputSource{
		MeleePressedFunc: func() bool { return true },
	}
	g := newGroundedState(input)
	g.OnStart(0)

	next := g.Update()

	if next != gamestates.StateMeleeAttack {
		t.Errorf("melee pressed: Update() = %v, want StateMeleeAttack", next)
	}
}

func TestGroundedState_DashPressed_TakesPrecedenceOverMelee(t *testing.T) {
	input := &MockInputSource{
		DashPressedFunc:  func() bool { return true },
		MeleePressedFunc: func() bool { return true },
	}
	g := newGroundedState(input)
	g.OnStart(0)

	next := g.Update()

	if next != gamestates.StateDashing {
		t.Errorf("dash+melee pressed: Update() = %v, want StateDashing (dash precedence)", next)
	}
}

func TestMeleeTrigger_BlockedDuringCooldown(t *testing.T) {
	owner := newPlayerOwner(100, 100, animation.FaceDirectionRight)
	w := newTestMeleeWeaponForState(owner)

	w.SetCooldown(15)
	if w.CanFire() {
		t.Fatalf("precondition: weapon must not be able to fire during cooldown")
	}

	_, ok := gamestates.TryMeleeFromFalling(w, true /*meleePressed*/)
	if ok {
		t.Errorf("TryMeleeFromFalling returned ok=true while weapon on cooldown; want false")
	}
}

// AC5 — The state captures the step index at OnStart so the animation layer
// can pick MeleeAttack1 / MeleeAttack2 / MeleeAttack3.
func TestMeleeAttackState_UsesCurrentComboStep(t *testing.T) {
	owner := newPlayerOwner(100, 100, animation.FaceDirectionRight)
	sp := &mockSpace{}
	w := newThreeStepStateMeleeWeapon(owner)

	// Step 0.
	w.Fire(owner.pos.Min.X*16, owner.pos.Min.Y*16, owner.faceDir, body.ShootDirectionStraight, 0)
	st0 := gamestates.NewMeleeAttackState(owner, sp, w, gamestates.StateGrounded)
	st0.SetAnimationFrames(8)
	st0.OnStart(0)
	if got := st0.StepUsed(); got != 0 {
		t.Errorf("step 0: StepUsed() = %d, want 0", got)
	}

	// Finish swing so combo window opens, then advance to step 1.
	for i := 0; i <= 5+1; i++ {
		w.Update()
	}
	if !w.AdvanceCombo() {
		t.Fatalf("AdvanceCombo step 0→1 returned false; want true")
	}
	w.Fire(owner.pos.Min.X*16, owner.pos.Min.Y*16, owner.faceDir, body.ShootDirectionStraight, 0)
	st1 := gamestates.NewMeleeAttackState(owner, sp, w, gamestates.StateGrounded)
	st1.SetAnimationFrames(8)
	st1.OnStart(0)
	if got := st1.StepUsed(); got != 1 {
		t.Errorf("step 1: StepUsed() = %d, want 1", got)
	}

	// Advance to step 2.
	for i := 0; i <= 5+1; i++ {
		w.Update()
	}
	if !w.AdvanceCombo() {
		t.Fatalf("AdvanceCombo step 1→2 returned false; want true")
	}
	w.Fire(owner.pos.Min.X*16, owner.pos.Min.Y*16, owner.faceDir, body.ShootDirectionStraight, 0)
	st2 := gamestates.NewMeleeAttackState(owner, sp, w, gamestates.StateGrounded)
	st2.SetAnimationFrames(8)
	st2.OnStart(0)
	if got := st2.StepUsed(); got != 2 {
		t.Errorf("step 2: StepUsed() = %d, want 2", got)
	}
}

// AC3 bullet 3 — Pressing Dash or Jump while the combo window is open resets
// the combo. Verified via the shared helper used by ClimberPlayer.Update to
// keep the unit test independent of the full climber harness.
func TestResetComboOnInterrupt_ResetsWhenDashOrJumpPressedDuringWindow(t *testing.T) {
	owner := newPlayerOwner(100, 100, animation.FaceDirectionRight)

	tests := []struct {
		name        string
		dashPressed bool
		jumpPressed bool
		wantReset   bool
	}{
		{name: "dash pressed during window resets combo", dashPressed: true, wantReset: true},
		{name: "jump pressed during window resets combo", jumpPressed: true, wantReset: true},
		{name: "neither pressed leaves combo intact", wantReset: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := newThreeStepStateMeleeWeapon(owner)

			// Open the combo window: fire step 0, advance past active frames.
			w.Fire(owner.pos.Min.X*16, owner.pos.Min.Y*16, owner.faceDir, body.ShootDirectionStraight, 0)
			for i := 0; i <= 5+1; i++ {
				w.Update()
			}
			// Advance to step 1 so we have something to reset from.
			if !w.AdvanceCombo() {
				t.Fatalf("precondition: AdvanceCombo failed; window not open?")
			}
			// Reopen the window by running step 1's swing to completion.
			w.Fire(owner.pos.Min.X*16, owner.pos.Min.Y*16, owner.faceDir, body.ShootDirectionStraight, 0)
			for i := 0; i <= 5+1; i++ {
				w.Update()
			}
			if w.ComboWindowRemaining() == 0 {
				t.Fatalf("precondition: combo window must be open")
			}
			if w.StepIndex() != 1 {
				t.Fatalf("precondition: StepIndex = %d, want 1", w.StepIndex())
			}

			gamestates.ResetComboOnInterrupt(w, tc.dashPressed, tc.jumpPressed)

			if tc.wantReset {
				if w.StepIndex() != 0 {
					t.Errorf("StepIndex() = %d, want 0 (reset)", w.StepIndex())
				}
				if w.ComboWindowRemaining() != 0 {
					t.Errorf("ComboWindowRemaining() = %d, want 0 (reset)", w.ComboWindowRemaining())
				}
			} else {
				if w.StepIndex() != 1 {
					t.Errorf("StepIndex() = %d, want 1 (no reset)", w.StepIndex())
				}
				if w.ComboWindowRemaining() == 0 {
					t.Errorf("ComboWindowRemaining() = 0, want > 0 (no reset)")
				}
			}
		})
	}
}
