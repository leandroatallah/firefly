package input

type HorizontalAxis struct {
	held [2]bool // index 0 = left (-1), index 1 = right (+1)
	last int     // -1 or +1
}

func NewHorizontalAxis() *HorizontalAxis {
	return &HorizontalAxis{}
}

func (h *HorizontalAxis) Press(dir int) {
	h.held[dirIndex(dir)] = true
	h.last = dir
}

func (h *HorizontalAxis) Release(dir int) {
	h.held[dirIndex(dir)] = false
}

func (h *HorizontalAxis) Value() int {
	left, right := h.held[0], h.held[1]
	switch {
	case left && right:
		return h.last
	case left:
		return -1
	case right:
		return 1
	default:
		return 0
	}
}

func dirIndex(dir int) int {
	if dir < 0 {
		return 0
	}
	return 1
}
