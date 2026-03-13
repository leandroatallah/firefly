package assets

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
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

func createTestFsFromPath(relPath string) fstest.MapFS {
	moduleRoot := getModuleRoot()
	data, err := os.ReadFile(filepath.Join(moduleRoot, relPath))
	if err != nil {
		panic(err)
	}
	return fstest.MapFS{
		"test.png": &fstest.MapFile{
			Data: data,
		},
	}
}

func TestLoadImageFromFs(t *testing.T) {
	cfg := &config.AppConfig{
		ScreenWidth:  320,
		ScreenHeight: 240,
	}

	testFs := createTestFsFromPath("assets/images/default-idle.png")

	ctx := &app.AppContext{
		Config: cfg,
		Assets: testFs,
	}

	img := LoadImageFromFs(ctx, "test.png")

	if img == nil {
		t.Fatalf("LoadImageFromFs returned nil")
	}

	bounds := img.Bounds()
	if bounds.Dx() <= 0 || bounds.Dy() <= 0 {
		t.Fatalf("expected non-zero image bounds, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestLoadImageFromFsNotFound(t *testing.T) {
	// Note: LoadImageFromFs uses log.Fatal on error, which calls os.Exit(1).
	// This test documents the expected behavior: an error is logged when file is not found.
	// In production, this would terminate the application.
	// We skip this test in automated testing since we can't recover from log.Fatal.
	t.Skip("LoadImageFromFs uses log.Fatal which cannot be tested without modifying the function")
}

func TestLoadImageFromFsValidatesImage(t *testing.T) {
	cfg := &config.AppConfig{
		ScreenWidth:  320,
		ScreenHeight: 240,
	}

	testFs := createTestFsFromPath("assets/images/default-idle.png")

	ctx := &app.AppContext{
		Config: cfg,
		Assets: testFs,
	}

	img := LoadImageFromFs(ctx, "test.png")

	if img == nil {
		t.Fatal("expected valid image, got nil")
	}

	// Verify it's a proper ebiten image by checking we can get its size
	w, h := img.Bounds().Dx(), img.Bounds().Bounds().Dy()
	if w <= 0 || h <= 0 {
		t.Fatalf("expected positive image dimensions, got %dx%d", w, h)
	}
}
