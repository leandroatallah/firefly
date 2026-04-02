package gamestates

import (
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	contractsbody "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
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

func (s *ShootingSkill) Update(body shootingBody) {
	if s.cooldown > 0 {
		s.cooldown--
	}
	
	if s.cooldown > 0 {
		return
	}

	x16, y16 := body.GetPosition16()
	
	if body.FaceDirection() == animation.FaceDirectionRight {
		x16 += s.cfg.SpawnOffsetX16
	} else {
		x16 -= s.cfg.SpawnOffsetX16
	}
	
	y16 += s.toggler.Next()

	speedX16 := s.cfg.BulletSpeedX16
	if body.FaceDirection() == animation.FaceDirectionLeft {
		speedX16 = -speedX16
	}

	s.shooter.SpawnBullet(x16, y16, speedX16, body.Owner())
	s.cooldown = s.cfg.CooldownFrames
}
