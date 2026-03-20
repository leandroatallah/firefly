package gamescene

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/render/screenutil"
	"github.com/leandroatallah/firefly/internal/engine/scene"
	"github.com/leandroatallah/firefly/internal/engine/scene/transition"
)

type PhaseTitleScene struct {
	scene.BaseScene

	fontText     *font.FontText
	title        string
	showTitle    bool
	musicStarted bool

	shouldInitMusic bool
}

func NewPhaseTitleScene(ctx *app.AppContext) *PhaseTitleScene {
	s := &PhaseTitleScene{
		fontText: ctx.Font,
	}
	s.SetAppContext(ctx)
	return s
}

func (s *PhaseTitleScene) OnStart() {
	phase, err := s.AppContext().PhaseManager.GetCurrentPhase()
	if err != nil {
		log.Printf("PhaseTitleScene: failed to get current phase: %v", err)
	}
	s.title = phase.Title
	s.AppContext().AudioManager.FadeOutAll(time.Second)
	s.shouldInitMusic = true
}

func (s *PhaseTitleScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	if s.showTitle {
		text := s.title
		if text == "" {
			text = "Phase"
		}
		screenutil.DrawCenteredText(screen, s.fontText, text, 16, color.White)
	}
}

func (s *PhaseTitleScene) Update() error {
	if err := s.BaseScene.Update(); err != nil {
		return err
	}

	ctx := s.AppContext()

	if s.shouldInitMusic {
		if am := ctx.AudioManager; am != nil {
			s.shouldInitMusic = false
			s.Schedule(2*time.Second, func() {
				s.showTitle = true
				am.SetVolume(1.0)
				am.PlayMusic(TitleSound, true) // Loop title music
			})
		}
	}

	if s.showTitle && inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		ctx.CompleteCurrentPhase(transition.NewFader(0, config.Get().FadeVisibleDuration), true)
	}
	return nil
}

func (s *PhaseTitleScene) OnFinish() {}
