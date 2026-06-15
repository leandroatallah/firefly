// Package platformerphasescene provides a genre-reusable platformer phase scene
// that assembles the engine scene base with platformer-specific actor handling.
package platformerphasescene

import (
	"image"
	"image/color"
	"log"
	"time"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/navigation"
	sequencestypes "github.com/boilerplate/ebiten-template/internal/engine/contracts/sequences"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/tilemaplayer"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	actorevents "github.com/boilerplate/ebiten-template/internal/engine/entity/actors/events"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/items"
	"github.com/boilerplate/ebiten-template/internal/engine/event"
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
	enginecamera "github.com/boilerplate/ebiten-template/internal/engine/render/camera"
	"github.com/boilerplate/ebiten-template/internal/engine/render/draworder"
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
	"github.com/boilerplate/ebiten-template/internal/kit/actors/platformer"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Player is the minimal interface the scene requires from a platformer player.
// It is a subset of platformer.PlatformerActorEntity that mock and production
// types both satisfy.
type Player interface {
	body.Collidable
	body.Drawable

	GetPositionMin() (x, y int)
	GetShape() body.Shape

	State() actors.ActorStateEnum
	SetImmobile(bool)
	GetCharacter() *actors.Character
	Update(body.BodiesSpace) error
}

// PlatformerPhaseScene is the kit-level platformer scene. It manages a
// platformer player, fall-death detection, and a screen flipper.
// Game layers embed or compose this type to wire concrete factories.
type PlatformerPhaseScene struct {
	camera       *enginecamera.Controller
	space        body.BodiesSpace
	hasPlayer    bool
	cameraMode   scene.CameraMode
	screenWidth  float64
	screenHeight float64

	// Player — holds the minimal interface; production code uses platformer.PlatformerActorEntity.
	player     Player
	dyingState actors.ActorStateEnum
	deadState  actors.ActorStateEnum

	// Death sequence
	deathActive bool

	// OnDeathStarted is called after the death sequence activates. Game layers
	// set this to spawn VFX and enable navigation triggers.
	OnDeathStarted func()

	// Screen flipper callbacks. Set up by OnStart when hasFlipper=true.
	onFlipStart  func()
	onFlipFinish func()
	hasFlipper   bool

	// Hook for overriding SetNewStateFatal calls (used in tests for assertion).
	setNewStateFatalHook func(actors.ActorStateEnum)

	// Debug hook invoked at the end of Draw.
	debugDrawHook func(*ebiten.Image)

	// Optional game-layer overrides.
	customGoal      phases.Goal
	endpointHandler func(id string) bool
	onStarted       func()

	// VFX vignette (may be nil)
	vfx        *enginevfx.Vignette
	flashCount int

	// Full-loop fields (nil when created via NewForTest)
	tilemapScene      *scene.TilemapScene
	appCtx            *app.AppContext
	goal              phases.Goal
	sequencePlayer    sequencestypes.Player
	allowPause        bool
	pauseScreen       *pause.PauseScreen
	pauseMenu         *menu.Menu
	screenFlipper     *scene.ScreenFlipper
	completionTrigger utils.DelayTrigger
	deathTrigger      utils.DelayTrigger
	rebootScene       navigation.SceneType
	menuScene         navigation.SceneType
	playerFactory     func(*app.AppContext) (Player, error)
	initActors        func(*scene.TilemapScene)
	count             int
}

// New creates a PlatformerPhaseScene for production use. dyingState and deadState
// are the genre-specific actor state values used to trigger/detect death.
func New(
	cam *enginecamera.Controller,
	space body.BodiesSpace,
	sw, sh float64,
	dyingState, deadState actors.ActorStateEnum,
) *PlatformerPhaseScene {
	return newScene(cam, space, sw, sh, dyingState, deadState)
}

// newScene creates a new PlatformerPhaseScene with the given camera and space.
func newScene(
	cam *enginecamera.Controller,
	space body.BodiesSpace,
	sw, sh float64,
	dyingState, deadState actors.ActorStateEnum,
) *PlatformerPhaseScene {
	return &PlatformerPhaseScene{
		camera:       cam,
		space:        space,
		screenWidth:  sw,
		screenHeight: sh,
		dyingState:   dyingState,
		deadState:    deadState,
		vfx:          enginevfx.NewVignette(),
	}
}

// Camera returns the underlying engine camera controller.
func (s *PlatformerPhaseScene) Camera() *enginecamera.Controller { return s.camera }

// SetAppContext implements navigation.Scene. The kit scene context is set at
// construction time via NewWithOptions; this method is provided to satisfy the
// navigation.Scene interface when the scene manager calls it after factory construction.
func (s *PlatformerPhaseScene) SetAppContext(_ any) {}

// OnStart initialises the scene for the current phase.
func (s *PlatformerPhaseScene) OnStart() {
	if s.appCtx != nil {
		s.fullOnStart()
		return
	}
	// minimal path (NewForTest) — keep existing behavior
	if s.hasPlayer {
		s.cameraMode = scene.CameraModeFollow
		s.camera.SetFollowing(true)
		s.camera.SetVerticalOnlyUpward(false)
		if s.player != nil {
			s.camera.SetFollowTarget(s.player)
			if s.hasFlipper {
				s.wireFlipperCallbacks(s.player)
			}
		}
	} else {
		s.cameraMode = scene.CameraModeFixed
		s.camera.SetFollowing(false)
	}
}

func (s *PlatformerPhaseScene) fullOnStart() {
	ts := scene.NewTilemapScene(s.appCtx)
	ts.OnStart()
	s.tilemapScene = ts
	s.camera = ts.Camera()
	s.space = s.appCtx.Space
	s.count = 0

	s.hasPlayer = ts.Tilemap().HasPlayerStartPosition()

	if s.hasPlayer {
		p, err := s.playerFactory(s.appCtx)
		if err != nil {
			log.Fatal(err)
		}
		s.player = p
		// Production players implement actors.ActorEntity; use type assertion.
		if ae, ok := any(p).(actors.ActorEntity); ok {
			s.appCtx.ActorManager.Register(ae)
			s.appCtx.ActorManager.RegisterPrimary(ae)
		}
		s.space.AddBody(p)
		s.OnDeathStarted = func() {
			if s.appCtx.VFX != nil {
				deathX, deathY := s.player.GetPositionMin()
				deathW, deathH := s.player.GetShape().Width(), s.player.GetShape().Height()
				s.appCtx.VFX.SpawnDeathExplosion(
					float64(deathX)+float64(deathW)/2,
					float64(deathY)+float64(deathH)/2,
					50,
				)
			}
			s.deathTrigger.Enable(timing.FromDuration(time.Second))
		}
		if phase, err := s.appCtx.PhaseManager.GetCurrentPhase(); err == nil && phase.BlockPlayerMovement {
			if pl, ok := s.appCtx.ActorManager.GetPlayer(); ok {
				pl.BlockMovement()
			}
		}
	}

	if s.initActors != nil {
		s.initActors(ts)
	}
	if s.hasPlayer {
		// Production players implement actors.ActorEntity for SetPlayerStartPosition.
		if ae, ok := any(s.player).(actors.ActorEntity); ok {
			ts.SetPlayerStartPosition(ae)
		}
	}

	ts.Tilemap().CreateCollisionBodies(s.space, func(id string) body.Touchable {
		return bodyphysics.NewTouchTrigger(func() {
			s.endpointTrigger(id)
		}, s.player)
	})

	if s.hasPlayer {
		ts.SetCameraConfig(scene.CameraConfig{Mode: scene.CameraModeFollow})
		s.camera.SetFollowing(true)
		s.camera.SetVerticalOnlyUpward(false)
		s.camera.SetFollowTarget(s.player)
		// Production players implement body.Movable.
		if mv, ok := any(s.player).(body.Movable); ok {
			s.screenFlipper = scene.NewScreenFlipper(s.camera, mv, ts.Tilemap(), s.appCtx)
			s.screenFlipper.PlayerPushDistance = float64(ts.Tilemap().Tilewidth / 2)
			s.screenFlipper.FlipStrategy = func(dx, dy int) scene.FlipType {
				if dy != 0 {
					return scene.FlipTypeInstant
				}
				return scene.FlipTypeSmooth
			}
			s.screenFlipper.OnFlipStart = func() { s.player.SetImmobile(true) }
			s.screenFlipper.OnFlipFinish = func() { s.player.SetImmobile(false) }
			s.screenFlipper.SnapToCurrentRoom()
		}
	} else {
		ts.SetCameraConfig(scene.CameraConfig{Mode: scene.CameraModeFixed})
		if x, y, found := ts.Tilemap().GetCameraStartPosition(); found {
			s.camera.SetPositionTopLeft(float64(x), float64(y))
		} else {
			s.camera.SetPositionTopLeft(0, 0)
		}
	}

	s.buildPauseScreen()
	s.buildSequencePlayer()
	s.initGoal()
	s.subscribeEvents()

	if s.onStarted != nil {
		s.onStarted()
	}
}

func (s *PlatformerPhaseScene) buildPauseScreen() {
	s.pauseScreen = pause.NewPauseScreen(ebiten.KeyEnter, 250*time.Millisecond)
	s.pauseMenu = menu.NewMenu()
	s.pauseMenu.SetFontSize(8)
	s.pauseMenu.AddItem("", func() { s.pauseScreen.Toggle() })
	s.pauseMenu.AddItem("", func() {
		s.pauseScreen.Toggle()
		s.freezeAllActors()
		s.appCtx.AudioManager.PauseCurrentMusic()
		s.appCtx.SceneManager.NavigateTo(
			s.menuScene,
			transition.NewFader(0, 2*time.Second),
			true,
		)
	})
	s.pauseMenu.SetOnNavigate(func() {
		s.appCtx.AudioManager.PlaySound("assets/audio/Menu_Click.ogg")
	})
	s.pauseMenu.SetOnSelect(func() {
		s.appCtx.AudioManager.PlaySound("assets/audio/Menu_Select2.ogg")
	})
	s.pauseScreen.SetMenu(s.pauseMenu)
	s.pauseScreen.SetFont(s.appCtx.Font)
	s.refreshPauseMenuLabels()
	s.pauseScreen.SetOnStart(func(p *pause.PauseScreen) {
		p.SetMenu(s.pauseMenu)
		if s.appCtx.AudioManager != nil {
			s.appCtx.AudioManager.PauseCurrentMusic()
		}
	})
	s.pauseScreen.SetOnFinish(func(p *pause.PauseScreen) {
		if s.appCtx.AudioManager != nil {
			s.appCtx.AudioManager.ResumeCurrentMusic()
		}
	})
}

func (s *PlatformerPhaseScene) buildSequencePlayer() {
	phase, err := s.appCtx.PhaseManager.GetCurrentPhase()
	if err == nil && phase.SequencePath != "" {
		s.sequencePlayer = sequences.NewSequencePlayer(s.appCtx)
		s.allowPause = phase.GoalType != phases.SequenceGoalType
		s.sequencePlayer.PlaySequence(phase.SequencePath)
	} else {
		s.allowPause = true
	}
}

func (s *PlatformerPhaseScene) initGoal() {
	if s.customGoal != nil {
		s.goal = s.customGoal
		return
	}
	phase, _ := s.appCtx.PhaseManager.GetCurrentPhase()
	switch phase.GoalType {
	case phases.ReactEndpointType:
		goal := &phases.ReachEndpointGoal{}
		goal.OnCompletion_ = func() {
			s.freezeAllActors()
			if s.appCtx.AudioManager != nil {
				s.appCtx.AudioManager.FadeOutCurrentTrack(time.Second)
			}
			s.completionTrigger.Enable(timing.FromDuration(time.Second))
		}
		s.goal = goal
	case phases.SequenceGoalType:
		s.goal = &phases.SequenceGoal{
			Player: s.sequencePlayer,
			OnCompleteFunc: func() {
				s.completionTrigger.Enable(timing.FromDuration(time.Second))
			},
		}
	default:
		s.goal = &phases.NoGoal{}
	}
}

func (s *PlatformerPhaseScene) freezeAllActors() {
	if s.appCtx == nil || s.appCtx.ActorManager == nil {
		return
	}
	s.appCtx.ActorManager.ForEach(func(actor actors.ActorEntity) {
		actor.SetImmobile(true)
		actor.SetFreeze(true)
	})
}

func (s *PlatformerPhaseScene) refreshPauseMenuLabels() {
	if s.pauseMenu == nil {
		return
	}
	i18n := s.appCtx.I18n
	s.pauseMenu.UpdateItemLabel(0, i18n.T("menu.start"))
	s.pauseMenu.UpdateItemLabel(1, i18n.T("menu.exit"))
}

func (s *PlatformerPhaseScene) subscribeEvents() {
	em := s.appCtx.EventManager
	em.Subscribe(actorevents.ActorJumpedType, func(e event.Event) {
		if s.appCtx.VFX == nil {
			return
		}
		if evt, ok := e.(*actorevents.ActorJumpedEvent); ok {
			s.appCtx.VFX.SpawnJumpPuff(evt.X, evt.Y+1.0, 1)
			s.appCtx.AudioManager.PlaySoundAtVolume("assets/audio/Menu_Select.ogg", 0.3)
		}
	})
	em.Subscribe(actorevents.ActorLandedType, func(e event.Event) {
		if s.appCtx.VFX == nil {
			return
		}
		if evt, ok := e.(*actorevents.ActorLandedEvent); ok {
			s.appCtx.VFX.SpawnLandingPuff(evt.X, evt.Y+1.0, 1)
		}
	})
}

func (s *PlatformerPhaseScene) canPause() bool {
	return s.allowPause && (s.sequencePlayer == nil || !s.sequencePlayer.IsPlaying())
}

func (s *PlatformerPhaseScene) endpointTrigger(id string) {
	if !s.hasPlayer {
		return
	}
	if s.deathActive {
		return
	}
	if s.endpointHandler != nil && s.endpointHandler(id) {
		return
	}
	switch id {
	case "SPIKE":
		s.startDeathSequence()
	case "CUTSCENE":
		// reserved
	default:
		if g, ok := s.goal.(*phases.ReachEndpointGoal); ok {
			g.Reach()
		}
	}
}

// OnFinish is called when the scene is about to be replaced.
func (s *PlatformerPhaseScene) OnFinish() {
	if s.appCtx == nil {
		return
	}
	if s.tilemapScene != nil {
		s.tilemapScene.OnFinish()
	}
	if s.appCtx.ProjectileManager != nil {
		s.appCtx.ProjectileManager.Clear()
	}
	if s.hasPlayer {
		if p, ok := s.appCtx.ActorManager.GetPlayer(); ok {
			p.UnblockMovement()
		}
		if ae, ok := any(s.player).(actors.ActorEntity); ok {
			s.appCtx.ActorManager.Unregister(ae)
		}
	}
}

// wireFlipperCallbacks sets up the screen-flipper immobility callbacks for p.
func (s *PlatformerPhaseScene) wireFlipperCallbacks(p Player) {
	s.onFlipStart = func() { p.SetImmobile(true) }
	s.onFlipFinish = func() { p.SetImmobile(false) }
}

// Update advances the scene by one frame.
func (s *PlatformerPhaseScene) Update() error {
	// full mode
	if s.tilemapScene != nil && s.tilemapScene.Tilemap() != nil {
		return s.fullUpdate()
	}
	// minimal mode (existing Update body)
	if s.space == nil {
		return nil
	}
	for _, i := range s.space.Bodies() {
		switch b := i.(type) {
		case platformer.PlatformerActorEntity:
			if b.State() == actors.Dead {
				s.space.RemoveBody(i)
				continue
			}
			if err := b.Update(s.space); err != nil {
				return err
			}
		case body.Obstacle:
			continue
		}
	}
	s.space.ProcessRemovals()
	return nil
}

func (s *PlatformerPhaseScene) fullUpdate() error {
	if s.pauseScreen != nil && s.canPause() {
		s.pauseScreen.Update()
		if s.pauseScreen.IsPaused() {
			return nil
		}
	}
	if s.sequencePlayer != nil {
		s.sequencePlayer.Update()
		if s.sequencePlayer.IsDebugPaused() {
			return nil
		}
	}
	if s.appCtx.VFX != nil {
		s.appCtx.VFX.Update()
	}
	if s.screenFlipper != nil {
		s.screenFlipper.Update()
		if s.screenFlipper.IsFlipping() {
			return nil
		}
	}
	if s.hasPlayer && s.player != nil && !s.deathActive &&
		(s.player.State() == s.dyingState || s.player.State() == s.deadState) {
		s.startDeathSequence()
	}
	s.completionTrigger.Update()
	s.deathTrigger.Update()
	if s.deathTrigger.Trigger() {
		s.appCtx.SceneManager.NavigateTo(
			s.rebootScene,
			transition.NewFader(0, config.Get().FadeVisibleDuration),
			false,
		)
	}
	if s.goal != nil && s.goal.IsCompleted() && !s.completionTrigger.IsEnabled() {
		s.goal.OnCompletion()
	}
	if config.Get().CamDebug {
		s.camera.CamDebug()
	}
	if s.completionTrigger.Trigger() {
		s.appCtx.CompleteCurrentPhase(nil, true)
	}
	s.camera.Update()
	if err := s.tilemapScene.BaseScene.Update(); err != nil {
		return err
	}
	s.count++
	space := s.space
	for _, i := range space.Bodies() {
		switch b := i.(type) {
		case platformer.PlatformerActorEntity:
			if b.State() == actors.Dead {
				if s.appCtx.VFX != nil {
					x, y := b.GetPositionMin()
					w, h := b.GetShape().Width(), b.GetShape().Height()
					s.appCtx.VFX.SpawnDeathExplosion(float64(x)+float64(w)/2, float64(y)+float64(h)/2, 30)
				}
				space.RemoveBody(i)
				continue
			}
			if err := b.Update(space); err != nil {
				return err
			}
		case items.Item:
			if b.IsRemoved() {
				space.RemoveBody(i)
				continue
			}
			if err := b.Update(space); err != nil {
				return err
			}
		case body.Obstacle:
			continue
		}
	}
	if s.appCtx.ProjectileManager != nil {
		s.appCtx.ProjectileManager.Update()
	}
	if s.hasPlayer && s.player != nil {
		space.ResolveCollisions(s.player)
	}
	space.ProcessRemovals()
	return nil
}

// SetPlayer wires a player into the scene. Calling this with a non-nil value
// sets hasPlayer=true; calling with nil clears the player.
func (s *PlatformerPhaseScene) SetPlayer(p Player) {
	s.player = p
	s.hasPlayer = p != nil
	if s.hasFlipper && p != nil {
		s.wireFlipperCallbacks(p)
	}
}

// DeathActive reports whether the death sequence has been triggered.
func (s *PlatformerPhaseScene) DeathActive() bool { return s.deathActive }

// StartDeathSequence triggers the death sequence programmatically (e.g., from
// a state check in the game layer).
func (s *PlatformerPhaseScene) StartDeathSequence() { s.startDeathSequence() }

// DrawActors renders the actor bodies (without filling the background). Game
// layers call this after drawing the tilemap so actors sort correctly with
// game-specific bodies (items, etc.) that share the same sorted pass.
func (s *PlatformerPhaseScene) DrawActors(screen *ebiten.Image) {
	if s.space != nil {
		for _, b := range draworder.SortByGroundY(s.space.Bodies()) {
			switch sb := b.(type) {
			case platformer.PlatformerActorEntity:
				opts := sb.ImageOptions()
				sb.UpdateImageOptions()
				s.camera.Draw(sb.Image(), opts, screen)
				if config.Get().CollisionBox {
					s.camera.DrawCollisionBox(screen, sb)
				}
			case body.Obstacle:
				if config.Get().CollisionBox {
					s.camera.DrawCollisionBox(screen, sb)
				}
			}
		}
	}
	if s.flashCount > 0 {
		screenutil.DrawScreenFlash(screen)
		s.flashCount--
	}
	if s.debugDrawHook != nil {
		s.debugDrawHook(screen)
	}
}

// Draw renders the scene (background fill + actors).
func (s *PlatformerPhaseScene) Draw(screen *ebiten.Image) {
	if s.tilemapScene != nil && s.tilemapScene.Tilemap() != nil {
		s.fullDraw(screen)
		return
	}
	// minimal draw (existing body)
	screen.Fill(color.RGBA{0, 0, 0, 0xff})
	s.DrawActors(screen)
}

func (s *PlatformerPhaseScene) DrawOver(screen *ebiten.Image) {
	s.sequencePlayer.DrawOver(screen)
}

func (s *PlatformerPhaseScene) fullDraw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 0xff})
	tilemapImg, _ := s.tilemapScene.Tilemap().Image(screen)
	s.camera.Draw(tilemapImg, s.tilemapScene.Tilemap().ImageOptions(), screen)
	space := s.space
	for _, b := range draworder.SortByGroundY(space.Bodies()) {
		switch sb := b.(type) {
		case platformer.PlatformerActorEntity:
			opts := sb.ImageOptions()
			sb.UpdateImageOptions()
			s.camera.Draw(sb.Image(), opts, screen)
			if config.Get().CollisionBox {
				s.camera.DrawCollisionBox(screen, sb)
			}
		case items.Item:
			if sb.IsRemoved() {
				continue
			}
			opts := sb.ImageOptions()
			sb.UpdateImageOptions()
			s.camera.Draw(sb.Image(), opts, screen)
			if config.Get().CollisionBox {
				s.camera.DrawCollisionBox(screen, sb)
			}
		case body.Obstacle:
			if config.Get().CollisionBox {
				s.camera.DrawCollisionBox(screen, sb)
			}
		}
	}
	if s.appCtx.ProjectileManager != nil {
		camX, camY := s.camera.GetActualCenter()
		camX -= float64(config.Get().ScreenWidth) / 2
		camY -= float64(config.Get().ScreenHeight) / 2
		s.appCtx.ProjectileManager.DrawWithOffset(screen, camX, camY)
	}
	if config.Get().CollisionBox {
		if pm := s.appCtx.ProjectileManager; pm != nil {
			pm.DrawCollisionBoxesWithOffset(func(b body.Collidable) {
				s.camera.DrawCollisionBox(screen, b)
			})
		}
	}
	if s.flashCount > 0 {
		screenutil.DrawScreenFlash(screen)
		s.flashCount--
	}
	if s.appCtx.VFX != nil {
		s.appCtx.VFX.Draw(screen, s.camera)
	}
	if s.vfx != nil && s.hasPlayer && s.player != nil {
		s.vfx.Draw(screen, s.camera, s.player)
	}
	if s.debugDrawHook != nil {
		s.debugDrawHook(screen)
	}
	if s.pauseScreen != nil && s.pauseScreen.IsPaused() {
		s.drawPause(screen)
	}
	if s.sequencePlayer != nil {
		s.sequencePlayer.Draw(screen)
	}
}

func (s *PlatformerPhaseScene) drawPause(screen *ebiten.Image) {
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
	if m := s.pauseScreen.Menu(); m != nil {
		m.Draw(screen, s.pauseScreen.Font(), cfg.ScreenWidth/2, cfg.ScreenHeight/2)
	}
}

// SetDebugDrawHook sets a function invoked at the end of every Draw call.
func (s *PlatformerPhaseScene) SetDebugDrawHook(f func(*ebiten.Image)) {
	s.debugDrawHook = f
}

// SetGoal installs a game-layer goal. When set, the kit's GoalType switch is
// skipped and the provided goal is used as-is.
func (s *PlatformerPhaseScene) SetGoal(g phases.Goal) { s.customGoal = g }

// SetEndpointHandler installs a game-layer handler invoked on endpoint touch
// before the kit's built-in SPIKE/CUTSCENE/default logic. Return true to
// indicate the touch was handled and suppress the default behavior.
func (s *PlatformerPhaseScene) SetEndpointHandler(f func(id string) bool) {
	s.endpointHandler = f
}

// SetOnStarted installs a callback invoked at the end of OnStart, after the
// player, tilemap actors, goal and events have been wired. Useful for
// post-init scans of the physics space.
func (s *PlatformerPhaseScene) SetOnStarted(f func()) { s.onStarted = f }

// Space exposes the physics space for game-layer body counters / queries.
func (s *PlatformerPhaseScene) Space() body.BodiesSpace { return s.space }

// Player exposes the active player (may be nil if the phase has no player).
func (s *PlatformerPhaseScene) Player() Player { return s.player }

// PlaySequence plays a sequence file on the scene's sequence player,
// constructing one on-demand if the phase did not have a default sequence.
// Game-layer code uses this for ad-hoc cutscenes triggered by gameplay events.
func (s *PlatformerPhaseScene) PlaySequence(path string) {
	if s.sequencePlayer == nil {
		s.sequencePlayer = sequences.NewSequencePlayer(s.appCtx)
	}
	s.sequencePlayer.PlaySequence(path)
}

// CompletionTrigger returns the kit's completion trigger so game-layer goals
// can schedule phase completion through the same path as built-in goals.
func (s *PlatformerPhaseScene) EnableCompletionTrigger(frames int) {
	s.completionTrigger.Enable(frames)
}

// FreezeAllActors freezes every registered actor. Exposed so game-layer goals
// can reuse the same freeze logic as the built-in completion path.
func (s *PlatformerPhaseScene) FreezeAllActors() { s.freezeAllActors() }

// TriggerScreenFlash triggers a white flash overlay for the next two frames.
func (s *PlatformerPhaseScene) TriggerScreenFlash() {
	s.flashCount = 2
}

// EnableVignetteDarkness enables the world darkness overlay with the given radius in screen pixels.
func (s *PlatformerPhaseScene) EnableVignetteDarkness(radiusPx float64) {
	if s.vfx == nil {
		s.vfx = enginevfx.NewVignette()
	}
	s.vfx.Enable(radiusPx)
}

// DisableVignetteDarkness disables the world darkness overlay.
func (s *PlatformerPhaseScene) DisableVignetteDarkness() {
	if s.vfx == nil {
		return
	}
	s.vfx.Disable()
}

// startDeathSequence activates the death state on the player.
func (s *PlatformerPhaseScene) startDeathSequence() {
	if s.deathActive {
		return
	}
	s.deathActive = true
	if s.player == nil {
		return
	}
	if s.setNewStateFatalHook != nil {
		s.setNewStateFatalHook(s.dyingState)
	} else if ch := s.player.GetCharacter(); ch != nil {
		ch.SetNewStateFatal(s.dyingState)
	}
	s.player.SetImmobile(true)
	if s.OnDeathStarted != nil {
		s.OnDeathStarted()
	}
}

// --- test-support API -------------------------------------------------------

// TestOptions configures a scene created via NewForTest.
type TestOptions struct {
	CameraCenterX          float64
	CameraCenterY          float64
	ScreenWidth            float64
	ScreenHeight           float64
	PlayerY                float64
	DyingState             actors.ActorStateEnum
	DeadState              actors.ActorStateEnum
	HasFlipper             bool
	HasPlayerStartPosition bool
}

// NewForTest creates a PlatformerPhaseScene in a minimal headless state
// suitable for unit tests. It must not be used in production code.
func NewForTest(opts TestOptions) *PlatformerPhaseScene {
	cfg := config.Get()
	sw := opts.ScreenWidth
	if sw == 0 {
		sw = float64(cfg.ScreenWidth)
	}
	sh := opts.ScreenHeight
	if sh == 0 {
		sh = float64(cfg.ScreenHeight)
	}

	cam := enginecamera.NewController(opts.CameraCenterX, opts.CameraCenterY)
	cam.DisableSmoothing()
	cam.SetCenter(opts.CameraCenterX, opts.CameraCenterY)

	space := &testSpace{}

	s := newScene(cam, space, sw, sh, opts.DyingState, opts.DeadState)
	s.hasPlayer = opts.HasPlayerStartPosition
	s.hasFlipper = opts.HasFlipper

	return s
}

// SetPlayerForTest injects a mock player for testing. Delegates to SetPlayer.
func (s *PlatformerPhaseScene) SetPlayerForTest(p Player) {
	s.SetPlayer(p)
}

// SetSetNewStateFatalRecorder overrides the SetNewStateFatal call path so tests
// can assert it was invoked with the correct state.
func (s *PlatformerPhaseScene) SetSetNewStateFatalRecorder(f func(actors.ActorStateEnum)) {
	s.setNewStateFatalHook = f
}

// DeathActiveForTest returns whether the death sequence is active.
func (s *PlatformerPhaseScene) DeathActiveForTest() bool {
	return s.deathActive
}

// SetDeathActiveForTest sets the death-active flag for test setup.
func (s *PlatformerPhaseScene) SetDeathActiveForTest(v bool) {
	s.deathActive = v
}

// InvokeScreenFlipperOnFlipStartForTest simulates the OnFlipStart callback.
func (s *PlatformerPhaseScene) InvokeScreenFlipperOnFlipStartForTest() {
	if s.onFlipStart != nil {
		s.onFlipStart()
	} else {
		log.Println("InvokeScreenFlipperOnFlipStartForTest: onFlipStart not wired")
	}
}

// InvokeScreenFlipperOnFlipFinishForTest simulates the OnFlipFinish callback.
func (s *PlatformerPhaseScene) InvokeScreenFlipperOnFlipFinishForTest() {
	if s.onFlipFinish != nil {
		s.onFlipFinish()
	} else {
		log.Println("InvokeScreenFlipperOnFlipFinishForTest: onFlipFinish not wired")
	}
}

// ScreenFlipperForTest returns a non-nil sentinel when a screen flipper is
// configured, and nil when there is none. Tests may nil-check the return value.
func (s *PlatformerPhaseScene) ScreenFlipperForTest() interface{} {
	if !s.hasFlipper {
		return nil
	}
	return struct{}{}
}

// CameraIsFixedModeForTest returns true when the camera is in non-follow mode.
func (s *PlatformerPhaseScene) CameraIsFixedModeForTest() bool {
	return !s.camera.IsFollowing()
}

// --- minimal in-test space implementation -----------------------------------

type testSpace struct {
	bodies []body.Collidable
}

func (s *testSpace) AddBody(b body.Collidable) { s.bodies = append(s.bodies, b) }
func (s *testSpace) Bodies() []body.Collidable { return s.bodies }
func (s *testSpace) RemoveBody(b body.Collidable) {
	for i, c := range s.bodies {
		if c == b {
			s.bodies = append(s.bodies[:i], s.bodies[i+1:]...)
			return
		}
	}
}
func (s *testSpace) QueueForRemoval(body.Collidable)                                     {}
func (s *testSpace) ProcessRemovals()                                                    {}
func (s *testSpace) Clear()                                                              {}
func (s *testSpace) ResolveCollisions(body.Collidable) (bool, bool)                      { return false, false }
func (s *testSpace) SetTilemapDimensionsProvider(tilemaplayer.TilemapDimensionsProvider) {}
func (s *testSpace) GetTilemapDimensionsProvider() tilemaplayer.TilemapDimensionsProvider {
	return nil
}
func (s *testSpace) Find(id string) body.Collidable {
	for _, b := range s.bodies {
		if b.ID() == id {
			return b
		}
	}
	return nil
}
func (s *testSpace) Query(image.Rectangle) []body.Collidable { return nil }
