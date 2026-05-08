// Package beatemupphasescene provides a genre-reusable beat-em-up phase scene
// that assembles the engine scene base with beat-em-up-specific actor handling.
package beatemupphasescene

import (
	"image"
	"image/color"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/tilemaplayer"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	enginecamera "github.com/boilerplate/ebiten-template/internal/engine/render/camera"
	"github.com/boilerplate/ebiten-template/internal/engine/render/draworder"
	"github.com/boilerplate/ebiten-template/internal/engine/scene"
	beatemupkit "github.com/boilerplate/ebiten-template/internal/kit/actors/beatemup"
	"github.com/hajimehoshi/ebiten/v2"
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
}

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

// Update advances the scene by one frame.
// Beat-em-up scenes do NOT perform fall-death checks.
func (s *BeatemupPhaseScene) Update() error {
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
	screen.Fill(color.RGBA{0, 0, 0, 0xff})
	s.DrawActors(screen)
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
