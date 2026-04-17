package builder

import (
	"fmt"

	"github.com/boilerplate/ebiten-template/internal/engine/combat/weapon"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
	"github.com/boilerplate/ebiten-template/internal/engine/data/schemas"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
)

// ConfigureEnemyWeapon reads the weapon block from cfg and wires an EnemyShooter
// onto the character's embedded body.
// Returns (nil, nil) if cfg is nil.
// Returns (nil, error) if manager is nil when cfg is non-nil.
func ConfigureEnemyWeapon(
	character actors.ActorEntity,
	cfg *schemas.EnemyWeaponConfig,
	manager combat.ProjectileManager,
) (combat.EnemyShooter, error) {
	if cfg == nil {
		return nil, nil
	}
	if manager == nil {
		return nil, fmt.Errorf("projectile manager must not be nil")
	}

	// Parse ShootMode
	var mode combat.ShootMode
	switch cfg.ShootMode {
	case "", "on_sight":
		mode = combat.ShootModeOnSight
	case "always":
		mode = combat.ShootModeAlways
	default:
		return nil, fmt.Errorf("unknown shoot_mode: %q", cfg.ShootMode)
	}

	// Parse ShootDirection
	var dir body.ShootDirection
	switch cfg.ShootDirection {
	case "", "horizontal":
		dir = body.ShootDirectionStraight
	case "vertical":
		dir = body.ShootDirectionDown
	default:
		return nil, fmt.Errorf("unknown shoot_direction: %q", cfg.ShootDirection)
	}

	// Parse ShootState gate
	var (
		stateEnum actors.ActorStateEnum
		stateGate bool
	)
	if cfg.ShootState != "" {
		enum, ok := actors.GetStateEnum(cfg.ShootState)
		if !ok {
			return nil, fmt.Errorf("unknown shoot_state: %q", cfg.ShootState)
		}
		stateEnum = enum
		stateGate = true
	}

	// Build weapon; id = character.ID() + "_weapon"
	weaponID := character.ID() + "_weapon"
	w := weapon.NewProjectileWeapon(
		weaponID,
		cfg.Cooldown,
		cfg.ProjectileType,
		cfg.Speed,
		manager,
		"",
		0, 0,
	)
	w.SetDamage(cfg.Damage)

	// Owner is the character's embedded MovableBody
	ch := character.GetCharacter()
	w.SetOwner(ch)

	shooter := weapon.NewEnemyShooting(ch, w, cfg.Range, mode, dir, stateEnum, stateGate)
	return shooter, nil
}
