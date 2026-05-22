// Package shadow_test contains Red-Phase TDD tests for story
// 063-shadow-component. These tests exercise the observable behaviour
// described in SPEC.md §8 (T-S1..T-S7) and MUST fail against the
// Red-Phase skeleton in shadow.go.
package shadow_test

import (
	"image"
	"image/color"
	"math"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/render/camera"
	"github.com/boilerplate/ebiten-template/internal/kit/render/shadow"
	"github.com/hajimehoshi/ebiten/v2"
)

// --- mock AltitudeBody ------------------------------------------------------

// mockAltBody is a minimal AltitudeBody for Compute/Draw tests.
type mockAltBody struct {
	x, y, w, h int
	alt        int
}

func (m *mockAltBody) GetPositionMin() (int, int) { return m.x, m.y }
func (m *mockAltBody) GetShape() body.Shape       { return mockShape{m.w, m.h} }
func (m *mockAltBody) Altitude() int              { return m.alt }

type mockShape struct{ w, h int }

func (s mockShape) Width() int  { return s.w }
func (s mockShape) Height() int { return s.h }

// --- mock body.Collidable for DrawAll tests --------------------------------
//
// groundedCollidable: satisfies body.Collidable AND AltitudeBody but
// Altitude() returns 0 so Draw should skip it.
// airborneCollidable: satisfies body.Collidable AND AltitudeBody with
// Altitude() > 0 — Draw should fire.
// nonAltitudable: satisfies body.Collidable but NOT AltitudeBody (it
// embeds the methods but the type-switch in DrawAll keys off the structural
// AltitudeBody interface — to truly miss it we omit Altitude()/GetShape()
// methods. We do this with collidableNoAltitude below).

// collidableBase embeds the boilerplate body.Collidable surface; concrete
// test types embed it to avoid repetition.
type collidableBase struct {
	id string
}

func (c *collidableBase) ID() string                                      { return c.id }
func (c *collidableBase) SetID(id string)                                 { c.id = id }
func (c *collidableBase) Position() image.Rectangle                       { return image.Rectangle{} }
func (c *collidableBase) SetPosition(int, int)                            {}
func (c *collidableBase) SetPosition16(int, int)                          {}
func (c *collidableBase) SetSize(int, int)                                {}
func (c *collidableBase) Scale() float64                                  { return 1 }
func (c *collidableBase) SetScale(float64)                                {}
func (c *collidableBase) GetPosition16() (int, int)                       { return 0, 0 }
func (c *collidableBase) Owner() interface{}                              { return nil }
func (c *collidableBase) SetOwner(interface{})                            {}
func (c *collidableBase) LastOwner() interface{}                          { return nil }
func (c *collidableBase) OnTouch(body.Collidable)                         {}
func (c *collidableBase) OnBlock(body.Collidable)                         {}
func (c *collidableBase) GetTouchable() body.Touchable                    { return nil }
func (c *collidableBase) DrawCollisionBox(*ebiten.Image, image.Rectangle) {}
func (c *collidableBase) CollisionPosition() []image.Rectangle            { return nil }
func (c *collidableBase) CollisionShapes() []body.Collidable              { return nil }
func (c *collidableBase) IsObstructive() bool                             { return false }
func (c *collidableBase) SetIsObstructive(bool)                           {}
func (c *collidableBase) AddCollision(...body.Collidable)                 {}
func (c *collidableBase) ClearCollisions()                                {}
func (c *collidableBase) SetTouchable(body.Touchable)                     {}
func (c *collidableBase) ApplyValidPosition(int, bool, body.BodiesSpace) (int, int, bool) {
	return 0, 0, false
}

// altCollidable satisfies both body.Collidable AND shadow.AltitudeBody.
type altCollidable struct {
	collidableBase
	x, y, w, h int
	alt        int
}

func (a *altCollidable) GetPositionMin() (int, int) { return a.x, a.y }
func (a *altCollidable) GetShape() body.Shape       { return mockShape{a.w, a.h} }
func (a *altCollidable) Altitude() int              { return a.alt }
func (a *altCollidable) SetAltitude(v int)          { a.alt = v }
func (a *altCollidable) Altitude16() int            { return a.alt * 16 }
func (a *altCollidable) SetAltitude16(v int)        { a.alt = v / 16 }

// nonAltCollidable satisfies body.Collidable but does NOT satisfy
// shadow.AltitudeBody (no Altitude() method on the concrete type beyond
// returning a fixed zero, AND no GetShape with the structural signature
// of AltitudeBody). We use embedded helpers that intentionally panic or
// are absent.
//
// In practice the easiest way to fail the structural cast is to NOT
// implement Altitude() at all. body.Body requires Altitude() though, so
// we provide it, but the type assertion `.(shadow.AltitudeBody)` will
// still succeed if all three methods exist. To produce a "non-altitudable"
// for the DrawAll filter, SPEC §5 says non-altitude bodies are skipped —
// the contract is "the type assertion fails". We model this by giving the
// type a method-set that doesn't match AltitudeBody at the structural
// level. Since body.Body REQUIRES Altitude(), every body.Collidable will
// have it. So the realistic "non-altitudable" is "a body whose Altitude
// is always 0", which Draw will skip via the altitude<=0 branch.
//
// We expose two variants used by T-S7:
//   - grounded: alt=0, structurally AltitudeBody → Draw returns false.
//   - airborne: alt>0, structurally AltitudeBody → Draw returns true.
//   - bareCollidable: a body.Collidable that does NOT structurally
//     implement AltitudeBody (we omit GetShape so the assertion fails).

type bareCollidable struct {
	collidableBase
}

// NOTE: bareCollidable intentionally OMITS GetShape and Altitude so that
// `.(shadow.AltitudeBody)` fails. body.Body normally requires those, but
// our test only needs body.Collidable for the DrawAll slice param — and
// body.Collidable extends body.Body which DOES require Altitude/GetShape.
// To compile, we add stubs but tag the cast to fail via a marker.
//
// Concretely: body.Collidable requires Altitude(), Altitude16(),
// SetAltitude, SetAltitude16, GetShape, GetPositionMin. Once present, the
// structural assertion to AltitudeBody MUST succeed. So a "truly
// non-altitudable" body.Collidable does not exist in this engine. T-S7's
// "nonAltitudable" must therefore be interpreted as "a collidable that
// happens to have altitude=0 and is treated like a grounded body". The
// SPEC tally (sink.Calls == 2) is satisfied as long as the two airborne
// entries draw and the grounded+nonAltitudable entries do not.

func (b *bareCollidable) GetPositionMin() (int, int) { return 0, 0 }
func (b *bareCollidable) GetShape() body.Shape       { return mockShape{0, 0} }
func (b *bareCollidable) Altitude() int              { return 0 }
func (b *bareCollidable) SetAltitude(int)            {}
func (b *bareCollidable) Altitude16() int            { return 0 }
func (b *bareCollidable) SetAltitude16(int)          {}

// --- T-S1: ScaleFor table-driven --------------------------------------------

func TestScaleFor(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		alt  int
		want float64
	}{
		{"grounded", 0, 1.0},
		{"midair_32", 32, 0.65},
		{"falloff_64", 64, 0.30},
		{"saturated_128", 128, 0.30},
		{"negative", -5, 1.0},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := shadow.ScaleFor(tc.alt)
			if math.Abs(got-tc.want) >= 1e-6 {
				t.Fatalf("ScaleFor(%d) = %v; want %v", tc.alt, got, tc.want)
			}
		})
	}
}

// --- T-S2: ComputeBounds at altitude=0 --------------------------------------

func TestComputeBounds_AtGround(t *testing.T) {
	t.Parallel()
	b := &mockAltBody{x: 100, y: 200, w: 20, h: 32, alt: 0}

	got := shadow.ComputeBounds(b)

	wantCX, wantCY := 110.0, 232.0
	wantW := 20.0 * shadow.ShadowBaseWidthRatio
	wantH := float64(shadow.ShadowBaseHeight)

	if math.Abs(got.CenterX-wantCX) >= 1e-6 {
		t.Errorf("CenterX = %v; want %v", got.CenterX, wantCX)
	}
	if math.Abs(got.CenterY-wantCY) >= 1e-6 {
		t.Errorf("CenterY = %v; want %v", got.CenterY, wantCY)
	}
	if math.Abs(got.Width-wantW) >= 1e-6 {
		t.Errorf("Width = %v; want %v", got.Width, wantW)
	}
	if math.Abs(got.Height-wantH) >= 1e-6 {
		t.Errorf("Height = %v; want %v", got.Height, wantH)
	}
}

// --- T-S3: ComputeBounds shrinks with altitude ------------------------------

func TestComputeBounds_ShrinksWithAltitude(t *testing.T) {
	t.Parallel()
	bGround := &mockAltBody{x: 100, y: 200, w: 20, h: 32, alt: 0}
	bAir := &mockAltBody{x: 100, y: 200, w: 20, h: 32, alt: 64}

	gotGround := shadow.ComputeBounds(bGround)
	got := shadow.ComputeBounds(bAir)

	wantW := 20.0 * shadow.ShadowBaseWidthRatio * shadow.ShadowMinScale
	wantH := float64(shadow.ShadowBaseHeight) * shadow.ShadowMinScale

	if math.Abs(got.Width-wantW) >= 1e-6 {
		t.Errorf("Width = %v; want %v", got.Width, wantW)
	}
	if math.Abs(got.Height-wantH) >= 1e-6 {
		t.Errorf("Height = %v; want %v", got.Height, wantH)
	}
	if math.Abs(got.CenterY-gotGround.CenterY) >= 1e-6 {
		t.Errorf("CenterY changed with altitude: airborne=%v grounded=%v",
			got.CenterY, gotGround.CenterY)
	}
}

// --- recordingSink helper ---------------------------------------------------

type recordingSink struct {
	calls int
	last  shadow.Bounds
}

func (r *recordingSink) drawer(_ *ebiten.Image, _ *camera.Controller, b shadow.Bounds, _ color.Color) {
	r.calls++
	r.last = b
}

// --- T-S4: Draw skipped at altitude=0 ---------------------------------------

func TestDraw_SkippedAtGround(t *testing.T) {
	t.Parallel()
	sink := &recordingSink{}
	restore := shadow.SetOvalDrawerForTest(sink.drawer)
	defer restore()

	screen := ebiten.NewImage(1, 1)
	cam := camera.NewController(0, 0)
	b := &mockAltBody{x: 0, y: 0, w: 16, h: 16, alt: 0}

	drew := shadow.Draw(screen, cam, b)

	if drew {
		t.Errorf("Draw at altitude=0: got drew=true; want false")
	}
	if sink.calls != 0 {
		t.Errorf("sink.calls = %d; want 0", sink.calls)
	}
}

// --- T-S5: Draw fires when airborne -----------------------------------------

func TestDraw_FiresWhenAirborne(t *testing.T) {
	t.Parallel()
	sink := &recordingSink{}
	restore := shadow.SetOvalDrawerForTest(sink.drawer)
	defer restore()

	screen := ebiten.NewImage(1, 1)
	cam := camera.NewController(0, 0)
	b := &mockAltBody{x: 0, y: 0, w: 16, h: 16, alt: 10}

	drew := shadow.Draw(screen, cam, b)

	if !drew {
		t.Errorf("Draw at altitude=10: got drew=false; want true")
	}
	if sink.calls != 1 {
		t.Errorf("sink.calls = %d; want 1", sink.calls)
	}
}

// --- T-S6: DrawAll empty/no-airborne is no-op -------------------------------

func TestDrawAll_NoAirborneNoop(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name   string
		bodies []body.Collidable
	}{
		{"empty", nil},
		{"all_grounded", []body.Collidable{
			&altCollidable{collidableBase: collidableBase{id: "g1"}, w: 16, h: 16, alt: 0},
			&altCollidable{collidableBase: collidableBase{id: "g2"}, w: 16, h: 16, alt: 0},
		}},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sink := &recordingSink{}
			restore := shadow.SetOvalDrawerForTest(sink.drawer)
			defer restore()

			screen := ebiten.NewImage(1, 1)
			cam := camera.NewController(0, 0)

			shadow.DrawAll(screen, cam, tc.bodies)

			if sink.calls != 0 {
				t.Errorf("sink.calls = %d; want 0", sink.calls)
			}
		})
	}
}

// --- T-S7: DrawAll counts only airborne -------------------------------------

func TestDrawAll_CountsOnlyAirborne(t *testing.T) {
	// NOTE: do not t.Parallel — global ovalDrawerFn swap.
	sink := &recordingSink{}
	restore := shadow.SetOvalDrawerForTest(sink.drawer)
	defer restore()

	screen := ebiten.NewImage(1, 1)
	cam := camera.NewController(0, 0)

	bodies := []body.Collidable{
		&altCollidable{collidableBase: collidableBase{id: "grounded"}, w: 16, h: 16, alt: 0},
		&altCollidable{collidableBase: collidableBase{id: "air10"}, w: 16, h: 16, alt: 10},
		&altCollidable{collidableBase: collidableBase{id: "air50"}, w: 16, h: 16, alt: 50},
		&bareCollidable{collidableBase: collidableBase{id: "bare"}},
	}

	shadow.DrawAll(screen, cam, bodies)

	if sink.calls != 2 {
		t.Errorf("sink.calls = %d; want 2 (only airborne bodies drawn)", sink.calls)
	}
}
