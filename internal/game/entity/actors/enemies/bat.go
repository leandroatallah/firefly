package gameenemies

import (
	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/builder"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/movement"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/platformer"
	gameplayer "github.com/boilerplate/ebiten-template/internal/game/entity/actors/player"
	gamestates "github.com/boilerplate/ebiten-template/internal/game/entity/actors/states"
	kitactors "github.com/boilerplate/ebiten-template/internal/kit/actors"
	kitcombat "github.com/boilerplate/ebiten-template/internal/kit/combat"
	kitcombatweapon "github.com/boilerplate/ebiten-template/internal/kit/combat/weapon"
)

type BatEnemy struct {
	*platformer.PlatformerCharacter
	*kitactors.ShooterCharacter
}

// NewBatEnemy creates a new bat enemy.
func NewBatEnemy(ctx *app.AppContext, x, y int, id string) (*BatEnemy, error) {
	character, spriteData, statData, stateMap, err := builder.PreparePlatformer(ctx, "assets/entities/enemies/bat.json")
	if err != nil {
		return nil, err
	}

	enemy := &BatEnemy{PlatformerCharacter: character}
	// Set the owner on the embedded character so LastOwner() works correctly
	enemy.SetOwner(enemy)
	enemy.SetPosition(x, y)

	if err = builder.ConfigureCharacter(enemy, spriteData, statData, stateMap, id); err != nil {
		return nil, err
	}

	if err = builder.ApplyPlatformerPhysics(enemy, nil); err != nil {
		return nil, err
	}

	shooter, err := kitcombatweapon.ConfigureEnemy(enemy, spriteData.Weapon, ctx.ProjectileManager)
	if err != nil {
		return nil, err
	}
	enemy.ShooterCharacter = kitactors.NewShooterCharacter(shooter)

	enemy.GetCharacter().SetFaction(kitcombat.FactionEnemy)
	enemy.SetGravityEnabled(false)
	enemy.SetMovementState(movement.SideToSide, nil, movement.WithIgnoreLedges(true), movement.WithWaitBeforeTurn(60))

	return enemy, nil
}

func (e *BatEnemy) SetTarget(target body.MovableCollidable) {
	e.Character.MovementState().SetTarget(target)
	if e.Shooter() != nil {
		e.Shooter().SetTarget(target)
	}
}

// Character Methods
func (e *BatEnemy) Update(space body.BodiesSpace) error {
	e.UpdateShooter()
	return e.Character.Update(space)
}

func (e *BatEnemy) GetCharacter() *actors.Character {
	return e.Character
}

func (e *BatEnemy) OnTouch(other body.Collidable) {
	// Only react to actor bodies. A projectile's body.Owner() is a *MovableBody,
	// which does not implement PlatformerActorEntity — this prevents projectile
	// hits from being mistaken for player contact damage.
	if _, ok := other.Owner().(platformer.PlatformerActorEntity); !ok {
		return
	}
	owner := other.LastOwner()
	switch owner.(type) {
	case *gameplayer.ClimberPlayer:
		if owner.(platformer.PlatformerActorEntity).State() == gamestates.Dying {
			return
		}

		if alive, ok := owner.(platformer.AlivePlayer); ok {
			alive.Hurt(1)
		}
	}
}

func (e *BatEnemy) OnDie() {}

func (e *BatEnemy) IsEnemy() bool {
	return true
}
