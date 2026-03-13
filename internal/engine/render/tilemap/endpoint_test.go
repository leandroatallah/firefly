package tilemap

import (
	"testing"

	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/physics/space"
)

func TestCreateCollisionBodiesWithTwoEndpointLayers(t *testing.T) {
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
				Name:    "Endpoint",
				Type:    "objectgroup",
				Visible: true,
				Objects: []*Obstacle{
					{
						X: 16, Y: 16, Width: 8, Height: 8,
						Properties: []Property{{Name: "event_id", Value: "OBJ_ENDPOINT"}},
					},
				},
			},
		},
		Tilesets: []*Tileset{{Firstgid: 1, Columns: 1, Tilewidth: 8, Tileheight: 8}},
	}

	sp := space.NewSpace()
	capturedIDs := make(map[string]bool)
	tm.CreateCollisionBodies(sp, func(id string) body.Touchable {
		capturedIDs[id] = true
		return nil
	})

	bs := sp.Bodies()
	if len(bs) != 2 {
		t.Fatalf("expected 2 bodies created, got %d", len(bs))
	}

	if !capturedIDs["ENDPOINT_8_0"] {
		t.Errorf("expected tile endpoint ID ENDPOINT_8_0, got %v", capturedIDs)
	}
	if !capturedIDs["OBJ_ENDPOINT"] {
		t.Errorf("expected object endpoint ID OBJ_ENDPOINT, got %v", capturedIDs)
	}
}
