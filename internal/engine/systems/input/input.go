package input

import "github.com/hajimehoshi/ebiten/v2"

func IsSomeKeyPressed(keys ...ebiten.Key) bool {
	for _, k := range keys {
		if ebiten.IsKeyPressed(k) {
			return true
		}
	}
	return false
}
