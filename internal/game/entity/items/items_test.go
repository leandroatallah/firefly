package gameitems

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/entity/items"
	"github.com/leandroatallah/firefly/internal/engine/mocks"
)

type MockCoinCollectorActor struct {
	mocks.MockActor
	coins int
}

func (m *MockCoinCollectorActor) AddCoinCount(amount int) {
	m.coins += amount
}

func (m *MockCoinCollectorActor) CoinCount() int {
	return m.coins
}

func getModuleRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			panic("could not find go.mod")
		}
		dir = parent
	}
}

func TestMain(m *testing.M) {
	err := os.Chdir(getModuleRoot())
	if err != nil {
		panic(err)
	}

	cfg := &config.AppConfig{
		ScreenWidth:  320,
		ScreenHeight: 224,
	}
	config.Set(cfg)

	os.Exit(m.Run())
}

func TestNewCollectibleCoinItem(t *testing.T) {
	ctx := &app.AppContext{
		ActorManager: actors.NewManager(),
	}
	
	coin, err := NewCollectibleCoinItem(ctx, 100, 100, "coin-1")
	if err != nil {
		t.Fatalf("failed to create coin: %v", err)
	}

	if coin == nil {
		t.Fatal("NewCollectibleCoinItem returned nil")
	}

	if coin.ID() != "coin-1" {
		t.Errorf("expected ID coin-1, got %s", coin.ID())
	}
}

func TestCollectibleCoinItem_OnTouch(t *testing.T) {
	actorMgr := actors.NewManager()
	ctx := &app.AppContext{
		ActorManager: actorMgr,
	}
	
	coin, err := NewCollectibleCoinItem(ctx, 100, 100, "coin-1")
	if err != nil {
		t.Fatalf("failed to create coin: %v", err)
	}

	mockPlayer := &MockCoinCollectorActor{
		MockActor: mocks.MockActor{Id: "player"},
	}
	actorMgr.Register(mockPlayer)
	
	// Touch with non-player should do nothing
	mockOther := &mocks.MockActor{Id: "other"}
	coin.OnTouch(mockOther)
	if coin.IsRemoved() {
		t.Error("coin should not be removed by non-player")
	}

	// Touch with player should remove coin
	coin.OnTouch(mockPlayer)
	if !coin.IsRemoved() {
		t.Error("coin should be removed by player")
	}
	if mockPlayer.coins != 1 {
		t.Errorf("expected 1 coin, got %d", mockPlayer.coins)
	}
}

func TestNewFallingPlatformItem(t *testing.T) {
	ctx := &app.AppContext{
		ActorManager: actors.NewManager(),
	}
	
	fp, err := NewFallingPlatformItem(ctx, 100, 100, "fp-1")
	if err != nil {
		t.Fatalf("failed to create falling platform: %v", err)
	}

	if fp == nil {
		t.Fatal("NewFallingPlatformItem returned nil")
	}

	if fp.ID() != "fp-1" {
		t.Errorf("expected ID fp-1, got %s", fp.ID())
	}

	// Test OnTouch
	mockPlayer := &mocks.MockActor{Id: "player"}
	mockPlayer.SetPosition(100, 80) // Above platform
	ctx.ActorManager.Register(mockPlayer)
	
	fp.OnTouch(mockPlayer)
}

func TestNewPowerUpItems(t *testing.T) {
	ctx := &app.AppContext{
		ActorManager: actors.NewManager(),
	}

	tests := []struct {
		name    string
		factory func(ctx *app.AppContext, x, y int, id string) (items.Item, error)
	}{
		{"Freeze", func(ctx *app.AppContext, x, y int, id string) (items.Item, error) {
			item, err := NewFreezePowerItem(ctx, x, y, id)
			if err != nil {
				return nil, err
			}
			return item, nil
		}},
		{"Grow", func(ctx *app.AppContext, x, y int, id string) (items.Item, error) {
			item, err := NewGrowPowerItem(ctx, x, y, id)
			if err != nil {
				return nil, err
			}
			return item, nil
		}},
		{"Star", func(ctx *app.AppContext, x, y int, id string) (items.Item, error) {
			item, err := NewStarPowerItem(ctx, x, y, id)
			if err != nil {
				return nil, err
			}
			return item, nil
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := tt.name + "-1"
			item, err := tt.factory(ctx, 100, 100, id)
			if err != nil {
				t.Fatalf("failed to create %s power-up: %v", tt.name, err)
			}
			if item == nil {
				t.Fatalf("New%sPowerItem returned nil", tt.name)
			}
			if item.ID() != id {
				t.Errorf("expected ID %s, got %s", id, item.ID())
			}
		})
	}
}

func TestInitItemMap(t *testing.T) {
	ctx := &app.AppContext{}
	m := InitItemMap(ctx)
	if m == nil {
		t.Fatal("InitItemMap returned nil")
	}
	if _, ok := m[FreezePowerUpType]; !ok {
		t.Error("FreezePowerUpType missing from ItemMap")
	}
	if _, ok := m[GrowPowerUpType]; !ok {
		t.Error("GrowPowerUpType missing from ItemMap")
	}
	if _, ok := m[StarPowerUpType]; !ok {
		t.Error("StarPowerUpType missing from ItemMap")
	}
}
