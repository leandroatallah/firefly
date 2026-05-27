package debug

import "sort"

// Group identifies the source of a registered debug entry.
type Group int

const (
	// GroupFlags is the group for CLI-registered flags (via Register).
	GroupFlags Group = iota
	// GroupChannels is the group for JSON-loaded debug channels (via Init).
	GroupChannels
)

// String returns a human-readable label for the group.
func (g Group) String() string {
	switch g {
	case GroupFlags:
		return "CLI Flags"
	case GroupChannels:
		return "Debug Channels"
	default:
		return "Other"
	}
}

// Entry represents a single debug flag in the registry.
type Entry struct {
	Name  string
	Ptr   *bool
	Group Group
}

// Register adds or replaces a named debug flag pointer in the registry.
// Last-write-wins for duplicate names. Entries registered through Register
// belong to GroupFlags.
func Register(name string, ptr *bool) {
	mu.Lock()
	if registry == nil {
		registry = make(map[string]*bool)
	}
	registry[name] = ptr
	mu.Unlock()
}

// List returns all registered debug entries: GroupFlags entries first
// (sorted alphabetically), then GroupChannels entries (sorted alphabetically).
// Returns a non-nil empty slice when nothing is registered.
func List() []Entry {
	mu.Lock()
	flagNames := make([]string, 0, len(registry))
	for k := range registry {
		flagNames = append(flagNames, k)
	}
	sort.Strings(flagNames)

	var chanMap map[string]*bool
	if m := channels.Load(); m != nil {
		chanMap = *m
	}
	chanNames := make([]string, 0, len(chanMap))
	for k := range chanMap {
		chanNames = append(chanNames, k)
	}
	sort.Strings(chanNames)

	out := make([]Entry, 0, len(flagNames)+len(chanNames))
	for _, n := range flagNames {
		out = append(out, Entry{Name: n, Ptr: registry[n], Group: GroupFlags})
	}
	for _, n := range chanNames {
		out = append(out, Entry{Name: n, Ptr: chanMap[n], Group: GroupChannels})
	}
	mu.Unlock()
	return out
}
