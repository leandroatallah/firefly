package input

import "github.com/hajimehoshi/ebiten/v2"

// Swappable function var: allows injection in tests
//
//nolint:gochecknoglobals
var isKeyPressed = ebiten.IsKeyPressed

func IsSomeKeyPressed(keys ...ebiten.Key) bool {
	for _, k := range keys {
		if isKeyPressed(k) {
			return true
		}
	}
	return false
}
