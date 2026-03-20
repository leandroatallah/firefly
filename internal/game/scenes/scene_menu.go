package gamescene

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"
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
	// Game start
	s.mainMenu.AddItem("", func() {
		if !s.isNavigating {
			s.isNavigating = true
			s.navigationTrigger.Enable(timing.FromDuration(time.Second))
			s.mainMenu.SetVisible(false) // Hide menu to prevent double clicks
		}
	})
	// Options
	s.mainMenu.AddItem("", func() {
		s.mainMenu.SetVisible(false)
		s.optionsMenu.SetVisible(true)
	})
	// Exit
	s.mainMenu.AddItem("", func() {
		// os.Exit(0) is acceptable for a simple game.
		os.Exit(0)
	})

	// Options Menu
	// Back
	s.optionsMenu.AddItem("", func() {
		s.optionsMenu.SetVisible(false)
		s.mainMenu.SetVisible(true)
	})
	// Language
	s.optionsMenu.AddItem("", func() {
		cfg := config.Get()
		if cfg.Language == "en" {
			cfg.Language = "pt-br"
		} else {
			cfg.Language = "en"
		}
		s.AppContext().I18n.Load(cfg.Language)
		s.refreshMenuLabels()
	})

	// Fullscreen
	s.optionsMenu.AddItem("", func() {
		cfg := config.Get()
		cfg.Fullscreen = !cfg.Fullscreen
		ebiten.SetFullscreen(cfg.Fullscreen)
		s.refreshMenuLabels()
	})

	s.refreshMenuLabels()
}

func (s *MenuScene) refreshMenuLabels() {
	i18n := s.AppContext().I18n
	cfg := config.Get()

	// Main Menu
	s.mainMenu.UpdateItemLabel(0, i18n.T("menu_game_start"))
	s.mainMenu.UpdateItemLabel(1, i18n.T("menu_options"))
	s.mainMenu.UpdateItemLabel(2, i18n.T("menu_exit"))

	// Options Menu
	s.optionsMenu.UpdateItemLabel(0, i18n.T("options_back"))
	s.optionsMenu.UpdateItemLabel(1, fmt.Sprintf("%s: %s", i18n.T("options_language"), strings.ToUpper(cfg.Language)))

	fullscreenKey := "options_fullscreen_off"
	if cfg.Fullscreen {
		fullscreenKey = "options_fullscreen_on"
	}
	s.optionsMenu.UpdateItemLabel(2, i18n.T(fullscreenKey))
}

func (s *MenuScene) OnStart() {
	s.Schedule(2*time.Second, func() {
		am := s.AppContext().SceneManager.AudioManager()
		if am != nil {
			am.SetVolume(1)
			am.PlayMusic(TitleSound, true) // Loop menu music
		}
	})

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
	if err := s.BaseScene.Update(); err != nil {
		return err
	}

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
