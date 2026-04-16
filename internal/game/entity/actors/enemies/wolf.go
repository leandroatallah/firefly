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
)

type WolfEnemy struct {
	*platformer.PlatformerCharacter
}

// NewWolfEnemy creates a new wolf enemy.
func NewWolfEnemy(ctx *app.AppContext, x, y int, id string) (*WolfEnemy, error) {
	character, spriteData, statData, stateMap, err := builder.PreparePlatformer(ctx, "assets/entities/enemies/wolf.json")
	if err != nil {
		return nil, err
	}

	enemy := &WolfEnemy{PlatformerCharacter: character}
	// Set the owner on the embedded character so LastOwner() works correctly
	enemy.SetOwner(enemy)
	enemy.SetPosition(x, y)

	if err = builder.ConfigureCharacter(enemy, spriteData, statData, stateMap, id); err != nil {
		return nil, err
	}

	if err = builder.ApplyPlatformerPhysics(enemy, nil); err != nil {
		return nil, err
	}

	enemy.SetMovementState(
		movement.SideToSide,
		nil,
		movement.WithWaitBeforeTurn(60),
		movement.WithLimitToRoom(true),
	)

	return enemy, nil
}

func (e *WolfEnemy) SetTarget(target body.MovableCollidable) {
	e.Character.MovementState().SetTarget(target)
}

// Character Methods
func (e *WolfEnemy) Update(space body.BodiesSpace) error {
	return e.Character.Update(space)
}

func (e *WolfEnemy) GetCharacter() *actors.Character {
	return e.Character
}

func (e *WolfEnemy) OnTouch(other body.Collidable) {
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

func (e *WolfEnemy) OnDie() {}

func (e *WolfEnemy) IsEnemy() bool {
	return true
}
