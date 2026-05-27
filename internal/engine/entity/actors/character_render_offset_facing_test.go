package actors_test

import (
	"image"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
)

// intPtr is a small helper so tests can pass *int literals to SetRenderOffset.
// Go does not allow taking the address of an int literal directly.
func intPtr(v int) *int { return &v }

// T-C1 (story 070): SetRenderOffset / RenderOffset round-trip with facing-aware
// X resolution. The optional dxFlipped pointer becomes the X translation when
// the character faces left; nil falls back to dx for both facings. Verifies the
// behavior described in SPEC §2 and AC-3/AC-4/AC-9.
func TestCharacter_RenderOffset_FacingAware(t *testing.T) {
	tests := []struct {
		name      string
		dx        int
		dy        int
		dxFlipped *int
		facing    animation.FacingDirectionEnum
		wantPt    image.Point
	}{
		{
			name:      "nil flipped, facing right uses X",
			dx:        -4,
			dy:        2,
			dxFlipped: nil,
			facing:    animation.FaceDirectionRight,
			wantPt:    image.Pt(-4, 2),
		},
		{
			name:      "nil flipped, facing left falls back to X (story 068 regression)",
			dx:        -4,
			dy:        2,
			dxFlipped: nil,
			facing:    animation.FaceDirectionLeft,
			wantPt:    image.Pt(-4, 2),
		},
		{
			name:      "flipped set, facing right uses X (not XFlipped)",
			dx:        -4,
			dy:        2,
			dxFlipped: intPtr(6),
			facing:    animation.FaceDirectionRight,
			wantPt:    image.Pt(-4, 2),
		},
		{
			name:      "flipped set, facing left uses XFlipped",
			dx:        -4,
			dy:        2,
			dxFlipped: intPtr(6),
			facing:    animation.FaceDirectionLeft,
			wantPt:    image.Pt(6, 2),
		},
		{
			name:      "flipped=0 explicit override, facing left",
			dx:        -4,
			dy:        2,
			dxFlipped: intPtr(0),
			facing:    animation.FaceDirectionLeft,
			wantPt:    image.Pt(0, 2),
		},
		{
			name:      "flipped=0 explicit override, facing right still uses X",
			dx:        -4,
			dy:        2,
			dxFlipped: intPtr(0),
			facing:    animation.FaceDirectionRight,
			wantPt:    image.Pt(-4, 2),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newRenderOffsetCharacter()
			// Zero acceleration so UpdateImageOptions does not override the
			// facing we set below (see NOTES.md risk: accX sign override).
			c.SetAcceleration(0, 0)
			c.SetFaceDirection(tt.facing)

			c.SetRenderOffset(actors.Idle, tt.dx, tt.dy, tt.dxFlipped)

			got, ok := c.RenderOffset(actors.Idle)
			if !ok {
				t.Fatalf("RenderOffset(Idle) ok = false; want true after SetRenderOffset")
			}
			if got != tt.wantPt {
				t.Errorf("RenderOffset(Idle) facing=%v = %v, want %v",
					tt.facing, got, tt.wantPt)
			}
		})
	}
}

// T-C1b (story 070): An unregistered state returns ok=false regardless of
// what other states have been registered (facing-aware variant of the
// existing 068 unregistered-state check).
func TestCharacter_RenderOffset_FacingAware_UnregisteredState(t *testing.T) {
	c := newRenderOffsetCharacter()
	c.SetRenderOffset(actors.Idle, -4, 2, intPtr(6))

	got, ok := c.RenderOffset(actors.Walking)
	if ok {
		t.Errorf("RenderOffset(Walking) ok = true; want false (never registered); got %v", got)
	}
	if got != (image.Point{}) {
		t.Errorf("RenderOffset(Walking) = %v, want zero Point", got)
	}
}

// T-C2 (story 070): UpdateImageOptions applies the facing-resolved render
// offset as the final additive translation. The X picked must depend on the
// current facing direction at every call (no caching). Verifies AC-3 and AC-9.
//
// Baseline is recaptured per-facing because facing-left mirroring shifts the
// GeoM translation row independent of any per-state offset.
func TestCharacter_UpdateImageOptions_AppliesFacingAwareRenderOffset(t *testing.T) {
	// Right-facing baseline.
	baseR := newRenderOffsetCharacter()
	baseR.SetAcceleration(0, 0)
	baseR.SetFaceDirection(animation.FaceDirectionRight)
	baseRtx, baseRty := translation(baseR)

	// Left-facing baseline.
	baseL := newRenderOffsetCharacter()
	baseL.SetAcceleration(0, 0)
	baseL.SetFaceDirection(animation.FaceDirectionLeft)
	baseLtx, baseLty := translation(baseL)

	tests := []struct {
		name      string
		dx        int
		dy        int
		dxFlipped *int
		facing    animation.FacingDirectionEnum
		wantDx    float64
		wantDy    float64
	}{
		{
			name:   "X only, facing right",
			dx:     -4,
			dy:     2,
			facing: animation.FaceDirectionRight,
			wantDx: -4,
			wantDy: 2,
		},
		{
			name:   "X only, facing left -> same X (story 068 regression)",
			dx:     -4,
			dy:     2,
			facing: animation.FaceDirectionLeft,
			wantDx: -4,
			wantDy: 2,
		},
		{
			name:      "X + XFlipped, facing right -> uses X",
			dx:        -4,
			dy:        2,
			dxFlipped: intPtr(6),
			facing:    animation.FaceDirectionRight,
			wantDx:    -4,
			wantDy:    2,
		},
		{
			name:      "X + XFlipped, facing left -> uses XFlipped",
			dx:        -4,
			dy:        2,
			dxFlipped: intPtr(6),
			facing:    animation.FaceDirectionLeft,
			wantDx:    6,
			wantDy:    2,
		},
		{
			name:      "XFlipped=0 explicit, facing left -> 0",
			dx:        -4,
			dy:        2,
			dxFlipped: intPtr(0),
			facing:    animation.FaceDirectionLeft,
			wantDx:    0,
			wantDy:    2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newRenderOffsetCharacter()
			c.SetAcceleration(0, 0)
			c.SetFaceDirection(tt.facing)
			c.SetRenderOffset(actors.Idle, tt.dx, tt.dy, tt.dxFlipped)

			gotTx, gotTy := translation(c)

			var wantTx, wantTy float64
			if tt.facing == animation.FaceDirectionLeft {
				wantTx = baseLtx + tt.wantDx
				wantTy = baseLty + tt.wantDy
			} else {
				wantTx = baseRtx + tt.wantDx
				wantTy = baseRty + tt.wantDy
			}

			if gotTx != wantTx || gotTy != wantTy {
				t.Errorf("translation facing=%v = (%v, %v), want (%v, %v); delta=(%v, %v)",
					tt.facing, gotTx, gotTy, wantTx, wantTy, tt.wantDx, tt.wantDy)
			}
		})
	}
}

// T-C2b (story 070): Y is identical for both facings — only X depends on
// facing direction. SPEC §3 / NOTES.md: AC-9 row "Y identical for both facings".
func TestCharacter_UpdateImageOptions_YIdenticalAcrossFacings(t *testing.T) {
	const dy = 3

	// Capture per-facing baselines (Y component only).
	baseR := newRenderOffsetCharacter()
	baseR.SetAcceleration(0, 0)
	baseR.SetFaceDirection(animation.FaceDirectionRight)
	_, baseRty := translation(baseR)

	baseL := newRenderOffsetCharacter()
	baseL.SetAcceleration(0, 0)
	baseL.SetFaceDirection(animation.FaceDirectionLeft)
	_, baseLty := translation(baseL)

	flipped := 6

	right := newRenderOffsetCharacter()
	right.SetAcceleration(0, 0)
	right.SetFaceDirection(animation.FaceDirectionRight)
	right.SetRenderOffset(actors.Idle, -4, dy, &flipped)
	_, rightTy := translation(right)

	left := newRenderOffsetCharacter()
	left.SetAcceleration(0, 0)
	left.SetFaceDirection(animation.FaceDirectionLeft)
	left.SetRenderOffset(actors.Idle, -4, dy, &flipped)
	_, leftTy := translation(left)

	gotRightDelta := rightTy - baseRty
	gotLeftDelta := leftTy - baseLty
	if gotRightDelta != float64(dy) {
		t.Errorf("facing-right Y delta = %v, want %v", gotRightDelta, dy)
	}
	if gotLeftDelta != float64(dy) {
		t.Errorf("facing-left Y delta = %v, want %v", gotLeftDelta, dy)
	}
	if gotRightDelta != gotLeftDelta {
		t.Errorf("Y delta differs by facing (right=%v, left=%v); Y must be facing-independent",
			gotRightDelta, gotLeftDelta)
	}
}
