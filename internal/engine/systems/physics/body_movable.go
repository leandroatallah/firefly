package physics

import (
	"github.com/leandroatallah/firefly/internal/engine/config"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
)

type MovableBody struct {
	body.Shape

	vx16          int
	vy16          int
	accelerationX int
	accelerationY int
	speed         int
	maxSpeed      int
	immobile      bool
	faceDirection FacingDirectionEnum
}

func (b *MovableBody) Move() {
	panic("You should implement this method in derivated structs")
}

func (b *MovableBody) MoveX(distance int) {
	b.accelerationX = distance * config.Unit
}

func (b *MovableBody) MoveY(distance int) {
	b.accelerationY = distance * config.Unit
}

// TODO: Should it be moved to Movable?
func (b *MovableBody) OnMoveLeft(distance int) {
	b.MoveX(-distance)
}
func (b *MovableBody) OnMoveUpLeft(distance int) {
	b.MoveX(-distance)
	b.MoveY(-distance)
}
func (b *MovableBody) OnMoveDownLeft(distance int) {
	b.MoveX(-distance)
	b.MoveY(distance)
}
func (b *MovableBody) OnMoveRight(distance int) {
	b.MoveX(distance)
}
func (b *MovableBody) OnMoveUpRight(distance int) {
	b.MoveX(distance)
	b.MoveY(-distance)
}
func (b *MovableBody) OnMoveDownRight(distance int) {
	b.MoveX(distance)
	b.MoveY(distance)
}
func (b *MovableBody) OnMoveUp(distance int) {
	b.MoveY(-distance)
}
func (b *MovableBody) OnMoveDown(distance int) {
	b.MoveY(distance)
}

// TODO: Improve this method (split of find out a better approach)
func (b *MovableBody) SetSpeedAndMaxSpeed(speed, maxSpeed int) {
	b.speed = speed
	b.maxSpeed = maxSpeed
}

func (b *MovableBody) Speed() int {
	return b.speed
}

func (b *MovableBody) Immobile() bool {
	return b.immobile
}

func (b *MovableBody) SetImmobile(immobile bool) {
	b.immobile = immobile
}

func (b *MovableBody) FaceDirection() FacingDirectionEnum {
	return b.faceDirection
}

func (b *MovableBody) IsWalking() bool {
	return b.vx16 != 0 || b.vy16 != 0
}

// Platform methods
func (b *MovableBody) TryJump(force int) {
	b.vy16 = -force * config.Unit
}
