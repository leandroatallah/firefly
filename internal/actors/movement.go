package actors

type MovementState interface {
	State() MovementStateEnum
	OnStart()
}

type MovementStateEnum int

const (
	Input MovementStateEnum = iota
	Rand
	Chase
	DumbChase
	Patrol
	Avoid
)

type BaseMovementState struct {
	state MovementStateEnum
	actor ActorEntity
}

func (s *BaseMovementState) State() MovementStateEnum {
	return s.state
}

func (s *BaseMovementState) OnStart() {}

// Movement States
type RandMovementState struct {
	BaseMovementState
}

type ChaseMovementState struct {
	BaseMovementState
}

type DumbChaseMovementState struct {
	BaseMovementState
}

type PatrolMovementState struct {
	BaseMovementState
}

type AvoidMovementState struct {
	BaseMovementState
}
