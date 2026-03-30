package gamescene

import (
	"fmt"
	"image/color"
	"os"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/assets/font"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/render/camera"
	"github.com/boilerplate/ebiten-template/internal/engine/scene"
	"github.com/boilerplate/ebiten-template/internal/engine/scene/transition"
	"github.com/boilerplate/ebiten-template/internal/engine/ui/menu"
	"github.com/boilerplate/ebiten-template/internal/engine/utils"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/timing"
	scenestypes "github.com/boilerplate/ebiten-template/internal/game/scenes/types"
)

type MenuScene struct {
	scene.BaseScene

	fontText  *font.FontText
	fontSmall *font.FontText

	count             int
	isNavigating      bool
	navigationTrigger utils.DelayTrigger

	mainMenu    *menu.Menu
	optionsMenu *menu.Menu
	camera      *camera.Controller
}

func NewMenuScene(context *app.AppContext) *MenuScene {
	scene := MenuScene{
		fontText:  context.Font,
		fontSmall: context.Font,
		camera:    camera.NewController(0, 0),
	}
	scene.SetAppContext(context)
	scene.initMenus()
	return &scene
}

func (s *MenuScene) initMenus() {
	s.mainMenu = menu.NewMenu()
	s.mainMenu.SetFontSize(8)
	s.mainMenu.SetItemSpacing(8)

	s.optionsMenu = menu.NewMenu()
	s.optionsMenu.SetFontSize(8)
	s.optionsMenu.SetItemSpacing(8)

	s.mainMenu.SetOnNavigate(func() { s.camera.AddTrauma(0.2) })
	s.optionsMenu.SetOnNavigate(func() { s.camera.AddTrauma(0.2) })

	// Main Menu Items
	s.mainMenu.AddItem("", func() {
		if !s.isNavigating {
			s.isNavigating = true
			s.navigationTrigger.Enable(timing.FromDuration(time.Second))
			s.mainMenu.SetVisible(false)
			s.AppContext().AudioManager.FadeOutCurrentTrack(time.Second)
		}
	})
	s.mainMenu.AddItem("", func() {
		s.mainMenu.SetVisible(false)
		s.optionsMenu.SetVisible(true)
	})
	s.mainMenu.AddItem("", func() {
		os.Exit(0)
	})

	// Options Menu Items
	s.optionsMenu.AddItem("", func() {
		s.optionsMenu.SetVisible(false)
		s.mainMenu.SetVisible(true)
	})
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

	s.mainMenu.UpdateItemLabel(0, i18n.T("menu.start"))
	s.mainMenu.UpdateItemLabel(1, i18n.T("menu.options"))
	s.mainMenu.UpdateItemLabel(2, i18n.T("menu.exit"))

	s.optionsMenu.UpdateItemLabel(0, i18n.T("options.back"))
	s.optionsMenu.UpdateItemLabel(1, fmt.Sprintf("%s: %s", i18n.T("options.language"), strings.ToUpper(cfg.Language)))

	fullscreenKey := "options.fullscreen_off"
	if cfg.Fullscreen {
		fullscreenKey = "options.fullscreen_on"
	}
	s.optionsMenu.UpdateItemLabel(2, i18n.T(fullscreenKey))
}

func (s *MenuScene) OnStart() {
	s.count = 0
	s.isNavigating = false
	s.navigationTrigger = utils.DelayTrigger{}

	s.mainMenu.SetVisible(false)
	s.optionsMenu.SetVisible(false)

	s.camera.SetCenter(float64(config.Get().ScreenWidth/2), float64(config.Get().ScreenHeight/2))

	s.PlayMusicWithLoop("assets/audio/music/Goblins_Den_Regular.ogg", true, false)
}

func (s *MenuScene) Update() error {
	if err := s.BaseScene.Update(); err != nil {
		return err
	}

	if s.count == timing.FromDuration(time.Second) {
		s.mainMenu.SetVisible(true)
	}

	if !s.isNavigating {
		if s.optionsMenu.Visible() {
			s.optionsMenu.Update()
		} else if s.mainMenu.Visible() {
			s.mainMenu.Update()
		}
	}

	s.navigationTrigger.Update()
	if s.navigationTrigger.Trigger() {
		s.AppContext().SceneManager.NavigateTo(
			scenestypes.ScenePhases, transition.NewFader(0, config.Get().FadeVisibleDuration), true,
		)
	}

	s.camera.Update()
	s.count++

	return nil
}

func (s *MenuScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	centerX := config.Get().ScreenWidth / 2
	centerY := config.Get().ScreenHeight / 2

	sceneLayer := ebiten.NewImage(config.Get().ScreenWidth, config.Get().ScreenHeight)

	if s.fontSmall != nil {
		titleOp := &text.DrawOptions{}
		titleOp.ColorScale.ScaleWithColor(color.White)
		titleFace := s.fontSmall.NewFace(32)
		titleText := "Boilerplate"
		textWidth, _ := text.Measure(titleText, titleFace, 0)
		titleOp.GeoM.Translate(float64(centerX)-textWidth/2, float64(centerY)-90)
		s.fontSmall.Draw(sceneLayer, titleText, 32, titleOp)
	}

	menuY := centerY + 20
	if s.optionsMenu.Visible() {
		s.optionsMenu.Draw(sceneLayer, s.fontText, centerX, menuY)
	} else if s.mainMenu.Visible() {
		s.mainMenu.Draw(sceneLayer, s.fontText, centerX, menuY)
	}

	opts := &ebiten.DrawImageOptions{}
	s.camera.Draw(sceneLayer, opts, screen)
}

func (s *MenuScene) OnFinish() {}
