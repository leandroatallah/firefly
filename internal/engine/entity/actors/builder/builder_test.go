package builder

import (
	"image"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/data/schemas"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/movement"
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
	"github.com/boilerplate/ebiten-template/internal/engine/render/sprites"
	"github.com/hajimehoshi/ebiten/v2"
)

// mockActorWithCollision implements actors.ActorEntity with collisionRectSetter
type mockActorWithCollision struct {
	id                      string
	speed                   int
	maxSpeed                int
	health                  int
	maxHealth               int
	movementMdl             physicsmovement.MovementModel
	character               *actors.Character
	addCollisionRectCalled  int
	refreshCollisionsCalled int
	lastState               actors.ActorStateEnum
	lastRect                body.Collidable

	// Body components
	x, y      int
	width     int
	height    int
	velocity  struct{ x, y int }
	accel     struct{ x, y int }
	immobile  bool
	frozen    bool
	owner     interface{}
	lastOwner interface{}
	touchable body.Touchable
}

func newMockActorWithCollision() *mockActorWithCollision {
	return &mockActorWithCollision{
		width: 16, height: 16,
		maxHealth: 100,
	}
}

// ActorEntity implementation
func (m *mockActorWithCollision) ID() string      { return m.id }
func (m *mockActorWithCollision) SetID(id string) { m.id = id }
func (m *mockActorWithCollision) Position() image.Rectangle {
	return image.Rect(m.x, m.y, m.x+m.width, m.y+m.height)
}
func (m *mockActorWithCollision) SetPosition(x, y int)       { m.x, m.y = x, y }
func (m *mockActorWithCollision) SetSize(w, h int)           { m.width, m.height = w, h }
func (m *mockActorWithCollision) Scale() float64             { return 1.0 }
func (m *mockActorWithCollision) SetScale(s float64)         {}
func (m *mockActorWithCollision) SetPosition16(x16, y16 int) { m.x, m.y = x16/16, y16/16 }
func (m *mockActorWithCollision) GetPosition16() (int, int)  { return m.x * 16, m.y * 16 }
func (m *mockActorWithCollision) GetPositionMin() (int, int) { return m.x, m.y }
func (m *mockActorWithCollision) GetShape() body.Shape {
	return bodyphysics.NewRect(m.x, m.y, m.width, m.height)
}
func (m *mockActorWithCollision) Speed() int                                   { return m.speed }
func (m *mockActorWithCollision) MaxSpeed() int                                { return m.maxSpeed }
func (m *mockActorWithCollision) SetSpeed(s int) error                         { m.speed = s; return nil }
func (m *mockActorWithCollision) SetMaxSpeed(s int) error                      { m.maxSpeed = s; return nil }
func (m *mockActorWithCollision) MovementModel() physicsmovement.MovementModel { return m.movementMdl }
func (m *mockActorWithCollision) SetMovementModel(model physicsmovement.MovementModel) {
	m.movementMdl = model
}
func (m *mockActorWithCollision) GetCharacter() *actors.Character        { return m.character }
func (m *mockActorWithCollision) SetCharacter(c *actors.Character)       { m.character = c }
func (m *mockActorWithCollision) Health() int                            { return m.health }
func (m *mockActorWithCollision) MaxHealth() int                         { return m.maxHealth }
func (m *mockActorWithCollision) SetHealth(h int)                        { m.health = h }
func (m *mockActorWithCollision) SetMaxHealth(h int)                     { m.maxHealth = h }
func (m *mockActorWithCollision) LoseHealth(d int)                       { m.health -= d }
func (m *mockActorWithCollision) RestoreHealth(h int)                    { m.health += h }
func (m *mockActorWithCollision) Update(space body.BodiesSpace) error    { return nil }
func (m *mockActorWithCollision) Image() *ebiten.Image                   { return nil }
func (m *mockActorWithCollision) ImageOptions() *ebiten.DrawImageOptions { return nil }
func (m *mockActorWithCollision) UpdateImageOptions()                    {}

// Stateful
func (m *mockActorWithCollision) State() actors.ActorStateEnum     { return actors.Idle }
func (m *mockActorWithCollision) SetState(state actors.ActorState) {}
func (m *mockActorWithCollision) SetMovementState(state movement.MovementStateEnum, target body.MovableCollidable, options ...movement.MovementStateOption) {
}
func (m *mockActorWithCollision) SwitchMovementState(state movement.MovementStateEnum) {}
func (m *mockActorWithCollision) MovementState() movement.MovementState                { return nil }
func (m *mockActorWithCollision) NewState(state actors.ActorStateEnum) (actors.ActorState, error) {
	return nil, nil
}

// Controllable
func (m *mockActorWithCollision) OnMoveLeft(force int)    {}
func (m *mockActorWithCollision) OnMoveRight(force int)   {}
func (m *mockActorWithCollision) BlockMovement()          {}
func (m *mockActorWithCollision) UnblockMovement()        {}
func (m *mockActorWithCollision) IsMovementBlocked() bool { return false }

// Damageable
func (m *mockActorWithCollision) Hurt(damage int) {}

// Ownable
func (m *mockActorWithCollision) Owner() interface{}         { return m.owner }
func (m *mockActorWithCollision) SetOwner(owner interface{}) { m.owner = owner }
func (m *mockActorWithCollision) LastOwner() interface{}     { return m.lastOwner }

// Movable
func (m *mockActorWithCollision) Velocity() (int, int)       { return m.velocity.x, m.velocity.y }
func (m *mockActorWithCollision) SetVelocity(vx, vy int)     { m.velocity.x, m.velocity.y = vx, vy }
func (m *mockActorWithCollision) Acceleration() (int, int)   { return m.accel.x, m.accel.y }
func (m *mockActorWithCollision) SetAcceleration(ax, ay int) { m.accel.x, m.accel.y = ax, ay }
func (m *mockActorWithCollision) Immobile() bool             { return m.immobile }
func (m *mockActorWithCollision) SetImmobile(i bool)         { m.immobile = i }
func (m *mockActorWithCollision) Freeze() bool               { return m.frozen }
func (m *mockActorWithCollision) SetFreeze(f bool)           { m.frozen = f }
func (m *mockActorWithCollision) MoveX(d int)                { m.x += d }
func (m *mockActorWithCollision) MoveY(d int)                { m.y += d }
func (m *mockActorWithCollision) OnMoveUpLeft(d int)         {}
func (m *mockActorWithCollision) OnMoveDownLeft(d int)       {}
func (m *mockActorWithCollision) OnMoveUpRight(d int)        {}
func (m *mockActorWithCollision) OnMoveDownRight(d int)      {}
func (m *mockActorWithCollision) OnMoveUp(d int)             {}
func (m *mockActorWithCollision) OnMoveDown(d int)           {}
func (m *mockActorWithCollision) ApplyValidPosition(d int, ax bool, sp body.BodiesSpace) (int, int, bool) {
	return m.x, m.y, false
}
func (m *mockActorWithCollision) CheckMovementDirectionX()                         {}
func (m *mockActorWithCollision) TryJump(f int)                                    {}
func (m *mockActorWithCollision) SetJumpForceMultiplier(mu float64)                {}
func (m *mockActorWithCollision) JumpForceMultiplier() float64                     { return 1.0 }
func (m *mockActorWithCollision) SetHorizontalInertia(i float64)                   {}
func (m *mockActorWithCollision) HorizontalInertia() float64                       { return 1.0 }
func (m *mockActorWithCollision) FaceDirection() animation.FacingDirectionEnum     { return 0 }
func (m *mockActorWithCollision) SetFaceDirection(v animation.FacingDirectionEnum) {}
func (m *mockActorWithCollision) IsIdle() bool                                     { return true }
func (m *mockActorWithCollision) IsWalking() bool                                  { return false }
func (m *mockActorWithCollision) IsFalling() bool                                  { return false }
func (m *mockActorWithCollision) IsGoingUp() bool                                  { return false }

// Touchable
func (m *mockActorWithCollision) SetTouchable(t body.Touchable) { m.touchable = t }
func (m *mockActorWithCollision) GetTouchable() body.Touchable  { return m.touchable }
func (m *mockActorWithCollision) OnTouch(other body.Collidable) {}
func (m *mockActorWithCollision) OnBlock(other body.Collidable) {}

// Collidable
func (m *mockActorWithCollision) DrawCollisionBox(s *ebiten.Image, p image.Rectangle) {}
func (m *mockActorWithCollision) CollisionPosition() []image.Rectangle {
	return []image.Rectangle{m.Position()}
}
func (m *mockActorWithCollision) CollisionShapes() []body.Collidable   { return nil }
func (m *mockActorWithCollision) IsObstructive() bool                  { return true }
func (m *mockActorWithCollision) SetIsObstructive(v bool)              {}
func (m *mockActorWithCollision) AddCollision(list ...body.Collidable) {}
func (m *mockActorWithCollision) ClearCollisions()                     {}

// Alive
func (m *mockActorWithCollision) Invulnerable() bool        { return false }
func (m *mockActorWithCollision) SetInvulnerability(v bool) {}

// collisionRectSetter
func (m *mockActorWithCollision) AddCollisionRect(state actors.ActorStateEnum, rect body.Collidable) {
	m.addCollisionRectCalled++
	m.lastState = state
	m.lastRect = rect
}
func (m *mockActorWithCollision) RefreshCollisions() { m.refreshCollisionsCalled++ }

// minimalActor is a minimal ActorEntity implementation that doesn't implement collisionRectSetter
type minimalActor struct {
	id string
}

func (m *minimalActor) ID() string                                           { return m.id }
func (m *minimalActor) SetID(id string)                                      { m.id = id }
func (m *minimalActor) Position() image.Rectangle                            { return image.Rect(0, 0, 16, 16) }
func (m *minimalActor) SetPosition(x, y int)                                 {}
func (m *minimalActor) SetSize(w, h int)                                     {}
func (m *minimalActor) Scale() float64                                       { return 1.0 }
func (m *minimalActor) SetScale(s float64)                                   {}
func (m *minimalActor) SetPosition16(x16, y16 int)                           {}
func (m *minimalActor) GetPosition16() (int, int)                            { return 0, 0 }
func (m *minimalActor) GetPositionMin() (int, int)                           { return 0, 0 }
func (m *minimalActor) GetShape() body.Shape                                 { return bodyphysics.NewRect(0, 0, 16, 16) }
func (m *minimalActor) Speed() int                                           { return 0 }
func (m *minimalActor) MaxSpeed() int                                        { return 0 }
func (m *minimalActor) SetSpeed(s int) error                                 { return nil }
func (m *minimalActor) SetMaxSpeed(s int) error                              { return nil }
func (m *minimalActor) MovementModel() physicsmovement.MovementModel         { return nil }
func (m *minimalActor) SetMovementModel(model physicsmovement.MovementModel) {}
func (m *minimalActor) GetCharacter() *actors.Character                      { return nil }
func (m *minimalActor) Health() int                                          { return 0 }
func (m *minimalActor) MaxHealth() int                                       { return 0 }
func (m *minimalActor) SetHealth(h int)                                      {}
func (m *minimalActor) SetMaxHealth(h int)                                   {}
func (m *minimalActor) LoseHealth(d int)                                     {}
func (m *minimalActor) RestoreHealth(h int)                                  {}
func (m *minimalActor) Update(space body.BodiesSpace) error                  { return nil }
func (m *minimalActor) Image() *ebiten.Image                                 { return nil }
func (m *minimalActor) ImageOptions() *ebiten.DrawImageOptions               { return nil }
func (m *minimalActor) UpdateImageOptions()                                  {}
func (m *minimalActor) State() actors.ActorStateEnum                         { return actors.Idle }
func (m *minimalActor) SetState(state actors.ActorState)                     {}
func (m *minimalActor) SetMovementState(state movement.MovementStateEnum, target body.MovableCollidable, options ...movement.MovementStateOption) {
}
func (m *minimalActor) SwitchMovementState(state movement.MovementStateEnum) {}
func (m *minimalActor) MovementState() movement.MovementState                { return nil }
func (m *minimalActor) NewState(state actors.ActorStateEnum) (actors.ActorState, error) {
	return nil, nil
}
func (m *minimalActor) OnMoveLeft(force int)       {}
func (m *minimalActor) OnMoveRight(force int)      {}
func (m *minimalActor) BlockMovement()             {}
func (m *minimalActor) UnblockMovement()           {}
func (m *minimalActor) IsMovementBlocked() bool    { return false }
func (m *minimalActor) Hurt(damage int)            {}
func (m *minimalActor) Owner() interface{}         { return nil }
func (m *minimalActor) SetOwner(owner interface{}) {}
func (m *minimalActor) LastOwner() interface{}     { return nil }
func (m *minimalActor) Velocity() (int, int)       { return 0, 0 }
func (m *minimalActor) SetVelocity(vx, vy int)     {}
func (m *minimalActor) Acceleration() (int, int)   { return 0, 0 }
func (m *minimalActor) SetAcceleration(ax, ay int) {}
func (m *minimalActor) Immobile() bool             { return false }
func (m *minimalActor) SetImmobile(i bool)         {}
func (m *minimalActor) Freeze() bool               { return false }
func (m *minimalActor) SetFreeze(f bool)           {}
func (m *minimalActor) MoveX(d int)                {}
func (m *minimalActor) MoveY(d int)                {}
func (m *minimalActor) OnMoveUpLeft(d int)         {}
func (m *minimalActor) OnMoveDownLeft(d int)       {}
func (m *minimalActor) OnMoveUpRight(d int)        {}
func (m *minimalActor) OnMoveDownRight(d int)      {}
func (m *minimalActor) OnMoveUp(d int)             {}
func (m *minimalActor) OnMoveDown(d int)           {}
func (m *minimalActor) ApplyValidPosition(d int, ax bool, sp body.BodiesSpace) (int, int, bool) {
	return 0, 0, false
}
func (m *minimalActor) CheckMovementDirectionX()                            {}
func (m *minimalActor) TryJump(f int)                                       {}
func (m *minimalActor) SetJumpForceMultiplier(mu float64)                   {}
func (m *minimalActor) JumpForceMultiplier() float64                        { return 1.0 }
func (m *minimalActor) SetHorizontalInertia(i float64)                      {}
func (m *minimalActor) HorizontalInertia() float64                          { return 1.0 }
func (m *minimalActor) FaceDirection() animation.FacingDirectionEnum        { return 0 }
func (m *minimalActor) SetFaceDirection(v animation.FacingDirectionEnum)    {}
func (m *minimalActor) IsIdle() bool                                        { return true }
func (m *minimalActor) IsWalking() bool                                     { return false }
func (m *minimalActor) IsFalling() bool                                     { return false }
func (m *minimalActor) IsGoingUp() bool                                     { return false }
func (m *minimalActor) SetTouchable(t body.Touchable)                       {}
func (m *minimalActor) GetTouchable() body.Touchable                        { return nil }
func (m *minimalActor) OnTouch(other body.Collidable)                       {}
func (m *minimalActor) OnBlock(other body.Collidable)                       {}
func (m *minimalActor) DrawCollisionBox(s *ebiten.Image, p image.Rectangle) {}
func (m *minimalActor) CollisionPosition() []image.Rectangle                { return nil }
func (m *minimalActor) CollisionShapes() []body.Collidable                  { return nil }
func (m *minimalActor) IsObstructive() bool                                 { return false }
func (m *minimalActor) SetIsObstructive(v bool)                             {}
func (m *minimalActor) AddCollision(list ...body.Collidable)                {}
func (m *minimalActor) ClearCollisions()                                    {}
func (m *minimalActor) Invulnerable() bool                                  { return false }
func (m *minimalActor) SetInvulnerability(v bool)                           {}

func TestBuildStateMap(t *testing.T) {
	t.Run("valid states", func(t *testing.T) {
		data := schemas.SpriteData{
			Assets: map[string]schemas.AssetData{
				"idle": {Path: "idle.png"},
				"walk": {Path: "walk.png"},
				"jump": {Path: "jump.png"},
			},
		}

		stateMap, err := BuildStateMap(data)
		if err != nil {
			t.Fatalf("BuildStateMap returned error: %v", err)
		}

		if len(stateMap) != 3 {
			t.Errorf("Expected 3 states, got %d", len(stateMap))
		}

		if _, ok := stateMap["idle"]; !ok {
			t.Error("Expected 'idle' state in map")
		}
		if _, ok := stateMap["walk"]; !ok {
			t.Error("Expected 'walk' state in map")
		}
		if _, ok := stateMap["jump"]; !ok {
			t.Error("Expected 'jump' state in map")
		}
	})

	t.Run("empty assets", func(t *testing.T) {
		data := schemas.SpriteData{
			Assets: map[string]schemas.AssetData{},
		}

		stateMap, err := BuildStateMap(data)
		if err != nil {
			t.Fatalf("BuildStateMap returned error: %v", err)
		}

		if len(stateMap) != 0 {
			t.Errorf("Expected 0 states, got %d", len(stateMap))
		}
	})

	t.Run("unregistered state", func(t *testing.T) {
		data := schemas.SpriteData{
			Assets: map[string]schemas.AssetData{
				"unknown_state": {Path: "unknown.png"},
			},
		}

		_, err := BuildStateMap(data)
		if err == nil {
			t.Error("Expected error for unregistered state, got nil")
		}
	})
}

func TestBodyRectFromSpriteData(t *testing.T) {
	data := schemas.SpriteData{
		BodyRect: schemas.ShapeRect{X: 10, Y: 20, Width: 32, Height: 48},
	}

	rect := BodyRectFromSpriteData(data)
	if rect == nil {
		t.Fatal("BodyRectFromSpriteData returned nil")
	}

	if rect.Width() != 32 {
		t.Errorf("Expected width 32, got %d", rect.Width())
	}
	if rect.Height() != 48 {
		t.Errorf("Expected height 48, got %d", rect.Height())
	}
}

func TestSetCharacterStats(t *testing.T) {
	actor := newMockActorWithCollision()
	statData := actors.StatData{
		Health:   100,
		Speed:    5,
		MaxSpeed: 10,
	}

	err := SetCharacterStats(actor, statData)
	if err != nil {
		t.Fatalf("SetCharacterStats returned error: %v", err)
	}

	if actor.MaxHealth() != 100 {
		t.Errorf("Expected max health 100, got %d", actor.MaxHealth())
	}
	if actor.Speed() != 5 {
		t.Errorf("Expected speed 5, got %d", actor.Speed())
	}
	if actor.MaxSpeed() != 10 {
		t.Errorf("Expected max speed 10, got %d", actor.MaxSpeed())
	}
}

func TestSetCharacterBodies(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		actor := newMockActorWithCollision()
		character := actors.NewCharacter(sprites.SpriteMap{}, bodyphysics.NewRect(0, 0, 16, 16))
		actor.SetCharacter(character)

		data := schemas.SpriteData{
			BodyRect: schemas.ShapeRect{X: 0, Y: 0, Width: 16, Height: 16},
			Assets: map[string]schemas.AssetData{
				"idle": {Path: "idle.png", CollisionRects: []schemas.ShapeRect{{X: 0, Y: 0, Width: 10, Height: 10}}},
			},
		}

		stateMap := map[string]animation.SpriteState{
			"idle": actors.Idle,
		}

		err := SetCharacterBodies(actor, data, stateMap, "test_actor")
		if err != nil {
			t.Fatalf("SetCharacterBodies returned error: %v", err)
		}

		if actor.ID() != "test_actor" {
			t.Errorf("Expected ID 'test_actor', got '%s'", actor.ID())
		}

		if actor.addCollisionRectCalled == 0 {
			t.Error("Expected AddCollisionRect to be called")
		}

		if actor.refreshCollisionsCalled == 0 {
			t.Error("Expected RefreshCollisions to be called")
		}
	})

	t.Run("missing collisionRectSetter", func(t *testing.T) {
		actor := &minimalActor{id: "test"}
		data := schemas.SpriteData{}
		stateMap := map[string]animation.SpriteState{}

		err := SetCharacterBodies(actor, data, stateMap, "test")
		if err == nil {
			t.Error("Expected error for missing collisionRectSetter, got nil")
		}
	})
}

func TestConfigureCharacter(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		actor := newMockActorWithCollision()
		character := actors.NewCharacter(sprites.SpriteMap{}, bodyphysics.NewRect(0, 0, 16, 16))
		actor.SetCharacter(character)

		spriteData := schemas.SpriteData{
			BodyRect: schemas.ShapeRect{X: 0, Y: 0, Width: 16, Height: 16},
			Assets: map[string]schemas.AssetData{
				"idle": {Path: "idle.png", CollisionRects: []schemas.ShapeRect{{X: 0, Y: 0, Width: 10, Height: 10}}},
			},
		}

		statData := actors.StatData{
			Health:   100,
			Speed:    5,
			MaxSpeed: 10,
		}

		stateMap := map[string]animation.SpriteState{
			"idle": actors.Idle,
		}

		err := ConfigureCharacter(actor, spriteData, statData, stateMap, "test_actor")
		if err != nil {
			t.Fatalf("ConfigureCharacter returned error: %v", err)
		}

		if actor.ID() != "test_actor" {
			t.Errorf("Expected ID 'test_actor', got '%s'", actor.ID())
		}
		if actor.MaxHealth() != 100 {
			t.Errorf("Expected max health 100, got %d", actor.MaxHealth())
		}
	})

	t.Run("SetCharacterBodies fails", func(t *testing.T) {
		actor := &minimalActor{id: "test"}
		spriteData := schemas.SpriteData{}
		statData := actors.StatData{}
		stateMap := map[string]animation.SpriteState{}

		err := ConfigureCharacter(actor, spriteData, statData, stateMap, "test")
		if err == nil {
			t.Error("Expected error from ConfigureCharacter, got nil")
		}
	})
}

func TestApplyPlatformerPhysics(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		actor := newMockActorWithCollision()
		character := actors.NewCharacter(sprites.SpriteMap{}, bodyphysics.NewRect(0, 0, 16, 16))
		actor.SetCharacter(character)

		// Create a mock blocker that implements PlayerMovementBlocker
		blocker := &mockBlocker{}

		err := ApplyPlatformerPhysics(actor, blocker)
		if err != nil {
			t.Fatalf("ApplyPlatformerPhysics returned error: %v", err)
		}

		if actor.MovementModel() == nil {
			t.Error("Expected movement model to be set")
		}

		// Note: SetTouchable is called on the character, but GetTouchable()
		// returns the CollidableBody's Touchable which is separate from Character's Touchable field
		// The important thing is that no panic occurred
	})
}

// mockBlocker implements physicsmovement.PlayerMovementBlocker
type mockBlocker struct{}

func (m *mockBlocker) IsMovementBlocked() bool { return false }
