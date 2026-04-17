package combat

import (
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
)

// Faction identifies which side an actor or projectile belongs to.
// The zero value is FactionNeutral, which is safe by default.
type Faction int

const (
	FactionNeutral Faction = iota
	FactionPlayer
	FactionEnemy
)

// Factioned is an interface for entities that belong to a Faction.
type Factioned interface {
	Faction() Faction
}

// Weapon represents a combat weapon with cooldown-based firing.
type Weapon interface {
	ID() string
	Fire(x16, y16 int, faceDir animation.FacingDirectionEnum, direction body.ShootDirection, state int)
	CanFire() bool
	Update()
	Cooldown() int
	SetCooldown(frames int)
	SetOwner(owner interface{})
}

// ShootMode determines when an enemy fires.
type ShootMode int

const (
	ShootModeOnSight ShootMode = iota
	ShootModeAlways
)

// TargetBody is the minimal set of methods needed for targeting.
type TargetBody interface {
	GetPosition16() (int, int)
}

// EnemyShooter manages automatic firing logic for an enemy.
type EnemyShooter interface {
	SetTarget(TargetBody)
	Target() TargetBody
	Range() int
	Mode() ShootMode
	Direction() body.ShootDirection
	ShootState() (interface{}, bool)
	TryFire() bool
	Update()
}
