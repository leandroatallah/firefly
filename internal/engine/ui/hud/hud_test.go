package hud

import (
	"errors"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestBaseHUD(t *testing.T) {
	var h BaseHUD

	if h.Visible() {
		t.Error("expected visible to be false by default")
	}

	h.SetVisible(true)
	if !h.Visible() {
		t.Error("expected visible to be true after SetVisible(true)")
	}

	err := h.Update()
	if err != nil {
		t.Errorf("unexpected error from Update: %v", err)
	}
}

type mockHUD struct {
	BaseHUD
	updateCalled bool
	drawCalled   bool
}

func (m *mockHUD) Update() error {
	m.updateCalled = true
	return nil
}

func (m *mockHUD) Draw(screen *ebiten.Image) {
	m.drawCalled = true
}

func TestManager(t *testing.T) {
	m1 := &mockHUD{}
	m2 := &mockHUD{}
	m2.SetVisible(true)

	mgr := NewManager(m1)
	mgr.Add(m2)

	err := mgr.Update()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !m1.updateCalled || !m2.updateCalled {
		t.Error("expected Update to be called on all elements")
	}

	mgr.Draw(nil)

	if m1.drawCalled {
		t.Error("expected Draw NOT to be called on invisible element")
	}
	if !m2.drawCalled {
		t.Error("expected Draw to be called on visible element")
	}
}

type errorHUD struct {
	BaseHUD
}

func (e *errorHUD) Update() error {
	return errors.New("test error")
}

func (e *errorHUD) Draw(screen *ebiten.Image) {}

func TestManager_UpdateError(t *testing.T) {
	mgr := NewManager(&errorHUD{})
	if err := mgr.Update(); err == nil {
		t.Error("expected error from Update")
	}
}
