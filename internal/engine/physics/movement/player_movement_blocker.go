package physicsmovement

// PlayerMovementBlocker defines the interface for checking if player movement is blocked.
type PlayerMovementBlocker interface {
	IsMovementBlocked() bool
}
