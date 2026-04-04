package movement

import (
	"image"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/hajimehoshi/ebiten/v2"
)

// mockInputSource covers the three input signals needed for drop-through disambiguation.
type mockInputSource struct {
	dropHeld bool
	duckHeld bool
	jumpHeld bool
}

func (m *mockInputSource) DropHeld() bool { return m.dropHeld }
func (m *mockInputSource) DuckHeld() bool { return m.duckHeld }
func (m *mockInputSource) JumpHeld() bool { return m.jumpHeld }

// mockOneWayPlatform is a minimal OneWayPlatform that tracks pass-through state
// with a real frame countdown so the expiry case can be tested without the
// full implementation.
type mockOneWayPlatform struct {
	id       string
	pos      image.Rectangle
	oneWay   bool
	counters map[string]int
}

func newMockOneWayPlatform(id string, pos image.Rectangle, oneWay bool) *mockOneWayPlatform {
	return &mockOneWayPlatform{id: id, pos: pos, oneWay: oneWay, counters: make(map[string]int)}
}

func (p *mockOneWayPlatform) IsOneWay() bool { return p.oneWay }

func (p *mockOneWayPlatform) SetPassThrough(actor body.Collidable, frames int) {
	p.counters[actor.ID()] = frames
}

func (p *mockOneWayPlatform) IsPassThrough(actor body.Collidable) bool {
	return p.counters[actor.ID()] > 0
}

func (p *mockOneWayPlatform) Update() {
	for k, v := range p.counters {
		v--
		if v <= 0 {
			delete(p.counters, k)
		} else {
			p.counters[k] = v
		}
	}
}

// body.Body stubs
func (p *mockOneWayPlatform) ID() string                 { return p.id }
func (p *mockOneWayPlatform) SetID(id string)            { p.id = id }
func (p *mockOneWayPlatform) Position() image.Rectangle  { return p.pos }
func (p *mockOneWayPlatform) SetPosition(x, y int)       {}
func (p *mockOneWayPlatform) SetPosition16(x16, y16 int) {}
func (p *mockOneWayPlatform) SetSize(w, h int)           {}
func (p *mockOneWayPlatform) Scale() float64             { return 1 }
func (p *mockOneWayPlatform) SetScale(float64)           {}
func (p *mockOneWayPlatform) GetPosition16() (int, int)  { return 0, 0 }
func (p *mockOneWayPlatform) GetPositionMin() (int, int) { return p.pos.Min.X, p.pos.Min.Y }
func (p *mockOneWayPlatform) GetShape() body.Shape       { return p }
func (p *mockOneWayPlatform) Width() int                 { return p.pos.Dx() }
func (p *mockOneWayPlatform) Height() int                { return p.pos.Dy() }
func (p *mockOneWayPlatform) Owner() interface{}         { return nil }
func (p *mockOneWayPlatform) SetOwner(interface{})       {}
func (p *mockOneWayPlatform) LastOwner() interface{}     { return nil }

// body.Collidable stubs
func (p *mockOneWayPlatform) GetTouchable() body.Touchable                        { return p }
func (p *mockOneWayPlatform) DrawCollisionBox(_ *ebiten.Image, _ image.Rectangle) {}
func (p *mockOneWayPlatform) CollisionPosition() []image.Rectangle                { return []image.Rectangle{p.pos} }
func (p *mockOneWayPlatform) CollisionShapes() []body.Collidable                  { return nil }
func (p *mockOneWayPlatform) IsObstructive() bool                                 { return true }
func (p *mockOneWayPlatform) SetIsObstructive(bool)                               {}
func (p *mockOneWayPlatform) AddCollision(...body.Collidable)                     {}
func (p *mockOneWayPlatform) ClearCollisions()                                    {}
func (p *mockOneWayPlatform) SetTouchable(body.Touchable)                         {}
func (p *mockOneWayPlatform) OnTouch(body.Collidable)                             {}
func (p *mockOneWayPlatform) OnBlock(body.Collidable)                             {}
func (p *mockOneWayPlatform) ApplyValidPosition(_ int, _ bool, _ body.BodiesSpace) (int, int, bool) {
	return 0, 0, false
}

func TestDropThrough(t *testing.T) {
	tests := []struct {
		name             string
		input            *mockInputSource
		platformIsOneWay bool
		vyBefore         int
		wantPassThrough  bool
		wantVyUnchanged  bool
	}{
		{
			name:             "down+jump on one-way triggers drop-through",
			input:            &mockInputSource{dropHeld: true},
			platformIsOneWay: true,
			vyBefore:         0,
			wantPassThrough:  true,
			wantVyUnchanged:  true,
		},
		{
			name:             "down alone on one-way does not trigger drop-through",
			input:            &mockInputSource{duckHeld: true, jumpHeld: false},
			platformIsOneWay: true,
			vyBefore:         0,
			wantPassThrough:  false,
			wantVyUnchanged:  true,
		},
		{
			name:             "jump alone on one-way does not trigger drop-through",
			input:            &mockInputSource{jumpHeld: true},
			platformIsOneWay: true,
			vyBefore:         0,
			wantPassThrough:  false,
			wantVyUnchanged:  false, // normal jump changes vy
		},
		{
			name:             "down+jump on solid platform does not trigger drop-through",
			input:            &mockInputSource{dropHeld: true},
			platformIsOneWay: false,
			vyBefore:         0,
			wantPassThrough:  false,
			wantVyUnchanged:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actor := newMockMovableCollidable()
			actor.SetID("actor")
			actor.SetVelocity(0, tt.vyBefore)

			platform := newMockOneWayPlatform("platform", image.Rect(0, 20, 100, 30), tt.platformIsOneWay)

			tryDropThrough(actor, platform, tt.input)

			got := platform.IsPassThrough(actor)
			if got != tt.wantPassThrough {
				t.Errorf("IsPassThrough = %v; want %v", got, tt.wantPassThrough)
			}

			if tt.wantVyUnchanged {
				_, vy := actor.Velocity()
				if vy != tt.vyBefore {
					t.Errorf("vy = %d; want %d (unchanged)", vy, tt.vyBefore)
				}
			}
		})
	}
}

func TestDropThrough_PassThroughExpires(t *testing.T) {
	actor := newMockMovableCollidable()
	actor.SetID("actor")

	platform := newMockOneWayPlatform("platform", image.Rect(0, 20, 100, 30), true)
	platform.SetPassThrough(actor, 2)

	if !platform.IsPassThrough(actor) {
		t.Fatal("expected IsPassThrough=true immediately after SetPassThrough")
	}

	platform.Update()
	if !platform.IsPassThrough(actor) {
		t.Error("expected IsPassThrough=true after 1 tick (countdown=1 remaining)")
	}

	platform.Update()
	if platform.IsPassThrough(actor) {
		t.Error("expected IsPassThrough=false after 2 ticks (expired)")
	}
}
