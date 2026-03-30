package camera

import (
	"image"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	enginecamera "github.com/boilerplate/ebiten-template/internal/engine/render/camera"
	"github.com/hajimehoshi/ebiten/v2"
)

// Controller wraps the engine's camera.Controller to provide a game-layer abstraction.
type Controller struct {
	base *enginecamera.Controller
}

// Base returns the underlying engine camera controller.
func (c *Controller) Base() *enginecamera.Controller {
	return c.base
}

// NewController creates a new game-layer camera controller wrapping the engine controller.
func NewController(base *enginecamera.Controller) *Controller {
	return &Controller{
		base: base,
	}
}

// Update delegates to the base camera controller.
func (c *Controller) Update() {
	c.base.Update()
}

// Draw delegates to the base camera controller.
func (c *Controller) Draw(
	src *ebiten.Image, options *ebiten.DrawImageOptions, dst *ebiten.Image,
) {
	c.base.Draw(src, options, dst)
}

// DrawCollisionBox delegates to the base camera controller.
func (c *Controller) DrawCollisionBox(screen *ebiten.Image, b body.Collidable) {
	c.base.DrawCollisionBox(screen, b)
}

// SetFollowTarget sets the follow target.
func (c *Controller) SetFollowTarget(b body.Body) {
	c.base.SetFollowTarget(b)
}

// SetBounds delegates to the base camera controller.
func (c *Controller) SetBounds(bounds *image.Rectangle) {
	c.base.SetBounds(bounds)
}

// Bounds delegates to the base camera controller.
func (c *Controller) Bounds() *image.Rectangle {
	return c.base.Bounds()
}

// SetCenter delegates to the base camera controller.
func (c *Controller) SetCenter(x, y float64) {
	c.base.SetCenter(x, y)
}

// SetPositionTopLeft delegates to the base camera controller.
func (c *Controller) SetPositionTopLeft(x, y float64) {
	c.base.SetPositionTopLeft(x, y)
}

// IsFollowing delegates to the base camera controller.
func (c *Controller) IsFollowing() bool {
	return c.base.IsFollowing()
}

// FollowTarget delegates to the base camera controller.
func (c *Controller) FollowTarget() body.Body {
	return c.base.FollowTarget()
}

// SetFollowing delegates to the base camera controller.
func (c *Controller) SetFollowing(following bool) {
	c.base.SetFollowing(following)
}

// Position delegates to the base camera controller.
func (c *Controller) Position() image.Rectangle {
	return c.base.Position()
}

// Target delegates to the base camera controller.
func (c *Controller) Target() body.Body {
	return c.base.Target()
}

// Width delegates to the base camera controller.
func (c *Controller) Width() float64 {
	return c.base.Width()
}

// Height delegates to the base camera controller.
func (c *Controller) Height() float64 {
	return c.base.Height()
}

// AddTrauma delegates to the base camera controller.
func (c *Controller) AddTrauma(amount float64) {
	c.base.AddTrauma(amount)
}

// CamDebug delegates to the base camera controller for debug purposes.
func (c *Controller) CamDebug() {
	c.base.CamDebug()
}

// Kamera returns nil; the game camera layer does not expose the underlying kamera instance.
func (c *Controller) Kamera() interface{} {
	return nil
}

// DisableSmoothing delegates to the base camera controller.
func (c *Controller) DisableSmoothing() {
	c.base.DisableSmoothing()
}
