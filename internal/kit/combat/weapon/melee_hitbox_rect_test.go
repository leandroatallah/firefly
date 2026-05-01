package weapon_test

import (
	"image"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
)

// expectedRectFromTestStep computes the hitbox rect for the single-step
// weapon used by newTestMeleeWeapon (HitboxW16=24*16, HitboxH16=16*16,
// HitboxOffsetX16=12*16, HitboxOffsetY16=0). It mirrors the production
// pixel-space math so the test is independent of the implementation.
func expectedRectFromTestStep(originX16, originY16 int, face animation.FacingDirectionEnum) image.Rectangle {
	const (
		hitboxW16       = 24 * 16
		hitboxH16       = 16 * 16
		hitboxOffsetX16 = 12 * 16
		hitboxOffsetY16 = 0
	)
	halfW16 := hitboxW16 / 2
	var centerX16 int
	if face == animation.FaceDirectionLeft {
		centerX16 = originX16 - hitboxOffsetX16
	} else {
		centerX16 = originX16 + hitboxOffsetX16
	}
	x0 := (centerX16 - halfW16) / 16
	x1 := (centerX16 + halfW16) / 16
	y0 := (originY16 + hitboxOffsetY16) / 16
	y1 := y0 + hitboxH16/16
	return image.Rect(x0, y0, x1, y1)
}

// TestMeleeWeapon_ActiveHitboxRect verifies the (rect, active) accessor
// the Phase Scene uses to gate orange debug-rect rendering.
func TestMeleeWeapon_ActiveHitboxRect(t *testing.T) {
	const (
		ownerPxX = 100
		ownerPxY = 100
	)

	tests := []struct {
		name        string
		face        animation.FacingDirectionEnum
		framesToRun int  // number of Update() ticks after Fire
		fire        bool // whether to call Fire()
		wantActive  bool
		wantRect    image.Rectangle
	}{
		{
			name:        "not swinging returns zero rect and false",
			face:        animation.FaceDirectionRight,
			framesToRun: 0,
			fire:        false,
			wantActive:  false,
			wantRect:    image.Rectangle{},
		},
		{
			name:        "in startup window returns zero rect and false",
			face:        animation.FaceDirectionRight,
			framesToRun: 1, // ActiveFrames[0]==3, so frame 1 is pre-active
			fire:        true,
			wantActive:  false,
			wantRect:    image.Rectangle{},
		},
		{
			name:        "active window facing right returns rect and true",
			face:        animation.FaceDirectionRight,
			framesToRun: 3, // ActiveFrames[0]
			fire:        true,
			wantActive:  true,
			wantRect:    expectedRectFromTestStep(ownerPxX*16, ownerPxY*16, animation.FaceDirectionRight),
		},
		{
			name:        "active window facing left returns mirrored rect and true",
			face:        animation.FaceDirectionLeft,
			framesToRun: 3,
			fire:        true,
			wantActive:  true,
			wantRect:    expectedRectFromTestStep(ownerPxX*16, ownerPxY*16, animation.FaceDirectionLeft),
		},
		{
			name:        "past active window returns zero rect and false",
			face:        animation.FaceDirectionRight,
			framesToRun: 7, // ActiveFrames[1]==5, 6 ends swing; 7 is well past
			fire:        true,
			wantActive:  false,
			wantRect:    image.Rectangle{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			owner := newMeleeOwner(ownerPxX, ownerPxY, combat.FactionPlayer, tc.face)
			w := newTestMeleeWeapon(owner)

			if tc.fire {
				w.Fire(owner.x16, owner.y16, owner.face, body.ShootDirectionStraight, 0)
				for i := 0; i < tc.framesToRun; i++ {
					w.Update()
				}
			}

			rect, active := w.ActiveHitboxRect()
			if active != tc.wantActive {
				t.Fatalf("ActiveHitboxRect() active = %v, want %v", active, tc.wantActive)
			}
			if rect != tc.wantRect {
				t.Errorf("ActiveHitboxRect() rect = %+v, want %+v", rect, tc.wantRect)
			}
		})
	}
}

// TestMeleeWeapon_HitboxRect_ParityWithActiveHitboxRect verifies that during
// the active window, HitboxRect() and ActiveHitboxRect() return byte-identical
// rectangles. This is the single source-of-truth invariant the debug overlay
// relies on (orange box must align with the same rect ApplyHitbox queries).
func TestMeleeWeapon_HitboxRect_ParityWithActiveHitboxRect(t *testing.T) {
	tests := []struct {
		name string
		face animation.FacingDirectionEnum
	}{
		{"facing right", animation.FaceDirectionRight},
		{"facing left", animation.FaceDirectionLeft},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			owner := newMeleeOwner(100, 100, combat.FactionPlayer, tc.face)
			w := newTestMeleeWeapon(owner)
			w.Fire(owner.x16, owner.y16, owner.face, body.ShootDirectionStraight, 0)
			for i := 0; i < 3; i++ {
				w.Update()
			}
			if !w.IsHitboxActive() {
				t.Fatalf("precondition: IsHitboxActive() = false at frame 3, want true")
			}

			plain := w.HitboxRect()
			active, ok := w.ActiveHitboxRect()
			if !ok {
				t.Fatalf("ActiveHitboxRect() ok = false during active window, want true")
			}
			if plain != active {
				t.Errorf("HitboxRect()=%+v != ActiveHitboxRect()=%+v (must be byte-identical)", plain, active)
			}
		})
	}
}
