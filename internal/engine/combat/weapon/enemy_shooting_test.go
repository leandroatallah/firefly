package weapon_test

import (
	"image"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/combat/weapon"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
)

// fakeShape implements body.Shape.
type fakeShape struct {
	w, h int
}

func (s *fakeShape) Width() int  { return s.w }
func (s *fakeShape) Height() int { return s.h }

// fakeOwner implements the subset of body.MovableCollidable that EnemyShooting needs:
// GetPosition16, FaceDirection, SetFaceDirection, GetShape, State, plus Faction().
// The remaining MovableCollidable methods are stubbed so the type satisfies the
// interface the constructor signature requires (body.MovableCollidable).
type fakeOwner struct {
	x16, y16       int
	shape          *fakeShape
	face           animation.FacingDirectionEnum
	state          actors.ActorStateEnum
	faction        combat.Faction
	setFaceDirCnt  int
	lastFaceDirSet animation.FacingDirectionEnum
}

func newFakeOwner(xPx, yPx int) *fakeOwner {
	return &fakeOwner{
		x16:   xPx * 16,
		y16:   yPx * 16,
		shape: &fakeShape{w: 16, h: 16},
		face:  animation.FaceDirectionRight,
	}
}

// ... Body ---
func (o *fakeOwner) ID() string   { return "fake-owner" }
func (o *fakeOwner) SetID(string) {}
func (o *fakeOwner) Position() image.Rectangle {
	return image.Rect(o.x16/16, o.y16/16, o.x16/16+o.shape.w, o.y16/16+o.shape.h)
}
func (o *fakeOwner) SetPosition(x, y int)       { o.x16, o.y16 = x*16, y*16 }
func (o *fakeOwner) SetPosition16(x16, y16 int) { o.x16, o.y16 = x16, y16 }
func (o *fakeOwner) SetSize(int, int)           {}
func (o *fakeOwner) Scale() float64             { return 1 }
func (o *fakeOwner) SetScale(float64)           {}
func (o *fakeOwner) GetPosition16() (int, int)  { return o.x16, o.y16 }
func (o *fakeOwner) GetPositionMin() (int, int) { return o.x16 / 16, o.y16 / 16 }
func (o *fakeOwner) GetShape() body.Shape       { return o.shape }

// ... Ownable ---
func (o *fakeOwner) Owner() interface{}     { return nil }
func (o *fakeOwner) SetOwner(interface{})   {}
func (o *fakeOwner) LastOwner() interface{} { return nil }

// ... Movable ---
func (o *fakeOwner) MoveX(int)                {}
func (o *fakeOwner) MoveY(int)                {}
func (o *fakeOwner) OnMoveLeft(int)           {}
func (o *fakeOwner) OnMoveUpLeft(int)         {}
func (o *fakeOwner) OnMoveDownLeft(int)       {}
func (o *fakeOwner) OnMoveRight(int)          {}
func (o *fakeOwner) OnMoveUpRight(int)        {}
func (o *fakeOwner) OnMoveDownRight(int)      {}
func (o *fakeOwner) OnMoveUp(int)             {}
func (o *fakeOwner) OnMoveDown(int)           {}
func (o *fakeOwner) Velocity() (int, int)     { return 0, 0 }
func (o *fakeOwner) SetVelocity(int, int)     {}
func (o *fakeOwner) Acceleration() (int, int) { return 0, 0 }
func (o *fakeOwner) SetAcceleration(int, int) {}
func (o *fakeOwner) SetSpeed(int) error       { return nil }
func (o *fakeOwner) SetMaxSpeed(int) error    { return nil }
func (o *fakeOwner) Speed() int               { return 0 }
func (o *fakeOwner) MaxSpeed() int            { return 0 }
func (o *fakeOwner) Immobile() bool           { return false }
func (o *fakeOwner) SetImmobile(bool)         {}
func (o *fakeOwner) SetFreeze(bool)           {}
func (o *fakeOwner) Freeze() bool             { return false }
func (o *fakeOwner) FaceDirection() animation.FacingDirectionEnum {
	return o.face
}
func (o *fakeOwner) SetFaceDirection(v animation.FacingDirectionEnum) {
	o.setFaceDirCnt++
	o.lastFaceDirSet = v
	o.face = v
}
func (o *fakeOwner) IsIdle() bool                   { return true }
func (o *fakeOwner) IsWalking() bool                { return false }
func (o *fakeOwner) IsFalling() bool                { return false }
func (o *fakeOwner) IsGoingUp() bool                { return false }
func (o *fakeOwner) CheckMovementDirectionX()       {}
func (o *fakeOwner) TryJump(int)                    {}
func (o *fakeOwner) SetJumpForceMultiplier(float64) {}
func (o *fakeOwner) JumpForceMultiplier() float64   { return 1 }
func (o *fakeOwner) SetHorizontalInertia(float64)   {}
func (o *fakeOwner) HorizontalInertia() float64     { return 1 }

// ... Collidable ---
func (o *fakeOwner) GetTouchable() body.Touchable                  { return nil }
func (o *fakeOwner) DrawCollisionBox(interface{}, image.Rectangle) {}
func (o *fakeOwner) CollisionPosition() []image.Rectangle {
	return []image.Rectangle{o.Position()}
}
func (o *fakeOwner) CollisionShapes() []body.Collidable { return nil }
func (o *fakeOwner) IsObstructive() bool                { return false }
func (o *fakeOwner) SetIsObstructive(bool)              {}
func (o *fakeOwner) AddCollision(...body.Collidable)    {}
func (o *fakeOwner) ClearCollisions()                   {}
func (o *fakeOwner) SetTouchable(body.Touchable)        {}
func (o *fakeOwner) OnTouch(body.Collidable)            {}
func (o *fakeOwner) OnBlock(body.Collidable)            {}
func (o *fakeOwner) ApplyValidPosition(int, bool, body.BodiesSpace) (int, int, bool) {
	return o.x16 / 16, o.y16 / 16, false
}

// State gate / faction helpers consumed by EnemyShooting.
func (o *fakeOwner) State() actors.ActorStateEnum { return o.state }
func (o *fakeOwner) Faction() combat.Faction      { return o.faction }

// Matches body.Obstacle DrawCollisionBox signature strictness only when needed; we
// satisfy the simpler Collidable contract used by EnemyShooting.

// recordedSpawn captures arguments of a SpawnProjectile call.
type recordedSpawn struct {
	projectileType string
	x16, y16       int
	vx16, vy16     int
	damage         int
	owner          interface{}
}

func newRecordingManager(calls *[]recordedSpawn) *mockProjectileManager {
	return &mockProjectileManager{
		SpawnProjectileFunc: func(projectileType string, x16, y16, vx16, vy16, damage int, owner interface{}) {
			*calls = append(*calls, recordedSpawn{projectileType, x16, y16, vx16, vy16, damage, owner})
		},
	}
}

// ------------------------------------------------------------------
// §8.1 Firing gates
// ------------------------------------------------------------------

func TestEnemyShooting_TryFire_Gates(t *testing.T) {
	const cooldownFrames = 10
	const projectileSpeed = 100 // fp16 per frame (output value for vx16/vy16)

	type tcase struct {
		name           string
		mode           combat.ShootMode
		direction      body.ShootDirection
		ownerPos       [2]int // pixel coords
		hasTarget      bool
		targetOffset   [2]int
		rangePx        int
		initialCooldwn int // value to force on weapon before TryFire
		stateGate      bool
		shootState     actors.ActorStateEnum
		ownerState     actors.ActorStateEnum
		wantFire       bool
		wantFaceDir    animation.FacingDirectionEnum
		wantFaceDirSet bool
	}

	tests := []tcase{
		{
			name:           "on_sight fires when target in range, left",
			mode:           combat.ShootModeOnSight,
			direction:      body.ShootDirectionStraight,
			ownerPos:       [2]int{100, 100},
			hasTarget:      true,
			targetOffset:   [2]int{-50, 0},
			rangePx:        160,
			initialCooldwn: 0,
			wantFire:       true,
			wantFaceDir:    animation.FaceDirectionLeft,
			wantFaceDirSet: true,
		},
		{
			name:           "on_sight fires when target in range, right",
			mode:           combat.ShootModeOnSight,
			direction:      body.ShootDirectionStraight,
			ownerPos:       [2]int{100, 100},
			hasTarget:      true,
			targetOffset:   [2]int{40, 0},
			rangePx:        160,
			initialCooldwn: 0,
			wantFire:       true,
			wantFaceDir:    animation.FaceDirectionRight,
			wantFaceDirSet: true,
		},
		{
			name:           "on_sight skips when out of range",
			mode:           combat.ShootModeOnSight,
			direction:      body.ShootDirectionStraight,
			ownerPos:       [2]int{100, 100},
			hasTarget:      true,
			targetOffset:   [2]int{500, 0},
			rangePx:        160,
			initialCooldwn: 0,
			wantFire:       false,
		},
		{
			name:           "on_sight skips during cooldown",
			mode:           combat.ShootModeOnSight,
			direction:      body.ShootDirectionStraight,
			ownerPos:       [2]int{100, 100},
			hasTarget:      true,
			targetOffset:   [2]int{40, 0},
			rangePx:        160,
			initialCooldwn: 5,
			wantFire:       false,
		},
		{
			name:           "on_sight skips when target nil",
			mode:           combat.ShootModeOnSight,
			direction:      body.ShootDirectionStraight,
			ownerPos:       [2]int{100, 100},
			hasTarget:      false,
			rangePx:        160,
			initialCooldwn: 0,
			wantFire:       false,
		},
		{
			name:           "always fires with no target",
			mode:           combat.ShootModeAlways,
			direction:      body.ShootDirectionDown,
			ownerPos:       [2]int{50, 50},
			hasTarget:      false,
			rangePx:        0,
			initialCooldwn: 0,
			wantFire:       true,
		},
		{
			name:           "always skips during cooldown",
			mode:           combat.ShootModeAlways,
			direction:      body.ShootDirectionDown,
			ownerPos:       [2]int{50, 50},
			hasTarget:      false,
			rangePx:        0,
			initialCooldwn: 5,
			wantFire:       false,
		},
		{
			name:           "state gate blocks fire when state mismatches (always)",
			mode:           combat.ShootModeAlways,
			direction:      body.ShootDirectionDown,
			ownerPos:       [2]int{50, 50},
			hasTarget:      false,
			rangePx:        0,
			initialCooldwn: 0,
			stateGate:      true,
			shootState:     actors.Walking,
			ownerState:     actors.Idle,
			wantFire:       false,
		},
		{
			name:           "state gate permits fire when state matches (always)",
			mode:           combat.ShootModeAlways,
			direction:      body.ShootDirectionDown,
			ownerPos:       [2]int{50, 50},
			hasTarget:      false,
			rangePx:        0,
			initialCooldwn: 0,
			stateGate:      true,
			shootState:     actors.Idle,
			ownerState:     actors.Idle,
			wantFire:       true,
		},
		{
			name:           "state gate also applies to on_sight",
			mode:           combat.ShootModeOnSight,
			direction:      body.ShootDirectionStraight,
			ownerPos:       [2]int{100, 100},
			hasTarget:      true,
			targetOffset:   [2]int{40, 0},
			rangePx:        160,
			initialCooldwn: 0,
			stateGate:      true,
			shootState:     actors.Walking,
			ownerState:     actors.Idle,
			wantFire:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner := newFakeOwner(tt.ownerPos[0], tt.ownerPos[1])
			owner.state = tt.ownerState

			var calls []recordedSpawn
			mgr := newRecordingManager(&calls)

			w := weapon.NewProjectileWeapon("enemy_w", cooldownFrames, "bullet", projectileSpeed, mgr, "", 0, 0)
			w.SetOwner(owner)
			if tt.initialCooldwn > 0 {
				w.SetCooldown(tt.initialCooldwn)
			}

			shooter := weapon.NewEnemyShooting(owner, w, tt.rangePx, tt.mode, tt.direction, tt.shootState, tt.stateGate)

			if tt.hasTarget {
				target := newFakeOwner(tt.ownerPos[0]+tt.targetOffset[0], tt.ownerPos[1]+tt.targetOffset[1])
				shooter.SetTarget(target)
			}

			got := shooter.TryFire()

			if got != tt.wantFire {
				t.Errorf("TryFire() = %v, want %v", got, tt.wantFire)
			}
			if tt.wantFire && len(calls) != 1 {
				t.Errorf("SpawnProjectile calls: got %d, want 1", len(calls))
			}
			if !tt.wantFire && len(calls) != 0 {
				t.Errorf("SpawnProjectile calls: got %d, want 0", len(calls))
			}
			if tt.wantFaceDirSet && owner.lastFaceDirSet != tt.wantFaceDir {
				t.Errorf("owner faceDir: got %v, want %v", owner.lastFaceDirSet, tt.wantFaceDir)
			}
		})
	}
}

// TestEnemyShooting_Update_CooldownTicks verifies weapon.Update is called each
// frame (cooldown decrements) even when all other gates block firing.
func TestEnemyShooting_Update_CooldownTicks(t *testing.T) {
	owner := newFakeOwner(100, 100)
	var calls []recordedSpawn
	mgr := newRecordingManager(&calls)

	w := weapon.NewProjectileWeapon("enemy_w", 10, "bullet", 100, mgr, "", 0, 0)
	w.SetOwner(owner)
	w.SetCooldown(3)

	// OnSight with no target — blocked at target gate, but cooldown must still tick.
	shooter := weapon.NewEnemyShooting(owner, w, 160, combat.ShootModeOnSight, body.ShootDirectionStraight, 0, false)
	shooter.Update()
	if w.Cooldown() != 2 {
		t.Errorf("cooldown after Update: got %d, want 2 (must tick even when firing is gated)", w.Cooldown())
	}
}

// TestEnemyShooting_OnSight_FiresAfterCooldown verifies repeated Update calls
// eventually let on_sight fire when cooldown elapses.
func TestEnemyShooting_OnSight_FiresAfterCooldown(t *testing.T) {
	owner := newFakeOwner(100, 100)
	target := newFakeOwner(140, 100)

	var calls []recordedSpawn
	mgr := newRecordingManager(&calls)

	w := weapon.NewProjectileWeapon("enemy_w", 10, "bullet", 100, mgr, "", 0, 0)
	w.SetOwner(owner)
	w.SetCooldown(1)

	shooter := weapon.NewEnemyShooting(owner, w, 160, combat.ShootModeOnSight, body.ShootDirectionStraight, 0, false)
	shooter.SetTarget(target)

	// First update: cooldown ticks from 1 -> 0, then TryFire should succeed.
	shooter.Update()

	if len(calls) != 1 {
		t.Fatalf("expected 1 spawn after cooldown elapse, got %d", len(calls))
	}
}

// ------------------------------------------------------------------
// §8.2 Direction axis mapping
// ------------------------------------------------------------------

func TestEnemyShooting_DirectionAxis(t *testing.T) {
	tests := []struct {
		name      string
		direction body.ShootDirection
		mode      combat.ShootMode
		faceDir   animation.FacingDirectionEnum
		wantVxSgn int // -1, 0, or +1
		wantVy    int // expected vy16 (>0, =0, or <0)
	}{
		{"horizontal right (wolf archetype)", body.ShootDirectionStraight, combat.ShootModeAlways, animation.FaceDirectionRight, +1, 0},
		{"horizontal left (wolf archetype)", body.ShootDirectionStraight, combat.ShootModeAlways, animation.FaceDirectionLeft, -1, 0},
		{"vertical down facing right (bat archetype)", body.ShootDirectionDown, combat.ShootModeAlways, animation.FaceDirectionRight, 0, +1},
		{"vertical down facing left (bat archetype)", body.ShootDirectionDown, combat.ShootModeAlways, animation.FaceDirectionLeft, 0, +1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner := newFakeOwner(50, 50)
			owner.face = tt.faceDir

			var calls []recordedSpawn
			mgr := newRecordingManager(&calls)
			w := weapon.NewProjectileWeapon("enemy_w", 10, "bullet", 100, mgr, "", 0, 0)
			w.SetOwner(owner)

			shooter := weapon.NewEnemyShooting(owner, w, 0, tt.mode, tt.direction, 0, false)
			if !shooter.TryFire() {
				t.Fatal("expected TryFire to succeed in always mode with cooldown ready")
			}
			if len(calls) != 1 {
				t.Fatalf("expected 1 spawn, got %d", len(calls))
			}
			got := calls[0]

			switch tt.wantVxSgn {
			case 0:
				if got.vx16 != 0 {
					t.Errorf("vx16: got %d, want 0", got.vx16)
				}
			case 1:
				if got.vx16 <= 0 {
					t.Errorf("vx16: got %d, want >0", got.vx16)
				}
			case -1:
				if got.vx16 >= 0 {
					t.Errorf("vx16: got %d, want <0", got.vx16)
				}
			}
			switch {
			case tt.wantVy == 0 && got.vy16 != 0:
				t.Errorf("vy16: got %d, want 0", got.vy16)
			case tt.wantVy > 0 && got.vy16 <= 0:
				t.Errorf("vy16: got %d, want >0", got.vy16)
			case tt.wantVy < 0 && got.vy16 >= 0:
				t.Errorf("vy16: got %d, want <0", got.vy16)
			}
		})
	}
}

// ------------------------------------------------------------------
// §8.3 Faction propagation (owner identity preserved on spawn)
// ------------------------------------------------------------------

func TestEnemyShooting_OwnerIdentityPreserved(t *testing.T) {
	owner := newFakeOwner(50, 50)
	owner.faction = combat.FactionEnemy

	var calls []recordedSpawn
	mgr := newRecordingManager(&calls)
	w := weapon.NewProjectileWeapon("enemy_w", 10, "bullet", 100, mgr, "", 0, 0)
	w.SetOwner(owner)

	shooter := weapon.NewEnemyShooting(owner, w, 0, combat.ShootModeAlways, body.ShootDirectionDown, 0, false)
	if !shooter.TryFire() {
		t.Fatal("expected TryFire to succeed")
	}
	if len(calls) != 1 {
		t.Fatalf("expected 1 spawn, got %d", len(calls))
	}
	if calls[0].owner != owner {
		t.Errorf("spawn owner: got %v, want fake enemy owner %p", calls[0].owner, owner)
	}
}

// ------------------------------------------------------------------
// Accessors returned by EnemyShooting for its configuration.
// ------------------------------------------------------------------

func TestEnemyShooting_Accessors(t *testing.T) {
	owner := newFakeOwner(0, 0)
	var calls []recordedSpawn
	mgr := newRecordingManager(&calls)
	w := weapon.NewProjectileWeapon("enemy_w", 10, "bullet", 100, mgr, "", 0, 0)

	shooter := weapon.NewEnemyShooting(owner, w, 160, combat.ShootModeOnSight, body.ShootDirectionStraight, actors.Walking, true)

	if shooter.Mode() != combat.ShootModeOnSight {
		t.Errorf("Mode(): got %v, want ShootModeOnSight", shooter.Mode())
	}
	if shooter.Direction() != body.ShootDirectionStraight {
		t.Errorf("Direction(): got %v, want ShootDirectionStraight", shooter.Direction())
	}
	if shooter.Range() != 160 {
		t.Errorf("Range(): got %d, want 160", shooter.Range())
	}
	enum, ok := shooter.ShootState()
	if !ok || enum != actors.Walking {
		t.Errorf("ShootState(): got (%v,%v), want (Walking,true)", enum, ok)
	}

	// Contract: Target() returns the last value passed to SetTarget, or nil.
	if shooter.Target() != nil {
		t.Errorf("Target() before SetTarget: got %v, want nil", shooter.Target())
	}
	target := newFakeOwner(10, 10)
	shooter.SetTarget(target)
	if shooter.Target() != target {
		t.Errorf("Target() after SetTarget: got %v, want %v", shooter.Target(), target)
	}
}
