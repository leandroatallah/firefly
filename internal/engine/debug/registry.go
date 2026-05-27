package debug

import "sort"

// Entry represents a single debug flag in the registry.
type Entry struct {
	Name string
	Ptr  *bool
}

// Register adds or replaces a named debug flag pointer in the registry.
// Last-write-wins for duplicate names.
func Register(name string, ptr *bool) {
	mu.Lock()
	if registry == nil {
		registry = make(map[string]*bool)
	}
	registry[name] = ptr
	mu.Unlock()
}

// List returns all registered debug flags sorted alphabetically by name.
// Returns a non-nil empty slice when the registry is empty.
func List() []Entry {
	mu.Lock()
	names := make([]string, 0, len(registry))
	for k := range registry {
		names = append(names, k)
	}
	sort.Strings(names)
	out := make([]Entry, 0, len(names))
	for _, n := range names {
		out = append(out, Entry{n, registry[n]})
	}
	mu.Unlock()
	return out
}
