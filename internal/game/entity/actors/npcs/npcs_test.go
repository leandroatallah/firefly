package gamenpcs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
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

func TestNewPrincess(t *testing.T) {
	ctx := &app.AppContext{
		ActorManager: actors.NewManager(),
	}

	s, err := NewPrincess(ctx, 100, 100, "princess-1")
	if err != nil {
		t.Fatalf("failed to create princess: %v", err)
	}

	if s == nil {
		t.Fatal("NewPrincess returned nil")
	}

	if s.ID() != "princess-1" {
		t.Errorf("expected ID princess-1, got %s", s.ID())
	}
}
