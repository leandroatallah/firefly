package transition

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestNewBaseTransition(t *testing.T) {
	bt := NewBaseTransition()
	if bt == nil {
		t.Fatal("NewBaseTransition returned nil")
	}
	// Verify it has the expected fields
	if bt.active || bt.starting || bt.exiting || bt.onExitCb != nil {
		t.Error("new BaseTransition should be initialized with zero values")
	}

	// Smoke test methods
	bt.Update()
	bt.Draw(ebiten.NewImage(1, 1))
	bt.StartTransition(func() {})
	bt.EndTransition(func() {})
}
