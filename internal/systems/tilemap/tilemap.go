package tilemap

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	"io"
	"math"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/systems/physics"
)

type Tilemap struct {
	Height       int        `json:"height"`
	Infinite     bool       `json:"infinite"`
	Layers       []*Layer   `json:"layers"`
	Tileheight   int        `json:"tileheight"`
	Tilesets     []*Tileset `json:"tilesets"`
	imageOptions *ebiten.DrawImageOptions
}

type Layer struct {
	Data    []int       `json:"data"`
	Height  int         `json:"height"`
	Id      int         `json:"id"`
	Name    string      `json:"name"`
	Opacity int         `json:"opacity"`
	Type    string      `json:"type"`
	Visible bool        `json:"visible"`
	Width   int         `json:"width"`
	X       int         `json:"x"`
	Y       int         `json:"y"`
	Objects []*Obstacle `json:"objects"`
}

type Obstacle struct {
	Height   float64 `json:"height"`
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Rotation float64 `json:"rotation"`
	Type     string  `json:"type"`
	Visible  bool    `json:"visible"`
	Width    float64 `json:"width"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
}

type Tileset struct {
	Columns          int           `json:"columns"`
	Firstgid         int           `json:"firstgid"`
	Image            string        `json:"image"`
	Imageheight      int           `json:"imageheight"`
	Imagewidth       int           `json:"imagewidth"`
	Margin           int           `json:"margin"`
	Name             string        `json:"name"`
	Spacing          int           `json:"spacing"`
	Tilecount        int           `json:"tilecount"`
	Tileheight       int           `json:"tileheight"`
	Tilewidth        int           `json:"tilewidth"`
	Transparentcolor string        `json:"transparentcolor"`
	EbitenImage      *ebiten.Image `json:"-"`
}

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
		if layer.Type != "tilelayer" || !layer.Visible {
			continue
		}

		for i, tileID := range layer.Data {
			if tileID == 0 {
				continue
			}

			// Find the correct tileset for this tile ID.
			// This assumes tilesets are sorted by Firstgid in the JSON.
			var tileset *Tileset
			for _, ts := range t.Tilesets {
				if tileID >= ts.Firstgid {
					tileset = ts
				}
			}

			if tileset == nil || tileset.EbitenImage == nil {
				continue // No valid tileset found for this ID.
			}

			localTileID := tileID - tileset.Firstgid

			tileX := localTileID % tileset.Columns
			tileY := localTileID / tileset.Columns
			sx := tileset.Margin + tileX*(tileset.Tilewidth+tileset.Spacing)
			sy := tileset.Margin + tileY*(tileset.Tileheight+tileset.Spacing)

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64((i%layer.Width)*tileset.Tilewidth), float64((i/layer.Width)*tileset.Tileheight))

			tileRect := image.Rect(sx, sy, sx+tileset.Tilewidth, sy+tileset.Tileheight)
			tile := tileset.EbitenImage.SubImage(tileRect).(*ebiten.Image)
			result.DrawImage(tile, op)
		}
	}

	t.imageOptions.GeoM.Reset()
	sh := screen.Bounds().Dy()
	t.imageOptions.GeoM.Translate(0, float64(sh-mapHeight))

	return result, nil
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

func (t *Tilemap) CreateCollisionBodies(space *physics.Space) error {
	if t == nil {
		return fmt.Errorf("the tilemap was not initialized")
	}

	mapHeight := t.Height * t.Tileheight
	yOffset := config.ScreenHeight - mapHeight

	for _, layer := range t.Layers {
		if layer.Type != "objectgroup" || layer.Name == "PlayerStart" {
			continue
		}

		for _, obj := range layer.Objects {
			rect := physics.NewRect(int(obj.X), int(obj.Y)+yOffset, int(obj.Width), int(obj.Height))
			obstacle := physics.NewObstacleRect(rect).AddCollision()
			obstacle.SetIsObstructive(true)
			space.AddBody(obstacle)
		}
	}

	return nil
}

func (t *Tilemap) ImageOptions() *ebiten.DrawImageOptions {
	return t.imageOptions
}

// GetPlayerStartPosition searches for a layer named "PlayerStart" in the tilemap's object layers.
// It assumes there is only one object in this layer and returns its x, y coordinates.
// The y coordinate is adjusted to account for the tilemap's rendering offset.
func (t *Tilemap) GetPlayerStartPosition() (x, y int, found bool) {
	if t == nil {
		return 0, 0, false
	}

	mapHeight := t.Height * t.Tileheight
	yOffset := config.ScreenHeight - mapHeight

	for _, layer := range t.Layers {
		if layer.Name == "PlayerStart" && layer.Type == "objectgroup" && len(layer.Objects) > 0 {
			obj := layer.Objects[0]
			px := int(math.Round(obj.X))
			py := int(math.Round(obj.Y)) + yOffset
			return px * config.Unit, py * config.Unit, true
		}
	}

	return 0, 0, false
}
