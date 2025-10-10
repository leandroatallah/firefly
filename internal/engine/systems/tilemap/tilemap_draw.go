package tilemap

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	"io"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
)

func (t *Tilemap) ParseToImage(screen *ebiten.Image) (*ebiten.Image, error) {
	if t == nil {
		return nil, fmt.Errorf("the tilemap was not initialized")
	}
	if len(t.Layers) == 0 || len(t.Tilesets) == 0 {
		return nil, fmt.Errorf("tilemap is not valid")
	}

	// Use the first layer to determine map dimensions. This assumes all layers are the same size.
	mapWidth := t.Layers[0].Width * t.Tilewidth
	mapHeight := t.Layers[0].Height * t.Tileheight
	result := ebiten.NewImage(mapWidth, mapHeight)

	for _, layer := range t.Layers {
		if !layer.Visible {
			continue
		}

		if layer.Type == "tilelayer" {
			t.ParseBase(layer, result)
		}
	}

	// Reset to sync camera
	t.Reset(screen)

	return result, nil
}

// findTileset returns the tileset that applies to the given gid.
// Assumes t.Tilesets are sorted by Firstgid ascending.
func (t *Tilemap) findTileset(gid int) *Tileset {
	if gid <= 0 {
		return nil
	}
	var tileset *Tileset
	for _, ts := range t.Tilesets {
		if gid >= ts.Firstgid {
			tileset = ts
		} else {
			// Since sorted ascending, once gid < Firstgid, we can stop.
			break
		}
	}
	// tileset can be nil if no Firstgid <= gid
	return tileset
}

// tilesetSourceRect computes the source rectangle inside the tileset image for the given gid.
// Returns the rect and the local tile width/height for convenience.
func tilesetSourceRect(ts *Tileset, gid int) image.Rectangle {
	localTileID := gid - ts.Firstgid
	tileX := localTileID % ts.Columns
	tileY := localTileID / ts.Columns
	sx := ts.Margin + tileX*(ts.Tilewidth+ts.Spacing)
	sy := ts.Margin + tileY*(ts.Tileheight+ts.Spacing)
	return image.Rect(sx, sy, sx+ts.Tilewidth, sy+ts.Tileheight)
}

// tilesetSourceID returns the tile ID (zero-based numbering)
func tilesetSourceID(ts *Tileset, gid int) int {
	return gid - ts.Firstgid
}

// drawTileOpts returns a DrawImageOptions translated to the destination pixel coordinates.
func drawTileOpts(i, layerWidth, tileWidth, tileHeight int) *ebiten.DrawImageOptions {
	x := i % layerWidth
	y := i / layerWidth
	dx := x * tileWidth
	dy := y * tileHeight

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(dx), float64(dy))
	return op
}

// TODO: Prevent parse base image more than once. It's unchangeable.
func (t *Tilemap) ParseBase(layer *Layer, result *ebiten.Image) {
	if layer == nil || result == nil {
		return
	}

	for i, tileID := range layer.Data {
		if tileID == 0 {
			continue
		}

		ts := t.findTileset(tileID)
		if ts == nil || ts.EbitenImage == nil {
			continue
		}

		srcRect := tilesetSourceRect(ts, tileID)

		tile := ts.EbitenImage.SubImage(srcRect).(*ebiten.Image)
		op := drawTileOpts(i, layer.Width, ts.Tilewidth, ts.Tileheight)
		result.DrawImage(tile, op)
	}
}

func (t *Tilemap) ParseItems(layer *Layer, result *ebiten.Image) {
	if layer == nil || result == nil {
		return
	}

	for _, obj := range layer.Objects {
		gid := obj.Gid
		if gid == 0 {
			continue
		}

		ts := t.findTileset(gid)
		if ts == nil || ts.EbitenImage == nil {
			continue
		}

		srcRect := tilesetSourceRect(ts, gid)
		tileImg := ts.EbitenImage.SubImage(srcRect).(*ebiten.Image)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(obj.X, obj.Y-obj.Height)
		result.DrawImage(tileImg, op)
	}
}

func (t *Tilemap) Reset(screen *ebiten.Image) {
	t.imageOptions.GeoM.Reset()
	sh := screen.Bounds().Dy()
	mapHeight := t.Layers[0].Height * t.Tileheight
	t.imageOptions.GeoM.Translate(0, float64(sh-mapHeight))
}

func LoadTilemap(path string) (*Tilemap, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var tilemap Tilemap
	if err := json.Unmarshal(byteValue, &tilemap); err != nil {
		return nil, err
	}
	tilemap.imageOptions = &ebiten.DrawImageOptions{}

	// After loading the tilemap structure, load the associated tileset images.
	for _, ts := range tilemap.Tilesets {
		imagePath := filepath.Join(filepath.Dir(path), ts.Image)
		img, err := loadImage(imagePath)
		if err != nil {
			return nil, fmt.Errorf("failed to load tileset image %s: %w", imagePath, err)
		}
		ts.EbitenImage = img
	}

	return &tilemap, nil
}

// loadImage is a helper function to load an image from a file path.
func loadImage(path string) (*ebiten.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return ebiten.NewImageFromImage(img), nil
}

func (t *Tilemap) ImageOptions() *ebiten.DrawImageOptions {
	return t.imageOptions
}
