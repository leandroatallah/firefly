package gamescenelevels

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/enemies"
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
	"github.com/leandroatallah/firefly/internal/engine/camera"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/scene"
	"github.com/leandroatallah/firefly/internal/engine/scene/transition"
	"github.com/leandroatallah/firefly/internal/engine/entity/items"
	"github.com/leandroatallah/firefly/internal/engine/physics"
	gameenemies "github.com/leandroatallah/firefly/internal/game/actors/enemies"
	gamecamera "github.com/leandroatallah/firefly/internal/game/camera"
	gameitems "github.com/leandroatallah/firefly/internal/game/items"
	scenestypes "github.com/leandroatallah/firefly/internal/game/scenes/types"
)

const (
	bgSound = "assets/audio/Sketchbook.ogg"
)

type LevelsScene struct {
	scene.TilemapScene
	count          int
	player         actors.ActorEntity
	cam            *camera.Controller
	levelCompleted bool
	mainText       *font.FontText
}

func NewLevelsScene(context *app.AppContext) *LevelsScene {
	mainText, err := font.NewFontText(config.Get().MainFontFace)
	if err != nil {
		log.Fatal(err)
	}
	tilemapScene := scene.NewTilemapScene(context)
	scene := LevelsScene{
		TilemapScene: *tilemapScene,
		mainText:     mainText,
	}
	scene.SetAppContext(context)
	return &scene
}

func (s *LevelsScene) OnStart() {
	s.TilemapScene.OnStart()

	go func() {
		time.Sleep(1 * time.Second)
		am := s.AppContext().AudioManager
		am.SetVolume(0.25)
		am.PlayMusic(bgSound)
	}()

	// Create player and register to space and context
	p, err := createPlayer()
	if err != nil {
		log.Fatal(err)
	}
	s.player = p
	s.AppContext().ActorManager.Register(s.player)
	s.PhysicsSpace().AddBody(s.player)

	// Set items map to factory creation process
	itemsMap := map[int]items.ItemType{
		0: gameitems.CollectibleCoinType,
	}

	// Set items position from tilemap
	f := items.NewItemFactory(gameitems.InitItemMap(s.AppContext()))
	s.InitItems(itemsMap, f)

	// Set enemies position from tilemap
	enemyFactory := enemies.NewEnemyFactory(gameenemies.InitEnemyMap(s.AppContext()))
	s.InitEnemies(enemyFactory)

	s.SetPlayerStartPosition(s.player)

	// Init camera target
	pPos := s.player.Position().Min
	s.cam = gamecamera.New(pPos.X, pPos.Y)
	s.cam.SetFollowTarget(s.player)

	// Init collisions bodies and touch trigger for endpoints
	endpointTrigger := physics.NewTouchTrigger(s.finishLevel, s.player)
	s.Tilemap().CreateCollisionBodies(s.PhysicsSpace(), endpointTrigger)

	s.levelCompleted = false
}

func (s *LevelsScene) Update() error {
	if config.Get().CamDebug {
		s.CamDebug()
	}

	s.count++

	s.cam.Update()

	// Execute bodies updates
	space := s.PhysicsSpace()
	for _, i := range space.Bodies() {
		switch b := i.(type) {
		// ActorEntity case should came first. It can be confused with body.Obstacle
		case actors.ActorEntity:
			if err := b.Update(space); err != nil {
				return err
			}
		case items.Item:
			// Remove items marked as removed
			if b.IsRemoved() {
				s.PhysicsSpace().RemoveBody(i)
				continue
			}
			if err := b.Update(space); err != nil {
				return err
			}
		case body.Obstacle:
			continue
		}
	}

	return nil
}

func (s *LevelsScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x3c, 0xbc, 0xfc, 0xff})

	// Get tilemap image and draw based on camera
	tilemap, err := s.Tilemap().Image(screen)
	if err != nil {
		log.Fatal(err)
	}
	s.cam.Draw(tilemap, s.Tilemap().ImageOptions(), screen)

	// Draw bodies based on camera
	space := s.PhysicsSpace()
	for _, b := range space.Bodies() {
		switch sb := b.(type) {
		case actors.ActorEntity:
			opts := sb.ImageOptions()
			sb.UpdateImageOptions()
			if config.Get().CollisionBox {
				s.cam.Draw(sb.ImageCollisionBox(), opts, screen)
			} else {
				s.cam.Draw(sb.Image(), opts, screen)
			}
		case items.Item:
			if sb.IsRemoved() {
				continue
			}
			opts := sb.ImageOptions()
			sb.UpdateImageOptions()
			if config.Get().CollisionBox {
				s.cam.Draw(sb.ImageCollisionBox(), opts, screen)
			} else {
				s.cam.Draw(sb.Image(), opts, screen)
			}
		case body.Obstacle:
			if config.Get().CollisionBox {
				opts := sb.ImageOptions()
				sb.UpdateImageOptions()
				s.cam.Draw(sb.ImageCollisionBox(), opts, screen)
			}
		}
	}

	s.DrawHUD(screen)
}

func (s *LevelsScene) OnFinish() {
	s.TilemapScene.OnFinish()

	s.Audiomanager().PauseMusic(bgSound)
}

func (s *LevelsScene) finishLevel() {
	if s.levelCompleted {
		return
	}

	s.levelCompleted = true
	s.AppContext().SceneManager.NavigateTo(scenestypes.SceneSummary, transition.NewFader(), true)
}
