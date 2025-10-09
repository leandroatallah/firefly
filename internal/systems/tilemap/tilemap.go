package tilemap

import (
	_ "image/png"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/config"
)

type Tilemap struct {
	Height        int        `json:"height"`
	Infinite      bool       `json:"infinite"`
	Layers        []*Layer   `json:"layers"`
	Tileheight    int        `json:"tileheight"`
	Tilewidth     int        `json:"tilewidth"`
	Tilesets      []*Tileset `json:"tilesets"`
	image         *ebiten.Image
	imageBaseDone bool
	imageOptions  *ebiten.DrawImageOptions
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
	Gid      int     `json:"gid"`
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

func (t *Tilemap) Image(screen *ebiten.Image) (*ebiten.Image, error) {
	var err error
	// ParseToImage draw the base layer only one time.
	// TODO: Check a way to reduce other layers reset rate.
	t.image, err = t.ParseToImage(screen)
	if err != nil {
		return nil, err
	}

	t.Reset(screen)

	return t.image, nil
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
