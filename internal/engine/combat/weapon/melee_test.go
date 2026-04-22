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
// §4 RED-1 tests
// ---------------------------------------------------------------------------

// newTestMeleeWeapon constructs a MeleeWeapon directly (bypasses the factory)
// so that hitbox frame tests can run without JSON plumbing.
// damage=1, activeFrames=[3,5], cooldown=20, hitbox 24x16 offset (12,0).
func newTestMeleeWeapon(owner interface{}) *weapon.MeleeWeapon {
	w := weapon.NewMeleeWeapon(
		"player_melee",
		1,            // damage
		20,           // cooldownFrames
		[2]int{3, 5}, // activeFrames
		24*16, 16*16, // hitbox W/H in fp16
		12*16, 0, // hitbox offset in fp16
	)
	w.SetOwner(owner)
	return w
}

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

	// Hitbox at frame 3 extends from owner origin + offset (12*16 fp16 = 12px) with
	// size 24x16. Place overlapping ally/enemy at x≈110, and far enemy at x=400.
	enemy := newMeleeTarget("enemy", 110, 100, 8, 8, combat.FactionEnemy)
	ally := newMeleeTarget("ally", 110, 100, 8, 8, combat.FactionPlayer)
	farEnemy := newMeleeTarget("far_enemy", 400, 100, 8, 8, combat.FactionEnemy)

	space := &fakeSpace{}
	space.AddBody(enemy)
	space.AddBody(ally)
	space.AddBody(farEnemy)

	w.Fire(owner.x16, owner.y16, owner.face, body.ShootDirectionStraight, 0)
	// Advance to first active frame (swingFrame == 3).
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
	// Advance to frame 3 and apply hitbox across the full active window (frames 3..5).
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

	// Target on the LEFT of the owner. If offset is correctly mirrored, the
	// query rect will include (88, 100); if NOT mirrored, a same-position target
	// on the RIGHT (x=110) would be hit but the left target would not.
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

	// Cross-check: the recorded query rect must be to the LEFT of the owner
	// origin (centerX ≈ 100), not the right.
	ownerCenterX := 100
	if space.lastQuery.Max.X > ownerCenterX+4 {
		t.Errorf("query rect %+v extends to the right of owner origin (x=%d); expected mirrored to the left", space.lastQuery, ownerCenterX)
	}
}
