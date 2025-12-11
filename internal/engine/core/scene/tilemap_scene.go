package scene

import (
	"fmt"
	"log"
	"time"

	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/camera"
	"github.com/leandroatallah/firefly/internal/engine/core"
	"github.com/leandroatallah/firefly/internal/engine/items"
	"github.com/leandroatallah/firefly/internal/engine/systems/audiomanager"
	"github.com/leandroatallah/firefly/internal/engine/systems/tilemap"
)

type TilemapScene struct {
	BaseScene
	tilemap *tilemap.Tilemap
	cam     *camera.Controller
}

func NewTilemapScene(context *core.AppContext) *TilemapScene {
	scene := TilemapScene{}
	scene.SetAppContext(context)
	return &scene
}

func (s *TilemapScene) OnStart() {
	s.BaseScene.OnStart()

	// Init audio manager
	s.audiomanager = s.AppContext.AudioManager

	// Load levels from context
	level, err := s.AppContext.LevelManager.GetCurrentLevel()
	if err != nil {
		log.Fatalf("failed to get current level: %v", err)
	}

	// Init tilemap
	tm, err := tilemap.LoadTilemap(level.TilemapPath)
	if err != nil {
		log.Fatal(err)
	}
	s.tilemap = tm

	// Init space
	s.PhysicsSpace().SetTilemapDimensionsProvider(s)
}

func (s *TilemapScene) GetTilemapWidth() int {
	if s.tilemap != nil && len(s.tilemap.Layers) > 0 {
		return s.tilemap.Layers[0].Width * s.tilemap.Tileheight
	}
	return config.Get().ScreenWidth
}

func (s *TilemapScene) GetTilemapHeight() int {
	if s.tilemap != nil && len(s.tilemap.Layers) > 0 {
		return s.tilemap.Layers[0].Height * s.tilemap.Tileheight
	}
	return config.Get().ScreenHeight
}

func (s *TilemapScene) Tilemap() *tilemap.Tilemap {
	return s.tilemap
}

func (s *TilemapScene) Audiomanager() *audiomanager.AudioManager {
	return s.audiomanager
}

func (s *TilemapScene) InitItems(items map[int]items.ItemType, factory *items.ItemFactory) error {
	itemsPos := s.tilemap.GetItemsPositionID()

	if len(itemsPos) > 0 {
		for _, i := range itemsPos {
			itemType, found := items[i.ID]
			if !found {
				return fmt.Errorf("Unable to find item by ID.")
			}

			item, err := factory.Create(itemType, i.X, i.Y)
			if err != nil {
				return err
			}
			// TODO: Improve this
			if item.ID() == "" {
				item.SetID(fmt.Sprintf("ITEM-%d", time.Now().Nanosecond()))
			} else {
				item.SetID(fmt.Sprintf("%vITEM-%d", item.ID(), time.Now().Nanosecond()))
			}
			s.PhysicsSpace().AddBody(item)
		}
	}

	return nil
}

func (s *TilemapScene) SetPlayerStartPosition(p actors.ActorEntity) {
	// Set player initial position from tilemap
	if x, y, found := s.tilemap.GetPlayerStartPosition(); found {
		// Update Y position based on player height
		y -= p.Position().Dy() * config.Get().Unit
		p.SetPosition(x, y)
	}
}
