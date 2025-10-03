package physics

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/systems/input"
)

type PlatformMovementModel struct {
	onGround     bool
	maxFallSpeed int
}

func NewPlatformMovementModel() *PlatformMovementModel {
	return &PlatformMovementModel{
		maxFallSpeed: 12 * config.Unit,
	}
}

func (m *PlatformMovementModel) Update(body *PhysicsBody, space *Space) error {
	// Handle input for player movement
	m.InputHandler(body)

	if applyAxisMovement(body, body.vx16, true, space) {
		body.vx16 = 0
	}

	verticalBlocked := applyAxisMovement(body, body.vy16, false, space)
	landed := false
	if verticalBlocked {
		if body.vy16 > 0 {
			landed = true
		}
		body.vy16 = 0
	}

	if m.clampToPlayArea(body) {
		landed = true
		body.vy16 = 0
	}

	if landed {
		m.onGround = true
	} else if body.vy16 != 0 {
		m.onGround = false
	}

	scaledAccX, _ := smoothDiagonalMovement(body.accelerationX, 0)
	body.vx16 = increaseVelocity(body.vx16, scaledAccX)
	body.vx16 = clampAxisVelocity(body.vx16, body.maxSpeed*config.Unit)

	body.CheckMovementDirectionX()

	body.accelerationX, body.accelerationY = 0, 0
	body.vx16 = reduceVelocity(body.vx16)

	if m.onGround {
		body.vy16 = 0
	} else {
		applyGravity(body, m.maxFallSpeed)
	}

	return nil
}

func (m *PlatformMovementModel) InputHandler(body *PhysicsBody) {
	if body.Immobile() {
		return
	}

	if input.IsSomeKeyPressed(ebiten.KeyA, ebiten.KeyLeft) {
		body.OnMoveLeft(body.Speed())
	}
	if input.IsSomeKeyPressed(ebiten.KeyD, ebiten.KeyRight) {
		body.OnMoveRight(body.Speed())
	}
	if m.onGround && input.IsSomeKeyPressed(ebiten.KeySpace) {
		body.TryJump()
		m.onGround = false
	}
}

func (m *PlatformMovementModel) clampToPlayArea(body *PhysicsBody) bool {
	rect, ok := body.Shape.(*Rect)
	if !ok {
		return false
	}

	if rect.x16 < 0 {
		body.ApplyValidMovement(-rect.x16, true, nil)
	}

	rightEdge := rect.x16 + rect.width*config.Unit
	maxRight := config.ScreenWidth * config.Unit
	if rightEdge > maxRight {
		body.ApplyValidMovement(maxRight-rightEdge, true, nil)
	}

	if rect.y16 < 0 {
		body.ApplyValidMovement(-rect.y16, false, nil)
	}

	bottom := rect.y16 + rect.height*config.Unit
	maxBottom := config.ScreenHeight * config.Unit
	if bottom >= maxBottom {
		if bottom > maxBottom {
			body.ApplyValidMovement(maxBottom-bottom, false, nil)
		}
		return true
	}

	return false
}
