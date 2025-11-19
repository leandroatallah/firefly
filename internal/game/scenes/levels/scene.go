package gamescenelevels

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
	"github.com/leandroatallah/firefly/internal/engine/camera"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/core"
	"github.com/leandroatallah/firefly/internal/engine/core/scene"
	"github.com/leandroatallah/firefly/internal/engine/core/transition"
	"github.com/leandroatallah/firefly/internal/engine/items"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
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
	player         actors.PlayerEntity
	cam            *camera.Controller
	levelCompleted bool
	mainText       *font.FontText
}

func NewLevelsScene(context *core.AppContext) *LevelsScene {
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
		s.Audiomanager().SetVolume(0.25)
		s.Audiomanager().PlayMusic(bgSound)
	}()

	// Create player and register to space and context
	p, err := createPlayer(s.AppContext)
	if err != nil {
		log.Fatal(err)
	}
	s.player = p
	s.player.SetID("player")
	s.AppContext.ActorManager.Register(s.player)
	s.PhysicsSpace().AddBody(s.player)

	// Set items map to factory creation process
	itemsMap := map[int]items.ItemType{
		0: gameitems.CollectibleCoinType,
		1: gameitems.SignpostType,
	}

	// Set items position from tilemap
	f := items.NewItemFactory(gameitems.InitItemMap())
	s.InitItems(itemsMap, f)

	s.SetPlayerStartPosition(s.player)

	// Init camera target
	pPos := s.player.Position().Min
	s.cam = gamecamera.New(pPos.X, pPos.Y)
	s.cam.SetFollowTarget(s.player)

	s.SetCamTargetPointToSpace()

	// Init collisions bodies and touch trigger for endpoints
	endpointTrigger := physics.NewTouchTrigger(s.finishLevel, s.player)
	s.Tilemap().CreateCollisionBodies(s.PhysicsSpace(), endpointTrigger)

	s.levelCompleted = false
}

func (s *LevelsScene) Update() error {
	// Remove this
	s.CamDebug()

	s.count++

	s.cam.Update()

	// Execute bodies updates
	space := s.PhysicsSpace()
	for _, i := range space.Bodies() {
		if item, ok := i.(items.Item); ok {
			// Remove items marked as removed
			if item.IsRemoved() {
				s.PhysicsSpace().RemoveBody(i)
				continue
			}
		}

		// Update actors and items that are not actors
		if actor, ok := i.(actors.ActorEntity); ok {
			if err := actor.Update(space); err != nil {
				return err
			}
		} else if item, ok := i.(items.Item); ok {
			item.Update(space)
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

	// Draw collisions based on camera
	space := s.PhysicsSpace()
	for _, b := range space.Bodies() {
		switch sb := b.(type) {
		case actors.PlayerEntity:
			// Draw player based on camera
			if img := s.player.Image(); img != nil {
				opts := *s.player.ImageOptions()
				s.cam.Draw(img, &opts, screen)
			}
		case items.Item:
			if sb.IsRemoved() {
				continue
			}
			opts := sb.ImageOptions()
			opts.GeoM.Reset()
			pos := sb.Position().Min
			opts.GeoM.Translate(float64(pos.X), float64(pos.Y))
			s.cam.Draw(sb.Image(), opts, screen)
		case body.Obstacle:
			opts := sb.ImageOptions()
			opts.GeoM.Reset()
			pos := sb.Position().Min
			opts.GeoM.Translate(float64(pos.X), float64(pos.Y))
			s.cam.Draw(sb.Image(), opts, screen)
		default:
			if b.ID() == "TARGET" {
				s.DrawCamTargetPoint(screen)
			}
			continue
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
	s.AppContext.SceneManager.NavigateTo(scenestypes.SceneSummary, transition.NewFader(), true)
}
