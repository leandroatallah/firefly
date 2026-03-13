package gamescene

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/scene"
	"github.com/leandroatallah/firefly/internal/engine/scene/transition"
	"github.com/leandroatallah/firefly/internal/engine/utils"
	"github.com/leandroatallah/firefly/internal/engine/utils/timing"
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
	s.navigationTrigger.Enable(timing.FromDuration(167 * time.Millisecond))
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
