package actors_test

import (
	"image"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
)

// T-C1 (story 070): SetRenderOffset / RenderOffset round-trip with facing-aware
// X resolution. X is auto-mirrored (negated) when the character faces left so
// the screen-space nudge follows the mirrored sprite content. Y is identical
// for both facings. Verifies SPEC §2 and AC-3/AC-4/AC-9.
func TestCharacter_RenderOffset_FacingAware(t *testing.T) {
	tests := []struct {
		name   string
		dx     int
		dy     int
		facing animation.FacingDirectionEnum
		wantPt image.Point
	}{
		{
			name:   "facing right uses X as-is",
			dx:     10,
			dy:     2,
			facing: animation.FaceDirectionRight,
			wantPt: image.Pt(10, 2),
		},
		{
			name:   "facing left auto-mirrors X",
			dx:     10,
			dy:     2,
			facing: animation.FaceDirectionLeft,
			wantPt: image.Pt(-10, 2),
		},
		{
			name:   "facing left auto-mirrors negative X",
			dx:     -4,
			dy:     2,
			facing: animation.FaceDirectionLeft,
			wantPt: image.Pt(4, 2),
		},
		{
			name:   "facing left with X=0 stays 0",
			dx:     0,
			dy:     3,
			facing: animation.FaceDirectionLeft,
			wantPt: image.Pt(0, 3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newRenderOffsetCharacter()
			// Zero acceleration so UpdateImageOptions does not override the
			// facing we set below (see NOTES.md risk: accX sign override).
			c.SetAcceleration(0, 0)
			c.SetFaceDirection(tt.facing)

			c.SetRenderOffset(actors.Idle, tt.dx, tt.dy)

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
// what other states have been registered.
func TestCharacter_RenderOffset_FacingAware_UnregisteredState(t *testing.T) {
	c := newRenderOffsetCharacter()
	c.SetRenderOffset(actors.Idle, -4, 2)

	got, ok := c.RenderOffset(actors.Walking)
	if ok {
		t.Errorf("RenderOffset(Walking) ok = true; want false (never registered); got %v", got)
	}
	if got != (image.Point{}) {
		t.Errorf("RenderOffset(Walking) = %v, want zero Point", got)
	}
}

// T-C2 (story 070): UpdateImageOptions applies the facing-resolved render
// offset as the final additive translation. X is auto-mirrored when facing
// left; Y is facing-independent. AC-3 / AC-9.
func TestCharacter_UpdateImageOptions_AppliesFacingAwareRenderOffset(t *testing.T) {
	baseR := newRenderOffsetCharacter()
	baseR.SetAcceleration(0, 0)
	baseR.SetFaceDirection(animation.FaceDirectionRight)
	baseRtx, baseRty := translation(baseR)

	baseL := newRenderOffsetCharacter()
	baseL.SetAcceleration(0, 0)
	baseL.SetFaceDirection(animation.FaceDirectionLeft)
	baseLtx, baseLty := translation(baseL)

	tests := []struct {
		name   string
		dx     int
		dy     int
		facing animation.FacingDirectionEnum
		wantDx float64
		wantDy float64
	}{
		{
			name:   "X positive, facing right",
			dx:     10,
			dy:     2,
			facing: animation.FaceDirectionRight,
			wantDx: 10,
			wantDy: 2,
		},
		{
			name:   "X positive, facing left auto-mirrors",
			dx:     10,
			dy:     2,
			facing: animation.FaceDirectionLeft,
			wantDx: -10,
			wantDy: 2,
		},
		{
			name:   "X negative, facing left auto-mirrors",
			dx:     -4,
			dy:     2,
			facing: animation.FaceDirectionLeft,
			wantDx: 4,
			wantDy: 2,
		},
		{
			name:   "Y only, facing left",
			dx:     0,
			dy:     3,
			facing: animation.FaceDirectionLeft,
			wantDx: 0,
			wantDy: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newRenderOffsetCharacter()
			c.SetAcceleration(0, 0)
			c.SetFaceDirection(tt.facing)
			c.SetRenderOffset(actors.Idle, tt.dx, tt.dy)

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
// facing direction.
func TestCharacter_UpdateImageOptions_YIdenticalAcrossFacings(t *testing.T) {
	const dy = 3

	baseR := newRenderOffsetCharacter()
	baseR.SetAcceleration(0, 0)
	baseR.SetFaceDirection(animation.FaceDirectionRight)
	_, baseRty := translation(baseR)

	baseL := newRenderOffsetCharacter()
	baseL.SetAcceleration(0, 0)
	baseL.SetFaceDirection(animation.FaceDirectionLeft)
	_, baseLty := translation(baseL)

	right := newRenderOffsetCharacter()
	right.SetAcceleration(0, 0)
	right.SetFaceDirection(animation.FaceDirectionRight)
	right.SetRenderOffset(actors.Idle, -4, dy)
	_, rightTy := translation(right)

	left := newRenderOffsetCharacter()
	left.SetAcceleration(0, 0)
	left.SetFaceDirection(animation.FaceDirectionLeft)
	left.SetRenderOffset(actors.Idle, -4, dy)
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
