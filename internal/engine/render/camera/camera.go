package camera

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
	"github.com/setanarut/kamera/v2"
)

var collisionBoxImage *ebiten.Image

func init() {
	collisionBoxImage = ebiten.NewImage(1, 1)
	collisionBoxImage.Fill(color.White)
}

type Controller struct {
	cam              *kamera.Camera
	target           body.Collidable
	followTarget     body.Body
	DeadZoneRadius   float64
	SmoothingFactor  float64
	isFollowing      bool
	centerX, centerY float64
	screenWidth      float64
	screenHeight     float64
	bounds           *image.Rectangle
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
		cam:          cam,
		target:       targetBody,
		isFollowing:  false,
		centerX:      x,
		centerY:      y,
		screenWidth:  float64(cfg.ScreenWidth),
		screenHeight: float64(cfg.ScreenHeight),
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

func (c *Controller) IsFollowing() bool {
	return c.isFollowing
}

func (c *Controller) FollowTarget() body.Body {
	return c.followTarget
}

func (c *Controller) SetFollowing(following bool) {
	c.isFollowing = following
}

// SetBounds restricts the camera movement to the specified rectangle.
func (c *Controller) SetBounds(bounds *image.Rectangle) {
	c.bounds = bounds
}

// Bounds returns the current camera movement restrictions.
func (c *Controller) Bounds() *image.Rectangle {
	return c.bounds
}

func (c *Controller) SetCenter(x, y float64) {
	c.centerX = x
	c.centerY = y
	c.Kamera().SetCenter(x, y)
}

// GetCenter returns the current camera center position.
func (c *Controller) GetCenter() (float64, float64) {
	return c.centerX, c.centerY
}

// GetActualCenter returns the actual smoothed camera center position from kamera.
func (c *Controller) GetActualCenter() (float64, float64) {
	return c.cam.Center()
}

// SetPositionTopLeft positions the camera so that (x, y) appears at the top-left corner of the screen.
func (c *Controller) SetPositionTopLeft(x, y float64) {
	centerX := x + c.screenWidth/2
	centerY := y + c.screenHeight/2
	c.centerX = centerX
	c.centerY = centerY
	c.Kamera().SetCenter(centerX, centerY)
}

func (c *Controller) SetFollowTarget(b body.Body) {
	c.followTarget = b
	x, y := b.GetPositionMin()
	w, h := b.GetShape().Width(), b.GetShape().Height()
	targetX := float64(x) + float64(w)/2
	targetY := float64(y) + float64(h)/2
	c.Kamera().SetCenter(targetX, targetY)
}

func (c *Controller) Update() {
	var targetX, targetY float64
	if c.isFollowing && c.followTarget != nil {
		x, y := c.followTarget.GetPositionMin()
		// Use center of target for following
		w, h := c.followTarget.GetShape().Width(), c.followTarget.GetShape().Height()
		targetX = float64(x) + float64(w)/2
		targetY = float64(y) + float64(h)/2
	} else {
		targetX = c.centerX
		targetY = c.centerY
	}

	if c.bounds != nil {
		// Calculate viewport half-dimensions
		halfW := c.screenWidth / 2
		halfH := c.screenHeight / 2

		// Calculate min and max center positions
		minX := float64(c.bounds.Min.X) + halfW
		maxX := float64(c.bounds.Max.X) - halfW
		minY := float64(c.bounds.Min.Y) + halfH
		maxY := float64(c.bounds.Max.Y) - halfH

		// Clamp targetX
		if targetX < minX {
			targetX = minX
		}
		if targetX > maxX {
			targetX = maxX
		}

		// Clamp targetY
		if targetY < minY {
			targetY = minY
		}
		if targetY > maxY {
			targetY = maxY
		}
	}

	c.cam.LookAt(targetX, targetY)
	// Update stored center to match actual camera position
	c.centerX = targetX
	c.centerY = targetY
}

func (c *Controller) Draw(
	src *ebiten.Image, options *ebiten.DrawImageOptions, dst *ebiten.Image,
) {
	c.cam.Draw(src, options, dst)
}

func (c *Controller) DrawCollisionBox(screen *ebiten.Image, b body.Collidable) {
	isObstructive := b.IsObstructive()
	for _, rect := range b.CollisionPosition() {
		// Draw outer rect
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(float64(rect.Dx()), float64(rect.Dy()))
		opts.GeoM.Translate(float64(rect.Min.X), float64(rect.Min.Y))

		if isObstructive {
			opts.ColorScale.Scale(0.66, 0, 0, 1) // Dark Red
		} else {
			opts.ColorScale.Scale(0, 0.66, 0, 1) // Dark Green
		}
		c.Draw(collisionBoxImage, opts, screen)

		// Draw inner rect
		if rect.Dx() > 2 && rect.Dy() > 2 {
			opts = &ebiten.DrawImageOptions{}
			opts.GeoM.Scale(float64(rect.Dx()-2), float64(rect.Dy()-2))
			opts.GeoM.Translate(float64(rect.Min.X+1), float64(rect.Min.Y+1))

			if isObstructive {
				opts.ColorScale.Scale(1, 0, 0, 1) // Red
			} else {
				opts.ColorScale.Scale(0, 1, 0, 1) // Green
			}
			c.Draw(collisionBoxImage, opts, screen)
		}
	}
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

func (c *Controller) Width() float64 {
	return c.screenWidth
}

func (c *Controller) Height() float64 {
	return c.screenHeight
}

func (c *Controller) AddTrauma(amount float64) {
	c.cam.AddTrauma(amount)
}
