package gameplayer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/event"
	"github.com/boilerplate/ebiten-template/internal/engine/physics/space"
	_ "github.com/boilerplate/ebiten-template/internal/kit/states"
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

func TestClimberPlayer_Skills(t *testing.T) {
	ctx := &app.AppContext{
		Assets:       os.DirFS("."),
		ActorManager: actors.NewManager(),
		Space:        space.NewSpace(),
		EventManager: event.NewManager(),
	}

	p, err := NewClimberPlayer(ctx)
	if err != nil {
		t.Fatalf("failed to create climber player: %v", err)
	}

	if p == nil {
		t.Fatal("NewClimberPlayer returned nil")
	}
}
