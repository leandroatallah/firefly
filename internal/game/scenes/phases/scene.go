package gamescenephases

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/assets/font"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	sequencestypes "github.com/boilerplate/ebiten-template/internal/engine/contracts/sequences"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/enemies"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/npcs"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/platformer"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/items"
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
	enginecamera "github.com/boilerplate/ebiten-template/internal/engine/render/camera"
	"github.com/boilerplate/ebiten-template/internal/engine/render/screenutil"
	enginevfx "github.com/boilerplate/ebiten-template/internal/engine/render/vfx"
	"github.com/boilerplate/ebiten-template/internal/engine/scene"
	"github.com/boilerplate/ebiten-template/internal/engine/scene/pause"
	"github.com/boilerplate/ebiten-template/internal/engine/scene/phases"
	"github.com/boilerplate/ebiten-template/internal/engine/scene/transition"
	"github.com/boilerplate/ebiten-template/internal/engine/sequences"
	"github.com/boilerplate/ebiten-template/internal/engine/ui/menu"
	"github.com/boilerplate/ebiten-template/internal/engine/utils"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/timing"
	gameenemies "github.com/boilerplate/ebiten-template/internal/game/entity/actors/enemies"
	gamenpcs "github.com/boilerplate/ebiten-template/internal/game/entity/actors/npcs"
	gamestates "github.com/boilerplate/ebiten-template/internal/game/entity/actors/states"
	gameitems "github.com/boilerplate/ebiten-template/internal/game/entity/items"
	gameentitytypes "github.com/boilerplate/ebiten-template/internal/game/entity/types"
	gamecamera "github.com/boilerplate/ebiten-template/internal/game/render/camera"
	scenestypes "github.com/boilerplate/ebiten-template/internal/game/scenes/types"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
	hasPlayer       bool
	goal            phases.Goal

	// Navigation triggers
	completionTrigger utils.DelayTrigger
	deathTrigger      utils.DelayTrigger

	// UI effects
	ShowDrawScreenFlash int

	screenFlipper  *scene.ScreenFlipper
	sequencePlayer sequencestypes.Player
	pauseScreen    *pause.PauseScreen
	pauseMenu      *menu.Menu

	// Game-layer camera controller with vertical-only-upward constraint
	gameCamera *gamecamera.Controller

	vignette *enginevfx.Vignette

	death         deathSequence
	bullets       []*gamestates.Bullet
	bulletImg     *ebiten.Image
	bulletCounter int
}

// SpawnBullet implements body.Shooter interface for shooting skill.
func (s *PhasesScene) SpawnBullet(x16, y16, vx16, vy16 int, owner interface{}) {
	baseBody := bodyphysics.NewBody(bodyphysics.NewRect(0, 0, 2, 1))
	baseBody.SetPosition16(x16, y16)

	s.bulletCounter++
	bulletID := fmt.Sprintf("bullet_%d", s.bulletCounter)
	baseBody.SetID(bulletID)

	movableBody := bodyphysics.NewMovableBody(baseBody)
	collidableBody := bodyphysics.NewCollidableBody(baseBody)
	collidableBody.SetOwner(owner)

	bullet := gamestates.NewBullet(movableBody, collidableBody, s.PhysicsSpace(), vx16, vy16)
	collidableBody.SetTouchable(bullet)
	s.PhysicsSpace().AddBody(collidableBody)
	s.bullets = append(s.bullets, bullet)
}

func NewPhasesScene(ctx *app.AppContext) *PhasesScene {
	tilemapScene := scene.NewTilemapScene(ctx)

	bulletImg := ebiten.NewImage(2, 1)
	bulletImg.Fill(color.RGBA{255, 255, 255, 255})

	scene := &PhasesScene{
		TilemapScene: tilemapScene,
		mainText:     ctx.Font,
		bodyCounter:  &BodyCounter{},
		vignette:     enginevfx.NewVignette(),
		bulletImg:    bulletImg,
	}
	scene.SetAppContext(ctx)

	subscribeEvents(ctx, scene)

	return scene
}

func (s *PhasesScene) OnStart() {
	s.TilemapScene.OnStart()
	s.count = 0
	s.death.active = false

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
		ctx.ActorManager.RegisterPrimary(s.player)
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

	// Create pause menu
	s.pauseMenu = menu.NewMenu()
	s.pauseMenu.SetFontSize(8)
	// Resume
	s.pauseMenu.AddItem("", func() {
		s.pauseScreen.Toggle()
	})
	// Exit to Menu
	s.pauseMenu.AddItem("", func() {
		s.pauseScreen.Toggle()
		s.freezeAllActors()
		ctx.AudioManager.PauseCurrentMusic()
		ctx.SceneManager.NavigateTo(
			scenestypes.SceneMenu,
			transition.NewFader(0, time.Duration(2*time.Second)),
			true,
		)
	})

	// Set up pause menu navigation and selection callbacks
	s.pauseMenu.SetOnNavigate(func() {
		ctx.AudioManager.PlaySound("assets/audio/Menu_Click.ogg")
	})
	s.pauseMenu.SetOnSelect(func() {
		ctx.AudioManager.PlaySound("assets/audio/Menu_Select2.ogg")
	})

	s.pauseScreen.SetMenu(s.pauseMenu)
	s.pauseScreen.SetFont(s.mainText)
	s.refreshPauseMenuLabels()

	s.pauseScreen.SetOnStart(func(p *pause.PauseScreen) {
		p.SetMenu(s.pauseMenu) // Reset to main pause menu on open
		if ctx.AudioManager != nil {
			ctx.AudioManager.PauseCurrentMusic()
		}
	})
	s.pauseScreen.SetOnFinish(func(p *pause.PauseScreen) {
		if ctx.AudioManager != nil {
			ctx.AudioManager.ResumeCurrentMusic()
		}
	})

	phase, err := ctx.PhaseManager.GetCurrentPhase()
	if err == nil && phase.SequencePath != "" {
		s.sequencePlayer = sequences.NewSequencePlayer(ctx)
		s.allowPause = phase.GoalType != SequenceGoalType
		s.sequencePlayer.PlaySequence(phase.SequencePath)
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

// startDeathSequence triggers the death VFX and navigates to PhaseRebootScene,
// which fades to black and NavigateBack to restart the phase via OnStart.
func (s *PhasesScene) startDeathSequence() {
	if s.death.active {
		return
	}

	if s.player == nil {
		return
	}

	s.death.active = true

	// Spawn explosion VFX at player position
	if s.AppContext().VFX != nil {
		deathX, deathY := s.player.GetPositionMin()
		deathW, deathH := s.player.GetShape().Width(), s.player.GetShape().Height()
		s.AppContext().VFX.SpawnDeathExplosion(
			float64(deathX)+float64(deathW)/2,
			float64(deathY)+float64(deathH)/2,
			50,
		)
	}

	s.player.GetCharacter().SetNewStateFatal(gamestates.Dying)
	s.player.SetImmobile(true)

	s.deathTrigger.Enable(timing.FromDuration(time.Second))
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

	// Check if player died (from any cause) and death sequence hasn't started
	if s.hasPlayer && !s.death.active && (s.player.State() == gamestates.Dying || s.player.State() == gamestates.Dead) {
		s.startDeathSequence()
	}

	// Update navigation triggers
	s.completionTrigger.Update()
	s.deathTrigger.Update()

	if s.deathTrigger.Trigger() {
		s.AppContext().SceneManager.NavigateTo(
			scenestypes.ScenePhaseReboot,
			transition.NewFader(0, config.Get().FadeVisibleDuration),
			false,
		)
	}

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
	if s.gameCamera != nil {
		s.gameCamera.Update()
	} else {
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

	// Update bullets
	for _, bullet := range s.bullets {
		bullet.Update()
	}

	// Check for trigger collisions (spikes, endpoints, etc.)
	// This ensures non-blocking trigger bodies call their OnTouch() callbacks
	// Must happen after actor updates to detect collisions at current positions
	if s.hasPlayer && s.player != nil {
		space.ResolveCollisions(s.player)
	}

	// Remove bodies queued for removal
	space.ProcessRemovals()

	// Clean up removed bullets by checking if body still exists in space
	activeBullets := make([]*gamestates.Bullet, 0, len(s.bullets))
	bodies := space.Bodies()
	for _, bullet := range s.bullets {
		found := false
		for _, b := range bodies {
			if b == bullet.Body() {
				found = true
				break
			}
		}
		if found {
			activeBullets = append(activeBullets, bullet)
		}
	}
	s.bullets = activeBullets

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

	// Draw bullets
	for _, bullet := range s.bullets {
		opts := &ebiten.DrawImageOptions{}
		x, y := bullet.Body().GetPositionMin()
		opts.GeoM.Translate(float64(x), float64(y))
		s.Camera().Draw(s.bulletImg, opts, screen)
	}

	if s.ShowDrawScreenFlash > 0 {
		screenutil.DrawScreenFlash(screen)
		s.ShowDrawScreenFlash--
	}

	if s.AppContext().VFX != nil {
		s.AppContext().VFX.Draw(screen, s.gameCamera.Base())
	}

	// Darkness vignette should cover world (including VFX) but not UI.
	if s.vignette != nil && s.hasPlayer && s.player != nil && s.gameCamera != nil {
		s.vignette.Draw(screen, s.gameCamera.Base(), s.player)
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

// TriggerScreenFlash triggers a white screen flash effect for feedback.
func (s *PhasesScene) TriggerScreenFlash() {
	s.ShowDrawScreenFlash = 2
}

func (s *PhasesScene) endpointTrigger(eventType string) {
	if !s.hasPlayer {
		return
	}

	// Prevent multiple triggers (e.g., from continuous spike collision)
	if s.death.active {
		return
	}

	switch eventType {
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

func (s *PhasesScene) refreshPauseMenuLabels() {
	i18n := s.AppContext().I18n

	if s.pauseMenu != nil {
		s.pauseMenu.UpdateItemLabel(0, i18n.T("menu.start"))
		s.pauseMenu.UpdateItemLabel(1, i18n.T("menu.exit"))
	}
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

	// Draw menu on top of background
	if menu := s.pauseScreen.Menu(); menu != nil {
		menu.Draw(screen, s.pauseScreen.Font(), cfg.ScreenWidth/2, cfg.ScreenHeight/2)
	}
}
