package physics

type CollisionArea struct {
	Shape
}

func NewCollisionArea(element Shape) *CollisionArea {
	return &CollisionArea{
		Shape: element,
	}
}

func rectToCollisionArea(element Shape) *CollisionArea {
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
