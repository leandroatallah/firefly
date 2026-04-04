package input

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestReadPlayerCommands(t *testing.T) {
	tests := []struct {
		name     string
		stubKeys map[ebiten.Key]bool
		want     PlayerCommands
	}{
		{
			name:     "no keys pressed",
			stubKeys: map[ebiten.Key]bool{},
			want:     PlayerCommands{},
		},
		{
			name:     "shoot only",
			stubKeys: map[ebiten.Key]bool{ebiten.KeyX: true},
			want:     PlayerCommands{Shoot: true},
		},
		{
			name:     "up via KeyUp",
			stubKeys: map[ebiten.Key]bool{ebiten.KeyUp: true},
			want:     PlayerCommands{Up: true},
		},
		{
			name:     "up via KeyW",
			stubKeys: map[ebiten.Key]bool{ebiten.KeyW: true},
			want:     PlayerCommands{Up: true},
		},
		{
			name:     "up via both KeyUp and KeyW",
			stubKeys: map[ebiten.Key]bool{ebiten.KeyUp: true, ebiten.KeyW: true},
			want:     PlayerCommands{Up: true},
		},
		{
			name:     "down via KeyDown",
			stubKeys: map[ebiten.Key]bool{ebiten.KeyDown: true},
			want:     PlayerCommands{Down: true},
		},
		{
			name:     "down via KeyS",
			stubKeys: map[ebiten.Key]bool{ebiten.KeyS: true},
			want:     PlayerCommands{Down: true},
		},
		{
			name:     "down via both KeyDown and KeyS",
			stubKeys: map[ebiten.Key]bool{ebiten.KeyDown: true, ebiten.KeyS: true},
			want:     PlayerCommands{Down: true},
		},
		{
			name:     "left via KeyLeft",
			stubKeys: map[ebiten.Key]bool{ebiten.KeyLeft: true},
			want:     PlayerCommands{Left: true},
		},
		{
			name:     "left via KeyA",
			stubKeys: map[ebiten.Key]bool{ebiten.KeyA: true},
			want:     PlayerCommands{Left: true},
		},
		{
			name:     "left via both KeyLeft and KeyA",
			stubKeys: map[ebiten.Key]bool{ebiten.KeyLeft: true, ebiten.KeyA: true},
			want:     PlayerCommands{Left: true},
		},
		{
			name:     "right via KeyRight",
			stubKeys: map[ebiten.Key]bool{ebiten.KeyRight: true},
			want:     PlayerCommands{Right: true},
		},
		{
			name:     "right via KeyD",
			stubKeys: map[ebiten.Key]bool{ebiten.KeyD: true},
			want:     PlayerCommands{Right: true},
		},
		{
			name:     "right via both KeyRight and KeyD",
			stubKeys: map[ebiten.Key]bool{ebiten.KeyRight: true, ebiten.KeyD: true},
			want:     PlayerCommands{Right: true},
		},
		{
			name:     "jump",
			stubKeys: map[ebiten.Key]bool{ebiten.KeySpace: true},
			want:     PlayerCommands{Jump: true},
		},
		{
			name:     "dash",
			stubKeys: map[ebiten.Key]bool{ebiten.KeyShift: true},
			want:     PlayerCommands{Dash: true},
		},
		{
			name:     "confirm",
			stubKeys: map[ebiten.Key]bool{ebiten.KeyEnter: true},
			want:     PlayerCommands{Confirm: true},
		},
		{
			name:     "cancel",
			stubKeys: map[ebiten.Key]bool{ebiten.KeyEscape: true},
			want:     PlayerCommands{Cancel: true},
		},
		{
			name: "all keys pressed",
			stubKeys: map[ebiten.Key]bool{
				ebiten.KeyUp:     true,
				ebiten.KeyDown:   true,
				ebiten.KeyLeft:   true,
				ebiten.KeyRight:  true,
				ebiten.KeyX:      true,
				ebiten.KeySpace:  true,
				ebiten.KeyShift:  true,
				ebiten.KeyEnter:  true,
				ebiten.KeyEscape: true,
			},
			want: PlayerCommands{
				Up:      true,
				Down:    true,
				Left:    true,
				Right:   true,
				Shoot:   true,
				Jump:    true,
				Dash:    true,
				Confirm: true,
				Cancel:  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Stub isKeyPressed to return values from the test case
			oldIsKeyPressed := isKeyPressed
			defer func() { isKeyPressed = oldIsKeyPressed }()

			isKeyPressed = func(key ebiten.Key) bool {
				return tt.stubKeys[key]
			}

			got := ReadPlayerCommands()
			if got != tt.want {
				t.Errorf("ReadPlayerCommands() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
