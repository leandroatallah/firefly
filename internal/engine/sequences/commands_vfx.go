package sequences

import (
	"log"

	"github.com/leandroatallah/firefly/internal/engine/app"
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
