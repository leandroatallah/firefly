package physics

type CollisionArea struct {
	Rect
}

func NewCollisionArea(element Rect) *CollisionArea {
	return &CollisionArea{
		Rect: element,
	}
}

func elementToCollisionArea(element Rect) *CollisionArea {
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
