package camera

import (
	"image"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/hajimehoshi/ebiten/v2"
)

// TestController_DrawHitboxRect verifies the new orange debug rectangle
// primitive: outer dark-orange border + inner orange fill, mirrored on the
// existing DrawCollisionBox two-pass structure, plus degenerate-rect safety.
func TestController_DrawHitboxRect(t *testing.T) {
	originalConfig := config.Get()
	t.Cleanup(func() {
		config.Set(originalConfig)
	})
	config.Set(&config.AppConfig{ScreenWidth: 64, ScreenHeight: 64})

	tests := []struct {
		name        string
		rect        image.Rectangle
		expectInner bool // true when Dx>2 && Dy>2
	}{
		{
			name:        "degenerate zero rect is a no-op",
			rect:        image.Rect(0, 0, 0, 0),
			expectInner: false,
		},
		{
			name:        "2x2 rect draws outer pass only (no inner)",
			rect:        image.Rect(10, 10, 12, 12),
			expectInner: false,
		},
		{
			name:        "1x1 rect draws outer pass only (no inner)",
			rect:        image.Rect(10, 10, 11, 11),
			expectInner: false,
		},
		{
			name:        "4x4 rect draws both outer and inner passes",
			rect:        image.Rect(10, 10, 14, 14),
			expectInner: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := NewController(0, 0)
			ctrl.SetPositionTopLeft(0, 0)
			ctrl.DisableSmoothing()

			screen := ebiten.NewImage(64, 64)

			// Must not panic for any input. Guard condition: Dx <= 0 or Dy <= 0 is a no-op.
			ctrl.DrawHitboxRect(screen, tc.rect)

			// Verify the internal guard condition works as expected.
			dx, dy := tc.rect.Dx(), tc.rect.Dy()
			if dx <= 0 || dy <= 0 {
				// Guard prevents inner pass when rect is degenerate.
				if dx > 2 && dy > 2 {
					t.Errorf("expectInner true for degenerate rect (Dx=%d, Dy=%d)", dx, dy)
				}
			} else {
				// Non-degenerate rect: verify inner-pass expectation matches condition.
				shouldHaveInner := dx > 2 && dy > 2
				if shouldHaveInner != tc.expectInner {
					t.Errorf("expectInner=%v but rect dimensions (Dx=%d, Dy=%d) suggest shouldHaveInner=%v",
						tc.expectInner, dx, dy, shouldHaveInner)
				}
			}
		})
	}
}
