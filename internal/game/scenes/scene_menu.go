package gamescene

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/render/camera"
	"github.com/leandroatallah/firefly/internal/engine/render/particles"
	"github.com/leandroatallah/firefly/internal/engine/scene"
	"github.com/leandroatallah/firefly/internal/engine/scene/transition"
	"github.com/leandroatallah/firefly/internal/engine/ui/menu"
	"github.com/leandroatallah/firefly/internal/engine/utils"
	"github.com/leandroatallah/firefly/internal/engine/utils/timing"
	scenestypes "github.com/leandroatallah/firefly/internal/game/scenes/types"
)

type MenuScene struct {
	scene.BaseScene

	fontText  *font.FontText
	fontSmall *font.FontText

	count              int
	isNavigating       bool
	navigationTrigger  utils.DelayTrigger
	shouldFadeOutSound bool
	isFadingOutSound   bool
	musicStarted       bool

	mainMenu    *menu.Menu
	optionsMenu *menu.Menu

	camera          *camera.Controller
	particles       *particles.System
	bgParticles     *particles.System
	bgParticleTimer int
}

func NewMenuScene(context *app.AppContext) *MenuScene {
	fontText, err := font.NewFontText(config.Get().MainFontFace)
	if err != nil {
		log.Fatal(err)
	}

	fontSmall, err := font.NewFontText(config.Get().SmallFontFace)
	if err != nil {
		log.Fatal(err)
	}

	scene := MenuScene{
		fontText:        fontText,
		fontSmall:       fontSmall,
		camera:          camera.NewController(0, 0),
		particles:       particles.NewSystem(),
		bgParticles:     particles.NewSystem(),
		bgParticleTimer: 0,
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

	// Set up navigation callback for screen shake
	s.mainMenu.SetOnNavigate(func() {
		s.camera.AddTrauma(0.5)
	})
	s.optionsMenu.SetOnNavigate(func() {
		s.camera.AddTrauma(0.5)
	})

	// Set up selection callback for particle effect (only on Game Start)
	s.mainMenu.SetOnSelect(func() {
		// Only spawn particles for Game Start (first item)
		if s.mainMenu.SelectedIndex() == 0 {
			s.spawnParticles()
		}
	})

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
	// Reset state
	s.count = 0
	s.isNavigating = false
	s.isFadingOutSound = false
	s.shouldFadeOutSound = false
	s.navigationTrigger = utils.DelayTrigger{} // Reset trigger state
	s.musicStarted = false

	// Reset menus
	s.mainMenu.SetVisible(false)
	s.optionsMenu.SetVisible(false)

	// Reset camera
	s.camera.SetCenter(float64(config.Get().ScreenWidth/2), float64(config.Get().ScreenHeight/2))
}

// spawnParticles creates a burst of particles from the center of the screen
func (s *MenuScene) spawnParticles() {
	centerX := float64(config.Get().ScreenWidth / 2)
	centerY := float64(config.Get().ScreenHeight / 2)

	// Create particle burst with random movement
	for i := 0; i < 40; i++ {
		// More random velocity distribution
		velX := (rand.Float64() - 0.5) * 6.0
		velY := (rand.Float64() - 0.5) * 6.0

		p := &particles.Particle{
			X:           centerX + (rand.Float64()-0.5)*20,
			Y:           centerY + (rand.Float64()-0.5)*20,
			VelX:        velX,
			VelY:        velY,
			Duration:    40 + rand.Intn(40),
			MaxDuration: 80,
			Scale:       1.0 + rand.Float64()*3.0,
			ScaleSpeed:  -0.01 - rand.Float64()*0.02,
			Config:      s.getPixelConfig(),
		}
		// White color with varying opacity
		opacity := uint8(150 + rand.Intn(105))
		p.ColorScale.ScaleWithColor(color.RGBA{opacity, opacity, opacity, 255})
		s.particles.Add(p)
	}
}

// spawnBackgroundParticle creates a single subtle background particle
func (s *MenuScene) spawnBackgroundParticle() {
	p := &particles.Particle{
		X:           rand.Float64() * float64(config.Get().ScreenWidth),
		Y:           float64(config.Get().ScreenHeight) + 5,
		VelX:        (rand.Float64() - 0.5) * 0.3,
		VelY:        -0.2 - rand.Float64()*0.3,
		Duration:    180 + rand.Intn(120),
		MaxDuration: 300,
		Scale:       0.5 + rand.Float64()*1.0,
		ScaleSpeed:  0,
		Config:      s.getPixelConfig(),
	}
	// Subtle white/gray with visible opacity
	opacity := uint8(60 + rand.Intn(80))
	p.ColorScale.ScaleWithColor(color.RGBA{opacity, opacity, opacity, 255})
	s.bgParticles.Add(p)
}

// getPixelConfig returns a 1x1 pixel particle configuration
func (s *MenuScene) getPixelConfig() *particles.Config {
	pixelImg := ebiten.NewImage(1, 1)
	pixelImg.Fill(color.White)
	return &particles.Config{
		Image:       pixelImg,
		FrameWidth:  1,
		FrameHeight: 1,
		FrameCount:  1,
	}
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
			scenestypes.SceneStory, transition.NewFader(0, config.Get().FadeVisibleDuration), true,
		)
	}

	// Start music only after transition is finished
	if !s.musicStarted && !s.AppContext().SceneManager.IsTransitioning() {
		s.musicStarted = true
		am := s.AppContext().SceneManager.AudioManager()
		if am != nil {
			am.SetVolume(1)
			am.PlayMusic(TitleSound, true) // Loop menu music
		}
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

	// Update particles
	s.particles.Update()
	s.bgParticles.Update()

	// Spawn background particles periodically (every 10 frames)
	s.bgParticleTimer++
	if s.bgParticleTimer >= 10 {
		s.bgParticleTimer = 0
		s.spawnBackgroundParticle()
	}

	// Update camera
	s.camera.Update()

	s.count++

	return nil
}

func (s *MenuScene) Draw(screen *ebiten.Image) {
	// Fill screen with black first
	screen.Fill(color.Black)

	centerX := config.Get().ScreenWidth / 2
	centerY := config.Get().ScreenHeight / 2

	// Draw background particles directly to screen (no camera shake, floating effect)
	s.drawBackgroundParticles(screen)

	// Create a layer for scene content that will have camera shake applied
	sceneLayer := ebiten.NewImage(config.Get().ScreenWidth, config.Get().ScreenHeight)

	// Draw title "Growbel" using small font (32px) to scene layer
	if s.fontSmall != nil {
		titleOp := &text.DrawOptions{}
		titleOp.ColorScale.ScaleWithColor(color.White)
		titleFace := s.fontSmall.NewFace(32)
		textWidth, _ := text.Measure("Growbel", titleFace, 0)
		titleOp.GeoM.Translate(float64(centerX)-textWidth/2, float64(centerY)-90)
		s.fontSmall.Draw(sceneLayer, "Growbel", 32, titleOp)
	}

	// Draw menus to scene layer
	menuY := centerY + 20
	if s.optionsMenu.Visible() {
		s.optionsMenu.Draw(sceneLayer, s.fontText, centerX, menuY)
	} else if s.mainMenu.Visible() {
		s.mainMenu.Draw(sceneLayer, s.fontText, centerX, menuY)
	}

	// Draw selection particles to scene layer
	s.particles.Draw(sceneLayer, s.camera)

	// Apply camera shake when drawing scene layer to final screen
	opts := &ebiten.DrawImageOptions{}
	s.camera.Draw(sceneLayer, opts, screen)
}

// drawBackgroundParticles draws background particles without camera shake
func (s *MenuScene) drawBackgroundParticles(screen *ebiten.Image) {
	for _, p := range s.bgParticles.Particles() {
		if p.Config.Image == nil {
			continue
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(p.Scale, p.Scale)
		op.GeoM.Translate(p.X, p.Y)

		// Center the particle
		w, h := p.Config.FrameWidth, p.Config.FrameHeight
		if w == 0 {
			w = p.Config.Image.Bounds().Dx()
		}
		if h == 0 {
			h = p.Config.Image.Bounds().Dy()
		}
		op.GeoM.Translate(-float64(w)*p.Scale/2, -float64(h)*p.Scale)
		op.ColorScale = p.ColorScale

		screen.DrawImage(p.Config.Image, op)
	}
}

func (s *MenuScene) OnFinish() {}
