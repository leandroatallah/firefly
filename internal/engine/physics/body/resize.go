package body

import "image"

func ResizeFixedBottom(rect image.Rectangle, newHeight int) image.Rectangle {
	if newHeight < 0 {
		newHeight = 0
	}
	return image.Rectangle{
		Min: image.Point{X: rect.Min.X, Y: rect.Max.Y - newHeight},
		Max: rect.Max,
	}
}

func ResizeFixedTop(rect image.Rectangle, newHeight int) image.Rectangle {
	if newHeight < 0 {
		newHeight = 0
	}
	return image.Rectangle{
		Min: rect.Min,
		Max: image.Point{X: rect.Max.X, Y: rect.Min.Y + newHeight},
	}
}
