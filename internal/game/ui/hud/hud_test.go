package gamehud

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	gamesetup "github.com/boilerplate/ebiten-template/internal/game/app"
	"github.com/hajimehoshi/ebiten/v2"
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
	// Change working directory to module root so assets can be found
	err := os.Chdir(getModuleRoot())
	if err != nil {
		panic(err)
	}

	// Initialize config
	cfg := gamesetup.NewConfig()
	config.Set(cfg)

	os.Exit(m.Run())
}

func heartFS(t *testing.T) fs.FS {
	t.Helper()
	data, err := os.ReadFile("assets/images/item-power-grow.png")
	if err != nil {
		t.Fatalf("could not read test image: %v", err)
	}
	return fstest.MapFS{
		"assets/images/heart.png": &fstest.MapFile{Data: data},
	}
}

func TestNewStatusBar(t *testing.T) {
	sb, err := NewStatusBar(nil, 100, nil, heartFS(t))
	if err != nil {
		t.Fatalf("failed to create StatusBar: %v", err)
	}

	if sb == nil {
		t.Fatal("NewStatusBar returned nil")
	}

	if sb.score != 100 {
		t.Errorf("expected score 100, got %d", sb.score)
	}
}

func TestStatusBar_Update(t *testing.T) {
	sb, err := NewStatusBar(nil, 100, nil, heartFS(t))
	if err != nil {
		t.Fatalf("failed to create StatusBar: %v", err)
	}

	if err := sb.Update(); err != nil {
		t.Errorf("Update returned error: %v", err)
	}
}

func TestStatusBar_Draw(t *testing.T) {
	sb, err := NewStatusBar(nil, 100, nil, heartFS(t))
	if err != nil {
		t.Fatalf("failed to create StatusBar: %v", err)
	}

	screen := ebiten.NewImage(320, 240)
	sb.player = nil
	sb.Draw(screen)
}
