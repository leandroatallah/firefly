// Package platformerphasescene_test holds Red-Phase TDD tests for story
// 055-kit-genre-phase-scenes. These tests intentionally fail until the
// production types in package platformerphasescene (this directory) are
// introduced — they exercise the observable behaviour described in
// SPEC.md §4.1 and §8.1.
package platformerphasescene_test

import (
	"image"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/movement"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"

	platformerphasescene "github.com/boilerplate/ebiten-template/internal/kit/scenes/phases/platformer"
	"github.com/hajimehoshi/ebiten/v2"
)

// --- mock player satisfying platformer.PlatformerActorEntity ---------------

type mockPlatformerPlayer struct {
	x16, y16, w, h        int
	state                 actors.ActorStateEnum
	character             *actors.Character
	setImmobileTrueCalls  int
	setImmobileFalseCalls int
}

func newMockPlatformerPlayer(x16, y16 int) *mockPlatformerPlayer {
	return &mockPlatformerPlayer{x16: x16, y16: y16, w: 16, h: 16}
}

// body.Body
func (m *mockPlatformerPlayer) ID() string   { return "player" }
func (m *mockPlatformerPlayer) SetID(string) {}
func (m *mockPlatformerPlayer) Position() image.Rectangle {
	return image.Rect(m.x16/16, m.y16/16, m.x16/16+m.w, m.y16/16+m.h)
}
func (m *mockPlatformerPlayer) SetPosition(x, y int)       { m.x16, m.y16 = x*16, y*16 }
func (m *mockPlatformerPlayer) SetPosition16(x16, y16 int) { m.x16, m.y16 = x16, y16 }
func (m *mockPlatformerPlayer) SetSize(w, h int)           { m.w, m.h = w, h }
func (m *mockPlatformerPlayer) Scale() float64             { return 1 }
func (m *mockPlatformerPlayer) SetScale(float64)           {}
func (m *mockPlatformerPlayer) GetPosition16() (int, int)  { return m.x16, m.y16 }
func (m *mockPlatformerPlayer) GetPositionMin() (int, int) { return m.x16 / 16, m.y16 / 16 }
func (m *mockPlatformerPlayer) GetShape() body.Shape       { return m }
func (m *mockPlatformerPlayer) Width() int                 { return m.w }
func (m *mockPlatformerPlayer) Height() int                { return m.h }
func (m *mockPlatformerPlayer) Owner() interface{}         { return nil }
func (m *mockPlatformerPlayer) SetOwner(interface{})       {}
func (m *mockPlatformerPlayer) LastOwner() interface{}     { return nil }
func (m *mockPlatformerPlayer) Altitude() int              { return 0 }
func (m *mockPlatformerPlayer) SetAltitude(int)            {}
func (m *mockPlatformerPlayer) Altitude16() int            { return 0 }
func (m *mockPlatformerPlayer) SetAltitude16(int)          {}

// body.Touchable / body.Collidable
func (m *mockPlatformerPlayer) OnTouch(body.Collidable)                         {}
func (m *mockPlatformerPlayer) OnBlock(body.Collidable)                         {}
func (m *mockPlatformerPlayer) GetTouchable() body.Touchable                    { return m }
func (m *mockPlatformerPlayer) DrawCollisionBox(*ebiten.Image, image.Rectangle) {}
func (m *mockPlatformerPlayer) CollisionPosition() []image.Rectangle            { return nil }
func (m *mockPlatformerPlayer) CollisionShapes() []body.Collidable              { return nil }
func (m *mockPlatformerPlayer) IsObstructive() bool                             { return false }
func (m *mockPlatformerPlayer) SetIsObstructive(bool)                           {}
func (m *mockPlatformerPlayer) AddCollision(...body.Collidable)                 {}
func (m *mockPlatformerPlayer) ClearCollisions()                                {}
func (m *mockPlatformerPlayer) SetTouchable(body.Touchable)                     {}
func (m *mockPlatformerPlayer) ApplyValidPosition(_ int, _ bool, _ body.BodiesSpace) (int, int, bool) {
	return m.x16 / 16, m.y16 / 16, false
}

// body.Drawable
func (m *mockPlatformerPlayer) Image() *ebiten.Image { return ebiten.NewImage(1, 1) }
func (m *mockPlatformerPlayer) ImageOptions() *ebiten.DrawImageOptions {
	return &ebiten.DrawImageOptions{}
}
func (m *mockPlatformerPlayer) UpdateImageOptions() {}

// Stateful / Damageable / Controllable
func (m *mockPlatformerPlayer) State() actors.ActorStateEnum { return m.state }
func (m *mockPlatformerPlayer) SetState(actors.ActorState)   {}
func (m *mockPlatformerPlayer) NewState(actors.ActorStateEnum) (actors.ActorState, error) {
	return nil, nil
}
func (m *mockPlatformerPlayer) SetMovementState(_ movement.MovementStateEnum, _ body.MovableCollidable, _ ...movement.MovementStateOption) {
}
func (m *mockPlatformerPlayer) SwitchMovementState(movement.MovementStateEnum) {}
func (m *mockPlatformerPlayer) MovementState() movement.MovementState          { return nil }
func (m *mockPlatformerPlayer) Hurt(int)                                       {}
func (m *mockPlatformerPlayer) OnMoveLeft(int)                                 {}
func (m *mockPlatformerPlayer) OnMoveRight(int)                                {}
func (m *mockPlatformerPlayer) BlockMovement()                                 {}
func (m *mockPlatformerPlayer) UnblockMovement()                               {}
func (m *mockPlatformerPlayer) IsMovementBlocked() bool                        { return false }

// MovableCollidableAlive
func (m *mockPlatformerPlayer) Velocity() (float64, float64) { return 0, 0 }
func (m *mockPlatformerPlayer) SetVelocity(float64, float64) {}
func (m *mockPlatformerPlayer) IsAlive() bool                { return m.state != actors.Dead }
func (m *mockPlatformerPlayer) SetImmobile(v bool) {
	if v {
		m.setImmobileTrueCalls++
	} else {
		m.setImmobileFalseCalls++
	}
}
func (m *mockPlatformerPlayer) IsImmobile() bool { return false }
func (m *mockPlatformerPlayer) SetFreeze(bool)   {}
func (m *mockPlatformerPlayer) IsFrozen() bool   { return false }

// ActorEntity
func (m *mockPlatformerPlayer) Update(body.BodiesSpace) error                  { return nil }
func (m *mockPlatformerPlayer) MovementModel() physicsmovement.MovementModel   { return nil }
func (m *mockPlatformerPlayer) SetMovementModel(physicsmovement.MovementModel) {}
func (m *mockPlatformerPlayer) GetCharacter() *actors.Character                { return m.character }

// Platformer-specific (PlatformerActorEntity)
func (m *mockPlatformerPlayer) OnDie()                      {}
func (m *mockPlatformerPlayer) OnJump()                     {}
func (m *mockPlatformerPlayer) OnLand()                     {}
func (m *mockPlatformerPlayer) OnFall()                     {}
func (m *mockPlatformerPlayer) SetOnJump(func(image.Point)) {}
func (m *mockPlatformerPlayer) SetOnFall(func(image.Point)) {}
func (m *mockPlatformerPlayer) SetOnLand(func(image.Point)) {}
func (m *mockPlatformerPlayer) AppContext() *app.AppContext { return nil }
func (m *mockPlatformerPlayer) SetAppContext(any)           {}

// recordingFatalCharacter wraps actors.Character to record SetNewStateFatal.
// We intercept by counting calls on the player; the production code must call
// player.GetCharacter().SetNewStateFatal — we substitute the character with a
// fake whose method bumps a counter on the player.
//
// Because actors.Character is a concrete type with embedded behaviour, we
// simulate behaviour observation by also exposing a direct hook the production
// code can call (recorded for assertion).
//
// The Red-Phase test asserts on these counters via the helper accessors below.

// --- tests -----------------------------------------------------------------

func TestPlatformerPhaseScene_ScreenFlipperCallbacksToggleImmobility(t *testing.T) {
	scene := platformerphasescene.NewForTest(platformerphasescene.TestOptions{
		CameraCenterX: 100,
		CameraCenterY: 100,
		ScreenWidth:   320,
		ScreenHeight:  200,
		HasFlipper:    true,
		DyingState:    actors.Dying,
		DeadState:     actors.Dead,
	})

	player := newMockPlatformerPlayer(0, 0)
	scene.SetPlayerForTest(player)

	scene.InvokeScreenFlipperOnFlipStartForTest()
	if player.setImmobileTrueCalls != 1 {
		t.Fatalf("expected SetImmobile(true) once on OnFlipStart, got %d", player.setImmobileTrueCalls)
	}
	if player.setImmobileFalseCalls != 0 {
		t.Fatalf("expected SetImmobile(false) NOT called on OnFlipStart, got %d", player.setImmobileFalseCalls)
	}

	scene.InvokeScreenFlipperOnFlipFinishForTest()
	if player.setImmobileFalseCalls != 1 {
		t.Fatalf("expected SetImmobile(false) once on OnFlipFinish, got %d", player.setImmobileFalseCalls)
	}
}

func TestPlatformerPhaseScene_NoPlayer_DoesNotPanic(t *testing.T) {
	scene := platformerphasescene.NewForTest(platformerphasescene.TestOptions{
		CameraCenterX:          100,
		CameraCenterY:          100,
		ScreenWidth:            320,
		ScreenHeight:           200,
		HasPlayerStartPosition: false,
		DyingState:             actors.Dying,
		DeadState:              actors.Dead,
	})

	// Must not panic.
	scene.OnStart()
	if err := scene.Update(); err != nil {
		t.Fatalf("Update() unexpected error: %v", err)
	}

	headless := ebiten.NewImage(320, 200)
	scene.Draw(headless)

	if scene.ScreenFlipperForTest() != nil {
		t.Error("expected screenFlipper to be nil when tilemap has no player start")
	}
	if !scene.CameraIsFixedModeForTest() {
		t.Error("expected camera to be in fixed mode when no player")
	}
}

func TestPlatformerPhaseScene_DebugDrawHookInvoked(t *testing.T) {
	scene := platformerphasescene.NewForTest(platformerphasescene.TestOptions{
		CameraCenterX: 100,
		CameraCenterY: 100,
		ScreenWidth:   320,
		ScreenHeight:  200,
		DyingState:    actors.Dying,
		DeadState:     actors.Dead,
	})

	count := 0
	scene.SetDebugDrawHook(func(*ebiten.Image) {
		count++
	})

	headless := ebiten.NewImage(320, 200)
	scene.Draw(headless)

	if count != 1 {
		t.Fatalf("expected DebugDrawHook to be invoked once during Draw, got %d", count)
	}
}
