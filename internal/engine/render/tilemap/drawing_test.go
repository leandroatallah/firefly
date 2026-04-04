package tilemap

import (
	"image"
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestTilemapDrawing(t *testing.T) {
	// Create a mock tileset image (32x32, 2x2 tiles of 16x16)
	tsImg := ebiten.NewImage(32, 32)
	tsImg.Fill(color.White)

	ts := &Tileset{
		Firstgid:    1,
		Columns:     2,
		Tilewidth:   16,
		Tileheight:  16,
		EbitenImage: tsImg,
	}

	tm := &Tilemap{
		Width:      2,
		Height:     2,
		Tilewidth:  16,
		Tileheight: 16,
		Tilesets:   []*Tileset{ts},
		Layers: []*Layer{
			{
				Name:    "Base",
				Visible: true,
				Type:    "tilelayer",
				Width:   2,
				Height:  2,
				Data:    []int{1, 2, 3, 4},
			},
			{
				Name:    "Items",
				Visible: true,
				Type:    "objectgroup",
				Objects: []*Obstacle{
					{Gid: 1, X: 0, Y: 16, Width: 16, Height: 16},
				},
			},
		},
	}
	tm.imageOptions = &ebiten.DrawImageOptions{}

	screen := ebiten.NewImage(320, 240)

	// Test Image (this calls ParseToImage and sets tm.image)
	img, err := tm.Image(screen)
	if err != nil {
		t.Fatalf("Image failed: %v", err)
	}
	if img.Bounds().Dx() != 32 || img.Bounds().Dy() != 32 {
		t.Errorf("expected 32x32 image, got %dx%d", img.Bounds().Dx(), img.Bounds().Dy())
	}

	// Test caching
	img2, _ := tm.Image(screen)
	if img2 != img {
		t.Error("expected cached image")
	}

	// Test findTileset
	if tm.findTileset(0) != nil {
		t.Error("expected nil for gid 0")
	}
	if tm.findTileset(1) != ts {
		t.Error("expected ts for gid 1")
	}

	// Test tilesetSourceRect
	rect := tilesetSourceRect(ts, 1)
	if rect != image.Rect(0, 0, 16, 16) {
		t.Errorf("unexpected rect for gid 1: %v", rect)
	}

	// Test applyFlips (smoke test)
	op := &ebiten.DrawImageOptions{}
	applyFlips(op, true, true, true, 16, 16)

	// Test ParseItems
	tm.ParseItems(tm.Layers[1], img)

	// Test ImageOptions
	if tm.ImageOptions() == nil {
		t.Error("expected non-nil ImageOptions")
	}

	// Test Reset
	tm.Reset(screen)

	// Test isTilemapValid
	valid, _ := tm.isTilemapValid()
	if !valid {
		t.Error("expected tilemap to be valid")
	}
}

func TestTilemapHelpers(t *testing.T) {
	tm := &Tilemap{
		Layers: []*Layer{
			{Name: "Camera", Visible: true, Objects: []*Obstacle{{X: 10, Y: 20}}},
			{Name: "Hidden", Visible: false},
			{Name: "EmptyCamera", Visible: true, Type: "objectgroup", Objects: []*Obstacle{}},
		},
		Tilesets: []*Tileset{{Firstgid: 1}},
	}

	if !tm.HasCameraStartPosition() {
		t.Error("expected HasCameraStartPosition to be true")
	}
	x, y, found := tm.GetCameraStartPosition()
	if !found || x != 10 || y != 20 {
		t.Errorf("GetCameraStartPosition failed: %d, %d, %v", x, y, found)
	}

	if tm.HasPlayerStartPosition() {
		t.Error("expected HasPlayerStartPosition to be false")
	}
	_, _, found = tm.GetPlayerStartPosition()
	if found {
		t.Error("expected player start to not be found")
	}

	// Test empty camera layer
	tmEmptyCam := &Tilemap{
		Layers: []*Layer{{Name: "Camera", Visible: true, Objects: []*Obstacle{}}},
	}
	_, _, found = tmEmptyCam.GetCameraStartPosition()
	if found {
		t.Error("expected camera start to not be found in empty layer")
	}

	layer, ok := tm.FindLayerByName("Camera")
	if !ok || layer == nil {
		t.Error("FindLayerByName failed")
	}

	_, ok = tm.FindLayerByName("Hidden")
	if ok {
		t.Error("expected Hidden layer to be not found (invisible)")
	}

	// Nil tilemap checks
	var tmNil *Tilemap
	if tmNil.HasPlayerStartPosition() {
		t.Error("expected false for nil tilemap")
	}
	if tmNil.HasCameraStartPosition() {
		t.Error("expected false for nil tilemap")
	}
	_, _, found = tmNil.GetPlayerStartPosition()
	if found {
		t.Error("expected found to be false for nil tilemap")
	}
	_, _, found = tmNil.GetCameraStartPosition()
	if found {
		t.Error("expected found to be false for nil tilemap")
	}
	if tmNil.GetItemsPositionID() != nil {
		t.Error("expected nil items for nil tilemap")
	}
	if tmNil.GetEnemiesPositionID() != nil {
		t.Error("expected nil enemies for nil tilemap")
	}
	if tmNil.GetNpcsPositionID() != nil {
		t.Error("expected nil npcs for nil tilemap")
	}

	valid, _ := tmNil.isTilemapValid()
	if valid {
		t.Error("expected nil tilemap to be invalid")
	}
}
