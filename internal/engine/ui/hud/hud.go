package hud

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type HUD interface {
	Draw(screen *ebiten.Image)
	Update() error
	Visible() bool
	SetVisible(v bool)
}

// BaseHUD provides a common implementation for HUD elements.
type BaseHUD struct {
	visible bool
}

func (h *BaseHUD) Visible() bool {
	return h.visible
}

func (h *BaseHUD) SetVisible(v bool) {
	h.visible = v
}

func (h *BaseHUD) Update() error {
	return nil
}

// Manager handles multiple HUD elements.
type Manager struct {
	elements []HUD
}

func NewManager(elements ...HUD) *Manager {
	return &Manager{
		elements: elements,
	}
}

func (m *Manager) Add(h HUD) {
	m.elements = append(m.elements, h)
}

func (m *Manager) Update() error {
	for _, h := range m.elements {
		if err := h.Update(); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) Draw(screen *ebiten.Image) {
	for _, h := range m.elements {
		if h.Visible() {
			h.Draw(screen)
		}
	}
}
