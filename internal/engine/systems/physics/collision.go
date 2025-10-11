package physics

import "github.com/leandroatallah/firefly/internal/engine/contracts/body"

type CollisionArea struct {
	body.Shape
}

func NewCollisionArea(element body.Shape) *CollisionArea {
	return &CollisionArea{
		Shape: element,
	}
}

func rectToCollisionArea(element body.Shape) *CollisionArea {
	rect := element.Position()
	width := rect.Max.X - rect.Min.X
	height := rect.Max.Y - rect.Min.Y

	return NewCollisionArea(
		NewRect(
			rect.Min.X,
			rect.Min.Y,
			width,
			height,
		),
	)
}
