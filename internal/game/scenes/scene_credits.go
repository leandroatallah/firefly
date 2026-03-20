package gamescene

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/scene"
)

// CreditEntry represents a single credit entry
type CreditEntry struct {
	title string
	name  string
}

// CreditsScene displays the game credits
type CreditsScene struct {
	scene.BaseScene

	fontTitle *font.FontText
	fontText  *font.FontText

	credits      []CreditEntry
	scrollOffset float64
	scrollSpeed  float64
}

// NewCreditsScene creates a new credits scene
func NewCreditsScene(context *app.AppContext) *CreditsScene {
	fontTitle, err := font.NewFontText(config.Get().MainFontFace)
	if err != nil {
		log.Fatal(err)
	}

	fontText, err := font.NewFontText(config.Get().SmallFontFace)
	if err != nil {
		log.Fatal(err)
	}

	s := &CreditsScene{
		fontTitle:   fontTitle,
		fontText:    fontText,
		scrollSpeed: 0.5,
		credits:     buildCredits(),
	}
	s.SetAppContext(context)
	return s
}

// buildCredits builds the list of credits entries
func buildCredits() []CreditEntry {
	credits := []CreditEntry{
		// Game Title
		{"GROWBEL", ""},
		{"", ""},

		// Development
		{"DEVELOPMENT", ""},
		{"Game Design", "Leandro Atallah"},
		{"Programming", "Leandro Atallah"},
		{"", ""},

		// Music
		{"MUSIC", ""},
		{"Goblins Dance Battle", "Kevin MacLeod"},
		{"Goblins Den Regular", "Kevin MacLeod"},
		{"", ""},

		// Sound Effects
		{"SOUND EFFECTS", ""},
		{"Dialogue Bleeps", "dmochas"},
		{"", ""},

		// Fonts
		{"FONTS", ""},
		{"Press Start 2P", "Codenexus"},
		{"Tiny5", "Stefan Schmidt"},
		{"", ""},

		// Graphics
		{"GRAPHICS", ""},
		{"Player Sprites", "Kenney"},
		{"Enemy Sprites", "Leandro Atallah"},
		{"Tileset", "Leandro Atallah"},
		{"UI Elements", "Leandro Atallah"},
		{"", ""},

		// Tools & Libraries
		{"ENGINE & TOOLS", ""},
		{"Ebitengine", "Hajime Hoshi"},
		{"", ""},

		// Special Thanks
		{"SPECIAL THANKS", ""},
		{"You!", "For Playing"},
		{"", ""},
	}
	return credits
}

func (s *CreditsScene) OnStart() {
	s.BaseScene.OnStart()
	s.scrollOffset = float64(config.Get().ScreenHeight)

	am := s.AppContext().SceneManager.AudioManager()
	if am != nil {
		am.SetVolume(1)
		am.PlayMusic(StorySound, true) // Loop menu music
	}
}

func (s *CreditsScene) Update() error {
	s.BaseScene.Update()

	// Auto-scroll credits
	s.scrollOffset -= s.scrollSpeed

	return nil
}

func (s *CreditsScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{A: 255})

	centerX := config.Get().ScreenWidth / 2
	baseY := s.scrollOffset

	// Draw credits
	y := baseY
	for _, credit := range s.credits {
		if credit.title == "" && credit.name == "" {
			// Empty line for spacing
			y += 20
			continue
		}

		if credit.name == "" {
			// Section title - centered
			op := &text.DrawOptions{
				LayoutOptions: text.LayoutOptions{
					PrimaryAlign: text.AlignCenter,
				},
			}
			op.ColorScale.ScaleWithColor(color.White)
			op.GeoM.Translate(float64(centerX), y)
			s.fontTitle.Draw(screen, credit.title, 8, op)
			y += 30
		} else {
			// Credit entry - left and right columns
			// Measure text to calculate positions
			face := s.fontText.NewFace(8)

			// Title (left side) - right aligned
			titleWidth, _ := text.Measure(credit.title, face, 0)
			opLeft := &text.DrawOptions{}
			opLeft.ColorScale.ScaleWithColor(color.White)
			opLeft.GeoM.Translate(float64(centerX-10)-titleWidth, y)
			s.fontText.Draw(screen, credit.title, 8, opLeft)

			// Name (right side) - left aligned
			opRight := &text.DrawOptions{}
			opRight.ColorScale.ScaleWithColor(color.White)
			opRight.GeoM.Translate(float64(centerX+10), y)
			s.fontText.Draw(screen, credit.name, 8, opRight)
			y += 24
		}
	}
}

func (s *CreditsScene) OnFinish() {
	s.BaseScene.OnFinish()
}
