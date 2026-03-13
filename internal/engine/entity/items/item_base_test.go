package items

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
	"github.com/leandroatallah/firefly/internal/engine/render/sprites"
)

func TestBaseItem_Lifecycle(t *testing.T) {
	img := ebiten.NewImage(1, 1)
	sMap := sprites.SpriteMap{Idle: &sprites.Sprite{Image: img}}
	rect := bodyphysics.NewRect(0, 0, 8, 8)

	item := NewBaseItem("coin_01", sMap, rect)

	if item.ID() != "coin_01" {
		t.Errorf("expected ID coin_01, got %s", item.ID())
	}

	if item.IsRemoved() {
		t.Error("item should not be removed initially")
	}

	item.SetRemoved(true)
	if !item.IsRemoved() {
		t.Error("item should be marked as removed")
	}

	if item.State() != Idle {
		t.Errorf("expected state Idle, got %v", item.State())
	}
}
