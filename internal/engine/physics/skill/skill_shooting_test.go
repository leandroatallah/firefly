package skill_test

import (
	"image"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/mocks"
	"github.com/boilerplate/ebiten-template/internal/engine/physics/skill"
	"github.com/hajimehoshi/ebiten/v2"
)

func TestShootingSkill_CooldownGating(t *testing.T) {
	spawnCount := 0
	shooter := &mocks.MockShooter{
		SpawnBulletFunc: func(x16, y16, vx16, vy16 int, owner interface{}) {
			spawnCount++
		},
	}

	s := skill.NewShootingSkill(shooter, 3, 16<<4, 32<<4, 4)
	body := &mockMovableCollidable{
		getPosition16Func: func() (int, int) { return 100 << 4, 50 << 4 },
		faceDirectionFunc: func() animation.FacingDirectionEnum { return animation.FaceDirectionRight },
	}

	s.HandleInputWithDirection(body, nil, nil, false, false, false, false)
	s.Update(body, nil)
	s.HandleInputWithDirection(body, nil, nil, false, false, false, false)
	s.Update(body, nil)
	s.HandleInputWithDirection(body, nil, nil, false, false, false, false)
	s.Update(body, nil)
	s.HandleInputWithDirection(body, nil, nil, false, false, false, false)
	s.Update(body, nil)

	if spawnCount != 2 {
		t.Errorf("with cooldown=3, 4 frames should spawn 2 bullets (frame 1 and 4), got %d", spawnCount)
	}
}

func TestShootingSkill_AlternatingYOffset(t *testing.T) {
	var yOffsets []int
	shooter := &mocks.MockShooter{
		SpawnBulletFunc: func(x16, y16, vx16, vy16 int, owner interface{}) {
			yOffsets = append(yOffsets, y16)
		},
	}

	s := skill.NewShootingSkill(shooter, 0, 0, 32<<4, 4)
	body := &mockMovableCollidable{
		getPosition16Func: func() (int, int) { return 0, 50 << 4 },
		faceDirectionFunc: func() animation.FacingDirectionEnum { return animation.FaceDirectionRight },
	}

	for i := 0; i < 4; i++ {
		s.HandleInputWithDirection(body, nil, nil, false, false, false, false)
		s.Update(body, nil)
	}

	want := []int{(50 << 4) + 4, (50 << 4) - 4, (50 << 4) + 4, (50 << 4) - 4}
	if len(yOffsets) != len(want) {
		t.Fatalf("got %d offsets, want %d", len(yOffsets), len(want))
	}
	for i, got := range yOffsets {
		if got != want[i] {
			t.Errorf("offset[%d]: got %d, want %d", i, got, want[i])
		}
	}
}

func TestShootingSkill_StateTransitions(t *testing.T) {
	shooter := &mocks.MockShooter{
		SpawnBulletFunc: func(x16, y16, vx16, vy16 int, owner interface{}) {},
	}

	s := skill.NewShootingSkill(shooter, 2, 0, 32<<4, 4)
	body := &mockMovableCollidable{
		getPosition16Func: func() (int, int) { return 0, 0 },
		faceDirectionFunc: func() animation.FacingDirectionEnum { return animation.FaceDirectionRight },
	}

	if s.IsActive() {
		t.Error("initial state should be Ready (IsActive=false)")
	}

	s.HandleInputWithDirection(body, nil, nil, false, false, false, false)

	if !s.IsActive() {
		t.Error("after HandleInputWithDirection with spawn, state should be Cooldown (IsActive=true)")
	}

	s.Update(body, nil)
	s.Update(body, nil)

	if s.IsActive() {
		t.Error("after 2 Update calls, state should return to Ready (IsActive=false)")
	}
}

func TestShootingSkill_NoSpawnWhenNotReady(t *testing.T) {
	spawnCount := 0
	shooter := &mocks.MockShooter{
		SpawnBulletFunc: func(x16, y16, vx16, vy16 int, owner interface{}) {
			spawnCount++
		},
	}

	s := skill.NewShootingSkill(shooter, 3, 0, 32<<4, 4)
	body := &mockMovableCollidable{
		getPosition16Func: func() (int, int) { return 0, 0 },
		faceDirectionFunc: func() animation.FacingDirectionEnum { return animation.FaceDirectionRight },
	}

	s.HandleInputWithDirection(body, nil, nil, false, false, false, false)
	if spawnCount != 1 {
		t.Fatalf("first call should spawn, got %d spawns", spawnCount)
	}

	s.HandleInputWithDirection(body, nil, nil, false, false, false, false)
	if spawnCount != 1 {
		t.Errorf("second call during cooldown should not spawn, got %d spawns", spawnCount)
	}
}

type mockMovableCollidable struct {
	getPosition16Func func() (int, int)
	faceDirectionFunc func() animation.FacingDirectionEnum
	isDuckingFunc     func() bool
	width             int
}

type mockShape struct {
	width int
}

func (s *mockShape) Width() int  { return s.width }
func (s *mockShape) Height() int { return 0 }

func (m *mockMovableCollidable) GetPosition16() (int, int) {
	if m.getPosition16Func != nil {
		return m.getPosition16Func()
	}
	return 0, 0
}

func (m *mockMovableCollidable) FaceDirection() animation.FacingDirectionEnum {
	if m.faceDirectionFunc != nil {
		return m.faceDirectionFunc()
	}
	return animation.FaceDirectionRight
}

func (m *mockMovableCollidable) IsDucking() bool {
	if m.isDuckingFunc != nil {
		return m.isDuckingFunc()
	}
	return false
}

func (m *mockMovableCollidable) MoveX(distance int)                              {}
func (m *mockMovableCollidable) MoveY(distance int)                              {}
func (m *mockMovableCollidable) OnMoveLeft(distance int)                         {}
func (m *mockMovableCollidable) OnMoveUpLeft(distance int)                       {}
func (m *mockMovableCollidable) OnMoveDownLeft(distance int)                     {}
func (m *mockMovableCollidable) OnMoveRight(distance int)                        {}
func (m *mockMovableCollidable) OnMoveUpRight(distance int)                      {}
func (m *mockMovableCollidable) OnMoveDownRight(distance int)                    {}
func (m *mockMovableCollidable) OnMoveUp(distance int)                           {}
func (m *mockMovableCollidable) OnMoveDown(distance int)                         {}
func (m *mockMovableCollidable) Velocity() (vx16, vy16 int)                      { return 0, 0 }
func (m *mockMovableCollidable) SetVelocity(vx16, vy16 int)                      {}
func (m *mockMovableCollidable) Acceleration() (accX, accY int)                  { return 0, 0 }
func (m *mockMovableCollidable) SetAcceleration(accX, accY int)                  {}
func (m *mockMovableCollidable) SetSpeed(speed int) error                        { return nil }
func (m *mockMovableCollidable) SetMaxSpeed(maxSpeed int) error                  { return nil }
func (m *mockMovableCollidable) GetPosition() (int, int)                         { return 0, 0 }
func (m *mockMovableCollidable) SetPosition(x, y int)                            {}
func (m *mockMovableCollidable) GetBounds() (x, y, width, height int)            { return 0, 0, 0, 0 }
func (m *mockMovableCollidable) SetBounds(x, y, width, height int)               {}
func (m *mockMovableCollidable) IsCollidingWith(interface{}) bool                { return false }
func (m *mockMovableCollidable) OnCollision(interface{})                         {}
func (m *mockMovableCollidable) GetTouchable() body.Touchable                    { return m }
func (m *mockMovableCollidable) DrawCollisionBox(*ebiten.Image, image.Rectangle) {}
func (m *mockMovableCollidable) CollisionPosition() []image.Rectangle            { return nil }
func (m *mockMovableCollidable) CollisionShapes() []body.Collidable              { return nil }
func (m *mockMovableCollidable) IsObstructive() bool                             { return false }
func (m *mockMovableCollidable) SetIsObstructive(bool)                           {}
func (m *mockMovableCollidable) AddCollision(...body.Collidable)                 {}
func (m *mockMovableCollidable) ClearCollisions()                                {}
func (m *mockMovableCollidable) SetTouchable(body.Touchable)                     {}
func (m *mockMovableCollidable) OnTouch(body.Collidable)                         {}
func (m *mockMovableCollidable) OnBlock(body.Collidable)                         {}
func (m *mockMovableCollidable) ID() string                                      { return "" }
func (m *mockMovableCollidable) SetID(string)                                    {}
func (m *mockMovableCollidable) Position() image.Rectangle                       { return image.Rectangle{} }
func (m *mockMovableCollidable) SetPosition16(x16, y16 int)                      {}
func (m *mockMovableCollidable) SetSize(w, h int)                                {}
func (m *mockMovableCollidable) Scale() float64                                  { return 1.0 }
func (m *mockMovableCollidable) SetScale(float64)                                {}
func (m *mockMovableCollidable) GetPositionMin() (int, int)                      { return 0, 0 }
func (m *mockMovableCollidable) GetShape() body.Shape {
	if m.width > 0 {
		return &mockShape{width: m.width}
	}
	return &mockShape{width: 16} // default width
}
func (m *mockMovableCollidable) ApplyValidPosition(int, bool, body.BodiesSpace) (int, int, bool) {
	return 0, 0, false
}
func (m *mockMovableCollidable) Owner() interface{}                             { return nil }
func (m *mockMovableCollidable) SetOwner(interface{})                           {}
func (m *mockMovableCollidable) LastOwner() interface{}                         { return nil }
func (m *mockMovableCollidable) Speed() int                                     { return 0 }
func (m *mockMovableCollidable) MaxSpeed() int                                  { return 0 }
func (m *mockMovableCollidable) Immobile() bool                                 { return false }
func (m *mockMovableCollidable) SetImmobile(bool)                               {}
func (m *mockMovableCollidable) SetFreeze(bool)                                 {}
func (m *mockMovableCollidable) Freeze() bool                                   { return false }
func (m *mockMovableCollidable) SetFaceDirection(animation.FacingDirectionEnum) {}
func (m *mockMovableCollidable) IsIdle() bool                                   { return false }
func (m *mockMovableCollidable) IsWalking() bool                                { return false }
func (m *mockMovableCollidable) IsFalling() bool                                { return false }
func (m *mockMovableCollidable) IsGoingUp() bool                                { return false }
func (m *mockMovableCollidable) CheckMovementDirectionX()                       {}
func (m *mockMovableCollidable) TryJump(int)                                    {}
func (m *mockMovableCollidable) SetJumpForceMultiplier(float64)                 {}
func (m *mockMovableCollidable) JumpForceMultiplier() float64                   { return 1.0 }
func (m *mockMovableCollidable) SetHorizontalInertia(float64)                   {}
func (m *mockMovableCollidable) HorizontalInertia() float64                     { return 1.0 }
