package gamescene

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/core"
	"github.com/leandroatallah/firefly/internal/engine/core/scene"
	"github.com/leandroatallah/firefly/internal/engine/core/transition"
	"github.com/leandroatallah/firefly/internal/engine/items"
	"github.com/leandroatallah/firefly/internal/engine/sequences"
	"github.com/leandroatallah/firefly/internal/engine/systems/audiomanager"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
	"github.com/leandroatallah/firefly/internal/engine/systems/tilemap"
	gameitems "github.com/leandroatallah/firefly/internal/game/items"
	"github.com/setanarut/kamera/v2"
)

const (
	bgSound = "assets/Sketchbook 2024-06-19.ogg"
)

type LevelsScene struct {
	scene.BaseScene
	count          int
	player         actors.PlayerEntity
	space          *physics.Space
	tilemap        *tilemap.Tilemap
	cam            *kamera.Camera
	levelCompleted bool
	mainText       *font.FontText
	audiomanager   *audiomanager.AudioManager
	itemsMap       map[int]items.ItemType

	sequencePlayer *sequences.SequencePlayer
}

func NewLevelsScene(context *core.AppContext) *LevelsScene {
	mainText, err := font.NewFontText(config.Get().MainFontFace)
	if err != nil {
		log.Fatal(err)
	}
	scene := LevelsScene{mainText: mainText}
	scene.SetAppContext(context)
	return &scene
}

func (s *LevelsScene) OnStart() {
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

	// Init audio manager
	s.audiomanager = s.AppContext.AudioManager
	go func() {
		time.Sleep(1 * time.Second)
		s.audiomanager.SetVolume(0)
		s.audiomanager.PlaySound(bgSound)
	}()

	// Init space
	s.space = s.PhysicsSpace()
	s.space.SetTilemapDimensionsProvider(s)

	// Create player
	p, err := createPlayer(s.space, s.AppContext)
	if err != nil {
		log.Fatal(err)
	}
	s.player = p
	s.player.SetID("player")
	s.AppContext.ActorManager.Register(s.player)
	s.space.AddBody(s.player)

	// Set items map to factory creation process
	s.itemsMap = map[int]items.ItemType{
		0: gameitems.CollectibleCoinType,
		1: gameitems.SignpostType,
	}

	// Set items position from tilemap
	itemsPos := s.tilemap.GetItemsPositionID()
	if len(itemsPos) > 0 {
		f := items.NewItemFactory(gameitems.InitItemMap())
		for _, i := range itemsPos {
			itemType, found := s.itemsMap[i.ID]
			if !found {
				log.Fatal(err)
			}

			item, err := f.Create(itemType, i.X, i.Y)
			if err != nil {
				log.Fatal(err)
			}
			s.space.AddBody(item)
		}
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
		float64(config.Get().ScreenWidth),
		float64(config.Get().ScreenHeight),
	)
	s.cam.SmoothType = kamera.SmoothDamp
	s.cam.ShakeEnabled = true

	// Init collisions bodies and touch trigger for endpoints
	endpointTrigger := physics.NewTouchTrigger(s.finishLevel, s.player)
	s.tilemap.CreateCollisionBodies(s.space, endpointTrigger)

	s.levelCompleted = false

	// Init sequence
	s.sequencePlayer = sequences.NewSequencePlayer(s.AppContext)
}

func (s *LevelsScene) GetTilemapWidth() int {
	if s.tilemap != nil && len(s.tilemap.Layers) > 0 {
		return s.tilemap.Layers[0].Width * s.tilemap.Tileheight
	}
	return config.Get().ScreenWidth
}

func (s *LevelsScene) GetTilemapHeight() int {
	if s.tilemap != nil && len(s.tilemap.Layers) > 0 {
		return s.tilemap.Layers[0].Height * s.tilemap.Tileheight
	}
	return config.Get().ScreenHeight
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

	s.sequencePlayer.Update()

	// Press S to play the sequence
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		sequence, err := sequences.NewSequenceFromJSON("assets/sequences/leandro.json")
		if err != nil {
			log.Fatal(err)
		}
		s.sequencePlayer.Play(sequence)
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
		switch sb := b.(type) {
		case actors.PlayerEntity:
			continue
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
			continue
		}
	}

	// Draw player based on camera
	if img := s.player.Image(); img != nil {
		opts := *s.player.ImageOptions()
		s.cam.Draw(img, &opts, screen)
	}

	s.DrawHUD(screen)
}

func (s *LevelsScene) OnFinish() {
	// TODO: Should reset actor manager?
	s.audiomanager.PauseMusic(bgSound)
}

func (s *LevelsScene) finishLevel() {
	if s.levelCompleted {
		return
	}

	s.levelCompleted = true
	s.AppContext.SceneManager.NavigateTo(SceneSummary, transition.NewFader())
}

func createPlayer(space *physics.Space, appContext *core.AppContext) (actors.PlayerEntity, error) {
	player, err := actors.NewPlayer(actors.Platform, appContext)
	if err != nil {
		return nil, err
	}
	pos := player.Position()
	player.SetCollisionArea(physics.NewRect(pos.Min.X, pos.Min.Y, pos.Dx(), pos.Dy()))
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
