package scene

import (
	"fmt"
	"log"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/audio"
	"github.com/leandroatallah/firefly/internal/engine/render/camera"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/enemies"
	"github.com/leandroatallah/firefly/internal/engine/entity/items"
	"github.com/leandroatallah/firefly/internal/engine/render/tilemap"
)

type TilemapScene struct {
	BaseScene
	tilemap *tilemap.Tilemap
	cam     *camera.Controller
}

func NewTilemapScene(ctx *app.AppContext) *TilemapScene {
	scene := TilemapScene{}
	scene.SetAppContext(ctx)
	return &scene
}

func (s *TilemapScene) OnStart() {
	s.BaseScene.OnStart()

	// Load phases from context
	phase, err := s.AppContext().PhaseManager.GetCurrentPhase()
	if err != nil {
		log.Fatalf("failed to get current phase: %v", err)
	}

	// Init tilemap
	tm, err := tilemap.LoadTilemap(phase.TilemapPath)
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

func (s *TilemapScene) Audiomanager() *audio.AudioManager {
	return s.AppContext().AudioManager
}

func (s *TilemapScene) InitItems(items map[int]items.ItemType, factory *items.ItemFactory) error {
	itemsPos := s.tilemap.GetItemsPositionID()

	for _, i := range itemsPos {
		itemType, found := items[i.ItemType]
		if !found {
			return fmt.Errorf("Unable to find item by ID.")
		}

		item, err := factory.Create(itemType, i.X, i.Y, i.ID)
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
