package physics

import (
	"math"

	"github.com/leandroatallah/firefly/internal/engine/config"
)

const (
	gravityForce = 4
)

// increaseVelocity applies acceleration to the velocity for a single axis.
// Capping is handled in the Update loop to correctly manage the 2D vector's magnitude.
func increaseVelocity(velocity, acceleration int) int {
	velocity += acceleration
	return velocity
}

// reduceVelocity applies friction to the velocity for a single axis, slowing it down.
// It brings the velocity to zero if it's smaller than the friction value to prevent jitter.
func reduceVelocity(velocity int) int {
	friction := config.Unit / 4
	if velocity > friction {
		return velocity - friction
	}
	if velocity < -friction {
		return velocity + friction
	}
	return 0
}

// smoothDiagonalMovement converts raw input acceleration into a scaled and normalized vector.
// This ensures that the player's acceleration is consistent in all directions.
//
// Math:
//  1. The base acceleration from input (e.g., 2) is scaled up to a value that can
//     overcome friction.
//  2. If moving diagonally, the acceleration vector's magnitude would be `sqrt(ax² + ay²)`.
//     To ensure the magnitude is the same as for cardinal movement, we normalize it by
//     dividing each component by `sqrt(2)`.
func smoothDiagonalMovement(accX, accY int) (int, int) {
	// This factor determines the player's acceleration strength.
	// It should be large enough to overcome the friction in `reduceVelocity`.
	// Friction is `config.Unit / 4`. The base input acceleration is 2.
	// We'll use a factor of `config.Unit / 6` so that the final acceleration
	// (2 * config.Unit / 6 = config.Unit / 3) is greater than friction.
	accelerationFactor := float64(config.Unit / 6)

	fAccX := float64(accX) * accelerationFactor
	fAccY := float64(accY) * accelerationFactor

	isDiagonal := accX != 0 && accY != 0
	if isDiagonal {
		fAccX /= math.Sqrt2
		fAccY /= math.Sqrt2
	}

	return int(fAccX), int(fAccY)
}

// clampAxisVelocity ensures that the velocity on a single axis does not exceed a given limit.
// It handles both positive and negative velocities by comparing against the absolute limit.
func clampAxisVelocity(velocity, limit int) int {
	if limit <= 0 {
		return 0
	}
	switch {
	case velocity > limit:
		return limit
	case velocity < -limit:
		return -limit
	default:
		return velocity
	}
}

// applyAxisMovement moves the body by a given distance along a single axis (X or Y).
// It checks for collisions with other objects in the space.
// It returns true if the movement was blocked by a collision.
func applyAxisMovement(body *PhysicsBody, distance int, isXAxis bool, space *Space) bool {
	if distance == 0 {
		return false
	}

	rect, ok := body.Shape.(*Rect)
	var before int
	if ok {
		if isXAxis {
			before = rect.x16
		} else {
			before = rect.y16
		}
	}

	body.ApplyValidMovement(distance, isXAxis, space)

	if !ok {
		return false
	}

	if isXAxis {
		return rect.x16 == before
	}
	return rect.y16 == before
}

// applyGravity simulates gravity by increasing the vertical velocity of the body
// downwards, up to a terminal velocity specified by `ground`.
func applyGravity(body *PhysicsBody, ground int) {
	if body.vy16 < ground {
		body.vy16 += gravityForce
	}
}

// clampToPlayArea ensures the body stays within the screen boundaries.
// It adjusts the body's position if it goes beyond the edges of the screen.
// It returns true if the body is touching or has gone past the bottom of the screen,
// which can be interpreted as being on the ground for platformer.
func clampToPlayArea(body *PhysicsBody, space *Space) bool {
	rect, ok := body.Shape.(*Rect)
	if !ok {
		return false
	}

	if rect.x16 < 0 {
		body.ApplyValidMovement(-rect.x16, true, nil)
	}

	rightEdge := rect.x16 + rect.width*config.Unit
	maxRight := config.ScreenWidth * config.Unit
	provider := space.GetTilemapDimensionsProvider()
	if provider != nil {
		maxRight = provider.GetTilemapWidth() * config.Unit
	}
	if rightEdge > maxRight {
		body.ApplyValidMovement(maxRight-rightEdge, true, nil)
	}

	// Vertical clamping
	minTop := 0
	maxBottom := config.ScreenHeight * config.Unit
	if provider != nil {
		minTop = (config.ScreenHeight - provider.GetTilemapHeight()) * config.Unit
		maxBottom = provider.GetTilemapHeight() * config.Unit
	}

	if rect.y16 < minTop {
		body.ApplyValidMovement(minTop-rect.y16, false, nil)
	}

	bottom := rect.y16 + rect.height*config.Unit
	if bottom >= maxBottom {
		if bottom > maxBottom {
			body.ApplyValidMovement(maxBottom-bottom, false, nil)
		}
		return true
	}

	return false
}
