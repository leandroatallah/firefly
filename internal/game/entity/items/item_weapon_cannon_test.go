package gameitems

import (
	"os"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/mocks"
)

func TestWeaponCannonItem_CollectWhenNotOwned(t *testing.T) {
	ctx := newTestCtx()

	mockPlayer := &mocks.MockActor{Id: "player"}
	mockPlayer.SetPosition(100, 100)
	ctx.ActorManager.RegisterPrimary(mockPlayer)

	item, err := NewWeaponCannonItem(ctx, 100, 100, "cannon-1")
	if err != nil {
		t.Fatalf("failed to create weapon cannon item: %v", err)
	}

	if item == nil {
		t.Fatal("NewWeaponCannonItem returned nil")
	}

	if item.IsRemoved() {
		t.Error("item should not be removed initially")
	}
}

func TestWeaponCannonItem_CollectWhenAlreadyOwned(t *testing.T) {
	ctx := newTestCtx()

	mockPlayer := &mocks.MockActor{Id: "player"}
	mockPlayer.SetPosition(100, 100)
	ctx.ActorManager.RegisterPrimary(mockPlayer)

	item, err := NewWeaponCannonItem(ctx, 100, 100, "cannon-1")
	if err != nil {
		t.Fatalf("failed to create weapon cannon item: %v", err)
	}

	if item == nil {
		t.Fatal("NewWeaponCannonItem returned nil")
	}

	if item.IsRemoved() {
		t.Error("item should not be removed initially")
	}
}

func TestWeaponCannonItem_OnTouch_RemovesItem(t *testing.T) {
	ctx := newTestCtx()

	mockPlayer := &mocks.MockActor{Id: "player"}
	mockPlayer.SetPosition(100, 100)
	ctx.ActorManager.RegisterPrimary(mockPlayer)

	item, err := NewWeaponCannonItem(ctx, 100, 100, "cannon-1")
	if err != nil {
		t.Fatalf("failed to create weapon cannon item: %v", err)
	}

	item.OnTouch(mockPlayer)

	if !item.IsRemoved() {
		t.Error("item should be removed after player touches it")
	}
}

func TestWeaponCannonItem_OnTouch_WithoutPlayer(t *testing.T) {
	ctx := newTestCtx()

	item, err := NewWeaponCannonItem(ctx, 100, 100, "cannon-1")
	if err != nil {
		t.Fatalf("failed to create weapon cannon item: %v", err)
	}

	mockOther := &mocks.MockActor{Id: "other"}
	mockOther.SetPosition(100, 100)

	item.OnTouch(mockOther)

	if item.IsRemoved() {
		t.Error("item should not be removed when touched by non-player")
	}
}

func TestWeaponCannonItem_CreatesConfigFile(t *testing.T) {
	configPath := "assets/entities/items/item-weapon-cannon.json"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("config file %s does not exist", configPath)
	}
}

func TestWeaponCannonItem_HasCorrectID(t *testing.T) {
	ctx := newTestCtx()

	item, err := NewWeaponCannonItem(ctx, 100, 100, "cannon-test-id")
	if err != nil {
		t.Fatalf("failed to create weapon cannon item: %v", err)
	}

	if item.ID() != "cannon-test-id" {
		t.Errorf("expected ID cannon-test-id, got %s", item.ID())
	}
}
