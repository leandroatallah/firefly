package physics

// It holds a list of all physics bodies to help with bodies interaction.
type Space struct {
	bodies []Body
}

func NewSpace() *Space {
	return &Space{bodies: []Body{}}
}

// Add registers a new body with the space.
func (s *Space) Add(body Body) {
	s.bodies = append(s.bodies, body)
}

// Bodies returns the list of all bodies in the space.
func (s *Space) Bodies() []Body {
	return s.bodies
}
