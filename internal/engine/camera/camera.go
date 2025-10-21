package camera

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
	"github.com/setanarut/kamera/v2"
)

type Controller struct {
	cam             *kamera.Camera
	target          body.Body
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
	targetBody := physics.NewPhysicsBody(physics.NewRect(0, 0, 1, 1))

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
	cfg := config.Get()
	c.followTarget = b
	pPos := c.followTarget.Position().Min
	c.target.SetPosition(pPos.X, pPos.Y) // game coords
	c.cam.LookAt(float64(pPos.X*cfg.Unit), float64(pPos.Y*cfg.Unit))
}

func (c *Controller) Update() {
	// Update cam target to smoothly follow the player
	pPos := c.followTarget.Position().Min
	targetPos := c.target.Position().Min

	// A smaller factor makes the movement smoother (and slower).
	newX := float64(targetPos.X) + (float64(pPos.X)-float64(targetPos.X))*c.SmoothingFactor
	newY := float64(targetPos.Y) + (float64(pPos.Y)-float64(targetPos.Y))*c.SmoothingFactor

	c.target.SetPosition(int(newX)*config.Get().Unit, int(newY)*config.Get().Unit)

	// Update camera to look at the now smoothly moving camTarget
	finalTargetPos := c.target.Position().Min
	targetWidth := c.target.Position().Dx()
	targetHeight := c.target.Position().Dy()
	c.cam.LookAt(
		float64(finalTargetPos.X+(targetWidth/2)),
		float64(finalTargetPos.Y+(targetHeight/2)),
	)
}

func (c *Controller) Draw(
	dst *ebiten.Image, options *ebiten.DrawImageOptions, src *ebiten.Image,
) {
	c.cam.Draw(dst, options, src)
}

// Useful for debugging
func (c *Controller) Kamera() *kamera.Camera {
	return c.cam
}

func (c *Controller) Position() image.Rectangle {
	return c.target.Position()
}
