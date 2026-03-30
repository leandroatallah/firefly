package gamescene

import (
	"image/color"
	"time"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/scene"
	"github.com/boilerplate/ebiten-template/internal/engine/scene/transition"
	"github.com/boilerplate/ebiten-template/internal/engine/utils"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/timing"
	"github.com/hajimehoshi/ebiten/v2"
)

type PhaseRebootScene struct {
	scene.BaseScene

	count             int
	navigationTrigger utils.DelayTrigger
}

func NewPhaseRebootScene(context *app.AppContext) *PhaseRebootScene {
	overlay := ebiten.NewImage(config.Get().ScreenWidth, config.Get().ScreenHeight)
	overlay.Fill(color.Black)
	scene := PhaseRebootScene{}
	scene.SetAppContext(context)
	return &scene
}

func (s *PhaseRebootScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{A: 255})
}

func (s *PhaseRebootScene) OnStart() {
	// Freeze all actors and items to preserve state during reboot
	s.freezeAllActors()
	s.navigationTrigger.Enable(timing.FromDuration(167 * time.Millisecond))
}

func (s *PhaseRebootScene) freezeAllActors() {
	if s.AppContext().ActorManager != nil {
		s.AppContext().ActorManager.ForEach(func(actor actors.ActorEntity) {
			actor.SetFreeze(true)
		})
	}
}

func (s *PhaseRebootScene) Update() error {
	s.count++
	s.navigationTrigger.Update()

	if s.navigationTrigger.Trigger() {
		s.AppContext().SceneManager.NavigateBack(transition.NewFader(0, config.Get().FadeVisibleDuration))
	}

	return nil
}

func (s *PhaseRebootScene) OnFinish() {}
