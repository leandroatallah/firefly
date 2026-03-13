package imagemanager

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestNewImageManager(t *testing.T) {
	im := NewImageManager()
	if im == nil {
		t.Fatal("NewImageManager returned nil")
	}
	if im.images == nil {
		t.Fatal("NewImageManager created with nil images map")
	}
	if len(im.images) != 0 {
		t.Fatalf("expected empty images map, got %d items", len(im.images))
	}
}

func TestImageManagerAddAndGet(t *testing.T) {
	im := NewImageManager()

	testImg := ebiten.NewImage(32, 32)
	im.Add("test", testImg)

	got := im.GetImage("test")
	if got == nil {
		t.Fatal("GetImage returned nil for existing image")
	}

	if got != testImg {
		t.Fatal("GetImage returned different image than added")
	}
}

func TestImageManagerGetNotFound(t *testing.T) {
	im := NewImageManager()

	got := im.GetImage("nonexistent")
	if got != nil {
		t.Fatalf("expected nil for nonexistent image, got %v", got)
	}
}

func TestImageManagerAddMultiple(t *testing.T) {
	im := NewImageManager()

	img1 := ebiten.NewImage(16, 16)
	img2 := ebiten.NewImage(32, 32)
	img3 := ebiten.NewImage(64, 64)

	im.Add("small", img1)
	im.Add("medium", img2)
	im.Add("large", img3)

	if len(im.images) != 3 {
		t.Fatalf("expected 3 images, got %d", len(im.images))
	}

	if im.GetImage("small") != img1 {
		t.Fatal("small image mismatch")
	}
	if im.GetImage("medium") != img2 {
		t.Fatal("medium image mismatch")
	}
	if im.GetImage("large") != img3 {
		t.Fatal("large image mismatch")
	}
}

func TestImageManagerAddOverwrite(t *testing.T) {
	im := NewImageManager()

	img1 := ebiten.NewImage(16, 16)
	img2 := ebiten.NewImage(32, 32)

	im.Add("test", img1)
	if im.GetImage("test") != img1 {
		t.Fatal("first image mismatch")
	}

	im.Add("test", img2)
	if im.GetImage("test") != img2 {
		t.Fatal("overwritten image mismatch")
	}

	if len(im.images) != 1 {
		t.Fatalf("expected 1 image after overwrite, got %d", len(im.images))
	}
}

func TestImageItem(t *testing.T) {
	testImg := ebiten.NewImage(24, 24)
	item := &ImageItem{
		name:  "test-item",
		image: testImg,
	}

	if item.Name() != "test-item" {
		t.Fatalf("expected name 'test-item', got '%s'", item.Name())
	}

	if item.Image() != testImg {
		t.Fatal("Image() returned different image than set")
	}
}

func TestImageItemNilImage(t *testing.T) {
	item := &ImageItem{
		name:  "nil-item",
		image: nil,
	}

	if item.Name() != "nil-item" {
		t.Fatalf("expected name 'nil-item', got '%s'", item.Name())
	}

	if item.Image() != nil {
		t.Fatal("expected nil image, got non-nil")
	}
}
