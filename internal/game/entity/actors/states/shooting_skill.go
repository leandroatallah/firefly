package gamestates

import (
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	contractsbody "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
	"github.com/hajimehoshi/ebiten/v2"
)

type ShootingConfig struct {
	CooldownFrames int
	SpawnOffsetX16 int
	BulletSpeedX16 int
	YOffset        int
}

type ShootingSkill struct {
	cfg      ShootingConfig
	toggler  *OffsetToggler
	shooter  contractsbody.Shooter
	cooldown int
}

type shootingBody interface {
	GetPosition16() (int, int)
	FaceDirection() animation.FacingDirectionEnum
	Owner() interface{}
}

func NewShootingSkill(cfg ShootingConfig, shooter contractsbody.Shooter) *ShootingSkill {
	return &ShootingSkill{
		cfg:     cfg,
		toggler: NewOffsetToggler(cfg.YOffset),
		shooter: shooter,
	}
}

func (s *ShootingSkill) Update(body contractsbody.MovableCollidable, model *physicsmovement.PlatformMovementModel) {
	if s.cooldown > 0 {
		s.cooldown--
	}
}

func (s *ShootingSkill) HandleInput(body contractsbody.MovableCollidable, model *physicsmovement.PlatformMovementModel, space contractsbody.BodiesSpace) {
	if s.cooldown > 0 {
		return
	}

	if !ebiten.IsKeyPressed(s.ActivationKey()) {
		return
	}

	x16, y16 := body.GetPosition16()
	
	// Apply horizontal spawn offset based on facing direction
	if body.FaceDirection() == animation.FaceDirectionRight {
		x16 += s.cfg.SpawnOffsetX16
	} else {
		x16 -= s.cfg.SpawnOffsetX16
	}
	
	// Apply vertical offset toggle
	y16 += s.toggler.Next()

	// Bullet speed based on facing direction (horizontal only)
	speedX16 := s.cfg.BulletSpeedX16
	if body.FaceDirection() == animation.FaceDirectionLeft {
		speedX16 = -speedX16
	}

	s.shooter.SpawnBullet(x16, y16, speedX16, 0, body.Owner())
	s.cooldown = s.cfg.CooldownFrames
}

func (s *ShootingSkill) ActivationKey() ebiten.Key {
	return ebiten.KeyX
}

func (s *ShootingSkill) IsActive() bool {
	return s.cooldown > 0
}
