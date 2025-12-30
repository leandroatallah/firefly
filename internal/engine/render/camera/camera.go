package camera

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
	"github.com/setanarut/kamera/v2"
)

type Controller struct {
	cam             *kamera.Camera
	target          body.Collidable
	followTarget    body.Body
	DeadZoneRadius  float64
	SmoothingFactor float64
}

func NewController(x, y float64) *Controller {
	cfg := config.Get()
	cam := kamera.NewCamera(x, y, float64(cfg.ScreenWidth), float64(cfg.ScreenHeight))
	cam.SmoothType = kamera.SmoothDamp
	cam.ShakeEnabled = true

	// Create a body to be the camera's direct target
	// targetBody := physics.NewPhysicsBody(physics.NewRect(0, 0, 1, 1))
	targetBody := bodyphysics.NewCollidableBodyFromRect(bodyphysics.NewRect(0, 0, 1, 1))

	return &Controller{
		cam:    cam,
		target: targetBody,
	}
}

func NewCamera(x, y int) *kamera.Camera {
	cfg := config.Get()
	c := kamera.NewCamera(
		float64(x),
		float64(x),
		float64(cfg.ScreenWidth),
		float64(cfg.ScreenHeight),
	)
	c.SmoothType = kamera.SmoothDamp
	c.ShakeEnabled = true
	return c
}

func (c *Controller) SetFollowTarget(b body.Body) {
	c.followTarget = b
}

func (c *Controller) Update() {
	c.cam.LookAt(
		float64(c.followTarget.Position().Min.X),
		float64(c.followTarget.Position().Min.Y),
	)
}

func (c *Controller) Draw(
	src *ebiten.Image, options *ebiten.DrawImageOptions, dst *ebiten.Image,
) {
	c.cam.Draw(src, options, dst)
}

// Useful for debugging
func (c *Controller) Kamera() *kamera.Camera {
	return c.cam
}

func (c *Controller) Position() image.Rectangle {
	// return c.target.Position()
	return c.followTarget.Position()
}

func (c *Controller) Target() body.Body {
	// return c.target
	return c.followTarget
}
