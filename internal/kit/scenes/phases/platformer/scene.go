// Package platformerphasescene provides a genre-reusable platformer phase scene
// that assembles the engine scene base with platformer-specific actor handling.
package platformerphasescene

import (
	"image"
	"image/color"
	"log"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/tilemaplayer"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	enginecamera "github.com/boilerplate/ebiten-template/internal/engine/render/camera"
	"github.com/boilerplate/ebiten-template/internal/engine/render/draworder"
	"github.com/boilerplate/ebiten-template/internal/engine/render/screenutil"
	enginevfx "github.com/boilerplate/ebiten-template/internal/engine/render/vfx"
	"github.com/boilerplate/ebiten-template/internal/engine/scene"
	"github.com/boilerplate/ebiten-template/internal/kit/actors/platformer"
	"github.com/hajimehoshi/ebiten/v2"
)

// platformerPlayer is the minimal interface the scene requires from a platformer player.
// It is a subset of platformer.PlatformerActorEntity that the mock and production
// types both satisfy.
type platformerPlayer interface {
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
	player     platformerPlayer
	dyingState actors.ActorStateEnum
	deadState  actors.ActorStateEnum

	// Death sequence
	deathActive bool

	// Screen flipper callbacks. Set up by OnStart when hasFlipper=true.
	onFlipStart  func()
	onFlipFinish func()
	hasFlipper   bool

	// Hook for overriding SetNewStateFatal calls (used in tests for assertion).
	setNewStateFatalHook func(actors.ActorStateEnum)

	// Debug hook invoked at the end of Draw.
	debugDrawHook func(*ebiten.Image)

	// VFX vignette (may be nil)
	vfx        *enginevfx.Vignette
	flashCount int
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

// OnStart initialises the scene for the current phase.
func (s *PlatformerPhaseScene) OnStart() {
	if s.hasPlayer {
		s.cameraMode = scene.CameraModeFollow
		s.camera.SetFollowing(true)
		s.camera.SetVerticalOnlyUpward(true)
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

// wireFlipperCallbacks sets up the screen-flipper immobility callbacks for p.
func (s *PlatformerPhaseScene) wireFlipperCallbacks(p platformerPlayer) {
	s.onFlipStart = func() { p.SetImmobile(true) }
	s.onFlipFinish = func() { p.SetImmobile(false) }
}

// Update advances the scene by one frame.
func (s *PlatformerPhaseScene) Update() error {
	if s.hasPlayer && s.player != nil {
		s.checkPlayerFallDeath()
	}

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

// Draw renders the scene.
func (s *PlatformerPhaseScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 0xff})

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

// SetDebugDrawHook sets a function invoked at the end of every Draw call.
func (s *PlatformerPhaseScene) SetDebugDrawHook(f func(*ebiten.Image)) {
	s.debugDrawHook = f
}

// TriggerScreenFlash triggers a white flash overlay for the next two frames.
func (s *PlatformerPhaseScene) TriggerScreenFlash() {
	s.flashCount = 2
}

// checkPlayerFallDeath fires startDeathSequence when the player's top Y
// exceeds the camera's bottom edge.
func (s *PlatformerPhaseScene) checkPlayerFallDeath() {
	if s.camera == nil || s.player == nil {
		return
	}
	if s.deathActive {
		return
	}

	_, camY := s.camera.GetActualCenter()
	_, playerY := s.player.GetPositionMin()

	cameraBottom := camY + s.screenHeight/2
	playerTop := float64(playerY)
	if playerTop > cameraBottom {
		s.startDeathSequence()
	}
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

// SetPlayerForTest injects a mock player for testing. If hasFlipper is set and
// the player is non-nil, it also wires the flip callbacks so tests can invoke them.
func (s *PlatformerPhaseScene) SetPlayerForTest(p platformerPlayer) {
	s.player = p
	s.hasPlayer = p != nil
	if s.hasFlipper && p != nil {
		s.wireFlipperCallbacks(p)
	}
}

// SetSetNewStateFatalRecorder overrides the SetNewStateFatal call path so tests
// can assert it was invoked with the correct state.
func (s *PlatformerPhaseScene) SetSetNewStateFatalRecorder(f func(actors.ActorStateEnum)) {
	s.setNewStateFatalHook = f
}

// CheckPlayerFallDeathForTest exposes checkPlayerFallDeath for white-box testing.
func (s *PlatformerPhaseScene) CheckPlayerFallDeathForTest() {
	s.checkPlayerFallDeath()
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
