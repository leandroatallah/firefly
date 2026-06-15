package app

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/debug"
	"github.com/boilerplate/ebiten-template/internal/engine/ui/actorinspector"
	"github.com/boilerplate/ebiten-template/internal/engine/ui/debugoverlay"
	"github.com/boilerplate/ebiten-template/internal/engine/ui/phaseoverlay"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Game struct {
	AppContext      *AppContext
	debugFontFace   *text.GoTextFace
	slowMoApplied   bool
	lastSlowMo      bool
	lastFastForward bool
	debugOverlay    *debugoverlay.DebugOverlay
	phaseOverlay    *phaseoverlay.PhaseOverlay
	actorInspector  *actorinspector.Overlay
}

func NewGame(ctx *AppContext) *Game {
	debug.Init("assets/data/debug.json")
	return &Game{
		AppContext:   ctx,
		debugOverlay: debugoverlay.New(),
		phaseOverlay: phaseoverlay.New(),
		actorInspector: actorinspector.New(func() actorinspector.ActorSource {
			if ctx.ActorManager == nil {
				return nil
			}
			return ctx.ActorManager
		}),
	}
}

// DebugOverlay returns the debug overlay instance for external control (e.g. tests).
func (g *Game) DebugOverlay() *debugoverlay.DebugOverlay {
	return g.debugOverlay
}

// PhaseOverlay returns the phase-jump overlay instance for external wiring
// (entry list + select handler) and tests.
func (g *Game) PhaseOverlay() *phaseoverlay.PhaseOverlay {
	return g.phaseOverlay
}

// ActorInspector returns the actor inspector overlay instance for external wiring.
func (g *Game) ActorInspector() *actorinspector.Overlay {
	return g.actorInspector
}

func (g *Game) Update() error {
	cfg := g.AppContext.Config
	if !g.slowMoApplied || g.lastSlowMo != cfg.SlowMo || g.lastFastForward != cfg.FastForward {
		g.slowMoApplied = true
		g.lastSlowMo = cfg.SlowMo
		g.lastFastForward = cfg.FastForward
		// Slow-mo takes precedence over fast-forward when both are enabled.
		if tps, ok := EffectiveTPS(cfg.SlowMo, cfg.SlowMoFactor, ebiten.DefaultTPS); ok {
			ebiten.SetTPS(tps)
		} else if tps, ok := FastForwardTPS(cfg.FastForward, cfg.FastForwardFactor, ebiten.DefaultTPS); ok {
			ebiten.SetTPS(tps)
		} else {
			ebiten.SetTPS(ebiten.DefaultTPS)
		}
	}

	g.AppContext.FrameCount++

	f1Toggled := inpututil.IsKeyJustPressed(ebiten.KeyF1)
	if f1Toggled && !g.phaseOverlay.IsOpen() && !g.actorInspector.IsOpen() {
		if g.debugOverlay.IsOpen() {
			g.debugOverlay.Close()
		} else {
			g.debugOverlay.Open()
		}
	}

	if g.debugOverlay.IsOpen() {
		if !f1Toggled {
			g.debugOverlay.Update()
		}
		return nil
	}

	f2Toggled := inpututil.IsKeyJustPressed(ebiten.KeyF2)
	if f2Toggled && !g.actorInspector.IsOpen() {
		if g.phaseOverlay.IsOpen() {
			g.phaseOverlay.Close()
		} else {
			g.phaseOverlay.Open()
		}
	}

	if g.phaseOverlay.IsOpen() {
		if !f2Toggled {
			g.phaseOverlay.Update()
		}
		return nil
	}

	f5Toggled := inpututil.IsKeyJustPressed(ebiten.KeyF5)
	if f5Toggled {
		if g.actorInspector.IsOpen() {
			g.actorInspector.Close()
		} else {
			g.actorInspector.Open()
		}
	}

	if g.actorInspector.IsOpen() {
		if !f5Toggled {
			g.actorInspector.Update()
		}
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		cfg.SlowMo = !cfg.SlowMo
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF4) {
		cfg.FastForward = !cfg.FastForward
	}

	// Update Dialogue Manager
	if g.AppContext.DialogueManager != nil {
		g.AppContext.DialogueManager.Update()
	}

	// Then, update the current scene
	g.AppContext.SceneManager.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.AppContext.SceneManager.Draw(screen)

	// Draw Dialogue Manager
	if g.AppContext.DialogueManager != nil {
		g.AppContext.DialogueManager.Draw(screen)
	}

	g.AppContext.SceneManager.DrawOver(screen)

	g.debugOverlay.Draw(screen)
	g.phaseOverlay.Draw(screen)
	g.actorInspector.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	cfg := g.AppContext.Config
	return cfg.ScreenWidth, cfg.ScreenHeight
}

func (g *Game) DebugPhysics(screen *ebiten.Image) {
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

	if g.debugFontFace == nil {
		return
	}
	op := &text.DrawOptions{}
	op.GeoM.Translate(5, 15)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, b.String(), g.debugFontFace, op)
}
