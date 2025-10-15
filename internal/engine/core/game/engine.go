package game

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/engine/core"
	"github.com/leandroatallah/firefly/internal/engine/core/game/state"
	"golang.org/x/image/font"
)

type Game struct {
	AppContext    *core.AppContext
	state         state.GameState
	debugVisible  bool
	debugFontFace font.Face
}

func NewGame(ctx *core.AppContext) *Game {
	fontData, err := os.ReadFile(config.Get().MainFontFace)
	if err != nil {
		log.Fatal(err)
	}
	tt, err := truetype.Parse(fontData)
	if err != nil {
		log.Fatal(err)
	}

	return &Game{
		AppContext: ctx,
		debugFontFace: truetype.NewFace(tt, &truetype.Options{
			Size:    8,
			DPI:     72,
			Hinting: font.HintingFull,
		}),
	}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		g.debugVisible = !g.debugVisible
	}

	// First, update the input manager
	g.AppContext.InputManager.Update()

	// Update Dialogue Manager
	g.AppContext.DialogueManager.Update()

	// Then, update the current scene
	g.AppContext.SceneManager.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.AppContext.SceneManager.Draw(screen)

	// Draw Dialogue Manager
	g.AppContext.DialogueManager.Draw(screen)

	if g.debugVisible {
		cfg := config.Get().Physics
		var b strings.Builder
		fmt.Fprintf(&b, "--- Physics Debug ---\n")
		fmt.Fprintf(&b, "HorizontalInertia: %.2f\n", cfg.HorizontalInertia)
		fmt.Fprintf(&b, "AirFrictionMultiplier: %.2f\n", cfg.AirFrictionMultiplier)
		fmt.Fprintf(&b, "AirControlMultiplier: %.2f\n", cfg.AirControlMultiplier)
		fmt.Fprintf(&b, "CoyoteTimeFrames: %d\n", cfg.CoyoteTimeFrames)
		fmt.Fprintf(&b, "JumpBufferFrames: %d\n", cfg.JumpBufferFrames)
		fmt.Fprintf(&b, "JumpForce: %d\n", cfg.JumpForce)
		fmt.Fprintf(&b, "JumpCutMultiplier: %.2f\n", cfg.JumpCutMultiplier)
		fmt.Fprintf(&b, "UpwardGravity: %d\n", cfg.UpwardGravity)
		fmt.Fprintf(&b, "DownwardGravity: %d\n", cfg.DownwardGravity)
		fmt.Fprintf(&b, "MaxFallSpeed: %d\n", cfg.MaxFallSpeed)

		text.Draw(screen, b.String(), g.debugFontFace, 5, 15, color.White)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.Get().ScreenWidth, config.Get().ScreenHeight
}

func (g *Game) SetState(stateID state.GameStateEnum) error {
	// state, err := state.NewGameState(stateID, g.AppContext)
	// if err != nil {
	// 	return err
	// }
	//
	// g.state = state
	// g.state.OnStart()

	return nil
}
