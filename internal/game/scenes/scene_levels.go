package gamescene

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
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
	gameplayer "github.com/leandroatallah/firefly/internal/game/actors/player"
	gamecamera "github.com/leandroatallah/firefly/internal/game/camera"
	gameitems "github.com/leandroatallah/firefly/internal/game/items"
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
	s.AppContext.SceneManager.NavigateTo(SceneSummary, transition.NewFader(), true)
}

func createPlayer(appContext *core.AppContext) (actors.PlayerEntity, error) {
	p, err := gameplayer.NewCherryPlayer(appContext)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// TODO: REMOVE this method
func (s *LevelsScene) CamDebug() {
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		s.cam.Kamera().Angle += 0.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		s.cam.Kamera().Angle -= 0.02
	}

	if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
		s.cam.Kamera().Reset()
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) { // zoom out
		s.cam.Kamera().ZoomFactor /= 1.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) { // zoom in
		s.cam.Kamera().ZoomFactor *= 1.02
	}
}

func (s *LevelsScene) DrawHUD(screen *ebiten.Image) {
	coinCount := 0

	if p, ok := s.player.(*actors.PlayerPlatform); ok {
		coinCount = p.CoinCount()
	}

	hud := ebiten.NewImage(74, 12)
	hud.Fill(color.White)
	hudOp := &ebiten.DrawImageOptions{}
	hudOp.GeoM.Translate(4, 5)
	textOp := &text.DrawOptions{}
	textOp.ColorScale.Scale(0, 0, 0, 255)
	textOp.GeoM.Translate(2, 2)
	s.mainText.Draw(hud, fmt.Sprintf("Score: %d", coinCount), 8, textOp)

	// Draw simple HUD score
	// HUD need to be drawed on screen and not on the camera.
	screen.DrawImage(hud, hudOp)
}

func (s *LevelsScene) SetCamTargetPointToSpace() {
	tPos := s.cam.Target().Position()
	targetRect := physics.NewObstacleRect(physics.NewRect(tPos.Min.X, tPos.Min.Y, tPos.Dx(), tPos.Dy()))
	targetBody := physics.NewPhysicsBody(targetRect)
	targetBody.SetID("TARGET")
	s.PhysicsSpace().AddBody(targetBody)
}

func (s *LevelsScene) DrawCamTargetPoint(screen *ebiten.Image) {
	tPos := s.cam.Target().Position()
	targetImage := ebiten.NewImage(tPos.Dx(), tPos.Dy())
	targetImage.Fill(color.RGBA{0xff, 0, 0, 0xff})
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Reset()
	opts.GeoM.Translate(float64(tPos.Min.X), float64(tPos.Min.Y))
	s.cam.Draw(targetImage, opts, screen)
}
