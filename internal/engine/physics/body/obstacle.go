package body

import (
	"fmt"
	"image"
	"image/color"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// defaultLaneHalfWidth mirrors space.DefaultLaneHalfWidth (8 px).
// Defined locally to avoid an import cycle: the space package already
// imports this (body) package via state_collision_manager.go.
const defaultLaneHalfWidth = 8

type ObstacleRect struct {
	Ownership
	*MovableBody
	*CollidableBody

	imageOptions *ebiten.DrawImageOptions
}

func NewObstacleRect(bodyRect *Rect) *ObstacleRect {
	b := NewBody(bodyRect)
	movable := NewMovableBody(b)
	collidable := NewCollidableBody(b)
	obs := &ObstacleRect{
		MovableBody:    movable,
		CollidableBody: collidable,

		imageOptions: &ebiten.DrawImageOptions{},
	}
	// Set the owner for all body components to this ObstacleRect
	b.SetOwner(movable)
	movable.SetOwner(obs)
	collidable.SetOwner(obs)

	return obs
}

// Forwarding methods for Body to avoid ambiguous selector
// Always route via the MovableBody component
func (o *ObstacleRect) ID() string {
	return o.MovableBody.ID()
}
func (o *ObstacleRect) SetID(id string) {
	o.MovableBody.SetID(id)
}
func (o *ObstacleRect) Position() image.Rectangle {
	return o.MovableBody.Position()
}
func (o *ObstacleRect) SetPosition(x, y int) {
	o.MovableBody.SetPosition(x, y)
}
func (o *ObstacleRect) GetPositionMin() (int, int) {
	return o.MovableBody.GetPositionMin()
}

func (o *ObstacleRect) SetPosition16(x16, y16 int) {
	o.MovableBody.SetPosition16(x16, y16)
}

func (o *ObstacleRect) SetSize(width, height int) {
	o.MovableBody.SetSize(width, height)
}

func (o *ObstacleRect) Scale() float64 {
	return o.MovableBody.Scale()
}

func (o *ObstacleRect) SetScale(scale float64) {
	o.MovableBody.SetScale(scale)
}

func (o *ObstacleRect) GetPosition16() (int, int) {
	return o.MovableBody.GetPosition16()
}

func (o *ObstacleRect) GetShape() body.Shape {
	return o.MovableBody.GetShape()
}

func (o *ObstacleRect) AddCollisionBodies(list ...body.Collidable) {
	if len(list) == 0 {
		b := NewCollidableBodyFromRect(o.GetShape())
		x, y := o.GetPositionMin()
		b.SetPosition(x, y)
		b.SetID(fmt.Sprintf("%v_COLLISION_0", o.ID()))
		list = []body.Collidable{b}
	}
	o.AddCollision(list...)
}

func (o *ObstacleRect) Draw(screen *ebiten.Image) {
	rect := o.GetShape().(*Rect)
	x, y := o.GetPositionMin()
	vector.DrawFilledRect(
		screen,
		float32(x),
		float32(y),
		float32(rect.width),
		float32(rect.height),
		color.Transparent,
		false,
	)
}

func (o *ObstacleRect) Image() *ebiten.Image {
	w := o.Position().Dx()
	h := o.Position().Dy()
	i := ebiten.NewImage(w, h)
	return i
}

func (o *ObstacleRect) ImageOptions() *ebiten.DrawImageOptions {
	return o.imageOptions
}

func (o *ObstacleRect) UpdateImageOptions() {
	o.imageOptions.GeoM.Reset()
	x, y := o.GetPositionMin()
	o.imageOptions.GeoM.Translate(float64(x), float64(y))
}

func (o *ObstacleRect) Altitude() int       { return o.MovableBody.Altitude() }
func (o *ObstacleRect) Altitude16() int     { return o.MovableBody.Altitude16() }
func (o *ObstacleRect) SetAltitude(a int)   { o.MovableBody.SetAltitude(a) }
func (o *ObstacleRect) SetAltitude16(a int) { o.MovableBody.SetAltitude16(a) }

// GroundY implements space.DepthLaneBody.
// Returns the pre-altitude bottom edge in world (depth) coordinates.
// Uses the raw y16 field (not Position().Max.Y) so that airborne bodies
// return their floor-projected bottom rather than their screen-offset bottom.
// For obstacles (altitude always 0), the result is identical to Position().Max.Y.
func (o *ObstacleRect) GroundY() int {
	_, y16 := o.GetPosition16()
	return y16/16 + o.GetShape().(*Rect).Height()
}

// LaneHalfWidth implements space.DepthLaneBody.
// Returns max(height, defaultLaneHalfWidth) so that a character whose feet line
// falls anywhere across the obstacle's full Y-extent triggers the gate.
// Zero-height obstacles fall back to defaultLaneHalfWidth to avoid an
// exact-equal-match requirement.
func (o *ObstacleRect) LaneHalfWidth() int {
	h := o.GetShape().(*Rect).Height()
	if h < defaultLaneHalfWidth {
		return defaultLaneHalfWidth
	}
	return h
}
