package physics

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/leandroatallah/firefly/internal/config"
)

// Movable is a Shape but with movement
type Movable interface {
	// TODO: Make all Position return the same data type
	Position() (minX, minY, maxX, maxY int)
	ApplyValidMovement(velocity int, isXAxis bool, boundaries []Body)
}

// Body is a Shape but with collision
type Body interface {
	// TODO: Make all Position return the same data type
	Position() (minX, minY, maxX, maxY int)
	DrawCollisionBox(screen *ebiten.Image)
	CollisionPosition() []image.Rectangle
	IsColliding(boundaries []Body) bool
}

type PhysicsBody struct {
	Shape
	vx16          int
	vy16          int
	collisionList []*CollisionArea
}

func NewPhysicsBody(shape Shape) *PhysicsBody {
	return &PhysicsBody{Shape: shape}
}

func (b *PhysicsBody) Move() {
	panic("You should implement this method in derivated structs")
}

// TODO: Implement ease in movement
func (b *PhysicsBody) MoveY(distance int) {
	b.vy16 += distance * config.Unit
}

// TODO: Implement ease in movement
func (b *PhysicsBody) MoveX(distance int) {
	b.vx16 += distance * config.Unit
}

func (b *PhysicsBody) Position() (minX, minY, maxX, maxY int) {
	// TODO: Replace switch with "polymorphism"
	switch b.Shape.(type) {
	case *Rect:
		rect := b.Shape.(*Rect)
		minX = rect.x16 / config.Unit
		minY = rect.y16 / config.Unit
		maxX = minX + rect.width
		maxY = minY + rect.height
	}
	return
}

func (b *PhysicsBody) DrawCollisionBox(screen *ebiten.Image) {
	for _, c := range b.CollisionPosition() {
		minX := c.Min.X
		minY := c.Min.Y
		maxX := c.Max.X
		maxY := c.Max.Y

		width := float32(maxX - minX)
		height := float32(maxY - minY)
		vector.DrawFilledRect(
			screen,
			float32(minX), float32(minY), width, height,
			color.RGBA{0, 0xaa, 0, 0xff}, false)
		vector.DrawFilledRect(
			screen,
			float32(minX)+1, float32(minY)+1, width-2, height-2,
			color.RGBA{0, 0xff, 0, 0xff}, false)
	}
}

func (b *PhysicsBody) AddCollision(list ...*CollisionArea) *PhysicsBody {
	for _, i := range list {
		b.collisionList = append(b.collisionList, i)
	}
	return b
}

func (b *PhysicsBody) CollisionPosition() []image.Rectangle {
	res := []image.Rectangle{}
	for _, c := range b.collisionList {
		res = append(res, c.Position())
	}
	return res
}

// TODO: Needs to be updated when dealing with different shapes (e.g. circle)
func (b *PhysicsBody) checkRectIntersect(obj1, obj2 Body) bool {
	rects1 := obj1.CollisionPosition()
	rects2 := obj2.CollisionPosition()

	for _, r1 := range rects1 {
		for _, r2 := range rects2 {
			if r1.Overlaps(r2) {
				return true
			}
		}
	}

	return false
}

func (b *PhysicsBody) IsColliding(boundaries []Body) bool {
	for _, o := range boundaries {
		if b.checkRectIntersect(b, o.(Body)) {
			return true
		}
	}
	return false
}

func (b *PhysicsBody) updatePosition(velocity int, isXAxis bool) {
	// TODO: Replace switch with "polymorphism"
	switch b.Shape.(type) {
	case *Rect:
		rect := b.Shape.(*Rect)
		if isXAxis {
			rect.x16 += velocity
			for _, c := range b.collisionList {
				c.Shape.(*Rect).x16 += velocity
			}
		} else {
			rect.y16 += velocity
			for _, c := range b.collisionList {
				c.Shape.(*Rect).y16 += velocity
			}
		}
	}
}

func (b *PhysicsBody) ApplyValidMovement(velocity int, isXAxis bool, boundaries []Body) {
	if velocity == 0 {
		return
	}

	b.updatePosition(velocity, isXAxis)

	isValid := !b.IsColliding(boundaries)
	if !isValid {
		b.updatePosition(-velocity, isXAxis)
	}
}

func normalizeMoveOffset(move int, normalize bool) int {
	if normalize {
		return int(float64(move*config.Unit) / math.Sqrt2 / config.Unit)
	}
	return move
}
