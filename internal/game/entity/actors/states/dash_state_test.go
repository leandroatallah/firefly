package gamestates_test

import (
	"image"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	contractsbody "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/tilemaplayer"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	gamestates "github.com/boilerplate/ebiten-template/internal/game/entity/actors/states"
)

// mockBody implements body.MovableCollidable for DashState tests.
type mockBody struct {
	pos           image.Rectangle
	vx, vy        int
	freeze        bool
	faceDir       animation.FacingDirectionEnum
	grounded      bool
	wallBlocked   bool
}

func (m *mockBody) ID() string                                                    { return "mock" }
func (m *mockBody) SetID(_ string)                                                {}
func (m *mockBody) Position() image.Rectangle                                     { return m.pos }
func (m *mockBody) SetPosition(x, y int)                                          { m.pos = image.Rect(x, y, x+m.pos.Dx(), y+m.pos.Dy()) }
func (m *mockBody) SetSize(w, h int)                                              { m.pos.Max.X = m.pos.Min.X + w; m.pos.Max.Y = m.pos.Min.Y + h }
func (m *mockBody) Scale() float64                                                { return 1 }
func (m *mockBody) SetScale(_ float64)                                            {}
func (m *mockBody) SetPosition16(_, _ int)                                        {}
func (m *mockBody) GetPosition16() (int, int)                                     { return 0, 0 }
func (m *mockBody) GetPositionMin() (int, int)                                    { return m.pos.Min.X, m.pos.Min.Y }
func (m *mockBody) GetShape() contractsbody.Shape                                 { return nil }
func (m *mockBody) Owner() interface{}                                            { return nil }
func (m *mockBody) SetOwner(_ interface{})                                        {}
func (m *mockBody) LastOwner() interface{}                                        { return nil }
func (m *mockBody) MoveX(_ int)                                                   {}
func (m *mockBody) MoveY(_ int)                                                   {}
func (m *mockBody) OnMoveLeft(_ int)                                              {}
func (m *mockBody) OnMoveUpLeft(_ int)                                            {}
func (m *mockBody) OnMoveDownLeft(_ int)                                          {}
func (m *mockBody) OnMoveRight(_ int)                                             {}
func (m *mockBody) OnMoveUpRight(_ int)                                           {}
func (m *mockBody) OnMoveDownRight(_ int)                                         {}
func (m *mockBody) OnMoveUp(_ int)                                                {}
func (m *mockBody) OnMoveDown(_ int)                                              {}
func (m *mockBody) Velocity() (int, int)                                          { return m.vx, m.vy }
func (m *mockBody) SetVelocity(vx, vy int)                                        { m.vx = vx; m.vy = vy }
func (m *mockBody) Acceleration() (int, int)                                      { return 0, 0 }
func (m *mockBody) SetAcceleration(_, _ int)                                      {}
func (m *mockBody) SetSpeed(_ int) error                                          { return nil }
func (m *mockBody) SetMaxSpeed(_ int) error                                       { return nil }
func (m *mockBody) Speed() int                                                    { return 0 }
func (m *mockBody) MaxSpeed() int                                                 { return 0 }
func (m *mockBody) Immobile() bool                                                { return false }
func (m *mockBody) SetImmobile(_ bool)                                            {}
func (m *mockBody) SetFreeze(f bool)                                              { m.freeze = f }
func (m *mockBody) Freeze() bool                                                  { return m.freeze }
func (m *mockBody) FaceDirection() animation.FacingDirectionEnum                  { return m.faceDir }
func (m *mockBody) SetFaceDirection(v animation.FacingDirectionEnum)              { m.faceDir = v }
func (m *mockBody) IsIdle() bool                                                  { return false }
func (m *mockBody) IsWalking() bool                                               { return false }
func (m *mockBody) IsFalling() bool                                               { return !m.grounded }
func (m *mockBody) IsGoingUp() bool                                               { return false }
func (m *mockBody) CheckMovementDirectionX()                                      {}
func (m *mockBody) TryJump(_ int)                                                 {}
func (m *mockBody) SetJumpForceMultiplier(_ float64)                              {}
func (m *mockBody) JumpForceMultiplier() float64                                  { return 1 }
func (m *mockBody) SetHorizontalInertia(_ float64)                                {}
func (m *mockBody) HorizontalInertia() float64                                    { return 0 }
func (m *mockBody) OnTouch(_ contractsbody.Collidable)                            {}
func (m *mockBody) OnBlock(_ contractsbody.Collidable)                            {}
func (m *mockBody) GetTouchable() contractsbody.Touchable                         { return nil }
func (m *mockBody) DrawCollisionBox(screen *ebiten.Image, _ image.Rectangle)             {}
func (m *mockBody) CollisionPosition() []image.Rectangle                          { return nil }
func (m *mockBody) CollisionShapes() []contractsbody.Collidable                   { return nil }
func (m *mockBody) IsObstructive() bool                                           { return false }
func (m *mockBody) SetIsObstructive(_ bool)                                       {}
func (m *mockBody) AddCollision(_ ...contractsbody.Collidable)                    {}
func (m *mockBody) ClearCollisions()                                              {}
func (m *mockBody) SetTouchable(_ contractsbody.Touchable)                        {}
func (m *mockBody) ApplyValidPosition(_ int, _ bool, _ contractsbody.BodiesSpace) (int, int, bool) {
	return m.pos.Min.X, m.pos.Min.Y, m.wallBlocked
}

// mockSpace implements body.BodiesSpace for DashState tests.
type mockSpace struct {
	queryResult []contractsbody.Collidable
}

func (m *mockSpace) Query(_ image.Rectangle) []contractsbody.Collidable           { return m.queryResult }
func (m *mockSpace) AddBody(_ contractsbody.Collidable)                           {}
func (m *mockSpace) Bodies() []contractsbody.Collidable                           { return nil }
func (m *mockSpace) RemoveBody(_ contractsbody.Collidable)                        {}
func (m *mockSpace) QueueForRemoval(_ contractsbody.Collidable)                   {}
func (m *mockSpace) ProcessRemovals()                                             {}
func (m *mockSpace) Clear()                                                       {}
func (m *mockSpace) ResolveCollisions(_ contractsbody.Collidable) (bool, bool)    { return false, false }
func (m *mockSpace) SetTilemapDimensionsProvider(_ tilemaplayer.TilemapDimensionsProvider) {}
func (m *mockSpace) GetTilemapDimensionsProvider() tilemaplayer.TilemapDimensionsProvider  { return nil }
func (m *mockSpace) Find(_ string) contractsbody.Collidable                       { return nil }

func defaultCfg() gamestates.DashConfig {
	return gamestates.DashConfig{
		Speed:          160,
		DurationFrames: 18,
		BlockDistance:  32,
		Cooldown:       30,
		DuckHeight:     12,
	}
}

func newGroundedBody() *mockBody {
	return &mockBody{
		pos:      image.Rect(0, 0, 16, 24),
		grounded: true,
	}
}

func TestDashStateUpdate(t *testing.T) {
	tests := []struct {
		name          string
		setupBody     func() *mockBody
		setupSpace    func() *mockSpace
		ticksBefore   int // Update() calls before the assertion
		wantState     actors.ActorStateEnum
	}{
		{
			name: "tween in progress returns StateDashing",
			setupBody:  newGroundedBody,
			setupSpace: func() *mockSpace { return &mockSpace{} },
			ticksBefore: defaultCfg().DurationFrames / 2,
			wantState:  gamestates.StateDashing,
		},
		{
			name: "tween complete and grounded returns StateIdle",
			setupBody:  newGroundedBody,
			setupSpace: func() *mockSpace { return &mockSpace{} },
			ticksBefore: defaultCfg().DurationFrames,
			wantState:  actors.Idle,
		},
		{
			name: "tween complete and airborne returns StateFalling",
			setupBody: func() *mockBody {
				return &mockBody{pos: image.Rect(0, 0, 16, 24), grounded: false}
			},
			setupSpace:  func() *mockSpace { return &mockSpace{} },
			ticksBefore: defaultCfg().DurationFrames,
			wantState:   actors.Falling,
		},
		{
			name:      "wall collision mid-dash returns Falling or Idle",
			setupBody: func() *mockBody {
				b := newGroundedBody()
				b.wallBlocked = true
				return b
			},
			setupSpace:  func() *mockSpace { return &mockSpace{} },
			ticksBefore: 1,
			wantState:   actors.Idle, // grounded + wall → Idle
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := tc.setupBody()
			sp := tc.setupSpace()
			cfg := defaultCfg()

			ds := gamestates.NewDashState(b, sp, cfg)
			ds.OnStart(0)

			var got actors.ActorStateEnum
			for i := 0; i < tc.ticksBefore; i++ {
				got = ds.Update()
			}

			if got != tc.wantState {
				t.Errorf("got state %v, want %v", got, tc.wantState)
			}
		})
	}

	t.Run("second OnStart is no-op while dashing", func(t *testing.T) {
		b := newGroundedBody()
		sp := &mockSpace{}
		cfg := defaultCfg()

		ds := gamestates.NewDashState(b, sp, cfg)
		ds.OnStart(0)

		// Advance a few frames so tween is in progress
		for i := 0; i < 3; i++ {
			ds.Update()
		}

		// Second OnStart should not restart the tween
		ds.OnStart(0)

		got := ds.Update()
		if got != gamestates.StateDashing {
			t.Errorf("expected StateDashing after ignored second OnStart, got %v", got)
		}
	})
}
