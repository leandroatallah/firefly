package weapon

import (
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
)

// enemyOwner is the minimal set of capabilities EnemyShooting needs from its owner.
type enemyOwner interface {
	GetPosition16() (int, int)
	FaceDirection() animation.FacingDirectionEnum
	SetFaceDirection(animation.FacingDirectionEnum)
	GetShape() body.Shape
	State() actors.ActorStateEnum
}

// EnemyShooting implements combat.EnemyShooter, encapsulating all gate logic for
// automatic enemy firing.
type EnemyShooting struct {
	owner      enemyOwner
	weapon     *ProjectileWeapon
	target     combat.TargetBody
	rangePx    int
	mode       combat.ShootMode
	direction  body.ShootDirection
	shootState actors.ActorStateEnum
	hasState   bool
}

// NewEnemyShooting constructs an EnemyShooting.
// stateGate=true activates the shoot-state gate; shootState is the required owner state.
func NewEnemyShooting(
	owner enemyOwner,
	w *ProjectileWeapon,
	rangePx int,
	mode combat.ShootMode,
	direction body.ShootDirection,
	shootState actors.ActorStateEnum,
	stateGate bool,
) *EnemyShooting {
	return &EnemyShooting{
		owner:      owner,
		weapon:     w,
		rangePx:    rangePx,
		mode:       mode,
		direction:  direction,
		shootState: shootState,
		hasState:   stateGate,
	}
}

func (e *EnemyShooting) SetTarget(t combat.TargetBody)  { e.target = t }
func (e *EnemyShooting) Target() combat.TargetBody      { return e.target }
func (e *EnemyShooting) Range() int                     { return e.rangePx }
func (e *EnemyShooting) Mode() combat.ShootMode         { return e.mode }
func (e *EnemyShooting) Direction() body.ShootDirection { return e.direction }

func (e *EnemyShooting) ShootState() (interface{}, bool) {
	return e.shootState, e.hasState
}

// TryFire runs the gate chain. Returns true if a projectile was actually spawned.
func (e *EnemyShooting) TryFire() bool {
	// Gate 1: state gate
	if e.hasState {
		if e.owner.State() != e.shootState {
			return false
		}
	}

	// Gate 2 (OnSight only): target must be set and within range
	if e.mode == combat.ShootModeOnSight {
		if e.target == nil {
			return false
		}
		ox16, _ := e.owner.GetPosition16()
		tx16, _ := e.target.GetPosition16()
		dx := (ox16 - tx16) / 16
		if dx < 0 {
			dx = -dx
		}
		if e.rangePx > 0 && dx > e.rangePx {
			return false
		}
		// Set face direction toward target
		if tx16 < ox16 {
			e.owner.SetFaceDirection(animation.FaceDirectionLeft)
		} else {
			e.owner.SetFaceDirection(animation.FaceDirectionRight)
		}
	}

	// Gate 3: cooldown
	if !e.weapon.CanFire() {
		return false
	}

	// Fire
	x16, y16 := e.owner.GetPosition16()
	faceDir := e.owner.FaceDirection()
	e.weapon.Fire(x16, y16, faceDir, e.direction, int(e.owner.State()))
	return true
}

// Update ticks the weapon cooldown each frame, then attempts to fire.
func (e *EnemyShooting) Update() {
	e.weapon.Update()
	e.TryFire()
}
