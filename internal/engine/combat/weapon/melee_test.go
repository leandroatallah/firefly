package weapon_test

import (
	"image"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/combat/weapon"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/tilemaplayer"
	"github.com/hajimehoshi/ebiten/v2"
)

// ---------------------------------------------------------------------------
// Test doubles (package-local).
// ---------------------------------------------------------------------------

// meleeTarget is a Damageable + Factioned + Collidable target.
type meleeTarget struct {
	id          string
	x16, y16    int
	w, h        int
	faction     combat.Faction
	damageCalls []int
}

func newMeleeTarget(id string, xPx, yPx, w, h int, f combat.Faction) *meleeTarget {
	return &meleeTarget{
		id: id, x16: xPx * 16, y16: yPx * 16, w: w, h: h, faction: f,
	}
}

func (t *meleeTarget) TakeDamage(amount int)   { t.damageCalls = append(t.damageCalls, amount) }
func (t *meleeTarget) Faction() combat.Faction { return t.faction }

// body.Collidable surface (minimal, only what ApplyHitbox needs via Query results).
func (t *meleeTarget) ID() string   { return t.id }
func (t *meleeTarget) SetID(string) {}
func (t *meleeTarget) Position() image.Rectangle {
	return image.Rect(t.x16/16, t.y16/16, t.x16/16+t.w, t.y16/16+t.h)
}
func (t *meleeTarget) SetPosition(int, int)                            {}
func (t *meleeTarget) SetPosition16(x16, y16 int)                      { t.x16, t.y16 = x16, y16 }
func (t *meleeTarget) SetSize(int, int)                                {}
func (t *meleeTarget) Scale() float64                                  { return 1 }
func (t *meleeTarget) SetScale(float64)                                {}
func (t *meleeTarget) GetPosition16() (int, int)                       { return t.x16, t.y16 }
func (t *meleeTarget) GetPositionMin() (int, int)                      { return t.x16 / 16, t.y16 / 16 }
func (t *meleeTarget) GetShape() body.Shape                            { return nil }
func (t *meleeTarget) Owner() interface{}                              { return nil }
func (t *meleeTarget) SetOwner(interface{})                            {}
func (t *meleeTarget) LastOwner() interface{}                          { return nil }
func (t *meleeTarget) GetTouchable() body.Touchable                    { return nil }
func (t *meleeTarget) DrawCollisionBox(*ebiten.Image, image.Rectangle) {}
func (t *meleeTarget) CollisionPosition() []image.Rectangle            { return []image.Rectangle{t.Position()} }
func (t *meleeTarget) CollisionShapes() []body.Collidable              { return nil }
func (t *meleeTarget) IsObstructive() bool                             { return false }
func (t *meleeTarget) SetIsObstructive(bool)                           {}
func (t *meleeTarget) AddCollision(...body.Collidable)                 {}
func (t *meleeTarget) ClearCollisions()                                {}
func (t *meleeTarget) SetTouchable(body.Touchable)                     {}
func (t *meleeTarget) OnTouch(body.Collidable)                         {}
func (t *meleeTarget) OnBlock(body.Collidable)                         {}
func (t *meleeTarget) ApplyValidPosition(int, bool, body.BodiesSpace) (int, int, bool) {
	return t.x16 / 16, t.y16 / 16, false
}

// meleeOwner is a Factioned owner with position/face direction.
type meleeOwner struct {
	meleeTarget
	face animation.FacingDirectionEnum
}

func newMeleeOwner(xPx, yPx int, f combat.Faction, face animation.FacingDirectionEnum) *meleeOwner { //nolint:unparam
	return &meleeOwner{
		meleeTarget: *newMeleeTarget("owner", xPx, yPx, 16, 16, f),
		face:        face,
	}
}

func (o *meleeOwner) FaceDirection() animation.FacingDirectionEnum { return o.face }

// fakeSpace is a minimal BodiesSpace that records Query calls and returns
// bodies whose Position() overlaps the query rect.
type fakeSpace struct {
	bodies     []body.Collidable
	lastQuery  image.Rectangle
	queryCalls int
}

func (s *fakeSpace) AddBody(b body.Collidable)                                           { s.bodies = append(s.bodies, b) }
func (s *fakeSpace) Bodies() []body.Collidable                                           { return s.bodies }
func (s *fakeSpace) RemoveBody(body.Collidable)                                          {}
func (s *fakeSpace) QueueForRemoval(body.Collidable)                                     {}
func (s *fakeSpace) ProcessRemovals()                                                    {}
func (s *fakeSpace) Clear()                                                              { s.bodies = nil }
func (s *fakeSpace) ResolveCollisions(body.Collidable) (bool, bool)                      { return false, false }
func (s *fakeSpace) SetTilemapDimensionsProvider(tilemaplayer.TilemapDimensionsProvider) {}
func (s *fakeSpace) GetTilemapDimensionsProvider() tilemaplayer.TilemapDimensionsProvider {
	return nil
}
func (s *fakeSpace) Find(string) body.Collidable { return nil }
func (s *fakeSpace) Query(rect image.Rectangle) []body.Collidable {
	s.queryCalls++
	s.lastQuery = rect
	var hits []body.Collidable
	for _, b := range s.bodies {
		if b.Position().Overlaps(rect) {
			hits = append(hits, b)
		}
	}
	return hits
}

// ---------------------------------------------------------------------------
// Helpers for US-041 combo-aware weapon construction.
// ---------------------------------------------------------------------------

// newTestMeleeWeapon constructs a single-step MeleeWeapon (US-040 parity).
// damage=1, activeFrames=[3,5], cooldown=20, hitbox 24x16 offset (12,0).
func newTestMeleeWeapon(owner interface{}) *weapon.MeleeWeapon {
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

// newThreeStepComboWeapon builds a 3-step weapon used by the combo tests.
// Cooldown is 0 so step-to-step transitions aren't gated by cooldown.
// Damage progression: 1, 1, 2 matches the AC6 example.
// Hitbox width also varies per step so the per-step-hitbox test has a
// distinguishing property to observe.
// All steps share ActiveFrames [3,5] so runSwingToCompletion can use a fixed
// advance count (lastTestStepActiveFrame = 5).
func newThreeStepComboWeapon(owner interface{}) *weapon.MeleeWeapon {
	steps := []weapon.ComboStep{
		{Damage: 1, ActiveFrames: [2]int{3, 5}, HitboxW16: 24 * 16, HitboxH16: 16 * 16, HitboxOffsetX16: 12 * 16, HitboxOffsetY16: 0},
		{Damage: 1, ActiveFrames: [2]int{3, 5}, HitboxW16: 28 * 16, HitboxH16: 16 * 16, HitboxOffsetX16: 14 * 16, HitboxOffsetY16: -4 * 16},
		{Damage: 2, ActiveFrames: [2]int{3, 5}, HitboxW16: 32 * 16, HitboxH16: 20 * 16, HitboxOffsetX16: 16 * 16, HitboxOffsetY16: 0},
	}
	w := weapon.NewMeleeWeapon("player_melee", 0 /*cooldown*/, 15 /*comboWindowFrames*/, steps)
	w.SetOwner(owner)
	return w
}

// lastTestStepActiveFrame is the shared ActiveFrames[1] for all steps in
// newThreeStepComboWeapon. runSwingToCompletion advances past this frame.
const lastTestStepActiveFrame = 5

// runSwingToCompletion fires and advances the weapon past the current step's
// active window (ActiveFrames[1] == lastTestStepActiveFrame for all test steps).
func runSwingToCompletion(w *weapon.MeleeWeapon, owner *meleeOwner) {
	w.Fire(owner.x16, owner.y16, owner.face, body.ShootDirectionStraight, 0)
	// Advance swingFrame from 0 to lastTestStepActiveFrame+1 (inclusive) so
	// swinging goes false and Update() opens the combo window (or resets after
	// the last step).
	for i := 0; i <= lastTestStepActiveFrame+1; i++ {
		w.Update()
	}
}

// ---------------------------------------------------------------------------
// Single-step (US-040) parity tests — now driven by ComboStep[0].
// ---------------------------------------------------------------------------

func TestMeleeWeapon_Fire_HitboxActivation(t *testing.T) {
	tests := []struct {
		frame      int
		wantActive bool
	}{
		{0, false},
		{2, false},
		{3, true},
		{4, true},
		{5, true},
		{6, false},
	}

	for _, tc := range tests {
		name := ""
		switch tc.wantActive {
		case true:
			name = "active window frame"
		default:
			name = "inactive frame"
		}
		t.Run(name, func(t *testing.T) {
			owner := newMeleeOwner(100, 100, combat.FactionPlayer, animation.FaceDirectionRight)
			w := newTestMeleeWeapon(owner)

			w.Fire(owner.x16, owner.y16, owner.face, body.ShootDirectionStraight, 0)
			for i := 0; i < tc.frame; i++ {
				w.Update()
			}

			if got := w.IsHitboxActive(); got != tc.wantActive {
				t.Errorf("frame %d IsHitboxActive() = %v, want %v", tc.frame, got, tc.wantActive)
			}
		})
	}
}

func TestMeleeWeapon_ApplyHitbox_FactionGating(t *testing.T) {
	owner := newMeleeOwner(100, 100, combat.FactionPlayer, animation.FaceDirectionRight)
	w := newTestMeleeWeapon(owner)

	enemy := newMeleeTarget("enemy", 110, 100, 8, 8, combat.FactionEnemy)
	ally := newMeleeTarget("ally", 110, 100, 8, 8, combat.FactionPlayer)
	farEnemy := newMeleeTarget("far_enemy", 400, 100, 8, 8, combat.FactionEnemy)

	space := &fakeSpace{}
	space.AddBody(enemy)
	space.AddBody(ally)
	space.AddBody(farEnemy)

	w.Fire(owner.x16, owner.y16, owner.face, body.ShootDirectionStraight, 0)
	for i := 0; i < 3; i++ {
		w.Update()
	}
	if !w.IsHitboxActive() {
		t.Fatalf("expected hitbox to be active at frame 3")
	}
	w.ApplyHitbox(space)

	if len(enemy.damageCalls) != 1 || enemy.damageCalls[0] != 1 {
		t.Errorf("enemy TakeDamage: got %v, want [1]", enemy.damageCalls)
	}
	if len(ally.damageCalls) != 0 {
		t.Errorf("ally TakeDamage: got %v, want no calls (same-faction gate)", ally.damageCalls)
	}
	if len(farEnemy.damageCalls) != 0 {
		t.Errorf("far enemy TakeDamage: got %v, want no calls (outside hitbox)", farEnemy.damageCalls)
	}
}

func TestMeleeWeapon_ApplyHitbox_SingleHitPerSwing(t *testing.T) {
	owner := newMeleeOwner(100, 100, combat.FactionPlayer, animation.FaceDirectionRight)
	w := newTestMeleeWeapon(owner)

	enemy := newMeleeTarget("enemy", 110, 100, 8, 8, combat.FactionEnemy)
	space := &fakeSpace{}
	space.AddBody(enemy)

	w.Fire(owner.x16, owner.y16, owner.face, body.ShootDirectionStraight, 0)
	for i := 0; i < 3; i++ {
		w.Update()
	}
	for f := 3; f <= 5; f++ {
		if !w.IsHitboxActive() {
			t.Fatalf("expected hitbox active at frame %d", f)
		}
		w.ApplyHitbox(space)
		w.Update()
	}

	if len(enemy.damageCalls) != 1 {
		t.Errorf("TakeDamage called %d times, want exactly 1 (single-hit per swing)", len(enemy.damageCalls))
	}
}

func TestMeleeWeapon_Cooldown_PreventsRefire(t *testing.T) {
	const cooldownFrames = 20
	owner := newMeleeOwner(100, 100, combat.FactionPlayer, animation.FaceDirectionRight)
	w := newTestMeleeWeapon(owner)

	if !w.CanFire() {
		t.Fatalf("weapon should be ready to fire before first Fire()")
	}

	w.Fire(owner.x16, owner.y16, owner.face, body.ShootDirectionStraight, 0)
	if w.CanFire() {
		t.Errorf("immediately after Fire, CanFire() = true, want false")
	}

	for i := 0; i < cooldownFrames-1; i++ {
		w.Update()
	}
	if w.CanFire() {
		t.Errorf("after %d of %d cooldown frames, CanFire() = true, want false", cooldownFrames-1, cooldownFrames)
	}

	w.Update() // cooldown complete
	if !w.CanFire() {
		t.Errorf("after full cooldown (%d frames), CanFire() = false, want true", cooldownFrames)
	}
}

func TestMeleeWeapon_Fire_MirrorsHitboxWhenFacingLeft(t *testing.T) {
	owner := newMeleeOwner(100, 100, combat.FactionPlayer, animation.FaceDirectionLeft)
	w := newTestMeleeWeapon(owner)

	leftTarget := newMeleeTarget("left", 88, 100, 8, 8, combat.FactionEnemy)
	rightTarget := newMeleeTarget("right", 118, 100, 8, 8, combat.FactionEnemy)

	space := &fakeSpace{}
	space.AddBody(leftTarget)
	space.AddBody(rightTarget)

	w.Fire(owner.x16, owner.y16, owner.face, body.ShootDirectionStraight, 0)
	for i := 0; i < 3; i++ {
		w.Update()
	}
	w.ApplyHitbox(space)

	if len(leftTarget.damageCalls) != 1 {
		t.Errorf("left target (facing-left owner) TakeDamage: got %v, want [1] (hitbox should mirror)", leftTarget.damageCalls)
	}
	if len(rightTarget.damageCalls) != 0 {
		t.Errorf("right target TakeDamage: got %v, want no calls", rightTarget.damageCalls)
	}

	ownerCenterX := 100
	if space.lastQuery.Max.X > ownerCenterX+4 {
		t.Errorf("query rect %+v extends to the right of owner origin (x=%d); expected mirrored to the left", space.lastQuery, ownerCenterX)
	}
}

// ---------------------------------------------------------------------------
// §4 RED-1 — Combo chain behaviour
// ---------------------------------------------------------------------------

// AC1 — Pressing Z within combo_window_frames after a hit advances to the
// next step; state machine advances step 0 → 1 → 2, and the last step wraps
// back to 0 (AC4).
func TestMeleeWeapon_Combo_AdvancesWhenPressedWithinWindow(t *testing.T) {
	owner := newMeleeOwner(100, 100, combat.FactionPlayer, animation.FaceDirectionRight)
	w := newThreeStepComboWeapon(owner)

	if w.StepIndex() != 0 {
		t.Fatalf("initial StepIndex() = %d, want 0", w.StepIndex())
	}

	// Step 1 swing.
	runSwingToCompletion(w, owner)
	if w.ComboWindowRemaining() <= 0 {
		t.Fatalf("after step 1 completes, ComboWindowRemaining() = %d, want > 0", w.ComboWindowRemaining())
	}
	if w.StepIndex() != 0 {
		t.Errorf("before AdvanceCombo, StepIndex() = %d, want 0 (advance happens on the press, not on swing end)", w.StepIndex())
	}

	// Press Z within the window → advance to step 2.
	if !w.AdvanceCombo() {
		t.Fatalf("AdvanceCombo() returned false within window; want true")
	}
	if w.StepIndex() != 1 {
		t.Errorf("after AdvanceCombo, StepIndex() = %d, want 1", w.StepIndex())
	}

	// Step 2 swing.
	runSwingToCompletion(w, owner)
	if !w.AdvanceCombo() {
		t.Fatalf("AdvanceCombo() to step 3 returned false; want true")
	}
	if w.StepIndex() != 2 {
		t.Errorf("after second AdvanceCombo, StepIndex() = %d, want 2", w.StepIndex())
	}

	// Step 3 swing → AC4: chain resets automatically, no fourth swing possible.
	runSwingToCompletion(w, owner)
	if w.StepIndex() != 0 {
		t.Errorf("after last step swing completes, StepIndex() = %d, want 0 (AC4 last-step wrap)", w.StepIndex())
	}
	if w.ComboWindowRemaining() != 0 {
		t.Errorf("after last step, ComboWindowRemaining() = %d, want 0 (no window after wrap)", w.ComboWindowRemaining())
	}
}

// AC3 bullet 1 — If the player does not press Z within combo_window_frames,
// the chain resets to step 1 (index 0).
func TestMeleeWeapon_Combo_ResetsOnWindowExpiry(t *testing.T) {
	owner := newMeleeOwner(100, 100, combat.FactionPlayer, animation.FaceDirectionRight)
	w := newThreeStepComboWeapon(owner)

	runSwingToCompletion(w, owner)
	if w.ComboWindowRemaining() <= 0 {
		t.Fatalf("precondition: combo window must be open; got %d", w.ComboWindowRemaining())
	}
	// Advance past the window without any AdvanceCombo call.
	for i := 0; i < 15; /*comboWindowFrames*/ i++ {
		w.Update()
	}

	if w.StepIndex() != 0 {
		t.Errorf("after window expiry, StepIndex() = %d, want 0", w.StepIndex())
	}
	if w.ComboWindowRemaining() != 0 {
		t.Errorf("after window expiry, ComboWindowRemaining() = %d, want 0", w.ComboWindowRemaining())
	}
	if w.AdvanceCombo() {
		t.Errorf("AdvanceCombo() returned true after window expired; want false")
	}
}

// AC3 — Explicit reset API (used by ClimberPlayer.Hurt and by dash/jump interrupts).
func TestMeleeWeapon_Combo_ResetsOnDemand(t *testing.T) {
	owner := newMeleeOwner(100, 100, combat.FactionPlayer, animation.FaceDirectionRight)
	w := newThreeStepComboWeapon(owner)

	runSwingToCompletion(w, owner)
	if !w.AdvanceCombo() {
		t.Fatalf("precondition: AdvanceCombo should succeed while window open")
	}
	if w.StepIndex() != 1 {
		t.Fatalf("precondition: StepIndex should be 1 before reset, got %d", w.StepIndex())
	}

	w.ResetCombo()

	if w.StepIndex() != 0 {
		t.Errorf("after ResetCombo, StepIndex() = %d, want 0", w.StepIndex())
	}
	if w.ComboWindowRemaining() != 0 {
		t.Errorf("after ResetCombo, ComboWindowRemaining() = %d, want 0", w.ComboWindowRemaining())
	}
}

// AC2 — Each combo step uses its own hitbox dimensions and damage value.
func TestMeleeWeapon_Combo_PerStepDamageAndHitbox(t *testing.T) {
	owner := newMeleeOwner(100, 100, combat.FactionPlayer, animation.FaceDirectionRight)
	w := newThreeStepComboWeapon(owner)

	// Enemy placed to overlap all three step hitboxes (the closest one —
	// right in front of the owner — is reachable from each step's forward
	// hitbox). A fresh enemy per step keeps single-hit-per-swing honest.
	type stepExpect struct {
		damage int
	}
	expects := []stepExpect{{1}, {1}, {2}}

	for stepIdx, exp := range expects {
		if w.StepIndex() != stepIdx {
			t.Fatalf("step %d: precondition StepIndex() = %d, want %d", stepIdx, w.StepIndex(), stepIdx)
		}

		enemy := newMeleeTarget("enemy", 110, 100, 8, 8, combat.FactionEnemy)
		space := &fakeSpace{}
		space.AddBody(enemy)

		w.Fire(owner.x16, owner.y16, owner.face, body.ShootDirectionStraight, 0)
		// Advance to the first active frame of this step (ActiveFrames[0]=3 for all steps).
		for i := 0; i < 3; i++ {
			w.Update()
		}
		if !w.IsHitboxActive() {
			t.Fatalf("step %d: expected hitbox active at frame 3", stepIdx)
		}
		w.ApplyHitbox(space)

		if len(enemy.damageCalls) != 1 {
			t.Errorf("step %d: TakeDamage called %d times, want 1", stepIdx, len(enemy.damageCalls))
			continue
		}
		if enemy.damageCalls[0] != exp.damage {
			t.Errorf("step %d: TakeDamage value = %d, want %d", stepIdx, enemy.damageCalls[0], exp.damage)
		}

		// Hitbox width grows per step (24 → 28 → 32). Verify the query rect width
		// reflects the active step's dimensions.
		wantWidth := []int{24, 28, 32}[stepIdx]
		if space.lastQuery.Dx() != wantWidth {
			t.Errorf("step %d: query rect width = %d, want %d (step-specific hitbox)", stepIdx, space.lastQuery.Dx(), wantWidth)
		}

		// Finish the swing and, for steps 0/1, advance to the next step.
		for i := 3; i <= lastTestStepActiveFrame+1; i++ {
			w.Update()
		}
		if stepIdx < 2 {
			if !w.AdvanceCombo() {
				t.Fatalf("step %d: AdvanceCombo() returned false; want true", stepIdx)
			}
		}
	}
}

// AC4 — After the 3rd hit, the chain always resets without any external call.
func TestMeleeWeapon_Combo_LastStepAlwaysResets(t *testing.T) {
	owner := newMeleeOwner(100, 100, combat.FactionPlayer, animation.FaceDirectionRight)
	w := newThreeStepComboWeapon(owner)

	// Step 1
	runSwingToCompletion(w, owner)
	if !w.AdvanceCombo() {
		t.Fatalf("step 1→2 AdvanceCombo failed")
	}
	// Step 2
	runSwingToCompletion(w, owner)
	if !w.AdvanceCombo() {
		t.Fatalf("step 2→3 AdvanceCombo failed")
	}
	// Step 3
	runSwingToCompletion(w, owner)

	if w.StepIndex() != 0 {
		t.Errorf("after final step swing, StepIndex() = %d, want 0 (auto-reset)", w.StepIndex())
	}
	if w.ComboWindowRemaining() != 0 {
		t.Errorf("after final step swing, ComboWindowRemaining() = %d, want 0", w.ComboWindowRemaining())
	}
	// AdvanceCombo from a fresh (reset) state with no open window must fail.
	if w.AdvanceCombo() {
		t.Errorf("AdvanceCombo after auto-reset returned true; want false (no open window)")
	}
}

// AdvanceCombo is a no-op when no window is open (press outside the window
// should start step 0 on the next Fire, not jump ahead).
func TestMeleeWeapon_Combo_AdvanceCombo_NoopWhenWindowClosed(t *testing.T) {
	owner := newMeleeOwner(100, 100, combat.FactionPlayer, animation.FaceDirectionRight)
	w := newThreeStepComboWeapon(owner)

	if w.AdvanceCombo() {
		t.Errorf("AdvanceCombo on fresh weapon returned true; want false (no window open)")
	}
	if w.StepIndex() != 0 {
		t.Errorf("after no-op AdvanceCombo, StepIndex() = %d, want 0", w.StepIndex())
	}
}

// Steps() getter exposes the configured combo steps for tests / introspection.
func TestMeleeWeapon_Steps_ReturnsConfiguredSlice(t *testing.T) {
	owner := newMeleeOwner(100, 100, combat.FactionPlayer, animation.FaceDirectionRight)
	w := newThreeStepComboWeapon(owner)

	got := w.Steps()
	if len(got) != 3 {
		t.Fatalf("Steps() len = %d, want 3", len(got))
	}
	if got[0].Damage != 1 || got[2].Damage != 2 {
		t.Errorf("Steps() damage = [%d,%d,%d], want [1,1,2]", got[0].Damage, got[1].Damage, got[2].Damage)
	}
	if got[1].HitboxOffsetY16 != -4*16 {
		t.Errorf("Steps()[1].HitboxOffsetY16 = %d, want %d", got[1].HitboxOffsetY16, -4*16)
	}
}
