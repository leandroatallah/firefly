package draworder_test

import (
	"image"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/render/draworder"
	"github.com/hajimehoshi/ebiten/v2"
)

// fakeCollidable is a minimal in-test stand-in for body.Collidable.
// Only ID() and GetPosition16() carry meaningful state for the sort tests;
// the rest are no-op stubs sufficient to satisfy the interface.
type fakeCollidable struct {
	id   string
	x16  int
	y16  int
	w, h int
}

func newFakeCollidable(id string, y16 int) *fakeCollidable {
	return &fakeCollidable{id: id, x16: 0, y16: y16, w: 16, h: 16}
}

// body.Body methods.
func (f *fakeCollidable) ID() string      { return f.id }
func (f *fakeCollidable) SetID(id string) { f.id = id }
func (f *fakeCollidable) Position() image.Rectangle {
	return image.Rect(f.x16/16, f.y16/16, f.x16/16+f.w, f.y16/16+f.h)
}
func (f *fakeCollidable) SetPosition(x, y int)       { f.x16, f.y16 = x*16, y*16 }
func (f *fakeCollidable) SetPosition16(x16, y16 int) { f.x16, f.y16 = x16, y16 }
func (f *fakeCollidable) SetSize(w, h int)           { f.w, f.h = w, h }
func (f *fakeCollidable) Scale() float64             { return 1 }
func (f *fakeCollidable) SetScale(float64)           {}
func (f *fakeCollidable) GetPosition16() (int, int)  { return f.x16, f.y16 }
func (f *fakeCollidable) GetPositionMin() (int, int) { return f.x16 / 16, f.y16 / 16 }
func (f *fakeCollidable) GetShape() body.Shape       { return f }
func (f *fakeCollidable) Width() int                 { return f.w }
func (f *fakeCollidable) Height() int                { return f.h }
func (f *fakeCollidable) Owner() interface{}         { return nil }
func (f *fakeCollidable) SetOwner(interface{})       {}
func (f *fakeCollidable) LastOwner() interface{}     { return nil }

// Altitude axis (Story 053).
func (f *fakeCollidable) Altitude() int     { return 0 }
func (f *fakeCollidable) SetAltitude(int)   {}
func (f *fakeCollidable) Altitude16() int   { return 0 }
func (f *fakeCollidable) SetAltitude16(int) {}

// body.Touchable.
func (f *fakeCollidable) OnTouch(body.Collidable) {}
func (f *fakeCollidable) OnBlock(body.Collidable) {}

// body.Collidable.
func (f *fakeCollidable) GetTouchable() body.Touchable                        { return f }
func (f *fakeCollidable) DrawCollisionBox(_ *ebiten.Image, _ image.Rectangle) {}
func (f *fakeCollidable) CollisionPosition() []image.Rectangle                { return nil }
func (f *fakeCollidable) CollisionShapes() []body.Collidable                  { return nil }
func (f *fakeCollidable) IsObstructive() bool                                 { return false }
func (f *fakeCollidable) SetIsObstructive(bool)                               {}
func (f *fakeCollidable) AddCollision(...body.Collidable)                     {}
func (f *fakeCollidable) ClearCollisions()                                    {}
func (f *fakeCollidable) SetTouchable(body.Touchable)                         {}
func (f *fakeCollidable) ApplyValidPosition(int, bool, body.BodiesSpace) (int, int, bool) {
	return f.x16 / 16, f.y16 / 16, false
}

// Compile-time assertion: fakeCollidable satisfies body.Collidable.
var _ body.Collidable = (*fakeCollidable)(nil)

// withAltitude sets the altitude on a fakeCollidable wrapper variant.
type fakeCollidableAlt struct {
	*fakeCollidable
	alt16 int
}

func (f *fakeCollidableAlt) Altitude() int       { return f.alt16 / 16 }
func (f *fakeCollidableAlt) Altitude16() int     { return f.alt16 }
func (f *fakeCollidableAlt) SetAltitude(a int)   { f.alt16 = a * 16 }
func (f *fakeCollidableAlt) SetAltitude16(a int) { f.alt16 = a }

func newFakeWithAltitude(id string, y16, alt16 int) *fakeCollidableAlt {
	return &fakeCollidableAlt{
		fakeCollidable: newFakeCollidable(id, y16),
		alt16:          alt16,
	}
}

func ids(cs []body.Collidable) []string {
	out := make([]string, len(cs))
	for i, c := range cs {
		out[i] = c.ID()
	}
	return out
}

func equalIDs(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestSortByGroundY_AscendingOrder(t *testing.T) {
	in := []body.Collidable{
		newFakeCollidable("a", 300),
		newFakeCollidable("b", 100),
		newFakeCollidable("c", 200),
	}

	out := draworder.SortByGroundY(in)

	got := ids(out)
	want := []string{"b", "c", "a"}
	if !equalIDs(got, want) {
		t.Fatalf("SortByGroundY ascending order: got %v; want %v", got, want)
	}
}

func TestSortByGroundY_StableForEqualY(t *testing.T) {
	in := []body.Collidable{
		newFakeCollidable("idA", 100),
		newFakeCollidable("idB", 100),
		newFakeCollidable("idC", 50),
	}

	out := draworder.SortByGroundY(in)

	got := ids(out)
	want := []string{"idC", "idA", "idB"}
	if !equalIDs(got, want) {
		t.Fatalf("SortByGroundY stable order: got %v; want %v", got, want)
	}
}

func TestSortByGroundY_AltitudeIgnored(t *testing.T) {
	// Same ground y16=200; different altitudes. Input order must be preserved.
	in := []body.Collidable{
		newFakeWithAltitude("first", 200, 0),
		newFakeWithAltitude("second", 200, 100*16),
	}

	out := draworder.SortByGroundY(in)

	got := ids(out)
	want := []string{"first", "second"}
	if !equalIDs(got, want) {
		t.Fatalf("SortByGroundY altitude must not affect order: got %v; want %v", got, want)
	}
}

func TestSortByGroundY_DoesNotMutateInput(t *testing.T) {
	a := newFakeCollidable("a", 300)
	b := newFakeCollidable("b", 100)
	c := newFakeCollidable("c", 200)
	in := []body.Collidable{a, b, c}

	snapshot := make([]body.Collidable, len(in))
	copy(snapshot, in)

	_ = draworder.SortByGroundY(in)

	if !equalIDs(ids(in), ids(snapshot)) {
		t.Fatalf("SortByGroundY mutated input: got %v; want %v", ids(in), ids(snapshot))
	}
	for i := range in {
		if in[i] != snapshot[i] {
			t.Fatalf("SortByGroundY swapped element at index %d", i)
		}
	}
}

func TestSortByGroundY_EmptyAndSingle(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		out := draworder.SortByGroundY(nil)
		if len(out) != 0 {
			t.Fatalf("expected empty slice; got len=%d", len(out))
		}
	})

	t.Run("single", func(t *testing.T) {
		only := newFakeCollidable("only", 42)
		in := []body.Collidable{only}

		out := draworder.SortByGroundY(in)
		if len(out) != 1 || out[0].ID() != "only" {
			t.Fatalf("expected single-element identity; got %v", ids(out))
		}
	})
}

// T-D1: SortByGroundYAltitude sorts by effective ground depth (y16 - altitude16).
// This Red-Phase test will fail to compile until draworder.SortByGroundYAltitude
// is introduced — that is the missing-behaviour signal.
func TestSortByGroundYAltitude(t *testing.T) {
	cases := []struct {
		name string
		in   []body.Collidable
		want []string
	}{
		{
			name: "all non-altitudable falls back to SortByGroundY",
			in: []body.Collidable{
				newFakeCollidable("a", 300),
				newFakeCollidable("b", 100),
				newFakeCollidable("c", 200),
			},
			want: []string{"b", "c", "a"},
		},
		{
			name: "mixed altitudable and non-altitudable: non-altitudable treated as alt=0",
			in: []body.Collidable{
				// Effective depth = y16 - alt16.
				newFakeCollidable("ground100", 100),         // 100 - 0 = 100
				newFakeWithAltitude("airHigh", 5*16, 8*16),  // 80 - 128 = -48
				newFakeWithAltitude("groundlowAlt0", 50, 0), // 50 - 0 = 50
				newFakeWithAltitude("airC", 20*16, 5*16),    // 320 - 80 = 240
			},
			want: []string{"airHigh", "groundlowAlt0", "ground100", "airC"},
		},
		{
			name: "stable order for equal effective depth",
			in: []body.Collidable{
				// All produce effective depth = 100.
				newFakeCollidable("first", 100),
				newFakeWithAltitude("second", 100, 0),
				newFakeWithAltitude("third", 116, 16), // 116 - 16 = 100
			},
			want: []string{"first", "second", "third"},
		},
		{
			name: "spec example — B before A before C",
			// A: y16=10, alt=0    -> 10
			// B: y16=5,  alt=8    -> -3
			// C: y16=20, alt=5    -> 15
			in: []body.Collidable{
				newFakeWithAltitude("A", 10, 0),
				newFakeWithAltitude("B", 5, 8),
				newFakeWithAltitude("C", 20, 5),
			},
			want: []string{"B", "A", "C"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := draworder.SortByGroundYAltitude(tc.in)
			got := ids(out)
			if !equalIDs(got, tc.want) {
				t.Fatalf("SortByGroundYAltitude: got %v; want %v", got, tc.want)
			}
		})
	}
}
