package skill

import (
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
	"github.com/hajimehoshi/ebiten/v2"
)

type ShootingSkill struct {
	SkillBase
	shooter       body.Shooter
	spawnOffsetX  int
	bulletSpeed   int
	toggler       *OffsetToggler
	shootHeld     bool
	handler       body.StateTransitionHandler
	lastDirection body.ShootDirection
	directionSet  bool
}

func NewShootingSkill(shooter body.Shooter, cooldownFrames, spawnOffsetX16, bulletSpeedX16, yOffset int) *ShootingSkill {
	return &ShootingSkill{
		SkillBase: SkillBase{
			state:    StateReady,
			cooldown: cooldownFrames,
		},
		shooter:      shooter,
		spawnOffsetX: spawnOffsetX16,
		bulletSpeed:  bulletSpeedX16,
		toggler:      NewOffsetToggler(yOffset),
	}
}

func (s *ShootingSkill) SetStateTransitionHandler(handler body.StateTransitionHandler) {
	s.handler = handler
}

func (s *ShootingSkill) HandleInputWithDirection(b body.MovableCollidable, model *physicsmovement.PlatformMovementModel, space body.BodiesSpace, up, down, left, right bool) {
	direction := s.detectShootDirection(b, model, up, down, left, right)
	
	directionChanged := s.directionSet && direction != s.lastDirection
	if directionChanged && s.handler != nil {
		s.handler.TransitionToShooting(direction)
		if s.state == StateActive {
			s.state = StateReady
			s.timer = 0
		}
	}
	s.lastDirection = direction
	s.directionSet = true

	if s.state == StateReady && !directionChanged {
		if s.handler != nil {
			s.handler.TransitionToShooting(direction)
		}
		
		x16, y16 := b.GetPosition16()
		vx16, vy16 := s.calculateBulletVelocity(direction, b.FaceDirection())
		offsetX, offsetY := s.calculateSpawnOffset(direction, b.FaceDirection())
		
		s.shooter.SpawnBullet(x16+offsetX, y16+offsetY, vx16, vy16, b)
		
		s.state = StateActive
		s.timer = s.cooldown
	}
}

func (s *ShootingSkill) HandleInput(b body.MovableCollidable, model *physicsmovement.PlatformMovementModel, space body.BodiesSpace) {
	wasHeld := s.shootHeld
	s.shootHeld = true

	direction := body.ShootDirectionStraight
	
	if s.shootHeld && !wasHeld && s.handler != nil {
		s.handler.TransitionToShooting(direction)
	}
	s.lastDirection = direction
	s.directionSet = true

	if s.shootHeld && s.state == StateReady {
		x16, y16 := b.GetPosition16()
		dir := b.FaceDirection()

		offsetX := s.spawnOffsetX
		speedX := s.bulletSpeed
		if dir == animation.FaceDirectionLeft {
			offsetX = -offsetX
			speedX = -speedX
		}

		yOffset := s.toggler.Next()
		s.shooter.SpawnBullet(x16+offsetX, y16+yOffset, speedX, 0, b)

		s.state = StateActive
		s.timer = s.cooldown
	}
}

func (s *ShootingSkill) Update(b body.MovableCollidable, model *physicsmovement.PlatformMovementModel) {
	if s.state == StateActive {
		s.timer--
		if s.timer <= 0 {
			s.state = StateReady
		}
	}

	wasHeld := s.shootHeld
	s.shootHeld = ebiten.IsKeyPressed(ebiten.KeyX)

	if !s.shootHeld && wasHeld && s.handler != nil {
		s.handler.TransitionFromShooting()
	}
}

func (s *ShootingSkill) ActivationKey() ebiten.Key {
	return ebiten.KeyX
}

func (s *ShootingSkill) detectShootDirection(b body.MovableCollidable, model *physicsmovement.PlatformMovementModel, up, down, left, right bool) body.ShootDirection {
	isDucking := false
	if duckable, ok := b.(interface{ IsDucking() bool }); ok {
		isDucking = duckable.IsDucking()
	}
	
	if isDucking {
		return body.ShootDirectionStraight
	}
	
	isGrounded := model != nil && model.OnGround()
	
	if down && !isGrounded {
		if left || right {
			return body.ShootDirectionDiagonalDownForward
		}
		return body.ShootDirectionDown
	}
	
	if up {
		if left || right {
			return body.ShootDirectionDiagonalUpForward
		}
		return body.ShootDirectionUp
	}
	
	return body.ShootDirectionStraight
}

func (s *ShootingSkill) calculateBulletVelocity(direction body.ShootDirection, faceDir animation.FacingDirectionEnum) (vx16, vy16 int) {
	speed := s.bulletSpeed
	sign := 1
	if faceDir == animation.FaceDirectionLeft {
		sign = -1
	}
	
	switch direction {
	case body.ShootDirectionStraight:
		return sign * speed, 0
	case body.ShootDirectionUp:
		return 0, -speed
	case body.ShootDirectionDown:
		return 0, speed
	case body.ShootDirectionDiagonalUpForward:
		return sign * speed * 707 / 1000, -speed * 707 / 1000
	case body.ShootDirectionDiagonalDownForward:
		return sign * speed * 707 / 1000, speed * 707 / 1000
	}
	return sign * speed, 0
}

func (s *ShootingSkill) calculateSpawnOffset(direction body.ShootDirection, faceDir animation.FacingDirectionEnum) (offsetX16, offsetY16 int) {
	offset := s.spawnOffsetX
	sign := 1
	if faceDir == animation.FaceDirectionLeft {
		sign = -1
	}
	
	switch direction {
	case body.ShootDirectionStraight:
		return sign * offset, s.toggler.Next()
	case body.ShootDirectionUp:
		return 0, -offset
	case body.ShootDirectionDown:
		return 0, offset
	case body.ShootDirectionDiagonalUpForward:
		return sign * offset * 707 / 1000, -offset * 707 / 1000
	case body.ShootDirectionDiagonalDownForward:
		return sign * offset * 707 / 1000, offset * 707 / 1000
	}
	return sign * offset, s.toggler.Next()
}
