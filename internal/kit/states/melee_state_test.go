package kitstates_test

import (
	"image"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/combat/weapon"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	kitstates "github.com/boilerplate/ebiten-template/internal/kit/states"
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

// owningMockBody wraps mockBody and adds Faction() + IsDucking() so it
// satisfies the meleeOwnerIface used by MeleeAttackState (after US-042).
type owningMockBody struct {
	*mockBody
	faction combat.Faction
	ducking bool
}

func (o *owningMockBody) Faction() combat.Faction { return o.faction }
func (o *owningMockBody) IsDucking() bool         { return o.ducking }

func newPlayerOwner(xPx, yPx int, face animation.FacingDirectionEnum) *owningMockBody { //nolint:unparam
	b := &mockBody{
		pos:      image.Rect(xPx, yPx, xPx+16, yPx+24),
		grounded: true,
		faceDir:  face,
	}
	return &owningMockBody{mockBody: b, faction: combat.FactionPlayer}
}

// vfxSpy records SpawnDirectionalPuff invocations. It satisfies the package-local
// meleeVFXSpawner interface introduced by US-042 SPEC §3.2.
type vfxSpy struct {
	calls []vfxSpyCall
}

type vfxSpyCall struct {
	typeKey   string
	x, y      float64
	faceRight bool
	count     int
	randRange float64
}

func (v *vfxSpy) SpawnDirectionalPuff(typeKey string, x, y float64, faceRight bool, count int, randRange float64) {
	v.calls = append(v.calls, vfxSpyCall{typeKey, x, y, faceRight, count, randRange})
}

// ---------------------------------------------------------------------------
// Test helpers — weapon factories.
// ---------------------------------------------------------------------------

// newTestMeleeWeaponForState builds a single-step combo-capable weapon.
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

// ---------------------------------------------------------------------------
// US-042 RED-1: OnStart fires the weapon (no caller-side w.Fire).
// ---------------------------------------------------------------------------

func TestMeleeAttackState_OnStart_FiresWeapon(t *testing.T) {
	owner := newPlayerOwner(100, 100, animation.FaceDirectionRight)
	sp := &mockSpace{}
	w := newTestMeleeWeaponForState(owner)

	if w.IsSwinging() {
		t.Fatalf("precondition: weapon must not be swinging before OnStart")
	}

	st := kitstates.NewMeleeAttackState(owner, sp, w, nil /*vfx*/)
	st.SetAnimationFrames(8)
	st.OnStart(0)

	if !w.IsSwinging() {
		t.Errorf("after OnStart: weapon.IsSwinging() = false, want true (state must own Fire)")
	}
}

// ---------------------------------------------------------------------------
// US-042 RED-2: OnStart spawns VFX with correct facing offset.
// ---------------------------------------------------------------------------

func TestMeleeAttackState_OnStart_SpawnsVFX(t *testing.T) {
	tests := []struct {
		name          string
		face          animation.FacingDirectionEnum
		wantFaceRight bool
		wantOffsetPx  int
	}{
		{name: "facing right spawns VFX +12px", face: animation.FaceDirectionRight, wantFaceRight: true, wantOffsetPx: 12},
		{name: "facing left spawns VFX -12px", face: animation.FaceDirectionLeft, wantFaceRight: false, wantOffsetPx: -12},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			owner := newPlayerOwner(100, 50, tc.face)
			// owner GetPosition16 returns (0, 0) per mockBody default; use known x16/y16.
			owner.pos = image.Rect(100, 50, 116, 74)
			sp := &mockSpace{}
			w := newTestMeleeWeaponForState(owner)
			spy := &vfxSpy{}

			st := kitstates.NewMeleeAttackState(owner, sp, w, spy)
			st.SetAnimationFrames(8)
			st.OnStart(0)

			if len(spy.calls) != 1 {
				t.Fatalf("SpawnDirectionalPuff calls = %d, want 1", len(spy.calls))
			}
			c := spy.calls[0]
			if c.typeKey != "melee_slash" {
				t.Errorf("typeKey = %q, want %q", c.typeKey, "melee_slash")
			}
			if c.faceRight != tc.wantFaceRight {
				t.Errorf("faceRight = %v, want %v", c.faceRight, tc.wantFaceRight)
			}

			// Owner mockBody.GetPosition16 returns (0, 0), so px/py reflect only the offset.
			// fp16.From16(0 + 12*16) = 12, fp16.From16(-12*16) = -12, fp16.From16(0) = 0.
			wantX := float64(tc.wantOffsetPx)
			wantY := 0.0
			if c.x != wantX {
				t.Errorf("x = %v, want %v", c.x, wantX)
			}
			if c.y != wantY {
				t.Errorf("y = %v, want %v", c.y, wantY)
			}
		})
	}
}

// Nil VFX must not panic.
func TestMeleeAttackState_OnStart_NilVFX_NoPanic(t *testing.T) {
	owner := newPlayerOwner(100, 100, animation.FaceDirectionRight)
	sp := &mockSpace{}
	w := newTestMeleeWeaponForState(owner)

	st := kitstates.NewMeleeAttackState(owner, sp, w, nil)
	st.SetAnimationFrames(8)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("OnStart with nil vfx panicked: %v", r)
		}
	}()
	st.OnStart(0)
}

// ---------------------------------------------------------------------------
// US-042 RED-3: OnStart aborts when owner is ducking.
// ---------------------------------------------------------------------------

func TestMeleeAttackState_OnStart_DuckingAborts(t *testing.T) {
	owner := newPlayerOwner(100, 100, animation.FaceDirectionRight)
	owner.ducking = true
	sp := &mockSpace{}
	w := newTestMeleeWeaponForState(owner)
	spy := &vfxSpy{}

	st := kitstates.NewMeleeAttackState(owner, sp, w, spy)
	st.SetAnimationFrames(8)
	st.OnStart(0)

	if w.IsSwinging() {
		t.Errorf("ducking abort: weapon.IsSwinging() = true, want false (Fire must not be called)")
	}
	if len(spy.calls) != 0 {
		t.Errorf("ducking abort: VFX spawn calls = %d, want 0", len(spy.calls))
	}

	// Next Update must resolve immediately to returnTo (StateGrounded for grounded owner).
	got := st.Update()
	if got != kitstates.StateGrounded {
		t.Errorf("ducking abort: next Update() = %v, want StateGrounded (immediate resolve)", got)
	}
}

// ---------------------------------------------------------------------------
// US-042 RED-4: returnTo computed dynamically from grounded state.
// ---------------------------------------------------------------------------

func TestMeleeAttackState_DynamicReturnTo(t *testing.T) {
	tests := []struct {
		name     string
		grounded bool
		want     actors.ActorStateEnum
	}{
		{name: "grounded owner returns to StateGrounded", grounded: true, want: kitstates.StateGrounded},
		{name: "airborne owner returns to actors.Falling", grounded: false, want: actors.Falling},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			owner := newPlayerOwner(100, 100, animation.FaceDirectionRight)
			owner.grounded = tc.grounded
			sp := &mockSpace{}
			w := newTestMeleeWeaponForState(owner)

			const animFrames = 6
			st := kitstates.NewMeleeAttackState(owner, sp, w, nil /*vfx*/)
			st.SetAnimationFrames(animFrames)
			st.OnStart(0)

			var got actors.ActorStateEnum
			for i := 0; i < animFrames; i++ {
				got = st.Update()
			}
			if got != tc.want {
				t.Errorf("final Update() = %v, want %v", got, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// US-042 RED-6: Hitbox is applied during the active window (the state owns Fire).
// ---------------------------------------------------------------------------

func TestMeleeAttackState_Update_AppliesHitboxDuringActiveWindow(t *testing.T) {
	owner := newPlayerOwner(100, 100, animation.FaceDirectionRight)
	enemy := newMeleeEnemy("enemy", image.Rect(110, 100, 118, 108), combat.FactionEnemy)

	sp := &mockSpace{queryResult: []body.Collidable{enemy}}
	w := newTestMeleeWeaponForState(owner)

	const animFrames = 10
	st := kitstates.NewMeleeAttackState(owner, sp, w, nil /*vfx*/)
	st.SetAnimationFrames(animFrames)
	st.OnStart(0) // OnStart now owns Fire — no caller-side w.Fire.

	for i := 0; i < animFrames; i++ {
		st.Update()
	}

	if len(enemy.damageCalls) != 1 {
		t.Errorf("TakeDamage called %d times over full swing, want exactly 1", len(enemy.damageCalls))
	}
}

// ---------------------------------------------------------------------------
// US-042 RED-7: Step capture order — OnStart captures StepIndex BEFORE Fire.
// ---------------------------------------------------------------------------

func TestMeleeAttackState_UsesCurrentComboStep(t *testing.T) {
	owner := newPlayerOwner(100, 100, animation.FaceDirectionRight)
	sp := &mockSpace{}
	w := newThreeStepStateMeleeWeapon(owner)

	// Step 0 — state OnStart fires.
	st0 := kitstates.NewMeleeAttackState(owner, sp, w, nil)
	st0.SetAnimationFrames(8)
	st0.OnStart(0)
	if got := st0.StepUsed(); got != 0 {
		t.Errorf("step 0: StepUsed() = %d, want 0", got)
	}

	// Drive the swing to completion via state Update so hitbox+window logic settles,
	// matching how production will run after the refactor.
	for i := 0; i < 8; i++ {
		st0.Update()
	}
	// Open the next window if not already (weapon.Update inside state advances frames).
	for w.IsSwinging() {
		w.Update()
	}
	if !w.AdvanceCombo() {
		t.Fatalf("AdvanceCombo step 0→1 returned false; window must be open")
	}

	st1 := kitstates.NewMeleeAttackState(owner, sp, w, nil)
	st1.SetAnimationFrames(8)
	st1.OnStart(0)
	if got := st1.StepUsed(); got != 1 {
		t.Errorf("step 1: StepUsed() = %d, want 1", got)
	}

	for i := 0; i < 8; i++ {
		st1.Update()
	}
	for w.IsSwinging() {
		w.Update()
	}
	if !w.AdvanceCombo() {
		t.Fatalf("AdvanceCombo step 1→2 returned false; window must be open")
	}

	st2 := kitstates.NewMeleeAttackState(owner, sp, w, nil)
	st2.SetAnimationFrames(8)
	st2.OnStart(0)
	if got := st2.StepUsed(); got != 2 {
		t.Errorf("step 2: StepUsed() = %d, want 2", got)
	}
}

// ---------------------------------------------------------------------------
// Existing regression tests (signature updated to drop static returnTo).
// ---------------------------------------------------------------------------

func TestMeleeAttackState_ReturnsToGrounded_WhenAnimationFinishes(t *testing.T) {
	owner := newPlayerOwner(100, 100, animation.FaceDirectionRight)
	sp := &mockSpace{}
	w := newTestMeleeWeaponForState(owner)

	const animFrames = 8
	st := kitstates.NewMeleeAttackState(owner, sp, w, nil)
	st.SetAnimationFrames(animFrames)
	st.OnStart(0)

	for i := 0; i < animFrames-1; i++ {
		if got := st.Update(); got != kitstates.StateMeleeAttack {
			t.Errorf("tick %d Update() = %v, want StateMeleeAttack", i, got)
		}
	}
	if got := st.Update(); got != kitstates.StateGrounded {
		t.Errorf("final tick Update() = %v, want StateGrounded", got)
	}
}

func TestMeleeAttackState_AirMelee_ReturnsToFalling(t *testing.T) {
	owner := newPlayerOwner(100, 100, animation.FaceDirectionRight)
	owner.grounded = false
	sp := &mockSpace{}
	w := newTestMeleeWeaponForState(owner)

	const animFrames = 6
	st := kitstates.NewMeleeAttackState(owner, sp, w, nil)
	st.SetAnimationFrames(animFrames)
	st.OnStart(0)

	var got actors.ActorStateEnum
	for i := 0; i < animFrames; i++ {
		got = st.Update()
	}
	if got != actors.Falling {
		t.Errorf("air-melee final Update() = %v, want actors.Falling", got)
	}
}

func TestGroundedState_MeleePressed_TransitionsToMeleeAttack(t *testing.T) {
	input := &MockInputSource{
		MeleePressedFunc: func() bool { return true },
	}
	g := newGroundedState(input)
	g.OnStart(0)

	next := g.Update()

	if next != kitstates.StateMeleeAttack {
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

	if next != kitstates.StateDashing {
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

	_, ok := kitstates.TryMeleeFromFalling(w, true /*meleePressed*/)
	if ok {
		t.Errorf("TryMeleeFromFalling returned ok=true while weapon on cooldown; want false")
	}
}

// AC3 bullet 3 — Pressing Dash or Jump while the combo window is open resets
// the combo. Verified via the shared helper used by ClimberPlayer.Update.
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
			if !w.AdvanceCombo() {
				t.Fatalf("precondition: AdvanceCombo failed; window not open?")
			}
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

			kitstates.ResetComboOnInterrupt(w, tc.dashPressed, tc.jumpPressed)

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
