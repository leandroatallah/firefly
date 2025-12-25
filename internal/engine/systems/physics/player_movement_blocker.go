package physics

// PlayerMovementBlocker defines the interface for checking if player movement is blocked.
// TODO: Review this. Is it necessary? Should the name be this?
type PlayerMovementBlocker interface {
	IsMovementBlocked() bool
}
