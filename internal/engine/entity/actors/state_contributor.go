package actors

// StateContributor is an optional interface a skill adapter may implement
// to request an actor state change each frame. The Character polls all
// registered contributors in handleState before falling through to its
// default movement-based transitions.
type StateContributor interface {
	ContributeState(current ActorStateEnum) (ActorStateEnum, bool)
}
