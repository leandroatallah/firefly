package gamescene

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/mocks"
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
		MainFontFace:  "assets/fonts/pressstart2p.ttf",
		SmallFontFace: "assets/fonts/tiny5.ttf",
	}
	config.Set(cfg)

	os.Exit(m.Run())
}

func TestIntroScene(t *testing.T) {
	mockNav := &mocks.MockSceneManager{}
	ctx := &app.AppContext{
		SceneManager: mockNav,
	}
	
	s := NewIntroScene(ctx)
	if s == nil {
		t.Fatal("NewIntroScene returned nil")
	}

	s.OnStart()
	
	screen := ebiten.NewImage(320, 240)
	s.Draw(screen)
	
	for i := 0; i < 100; i++ {
		s.Update()
	}
	
	s.NextScene()
	s.OnFinish()
}

func TestInitSceneMap(t *testing.T) {
	ctx := &app.AppContext{}
	m := InitSceneMap(ctx)
	if len(m) == 0 {
		t.Error("InitSceneMap returned empty map")
	}
}
