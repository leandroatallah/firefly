package gamehud

import (
	"image"
	"os"
	"path/filepath"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	gamesetup "github.com/boilerplate/ebiten-template/internal/game/app"
)

// localMockActor implements platformer.PlatformerActorEntity for testing
type localMockActor struct{}

func (m *localMockActor) ID() string                                           { return "mock" }
func (m *localMockActor) SetID(string)                                         {}
func (m *localMockActor) Position() image.Rectangle                            { return image.Rect(0, 0, 10, 10) }
func (m *localMockActor) SetPosition(int, int)                                 {}
func (m *localMockActor) SetPosition16(int, int)                               {}
func (m *localMockActor) GetPosition16() (int, int)                            { return 0, 0 }
func (m *localMockActor) GetPositionMin() (int, int)                           { return 0, 0 }
func (m *localMockActor) GetShape() interface{}                                { return nil }
func (m *localMockActor) Speed() int                                           { return 0 }
func (m *localMockActor) MaxSpeed() int                                        { return 0 }
func (m *localMockActor) SetSpeed(int) error                                   { return nil }
func (m *localMockActor) SetMaxSpeed(int) error                                { return nil }
func (m *localMockActor) MovementModel() interface{}                           { return nil }
func (m *localMockActor) SetMovementModel(interface{})                         {}
func (m *localMockActor) OnMoveLeft(int)                                       {}
func (m *localMockActor) OnMoveRight(int)                                      {}
func (m *localMockActor) SetMovementState(interface{}, interface{}, ...interface{}) {}
func (m *localMockActor) GetCharacter() interface{}                            { return nil }
func (m *localMockActor) Image() *ebiten.Image                                 { return nil }
func (m *localMockActor) ImageOptions() *ebiten.DrawImageOptions               { return nil }
func (m *localMockActor) UpdateImageOptions()                                  {}
func (m *localMockActor) BlockMovement()                                       {}
func (m *localMockActor) UnblockMovement()                                     {}
func (m *localMockActor) IsMovementBlocked() bool                              { return false }
func (m *localMockActor) State() interface{}                                   { return 0 }
func (m *localMockActor) SetState(interface{})                                 {}
func (m *localMockActor) SwitchMovementState(interface{})                      {}
func (m *localMockActor) MovementState() interface{}                           { return nil }
func (m *localMockActor) NewState(interface{}) (interface{}, error)            { return nil, nil }
func (m *localMockActor) Hurt(int)                                             {}
func (m *localMockActor) OnDie()                                               {}
func (m *localMockActor) OnJump()                                              {}
func (m *localMockActor) OnLand()                                              {}
func (m *localMockActor) OnFall()                                              {}
func (m *localMockActor) SetOnJump(func(interface{}))                          {}
func (m *localMockActor) SetOnFall(func(interface{}))                          {}
func (m *localMockActor) SetOnLand(func(interface{}))                          {}
func (m *localMockActor) SetAppContext(interface{})                            {}
func (m *localMockActor) AppContext() interface{}                              { return nil }
func (m *localMockActor) Owner() interface{}                                   { return nil }
func (m *localMockActor) SetOwner(interface{})                                 {}
func (m *localMockActor) LastOwner() interface{}                               { return nil }
func (m *localMockActor) Update(interface{}) error                             { return nil }
func (m *localMockActor) Health() int                                          { return 100 }
func (m *localMockActor) MaxHealth() int                                       { return 100 }
func (m *localMockActor) SetHealth(int)                                        {}
func (m *localMockActor) SetMaxHealth(int)                                     {}
func (m *localMockActor) LoseHealth(int)                                       {}
func (m *localMockActor) RestoreHealth(int)                                    {}
func (m *localMockActor) Invulnerable() bool                                   { return false }
func (m *localMockActor) SetInvulnerability(bool)                              {}
func (m *localMockActor) GetTouchable() interface{}                            { return nil }
func (m *localMockActor) OnTouch(interface{})                                  {}
func (m *localMockActor) OnBlock(interface{})                                  {}
func (m *localMockActor) DrawCollisionBox(*ebiten.Image, image.Rectangle)      {}
func (m *localMockActor) CollisionPosition() []image.Rectangle                 { return nil }
func (m *localMockActor) CollisionShapes() []interface{}                       { return nil }
func (m *localMockActor) IsObstructive() bool                                  { return false }
func (m *localMockActor) SetIsObstructive(bool)                                {}
func (m *localMockActor) AddCollision(...interface{})                          {}
func (m *localMockActor) ClearCollisions()                                     {}
func (m *localMockActor) SetTouchable(interface{})                             {}
func (m *localMockActor) ApplyValidPosition(int, bool, interface{}) (int, int, bool) { return 0, 0, false }
func (m *localMockActor) MoveX(int)                                            {}
func (m *localMockActor) MoveY(int)                                            {}
func (m *localMockActor) OnMoveUpLeft(int)                                     {}
func (m *localMockActor) OnMoveDownLeft(int)                                   {}
func (m *localMockActor) OnMoveUpRight(int)                                    {}
func (m *localMockActor) OnMoveDownRight(int)                                  {}
func (m *localMockActor) OnMoveUp(int)                                         {}
func (m *localMockActor) OnMoveDown(int)                                       {}
func (m *localMockActor) Velocity() (int, int)                                 { return 0, 0 }
func (m *localMockActor) SetVelocity(int, int)                                 {}
func (m *localMockActor) Acceleration() (int, int)                             { return 0, 0 }
func (m *localMockActor) SetAcceleration(int, int)                             {}
func (m *localMockActor) Immobile() bool                                       { return false }
func (m *localMockActor) SetImmobile(bool)                                     {}
func (m *localMockActor) SetFreeze(bool)                                       {}
func (m *localMockActor) Freeze() bool                                         { return false }
func (m *localMockActor) FaceDirection() interface{}                           { return 0 }
func (m *localMockActor) SetFaceDirection(interface{})                         {}
func (m *localMockActor) IsIdle() bool                                         { return true }
func (m *localMockActor) IsWalking() bool                                      { return false }
func (m *localMockActor) IsFalling() bool                                      { return false }
func (m *localMockActor) IsGoingUp() bool                                      { return false }
func (m *localMockActor) CheckMovementDirectionX()                             {}
func (m *localMockActor) TryJump(int)                                          {}
func (m *localMockActor) SetJumpForceMultiplier(float64)                       {}
func (m *localMockActor) JumpForceMultiplier() float64                         { return 1.0 }
func (m *localMockActor) SetHorizontalInertia(float64)                         {}
func (m *localMockActor) HorizontalInertia() float64                           { return 1.0 }
func (m *localMockActor) SetOnJumpPos(func(interface{}))                       {}
func (m *localMockActor) SetOnFallPos(func(interface{}))                       {}
func (m *localMockActor) SetOnLandPos(func(interface{}))                       {}

func getModuleRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			panic("could not find go.mod")
		}
		dir = parent
	}
}

func TestMain(m *testing.M) {
	// Change working directory to module root so assets can be found
	err := os.Chdir(getModuleRoot())
	if err != nil {
		panic(err)
	}

	// Initialize config
	cfg := gamesetup.NewConfig()
	config.Set(cfg)

	os.Exit(m.Run())
}

func TestNewStatusBar(t *testing.T) {
	// StatusBar expects platformer.PlatformerActorEntity
	// Let's see what that interface is.

	sb, err := NewStatusBar(nil, 100, nil, os.DirFS("."))
	if err != nil {
		t.Fatalf("failed to create StatusBar: %v", err)
	}

	if sb == nil {
		t.Fatal("NewStatusBar returned nil")
	}

	if sb.score != 100 {
		t.Errorf("expected score 100, got %d", sb.score)
	}
}

func TestStatusBar_Update(t *testing.T) {
	sb, err := NewStatusBar(nil, 100, nil, os.DirFS("."))
	if err != nil {
		t.Fatalf("failed to create StatusBar: %v", err)
	}

	if err := sb.Update(); err != nil {
		t.Errorf("Update returned error: %v", err)
	}
}

func TestStatusBar_Draw(t *testing.T) {
	// Note: Testing Draw with a real player requires full PlatformerActorEntity implementation
	// For now, test the nil player case which exercises the early return logic
	sb, err := NewStatusBar(nil, 100, nil, os.DirFS("."))
	if err != nil {
		t.Fatalf("failed to create StatusBar: %v", err)
	}

	screen := ebiten.NewImage(320, 240)

	// Test with no player (should return early)
	sb.player = nil
	sb.Draw(screen)
}
