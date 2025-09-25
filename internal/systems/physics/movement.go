package physics

import (
	"math"

	"github.com/leandroatallah/firefly/internal/config"
)

func increaseVelocity(velocity, acceleration int) int {
	// increaseVelocity applies acceleration to the velocity for a single axis.
	// v_new = v_old + a
	// Capping is handled in the Update loop to correctly manage the 2D vector's magnitude.
	velocity += acceleration
	return velocity
}

func reduceVelocity(velocity int) int {
	// reduceVelocity applies friction to the velocity for a single axis, slowing it down.
	// It brings the velocity to zero if it's smaller than the friction value to prevent jitter.
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
