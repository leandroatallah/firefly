package kitskills

// OffsetToggler alternates between a positive and negative offset value.
type OffsetToggler struct {
	current int
}

// NewOffsetToggler creates an OffsetToggler starting at -offset.
func NewOffsetToggler(offset int) *OffsetToggler {
	return &OffsetToggler{current: -offset}
}

// Next toggles the sign and returns the new value.
func (o *OffsetToggler) Next() int {
	o.current = -o.current
	return o.current
}
