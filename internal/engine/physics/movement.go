package physics

import (
	"math"

	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
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
	friction := config.Get().Unit / 4
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
	accelerationFactor := float64(config.Get().Unit / 6)

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

// clampToPlayArea ensures the body stays within the screen boundaries.
// It adjusts the body's position if it goes beyond the edges of the screen.
// It returns true if the body is touching or has gone past the bottom of the screen,
// which can be interpreted as being on the ground for platformer.
func clampToPlayArea(body body.MovableCollidable, space *Space) bool {
	rect, ok := body.GetShape().(*Rect)
	if !ok {
		return false
	}

	cfg := config.Get()
	x, y := body.GetPositionMin()

	if x < 0 {
		x, y, _ = body.ApplyValidPosition(-x*cfg.Unit, true, nil)
	}

	x16, y16 := x*cfg.Unit, y*cfg.Unit

	x16 = body.Position().Min.X * cfg.Unit
	rightEdge := x16 + rect.width*cfg.Unit
	maxRight := cfg.ScreenWidth * cfg.Unit
	provider := space.GetTilemapDimensionsProvider()
	if provider != nil {
		maxRight = provider.GetTilemapWidth() * cfg.Unit
	}
	if rightEdge > maxRight {
		x, y, _ = body.ApplyValidPosition(maxRight-rightEdge, true, nil)
		x16, y16 = x*cfg.Unit, y*cfg.Unit
	}

	// Vertical clamping
	minTop := 0
	maxBottom := cfg.ScreenHeight * cfg.Unit
	if provider != nil {
		minTop = (cfg.ScreenHeight - provider.GetTilemapHeight()) * cfg.Unit
		maxBottom = provider.GetTilemapHeight() * cfg.Unit
	}

	if y16 < minTop {
		x, y, _ = body.ApplyValidPosition(minTop-y16, false, nil)
		x16, y16 = x*cfg.Unit, y*cfg.Unit
	}

	y16 = body.Position().Min.Y * cfg.Unit
	bottom := y16 + rect.height*cfg.Unit
	if bottom >= maxBottom {
		if bottom > maxBottom {
			_, _, _ = body.ApplyValidPosition(maxBottom-bottom, false, nil)
		}
		return true
	}

	return false
}
