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
// §4 RED-3 tests
// ---------------------------------------------------------------------------

func newTestMeleeWeaponForState(owner interface{}) *weapon.MeleeWeapon {
	w := weapon.NewMeleeWeapon(
		"player_melee",
		1,            // damage
		20,           // cooldownFrames
		[2]int{3, 5}, // activeFrames
		24*16, 16*16, // hitbox W/H fp16
		12*16, 0, // offset fp16
	)
	w.SetOwner(owner)
	return w
}

func TestMeleeAttackState_ReturnsToGrounded_WhenAnimationFinishes(t *testing.T) {
	owner := newPlayerOwner(100, 100, animation.FaceDirectionRight)
	_ = owner
	sp := &mockSpace{}
	w := newTestMeleeWeaponForState(owner)

	const animFrames = 8
	st := gamestates.NewMeleeAttackState(owner, sp, w, gamestates.StateGrounded)
	st.SetAnimationFrames(animFrames)
	st.OnStart(0)

	// First animFrames-1 ticks should stay in StateMeleeAttack.
	for i := 0; i < animFrames-1; i++ {
		if got := st.Update(); got != gamestates.StateMeleeAttack {
			t.Errorf("tick %d Update() = %v, want StateMeleeAttack", i, got)
		}
	}
	// The Nth tick should transition back to returnTo (StateGrounded).
	if got := st.Update(); got != gamestates.StateGrounded {
		t.Errorf("final tick Update() = %v, want StateGrounded", got)
	}
}

func TestMeleeAttackState_AirMelee_ReturnsToFalling(t *testing.T) {
	owner := newPlayerOwner(100, 100, animation.FaceDirectionRight)
	_ = owner
	owner.grounded = false
	sp := &mockSpace{}
	w := newTestMeleeWeaponForState(owner)

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

func TestMeleeAttackState_Update_AppliesHitboxDuringActiveWindow(t *testing.T) {
	owner := newPlayerOwner(100, 100, animation.FaceDirectionRight)
	_ = owner
	// Enemy positioned inside the forward hitbox range
	// (owner origin ~100 + offset 12, size 24x16 → enemy at 110).
	enemy := newMeleeEnemy("enemy", image.Rect(110, 100, 118, 108), combat.FactionEnemy)

	sp := &mockSpace{queryResult: []body.Collidable{enemy}}
	w := newTestMeleeWeaponForState(owner)

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

	// Force weapon into cooldown.
	w.SetCooldown(15)
	if w.CanFire() {
		t.Fatalf("precondition: weapon must not be able to fire during cooldown")
	}

	// TryMeleeFromFalling is the shared trigger helper used by the Falling
	// wiring (per SPEC §1.4). While on cooldown it must not transition.
	_, ok := gamestates.TryMeleeFromFalling(w, true /*meleePressed*/)
	if ok {
		t.Errorf("TryMeleeFromFalling returned ok=true while weapon on cooldown; want false")
	}
}
