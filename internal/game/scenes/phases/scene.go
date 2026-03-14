package gamescenephases

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	sequencestypes "github.com/leandroatallah/firefly/internal/engine/contracts/sequences"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/enemies"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/npcs"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/platformer"
	"github.com/leandroatallah/firefly/internal/engine/entity/items"
	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
	"github.com/leandroatallah/firefly/internal/engine/render/screenutil"
	"github.com/leandroatallah/firefly/internal/engine/scene"
	"github.com/leandroatallah/firefly/internal/engine/scene/pause"
	"github.com/leandroatallah/firefly/internal/engine/scene/phases"
	"github.com/leandroatallah/firefly/internal/engine/scene/transition"
	"github.com/leandroatallah/firefly/internal/engine/sequences"
	"github.com/leandroatallah/firefly/internal/engine/utils"
	"github.com/leandroatallah/firefly/internal/engine/utils/timing"
	gameenemies "github.com/leandroatallah/firefly/internal/game/entity/actors/enemies"
	gamenpcs "github.com/leandroatallah/firefly/internal/game/entity/actors/npcs"
	gameitems "github.com/leandroatallah/firefly/internal/game/entity/items"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
	gamecamera "github.com/leandroatallah/firefly/internal/game/render/camera"
	scenestypes "github.com/leandroatallah/firefly/internal/game/scenes/types"
)

const (
	bgSound = "assets/audio/Goblins_Den_Regular.ogg"
)

type PhasesScene struct {
	*scene.TilemapScene

	count       int
	player      platformer.PlatformerActorEntity
	mainText    *font.FontText
	bodyCounter *BodyCounter
	allowPause  bool

	// Complete phase
	reachedEndpoint bool
	hasEndpoints    bool
	hasPlayer       bool
	goal            phases.Goal

	// Navigation triggers
	rebootTrigger     utils.DelayTrigger
	completionTrigger utils.DelayTrigger

	// UI effects
	ShowDrawScreenFlash int

	screenFlipper *scene.ScreenFlipper
	// sequencePlayer *sequences.SequencePlayer
	sequencePlayer sequencestypes.Player
	pauseScreen    *pause.PauseScreen

	// Game-layer camera controller with vertical-only-upward constraint
	gameCamera *gamecamera.Controller
}

func NewPhasesScene(ctx *app.AppContext) *PhasesScene {
	mainText, err := font.NewFontText(config.Get().MainFontFace)
	if err != nil {
		log.Fatal(err)
	}
	tilemapScene := scene.NewTilemapScene(ctx)
	scene := &PhasesScene{
		TilemapScene: tilemapScene,
		mainText:     mainText,
		bodyCounter:  &BodyCounter{},
	}
	scene.SetAppContext(ctx)

	subscribeEvents(ctx, scene)

	return scene
}

func (s *PhasesScene) OnStart() {
	s.TilemapScene.OnStart()
	s.count = 0

	ctx := s.AppContext()

	// Check if player should be created (based on PlayerStart layer existence)
	s.hasPlayer = s.Tilemap().HasPlayerStartPosition()

	if s.hasPlayer {
		// Create player and register to space and context
		p, err := createPlayer(ctx, gameentitytypes.ClimberPlayerType)
		if err != nil {
			log.Fatal(err)
		}
		s.player = p
		ctx.ActorManager.Register(s.player)
		s.PhysicsSpace().AddBody(s.player)

		// Optionally block input for the current player of this phase
		if phase, err := ctx.PhaseManager.GetCurrentPhase(); err == nil && phase.BlockPlayerMovement {
			if p, ok := ctx.ActorManager.GetPlayer(); ok {
				p.BlockMovement()
			}
		}
	}

	s.initTilemap()

	// After init bodies, set body counter
	s.bodyCounter.setBodyCounter(s.PhysicsSpace())

	s.PhysicsSpace().Bodies()

	// Init collision bodies (obstacles and endpoints) - always created regardless of player
	s.Tilemap().CreateCollisionBodies(s.PhysicsSpace(), func(id string) body.Touchable {
		return bodyphysics.NewTouchTrigger(func() {
			s.endpointTrigger(id)
		}, s.player)
	})

	if s.hasPlayer {
		s.SetCameraConfig(scene.CameraConfig{Mode: scene.CameraModeFollow})
		// Wrap base camera with game-layer controller that adds vertical-only-upward constraint
		s.gameCamera = gamecamera.NewController(s.TilemapScene.Camera())
		s.gameCamera.SetFollowTarget(s.player)

		s.screenFlipper = scene.NewScreenFlipper(s.gameCamera.Base(), s.player, s.Tilemap(), ctx)
		tileWidth := s.Tilemap().Tilewidth
		s.screenFlipper.PlayerPushDistance = float64(tileWidth / 2)
		s.screenFlipper.FlipStrategy = func(dx, dy int) scene.FlipType {
			if dy != 0 {
				return scene.FlipTypeInstant
			}
			return scene.FlipTypeSmooth
		}
		s.screenFlipper.OnFlipStart = func() {
			s.player.SetImmobile(true)
		}
		s.screenFlipper.OnFlipFinish = func() {
			s.player.SetImmobile(false)
		}
		s.screenFlipper.SnapToCurrentRoom()
	} else {
		// No player: set camera to fixed mode at CameraStart position or top-left
		s.SetCameraConfig(scene.CameraConfig{Mode: scene.CameraModeFixed})

		if x, y, found := s.Tilemap().GetCameraStartPosition(); found {
			s.Camera().SetPositionTopLeft(float64(x), float64(y))
		} else {
			s.Camera().SetPositionTopLeft(0, 0)
		}
	}

	s.pauseScreen = pause.NewPauseScreen(ebiten.KeyEnter, 250*time.Millisecond)

	phase, err := ctx.PhaseManager.GetCurrentPhase()
	if err == nil && phase.SequencePath != "" {
		s.sequencePlayer = sequences.NewSequencePlayer(ctx)
		s.allowPause = phase.GoalType != SequenceGoalType
		seq, err := sequences.NewSequenceFromJSON(phase.SequencePath)
		if err != nil {
			log.Printf("Failed to load sequence: %v", err)
		} else {
			s.sequencePlayer.Play(seq)
		}
	}

	s.initGoal()
}

func (s *PhasesScene) initGoal() {
	phase, _ := s.AppContext().PhaseManager.GetCurrentPhase()
	switch phase.GoalType {
	case ReactEndpointType:
		s.goal = &ReachEndpointGoal{scene: s}
	case SequenceGoalType:
		s.goal = &phases.SequenceGoal{
			Player:         s.sequencePlayer,
			OnCompleteFunc: s.defaultCompletion,
		}
	case NoGoalType:
		s.goal = &phases.NoGoal{}
	default:
		s.goal = &phases.NoGoal{}
	}
}

func (s *PhasesScene) freezeAllActors() {
	if s.AppContext().ActorManager != nil {
		s.AppContext().ActorManager.ForEach(func(actor actors.ActorEntity) {
			actor.SetImmobile(true)
			actor.SetFreeze(true)
		})
	}
}

func (s *PhasesScene) defaultCompletion() {
	s.completionTrigger.Enable(timing.FromDuration(time.Second))
}

// Camera returns the game-layer camera controller with vertical-only-upward constraint.
// Falls back to base camera if gameCamera is not set (e.g., when hasPlayer is false).
func (s *PhasesScene) Camera() *gamecamera.Controller {
	if s.gameCamera != nil {
		return s.gameCamera
	}
	// Return a wrapper for base camera when gameCamera is not set
	return gamecamera.NewController(s.TilemapScene.Camera())
}

func (s *PhasesScene) Update() error {
	if s.pauseScreen != nil && s.canPause() {
		s.pauseScreen.Update()
		if s.pauseScreen.IsPaused() {
			return nil
		}
	}

	if s.sequencePlayer != nil {
		s.sequencePlayer.Update()
	}

	if s.AppContext().VFX != nil {
		s.AppContext().VFX.Update()
	}

	if s.screenFlipper != nil {
		s.screenFlipper.Update()
		if s.screenFlipper.IsFlipping() {
			return nil
		}
	}

	// Update navigation triggers
	s.rebootTrigger.Update()
	s.completionTrigger.Update()

	// Check goal completion
	if s.goal != nil && s.goal.IsCompleted() && !s.completionTrigger.IsEnabled() {
		s.goal.OnCompletion()
	}

	if config.Get().CamDebug {
		s.gameCamera.CamDebug()
	}

	if s.rebootTrigger.Trigger() {
		s.AppContext().SceneManager.NavigateTo(
			scenestypes.ScenePhaseReboot,
			transition.NewFader(0, 0),
			true,
		)
		return nil
	}

	if s.completionTrigger.Trigger() {
		s.AppContext().CompleteCurrentPhase(transition.NewFader(0, config.Get().FadeVisibleDuration), true)
	}

	// Update camera (use game-layer camera with vertical-only-upward constraint)
	if s.gameCamera != nil {
		s.gameCamera.Update()
	} else {
		// Fallback to base camera update
		s.TilemapScene.Camera().Update()
	}
	// Call BaseScene.Update directly for Schedule handling (skip TilemapScene.Update to avoid double camera update)
	if err := s.BaseScene.Update(); err != nil {
		return err
	}

	s.count++

	// Execute bodies updates
	space := s.PhysicsSpace()
	for _, i := range space.Bodies() {
		switch b := i.(type) {
		// ActorEntity case should came first. It can be confused with body.Obstacle
		case platformer.PlatformerActorEntity:
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

	// Remove bodies queued for removal
	space.ProcessRemovals()

	return nil
}

func (s *PhasesScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 0xff}) // force black

	// Get tilemap image and draw based on camera
	tilemap, err := s.Tilemap().Image(screen)
	if err != nil {
		log.Fatal(err)
	}
	s.Camera().Draw(tilemap, s.Tilemap().ImageOptions(), screen)

	// Draw bodies based on camera
	space := s.PhysicsSpace()
	for _, b := range space.Bodies() {
		switch sb := b.(type) {
		case platformer.PlatformerActorEntity:
			opts := sb.ImageOptions()
			sb.UpdateImageOptions()
			s.Camera().Draw(sb.Image(), opts, screen)
			if config.Get().CollisionBox {
				s.Camera().DrawCollisionBox(screen, sb)
			}
		case items.Item:
			if sb.IsRemoved() {
				continue
			}
			opts := sb.ImageOptions()
			sb.UpdateImageOptions()
			s.Camera().Draw(sb.Image(), opts, screen)
			if config.Get().CollisionBox {
				s.Camera().DrawCollisionBox(screen, sb)
			}
		case body.Obstacle:
			if config.Get().CollisionBox {
				s.Camera().DrawCollisionBox(screen, sb)
			}
		}
	}

	if s.ShowDrawScreenFlash > 0 {
		screenutil.DrawScreenFlash(screen)
		s.ShowDrawScreenFlash--
	}

	if s.AppContext().VFX != nil {
		s.AppContext().VFX.Draw(screen, s.gameCamera.Base())
	}

	if s.pauseScreen.IsPaused() {
		s.drawPause(screen)
	}
}

func (s *PhasesScene) Reboot() {
	s.ShowDrawScreenFlash = timing.FromDuration(67 * time.Millisecond) // 4 frames
	s.rebootTrigger.Enable(timing.FromDuration(1 * time.Second))       // 60 frames
}

func (s *PhasesScene) OnFinish() {
	s.TilemapScene.OnFinish()
	// Ensure we remove any movement block applied at phase start (for whichever actor is the current player)
	if s.hasPlayer {
		if p, ok := s.AppContext().ActorManager.GetPlayer(); ok {
			p.UnblockMovement()
		}
		s.AppContext().ActorManager.Unregister(s.player)
	}
}

func (s *PhasesScene) endpointTrigger(eventID string) {
	if !s.hasPlayer {
		return
	}

	switch eventID {
	case "SPIKE":
		s.player.OnDie()
		return
	case "CUTSCENE":
		// TODO: Implement this
	}

	s.reachedEndpoint = true
}

func (s *PhasesScene) initTilemap() {
	// Set items position from tilemap
	f := items.NewItemFactory(gameitems.InitItemMap(s.AppContext()))
	scene.InitItems(s.TilemapScene, f)

	// Set enemies position from tilemap
	enemyFactory := enemies.NewEnemyFactory(gameenemies.InitEnemyMap(s.AppContext()))
	scene.InitEnemies(s.TilemapScene, enemyFactory)

	// Set NPCs position from tilemap
	npcFactory := npcs.NewNpcFactory(gamenpcs.InitNpcMap(s.AppContext()))
	scene.InitNPCs(s.TilemapScene, npcFactory)

	if s.hasPlayer {
		s.SetPlayerStartPosition(s.player)
	}
}

func (s *PhasesScene) canPause() bool {
	return s.allowPause && !s.sequencePlayer.IsPlaying()
}

func (s *PhasesScene) drawPause(screen *ebiten.Image) {
	if !s.canPause() || s.pauseScreen == nil || !s.pauseScreen.IsPaused() {
		return
	}

	cfg := config.Get()
	for x := 0; x < cfg.ScreenWidth; x++ {
		for y := 0; y < cfg.ScreenWidth; y++ {
			if x%2 == 0 && y%2 == 0 {
				vector.DrawFilledRect(screen, float32(x), float32(y), 1, 1, color.Black, false)
			}
		}
	}

	speed := 10
	initialW, initialH := cfg.ScreenWidth/4, cfg.ScreenHeight/4
	w := max(min(initialW+s.pauseScreen.Count()*speed, cfg.ScreenWidth/2), 1)
	h := max(min(initialH+s.pauseScreen.Count()*speed, cfg.ScreenHeight/2), 1)
	container := ebiten.NewImage(w, h)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(cfg.ScreenWidth)/2, float64(cfg.ScreenHeight)/2)
	op.GeoM.Translate(-float64(w/2), -float64(h/2))
	container.Fill(color.Black)
	screen.DrawImage(container, op)
}
