package movement

import (
	"image"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
)

// FollowMovementState defines a movement behavior where an actor follows a target
// maintaining a certain distance or with a delay mechanism based on distance.
type FollowMovementState struct {
	BaseMovementState
	startDistance  int // Distance at which to start following
	stopDistance   int // Distance at which to stop following
	isMoving       bool
	stayOnPlatform bool // If true, won't move off edges
}

// NewFollowMovementState creates a new FollowMovementState.
func NewFollowMovementState(base BaseMovementState) *FollowMovementState {
	return &FollowMovementState{
		BaseMovementState: base,
		startDistance:     50, // Default to a noticeable delay/leash
		stopDistance:      20, // Stop reasonably close
	}
}

// WithFollowDistances sets the start and stop distances for following.
// start: Distance to target to trigger movement (the "leash" length).
// stop: Distance to target to stop movement.
func WithFollowDistances(start, stop int) MovementStateOption {
	return func(ms MovementState) {
		if s, ok := ms.(*FollowMovementState); ok {
			s.startDistance = start
			s.stopDistance = stop
		}
	}
}

// WithPlatformFollow enables platform-aware following that prevents falling off edges.
func WithPlatformFollow() MovementStateOption {
	return func(ms MovementState) {
		if s, ok := ms.(*FollowMovementState); ok {
			s.stayOnPlatform = true
		}
	}
}

func (s *FollowMovementState) Move(space body.BodiesSpace) {
	if s.actor.Immobile() {
		return
	}

	target := s.target
	if target == nil {
		return
	}

	dist := euclideanDistance(s.actor.Position().Min, target.Position().Min)

	if s.isMoving {
		if dist <= s.stopDistance {
			s.isMoving = false
		}
	} else {
		if dist >= s.startDistance {
			s.isMoving = true
		}
	}

	if s.isMoving {
		// Use direct movement calculation
		directions := calculateMovementDirections(s.actor, target, false)

		// If stayOnPlatform is enabled, filter out directions that would cause falling
		if s.stayOnPlatform {
			directions = s.filterSafeDirections(directions, space)
		}

		executeMovement(s.actor, directions)
	}
}

// filterSafeDirections removes directions that would cause the actor to fall off a platform.
func (s *FollowMovementState) filterSafeDirections(directions MovementDirections, space body.BodiesSpace) MovementDirections {
	safe := directions

	// If moving left, check if there's ground in that direction
	if directions.Left && !s.hasGroundInDirection(space, false) {
		safe.Left = false
	}

	// If moving right, check if there's ground in that direction
	if directions.Right && !s.hasGroundInDirection(space, true) {
		safe.Right = false
	}

	return safe
}

// hasGroundInDirection checks if there's solid ground in the specified direction.
// Uses the same approach as SideToSideMovementState for ledge detection.
func (s *FollowMovementState) hasGroundInDirection(space body.BodiesSpace, isRight bool) bool {
	actorPos := s.actor.Position()

	// Check point at bottom corner + 1 pixel down and 1 pixel in direction
	var groundCheckPoint image.Point
	if isRight {
		groundCheckPoint = image.Point{X: actorPos.Max.X, Y: actorPos.Max.Y + 1}
	} else {
		groundCheckPoint = image.Point{X: actorPos.Min.X - 1, Y: actorPos.Max.Y + 1}
	}

	groundCheckRect := image.Rectangle{Min: groundCheckPoint, Max: groundCheckPoint.Add(image.Point{X: 1, Y: 1})}

	// Query space for collidables at the check point
	colliders := space.Query(groundCheckRect)
	for _, c := range colliders {
		if c.IsObstructive() && c.ID() != s.actor.ID() {
			return true
		}
	}

	return false
}
