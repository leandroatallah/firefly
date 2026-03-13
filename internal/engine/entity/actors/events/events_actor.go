package events

const (
	ActorJumpedType = "actor_jumped"
	ActorLandedType = "actor_landed"
)

type ActorJumpedEvent struct {
	X, Y float64
}

func (e *ActorJumpedEvent) Type() string {
	return ActorJumpedType
}

type ActorLandedEvent struct {
	X, Y float64
}

func (e *ActorLandedEvent) Type() string {
	return ActorLandedType
}
