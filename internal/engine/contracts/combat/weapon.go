package combat

import (
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
)

// Weapon represents a combat weapon with cooldown-based firing.
type Weapon interface {
	ID() string
	Fire(x16, y16 int, faceDir animation.FacingDirectionEnum, direction body.ShootDirection, state int)
	CanFire() bool
	Update()
	Cooldown() int
	SetCooldown(frames int)
}
