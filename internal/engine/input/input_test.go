package input

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestIsSomeKeyPressed(t *testing.T) {
	tests := []struct {
		name      string
		keys      []ebiten.Key
		pressed   map[ebiten.Key]bool
		want      bool
	}{
		{
			name: "no keys given",
			keys: []ebiten.Key{},
			want: false,
		},
		{
			name: "one key pressed",
			keys: []ebiten.Key{ebiten.KeyA, ebiten.KeyB},
			pressed: map[ebiten.Key]bool{ebiten.KeyA: true},
			want: true,
		},
		{
			name: "no keys pressed",
			keys: []ebiten.Key{ebiten.KeyA, ebiten.KeyB},
			pressed: map[ebiten.Key]bool{},
			want: false,
		},
		{
			name: "all keys pressed",
			keys: []ebiten.Key{ebiten.KeyA, ebiten.KeyB},
			pressed: map[ebiten.Key]bool{ebiten.KeyA: true, ebiten.KeyB: true},
			want: true,
		},
	}

	originalIsKeyPressed := isKeyPressed
	defer func() { isKeyPressed = originalIsKeyPressed }()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isKeyPressed = func(k ebiten.Key) bool {
				return tt.pressed[k]
			}
			if got := IsSomeKeyPressed(tt.keys...); got != tt.want {
				t.Errorf("IsSomeKeyPressed() = %v, want %v", got, tt.want)
			}
		})
	}
}
