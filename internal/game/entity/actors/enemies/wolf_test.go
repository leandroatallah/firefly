package gameenemies

import (
	"image"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
	kitcombat "github.com/boilerplate/ebiten-template/internal/kit/combat"

	"github.com/hajimehoshi/ebiten/v2"
)

// fakePlayerTarget is a minimal body.MovableCollidable used as a target for
// the WolfEnemy on_sight shooter. Only the position / shape / facing API is
// exercised by the shooter's range-gate + face-direction logic.
type fakePlayerTarget struct {
	x16, y16 int
	shape    body.Shape
	face     animation.FacingDirectionEnum
}

func newFakePlayerTarget(xPx, yPx int) *fakePlayerTarget {
	return &fakePlayerTarget{
		x16:   xPx * 16,
		y16:   yPx * 16,
		shape: bodyphysics.NewRect(0, 0, 16, 16),
	}
}

func (p *fakePlayerTarget) ID() string   { return "fake-player" }
func (p *fakePlayerTarget) SetID(string) {}
func (p *fakePlayerTarget) Position() image.Rectangle {
	return image.Rect(p.x16/16, p.y16/16, p.x16/16+16, p.y16/16+16)
}
func (p *fakePlayerTarget) SetPosition(x, y int)       { p.x16, p.y16 = x*16, y*16 }
func (p *fakePlayerTarget) SetPosition16(x16, y16 int) { p.x16, p.y16 = x16, y16 }
func (p *fakePlayerTarget) SetSize(int, int)           {}
func (p *fakePlayerTarget) Scale() float64             { return 1 }
func (p *fakePlayerTarget) SetScale(float64)           {}
func (p *fakePlayerTarget) GetPosition16() (int, int)  { return p.x16, p.y16 }
func (p *fakePlayerTarget) GetPositionMin() (int, int) { return p.x16 / 16, p.y16 / 16 }
func (p *fakePlayerTarget) GetShape() body.Shape       { return p.shape }
func (p *fakePlayerTarget) Owner() interface{}         { return nil }
func (p *fakePlayerTarget) SetOwner(interface{})       {}
func (p *fakePlayerTarget) LastOwner() interface{}     { return nil }
func (p *fakePlayerTarget) MoveX(int)                  {}
func (p *fakePlayerTarget) MoveY(int)                  {}
func (p *fakePlayerTarget) OnMoveLeft(int)             {}
func (p *fakePlayerTarget) OnMoveUpLeft(int)           {}
func (p *fakePlayerTarget) OnMoveDownLeft(int)         {}
func (p *fakePlayerTarget) OnMoveRight(int)            {}
func (p *fakePlayerTarget) OnMoveUpRight(int)          {}
func (p *fakePlayerTarget) OnMoveDownRight(int)        {}
func (p *fakePlayerTarget) OnMoveUp(int)               {}
func (p *fakePlayerTarget) OnMoveDown(int)             {}
func (p *fakePlayerTarget) Velocity() (int, int)       { return 0, 0 }
func (p *fakePlayerTarget) SetVelocity(int, int)       {}
func (p *fakePlayerTarget) Acceleration() (int, int)   { return 0, 0 }
func (p *fakePlayerTarget) SetAcceleration(int, int)   {}
func (p *fakePlayerTarget) SetSpeed(int) error         { return nil }
func (p *fakePlayerTarget) SetMaxSpeed(int) error      { return nil }
func (p *fakePlayerTarget) Speed() int                 { return 0 }
func (p *fakePlayerTarget) MaxSpeed() int              { return 0 }
func (p *fakePlayerTarget) Immobile() bool             { return false }
func (p *fakePlayerTarget) SetImmobile(bool)           {}
func (p *fakePlayerTarget) SetFreeze(bool)             {}
func (p *fakePlayerTarget) Freeze() bool               { return false }
func (p *fakePlayerTarget) FaceDirection() animation.FacingDirectionEnum {
	return p.face
}
func (p *fakePlayerTarget) SetFaceDirection(v animation.FacingDirectionEnum) { p.face = v }
func (p *fakePlayerTarget) IsIdle() bool                                     { return true }
func (p *fakePlayerTarget) IsWalking() bool                                  { return false }
func (p *fakePlayerTarget) IsFalling() bool                                  { return false }
func (p *fakePlayerTarget) IsGoingUp() bool                                  { return false }
func (p *fakePlayerTarget) CheckMovementDirectionX()                         {}
func (p *fakePlayerTarget) TryJump(int)                                      {}
func (p *fakePlayerTarget) SetJumpForceMultiplier(float64)                   {}
func (p *fakePlayerTarget) JumpForceMultiplier() float64                     { return 1 }
func (p *fakePlayerTarget) SetHorizontalInertia(float64)                     {}
func (p *fakePlayerTarget) HorizontalInertia() float64                       { return 1 }
func (p *fakePlayerTarget) GetTouchable() body.Touchable                     { return nil }
func (p *fakePlayerTarget) DrawCollisionBox(screen *ebiten.Image, pos image.Rectangle) {
}
func (p *fakePlayerTarget) CollisionPosition() []image.Rectangle {
	return []image.Rectangle{p.Position()}
}
func (p *fakePlayerTarget) CollisionShapes() []body.Collidable { return nil }
func (p *fakePlayerTarget) IsObstructive() bool                { return false }
func (p *fakePlayerTarget) SetIsObstructive(bool)              {}
func (p *fakePlayerTarget) AddCollision(...body.Collidable)    {}
func (p *fakePlayerTarget) ClearCollisions()                   {}
func (p *fakePlayerTarget) SetTouchable(body.Touchable)        {}
func (p *fakePlayerTarget) OnTouch(body.Collidable)            {}
func (p *fakePlayerTarget) OnBlock(body.Collidable)            {}
func (p *fakePlayerTarget) ApplyValidPosition(int, bool, body.BodiesSpace) (int, int, bool) {
	return p.x16 / 16, p.y16 / 16, false
}

// TestWolfEnemy_FactionAndShooter verifies WolfEnemy construction wires
// FactionEnemy and a non-nil shooter from wolf.json's weapon block
// (shoot_mode=on_sight, shoot_direction=horizontal).
func TestWolfEnemy_FactionAndShooter(t *testing.T) {
	ctx := newEnemyTestContext()

	wolf, err := NewWolfEnemy(ctx, 100, 100, "wolf-1")
	if err != nil {
		t.Fatalf("NewWolfEnemy returned error: %v", err)
	}
	if wolf == nil {
		t.Fatal("NewWolfEnemy returned nil")
	}

	if got := wolf.GetCharacter().Faction(); got != kitcombat.FactionEnemy {
		t.Errorf("Faction() = %v, want FactionEnemy", got)
	}
	if wolf.Shooter() == nil {
		t.Error("expected WolfEnemy.Shooter() to be non-nil after construction")
	}
}

// TestWolfEnemy_OnSightFires_WhenTargetInRange verifies the Smart Patrol
// Shooter archetype: when SetTarget is called with a player within range,
// repeated Update calls spawn at least one projectile.
func TestWolfEnemy_OnSightFires_WhenTargetInRange(t *testing.T) {
	ctx := newEnemyTestContext()

	wolf, err := NewWolfEnemy(ctx, 100, 100, "wolf-fire")
	if err != nil {
		t.Fatalf("NewWolfEnemy returned error: %v", err)
	}

	// wolf.json range=160; place target 40px away → inside range.
	target := newFakePlayerTarget(140, 100)
	wolf.SetTarget(target)

	// wolf.json cooldown=90 frames; 240 frames guarantees at least one fire.
	firstBulletSeen := false
	for i := 0; i < 240; i++ {
		if err := wolf.Update(ctx.Space); err != nil {
			t.Fatalf("Update error at frame %d: %v", i, err)
		}
		if !firstBulletSeen {
			if bullets := bulletBodies(ctx.Space); len(bullets) > 0 {
				firstBulletSeen = true
			}
		}
	}

	if !firstBulletSeen {
		t.Fatal("expected WolfEnemy to fire at least one projectile when target is in range, got 0")
	}
}

// TestWolfEnemy_OnSightSkips_WhenTargetOutOfRange verifies that no projectile
// spawns when the target is beyond the configured range.
func TestWolfEnemy_OnSightSkips_WhenTargetOutOfRange(t *testing.T) {
	ctx := newEnemyTestContext()

	wolf, err := NewWolfEnemy(ctx, 100, 100, "wolf-far")
	if err != nil {
		t.Fatalf("NewWolfEnemy returned error: %v", err)
	}

	// wolf.json range=160; place target 500px away → out of range.
	target := newFakePlayerTarget(600, 100)
	wolf.SetTarget(target)

	for i := 0; i < 240; i++ {
		if err := wolf.Update(ctx.Space); err != nil {
			t.Fatalf("Update error at frame %d: %v", i, err)
		}
	}

	if bullets := bulletBodies(ctx.Space); len(bullets) != 0 {
		t.Errorf("expected 0 projectiles when target is out of range, got %d", len(bullets))
	}
}

// TestWolfEnemy_OnSightSkips_WhenNoTarget verifies no projectile spawns when
// SetTarget has never been called (on_sight gate without a target).
func TestWolfEnemy_OnSightSkips_WhenNoTarget(t *testing.T) {
	ctx := newEnemyTestContext()

	wolf, err := NewWolfEnemy(ctx, 100, 100, "wolf-notarget")
	if err != nil {
		t.Fatalf("NewWolfEnemy returned error: %v", err)
	}

	for i := 0; i < 240; i++ {
		if err := wolf.Update(ctx.Space); err != nil {
			t.Fatalf("Update error at frame %d: %v", i, err)
		}
	}

	if bullets := bulletBodies(ctx.Space); len(bullets) != 0 {
		t.Errorf("expected 0 projectiles without a target, got %d", len(bullets))
	}
}
