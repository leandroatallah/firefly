// Package gamebeatemupphase implements the game-layer beat-em-up phase scene.
// It wires the tilemap and delegates actor update/draw to the kit beatemup scene.
package gamebeatemupphase

import (
	"image/color"
	"log"
	"time"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/navigation"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/scene"
	"github.com/boilerplate/ebiten-template/internal/engine/scene/transition"
	"github.com/boilerplate/ebiten-template/internal/engine/utils"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/timing"
	gamestates "github.com/boilerplate/ebiten-template/internal/game/entity/actors/states"
	scenestypes "github.com/boilerplate/ebiten-template/internal/game/scenes/types"
	beatemupkit "github.com/boilerplate/ebiten-template/internal/kit/actors/beatemup"
	beatemupphasescene "github.com/boilerplate/ebiten-template/internal/kit/scenes/phases/beatemup"
	"github.com/hajimehoshi/ebiten/v2"
)

// BeatemupPhaseScene is the game-layer beat-em-up phase scene.
// It embeds TilemapScene for tilemap/camera/space management and delegates
// actor update and altitude-aware draw ordering to the kit beatemup scene.
type BeatemupPhaseScene struct {
	*scene.TilemapScene
	kitScene *beatemupphasescene.BeatemupPhaseScene

	player       beatemupkit.BeatEmUpActorEntity
	hasPlayer    bool
	deathTrigger utils.DelayTrigger
}

// NewBeatemupPhaseScene constructs a BeatemupPhaseScene wired to the given AppContext.
func NewBeatemupPhaseScene(ctx *app.AppContext) navigation.Scene {
	ts := scene.NewTilemapScene(ctx)
	s := &BeatemupPhaseScene{TilemapScene: ts}
	s.SetAppContext(ctx)

	cfg := config.Get()
	s.kitScene = beatemupphasescene.New(
		ts.Camera(),
		ctx.Space,
		float64(cfg.ScreenWidth),
		float64(cfg.ScreenHeight),
	)
	return s
}

func (s *BeatemupPhaseScene) OnStart() {
	s.TilemapScene.OnStart()

	s.Tilemap().CreateCollisionBodies(s.PhysicsSpace(), nil)

	ctx := s.AppContext()
	s.hasPlayer = s.Tilemap().HasPlayerStartPosition()

	if s.hasPlayer {
		p, err := createPlayer(ctx)
		if err != nil {
			log.Fatal(err)
		}
		s.player = p
		ctx.ActorManager.Register(s.player)
		ctx.ActorManager.RegisterPrimary(s.player)
		s.PhysicsSpace().AddBody(s.player)

		s.kitScene.SetPlayer(p)
		s.kitScene.OnDeathStarted = func() {
			if ctx.VFX != nil {
				deathX, deathY := s.player.GetPositionMin()
				deathW, deathH := s.player.GetShape().Width(), s.player.GetShape().Height()
				ctx.VFX.SpawnDeathExplosion(
					float64(deathX)+float64(deathW)/2,
					float64(deathY)+float64(deathH)/2,
					50,
				)
			}
			s.deathTrigger.Enable(timing.FromDuration(time.Second))
		}

		s.SetCameraConfig(scene.CameraConfig{Mode: scene.CameraModeFollow})
		s.Camera().SetVerticalOnlyUpward(false)
		s.Camera().SetFollowTarget(s.player)
		s.SetPlayerStartPosition(s.player)
	} else {
		s.SetCameraConfig(scene.CameraConfig{Mode: scene.CameraModeFixed})
		if x, y, found := s.Tilemap().GetCameraStartPosition(); found {
			s.Camera().SetPositionTopLeft(float64(x), float64(y))
		} else {
			s.Camera().SetPositionTopLeft(0, 0)
		}
	}
}

func (s *BeatemupPhaseScene) Update() error {
	if err := s.kitScene.Update(); err != nil {
		return err
	}

	// Detect player death triggered by combat (Hurt → Dying/Dead state).
	if s.hasPlayer && s.player != nil && !s.kitScene.DeathActive() &&
		(s.player.State() == gamestates.Dying || s.player.State() == actors.Dead) {
		s.kitScene.StartDeathSequence()
	}

	s.deathTrigger.Update()
	if s.deathTrigger.Trigger() {
		s.AppContext().SceneManager.NavigateTo(
			scenestypes.ScenePhaseReboot,
			transition.NewFader(0, config.Get().FadeVisibleDuration),
			false,
		)
	}

	return s.TilemapScene.Update()
}

func (s *BeatemupPhaseScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 0xff})

	tilemap, err := s.Tilemap().Image(screen)
	if err != nil {
		log.Fatal(err)
	}
	s.Camera().Draw(tilemap, s.Tilemap().ImageOptions(), screen)

	s.kitScene.DrawActors(screen)
}

func (s *BeatemupPhaseScene) OnFinish() {
	s.TilemapScene.OnFinish()
	if s.hasPlayer {
		s.AppContext().ActorManager.Unregister(s.player)
	}
}
