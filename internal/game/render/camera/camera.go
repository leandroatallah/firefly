package camera

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	enginecamera "github.com/leandroatallah/firefly/internal/engine/render/camera"
)

// Controller wraps the engine's camera.Controller to add game-specific behavior.
// Currently implements vertical-only-upward constraint: camera moves up but never down.
type Controller struct {
	base        *enginecamera.Controller
	lastCameraY float64
	initialized bool
}

// Base returns the underlying engine camera controller.
func (c *Controller) Base() *enginecamera.Controller {
	return c.base
}

// NewController creates a new game-layer camera controller wrapping the engine controller.
func NewController(base *enginecamera.Controller) *Controller {
	return &Controller{
		base:        base,
		lastCameraY: 0,
	}
}

// Update updates the camera position with game-specific constraints applied.
// Calculates target position and applies vertical-only-upward constraint BEFORE setting camera position.
func (c *Controller) Update() {
	if c.base.IsFollowing() && c.base.FollowTarget() != nil {
		target := c.base.FollowTarget()
		x, y := target.GetPositionMin()
		w, h := target.GetShape().Width(), target.GetShape().Height()
		targetX := float64(x) + float64(w)/2
		targetY := float64(y) + float64(h)/2

		// Apply vertical-only-upward constraint BEFORE setting position
		if c.initialized && targetY > c.lastCameraY {
			targetY = c.lastCameraY // Block downward movement
		} else {
			c.lastCameraY = targetY
		}
		c.initialized = true

		// Apply bounds clamping manually
		if bounds := c.base.Bounds(); bounds != nil {
			halfW := c.base.Width() / 2
			halfH := c.base.Height() / 2
			minX := float64(bounds.Min.X) + halfW
			maxX := float64(bounds.Max.X) - halfW
			minY := float64(bounds.Min.Y) + halfH
			maxY := float64(bounds.Max.Y) - halfH

			if targetX < minX {
				targetX = minX
			}
			if targetX > maxX {
				targetX = maxX
			}
			if targetY < minY {
				targetY = minY
			}
			if targetY > maxY {
				targetY = maxY
			}
		}

		c.base.SetCenter(targetX, targetY)
	}
	// Don't call base.Update() - we handled everything to prevent downward target calculation
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

// SetFollowTarget sets the follow target and initializes lastCameraY.
func (c *Controller) SetFollowTarget(b body.Body) {
	c.base.SetFollowTarget(b)
	_, y := b.GetPositionMin()
	h := b.GetShape().Height()
	c.lastCameraY = float64(y) + float64(h)/2
	c.initialized = true
}

// SetLastCameraY sets the last camera Y position for vertical constraint tracking.
func (c *Controller) SetLastCameraY(y float64) {
	c.lastCameraY = y
	c.initialized = true
}

// SetBounds delegates to the base camera controller.
func (c *Controller) SetBounds(bounds *image.Rectangle) {
	c.base.SetBounds(bounds)
}

// Bounds delegates to the base camera controller.
func (c *Controller) Bounds() *image.Rectangle {
	return c.base.Bounds()
}

// SetCenter delegates to the base camera controller and updates lastCameraY.
func (c *Controller) SetCenter(x, y float64) {
	c.lastCameraY = y
	c.initialized = true
	c.base.SetCenter(x, y)
}

// SetPositionTopLeft delegates to the base camera controller and updates lastCameraY.
func (c *Controller) SetPositionTopLeft(x, y float64) {
	c.lastCameraY = y + c.base.Height()/2
	c.initialized = true
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

// Kamera returns the underlying kamera.Camera for debug purposes.
func (c *Controller) Kamera() *ebiten.Image {
	return nil // Placeholder - should not be used directly
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
