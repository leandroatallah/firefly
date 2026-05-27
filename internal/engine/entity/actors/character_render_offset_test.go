package actors_test

import (
	"image"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
	"github.com/boilerplate/ebiten-template/internal/engine/render/sprites"
	"github.com/hajimehoshi/ebiten/v2"
)

// newRenderOffsetCharacter builds a Character with a single non-trivial sprite
// for Idle so UpdateImageOptions exercises the full draw transform pipeline.
func newRenderOffsetCharacter() *actors.Character {
	img := ebiten.NewImage(32, 32)
	sMap := sprites.SpriteMap{
		actors.Idle:    &sprites.Sprite{Image: img},
		actors.Walking: &sprites.Sprite{Image: img},
	}
	rect := bodyphysics.NewRect(100, 200, 16, 24)
	c := actors.NewCharacter(sMap, rect)
	c.SetMaxHealth(100)
	c.SetHealth(100)
	return c
}

// translation extracts the (tx, ty) translation row of the current image
// options GeoM after a call to UpdateImageOptions.
func translation(c *actors.Character) (float64, float64) {
	c.UpdateImageOptions()
	opts := c.ImageOptions()
	geom := opts.GeoM
	return geom.Element(0, 2), geom.Element(1, 2)
}

// T-C1: SetRenderOffset / RenderOffset round-trip.
// Registering an offset for a state must be retrievable; unregistered states
// must report ok=false and a zero Point.
func TestCharacter_RenderOffset_RoundTrip(t *testing.T) {
	c := newRenderOffsetCharacter()

	c.SetRenderOffset(actors.Idle, -4, 2)

	gotIdle, okIdle := c.RenderOffset(actors.Idle)
	if !okIdle {
		t.Error("RenderOffset(Idle) ok = false; want true after SetRenderOffset")
	}
	if gotIdle != image.Pt(-4, 2) {
		t.Errorf("RenderOffset(Idle) = %v, want (-4, 2)", gotIdle)
	}

	gotWalking, okWalking := c.RenderOffset(actors.Walking)
	if okWalking {
		t.Error("RenderOffset(Walking) ok = true; want false (never registered)")
	}
	if gotWalking != (image.Point{}) {
		t.Errorf("RenderOffset(Walking) = %v, want zero Point", gotWalking)
	}
}

// T-C1b: A fresh Character has no render offsets registered. RenderOffset
// for any state returns (Point{}, false). Calling UpdateImageOptions with
// the nil internal map must not panic.
func TestCharacter_RenderOffset_DefaultsNil(t *testing.T) {
	c := newRenderOffsetCharacter()

	got, ok := c.RenderOffset(actors.Idle)
	if ok {
		t.Error("fresh Character: RenderOffset(Idle) ok = true; want false")
	}
	if got != (image.Point{}) {
		t.Errorf("fresh Character: RenderOffset(Idle) = %v, want zero Point", got)
	}

	// Must not panic even when no offsets are registered.
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("UpdateImageOptions panicked with no offsets registered: %v", r)
		}
	}()
	c.UpdateImageOptions()
}

// T-C2: UpdateImageOptions applies the per-state offset as the final additive
// translation. Verified by comparing the GeoM translation row to a baseline
// captured before any offset is registered.
func TestCharacter_UpdateImageOptions_AppliesRenderOffset(t *testing.T) {
	// Baseline: no offsets registered → baseline translation.
	base := newRenderOffsetCharacter()
	baseTx, baseTy := translation(base)

	tests := []struct {
		name   string
		state  actors.ActorStateEnum
		setup  func(*actors.Character)
		wantDx float64
		wantDy float64
	}{
		{
			name:   "no offset registered yields baseline translation",
			state:  actors.Idle,
			setup:  func(c *actors.Character) {},
			wantDx: 0,
			wantDy: 0,
		},
		{
			name:  "offset (-4, 0) shifts X left by 4",
			state: actors.Idle,
			setup: func(c *actors.Character) {
				c.SetRenderOffset(actors.Idle, -4, 0)
			},
			wantDx: -4,
			wantDy: 0,
		},
		{
			name:  "offset (0, 3) shifts Y down by 3",
			state: actors.Idle,
			setup: func(c *actors.Character) {
				c.SetRenderOffset(actors.Idle, 0, 3)
			},
			wantDx: 0,
			wantDy: 3,
		},
		{
			name:  "offset registered for another state has no effect on current state",
			state: actors.Idle,
			setup: func(c *actors.Character) {
				c.SetRenderOffset(actors.Walking, -10, -10)
			},
			wantDx: 0,
			wantDy: 0,
		},
		{
			name:  "explicit zero offset behaves like no offset",
			state: actors.Idle,
			setup: func(c *actors.Character) {
				c.SetRenderOffset(actors.Idle, 0, 0)
			},
			wantDx: 0,
			wantDy: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newRenderOffsetCharacter()
			tt.setup(c)

			gotTx, gotTy := translation(c)

			wantTx := baseTx + tt.wantDx
			wantTy := baseTy + tt.wantDy
			if gotTx != wantTx || gotTy != wantTy {
				t.Errorf("translation = (%v, %v), want (%v, %v); baseline=(%v, %v) delta=(%v, %v)",
					gotTx, gotTy, wantTx, wantTy, baseTx, baseTy, tt.wantDx, tt.wantDy)
			}
		})
	}
}

// T-C3: Facing-left mirroring does NOT invert the offset's X. The offset is
// applied AFTER Scale(-1, 1), so a left-facing actor with X:-4 is still
// shifted -4 px on screen (additive, not mirrored). Decision documented in
// NOTES.md and SPEC.md §3.
func TestCharacter_UpdateImageOptions_RenderOffset_NotMirroredWhenFacingLeft(t *testing.T) {
	// Capture left-facing baseline (no offset registered).
	base := newRenderOffsetCharacter()
	base.SetFaceDirection(animation.FaceDirectionLeft)
	baseTx, baseTy := translation(base)

	c := newRenderOffsetCharacter()
	c.SetFaceDirection(animation.FaceDirectionLeft)
	c.SetRenderOffset(actors.Idle, -4, 0)

	gotTx, gotTy := translation(c)

	wantTx := baseTx + (-4) // additive, NOT mirrored to +4
	wantTy := baseTy
	if gotTx != wantTx || gotTy != wantTy {
		t.Errorf("facing-left translation = (%v, %v), want (%v, %v); left-baseline=(%v, %v)",
			gotTx, gotTy, wantTx, wantTy, baseTx, baseTy)
	}

	// Sanity guard: if the offset had been mirrored, the X would equal
	// baseTx + 4 instead. Make the regression explicit.
	if gotTx == baseTx+4 {
		t.Errorf("offset X was mirrored when facing left (got %v = baseline + 4); "+
			"AC-9 / SPEC §3 require additive, non-mirrored behavior", gotTx)
	}
}
