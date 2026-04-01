package gamestates

type OffsetToggler struct {
	current int
}

func NewOffsetToggler(offset int) *OffsetToggler {
	return &OffsetToggler{current: -offset}
}

func (o *OffsetToggler) Next() int {
	o.current = -o.current
	return o.current
}
