package actors

import (
	"fmt"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/systems/input"
	"github.com/leandroatallah/firefly/internal/systems/physics"
)

type MovementState interface {
	State() MovementStateEnum
	OnStart()
	Move()
	Target() physics.Body
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
	state  MovementStateEnum
	actor  ActorEntity
	target physics.Body
}

func NewBaseMovementState(
	state MovementStateEnum,
	actor ActorEntity,
	target physics.Body,
) *BaseMovementState {
	return &BaseMovementState{state: state, actor: actor, target: target}
}

func (s *BaseMovementState) State() MovementStateEnum {
	return s.state
}

func (s *BaseMovementState) OnStart() {}

func (s *BaseMovementState) Target() physics.Body {
	return s.target
}

// Movement States
type InputMovementState struct {
	BaseMovementState
}

func (s *InputMovementState) Move() {
	if input.IsSomeKeyPressed(ebiten.KeyA, ebiten.KeyLeft) {
		s.actor.OnMoveLeft(s.actor.Speed())
	}
	if input.IsSomeKeyPressed(ebiten.KeyD, ebiten.KeyRight) {
		s.actor.OnMoveRight(s.actor.Speed())
	}
	if input.IsSomeKeyPressed(ebiten.KeyW, ebiten.KeyUp) {
		s.actor.OnMoveUp(s.actor.Speed())
	}
	if input.IsSomeKeyPressed(ebiten.KeyS, ebiten.KeyDown) {
		s.actor.OnMoveDown(s.actor.Speed())
	}
}

type RandMovementState struct {
	BaseMovementState
}

func (s *RandMovementState) Move() {
	a := []func(){
		func() { s.actor.OnMoveLeft(s.actor.Speed()) },
		func() { s.actor.OnMoveRight(s.actor.Speed()) },
		func() { s.actor.OnMoveUp(s.actor.Speed()) },
		func() { s.actor.OnMoveDown(s.actor.Speed()) },
	}

	i := rand.IntN(4)
	a[i]()
}

type ChaseMovementState struct {
	BaseMovementState
}

func (s *ChaseMovementState) Move() {}

type DumbChaseMovementState struct {
	BaseMovementState
}

func (s *DumbChaseMovementState) Move() {
	directions := calculateMovementDirections(s.actor, s.target, false)
	executeMovement(s.actor, directions)
}

type PatrolMovementState struct {
	BaseMovementState
	waypoints            []Point
	currentWaypointIndex int
	reachedThreshold     float64
	waitTime             int
	waitCounter          int
}

type Point struct {
	X, Y int
}

func (s *PatrolMovementState) Move() {}

// reachedWaypoint checks if the actor is close enough to the target waypoint
func (s *PatrolMovementState) reachedWaypoint(target Point) bool {
	p0x, p0y, p1x, p1y := s.actor.Position()
	actorCenterX := (p0x + p1x) / 2
	actorCenterY := (p0y + p1y) / 2

	// Calculate distance to waypoint
	dx := actorCenterX - target.X
	dy := actorCenterY - target.Y
	distance := float64(dx*dx + dy*dy)

	return distance <= s.reachedThreshold*s.reachedThreshold
}

// advanceToNextWaypoint moves to the next waypoint in the patrol route
func (s *PatrolMovementState) advanceToNextWaypoint() {
	s.currentWaypointIndex = (s.currentWaypointIndex + 1) % len(s.waypoints)
}

// SetWaypoints sets the patrol route waypoints
func (s *PatrolMovementState) SetWaypoints(waypoints []Point) {
	s.waypoints = waypoints
	s.currentWaypointIndex = 0
	s.waitCounter = 0
}

// SetPatrolConfig sets patrol behavior configuration
func (s *PatrolMovementState) SetPatrolConfig(reachedThreshold float64, waitTime int) {
	s.reachedThreshold = reachedThreshold
	s.waitTime = waitTime
}

type AvoidMovementState struct {
	BaseMovementState
}

func (s *AvoidMovementState) Move() {
	directions := calculateMovementDirections(s.actor, s.target, true)
	executeMovement(s.actor, directions)
}

// Movement utility functions
type MovementDirections struct {
	Up    bool
	Down  bool
	Left  bool
	Right bool
}

// calculateMovementDirections determines which directions to move based on actor and target positions
func calculateMovementDirections(actorPos, targetPos physics.Body, isAvoid bool) MovementDirections {
	p0x, p0y, p1x, p1y := actorPos.Position()
	e0x, e0y, e1x, e1y := targetPos.Position()
	var up, down, left, right bool

	// Check direction to chase destination
	if p1x < e0x {
		right = true
	} else if p0x > e1x {
		left = true
	}

	if p1y < e0y {
		down = true
	} else if p0y > e1y {
		up = true
	}

	if isAvoid {
		// Invert to  move away from target
		up, down, left, right = !up, !down, !left, !right
	}

	return MovementDirections{Up: up, Down: down, Left: left, Right: right}
}

func executeMovement(actor ActorEntity, directions MovementDirections) {
	if !directions.Up && !directions.Down && !directions.Left && !directions.Right {
		return
	}

	speed := actor.Speed()

	if directions.Up {
		if directions.Left {
			actor.OnMoveUpLeft(speed)
		} else if directions.Right {
			actor.OnMoveUpRight(speed)
		} else {
			actor.OnMoveUp(speed)
		}
	} else if directions.Down {
		if directions.Left {
			actor.OnMoveDownLeft(speed)
		} else if directions.Right {
			actor.OnMoveDownRight(speed)
		} else {
			actor.OnMoveDown(speed)
		}
	} else if directions.Left {
		actor.OnMoveLeft(speed)
	} else if directions.Right {
		actor.OnMoveRight(speed)
	}
}

// State factory method
// TODO: Should it be a method?
func NewMovementState(actor ActorEntity, state MovementStateEnum, target physics.Body) (MovementState, error) {
	b := NewBaseMovementState(state, actor, target)

	switch state {
	case Input:
		return &InputMovementState{BaseMovementState: *b}, nil
	case Rand:
		return &RandMovementState{BaseMovementState: *b}, nil
	case Chase:
		return &ChaseMovementState{BaseMovementState: *b}, nil
	case DumbChase:
		return &DumbChaseMovementState{BaseMovementState: *b}, nil
	case Avoid:
		return &AvoidMovementState{BaseMovementState: *b}, nil
	// case Patrol:
	// 	// Create patrol state with default values
	// 	patrolState := &PatrolMovementState{
	// 		BaseMovementState:    *b,
	// 		waypoints:            []Point{},
	// 		currentWaypointIndex: 0,
	// 		reachedThreshold:     5.0, // 5 units threshold
	// 		waitTime:             30,  // 30 frames wait time
	// 		waitCounter:          0,
	// 	}
	// 	return patrolState, nil
	default:
		return nil, fmt.Errorf("unknown movement state type")
	}
}
