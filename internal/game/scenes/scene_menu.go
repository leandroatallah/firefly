package gamescene

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
	"github.com/leandroatallah/firefly/internal/engine/core"
	"github.com/leandroatallah/firefly/internal/engine/core/scene"
	"github.com/leandroatallah/firefly/internal/engine/core/transition"
	"github.com/leandroatallah/firefly/internal/engine/systems/audiomanager"
	"github.com/leandroatallah/firefly/internal/engine/systems/speech"
	gamespeech "github.com/leandroatallah/firefly/internal/game/speech"
)

const (
	kickBackBG = "assets/kick_backOGG.ogg"
)

type MenuScene struct {
	scene.BaseScene

	audiomanager *audiomanager.AudioManager
	fontText     *font.FontText
	dialogue     *speech.Dialogue
}

func NewMenuScene(context *core.AppContext) *MenuScene {
	fontText, err := font.NewFontText(config.Get().MainFontFace)
	if err != nil {
		log.Fatal(err)
	}

	scene := MenuScene{fontText: fontText}
	scene.SetAppContext(context)
	return &scene
}

func (s *MenuScene) OnStart() {
	// Init audio
	s.audiomanager = s.Manager.AudioManager()
	s.audiomanager.SetVolume(1)
	// s.audiomanager.PlayMusic(kickBackBG)

	// Init Dialogue
	speechFont := speech.NewSpeechFont(s.fontText, 8, 14)
	bubble := gamespeech.NewSpeechBubble(speechFont)
	s.dialogue = speech.NewDialogue(bubble)
	s.dialogue.SetSpellingDelay(60)
	s.dialogue.SetLines([]string{
		"Look at those birds in \nthe sky. Are they beautiful?",
		"To be honest... This is \nvery boring...",
	})
}

func (s *MenuScene) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		s.Manager.NavigateTo(SceneLevels, transition.NewFader())
	}

	// Draw dialogue
	if err := s.dialogue.Update(); err != nil {
		return err
	}

	// Debug dialogue
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		s.dialogue.Show()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		if !s.dialogue.IsSpellingComplete() {
			s.dialogue.CompleteSpelling()
		} else {
			if ok := s.dialogue.NextLine(); !ok {
				s.dialogue.Hide()
			}
		}
	}

	return nil
}

func (s *MenuScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0xCC, 0x24, 0x40, 0xff})

	textOp := &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign:   text.AlignCenter,
			SecondaryAlign: text.AlignCenter,
			LineSpacing:    0,
		},
	}
	textOp.GeoM.Translate(
		float64(config.Get().ScreenWidth/2),
		float64(config.Get().ScreenHeight/2),
	)
	textOp.ColorScale.Scale(1, 1, 1, float32(120))
	s.fontText.Draw(screen, "Press Enter to start", 8, textOp)

	// Draw dialogue
	s.dialogue.Draw(screen)
}

func (s *MenuScene) OnFinish() {
	s.audiomanager.PauseMusic(kickBackBG)
}
