package skill

import (
	combatweapon "github.com/boilerplate/ebiten-template/internal/engine/combat/weapon"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
	"github.com/boilerplate/ebiten-template/internal/engine/input"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
	"github.com/hajimehoshi/ebiten/v2"
)

type ShootingSkill struct {
	SkillBase
	inv            combat.Inventory
	shootHeld      bool
	weaponNextHeld bool
	weaponPrevHeld bool
	handler        body.StateTransitionHandler
	lastDirection  body.ShootDirection
	directionSet   bool
}

func NewShootingSkill(inv combat.Inventory) *ShootingSkill {
	return &ShootingSkill{
		SkillBase: SkillBase{
			state: StateReady,
		},
		inv: inv,
	}
}

func (s *ShootingSkill) SetStateTransitionHandler(handler body.StateTransitionHandler) {
	s.handler = handler
}

func (s *ShootingSkill) HandleInputWithDirection(b body.MovableCollidable, model *physicsmovement.PlatformMovementModel, space body.BodiesSpace, up, down, left, right bool) {
	direction := s.detectShootDirection(b, model, up, down, left, right)

	s.lastDirection = direction
	s.directionSet = true

	weapon := s.inv.ActiveWeapon()
	if weapon == nil || !weapon.CanFire() {
		return
	}

	if s.handler != nil {
		s.handler.TransitionToShooting(direction)
	}

	x16, y16 := b.GetPosition16()

	// Adjust spawn position to account for player width when facing right
	if b.FaceDirection() == animation.FaceDirectionRight {
		x16 += b.GetShape().Width() << 4
	}

	// Set owner to prevent projectile from immediately colliding with player
	if pw, ok := weapon.(*combatweapon.ProjectileWeapon); ok {
		pw.SetOwner(b)
	}

	weapon.Fire(x16, y16, b.FaceDirection(), direction)
}

func (s *ShootingSkill) HandleInput(b body.MovableCollidable, model *physicsmovement.PlatformMovementModel, space body.BodiesSpace) {
	cmds := input.CommandsReader()

	if cmds.WeaponNext && !s.weaponNextHeld {
		s.inv.SwitchNext()
	}
	s.weaponNextHeld = cmds.WeaponNext

	if cmds.WeaponPrev && !s.weaponPrevHeld {
		s.inv.SwitchPrev()
	}
	s.weaponPrevHeld = cmds.WeaponPrev

	if !cmds.Shoot {
		return
	}

	s.HandleInputWithDirection(b, model, space, cmds.Up, cmds.Down, cmds.Left, cmds.Right)
}

func (s *ShootingSkill) Update(b body.MovableCollidable, model *physicsmovement.PlatformMovementModel) {
	// Update inventory weapons (cooldowns)
	s.inv.Update()

	wasHeld := s.shootHeld
	s.shootHeld = input.CommandsReader().Shoot

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
