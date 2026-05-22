// Package beatemupphasescene_test holds Red-Phase TDD tests for story
// 055-kit-genre-phase-scenes. These tests intentionally fail until the
// production types in package beatemupphasescene are introduced — they
// exercise the observable behaviour described in SPEC.md §4.2–§4.3 and §8.2.
package beatemupphasescene_test

import (
	"image"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/movement"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"

	beatemupphasescene "github.com/boilerplate/ebiten-template/internal/kit/scenes/phases/beatemup"
	"github.com/hajimehoshi/ebiten/v2"
)

// --- mock BeatEmUp actor satisfying beatemupkit.BeatEmUpActorEntity --------

type mockBeatEmUpActor struct {
	id                      string
	x16, y16, w, h          int
	alt16                   int
	state                   actors.ActorStateEnum
	setNewStateFatalCalls   int
	lastSetNewStateFatalArg actors.ActorStateEnum
}

func newMockBeatEmUpActor(id string, y16, alt16 int) *mockBeatEmUpActor {
	return &mockBeatEmUpActor{id: id, x16: 0, y16: y16, alt16: alt16, w: 16, h: 16}
}

// body.Body
func (m *mockBeatEmUpActor) ID() string      { return m.id }
func (m *mockBeatEmUpActor) SetID(id string) { m.id = id }
func (m *mockBeatEmUpActor) Position() image.Rectangle {
	return image.Rect(m.x16/16, m.y16/16, m.x16/16+m.w, m.y16/16+m.h)
}
func (m *mockBeatEmUpActor) SetPosition(x, y int)       { m.x16, m.y16 = x*16, y*16 }
func (m *mockBeatEmUpActor) SetPosition16(x16, y16 int) { m.x16, m.y16 = x16, y16 }
func (m *mockBeatEmUpActor) SetSize(w, h int)           { m.w, m.h = w, h }
func (m *mockBeatEmUpActor) Scale() float64             { return 1 }
func (m *mockBeatEmUpActor) SetScale(float64)           {}
func (m *mockBeatEmUpActor) GetPosition16() (int, int)  { return m.x16, m.y16 }
func (m *mockBeatEmUpActor) GetPositionMin() (int, int) { return m.x16 / 16, m.y16 / 16 }
func (m *mockBeatEmUpActor) GetShape() body.Shape       { return m }
func (m *mockBeatEmUpActor) Width() int                 { return m.w }
func (m *mockBeatEmUpActor) Height() int                { return m.h }
func (m *mockBeatEmUpActor) Owner() interface{}         { return nil }
func (m *mockBeatEmUpActor) SetOwner(interface{})       {}
func (m *mockBeatEmUpActor) LastOwner() interface{}     { return nil }

// Altitude
func (m *mockBeatEmUpActor) Altitude() int       { return m.alt16 / 16 }
func (m *mockBeatEmUpActor) SetAltitude(a int)   { m.alt16 = a * 16 }
func (m *mockBeatEmUpActor) Altitude16() int     { return m.alt16 }
func (m *mockBeatEmUpActor) SetAltitude16(a int) { m.alt16 = a }

// Touchable / Collidable
func (m *mockBeatEmUpActor) OnTouch(body.Collidable)                         {}
func (m *mockBeatEmUpActor) OnBlock(body.Collidable)                         {}
func (m *mockBeatEmUpActor) GetTouchable() body.Touchable                    { return m }
func (m *mockBeatEmUpActor) DrawCollisionBox(*ebiten.Image, image.Rectangle) {}
func (m *mockBeatEmUpActor) CollisionPosition() []image.Rectangle            { return nil }
func (m *mockBeatEmUpActor) CollisionShapes() []body.Collidable              { return nil }
func (m *mockBeatEmUpActor) IsObstructive() bool                             { return false }
func (m *mockBeatEmUpActor) SetIsObstructive(bool)                           {}
func (m *mockBeatEmUpActor) AddCollision(...body.Collidable)                 {}
func (m *mockBeatEmUpActor) ClearCollisions()                                {}
func (m *mockBeatEmUpActor) SetTouchable(body.Touchable)                     {}
func (m *mockBeatEmUpActor) ApplyValidPosition(_ int, _ bool, _ body.BodiesSpace) (int, int, bool) {
	return m.x16 / 16, m.y16 / 16, false
}

// Drawable
func (m *mockBeatEmUpActor) Image() *ebiten.Image { return ebiten.NewImage(1, 1) }
func (m *mockBeatEmUpActor) ImageOptions() *ebiten.DrawImageOptions {
	return &ebiten.DrawImageOptions{}
}
func (m *mockBeatEmUpActor) UpdateImageOptions() {}

// Stateful / Damageable / Controllable
func (m *mockBeatEmUpActor) State() actors.ActorStateEnum { return m.state }
func (m *mockBeatEmUpActor) SetState(actors.ActorState)   {}
func (m *mockBeatEmUpActor) NewState(actors.ActorStateEnum) (actors.ActorState, error) {
	return nil, nil
}
func (m *mockBeatEmUpActor) SetMovementState(_ movement.MovementStateEnum, _ body.MovableCollidable, _ ...movement.MovementStateOption) {
}
func (m *mockBeatEmUpActor) SwitchMovementState(movement.MovementStateEnum) {}
func (m *mockBeatEmUpActor) MovementState() movement.MovementState          { return nil }
func (m *mockBeatEmUpActor) Hurt(int)                                       {}
func (m *mockBeatEmUpActor) OnMoveLeft(int)                                 {}
func (m *mockBeatEmUpActor) OnMoveRight(int)                                {}
func (m *mockBeatEmUpActor) BlockMovement()                                 {}
func (m *mockBeatEmUpActor) UnblockMovement()                               {}
func (m *mockBeatEmUpActor) IsMovementBlocked() bool                        { return false }

// MovableCollidableAlive
func (m *mockBeatEmUpActor) Velocity() (float64, float64) { return 0, 0 }
func (m *mockBeatEmUpActor) SetVelocity(float64, float64) {}
func (m *mockBeatEmUpActor) IsAlive() bool                { return m.state != actors.Dead }
func (m *mockBeatEmUpActor) SetImmobile(bool)             {}
func (m *mockBeatEmUpActor) IsImmobile() bool             { return false }
func (m *mockBeatEmUpActor) SetFreeze(bool)               {}
func (m *mockBeatEmUpActor) IsFrozen() bool               { return false }

// ActorEntity
func (m *mockBeatEmUpActor) Update(body.BodiesSpace) error                  { return nil }
func (m *mockBeatEmUpActor) MovementModel() physicsmovement.MovementModel   { return nil }
func (m *mockBeatEmUpActor) SetMovementModel(physicsmovement.MovementModel) {}
func (m *mockBeatEmUpActor) GetCharacter() *actors.Character                { return nil }

// recordingFatalCharacter behaviour: tests use a sidecar recorder hook on
// the scene so the production code can route SetNewStateFatal calls through
// it for assertion. (No-fall-death test relies on this never being called.)
func (m *mockBeatEmUpActor) recordSetNewStateFatal(s actors.ActorStateEnum) {
	m.setNewStateFatalCalls++
	m.lastSetNewStateFatalArg = s
}

// --- tests -----------------------------------------------------------------

func TestBeatemupPhaseScene_DrawOrderSortsByGroundYPlusAltitude(t *testing.T) {
	scene := beatemupphasescene.NewForTest(beatemupphasescene.TestOptions{
		ScreenWidth:  320,
		ScreenHeight: 200,
	})

	// Spec example: A(y16=10, alt=0)→eff 10, B(y16=5, alt=8)→eff -3, C(y16=20, alt=5)→eff 15.
	a := newMockBeatEmUpActor("A", 10, 0)
	b := newMockBeatEmUpActor("B", 5, 8)
	c := newMockBeatEmUpActor("C", 20, 5)
	scene.AddBodyForTest(a)
	scene.AddBodyForTest(b)
	scene.AddBodyForTest(c)

	var drawOrder []string
	scene.SetActorDrawHandlerForTest(func(_ *ebiten.Image, bd body.Collidable) bool {
		drawOrder = append(drawOrder, bd.ID())
		return true
	})

	headless := ebiten.NewImage(320, 200)
	scene.Draw(headless)

	want := []string{"B", "A", "C"}
	if len(drawOrder) != len(want) {
		t.Fatalf("draw order length = %d, want %d (got %v)", len(drawOrder), len(want), drawOrder)
	}
	for i := range want {
		if drawOrder[i] != want[i] {
			t.Fatalf("draw order = %v; want %v", drawOrder, want)
		}
	}
}

func TestBeatemupPhaseScene_RemovesDeadActor_NoAltitudePanic(t *testing.T) {
	scene := beatemupphasescene.NewForTest(beatemupphasescene.TestOptions{
		ScreenWidth:  320,
		ScreenHeight: 200,
	})

	dead := newMockBeatEmUpActor("dead", 0, 0)
	dead.state = actors.Dead
	scene.AddBodyForTest(dead)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Update() panicked while removing Dead beat-em-up actor: %v", r)
		}
	}()

	if err := scene.Update(); err != nil {
		t.Fatalf("Update() error: %v", err)
	}

	if scene.SpaceContainsBodyForTest(dead) {
		t.Error("expected RemoveBody(actor) for Dead beat-em-up actor; actor still in space")
	}
}

func TestBeatemupPhaseScene_NoFallDeathPath(t *testing.T) {
	scene := beatemupphasescene.NewForTest(beatemupphasescene.TestOptions{
		// Camera center (100,100); screen height 200 → bottom 200.
		CameraCenterX: 100,
		CameraCenterY: 100,
		ScreenWidth:   320,
		ScreenHeight:  200,
	})

	// Place player far below camera bottom — would trigger fall-death in platformer.
	player := newMockBeatEmUpActor("player", 1000*16, 0)
	scene.SetPlayerForTest(player)
	scene.SetSetNewStateFatalRecorder(player.recordSetNewStateFatal)

	if err := scene.Update(); err != nil {
		t.Fatalf("Update() error: %v", err)
	}

	if scene.DeathActiveForTest() {
		t.Error("beat-em-up scene must NOT trigger fall-death when player is below camera bottom")
	}
	if player.setNewStateFatalCalls != 0 {
		t.Errorf("beat-em-up must NOT call SetNewStateFatal on fall path; got %d calls",
			player.setNewStateFatalCalls)
	}
}

func TestBeatemupPhaseScene_CameraNotVerticalOnlyUpward(t *testing.T) {
	scene := beatemupphasescene.NewForTest(beatemupphasescene.TestOptions{
		CameraCenterX:          0,
		CameraCenterY:          0,
		ScreenWidth:            320,
		ScreenHeight:           200,
		HasPlayerStartPosition: true,
	})

	player := newMockBeatEmUpActor("player", 0, 0)
	scene.SetPlayerForTest(player)

	scene.OnStart()

	cam := scene.EngineCameraForTest()
	if cam == nil {
		t.Fatal("expected engine camera to be initialised after OnStart")
	}
	if cam.VerticalOnlyUpward {
		t.Error("beat-em-up camera must NOT have VerticalOnlyUpward set; got true")
	}
}

// --- T-S8 (063-shadow-component): shadow drawn before actor sprite ---------

func TestBeatemupPhaseScene_ShadowDrawnBeforeActor(t *testing.T) {
	scene := beatemupphasescene.NewForTest(beatemupphasescene.TestOptions{
		ScreenWidth:  320,
		ScreenHeight: 200,
	})

	// One airborne body (alt16 > 0).
	airborne := newMockBeatEmUpActor("air", 0, 16)
	scene.AddBodyForTest(airborne)

	var order []string
	scene.SetShadowDrawerForTest(func(_ *ebiten.Image, bodies []body.Collidable) {
		for _, b := range bodies {
			order = append(order, "shadow("+b.ID()+")")
		}
	})
	scene.SetActorDrawHandlerForTest(func(_ *ebiten.Image, b body.Collidable) bool {
		order = append(order, "actor("+b.ID()+")")
		return true
	})

	headless := ebiten.NewImage(320, 200)
	scene.DrawActors(headless)

	want := []string{"shadow(air)", "actor(air)"}
	if len(order) != len(want) {
		t.Fatalf("draw order length = %d, want %d (got %v)", len(order), len(want), order)
	}
	for i := range want {
		if order[i] != want[i] {
			t.Fatalf("draw order = %v; want %v", order, want)
		}
	}
}
