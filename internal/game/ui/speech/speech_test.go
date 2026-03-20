package gamespeech

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/data/i18n"
	"github.com/leandroatallah/firefly/internal/engine/ui/speech"
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

func TestSpeechBubble(t *testing.T) {
	fontMain, err := font.NewFontText(config.Get().MainFontFace)
	if err != nil {
		t.Fatalf("failed to load font: %v", err)
	}
	speechFont := speech.NewSpeechFont(fontMain, 8, 14)
	
	i18nManager := i18n.NewI18nManager(os.DirFS("."))
	sb := NewSpeechBubble(speechFont, i18nManager)
	if sb == nil {
		t.Fatal("NewSpeechBubble returned nil")
	}

	sb.Show()
	if !sb.Visible() {
		t.Error("expected visible after Show")
	}

	sb.Update()
	
	screen := ebiten.NewImage(320, 240)
	sb.Draw(screen, "Hello World")

	sb.Hide()
	sb.Update()
	
	sb.ResetText()
}

func TestStorySpeech(t *testing.T) {
	fontSmall, err := font.NewFontText(config.Get().SmallFontFace)
	if err != nil {
		t.Fatalf("failed to load font: %v", err)
	}
	speechFont := speech.NewSpeechFont(fontSmall, 8, 12)
	
	i18nManager := i18n.NewI18nManager(os.DirFS("."))
	ss := NewStorySpeech(speechFont, i18nManager)
	if ss == nil {
		t.Fatal("NewStorySpeech returned nil")
	}

	ss.Show()
	ss.Update()
	
	screen := ebiten.NewImage(320, 240)
	ss.Draw(screen, "Once upon a time...")
	
	ss.Hide()
}
