package gameplayer

import (
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/builder"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/platformer"
	gameplayermethods "github.com/leandroatallah/firefly/internal/game/entity/actors/methods"
	gamestates "github.com/leandroatallah/firefly/internal/game/entity/actors/states"
	gameskill "github.com/leandroatallah/firefly/internal/game/physics/skill"
)

// climberStateTransitionLogic provides custom state handling for the ClimberPlayer,
func climberStateTransitionLogic(c *actors.Character) bool {
	if gameplayermethods.StandardStateTransitionLogic(c) {
		return true
	}

	state := c.State()

	if state == gamestates.Rising && c.IsAnimationFinished() {
		c.SetNewStateFatal(actors.Idle)
		return true
	}

	if state == gamestates.Exiting || state == gamestates.Lying || state == gamestates.Rising {
		return true
	}

	return false
}

type ClimberPlayer struct {
	*platformer.PlatformerCharacter
	baseSpeed   int
	freezeSkill *gameskill.FreezeSkill
	growSkill   *gameskill.GrowSkill

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
		freezeSkill:         gameskill.NewFreezeSkill(),
		growSkill:           gameskill.NewGrowSkill(),
	}
	// Set the owner on the embedded character so LastOwner() works correctly
	player.SetOwner(player)
	// Ensure the original character pointer (referenced by physics bodies) also points to the player
	character.SetOwner(player)

	character.AddSkill(player.freezeSkill)
	character.AddSkill(player.growSkill)

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

func (p *ClimberPlayer) ActivateFreezeSkill() {
	if p.freezeSkill != nil {
		p.freezeSkill.RequestActivation()
	}
}

func (p *ClimberPlayer) ActivateGrowSkill() {
	if p.growSkill != nil {
		p.growSkill.RequestActivation()
	}
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

	p.SetNewStateFatal(gamestates.Dying)
}
