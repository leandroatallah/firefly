package gameitems

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/mocks"
)

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

func newTestCtx() *app.AppContext {
	return &app.AppContext{
		ActorManager: actors.NewManager(),
		Assets:       os.DirFS("."),
	}
}

func TestNewFallingPlatformItem(t *testing.T) {
	ctx := newTestCtx()

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
	mockPlayer.SetPosition(100, 80)
	ctx.ActorManager.Register(mockPlayer)

	fp.OnTouch(mockPlayer)
}

func TestInitItemMap(t *testing.T) {
	ctx := &app.AppContext{}
	m := InitItemMap(ctx)
	if m == nil {
		t.Fatal("InitItemMap returned nil")
	}
	if _, ok := m[FallingPlatformType]; !ok {
		t.Error("FallingPlatformType missing from ItemMap")
	}
}
