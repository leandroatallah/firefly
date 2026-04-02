package skill

import (
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
	"github.com/hajimehoshi/ebiten/v2"
)

type ActorStateEnum interface{}

type Stateful interface {
	State() ActorStateEnum
	SetState(interface{})
	NewState(ActorStateEnum) (interface{}, error)
}

type ShootingSkill struct {
	SkillBase
	shooter      body.Shooter
	spawnOffsetX int
	bulletSpeed  int
	toggler      *OffsetToggler
	shootHeld    bool

	idleState           ActorStateEnum
	walkingState        ActorStateEnum
	jumpingState        ActorStateEnum
	fallingState        ActorStateEnum
	idleShootingState   ActorStateEnum
	walkingShootingState ActorStateEnum
	jumpingShootingState ActorStateEnum
	fallingShootingState ActorStateEnum
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

func (s *ShootingSkill) SetStateEnums(idle, walking, jumping, falling, idleShooting, walkingShooting, jumpingShooting, fallingShooting ActorStateEnum) {
	s.idleState = idle
	s.walkingState = walking
	s.jumpingState = jumping
	s.fallingState = falling
	s.idleShootingState = idleShooting
	s.walkingShootingState = walkingShooting
	s.jumpingShootingState = jumpingShooting
	s.fallingShootingState = fallingShooting
}

func (s *ShootingSkill) HandleInput(b body.MovableCollidable, model *physicsmovement.PlatformMovementModel, space body.BodiesSpace) {
	wasHeld := s.shootHeld
	s.shootHeld = true

	if s.shootHeld && !wasHeld {
		s.transitionToShootingState(b)
	}

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
		s.shooter.SpawnBullet(x16+offsetX, y16+yOffset, speedX, b)

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

	if !s.shootHeld && wasHeld {
		s.transitionToBaseState(b)
	}
}

func (s *ShootingSkill) ActivationKey() ebiten.Key {
	return ebiten.KeyX
}

func (s *ShootingSkill) transitionToShootingState(b body.MovableCollidable) {
	if s.idleState == nil {
		return
	}

	actor, ok := b.(Stateful)
	if !ok {
		return
	}

	currentState := actor.State()
	var newState ActorStateEnum

	switch currentState {
	case s.idleState:
		newState = s.idleShootingState
	case s.walkingState:
		newState = s.walkingShootingState
	case s.jumpingState:
		newState = s.jumpingShootingState
	case s.fallingState:
		newState = s.fallingShootingState
	default:
		return
	}

	state, err := actor.NewState(newState)
	if err == nil {
		actor.SetState(state)
	}
}

func (s *ShootingSkill) transitionToBaseState(b body.MovableCollidable) {
	if s.idleState == nil {
		return
	}

	actor, ok := b.(Stateful)
	if !ok {
		return
	}

	currentState := actor.State()
	var newState ActorStateEnum

	switch currentState {
	case s.idleShootingState:
		newState = s.idleState
	case s.walkingShootingState:
		newState = s.walkingState
	case s.jumpingShootingState:
		newState = s.jumpingState
	case s.fallingShootingState:
		newState = s.fallingState
	default:
		return
	}

	state, err := actor.NewState(newState)
	if err == nil {
		actor.SetState(state)
	}
}
