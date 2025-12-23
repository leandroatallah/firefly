package tilemap

import (
	_ "image/png"
	"log"
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

type Property struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
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
	Gid        int        `json:"gid"`
	Height     float64    `json:"height"`
	Id         int        `json:"id"`
	Name       string     `json:"name"`
	Rotation   float64    `json:"rotation"`
	Type       string     `json:"type"`
	Visible    bool       `json:"visible"`
	Width      float64    `json:"width"`
	X          float64    `json:"x"`
	Y          float64    `json:"y"`
	Properties []Property `json:"properties"`
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
	if t.image == nil {
		img, err := t.ParseToImage(screen)
		if err != nil {
			return nil, err
		}
		t.image = img
	}

	t.Reset(screen)

	return t.image, nil
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

	cfg := config.Get()
	mapHeight := t.Height * t.Tileheight
	yOffset := mapHeight - cfg.ScreenHeight - 100

	layer, err := t.FindLayerByName("PlayerStart")
	if err != nil {
		log.Printf("failed to get player start position: %v", err)
		return 0, 0, false
	}

	obj := layer.Objects[0]
	px := int(math.Round(obj.X))
	py := int(math.Round(obj.Y)) + yOffset

	return px, py, true
}

type ItemPosition struct {
	X, Y     int
	ItemType int
	ID       string
}

func (t *Tilemap) GetItemsPositionID() []*ItemPosition {
	if t == nil {
		return nil
	}

	res := []*ItemPosition{}
	var firstgid int
	var ts *Tileset

	layer, err := t.FindLayerByName("Items")
	if err != nil {
		log.Printf("failed to get items position: %v", err)
		return nil
	}

	for _, obj := range layer.Objects {
		x16 := int(math.Round(obj.X))
		yValue := obj.Y
		if obj.Gid > 0 {
			yValue -= obj.Height
		}
		y16 := int(math.Round(yValue))
		if firstgid == 0 {
			firstgid = obj.Gid
			ts = t.findTileset(firstgid)
		}
		itemType := tilesetSourceID(ts, obj.Gid)

		var id string
		for _, p := range obj.Properties {
			if p.Name == "body_id" {
				id = p.Value
				break
			}
		}
		// o.SetID(fmt.Sprintf("%v_%v", prefix, id))
		res = append(res, &ItemPosition{X: x16, Y: y16, ItemType: itemType, ID: id})
	}

	return res
}

type EnemyPosition struct {
	X, Y      int
	EnemyType string
	ID        string
}

func (t *Tilemap) GetEnemiesPositionID() []*EnemyPosition {
	if t == nil {
		return nil
	}

	res := []*EnemyPosition{}

	layer, err := t.FindLayerByName("Enemies")
	if err != nil {
		log.Printf("failed to get enemies position: %v", err)
		return nil
	}

	for _, obj := range layer.Objects {
		x16 := int(math.Round(obj.X))
		yValue := obj.Y
		if obj.Gid > 0 {
			yValue -= obj.Height
		}
		y16 := int(math.Round(yValue))

		var id, enemyType string
		for _, p := range obj.Properties {
			if p.Name == "body_id" {
				id = p.Value
			}
			if p.Name == "enemy_type" {
				enemyType = p.Value
			}
		}
		res = append(res, &EnemyPosition{X: x16, Y: y16, EnemyType: enemyType, ID: id})
	}

	return res
}
