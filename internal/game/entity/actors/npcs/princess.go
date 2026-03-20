package gamenpcs

import (
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/builder"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/movement"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/platformer"
	gameplayermethods "github.com/leandroatallah/firefly/internal/game/entity/actors/methods"
)

type Princess struct {
	*platformer.PlatformerCharacter
}

// NewPrincess creates a new princess NPC.
func NewPrincess(ctx *app.AppContext, x, y int, id string) (*Princess, error) {
	character, spriteData, statData, stateMap, err := builder.PreparePlatformer(ctx, "assets/entities/npcs/princess.json")
	if err != nil {
		return nil, err
	}

	princess := &Princess{PlatformerCharacter: character}
	// Set the owner on the embedded character so LastOwner() works correctly
	princess.SetOwner(princess)
	princess.SetPosition(x, y)

	if err = builder.ConfigureCharacter(princess, spriteData, statData, stateMap, id); err != nil {
		return nil, err
	}

	if err = builder.ApplyPlatformerPhysics(princess, nil); err != nil {
		return nil, err
	}

	princess.Character.SetMovementState(movement.Idle, nil)
	princess.Character.SetStateTransitionHandler(gameplayermethods.StandardStateTransitionLogic)

	return princess, nil
}

func (s *Princess) SetTarget(target body.MovableCollidable) {
	s.Character.SetMovementState(movement.Wander, target)
}

// Character Methods
func (s *Princess) Update(space body.BodiesSpace) error {
	return s.Character.Update(space)
}

func (s *Princess) GetCharacter() *actors.Character {
	return s.Character
}

func (s *Princess) OnTouch(other body.Collidable) {}

func (s *Princess) Hurt(damage int) {}

func (s *Princess) OnDie() {}
