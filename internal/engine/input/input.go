package input

import "github.com/hajimehoshi/ebiten/v2"

var isKeyPressed = ebiten.IsKeyPressed

func IsSomeKeyPressed(keys ...ebiten.Key) bool {
	for _, k := range keys {
		if isKeyPressed(k) {
			return true
		}
	}
	return false
}
