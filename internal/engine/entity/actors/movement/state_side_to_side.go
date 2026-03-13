package movement

import (
	"image"

	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	physicsmovement "github.com/leandroatallah/firefly/internal/engine/physics/movement"
	"github.com/leandroatallah/firefly/internal/engine/utils/fp16"
)

// SideToSideMovementState defines a movement behavior where an actor moves
// back and forth (horizontally or vertically), changing direction upon detecting a ledge or a wall.
type SideToSideMovementState struct {
	BaseMovementState
	movingPositive bool
	vertical       bool
	waitDuration   int
	waitTimer      int
	isWaiting      bool
	ignoreLedges   bool
	limitToRoom    bool
}

// NewSideToSideMovementState creates a new SideToSideMovementState.
func NewSideToSideMovementState(base BaseMovementState) *SideToSideMovementState {
	return &SideToSideMovementState{
		BaseMovementState: base,
		movingPositive:    true, // Start by moving right/down
	}
}

// WithWaitBeforeTurn sets a delay before the actor turns to the other side.
func WithWaitBeforeTurn(duration int) MovementStateOption {
	return func(ms MovementState) {
		if s, ok := ms.(*SideToSideMovementState); ok {
			s.waitDuration = duration
		}
	}
}

// WithVerticalMovement sets whether the actor should move vertically.
func WithVerticalMovement(vertical bool) MovementStateOption {
	return func(ms MovementState) {
		if s, ok := ms.(*SideToSideMovementState); ok {
			s.vertical = vertical
		}
	}
}

// WithIgnoreLedges sets whether the actor should ignore ledges (e.g. for flying enemies).
func WithIgnoreLedges(ignore bool) MovementStateOption {
	return func(ms MovementState) {
		if s, ok := ms.(*SideToSideMovementState); ok {
			s.ignoreLedges = ignore
		}
	}
}

// WithLimitToRoom sets whether the actor should limit movement to the current room (camera bounds).
func WithLimitToRoom(limit bool) MovementStateOption {
	return func(ms MovementState) {
		if s, ok := ms.(*SideToSideMovementState); ok {
			s.limitToRoom = limit
		}
	}
}

// Move executes the side-to-side movement logic. It checks for ledges and walls
// to reverse direction and then applies movement.
func (s *SideToSideMovementState) Move(space body.BodiesSpace) {
	if s.actor.Immobile() {
		return
	}

	setGravity := func(enabled bool) {
		if provider, ok := s.actor.(interface {
			MovementModel() physicsmovement.MovementModel
		}); ok {
			if pm, ok := provider.MovementModel().(*physicsmovement.PlatformMovementModel); ok {
				pm.SetGravityEnabled(enabled)
			}
		}
	}

	if s.isWaiting {
		s.waitTimer--
		if s.waitTimer <= 0 {
			s.isWaiting = false
			s.movingPositive = !s.movingPositive
		} else {
			if s.vertical {
				setGravity(false)
				vx, _ := s.actor.Velocity()
				s.actor.SetVelocity(vx, 0)
			}
		}
		return
	}

	if s.shouldTurn(space) {
		if s.waitDuration > 0 {
			s.isWaiting = true
			s.waitTimer = s.waitDuration
			return
		}
		s.movingPositive = !s.movingPositive
	}

	if s.movingPositive {
		if s.vertical {
			setGravity(true)
			speed := s.actor.Speed()
			if m := config.Get().Physics.SpeedMultiplier; m != 0 {
				speed = int(float64(speed) * m)
			}
			vx, _ := s.actor.Velocity()
			s.actor.SetVelocity(vx, fp16.To16(speed))
		} else {
			s.actor.OnMoveRight(s.actor.Speed())
		}
	} else {
		if s.vertical {
			setGravity(true)
			speed := s.actor.Speed()
			if m := config.Get().Physics.SpeedMultiplier; m != 0 {
				speed = int(float64(speed) * m)
			}
			s.actor.TryJump(speed)
		} else {
			s.actor.OnMoveLeft(s.actor.Speed())
		}
	}
}

// shouldTurn checks for conditions that should make the actor reverse direction.
// It returns true if a wall is directly in front of the actor or if there is no
// ground just ahead of it (a ledge).
func (s *SideToSideMovementState) shouldTurn(space body.BodiesSpace) bool {
	if space == nil {
		return false
	}
	actorPos := s.actor.Position()

	// 1. Ledge detection (Only for horizontal movement)
	if !s.vertical && !s.ignoreLedges {
		var groundCheckPoint image.Point
		if s.movingPositive {
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
			if c.IsObstructive() && c.ID() != s.actor.ID() {
				hasGround = true
				break
			}
		}

		if !hasGround {
			return true // Turn at ledge
		}
	}

	// 2. Wall detection
	var wallCheckRect image.Rectangle
	if s.vertical {
		if s.movingPositive {
			// Check a 1-pixel-wide horizontal slice right below the actor.
			wallCheckRect = image.Rect(actorPos.Min.X, actorPos.Max.Y, actorPos.Max.X, actorPos.Max.Y+1)
		} else {
			// Check a 1-pixel-wide horizontal slice right above the actor.
			wallCheckRect = image.Rect(actorPos.Min.X, actorPos.Min.Y-1, actorPos.Max.X, actorPos.Min.Y)
		}
	} else {
		if s.movingPositive {
			// Check a 1-pixel-wide vertical slice right in front of the actor.
			wallCheckRect = image.Rect(actorPos.Max.X, actorPos.Min.Y, actorPos.Max.X+1, actorPos.Max.Y)
		} else {
			// Check a 1-pixel-wide vertical slice right in front of the actor.
			wallCheckRect = image.Rect(actorPos.Min.X-1, actorPos.Min.Y, actorPos.Min.X, actorPos.Max.Y)
		}
	}

	colliders := space.Query(wallCheckRect)
	for _, c := range colliders {
		if c.IsObstructive() && c.ID() != s.actor.ID() {
			return true // Turn at wall
		}
	}

	// 3. Room limit detection
	if s.limitToRoom {
		if provider := space.GetTilemapDimensionsProvider(); provider != nil {
			if bounds, ok := provider.GetCameraBounds(); ok {
				if s.vertical {
					if s.movingPositive {
						if actorPos.Max.Y >= bounds.Max.Y {
							return true
						}
					} else {
						if actorPos.Min.Y <= bounds.Min.Y {
							return true
						}
					}
				} else {
					if s.movingPositive {
						if actorPos.Max.X >= bounds.Max.X {
							return true
						}
					} else {
						if actorPos.Min.X <= bounds.Min.X {
							return true
						}
					}
				}
			}
		}
	}

	return false
}
