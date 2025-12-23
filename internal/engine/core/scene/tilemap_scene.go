package scene

import (
	"fmt"
	"log"

	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/actors/enemies"
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

func NewTilemapScene(ctx *core.AppContext) *TilemapScene {
	scene := TilemapScene{}
	scene.SetAppContext(ctx)
	return &scene
}

func (s *TilemapScene) OnStart() {
	s.BaseScene.OnStart()

	// Load levels from context
	level, err := s.AppContext().LevelManager.GetCurrentLevel()
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
	return s.AppContext().AudioManager
}

func (s *TilemapScene) InitItems(items map[int]items.ItemType, factory *items.ItemFactory) error {
	itemsPos := s.tilemap.GetItemsPositionID()

	for _, i := range itemsPos {
		itemType, found := items[i.ItemType]
		if !found {
			return fmt.Errorf("Unable to find item by ID.")
		}

		item, err := factory.Create(itemType, i.X, i.Y)
		if err != nil {
			return err
		}

		item.SetID(fmt.Sprintf("ITEM_%v", i.ID))
		s.PhysicsSpace().AddBody(item)
	}

	return nil
}

func (s *TilemapScene) InitEnemies(factory *enemies.EnemyFactory) error {
	enemiesPos := s.tilemap.GetEnemiesPositionID()

	for _, e := range enemiesPos {
		enemy, err := factory.Create(enemies.EnemyType(e.EnemyType), e.X, e.Y, e.ID)
		if err != nil {
			return err
		}

		s.PhysicsSpace().AddBody(enemy)
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
