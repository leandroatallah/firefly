package scene

import (
	"fmt"
	"image"
	"log"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/audio"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/enemies"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/npcs"
	"github.com/leandroatallah/firefly/internal/engine/entity/items"
	"github.com/leandroatallah/firefly/internal/engine/render/camera"
	"github.com/leandroatallah/firefly/internal/engine/render/tilemap"
)

type TilemapScene struct {
	BaseScene
	tilemap      *tilemap.Tilemap
	cam          *camera.Controller
	cameraConfig CameraConfig
}

func NewTilemapScene(ctx *app.AppContext) *TilemapScene {
	scene := TilemapScene{
		cam: camera.NewController(0, 0),
	}
	scene.SetAppContext(ctx)
	return &scene
}

func (s *TilemapScene) SetCameraConfig(config CameraConfig) {
	s.cameraConfig = config
	s.cam.SetFollowing(config.Mode == CameraModeFollow)
}

func (s *TilemapScene) Update() error {
	s.cam.Update()
	return s.BaseScene.Update()
}

func (s *TilemapScene) Camera() *camera.Controller {
	return s.cam
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

func (s *TilemapScene) GetCameraBounds() (image.Rectangle, bool) {
	if s.cam == nil || s.cam.Bounds() == nil {
		return image.Rectangle{}, false
	}
	return *s.cam.Bounds(), true
}

func (s *TilemapScene) Tilemap() *tilemap.Tilemap {
	return s.tilemap
}

func (s *TilemapScene) Audiomanager() *audio.AudioManager {
	return s.AppContext().AudioManager
}

func (s *TilemapScene) SetPlayerStartPosition(p actors.ActorEntity) {
	// Set player initial position from tilemap
	if x, y, found := s.tilemap.GetPlayerStartPosition(); found {
		// PlayerStart is a point object, so y is the ground level
		// Position player so their bottom aligns with y
		actorHeight := p.Position().Dy()
		p.SetPosition(x, y-actorHeight)
	}
}

func InitEnemies[T actors.ActorEntity](s *TilemapScene, factory *enemies.EnemyFactory[T]) error {
	enemiesPos := s.Tilemap().GetEnemiesPositionID()

	for _, e := range enemiesPos {
		enemy, err := factory.Create(enemies.EnemyType(e.EnemyType), e.X, e.Y, e.ID)
		pos := enemy.Position()
		// Adjust Y so enemy's base aligns with obstacle top
		enemy.SetPosition(pos.Min.X, pos.Min.Y-(pos.Dy()-s.tilemap.Tileheight))
		if err != nil {
			return err
		}

		s.PhysicsSpace().AddBody(enemy)
		if s.AppContext().ActorManager != nil {
			s.AppContext().ActorManager.Register(enemy)
		}
	}

	return nil
}

func InitNPCs[T actors.ActorEntity](s *TilemapScene, factory *npcs.NpcFactory[T]) error {
	npcsPos := s.Tilemap().GetNpcsPositionID()

	for _, n := range npcsPos {
		npc, err := factory.Create(npcs.NpcType(n.NpcType), n.X, n.Y, n.ID)
		pos := npc.Position()
		// Adjust Y so npc's base aligns with obstacle top
		npc.SetPosition(pos.Min.X, pos.Min.Y-(pos.Dy()-s.tilemap.Tileheight))
		if err != nil {
			return err
		}

		s.PhysicsSpace().AddBody(npc)
		if s.AppContext().ActorManager != nil {
			s.AppContext().ActorManager.Register(npc)
		}
	}

	return nil
}

func InitItems[T items.Item](s *TilemapScene, factory *items.ItemFactory[T]) error {
	itemsPos := s.Tilemap().GetItemsPositionID()

	for _, i := range itemsPos {
		item, err := factory.Create(items.ItemType(i.ItemType), i.X, i.Y, i.ID)
		// pos := item.Position()
		// item.SetPosition(pos.Min.X, pos.Min.Y-pos.Dy()) // Adjust Y position based on npc height
		if err != nil {
			return err
		}

		item.SetID(fmt.Sprintf("ITEM_%v", i.ID))
		s.PhysicsSpace().AddBody(item)
	}

	return nil
}
