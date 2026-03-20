package gameenemies

import (
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/builder"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/movement"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/platformer"
	gameplayer "github.com/leandroatallah/firefly/internal/game/entity/actors/player"
	gamestates "github.com/leandroatallah/firefly/internal/game/entity/actors/states"
)

type SwarmEnemy struct {
	*platformer.PlatformerCharacter
}

// NewSwarmEnemy creates a new swarm enemy.
func NewSwarmEnemy(ctx *app.AppContext, x, y int, id string) (*SwarmEnemy, error) {
	character, spriteData, statData, stateMap, err := builder.PreparePlatformer(ctx, "assets/entities/enemies/swarm.json")
	if err != nil {
		return nil, err
	}

	enemy := &SwarmEnemy{PlatformerCharacter: character}
	// Set the owner on the embedded character so LastOwner() works correctly
	enemy.SetOwner(enemy)
	enemy.SetPosition(x, y)

	if err = builder.ConfigureCharacter(enemy, spriteData, statData, stateMap, id); err != nil {
		return nil, err
	}

	if err = builder.ApplyPlatformerPhysics(enemy, nil); err != nil {
		return nil, err
	}

	enemy.SetHorizontalInertia(1.0)
	enemy.Character.SetMovementState(
		movement.SideToSide,
		nil,
		movement.WithIgnoreLedges(true),
		movement.WithWaitBeforeTurn(120),
		movement.WithVerticalMovement(true),
	)

	enemy.Character.SetStateTransitionHandler(func(c *actors.Character) bool {
		// Override the default platformer state transition logic.
		// Since the SwarmEnemy is a flying unit and lacks specific "Jumping" or "Falling" animations,
		// we force the state to remain either Idle or Walking based on horizontal movement.
		// This prevents the sprite from disappearing or flickering when the enemy moves vertically
		// (which would normally trigger a Jump/Fall state in the default platformer model).
		desiredState := actors.Idle
		if vx, _ := c.Velocity(); vx != 0 {
			desiredState = actors.Walking
		}

		if c.State() != desiredState {
			s, err := c.NewState(desiredState)
			if err == nil {
				c.SetState(s)
			}
		}
		return true
	})

	return enemy, nil
}

func (e *SwarmEnemy) SetTarget(target body.MovableCollidable) {
	e.Character.MovementState().SetTarget(target)
}

// Character Methods
func (e *SwarmEnemy) Update(space body.BodiesSpace) error {
	return e.Character.Update(space)
}

func (e *SwarmEnemy) GetCharacter() *actors.Character {
	return e.Character
}

func (e *SwarmEnemy) OnTouch(other body.Collidable) {
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

func (e *SwarmEnemy) OnDie() {}

func (e *SwarmEnemy) IsEnemy() bool {
	return true
}
