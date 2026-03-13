package movement

import (
	"image"
	"math/rand"

	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
)

var Wander MovementStateEnum

func init() {
	Wander = RegisterMovementState("Wander", NewWanderMovementState)
}

type wanderStateEnum int

const (
	wanderIdle wanderStateEnum = iota
	wanderMove
)

// WanderMovementState defines a movement behavior where an actor wanders around an anchor point.
// It alternates between moving and idling.
type WanderMovementState struct {
	BaseMovementState
	anchorX     int
	maxDistance int
	state       wanderStateEnum
	timer       int
	moveTime    int
	idleTime    int
	movingRight bool
}

// NewWanderMovementState creates a new WanderMovementState.
func NewWanderMovementState(base BaseMovementState) MovementState {
	return &WanderMovementState{
		BaseMovementState: base,
		maxDistance:       50, // Default max distance
	}
}

func (s *WanderMovementState) OnStart() {
	s.anchorX = s.Actor().Position().Min.X
	s.startIdle()
	s.BaseMovementState.OnStart()
}

func (s *WanderMovementState) Move(space body.BodiesSpace) {
	if s.Actor().Immobile() {
		return
	}

	s.timer++

	switch s.state {
	case wanderIdle:
		if s.timer > s.idleTime {
			s.pickNextMove()
		}
	case wanderMove:
		if s.shouldStop(space) {
			// Hit a wall or ledge, stop immediately
			s.startIdle()
			return
		}

		if s.movingRight {
			s.Actor().OnMoveRight(s.Actor().Speed())
		} else {
			s.Actor().OnMoveLeft(s.Actor().Speed())
		}

		if s.timer > s.moveTime {
			s.startIdle()
		}
	}
}

func (s *WanderMovementState) startIdle() {
	s.state = wanderIdle
	s.timer = 0
	s.idleTime = 60 + rand.Intn(120) // Random idle 1-3s (assuming 60 FPS)
}

func (s *WanderMovementState) pickNextMove() {
	s.state = wanderMove
	s.timer = 0
	s.moveTime = 30 + rand.Intn(60) // Random move 0.5-1.5s

	currentX := s.Actor().Position().Min.X

	// Decide direction based on anchor
	if currentX > s.anchorX+s.maxDistance {
		s.movingRight = false
	} else if currentX < s.anchorX-s.maxDistance {
		s.movingRight = true
	} else {
		s.movingRight = rand.Intn(2) == 0
	}
}

// shouldStop checks for conditions that should make the actor stop (ledge or wall).
func (s *WanderMovementState) shouldStop(space body.BodiesSpace) bool {
	if space == nil {
		return false
	}
	actorPos := s.Actor().Position()

	// 1. Ledge detection
	var groundCheckPoint image.Point
	if s.movingRight {
		// Check point is at the actor's bottom-right corner, plus one pixel down.
		groundCheckPoint = image.Point{X: actorPos.Max.X, Y: actorPos.Max.Y + 1}
	} else {
		// Check point is at the actor's bottom-left corner, minus one pixel left, plus one pixel down.
		groundCheckPoint = image.Point{X: actorPos.Min.X - 1, Y: actorPos.Max.Y + 1}
	}

	groundCheckRect := image.Rectangle{Min: groundCheckPoint, Max: groundCheckPoint.Add(image.Point{X: 1, Y: 1})}

	hasGround := false
	colliders := space.Query(groundCheckRect)
	for _, c := range colliders {
		if c.IsObstructive() && c.ID() != s.Actor().ID() {
			hasGround = true
			break
		}
	}

	if !hasGround {
		return true // Stop at ledge
	}

	// 2. Wall detection
	var wallCheckRect image.Rectangle
	if s.movingRight {
		// Check a 1-pixel-wide vertical slice right in front of the actor.
		wallCheckRect = image.Rect(actorPos.Max.X, actorPos.Min.Y, actorPos.Max.X+1, actorPos.Max.Y)
	} else {
		// Check a 1-pixel-wide vertical slice right in front of the actor.
		wallCheckRect = image.Rect(actorPos.Min.X-1, actorPos.Min.Y, actorPos.Min.X, actorPos.Max.Y)
	}

	colliders = space.Query(wallCheckRect)
	for _, c := range colliders {
		if c.IsObstructive() && c.ID() != s.Actor().ID() {
			return true // Stop at wall
		}
	}

	return false
}
