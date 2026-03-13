package gameplayer

import (
	"log"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/builder"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/platformer"
	gameplayermethods "github.com/leandroatallah/firefly/internal/game/entity/actors/methods"
	gamestates "github.com/leandroatallah/firefly/internal/game/entity/actors/states"
)

// climberStateTransitionLogic provides custom state handling for the ClimberPlayer,
func climberStateTransitionLogic(c *actors.Character) bool {
	if gameplayermethods.StandardStateTransitionLogic(c) {
		return true
	}

	setNewState := func(s actors.ActorStateEnum) {
		state, err := c.NewState(s)
		if err != nil {
			// Log the error instead of crashing if a state is not registered.
			log.Printf("Failed to create new state %v: %v", s, err)
			return
		}
		c.SetState(state)
	}

	state := c.State()

	if state == gamestates.Rising && c.IsAnimationFinished() {
		setNewState(actors.Idle)
		return true
	}

	if state == gamestates.Exiting || state == gamestates.Lying || state == gamestates.Rising {
		return true
	}

	return true // We've handled the state, so the engine shouldn't.
}

type ClimberPlayer struct {
	*platformer.PlatformerCharacter
	baseSpeed int

	*gameplayermethods.PlayerDeathBehavior
}

// NewClimberPlayer creates a new climber player.
func NewClimberPlayer(ctx *app.AppContext) (platformer.PlatformerActorEntity, error) {
	character, spriteData, statData, stateMap, err := builder.PreparePlatformer(ctx, "internal/game/entity/actors/player/climber.json")
	if err != nil {
		return nil, err
	}

	character.SetStateTransitionHandler(climberStateTransitionLogic)

	player := &ClimberPlayer{
		PlatformerCharacter: character,
	}
	// Set the owner on the embedded character so LastOwner() works correctly
	player.SetOwner(player)
	// Ensure the original character pointer (referenced by physics bodies) also points to the player
	character.SetOwner(player)

	if err = builder.ConfigureCharacter(player, spriteData, statData, stateMap, "player"); err != nil {
		return nil, err
	}
	player.baseSpeed = player.Speed()

	if err = builder.ApplyPlatformerPhysics(player, player); err != nil {
		return nil, err
	}

	character.StateCollisionManager.RefreshCollisions()
	player.PlayerDeathBehavior = gameplayermethods.NewPlayerDeathBehavior(player)

	return player, nil
}

func (p *ClimberPlayer) Update(space body.BodiesSpace) error {
	p.SetHorizontalInertia(-1.0)
	p.SetSpeed(p.baseSpeed)
	p.SetJumpForceMultiplier(1.0)
	return p.Character.Update(space)
}

func (p *ClimberPlayer) GetCharacter() *actors.Character {
	return p.Character
}

func (p *ClimberPlayer) Hurt(damage int) {
	if p.State() == gamestates.Dying {
		return
	}

	state, err := p.NewState(gamestates.Dying)
	if err != nil {
		return
	}
	p.SetState(state)
}
