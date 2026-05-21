package space

import (
	"image"
	"testing"

	contractsbody "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/hajimehoshi/ebiten/v2"
)

type testCollidable struct {
	id          string
	rect        image.Rectangle
	obstructive bool
	touchCount  int
	blockCount  int
	owner       interface{}
}

func newTestCollidable(id string, rect image.Rectangle, obstructive bool) *testCollidable {
	return &testCollidable{
		id:          id,
		rect:        rect,
		obstructive: obstructive,
	}
}

func (c *testCollidable) ID() string {
	return c.id
}

func (c *testCollidable) SetID(id string) {
	c.id = id
}

func (c *testCollidable) Position() image.Rectangle {
	return c.rect
}

func (c *testCollidable) SetPosition(x, y int) {
	dx := c.rect.Dx()
	dy := c.rect.Dy()
	c.rect = image.Rect(x, y, x+dx, y+dy)
}

func (c *testCollidable) SetSize(w, h int) {
	c.rect.Max.X = c.rect.Min.X + w
	c.rect.Max.Y = c.rect.Min.Y + h
}

func (c *testCollidable) Scale() float64 {
	return 1.0
}

func (c *testCollidable) SetScale(s float64) {
}

func (c *testCollidable) SetPosition16(x16, y16 int) {
}

func (c *testCollidable) GetPosition16() (int, int) {
	return 0, 0
}

func (c *testCollidable) GetPositionMin() (int, int) {
	return c.rect.Min.X, c.rect.Min.Y
}

type rectShape struct {
	w int
	h int
}

func (r rectShape) Width() int {
	return r.w
}

func (r rectShape) Height() int {
	return r.h
}

func (c *testCollidable) GetShape() contractsbody.Shape {
	return rectShape{w: c.rect.Dx(), h: c.rect.Dy()}
}

func (c *testCollidable) OnTouch(other contractsbody.Collidable) {
	c.touchCount++
}

func (c *testCollidable) OnBlock(other contractsbody.Collidable) {
	c.blockCount++
}

func (c *testCollidable) GetTouchable() contractsbody.Touchable {
	return c
}

func (c *testCollidable) DrawCollisionBox(screen *ebiten.Image, position image.Rectangle) {
}

func (c *testCollidable) CollisionPosition() []image.Rectangle {
	return []image.Rectangle{c.rect}
}

func (c *testCollidable) CollisionShapes() []contractsbody.Collidable {
	return nil
}

func (c *testCollidable) IsObstructive() bool {
	return c.obstructive
}

func (c *testCollidable) SetIsObstructive(value bool) {
	c.obstructive = value
}

func (c *testCollidable) AddCollision(list ...contractsbody.Collidable) {
}

func (c *testCollidable) ClearCollisions() {
}

func (c *testCollidable) SetTouchable(t contractsbody.Touchable) {
}

func (c *testCollidable) ApplyValidPosition(distance16 int, isXAxis bool, space contractsbody.BodiesSpace) (int, int, bool) {
	return c.rect.Min.X, c.rect.Min.Y, false
}

func (c *testCollidable) Owner() interface{} {
	return c.owner
}

func (c *testCollidable) SetOwner(o interface{}) {
	c.owner = o
}

func (c *testCollidable) LastOwner() interface{} {
	return c.owner
}

func TestSpaceResolveCollisionsAndQueryIntegration(t *testing.T) {
	sp := NewSpace()

	rectA := image.Rect(0, 0, 10, 10)
	rectB := image.Rect(5, 5, 15, 15)

	a := newTestCollidable("a", rectA, false)
	b := newTestCollidable("b", rectB, true)

	sp.AddBody(a)
	sp.AddBody(b)

	touching, blocking := sp.ResolveCollisions(a)
	if !touching {
		t.Fatalf("expected touching collision between a and b")
	}
	if !blocking {
		t.Fatalf("expected blocking collision when other body is obstructive")
	}
	if a.touchCount == 0 || b.touchCount == 0 {
		t.Fatalf("expected OnTouch to be called on both bodies")
	}
	if a.blockCount == 0 || b.blockCount == 0 {
		t.Fatalf("expected OnBlock to be called on both bodies")
	}

	bodies := sp.Bodies()
	if len(bodies) != 2 {
		t.Fatalf("expected 2 bodies in space, got %d", len(bodies))
	}
	if bodies[0].ID() != "a" || bodies[1].ID() != "b" {
		t.Fatalf("expected bodies to be sorted by ID, got %s and %s", bodies[0].ID(), bodies[1].ID())
	}

	hits := sp.Query(image.Rect(0, 0, 20, 20))
	if len(hits) != 2 {
		t.Fatalf("expected Query to return both bodies, got %d", len(hits))
	}

	hits = sp.Query(image.Rect(0, 0, 3, 3))
	if len(hits) != 1 || hits[0].ID() != "a" {
		t.Fatalf("expected Query to return only body a in small area, got %d hits", len(hits))
	}

	c := newTestCollidable("c", image.Rect(100, 100, 110, 110), false)
	sp.AddBody(c)
	sp.QueueForRemoval(c)
	sp.ProcessRemovals()

	for _, body := range sp.Bodies() {
		if body.ID() == "c" {
			t.Fatalf("expected queued body to be removed from space")
		}
	}

	b.SetIsObstructive(false)

	touching, blocking = sp.ResolveCollisions(a)
	if !touching {
		t.Fatalf("expected touching collision between a and d")
	}
	if blocking {
		t.Fatalf("expected no blocking collision when other body is not obstructive")
	}
}

func (m *testCollidable) Altitude() int           { return 0 }
func (m *testCollidable) SetAltitude(alt int)     {}
func (m *testCollidable) Altitude16() int         { return 0 }
func (m *testCollidable) SetAltitude16(alt16 int) {}

// depthLaneCollidable embeds *testCollidable and adds the DepthLaneBody
// opt-in surface (GroundY, LaneHalfWidth). It is used to drive the
// depth-aware collision gating tests for story 062.
type depthLaneCollidable struct {
	*testCollidable
	groundY       int
	laneHalfWidth int
}

func (d *depthLaneCollidable) GroundY() int       { return d.groundY }
func (d *depthLaneCollidable) LaneHalfWidth() int { return d.laneHalfWidth }

// TestHasCollisionDepthLane verifies that HasCollision applies depth-lane
// gating when BOTH bodies opt-in via the DepthLaneBody interface. When only
// one (or neither) side opts in, the legacy 2D bbox-overlap behavior must be
// preserved.
//
// Algorithm under test:
//  1. bbox overlap first (short-circuit on miss)
//  2. if BOTH bodies implement DepthLaneBody, additionally require
//     abs(a.GroundY() - b.GroundY()) <= max(a.LaneHalfWidth(), b.LaneHalfWidth())
func TestHasCollisionDepthLane(t *testing.T) {
	type bodySpec struct {
		rect    image.Rectangle
		groundY int
		halfW   int
		optIn   bool
	}

	makeBody := func(id string, s bodySpec) contractsbody.Collidable {
		base := newTestCollidable(id, s.rect, false)
		if !s.optIn {
			return base
		}
		return &depthLaneCollidable{
			testCollidable: base,
			groundY:        s.groundY,
			laneHalfWidth:  s.halfW,
		}
	}

	cases := []struct {
		name string
		a    bodySpec
		b    bodySpec
		want bool
	}{
		{
			name: "T-062-1 same-lane bbox overlap returns true",
			a:    bodySpec{rect: image.Rect(0, 0, 10, 10), groundY: 100, halfW: 8, optIn: true},
			b:    bodySpec{rect: image.Rect(5, 5, 15, 15), groundY: 104, halfW: 8, optIn: true},
			want: true,
		},
		{
			name: "T-062-2 different-lane bbox overlap returns false",
			a:    bodySpec{rect: image.Rect(0, 0, 10, 10), groundY: 100, halfW: 8, optIn: true},
			b:    bodySpec{rect: image.Rect(5, 5, 15, 15), groundY: 120, halfW: 8, optIn: true},
			want: false,
		},
		{
			name: "T-062-3 legacy 2D bbox overlap when neither opts in",
			a:    bodySpec{rect: image.Rect(0, 0, 10, 10), optIn: false},
			b:    bodySpec{rect: image.Rect(5, 5, 15, 15), optIn: false},
			want: true,
		},
		{
			name: "T-062-4 same-lane but no bbox overlap returns false",
			a:    bodySpec{rect: image.Rect(0, 0, 10, 10), groundY: 100, halfW: 8, optIn: true},
			b:    bodySpec{rect: image.Rect(50, 50, 60, 60), groundY: 100, halfW: 8, optIn: true},
			want: false,
		},
		{
			name: "T-062-5 halfW=0 with equal GroundY returns true",
			a:    bodySpec{rect: image.Rect(0, 0, 10, 10), groundY: 100, halfW: 0, optIn: true},
			b:    bodySpec{rect: image.Rect(5, 5, 15, 15), groundY: 100, halfW: 0, optIn: true},
			want: true,
		},
		{
			name: "T-062-6 halfW=0 with GroundY off by one returns false",
			a:    bodySpec{rect: image.Rect(0, 0, 10, 10), groundY: 100, halfW: 0, optIn: true},
			b:    bodySpec{rect: image.Rect(5, 5, 15, 15), groundY: 101, halfW: 0, optIn: true},
			want: false,
		},
		{
			name: "T-062-7 mixed opt-in falls back to legacy 2D overlap",
			a:    bodySpec{rect: image.Rect(0, 0, 10, 10), groundY: 100, halfW: 8, optIn: true},
			b:    bodySpec{rect: image.Rect(5, 5, 15, 15), optIn: false},
			want: true,
		},
		{
			name: "T-062-8A asymmetric halfW uses max and passes",
			a:    bodySpec{rect: image.Rect(0, 0, 10, 10), groundY: 100, halfW: 2, optIn: true},
			b:    bodySpec{rect: image.Rect(5, 5, 15, 15), groundY: 108, halfW: 10, optIn: true},
			want: true,
		},
		{
			name: "T-062-8B asymmetric halfW uses max and fails",
			a:    bodySpec{rect: image.Rect(0, 0, 10, 10), groundY: 100, halfW: 2, optIn: true},
			b:    bodySpec{rect: image.Rect(5, 5, 15, 15), groundY: 111, halfW: 10, optIn: true},
			want: false,
		},
		{
			name: "T-062-9 equal GroundY with bbox overlap returns true",
			a:    bodySpec{rect: image.Rect(0, 0, 10, 10), groundY: 100, halfW: 8, optIn: true},
			b:    bodySpec{rect: image.Rect(5, 5, 15, 15), groundY: 100, halfW: 8, optIn: true},
			want: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			a := makeBody("a", tc.a)
			b := makeBody("b", tc.b)

			if got := HasCollision(a, b); got != tc.want {
				t.Fatalf("HasCollision(a, b) = %v, want %v", got, tc.want)
			}
		})
	}
}
