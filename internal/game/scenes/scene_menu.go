package gamescene

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/scene"
	"github.com/leandroatallah/firefly/internal/engine/scene/transition"
	"github.com/leandroatallah/firefly/internal/engine/ui/menu"
	"github.com/leandroatallah/firefly/internal/engine/utils"
	"github.com/leandroatallah/firefly/internal/engine/utils/timing"
	scenestypes "github.com/leandroatallah/firefly/internal/game/scenes/types"
)

type MenuScene struct {
	scene.BaseScene

	fontText *font.FontText

	count              int
	isNavigating       bool
	navigationTrigger  utils.DelayTrigger
	shouldFadeOutSound bool
	isFadingOutSound   bool

	mainMenu    *menu.Menu
	optionsMenu *menu.Menu
}

func NewMenuScene(context *app.AppContext) *MenuScene {
	fontText, err := font.NewFontText(config.Get().MainFontFace)
	if err != nil {
		log.Fatal(err)
	}

	scene := MenuScene{fontText: fontText}
	scene.SetAppContext(context)
	scene.initMenus()
	return &scene
}

func (s *MenuScene) initMenus() {
	s.mainMenu = menu.NewMenu()
	s.mainMenu.SetFontSize(8)
	s.mainMenu.SetItemSpacing(12)
	s.optionsMenu = menu.NewMenu()
	s.optionsMenu.SetFontSize(8)
	s.optionsMenu.SetItemSpacing(12)

	// Main Menu
	s.mainMenu.AddItem("Game Start", func() {
		if !s.isNavigating {
			s.isNavigating = true
			s.navigationTrigger.Enable(timing.FromDuration(time.Second))
			s.mainMenu.SetVisible(false) // Hide menu to prevent double clicks
		}
	})
	s.mainMenu.AddItem("Options", func() {
		s.mainMenu.SetVisible(false)
		s.optionsMenu.SetVisible(true)
	})

	// Options Menu
	s.optionsMenu.AddItem("Back", func() {
		s.optionsMenu.SetVisible(false)
		s.mainMenu.SetVisible(true)
	})
	s.optionsMenu.AddItem("Language", func() {})
	s.optionsMenu.AddItem("Fullscreen", func() {})
}

func (s *MenuScene) OnStart() {
	am := s.AppContext().SceneManager.AudioManager()
	if am != nil {
		am.SetVolume(1)
		am.PlayMusic(TitleSound, true) // Loop menu music
	}

	// Reset state
	s.count = 0
	s.isNavigating = false
	s.isFadingOutSound = false
	s.shouldFadeOutSound = false
	s.navigationTrigger = utils.DelayTrigger{} // Reset trigger state

	// Reset menus
	s.mainMenu.SetVisible(true)
	s.optionsMenu.SetVisible(false)
}

func (s *MenuScene) Update() error {
	canSkipDelay := s.count > timing.FromDuration(time.Second)

	if canSkipDelay && !s.isNavigating {
		if s.optionsMenu.Visible() {
			s.optionsMenu.Update()
		} else if s.mainMenu.Visible() {
			s.mainMenu.Update()
		}
	}

	s.navigationTrigger.Update()
	if s.navigationTrigger.Trigger() {
		s.AppContext().SceneManager.NavigateTo(
			scenestypes.SceneStory, transition.NewFader(0, config.Get().FadeVisibleDuration), true,
		)
	}

	if s.isNavigating && s.shouldFadeOutSound && !s.isFadingOutSound {
		if s.AppContext().AudioManager != nil {
			s.AppContext().AudioManager.FadeOutAll(time.Second)
		}
		s.isFadingOutSound = true
	}

	if s.isNavigating {
		s.shouldFadeOutSound = true
	}

	s.count++

	return nil
}

func (s *MenuScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0xCC, 0x24, 0x40, 0xff})

	centerX := config.Get().ScreenWidth / 2
	centerY := config.Get().ScreenHeight / 2

	if s.optionsMenu.Visible() {
		s.optionsMenu.Draw(screen, s.fontText, centerX, centerY)
	} else if s.mainMenu.Visible() {
		s.mainMenu.Draw(screen, s.fontText, centerX, centerY)
	}
}

func (s *MenuScene) OnFinish() {}
