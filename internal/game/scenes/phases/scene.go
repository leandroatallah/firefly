package gamescenephases

import (
	"image"
	"image/color"
	"log"
	"math"
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
	enginecamera "github.com/leandroatallah/firefly/internal/engine/render/camera"
	"github.com/leandroatallah/firefly/internal/engine/render/screenutil"
	enginevfx "github.com/leandroatallah/firefly/internal/engine/render/vfx"
	"github.com/leandroatallah/firefly/internal/engine/scene"
	"github.com/leandroatallah/firefly/internal/engine/scene/pause"
	"github.com/leandroatallah/firefly/internal/engine/scene/phases"
	"github.com/leandroatallah/firefly/internal/engine/scene/transition"
	"github.com/leandroatallah/firefly/internal/engine/sequences"
	"github.com/leandroatallah/firefly/internal/engine/utils"
	"github.com/leandroatallah/firefly/internal/engine/utils/timing"
	gameenemies "github.com/leandroatallah/firefly/internal/game/entity/actors/enemies"
	gamenpcs "github.com/leandroatallah/firefly/internal/game/entity/actors/npcs"
	gamestates "github.com/leandroatallah/firefly/internal/game/entity/actors/states"
	gameitems "github.com/leandroatallah/firefly/internal/game/entity/items"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
	gamecamera "github.com/leandroatallah/firefly/internal/game/render/camera"
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
	completionTrigger utils.DelayTrigger

	// UI effects
	ShowDrawScreenFlash int

	screenFlipper  *scene.ScreenFlipper
	sequencePlayer sequencestypes.Player
	pauseScreen    *pause.PauseScreen

	// Game-layer camera controller with vertical-only-upward constraint
	gameCamera *gamecamera.Controller

	vignette *enginevfx.Vignette

	death deathSequence
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
		vignette:     enginevfx.NewVignette(),
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

// checkPlayerFallDeath checks if the player fell out of camera view and triggers death.
func (s *PhasesScene) checkPlayerFallDeath() {
	if s.gameCamera == nil || s.player == nil {
		return
	}

	// Don't trigger death during active death sequence
	if s.death.active {
		return
	}

	// Get camera center and player position
	_, camY := s.gameCamera.Base().GetActualCenter()
	_, playerY := s.player.GetPositionMin()

	// Calculate bottom of camera viewport
	// Camera center Y is the center of screen, so bottom is center + half screen height
	cameraBottom := camY + s.gameCamera.Height()/2

	// Check if player's top is below camera bottom (player fell out of view)
	playerTop := float64(playerY)
	if playerTop > cameraBottom {
		s.startDeathSequence()
	}
}

// clampCameraTarget clamps camera target position to camera bounds.
func (s *PhasesScene) clampCameraTarget(x, y float64, bounds *image.Rectangle) (float64, float64) {
	if bounds == nil {
		return x, y
	}

	halfW := s.gameCamera.Width() / 2
	halfH := s.gameCamera.Height() / 2
	minX := float64(bounds.Min.X) + halfW
	maxX := float64(bounds.Max.X) - halfW
	minY := float64(bounds.Min.Y) + halfH
	maxY := float64(bounds.Max.Y) - halfH

	if x < minX {
		x = minX
	}
	if x > maxX {
		x = maxX
	}
	if y < minY {
		y = minY
	}
	if y > maxY {
		y = maxY
	}
	return x, y
}

// startDeathSequence initiates the death sequence: player dies, camera moves to start position.
func (s *PhasesScene) startDeathSequence() {
	if s.player == nil || !s.Tilemap().HasPlayerStartPosition() {
		return
	}

	// Get player death position (for explosion VFX)
	deathX, deathY := s.player.GetPositionMin()
	deathW, deathH := s.player.GetShape().Width(), s.player.GetShape().Height()
	deathCenterX := float64(deathX) + float64(deathW)/2
	deathCenterY := float64(deathY) + float64(deathH)/2

	// Get player start position from tilemap
	startX, startY, found := s.Tilemap().GetPlayerStartPosition()
	if !found {
		return
	}

	// Call OnDie to set health to 0 and transition to Dying state
	s.player.GetCharacter().SetNewStateFatal(gamestates.Dying)

	// Spawn explosion VFX at death location
	if s.AppContext().VFX != nil {
		s.AppContext().VFX.SpawnDeathExplosion(deathCenterX, deathCenterY, 50)
	}

	// Disable player movement during death sequence
	s.player.SetImmobile(true)

	// Store player start position for teleport
	s.death.playerStartX = float64(startX)
	s.death.playerStartY = float64(startY)

	// Get current camera position
	baseCam := s.BaseCamera()
	s.death.cameraStartX, s.death.cameraStartY = baseCam.GetActualCenter()

	// Calculate camera target (clamped to bounds)
	s.death.cameraTargetX, s.death.cameraTargetY = s.clampCameraTarget(float64(startX), float64(startY), baseCam.Bounds())

	// Calculate dynamic duration based on distance
	dx := s.death.cameraTargetX - s.death.cameraStartX
	dy := s.death.cameraTargetY - s.death.cameraStartY
	distance := math.Sqrt(dx*dx + dy*dy)

	// Base duration + distance factor, capped between 30-60 frames
	s.death.duration = int(30 + distance*0.15)
	if s.death.duration > 60 {
		s.death.duration = 60
	}
	if s.death.duration < 30 {
		s.death.duration = 30
	}

	// Wait 1.5 seconds after explosion before moving camera
	s.death.waitDuration = 90 // 1.5 seconds at 60fps
	s.death.waitTimer = s.death.waitDuration
	s.death.phase = deathSequencePhaseWaiting
	s.death.timer = 0
	s.death.active = true
}

// updateDeathSequence updates the death sequence state machine.
func (s *PhasesScene) updateDeathSequence() {
	switch s.death.phase {
	case deathSequencePhaseWaiting:
		s.updateDeathWaitPhase()
	case deathSequencePhaseMoving:
		s.updateDeathCameraPhase()
	}
}

// updateDeathWaitPhase handles the waiting phase (explosion animation).
func (s *PhasesScene) updateDeathWaitPhase() {
	s.death.waitTimer--
	if s.death.waitTimer <= 0 {
		// Wait complete - teleport player to start position
		if s.player != nil {
			actorHeight := s.player.Position().Dy()
			s.player.SetPosition(int(s.death.playerStartX), int(s.death.playerStartY)-actorHeight)
		}
		// Start camera movement
		s.death.phase = deathSequencePhaseMoving
		s.death.timer = 0
	}
}

// updateDeathCameraPhase handles the camera movement phase.
func (s *PhasesScene) updateDeathCameraPhase() {
	s.death.timer++
	progress := float64(s.death.timer) / float64(s.death.duration)
	if progress >= 1.0 {
		// Camera reached start position - transition player to Rising and unfreeze
		if s.player != nil {
			s.player.GetCharacter().SetNewStateFatal(gamestates.Rising)
			s.player.SetImmobile(false)
		}
		// Snap camera to final position
		baseCam := s.BaseCamera()
		baseCam.SetCenter(s.death.cameraTargetX, s.death.cameraTargetY)

		// Sync game camera's lastCameraY to prevent oscillation
		_, camY := baseCam.GetActualCenter()
		s.gameCamera.SetLastCameraY(camY)

		s.death.active = false
	} else {
		// Ease-out quadratic interpolation: starts fast, slows down
		easedProgress := progress * (2 - progress)
		currentX := s.death.cameraStartX + (s.death.cameraTargetX-s.death.cameraStartX)*easedProgress
		currentY := s.death.cameraStartY + (s.death.cameraTargetY-s.death.cameraStartY)*easedProgress
		baseCam := s.BaseCamera()
		baseCam.SetCenter(currentX, currentY)
	}
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

// BaseCamera returns the underlying engine camera controller (bypasses game-layer constraint).
func (s *PhasesScene) BaseCamera() *enginecamera.Controller {
	if s.gameCamera != nil {
		return s.gameCamera.Base()
	}
	return s.TilemapScene.Camera()
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

	// Check if player fell out of camera view
	if s.hasPlayer {
		s.checkPlayerFallDeath()
	}

	// Update death sequence if active
	if s.death.active {
		s.updateDeathSequence()
	}

	// Update navigation triggers
	s.completionTrigger.Update()

	// Check goal completion
	if s.goal != nil && s.goal.IsCompleted() && !s.completionTrigger.IsEnabled() {
		s.goal.OnCompletion()
	}

	if config.Get().CamDebug {
		s.gameCamera.CamDebug()
	}

	if s.completionTrigger.Trigger() {
		s.AppContext().CompleteCurrentPhase(transition.NewFader(0, config.Get().FadeVisibleDuration), true)
	}

	// Update camera (use game-layer camera with vertical-only-upward constraint)
	// Skip during death sequence - camera is controlled by death sequence
	if !s.death.active {
		if s.gameCamera != nil {
			s.gameCamera.Update()
		} else {
			// Fallback to base camera update
			s.TilemapScene.Camera().Update()
		}
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

// EnableVignetteDarkness enables the world darkness overlay with the given radius in screen pixels.
// The effect follows the player and is applied after world rendering (so UI remains visible).
func (s *PhasesScene) EnableVignetteDarkness(radiusPx float64) {
	if s.vignette == nil {
		s.vignette = enginevfx.NewVignette()
	}
	s.vignette.Enable(radiusPx)
}

// DisableVignetteDarkness disables the world darkness overlay.
func (s *PhasesScene) DisableVignetteDarkness() {
	if s.vignette == nil {
		return
	}
	s.vignette.Disable()
}

func (s *PhasesScene) endpointTrigger(eventID string) {
	if !s.hasPlayer {
		return
	}

	switch eventID {
	case "SPIKE":
		s.startDeathSequence()
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
