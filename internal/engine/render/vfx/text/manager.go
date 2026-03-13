package text

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
	"github.com/leandroatallah/firefly/internal/engine/render/camera"
)

type Manager struct {
	texts       []FloatingText
	defaultFont *font.FontText
}

func NewManager() *Manager {
	return &Manager{
		texts: make([]FloatingText, 0),
	}
}

func (m *Manager) SetDefaultFont(f *font.FontText) {
	m.defaultFont = f
}

func (m *Manager) Add(ft FloatingText) {
	if ft != nil {
		ft.SetFont(m.defaultFont)
		m.texts = append(m.texts, ft)
	}
}

func (m *Manager) Update() {
	active := make([]FloatingText, 0, len(m.texts))
	for _, ft := range m.texts {
		ft.Update()
		if !ft.IsComplete() {
			active = append(active, ft)
		}
	}
	m.texts = active
}

func (m *Manager) Draw(screen *ebiten.Image, cam *camera.Controller) {
	for _, ft := range m.texts {
		ft.Draw(screen, cam)
	}
}
