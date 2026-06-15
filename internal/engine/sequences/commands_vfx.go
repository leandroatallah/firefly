package sequences

import (
	"image/color"
	"log"
	"math"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	contractvfx "github.com/boilerplate/ebiten-template/internal/engine/contracts/vfx"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/render/camera"
)

// FadeOutCommand fades the screen to black over a given number of frames and keeps it black.
// The fade persists until Reset() is called on the overlay or a FadeInCommand is used.
type FadeOutCommand struct {
	Frames int
	Color  color.RGBA
	ctx    *app.AppContext
}

func (c *FadeOutCommand) Init(appContext any) {
	c.ctx = appContext.(*app.AppContext)
	if c.Frames <= 0 {
		c.Frames = 17
	}
	if c.ctx != nil && c.ctx.FadeOverlay != nil {
		c.ctx.FadeOverlay.FadeOut(c.Frames)
	}
}

func (c *FadeOutCommand) Update() bool {
	if c.ctx == nil || c.ctx.FadeOverlay == nil {
		return true
	}
	// Command done when animation completes (not animating anymore)
	return !c.ctx.FadeOverlay.IsActive()
}

// FadeInCommand fades the screen from black over a given number of frames.
// Intended to be used after a FadeOutCommand or when the overlay is at full alpha.
type FadeInCommand struct {
	Frames int
	Color  color.RGBA
	ctx    *app.AppContext
}

func (c *FadeInCommand) Init(appContext any) {
	c.ctx = appContext.(*app.AppContext)
	if c.Frames <= 0 {
		c.Frames = 17
	}
	if c.ctx != nil && c.ctx.FadeOverlay != nil {
		c.ctx.FadeOverlay.FadeIn(c.Frames)
	}
}

func (c *FadeInCommand) Update() bool {
	if c.ctx == nil || c.ctx.FadeOverlay == nil {
		return true
	}
	return !c.ctx.FadeOverlay.IsActive()
}

// SolidColorCommand cover the screen to a solid color (default is black) over a given number of frames.
// The overlay persists until Reset() is called on the overlay.
type SolidColorCommand struct {
	Frames int
	ctx    *app.AppContext
	Color  color.RGBA
}

func (c *SolidColorCommand) Init(appContext any) {
	c.ctx = appContext.(*app.AppContext)
	if c.Frames <= 0 {
		c.Frames = 17
	}
	if c.ctx != nil && c.ctx.SolidColorOverlay != nil {
		c.ctx.SolidColorOverlay.SetColor(c.Color)
		c.ctx.SolidColorOverlay.FadeOut(c.Frames)
	}
}

func (c *SolidColorCommand) Update() bool {
	if c.ctx == nil || c.ctx.SolidColorOverlay == nil {
		return true
	}
	// Command done when animation completes (not animating anymore)
	return !c.ctx.SolidColorOverlay.IsActive()
}

type SpawnTextCommand struct {
	TargetID string `json:"target_id,omitempty"`
	Text     string `json:"text"`
	Duration int    `json:"duration"`
	Type     string `json:"type"`
	X, Y     float64
}

func (c *SpawnTextCommand) Init(appContext any) {
	ctx := appContext.(*app.AppContext)

	if c.Type == "screen" {
		ctx.VFX.SpawnFloatingText(c.Text, c.X, c.Y, c.Duration)
		return
	}

	if c.TargetID == "" {
		log.Printf("SpawnTextCommand: target_id required for overhead text")
		return
	}

	actor, found := ctx.ActorManager.Find(c.TargetID)
	if !found {
		log.Printf("SpawnTextCommand: actor not found: %s", c.TargetID)
		return
	}

	ctx.VFX.SpawnFloatingTextAbove(actor, c.Text, c.Duration)
}

func (c *SpawnTextCommand) Update() bool {
	return true
}

// QuakeCommand triggers a screen shake and falling rocks effect.
type QuakeCommand struct {
	Trauma   float64
	Duration int
	camera   interface{ AddTrauma(float64) }
	vfx      contractvfx.Manager
	timer    int
}

func (c *QuakeCommand) Init(appContext any) {
	ctx := appContext.(*app.AppContext)
	c.vfx = ctx.VFX
	c.timer = 0

	// Get camera
	currentScene := ctx.SceneManager.CurrentScene()
	if tilemapScene, ok := currentScene.(interface{ Camera() *camera.Controller }); ok {
		c.camera = tilemapScene.Camera()
	}
}

func (c *QuakeCommand) Update() bool {
	if c.camera != nil && c.timer%10 == 0 {
		c.camera.AddTrauma(c.Trauma)
	}

	c.timer++
	return c.timer >= c.Duration
}

// VignetteRadiusCommand animates the vignette radius over a given duration.
// The actual vignette implementation lives in the game layer; this command
// simply calls into scenes that expose the expected interface.
type VignetteRadiusCommand struct {
	InitialRadius float64
	FinalRadius   float64
	Duration      int

	frame      int
	controller interface {
		EnableVignetteDarkness(radiusPx float64)
		DisableVignetteDarkness()
	}
}

func (c *VignetteRadiusCommand) Init(appContext any) {
	ctx, ok := appContext.(*app.AppContext)
	if !ok || ctx == nil {
		return
	}

	// Interpret negative FinalRadius as a request for full-screen radius.
	// Compute this from the current config so it's resolution-independent.
	if c.FinalRadius < 0 {
		if cfg := config.Get(); cfg != nil {
			w := float64(cfg.ScreenWidth)
			h := float64(cfg.ScreenHeight)
			// Use the screen diagonal as a safe "cover everything" radius.
			c.FinalRadius = math.Hypot(w, h)
		}
	}

	currentScene := ctx.SceneManager.CurrentScene()

	if currentScene == nil {
		return
	}

	if v, ok := currentScene.(interface {
		EnableVignetteDarkness(radiusPx float64)
		DisableVignetteDarkness()
	}); ok {
		c.controller = v
		// Apply initial radius immediately if non-zero.
		if c.InitialRadius > 0 {
			c.controller.EnableVignetteDarkness(c.InitialRadius)
		} else if c.FinalRadius <= 0 {
			// Both initial and final are zero or negative: ensure vignette is off.
			c.controller.DisableVignetteDarkness()
		}
	}
}

func (c *VignetteRadiusCommand) Update() bool {
	if c.controller == nil {
		// Nothing to control; treat as instant.
		return true
	}

	// No duration or negative: jump directly to final state.
	if c.Duration <= 0 {
		if c.FinalRadius > 0 {
			c.controller.EnableVignetteDarkness(c.FinalRadius)
		} else {
			c.controller.DisableVignetteDarkness()
		}
		return true
	}

	if c.frame >= c.Duration {
		// Ensure we end exactly at the final radius and keep vignette enabled
		// when the final radius is non-zero.
		if c.FinalRadius > 0 {
			c.controller.EnableVignetteDarkness(c.FinalRadius)
		} else {
			c.controller.DisableVignetteDarkness()
		}
		return true
	}

	progress := float64(c.frame) / float64(c.Duration)
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}

	current := c.InitialRadius + (c.FinalRadius-c.InitialRadius)*progress
	if current > 0 {
		c.controller.EnableVignetteDarkness(current)
	} else {
		c.controller.DisableVignetteDarkness()
	}

	c.frame++
	return false
}
