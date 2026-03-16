package sequences

import (
	"log"
	"math"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/render/camera"
	"github.com/leandroatallah/firefly/internal/engine/render/particles/vfx"
)

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
	vfx      *vfx.Manager
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
