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

type Manager struct {
	// TODO: Implement this
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) Update() {
	// TODO: Implement this
}

