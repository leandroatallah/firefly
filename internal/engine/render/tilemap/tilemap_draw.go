package tilemap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
)

func (t *Tilemap) ParseToImage(screen *ebiten.Image) (*ebiten.Image, error) {
	if _, err := t.isTilemapValid(); err != nil {
		return nil, err
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

// applyFlips applies Tiled's flipping/rotation to the DrawImageOptions.
// w and h are the tile width and height.
func applyFlips(op *ebiten.DrawImageOptions, h, v, d bool, tileW, tileH float64) {
	if d {
		// Diagonal flip: Swap x and y axis
		op.GeoM.SetElement(0, 0, 0)
		op.GeoM.SetElement(0, 1, 1)
		op.GeoM.SetElement(1, 0, 1)
		op.GeoM.SetElement(1, 1, 0)
	}

	if h {
		op.GeoM.Scale(-1, 1)
		// Translate back to positive coordinates
		// For standard tiles, let's use the dimension along the X axis.
		wd := tileW
		if d {
			wd = tileH
		}
		op.GeoM.Translate(wd, 0)
	}

	if v {
		op.GeoM.Scale(1, -1)
		hd := tileH
		if d {
			hd = tileW
		}
		op.GeoM.Translate(0, hd)
	}
}

func (t *Tilemap) ParseBase(layer *Layer, result *ebiten.Image) {
	if layer == nil || result == nil {
		return
	}

	for i, rawID := range layer.Data {
		if rawID == 0 {
			continue
		}

		tileID, flipH, flipV, flipD := extractGIDAndFlags(rawID)
		if tileID == 0 {
			continue
		}

		ts := t.findTileset(tileID)
		if ts == nil || ts.EbitenImage == nil {
			continue
		}

		srcRect := tilesetSourceRect(ts, tileID)

		tile := ts.EbitenImage.SubImage(srcRect).(*ebiten.Image)

		op := &ebiten.DrawImageOptions{}
		applyFlips(op, flipH, flipV, flipD, float64(ts.Tilewidth), float64(ts.Tileheight))

		// Then translate to position
		x := i % layer.Width
		y := i / layer.Width
		dx := float64(x * ts.Tilewidth)
		dy := float64(y * ts.Tileheight)
		op.GeoM.Translate(dx, dy)

		result.DrawImage(tile, op)
	}
}

func (t *Tilemap) ParseItems(layer *Layer, result *ebiten.Image) {
	if layer == nil || result == nil || !layer.Visible {
		return
	}

	for _, obj := range layer.Objects {
		rawID := obj.Gid
		if rawID == 0 {
			continue
		}

		gid, flipH, flipV, flipD := extractGIDAndFlags(rawID)
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
		applyFlips(op, flipH, flipV, flipD, float64(ts.Tilewidth), float64(ts.Tileheight))
		op.GeoM.Translate(obj.X, obj.Y-obj.Height)

		result.DrawImage(tileImg, op)
	}
}

func (t *Tilemap) Reset(screen *ebiten.Image) {
	t.imageOptions.GeoM.Reset()
}

func LoadTilemap(fsys fs.FS, path string) (*Tilemap, error) {
	data, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, err
	}

	var tilemap Tilemap
	if err := json.Unmarshal(data, &tilemap); err != nil {
		return nil, err
	}
	tilemap.imageOptions = &ebiten.DrawImageOptions{}

	// After loading the tilemap structure, load the associated tileset images.
	for _, ts := range tilemap.Tilesets {
		imagePath := filepath.Join(filepath.Dir(path), ts.Image)
		img, err := loadImage(fsys, imagePath)
		if err != nil {
			return nil, fmt.Errorf("failed to load tileset image %s: %w", imagePath, err)
		}
		ts.EbitenImage = img
	}

	return &tilemap, nil
}

// loadImage is a helper function to load an image from a file path.
func loadImage(fsys fs.FS, path string) (*ebiten.Image, error) {
	data, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return ebiten.NewImageFromImage(img), nil
}

func (t *Tilemap) isTilemapValid() (bool, error) {
	if t == nil {
		return false, fmt.Errorf("the tilemap was not initialized")
	}
	if len(t.Layers) == 0 || len(t.Tilesets) == 0 {
		return false, fmt.Errorf("tilemap is not valid")
	}

	return true, nil
}

func (t *Tilemap) FindLayerByName(name string) (*Layer, bool) {
	if valid, err := t.isTilemapValid(); !valid {
		log.Printf("tilemap is not valid: %v", err)
		return nil, false
	}

	for _, layer := range t.Layers {
		if !layer.Visible {
			continue
		}

		if layer.Name == name {
			return layer, true
		}
	}

	return nil, false
}
