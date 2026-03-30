package tilemap

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/physics/space"
)

func TestExtractGIDAndFlags(t *testing.T) {
	raw := int(flippedHorizontallyFlag | flippedDiagonallyFlag | 123)
	gid, h, v, d := extractGIDAndFlags(raw)
	if gid != 123 || !h || v || !d {
		t.Fatalf("unexpected gid/flags: %d %v %v %v", gid, h, v, d)
	}
}

func TestFindLayerByNameAndPlayerStart(t *testing.T) {
	tm := &Tilemap{
		Width:      10,
		Height:     10,
		Tilewidth:  16,
		Tileheight: 16,
		Layers: []*Layer{
			{Name: "Background", Visible: true, Type: "tilelayer", Width: 10, Height: 10},
			{Name: "PlayerStart", Visible: true, Type: "objectgroup", Objects: []*Obstacle{
				{X: 32, Y: 48, Width: 16, Height: 16},
			}},
		},
		Tilesets: []*Tileset{{Firstgid: 1, Columns: 1, Tilewidth: 16, Tileheight: 16}},
	}

	layer, ok := tm.FindLayerByName("PlayerStart")
	if !ok || layer == nil {
		t.Fatalf("expected PlayerStart layer")
	}

	x, y, found := tm.GetPlayerStartPosition()
	if !found || x != 32 || y != 48 {
		t.Fatalf("GetPlayerStartPosition got (%d,%d,%v)", x, y, found)
	}
}

func TestCreateCollisionBodiesForEndpointAndObstacles(t *testing.T) {
	tm := &Tilemap{
		Width:      4,
		Height:     4,
		Tilewidth:  8,
		Tileheight: 8,
		Layers: []*Layer{
			{
				Name:    "Endpoint",
				Type:    "tilelayer",
				Visible: true,
				Width:   4,
				Height:  4,
				Data:    []int{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			{
				Name:    "Obstacles",
				Type:    "tilelayer",
				Visible: true,
				Width:   4,
				Height:  4,
				Data:    []int{0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
		Tilesets: []*Tileset{{Firstgid: 1, Columns: 1, Tilewidth: 8, Tileheight: 8}},
	}

	sp := space.NewSpace()
	var capturedID string
	tm.CreateCollisionBodies(sp, func(id string) body.Touchable {
		capturedID = id
		return nil
	})

	bs := sp.Bodies()
	if len(bs) == 0 {
		t.Fatalf("expected bodies created for endpoint/obstacles")
	}

	foundEndpoint := false
	foundObstacle := false
	for _, b := range bs {
		if b.IsObstructive() {
			foundObstacle = true
		} else {
			foundEndpoint = true
		}
	}
	if !foundEndpoint || !foundObstacle {
		t.Fatalf("expected both endpoint and obstacle bodies")
	}
	if capturedID == "" {
		t.Fatalf("expected endpoint factory to receive id")
	}
}

func TestGetItemsEnemiesNpcsPositions(t *testing.T) {
	tm := &Tilemap{
		Tilewidth:  16,
		Tileheight: 16,
		Layers: []*Layer{
			{
				Name:    "Items",
				Visible: true,
				Type:    "objectgroup",
				Objects: []*Obstacle{{
					Gid: 1, X: 10, Y: 20, Width: 16, Height: 16,
					Properties: []Property{{Name: "item_type", Value: "coin"}},
				}},
			},
			{
				Name:    "Enemies",
				Visible: true,
				Type:    "objectgroup",
				Objects: []*Obstacle{{
					Gid: 1, X: 30, Y: 40, Width: 16, Height: 16,
					Properties: []Property{{Name: "enemy_type", Value: "bat"}, {Name: "body_id", Value: "E1"}},
				}},
			},
			{
				Name:    "NPCs",
				Visible: true,
				Type:    "objectgroup",
				Objects: []*Obstacle{{
					Gid: 1, X: 50, Y: 60, Width: 16, Height: 16,
					Properties: []Property{{Name: "npc_type", Value: "guide"}, {Name: "body_id", Value: "N1"}},
				}},
			},
		},
		Tilesets: []*Tileset{{Firstgid: 1, Columns: 1, Tilewidth: 16, Tileheight: 16}},
	}

	items := tm.GetItemsPositionID()
	enemies := tm.GetEnemiesPositionID()
	npcs := tm.GetNpcsPositionID()

	if len(items) != 1 || items[0].ItemType != "coin" {
		t.Fatalf("unexpected items: %#v", items)
	}
	autoID := "ITEM_coin_0" // auto generated id
	if items[0].ID != autoID {
		t.Fatalf("unexpected item id: %s", items[0].ID)
	}
	if len(enemies) != 1 || enemies[0].ID != "E1" || enemies[0].EnemyType != "bat" {
		t.Fatalf("unexpected enemies: %#v", enemies)
	}
	if len(npcs) != 1 || npcs[0].ID != "N1" || npcs[0].NpcType != "guide" {
		t.Fatalf("unexpected npcs: %#v", npcs)
	}
}
