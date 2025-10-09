package scene

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/actors"
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
}

func NewLevelsScene() *LevelsScene {
	return &LevelsScene{}
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

	// Manually add Item to test
	coin := items.NewCollectibleCoinItem()
	s.space.AddBody(coin)

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
	itemsTrigger := physics.NewTouchTrigger(func() {
		fmt.Println("ITEM TRIGGER")
	}, s.player)

	// TODO: Review this implementation
	triggerMap := map[tilemap.LayerNameID]physics.Touchable{
		tilemap.EndpointLayer: endpointTrigger,
		tilemap.ItemsLayer:    itemsTrigger,
	}

	s.tilemap.CreateCollisionBodies(s.space, triggerMap)

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
		if item, ok := i.(items.Item); ok && item.IsRemoved() {
			s.space.RemoveBody(i)
			continue
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
	bodyOpts := &ebiten.DrawImageOptions{}
	for _, b := range space.Bodies() {
		switch body := b.(type) {
		case actors.PlayerEntity:
			continue
		case items.Item:
			if b.(items.Item).IsRemoved() {
				continue
			}
			bodyOpts.GeoM.Reset()
			pos := body.Position().Min
			bodyOpts.GeoM.Translate(float64(pos.X), float64(pos.Y))
			s.cam.Draw(body.Image(), bodyOpts, screen)
		case physics.Obstacle:
			bodyOpts.GeoM.Reset()
			pos := body.Position().Min
			bodyOpts.GeoM.Translate(float64(pos.X), float64(pos.Y))
			s.cam.Draw(body.ImageCollisionBox(), bodyOpts, screen)
		}
	}

	// Draw player based on camera
	pPos := s.player.Position().Min
	img := s.player.Image()
	s.player.ImageOptions().GeoM.Reset()
	s.player.ImageOptions().GeoM.Translate(float64(pPos.X), float64(pPos.Y))
	s.cam.Draw(img, s.player.ImageOptions(), screen)
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
