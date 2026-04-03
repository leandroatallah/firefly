package gameplayer

import (
	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/builder"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/platformer"
	gameplayermethods "github.com/boilerplate/ebiten-template/internal/game/entity/actors/methods"
	gamestates "github.com/boilerplate/ebiten-template/internal/game/entity/actors/states"
	"github.com/hajimehoshi/ebiten/v2"
)

// climberStateTransitionLogic provides custom state handling for the ClimberPlayer.
func climberStateTransitionLogic(c *actors.Character) bool {
	state := c.State()

	if state == gamestates.Exiting || state == gamestates.Dying || state == gamestates.Dead {
		return true
	}

	return false
}

type ClimberPlayer struct {
	*platformer.PlatformerCharacter
	baseSpeed int

	*gameplayermethods.PlayerDeathBehavior
}

// NewClimberPlayer creates a new climber player.
func NewClimberPlayer(ctx *app.AppContext) (platformer.PlatformerActorEntity, error) {
	character, spriteData, statData, stateMap, err := builder.PreparePlatformer(ctx, "assets/entities/player/climber.json")
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
	// Check for ducking input
	duckHeld := ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS)
	p.SetDucking(duckHeld && !p.IsFalling() && !p.IsGoingUp())
	
	// When ducking, prevent horizontal movement but allow facing direction change
	if p.IsDucking() {
		p.SetSpeed(0)
		p.SetHorizontalInertia(0)
		
		// Allow facing direction change while ducking
		if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
			p.SetFaceDirection(animation.FaceDirectionLeft)
		} else if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
			p.SetFaceDirection(animation.FaceDirectionRight)
		}
	} else {
		p.SetHorizontalInertia(-1.0)
		p.SetSpeed(p.baseSpeed)
	}
	
	p.SetJumpForceMultiplier(1.0)
	
	return p.Character.Update(space)
}

func (p *ClimberPlayer) GetCharacter() *actors.Character {
	return p.Character
}

func (p *ClimberPlayer) Hurt(damage int) {
	if p.State() == gamestates.Dying || p.State() == gamestates.Dead {
		return
	}

	p.SetNewStateFatal(gamestates.Dying)
}

func (p *ClimberPlayer) OnTouch(other body.Collidable) {
	// Standard player touch behavior
}

func (p *ClimberPlayer) OnBlock(other body.Collidable) {
	// Required to implement body.Touchable to avoid recursion if we rely on embedded CollidableBody.OnBlock
}
