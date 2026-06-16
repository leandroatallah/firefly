package movement

import (
	"math"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/fp16"
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
	friction := fp16.To16(1) / 4
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
	accelerationFactor := float64(fp16.To16(1) / 6)

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

// clampToPlayArea ensures the body stays within the world boundaries.
// It adjusts the body's position if it goes beyond the edges of the play area.
// It returns true if the body is touching or has gone past the bottom of the screen,
// which can be interpreted as being on the ground for platformer.
func clampToPlayArea(b body.MovableCollidable, space body.BodiesSpace) bool {
	// Use the union of all collision shapes as the effective bounding box.
	// This ensures clamping is based on the actual hitbox, not the sprite bounds.
	collisionRects := b.CollisionPosition()
	if len(collisionRects) == 0 {
		return false
	}

	union := collisionRects[0]
	for i := 1; i < len(collisionRects); i++ {
		union = union.Union(collisionRects[i])
	}

	cfg := config.Get()
	bx, by := b.GetPositionMin()

	// Offset from body origin to the collision bounding box top-left.
	offsetX := union.Min.X - bx
	offsetY := union.Min.Y - by

	collW := union.Dx()
	collH := union.Dy()
	newX, newY := bx, by

	provider := space.GetTilemapDimensionsProvider()
	maxX := cfg.ScreenWidth
	maxY := cfg.ScreenHeight
	if provider != nil {
		maxX = provider.GetTilemapWidth()
		maxY = provider.GetTilemapHeight()
	}

	// --- Horizontal clamping ---
	if union.Min.X < 0 {
		newX = -offsetX
	}
	if union.Min.X+collW > maxX {
		newX = maxX - collW - offsetX
	}

	// --- Vertical clamping ---
	if union.Min.Y < 0 {
		newY = -offsetY
	}

	isOnGround := false
	if union.Min.Y+collH >= maxY {
		if union.Min.Y+collH > maxY {
			newY = maxY - collH - offsetY
		}
		isOnGround = true
	}

	if newX != bx || newY != by {
		b.SetPosition(newX, newY)
	}

	return isOnGround
}
