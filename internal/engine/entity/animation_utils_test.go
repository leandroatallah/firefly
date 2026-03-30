package entity

import (
	"image"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/render/sprites"
)

// mockAnimatable implements Animatable for testing
type mockAnimatable struct {
	sprite    *sprites.Sprite
	pos       image.Rectangle
	frameRate int
}

func (m *mockAnimatable) GetSpriteByState(state animation.SpriteState) *sprites.Sprite {
	return m.sprite
}

func (m *mockAnimatable) Position() image.Rectangle {
	return m.pos
}

func (m *mockAnimatable) FrameRate() int {
	return m.frameRate
}

func TestIsAnimationFinished(t *testing.T) {
	tests := []struct {
		name     string
		entity   Animatable
		state    animation.SpriteState
		tick     int
		expected bool
	}{
		{
			name:     "nil entity returns true",
			entity:   nil,
			state:    0,
			tick:     0,
			expected: true,
		},
		{
			name: "nil sprite returns true",
			entity: &mockAnimatable{
				sprite:    nil,
				pos:       image.Rect(0, 0, 32, 32),
				frameRate: 1,
			},
			state:    0,
			tick:     0,
			expected: true,
		},
		{
			name: "sprite with nil image returns true",
			entity: &mockAnimatable{
				sprite:    &sprites.Sprite{Image: nil},
				pos:       image.Rect(0, 0, 32, 32),
				frameRate: 1,
			},
			state:    0,
			tick:     0,
			expected: true,
		},
		{
			name: "zero width rect returns true",
			entity: &mockAnimatable{
				sprite:    &sprites.Sprite{Image: ebiten.NewImage(32, 32)},
				pos:       image.Rect(0, 0, 0, 32),
				frameRate: 1,
			},
			state:    0,
			tick:     0,
			expected: true,
		},
		{
			name: "zero frame rate defaults to 1",
			entity: &mockAnimatable{
				sprite:    &sprites.Sprite{Image: ebiten.NewImage(64, 32)},
				pos:       image.Rect(0, 0, 32, 32),
				frameRate: 0,
			},
			state:    0,
			tick:     3, // 2 frames * 1 fps = 2, so tick 3 is finished
			expected: true,
		},
		{
			name: "animation not finished",
			entity: &mockAnimatable{
				sprite:    &sprites.Sprite{Image: ebiten.NewImage(64, 32)},
				pos:       image.Rect(0, 0, 32, 32),
				frameRate: 1,
			},
			state:    0,
			tick:     1, // 2 frames * 1 fps = 2, so tick 1 is not finished
			expected: false,
		},
		{
			name: "animation finished exactly at duration",
			entity: &mockAnimatable{
				sprite:    &sprites.Sprite{Image: ebiten.NewImage(96, 32)},
				pos:       image.Rect(0, 0, 32, 32),
				frameRate: 2,
			},
			state:    0,
			tick:     6, // 3 frames * 2 fps = 6
			expected: true,
		},
		{
			name: "animation finished past duration",
			entity: &mockAnimatable{
				sprite:    &sprites.Sprite{Image: ebiten.NewImage(96, 32)},
				pos:       image.Rect(0, 0, 32, 32),
				frameRate: 2,
			},
			state:    0,
			tick:     10, // 3 frames * 2 fps = 6, tick 10 > 6
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsAnimationFinished(tt.tick, tt.entity, tt.state)
			if result != tt.expected {
				t.Errorf("IsAnimationFinished(%d, entity, %v) = %v, want %v",
					tt.tick, tt.state, result, tt.expected)
			}
		})
	}
}
