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

type BatEnemy struct {
	*platformer.PlatformerCharacter
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

	enemy.SetGravityEnabled(false)
	enemy.Character.SetMovementState(movement.SideToSide, nil, movement.WithIgnoreLedges(true), movement.WithWaitBeforeTurn(60))

	return enemy, nil
}

func (e *BatEnemy) SetTarget(target body.MovableCollidable) {
	e.Character.MovementState().SetTarget(target)
}

// Character Methods
func (e *BatEnemy) Update(space body.BodiesSpace) error {
	return e.Character.Update(space)
}

func (e *BatEnemy) GetCharacter() *actors.Character {
	return e.Character
}

func (e *BatEnemy) OnTouch(other body.Collidable) {
	owner := other.LastOwner()
	switch owner.(type) {
	case *gameplayer.ClimberPlayer:
		if owner.(platformer.PlatformerActorEntity).State() == gamestates.Dying {
			return
		}

		// Check if player has invincibility skill active
		if invincible, ok := owner.(interface{ IsStarActive() bool; IsGrowActive() bool }); ok {
			if invincible.IsStarActive() || invincible.IsGrowActive() {
				return
			}
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
