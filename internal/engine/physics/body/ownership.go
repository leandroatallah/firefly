package body

// Ownership handles the hierarchical ownership of components.
// It allows components to traverse up to the root entity.
type Ownership struct {
	owner interface{}
}

// SetOwner sets the immediate owner of the component.
func (o *Ownership) SetOwner(owner interface{}) {
	o.owner = owner
}

// Owner returns the immediate owner of the component.
func (o *Ownership) Owner() interface{} {
	return o.owner
}

// LastOwner returns the top-most owner in the hierarchy.
// If the component has no owner, it returns nil.
func (o *Ownership) LastOwner() interface{} {
	current := o.owner
	visited := make(map[interface{}]struct{})

	for current != nil {
		if _, exists := visited[current]; exists {
			return current
		}
		visited[current] = struct{}{}

		if ownerGetter, ok := current.(interface{ Owner() interface{} }); ok {
			next := ownerGetter.Owner()
			if next == nil {
				return current
			}
			current = next
			continue
		}

		return current
	}
	return nil
}
