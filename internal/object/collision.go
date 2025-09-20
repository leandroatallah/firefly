package object

type CollisionArea struct {
	Element
}

func NewCollisionArea(element Element) *CollisionArea {
	return &CollisionArea{
		Element: element,
	}
}

func elementToCollisionArea(element Element) *CollisionArea {
	minX, minY, maxX, maxY := element.Position()
	width := maxX - minX
	height := maxY - minY

	return NewCollisionArea(
		NewElement(
			minX,
			minY,
			width,
			height,
		),
	)
}
