package scene

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/leandroatallah/firefly/internal/actors"
	"github.com/leandroatallah/firefly/internal/assets/font"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/core/transition"
	"github.com/leandroatallah/firefly/internal/items"
	"github.com/leandroatallah/firefly/internal/navigation"
	"github.com/leandroatallah/firefly/internal/systems/physics"
	"github.com/leandroatallah/firefly/internal/systems/tilemap"
	"github.com/setanarut/kamera/v2"
)

const (
	bgSound = "assets/Sketchbook 2024-06-19.ogg"
)

type LevelsScene struct {
	BaseScene
	count          int
	player         actors.PlayerEntity
	space          *physics.Space
	tilemap        *tilemap.Tilemap
	cam            *kamera.Camera
	levelCompleted bool
	mainText       *font.FontText
}

func NewLevelsScene() *LevelsScene {
	mainText, err := font.NewFontText(config.MainFontFace)
	if err != nil {
		log.Fatal(err)
	}
	return &LevelsScene{mainText: mainText}
}

func (s *LevelsScene) OnStart() {
	level, err := s.appContext.LevelManager.GetCurrentLevel()
	if err != nil {
		log.Fatalf("failed to get current level: %v", err)
	}

	// Init tilemap
	tm, err := tilemap.LoadTilemap(level.TilemapPath)
	if err != nil {
		log.Fatal(err)
	}
	s.tilemap = tm

	// Init audio manager
	s.audiomanager = s.appContext.AudioManager
	go func() {
		time.Sleep(1 * time.Second)
		s.audiomanager.SetVolume(0)
		s.audiomanager.PlaySound(bgSound)
	}()

	// Init space
	s.space = s.PhysicsSpace()
	s.space.SetTilemapDimensionsProvider(s)

	// Create player
	p, err := createPlayer(s.space)
	if err != nil {
		log.Fatal(err)
	}
	s.player = p
	s.space.AddBody(s.player)

	// Set items (and enemies?) position from tilemap
	itemsPos := s.tilemap.GetItemsPositionID()
	for _, i := range itemsPos {
		item, err := items.NewCollectibleCoinItem(i.X, i.Y)
		if err != nil {
			log.Fatal(err)
		}
		s.space.AddBody(item)
	}

	// Set player initial position from tilemap
	startX, startY, found := s.tilemap.GetPlayerStartPosition()
	if found {
		s.player.SetPosition(startX, startY)
	}

	// Init player camera
	pPos := s.player.Position().Min
	s.cam = kamera.NewCamera(
		float64(pPos.X),
		float64(pPos.Y),
		float64(config.ScreenWidth),
		float64(config.ScreenHeight),
	)
	s.cam.SmoothType = kamera.SmoothDamp
	s.cam.ShakeEnabled = true

	// Init collisions bodies and touch trigger for endpoints
	endpointTrigger := physics.NewTouchTrigger(s.finishLevel, s.player)
	s.tilemap.CreateCollisionBodies(s.space, endpointTrigger)

	s.levelCompleted = false
}

func (s *LevelsScene) GetTilemapWidth() int {
	if s.tilemap != nil && len(s.tilemap.Layers) > 0 {
		return s.tilemap.Layers[0].Width * s.tilemap.Tileheight
	}
	return config.ScreenWidth
}

func (s *LevelsScene) GetTilemapHeight() int {
	if s.tilemap != nil && len(s.tilemap.Layers) > 0 {
		return s.tilemap.Layers[0].Height * s.tilemap.Tileheight
	}
	return config.ScreenHeight
}

func (s *LevelsScene) Update() error {
	// Remove this
	s.CamDebug()

	s.count++

	// Update camera position
	pPos := s.player.Position().Min
	pWidth := s.player.Position().Dx()
	pHeight := s.player.Position().Dy()
	s.cam.LookAt(
		float64(pPos.X+(pWidth/2)),
		float64(pPos.Y+(pHeight/2)),
	)

	// Execute bodies updates
	space := s.PhysicsSpace()
	for _, i := range space.Bodies() {
		// Remove items marked as removed
		if item, ok := i.(items.Item); ok {
			if item.IsRemoved() {
				s.space.RemoveBody(i)
				continue
			}
			item.Update(space)
		}

		actor, ok := i.(actors.ActorEntity)
		if !ok {
			continue
		}
		if err := actor.Update(space); err != nil {
			return err
		}
	}

	return nil
}

func (s *LevelsScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x3c, 0xbc, 0xfc, 0xff})

	// Get tilemap image and draw based on camera
	tilemap, err := s.tilemap.Image(screen)
	if err != nil {
		log.Fatal(err)
	}
	s.cam.Draw(tilemap, s.tilemap.ImageOptions(), screen)

	// Draw collisions based on camera
	space := s.PhysicsSpace()
	for _, b := range space.Bodies() {
		switch body := b.(type) {
		case actors.PlayerEntity:
			continue
		case items.Item:
			if body.IsRemoved() {
				continue
			}
			opts := body.ImageOptions()
			opts.GeoM.Reset()
			pos := body.Position().Min
			opts.GeoM.Translate(float64(pos.X), float64(pos.Y))
			s.cam.Draw(body.Image(), opts, screen)
		case physics.Obstacle:
			opts := body.ImageOptions()
			opts.GeoM.Reset()
			pos := body.Position().Min
			opts.GeoM.Translate(float64(pos.X), float64(pos.Y))
			s.cam.Draw(body.Image(), opts, screen)
		default:
			continue
		}
	}

	// Draw player based on camera
	pPos := s.player.Position().Min
	img := s.player.Image()
	s.player.ImageOptions().GeoM.Reset()
	s.player.ImageOptions().GeoM.Translate(float64(pPos.X), float64(pPos.Y))
	s.cam.Draw(img, s.player.ImageOptions(), screen)

	s.DrawHUD(screen)
}

func (s *LevelsScene) OnFinish() {
	s.audiomanager.PauseMusic(bgSound)
}

func (s *LevelsScene) finishLevel() {
	if s.levelCompleted {
		return
	}

	s.levelCompleted = true
	s.appContext.SceneManager.NavigateTo(navigation.SceneSummary, transition.NewFader())
}

func createPlayer(space *physics.Space) (actors.PlayerEntity, error) {
	player, err := actors.NewPlayer(actors.Platform)
	if err != nil {
		return nil, err
	}
	space.AddBody(player)
	return player, nil
}

// TODO: REMOVE this method
func (s *LevelsScene) CamDebug() {

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		s.cam.Angle += 0.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		s.cam.Angle -= 0.02
	}

	if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
		s.cam.Reset()
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) { // zoom out
		s.cam.ZoomFactor /= 1.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) { // zoom in
		s.cam.ZoomFactor *= 1.02
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
