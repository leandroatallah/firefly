package tilemap

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
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
	mapWidth := t.Layers[0].Width * t.Tileheight
	mapHeight := t.Layers[0].Height * t.Tileheight
	result := ebiten.NewImage(mapWidth, mapHeight)

	for _, layer := range t.Layers {
		if !layer.Visible {
			continue
		}

		switch {
		case layer.Type == "objectgroup" && layer.Name == "Items":
			// TODO: Spawn items
			t.ParseItems(layer, result)
		case layer.Type == "tilelayer":
			if !t.imageBaseDone {
				t.ParseBase(layer, result)
				t.imageBaseDone = true
			}
		default:
			continue
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

// drawOptsAtPixel returns a DrawImageOptions translated to the destination pixel coordinates.
func drawOptsAtPixel(i int, ts *Tileset) *ebiten.DrawImageOptions {
	x := i % ts.Tilewidth
	y := i / ts.Tileheight
	dx := x * ts.Tilewidth
	dy := y * ts.Tileheight

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(dx), float64(dy))
	return op
}

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
		op := drawOptsAtPixel(i, ts)
		result.DrawImage(tile, op)
	}
}

func (t *Tilemap) ParseItems(layer *Layer, result *ebiten.Image) {
	if layer == nil || result == nil {
		return
	}

	for i, obj := range layer.Objects {
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
		tileImg.Fill(color.RGBA{0, 0, 0xff, 0xff})
		op := drawOptsAtPixel(i, ts)
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
