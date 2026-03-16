package movement

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/contracts/tilemaplayer"
)

type mockActor struct {
	pos            image.Rectangle
	speed          int
	moveLeftForce  int
	moveRightForce int
	moveUpForce    int
	moveDownForce  int
	immobile       bool
	obstructive    bool
	id             string
	velocity       struct{ x, y int }
}

func (m *mockActor) ID() string { 
	if m.id != "" {
		return m.id
	}
	return "mock" 
}
func (m *mockActor) SetID(id string) { m.id = id }
func (m *mockActor) Position() image.Rectangle { return m.pos }
func (m *mockActor) SetPosition(x, y int) { m.pos = image.Rect(x, y, x+10, y+10) }
func (m *mockActor) SetSize(w, h int) {
	m.pos.Max.X = m.pos.Min.X + w
	m.pos.Max.Y = m.pos.Min.Y + h
}
func (m *mockActor) Scale() float64 { return 1.0 }
func (m *mockActor) SetScale(s float64) {}
func (m *mockActor) SetPosition16(x16, y16 int) {}
func (m *mockActor) GetPosition16() (int, int) { return m.pos.Min.X * 16, m.pos.Min.Y * 16 }
func (m *mockActor) GetPositionMin() (int, int) { return m.pos.Min.X, m.pos.Min.Y }
func (m *mockActor) GetShape() body.Shape { return nil }
func (m *mockActor) Owner() interface{} { return nil }
func (m *mockActor) SetOwner(o interface{}) {}
func (m *mockActor) LastOwner() interface{} { return nil }

func (m *mockActor) MoveX(d int) { m.moveRightForce = d }
func (m *mockActor) MoveY(d int) { m.moveDownForce = d }
func (m *mockActor) OnMoveLeft(d int) { m.moveLeftForce = d }
func (m *mockActor) OnMoveUpLeft(d int) { m.moveLeftForce = d; m.moveUpForce = d }
func (m *mockActor) OnMoveDownLeft(d int) { m.moveLeftForce = d; m.moveDownForce = d }
func (m *mockActor) OnMoveRight(d int) { m.moveRightForce = d }
func (m *mockActor) OnMoveUpRight(d int) { m.moveRightForce = d; m.moveUpForce = d }
func (m *mockActor) OnMoveDownRight(d int) { m.moveRightForce = d; m.moveDownForce = d }
func (m *mockActor) OnMoveUp(d int) { m.moveUpForce = d }
func (m *mockActor) OnMoveDown(d int) { m.moveDownForce = d }

func (m *mockActor) Velocity() (int, int) { return m.velocity.x, m.velocity.y }
func (m *mockActor) SetVelocity(vx, vy int) { m.velocity.x, m.velocity.y = vx, vy }
func (m *mockActor) Acceleration() (ax, ay int) { return 0, 0 }
func (m *mockActor) SetAcceleration(ax, ay int) {}
func (m *mockActor) SetSpeed(s int) error { m.speed = s; return nil }
func (m *mockActor) SetMaxSpeed(s int) error { return nil }
func (m *mockActor) Speed() int { return m.speed }
func (m *mockActor) MaxSpeed() int { return 0 }
func (m *mockActor) Immobile() bool { return m.immobile }
func (m *mockActor) SetImmobile(i bool) { m.immobile = i }
func (m *mockActor) SetFreeze(f bool) {}
func (m *mockActor) Freeze() bool { return false }
func (m *mockActor) FaceDirection() animation.FacingDirectionEnum { return 0 }
func (m *mockActor) SetFaceDirection(v animation.FacingDirectionEnum) {}
func (m *mockActor) IsIdle() bool { return true }
func (m *mockActor) IsWalking() bool { return false }
func (m *mockActor) IsFalling() bool { return false }
func (m *mockActor) IsGoingUp() bool { return false }
func (m *mockActor) CheckMovementDirectionX() {}
func (m *mockActor) TryJump(f int) { m.moveRightForce = f }
func (m *mockActor) SetJumpForceMultiplier(mu float64) {}
func (m *mockActor) JumpForceMultiplier() float64 { return 1.0 }
func (m *mockActor) SetHorizontalInertia(i float64) {}
func (m *mockActor) HorizontalInertia() float64 { return 1.0 }

// Collidable methods
func (m *mockActor) OnTouch(other body.Collidable) {}
func (m *mockActor) OnBlock(other body.Collidable) {}
func (m *mockActor) GetTouchable() body.Touchable { return nil }
func (m *mockActor) DrawCollisionBox(screen *ebiten.Image, position image.Rectangle) {}
func (m *mockActor) CollisionPosition() []image.Rectangle { return nil }
func (m *mockActor) CollisionShapes() []body.Collidable { return nil }
func (m *mockActor) IsObstructive() bool { return m.obstructive }
func (m *mockActor) SetIsObstructive(value bool) { m.obstructive = value }
func (m *mockActor) AddCollision(list ...body.Collidable) {}
func (m *mockActor) ClearCollisions() {}
func (m *mockActor) SetTouchable(t body.Touchable) {}
func (m *mockActor) ApplyValidPosition(d int, x bool, s body.BodiesSpace) (int, int, bool) {
	return m.pos.Min.X, m.pos.Min.Y, false
}

// mockBodiesSpace implements body.BodiesSpace interface for testing
type mockBodiesSpace struct {
	queryFunc       func(rect image.Rectangle) []body.Collidable
	tilemapProvider tilemaplayer.TilemapDimensionsProvider
}

func (m *mockBodiesSpace) Query(rect image.Rectangle) []body.Collidable {
	if m.queryFunc != nil {
		return m.queryFunc(rect)
	}
	return nil
}

func (m *mockBodiesSpace) GetTilemapDimensionsProvider() tilemaplayer.TilemapDimensionsProvider {
	return m.tilemapProvider
}

// Implement other required methods with no-op stubs
func (m *mockBodiesSpace) AddBody(_ body.Collidable)                              {}
func (m *mockBodiesSpace) RemoveBody(_ body.Collidable)                           {}
func (m *mockBodiesSpace) Clear()                                                        {}
func (m *mockBodiesSpace) QueueForRemoval(_ body.Collidable)                      {}
func (m *mockBodiesSpace) ProcessRemovals()                                              {}
func (m *mockBodiesSpace) Bodies() []body.Collidable                              { return nil }
func (m *mockBodiesSpace) ResolveCollisions(_ body.Collidable) (bool, bool) { return false, false }
func (m *mockBodiesSpace) Find(_ string) body.Collidable                          { return nil }
func (m *mockBodiesSpace) HasCollision(_ body.Collidable) bool                           { return false }
func (m *mockBodiesSpace) SetTilemapDimensionsProvider(_ tilemaplayer.TilemapDimensionsProvider) {}

// mockTilemapProvider implements space.TilemapDimensionsProvider for testing
type mockTilemapProvider struct {
	width  int
	height int
	bounds image.Rectangle
	hasBounds bool
}

func (m *mockTilemapProvider) GetTilemapWidth() int                  { return m.width }
func (m *mockTilemapProvider) GetTilemapHeight() int                 { return m.height }
func (m *mockTilemapProvider) GetCameraBounds() (image.Rectangle, bool) { 
	return m.bounds, m.hasBounds 
}

// mockCollidable implements body.Collidable for testing
type mockCollidable struct {
	rect        image.Rectangle
	obstructive bool
	id          string
	pos16X      int
	pos16Y      int
}

func (m *mockCollidable) ID() string                                       { return m.id }
func (m *mockCollidable) SetID(id string)                                  { m.id = id }
func (m *mockCollidable) Position() image.Rectangle                        { return m.rect }
func (m *mockCollidable) SetPosition(x, y int)                             { m.rect = image.Rect(x, y, x+m.rect.Dx(), y+m.rect.Dy()) }
func (m *mockCollidable) SetSize(w, h int) {
	m.rect.Max.X = m.rect.Min.X + w
	m.rect.Max.Y = m.rect.Min.Y + h
}
func (m *mockCollidable) Scale() float64 { return 1.0 }
func (m *mockCollidable) SetScale(s float64) {}
func (m *mockCollidable) SetPosition16(x16, y16 int)                       { m.pos16X = x16; m.pos16Y = y16 }
func (m *mockCollidable) GetPosition16() (int, int)                        { return m.pos16X, m.pos16Y }
func (m *mockCollidable) GetPositionMin() (int, int)                       { return m.rect.Min.X, m.rect.Min.Y }
func (m *mockCollidable) GetShape() body.Shape                             { return nil }
func (m *mockCollidable) IsObstructive() bool                              { return m.obstructive }
func (m *mockCollidable) SetIsObstructive(v bool)                          { m.obstructive = v }
func (m *mockCollidable) OnTouch(_ body.Collidable)                        {}
func (m *mockCollidable) OnBlock(_ body.Collidable)                        {}
func (m *mockCollidable) GetTouchable() body.Touchable                     { return nil }
func (m *mockCollidable) DrawCollisionBox(_ *ebiten.Image, _ image.Rectangle) {}
func (m *mockCollidable) CollisionPosition() []image.Rectangle             { return []image.Rectangle{m.rect} }
func (m *mockCollidable) CollisionShapes() []body.Collidable               { return []body.Collidable{m} }
func (m *mockCollidable) AddCollision(_ ...body.Collidable)                {}
func (m *mockCollidable) ClearCollisions()                                 {}
func (m *mockCollidable) SetTouchable(_ body.Touchable)                    {}
func (m *mockCollidable) ApplyValidPosition(_ int, _ bool, _ body.BodiesSpace) (int, int, bool) {
	return m.rect.Min.X, m.rect.Min.Y, false
}
func (m *mockCollidable) Owner() interface{}                               { return nil }
func (m *mockCollidable) SetOwner(_ interface{})                           {}
func (m *mockCollidable) LastOwner() interface{}                           { return nil }

// Helper functions for creating test mocks

// mockMovableCollidable implements body.MovableCollidable for testing obstacles
type mockMovableCollidable struct {
	*mockCollidable
	speed          int
	moveLeftForce  int
	moveRightForce int
}

func newMockMovableCollidable(x, y, w, h int) *mockMovableCollidable {
	return &mockMovableCollidable{
		mockCollidable: &mockCollidable{
			rect:        image.Rect(x, y, x+w, y+h),
			obstructive: true,
		},
	}
}

func (m *mockMovableCollidable) MoveX(d int)                                 {}
func (m *mockMovableCollidable) MoveY(d int)                                 {}
func (m *mockMovableCollidable) OnMoveLeft(d int)                            {}
func (m *mockMovableCollidable) OnMoveUpLeft(d int)                          {}
func (m *mockMovableCollidable) OnMoveDownLeft(d int)                        {}
func (m *mockMovableCollidable) OnMoveRight(d int)                           {}
func (m *mockMovableCollidable) OnMoveUpRight(d int)                         {}
func (m *mockMovableCollidable) OnMoveDownRight(d int)                       {}
func (m *mockMovableCollidable) OnMoveUp(d int)                              {}
func (m *mockMovableCollidable) OnMoveDown(d int)                            {}
func (m *mockMovableCollidable) Velocity() (int, int)                        { return 0, 0 }
func (m *mockMovableCollidable) SetVelocity(_, _ int)                        {}
func (m *mockMovableCollidable) Acceleration() (int, int)                    { return 0, 0 }
func (m *mockMovableCollidable) SetAcceleration(_, _ int)                    {}
func (m *mockMovableCollidable) SetSpeed(s int) error                        { m.speed = s; return nil }
func (m *mockMovableCollidable) SetMaxSpeed(_ int) error                     { return nil }
func (m *mockMovableCollidable) Speed() int                                  { return m.speed }
func (m *mockMovableCollidable) MaxSpeed() int                               { return 0 }
func (m *mockMovableCollidable) Immobile() bool                              { return false }
func (m *mockMovableCollidable) SetImmobile(_ bool)                          {}
func (m *mockMovableCollidable) SetFreeze(_ bool)                            {}
func (m *mockMovableCollidable) Freeze() bool                                { return false }
func (m *mockMovableCollidable) FaceDirection() animation.FacingDirectionEnum { return 0 }
func (m *mockMovableCollidable) SetFaceDirection(_ animation.FacingDirectionEnum) {}
func (m *mockMovableCollidable) IsIdle() bool                                { return true }
func (m *mockMovableCollidable) IsWalking() bool                             { return false }
func (m *mockMovableCollidable) IsFalling() bool                             { return false }
func (m *mockMovableCollidable) IsGoingUp() bool                             { return false }
func (m *mockMovableCollidable) CheckMovementDirectionX()                    {}
func (m *mockMovableCollidable) TryJump(_ int)                               {}
func (m *mockMovableCollidable) SetJumpForceMultiplier(_ float64)            {}
func (m *mockMovableCollidable) JumpForceMultiplier() float64                { return 1.0 }
func (m *mockMovableCollidable) SetHorizontalInertia(_ float64)              {}
func (m *mockMovableCollidable) HorizontalInertia() float64                  { return 1.0 }

// newMockSpaceWithGround creates a mock space that returns ground at specified positions
func newMockSpaceWithGround(groundPositions []image.Point) *mockBodiesSpace {
	return &mockBodiesSpace{
		queryFunc: func(rect image.Rectangle) []body.Collidable {
			for _, pos := range groundPositions {
				groundRect := image.Rect(pos.X, pos.Y, pos.X+1, pos.Y+1)
				if rect.Overlaps(groundRect) {
					return []body.Collidable{
						&mockCollidable{rect: groundRect, obstructive: true, id: "ground"},
					}
				}
			}
			return nil
		},
	}
}

// newMockSpaceWithObstacles creates a mock space with obstacles at specified rectangles
func newMockSpaceWithObstacles(obstacles []image.Rectangle) *mockBodiesSpace {
	return &mockBodiesSpace{
		queryFunc: func(rect image.Rectangle) []body.Collidable {
			var result []body.Collidable
			for _, obs := range obstacles {
				if rect.Overlaps(obs) {
					result = append(result, &mockCollidable{rect: obs, obstructive: true, id: "obstacle"})
				}
			}
			return result
		},
	}
}

// newMockActor creates a test actor at the specified position with given speed
func newMockActor(x, y, speed int) *mockActor {
	return &mockActor{
		pos:   image.Rect(x, y, x+10, y+10),
		speed: speed,
	}
}
