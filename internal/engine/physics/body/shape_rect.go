package body

type Rect struct {
	width, height int
}

func NewRect(x, y, width, height int) *Rect {
	return &Rect{
		width:  width,
		height: height,
	}
}

func (r *Rect) Width() int {
	return r.width
}

func (r *Rect) Height() int {
	return r.height
}
