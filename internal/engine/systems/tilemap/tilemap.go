package tilemap

import (
	_ "image/png"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/config"
)

type Tilemap struct {
	Height       int        `json:"height"`
	Infinite     bool       `json:"infinite"`
	Layers       []*Layer   `json:"layers"`
	Tileheight   int        `json:"tileheight"`
	Tilewidth    int        `json:"tilewidth"`
	Tilesets     []*Tileset `json:"tilesets"`
	image        *ebiten.Image
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

	cfg := config.Get()

	mapHeight := t.Height * t.Tileheight
	yOffset := cfg.ScreenHeight - mapHeight

	for _, layer := range t.Layers {
		if layer.Name == "PlayerStart" && layer.Type == "objectgroup" && len(layer.Objects) > 0 {
			obj := layer.Objects[0]
			px := int(math.Round(obj.X))
			py := int(math.Round(obj.Y)) + yOffset
			return px * cfg.Unit, py * cfg.Unit, true
		}
	}

	return 0, 0, false
}

type ItemPosition struct {
	X, Y int
	ID   int
}

func (t *Tilemap) GetItemsPositionID() []*ItemPosition {
	if t == nil {
		return nil
	}

	mapHeight := t.Height * t.Tileheight
	yOffset := config.Get().ScreenHeight - mapHeight

	res := []*ItemPosition{}
	var firstgid int
	var ts *Tileset

	for _, layer := range t.Layers {
		if layer.Name == "Items" && layer.Type == "objectgroup" && len(layer.Objects) > 0 {
			for _, obj := range layer.Objects {
				x16 := int(math.Round(obj.X))
				yValue := obj.Y
				if obj.Gid > 0 {
					yValue -= obj.Height
				}
				y16 := int(math.Round(yValue)) + yOffset
				if firstgid == 0 {
					firstgid = obj.Gid
					ts = t.findTileset(firstgid)
				}
				id := tilesetSourceID(ts, obj.Gid)
				res = append(res, &ItemPosition{X: x16, Y: y16, ID: id})
			}
		}
	}

	return res
}
