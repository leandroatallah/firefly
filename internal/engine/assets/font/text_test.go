package font

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// getModuleRoot returns the absolute path to the module root
func getModuleRoot() string {
	// Find go.mod by walking up the directory tree
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

func TestNewFontText(t *testing.T) {
	moduleRoot := getModuleRoot()
	font, err := NewFontText(filepath.Join(moduleRoot, "assets/fonts/tiny5.ttf"))
	if err != nil {
		t.Fatalf("NewFontText failed: %v", err)
	}

	if font == nil {
		t.Fatal("NewFontText returned nil")
	}

	if font.source == nil {
		t.Fatal("NewFontText created FontText with nil source")
	}
}

func TestNewFontTextNotFound(t *testing.T) {
	font, err := NewFontText("nonexistent.ttf")
	if err == nil {
		t.Fatal("expected error for nonexistent font, got nil")
	}
	if font != nil {
		t.Fatal("expected nil font on error, got non-nil")
	}
}

func TestFontTextNewFace(t *testing.T) {
	moduleRoot := getModuleRoot()
	font, err := NewFontText(filepath.Join(moduleRoot, "assets/fonts/tiny5.ttf"))
	if err != nil {
		t.Fatalf("NewFontText failed: %v", err)
	}

	face := font.NewFace(16.0)
	if face == nil {
		t.Fatal("NewFace returned nil")
	}

	if face.Size != 16.0 {
		t.Fatalf("expected face size 16.0, got %f", face.Size)
	}

	if face.Source != font.source {
		t.Fatal("face source doesn't match font source")
	}
}

func TestFontTextNewFaceDifferentSizes(t *testing.T) {
	moduleRoot := getModuleRoot()
	font, err := NewFontText(filepath.Join(moduleRoot, "assets/fonts/tiny5.ttf"))
	if err != nil {
		t.Fatalf("NewFontText failed: %v", err)
	}

	sizes := []float64{8.0, 12.0, 16.0, 24.0, 32.0}
	for _, size := range sizes {
		face := font.NewFace(size)
		if face == nil {
			t.Fatalf("NewFace(%f) returned nil", size)
		}
		if face.Size != size {
			t.Fatalf("expected face size %f, got %f", size, face.Size)
		}
	}
}

func TestFontTextDraw(t *testing.T) {
	moduleRoot := getModuleRoot()
	font, err := NewFontText(filepath.Join(moduleRoot, "assets/fonts/tiny5.ttf"))
	if err != nil {
		t.Fatalf("NewFontText failed: %v", err)
	}

	screen := ebiten.NewImage(100, 50)
	op := &text.DrawOptions{}

	// Should not panic
	font.Draw(screen, "Hello", 16.0, op)
}

func TestFontTextDrawWithNilSource(t *testing.T) {
	font := &FontText{source: nil}

	screen := ebiten.NewImage(100, 50)
	op := &text.DrawOptions{}

	// Should not panic when source is nil
	font.Draw(screen, "Hello", 16.0, op)
}

func TestFontTextDrawEmptyString(t *testing.T) {
	moduleRoot := getModuleRoot()
	font, err := NewFontText(filepath.Join(moduleRoot, "assets/fonts/tiny5.ttf"))
	if err != nil {
		t.Fatalf("NewFontText failed: %v", err)
	}

	screen := ebiten.NewImage(100, 50)
	op := &text.DrawOptions{}

	// Should not panic with empty string
	font.Draw(screen, "", 16.0, op)
}
