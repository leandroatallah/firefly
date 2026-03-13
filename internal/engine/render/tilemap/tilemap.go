package tilemap

import (
	"fmt"
	_ "image/png"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	flippedHorizontallyFlag = 0x80000000
	flippedVerticallyFlag   = 0x40000000
	flippedDiagonallyFlag   = 0x20000000
	flipFlagsMask           = flippedHorizontallyFlag | flippedVerticallyFlag | flippedDiagonallyFlag
)

// extractGIDAndFlags parses the raw GID to separate the tile ID and flip flags.
func extractGIDAndFlags(rawID int) (int, bool, bool, bool) {
	raw := uint32(rawID)
	h := (raw & flippedHorizontallyFlag) != 0
	v := (raw & flippedVerticallyFlag) != 0
	d := (raw & flippedDiagonallyFlag) != 0
	gid := int(raw &^ flipFlagsMask)
	return gid, h, v, d
}

type Tilemap struct {
	Height       int        `json:"height"`
	Width        int        `json:"width"`
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

	layer, found := t.FindLayerByName("PlayerStart")
	if !found {
		log.Printf("PlayerStart layer not found in tilemap")
		return 0, 0, false
	}

	obj := layer.Objects[0]
	px := int(math.Round(obj.X))
	py := int(math.Round(obj.Y))

	return px, py, true
}

// HasPlayerStartPosition checks if a layer named "PlayerStart" exists in the tilemap and has at least one object.
func (t *Tilemap) HasPlayerStartPosition() bool {
	if t == nil {
		return false
	}

	layer, found := t.FindLayerByName("PlayerStart")
	if !found {
		return false
	}

	return len(layer.Objects) > 0
}

// GetCameraStartPosition searches for a layer named "Camera" in the tilemap's object layers.
// It assumes there is only one object in this layer and returns its x, y coordinates.
func (t *Tilemap) GetCameraStartPosition() (x, y int, found bool) {
	if t == nil {
		return 0, 0, false
	}

	layer, found := t.FindLayerByName("Camera")
	if !found {
		return 0, 0, false
	}

	if len(layer.Objects) == 0 {
		return 0, 0, false
	}

	obj := layer.Objects[0]
	px := int(math.Round(obj.X))
	py := int(math.Round(obj.Y))

	return px, py, true
}

// HasCameraStartPosition checks if a layer named "Camera" exists in the tilemap and has at least one object.
func (t *Tilemap) HasCameraStartPosition() bool {
	if t == nil {
		return false
	}

	layer, found := t.FindLayerByName("Camera")
	if !found {
		return false
	}

	return len(layer.Objects) > 0
}

type ItemPosition struct {
	X, Y     int
	ItemType string
	ID       string
}

func (t *Tilemap) GetItemsPositionID() []*ItemPosition {
	if t == nil {
		return nil
	}

	res := []*ItemPosition{}
	var firstgid int
	var ts *Tileset

	layer, found := t.FindLayerByName("Items")
	if !found {
		log.Printf("Items layer not found in tilemap")
		return nil
	}

	itemCount := 0
	for _, obj := range layer.Objects {
		x16 := int(math.Round(obj.X))
		yValue := obj.Y

		gid, _, _, _ := extractGIDAndFlags(obj.Gid)

		if gid > 0 {
			yValue -= (obj.Height * 1)
		}
		y16 := int(math.Round(yValue))
		if firstgid == 0 {
			firstgid = gid
			ts = t.findTileset(firstgid)
		}
		itemType := fmt.Sprintf("%d", tilesetSourceID(ts, gid))
		var id string
		for _, p := range obj.Properties {
			if p.Name == "body_id" {
				id = p.Value
			}
			if p.Name == "item_type" {
				itemType = p.Value
			}
		}

		if id == "" {
			id = fmt.Sprintf("ITEM_%s_%d", itemType, itemCount)
			itemCount++
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

	layer, found := t.FindLayerByName("Enemies")
	if !found {
		log.Printf("Enemies layer not found in tilemap")
		return nil
	}

	enemyCount := 0
	for _, obj := range layer.Objects {
		x16 := int(math.Round(obj.X))
		yValue := obj.Y

		gid, _, _, _ := extractGIDAndFlags(obj.Gid)

		if gid > 0 {
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

		if id == "" {
			if enemyType == "" {
				enemyType = "UNKNOWN"
			}
			id = fmt.Sprintf("%s_%d", enemyType, enemyCount)
			enemyCount++
		}

		res = append(res, &EnemyPosition{X: x16, Y: y16, EnemyType: enemyType, ID: id})
	}

	return res
}

type NpcPosition struct {
	X, Y    int
	NpcType string
	ID      string
}

func (t *Tilemap) GetNpcsPositionID() []*NpcPosition {
	if t == nil {
		return nil
	}

	res := []*NpcPosition{}

	layer, found := t.FindLayerByName("NPCs") // Assumes a layer named "NPCs"
	if !found {
		log.Printf("NPCs layer not found in tilemap")
		return nil
	}

	npcCount := 0
	for _, obj := range layer.Objects {
		x16 := int(math.Round(obj.X))
		yValue := obj.Y
		if obj.Gid > 0 {
			yValue -= obj.Height
		}
		y16 := int(math.Round(yValue))

		var id, npcType string
		for _, p := range obj.Properties {
			if p.Name == "body_id" {
				id = p.Value
			}
			if p.Name == "npc_type" {
				npcType = p.Value
			}
		}

		if id == "" {
			if npcType == "" {
				npcType = "UNKNOWN"
			}
			id = fmt.Sprintf("%s_%d", npcType, npcCount)
			npcCount++
		}

		res = append(res, &NpcPosition{X: x16, Y: y16, NpcType: npcType, ID: id})
	}

	return res
}
