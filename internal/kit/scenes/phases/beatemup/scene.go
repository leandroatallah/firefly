// Package beatemupphasescene provides a genre-reusable beat-em-up phase scene
// that assembles the engine scene base with beat-em-up-specific actor handling.
package beatemupphasescene

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
	"github.com/boilerplate/ebiten-template/internal/engine/scene"
	"github.com/boilerplate/ebiten-template/internal/engine/scene/pause"
	"github.com/boilerplate/ebiten-template/internal/engine/scene/phases"
	"github.com/boilerplate/ebiten-template/internal/engine/scene/transition"
	"github.com/boilerplate/ebiten-template/internal/engine/sequences"
	"github.com/boilerplate/ebiten-template/internal/engine/ui/menu"
	"github.com/boilerplate/ebiten-template/internal/engine/utils"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/timing"
	beatemupkit "github.com/boilerplate/ebiten-template/internal/kit/actors/beatemup"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Player is the minimal interface the scene requires from a beat-em-up player.
// It is a subset of beatemupkit.BeatEmUpActorEntity that mock and production
// types both satisfy.
type Player interface {
	body.Collidable
	body.Drawable

	GetPositionMin() (x, y int)
	GetShape() body.Shape
	Altitude16() int

	State() actors.ActorStateEnum
	SetImmobile(bool)
	GetCharacter() *actors.Character
	Update(body.BodiesSpace) error
}

// statefulBody is an internal interface for bodies that expose an actor state.
// The type switch uses this instead of the full BeatEmUpActorEntity so mock
// types in tests (which don't implement the full interface) are still matched.
type statefulBody interface {
	body.Collidable
	State() actors.ActorStateEnum
	Update(body.BodiesSpace) error
}

// BeatemupPhaseScene is the kit-level beat-em-up scene. It manages a
// beat-em-up player and uses altitude-aware draw ordering.
// Game layers embed or compose this type to wire concrete factories.
type BeatemupPhaseScene struct {
	camera       *enginecamera.Controller
	space        body.BodiesSpace
	hasPlayer    bool
	cameraMode   scene.CameraMode
	screenWidth  float64
	screenHeight float64

	// Player — holds the minimal interface; production code uses BeatEmUpActorEntity.
	player Player

	// OnDeathStarted is called after the death sequence activates. Game layers
	// set this to spawn VFX and enable navigation triggers.
	OnDeathStarted func()

	// Death state tracking.
	deathActive bool

	// Hook for overriding SetNewStateFatal calls (used in tests for assertion).
	setNewStateFatalHook func(actors.ActorStateEnum)

	// actorDrawHandler overrides the draw loop per actor (used in tests to record draw order).
	actorDrawHandler func(screen *ebiten.Image, b body.Collidable) bool

	// Debug hook invoked at the end of Draw.
	debugDrawHook func(*ebiten.Image)

	// flashCount tracks screen-flash frames.
	flashCount int

	// Full-loop fields (nil when created via NewForTest)
	tilemapScene      *scene.TilemapScene
	appCtx            *app.AppContext
	goal              phases.Goal
	sequencePlayer    sequencestypes.Player
	allowPause        bool
	pauseScreen       *pause.PauseScreen
	pauseMenu         *menu.Menu
	completionTrigger utils.DelayTrigger
	deathTrigger      utils.DelayTrigger
	rebootScene       navigation.SceneType
	menuScene         navigation.SceneType
	playerFactory     func(*app.AppContext) (Player, error)
	initActors        func(*scene.TilemapScene)
	count             int
}

// SetAppContext implements navigation.Scene. The kit scene context is set at
// construction time via NewWithOptions; this method satisfies the interface.
func (s *BeatemupPhaseScene) SetAppContext(_ any) {}

// New creates a BeatemupPhaseScene for production use.
func New(
	cam *enginecamera.Controller,
	space body.BodiesSpace,
	sw, sh float64,
) *BeatemupPhaseScene {
	return newScene(cam, space, sw, sh)
}

// newScene creates a new BeatemupPhaseScene with the given camera and space.
func newScene(
	cam *enginecamera.Controller,
	space body.BodiesSpace,
	sw, sh float64,
) *BeatemupPhaseScene {
	return &BeatemupPhaseScene{
		camera:       cam,
		space:        space,
		screenWidth:  sw,
		screenHeight: sh,
	}
}

// OnStart initialises the scene for the current phase.
func (s *BeatemupPhaseScene) OnStart() {
	if s.appCtx != nil {
		s.fullOnStart()
		return
	}
	// minimal path (NewForTest)
	if s.hasPlayer {
		s.cameraMode = scene.CameraModeFollow
		s.camera.SetFollowing(true)
		// Beat-em-up cameras follow freely in both directions.
		s.camera.SetVerticalOnlyUpward(false)
		if s.player != nil {
			s.camera.SetFollowTarget(s.player)
		}
	} else {
		s.cameraMode = scene.CameraModeFixed
		s.camera.SetFollowing(false)
	}
}

func (s *BeatemupPhaseScene) fullOnStart() {
	ts := scene.NewTilemapScene(s.appCtx)
	ts.OnStart()
	s.tilemapScene = ts
	s.camera = ts.Camera()
	s.space = s.appCtx.Space
	s.count = 0

	s.hasPlayer = ts.Tilemap().HasPlayerStartPosition()

	if s.hasPlayer {
		if s.player == nil {
			p, err := s.playerFactory(s.appCtx)
			if err != nil {
				log.Fatal(err)
			}
			s.player = p
		}
		p := s.player
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
		if ae, ok := any(s.player).(actors.ActorEntity); ok {
			ts.SetPlayerStartPosition(ae)
		}
	}

	// Bound the camera to the full tilemap rectangle.
	tilemapRect := image.Rect(0, 0, ts.GetTilemapWidth(), ts.GetTilemapHeight())
	s.camera.SetBounds(&tilemapRect)

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
}

func (s *BeatemupPhaseScene) buildPauseScreen() {
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

func (s *BeatemupPhaseScene) buildSequencePlayer() {
	phase, err := s.appCtx.PhaseManager.GetCurrentPhase()
	if err == nil && phase.SequencePath != "" {
		s.sequencePlayer = sequences.NewSequencePlayer(s.appCtx)
		s.allowPause = phase.GoalType != phases.SequenceGoalType
		s.sequencePlayer.PlaySequence(phase.SequencePath)
	} else {
		s.allowPause = true
	}
}

func (s *BeatemupPhaseScene) initGoal() {
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

func (s *BeatemupPhaseScene) freezeAllActors() {
	if s.appCtx == nil || s.appCtx.ActorManager == nil {
		return
	}
	s.appCtx.ActorManager.ForEach(func(actor actors.ActorEntity) {
		actor.SetImmobile(true)
		actor.SetFreeze(true)
	})
}

func (s *BeatemupPhaseScene) refreshPauseMenuLabels() {
	if s.pauseMenu == nil {
		return
	}
	i18n := s.appCtx.I18n
	s.pauseMenu.UpdateItemLabel(0, i18n.T("menu.start"))
	s.pauseMenu.UpdateItemLabel(1, i18n.T("menu.exit"))
}

func (s *BeatemupPhaseScene) subscribeEvents() {
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

func (s *BeatemupPhaseScene) canPause() bool {
	return s.allowPause && (s.sequencePlayer == nil || !s.sequencePlayer.IsPlaying())
}

func (s *BeatemupPhaseScene) endpointTrigger(id string) {
	if !s.hasPlayer {
		return
	}
	if s.deathActive {
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
func (s *BeatemupPhaseScene) OnFinish() {
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

// Update advances the scene by one frame.
// Beat-em-up scenes do NOT perform fall-death checks.
func (s *BeatemupPhaseScene) Update() error {
	// full mode
	if s.tilemapScene != nil && s.tilemapScene.Tilemap() != nil {
		return s.fullUpdate()
	}
	// minimal mode
	if s.space == nil {
		return nil
	}
	for _, i := range s.space.Bodies() {
		switch b := i.(type) {
		case statefulBody:
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

func (s *BeatemupPhaseScene) fullUpdate() error {
	if s.pauseScreen != nil && s.canPause() {
		s.pauseScreen.Update()
		if s.pauseScreen.IsPaused() {
			return nil
		}
	}
	if s.sequencePlayer != nil {
		s.sequencePlayer.Update()
	}
	if s.appCtx.VFX != nil {
		s.appCtx.VFX.Update()
	}
	// Beat-em-up: no fall-death check, no screenFlipper
	if s.hasPlayer && s.player != nil && !s.deathActive &&
		(s.player.State() == actors.Dying || s.player.State() == actors.Dead) {
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
		s.appCtx.CompleteCurrentPhase(transition.NewFader(0, config.Get().FadeVisibleDuration), true)
	}
	s.camera.Update()
	if err := s.tilemapScene.BaseScene.Update(); err != nil {
		return err
	}
	s.count++
	space := s.space
	for _, i := range space.Bodies() {
		switch b := i.(type) {
		case beatemupkit.BeatEmUpActorEntity:
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

// SetPlayer wires a player into the scene. Calling with non-nil sets
// hasPlayer=true; calling with nil clears the player.
func (s *BeatemupPhaseScene) SetPlayer(p Player) {
	s.player = p
	s.hasPlayer = p != nil
}

// DeathActive reports whether the death sequence has been triggered.
func (s *BeatemupPhaseScene) DeathActive() bool { return s.deathActive }

// StartDeathSequence triggers the death sequence programmatically (e.g., from
// a player-state check in the game layer).
func (s *BeatemupPhaseScene) StartDeathSequence() { s.startDeathSequence() }

// startDeathSequence activates the death state on the player.
func (s *BeatemupPhaseScene) startDeathSequence() {
	if s.deathActive {
		return
	}
	s.deathActive = true
	if s.player == nil {
		return
	}
	if s.setNewStateFatalHook != nil {
		s.setNewStateFatalHook(actors.Dying)
	} else if ch := s.player.GetCharacter(); ch != nil {
		ch.SetNewStateFatal(actors.Dying)
	}
	s.player.SetImmobile(true)
	if s.OnDeathStarted != nil {
		s.OnDeathStarted()
	}
}

// DrawActors renders the actor bodies using altitude-aware ordering without
// filling the background. Game layers call this after drawing the tilemap.
func (s *BeatemupPhaseScene) DrawActors(screen *ebiten.Image) {
	if s.space == nil {
		return
	}
	for _, b := range draworder.SortByGroundYAltitude(s.space.Bodies()) {
		if s.actorDrawHandler != nil {
			s.actorDrawHandler(screen, b)
			continue
		}
		switch sb := b.(type) {
		case beatemupkit.BeatEmUpActorEntity:
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

// Draw renders the scene (background fill + actors with altitude-aware ordering).
func (s *BeatemupPhaseScene) Draw(screen *ebiten.Image) {
	if s.tilemapScene != nil && s.tilemapScene.Tilemap() != nil {
		s.fullDraw(screen)
		return
	}
	screen.Fill(color.RGBA{0, 0, 0, 0xff})
	s.DrawActors(screen)
}

func (s *BeatemupPhaseScene) fullDraw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 0xff})
	tilemapImg, _ := s.tilemapScene.Tilemap().Image(screen)
	s.camera.Draw(tilemapImg, s.tilemapScene.Tilemap().ImageOptions(), screen)
	space := s.space
	for _, b := range draworder.SortByGroundYAltitude(space.Bodies()) {
		switch sb := b.(type) {
		case beatemupkit.BeatEmUpActorEntity:
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
	if s.debugDrawHook != nil {
		s.debugDrawHook(screen)
	}
	if s.pauseScreen != nil && s.pauseScreen.IsPaused() {
		s.drawPause(screen)
	}
}

func (s *BeatemupPhaseScene) drawPause(screen *ebiten.Image) {
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

// TriggerScreenFlash triggers a white screen flash effect for feedback.
func (s *BeatemupPhaseScene) TriggerScreenFlash() {
	s.flashCount = 2
}

// --- test-support API -------------------------------------------------------

// TestOptions configures a scene created via NewForTest.
type TestOptions struct {
	CameraCenterX          float64
	CameraCenterY          float64
	ScreenWidth            float64
	ScreenHeight           float64
	HasPlayerStartPosition bool
}

// NewForTest creates a BeatemupPhaseScene in a minimal headless state
// suitable for unit tests. It must not be used in production code.
func NewForTest(opts TestOptions) *BeatemupPhaseScene {
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

	s := newScene(cam, space, sw, sh)
	s.hasPlayer = opts.HasPlayerStartPosition

	return s
}

// AddBodyForTest adds a body directly to the space for testing.
func (s *BeatemupPhaseScene) AddBodyForTest(b body.Collidable) {
	s.space.AddBody(b)
}

// SetActorDrawHandlerForTest overrides the per-actor draw logic so tests can
// record the draw order without a full graphics context.
func (s *BeatemupPhaseScene) SetActorDrawHandlerForTest(f func(*ebiten.Image, body.Collidable) bool) {
	s.actorDrawHandler = f
}

// SpaceContainsBodyForTest reports whether the given body is still in the space.
func (s *BeatemupPhaseScene) SpaceContainsBodyForTest(b body.Collidable) bool {
	for _, c := range s.space.Bodies() {
		if c == b {
			return true
		}
	}
	return false
}

// SetPlayerForTest injects a mock player for testing. Delegates to SetPlayer.
func (s *BeatemupPhaseScene) SetPlayerForTest(p Player) {
	s.SetPlayer(p)
}

// SetSetNewStateFatalRecorder overrides the SetNewStateFatal call path so tests
// can assert it was invoked with the correct state.
func (s *BeatemupPhaseScene) SetSetNewStateFatalRecorder(f func(actors.ActorStateEnum)) {
	s.setNewStateFatalHook = f
}

// DeathActiveForTest returns whether the death sequence is active.
func (s *BeatemupPhaseScene) DeathActiveForTest() bool { return s.DeathActive() }

// EngineCameraForTest returns the underlying engine camera controller.
func (s *BeatemupPhaseScene) EngineCameraForTest() *enginecamera.Controller {
	return s.camera
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
