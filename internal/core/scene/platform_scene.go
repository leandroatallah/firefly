package scene

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/leandroatallah/firefly/internal/actors"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/navigation"
	"github.com/leandroatallah/firefly/internal/systems/physics"
	"github.com/leandroatallah/firefly/internal/systems/tilemap"
	"github.com/setanarut/kamera/v2"
)

const (
	bgSound = "assets/Sketchbook 2024-06-19.ogg"
)

type PlatformScene struct {
	BaseScene
	count   int
	player  actors.PlayerEntity
	space   *physics.Space
	tilemap *tilemap.Tilemap
	cam     *kamera.Camera
}

func (s *PlatformScene) OnStart() {
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

	// Init boundaries
	s.space = s.PhysicsSpace()
	s.space.SetTilemapDimensionsProvider(s)
	s.tilemap.CreateCollisionBodies(s.space)

	p, err := createPlayer(s.space)
	if err != nil {
		log.Fatal(err)
	}
	s.player = p

	// Set player initial position from tilemap
	startX, startY, found := s.tilemap.GetPlayerStartPosition()
	if found {
		s.player.SetPosition(startX, startY)
	}

	// Create player
	s.space.AddBody(s.player)

	pPos := s.player.Position().Min

	s.cam = kamera.NewCamera(
		float64(pPos.X),
		float64(pPos.Y),
		float64(config.ScreenWidth),
		float64(config.ScreenHeight),
	)
	s.cam.SmoothType = kamera.SmoothDamp
	s.cam.ShakeEnabled = true
}

func (s *PlatformScene) GetTilemapWidth() int {
	if s.tilemap != nil && len(s.tilemap.Layers) > 0 {
		return s.tilemap.Layers[0].Width * s.tilemap.Tileheight
	}
	return config.ScreenWidth // Fallback
}

func (s *PlatformScene) GetTilemapHeight() int {
	if s.tilemap != nil && len(s.tilemap.Layers) > 0 {
		return s.tilemap.Layers[0].Height * s.tilemap.Tileheight
	}
	return config.ScreenHeight // Fallback
}

func (s *PlatformScene) Update() error {
	// REMOVE
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
	// REMOVE

	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		s.appContext.LevelManager.AdvanceToNextLevel()
		s.appContext.SceneManager.NavigateTo(navigation.ScenePlatform, nil)
	}

	pPos := s.player.Position().Min
	pWidth := s.player.Position().Dx()
	pHeight := s.player.Position().Dy()
	s.cam.LookAt(
		float64(pPos.X+(pWidth/2)),
		float64(pPos.Y+(pHeight/2)),
	)

	space := s.PhysicsSpace()

	s.count++

	for _, i := range space.Bodies() {
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

func Translate(bx *[4]float64, x, y float64) {
	bx[0] += x
	bx[1] += y
}

func (s *PlatformScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x3c, 0xbc, 0xfc, 0xff})

	tilemapImg, err := s.tilemap.ParseToImage(screen)
	if err != nil {
		log.Fatal(err)
	}
	s.cam.Draw(tilemapImg, s.tilemap.ImageOptions(), screen)

	space := s.PhysicsSpace()

	bodyOpts := &ebiten.DrawImageOptions{}
	for _, b := range space.Bodies() {
		switch body := b.(type) {
		case actors.PlayerEntity:
			continue
		case physics.Obstacle:
			bodyOpts.GeoM.Reset()
			pos := body.Position().Min
			bodyOpts.GeoM.Translate(float64(pos.X), float64(pos.Y))
			s.cam.Draw(body.Image(), bodyOpts, screen)
		}
	}

	pPos := s.player.Position().Min
	img := s.player.Image()
	s.player.ImageOptions().GeoM.Reset()
	s.player.ImageOptions().GeoM.Translate(float64(pPos.X), float64(pPos.Y))
	s.cam.Draw(img, s.player.ImageOptions(), screen)
}

func (s *PlatformScene) OnFinish() {
	s.audiomanager.PauseMusic(bgSound)
}

func createPlatform(rect *physics.Rect, space *physics.Space) {
	o := physics.NewObstacleRect(rect).AddCollision()
	o.SetIsObstructive(true)

	space.AddBody(o)
}

func createPlayer(space *physics.Space) (actors.PlayerEntity, error) {
	player, err := actors.NewPlayer(actors.Platform)
	if err != nil {
		return nil, err
	}
	space.AddBody(player)
	return player, nil
}
