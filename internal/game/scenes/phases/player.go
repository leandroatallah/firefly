package gamescenephases

import (
	"image"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/events"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/platformer"
	"github.com/leandroatallah/firefly/internal/engine/physics/skill"
	gameplayer "github.com/leandroatallah/firefly/internal/game/entity/actors/player"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

func createPlayer(ctx *app.AppContext, playerType gameentitytypes.PlayerType) (platformer.PlatformerActorEntity, error) {
	var f func(*app.AppContext) (platformer.PlatformerActorEntity, error)

	switch playerType {
	case gameentitytypes.ClimberPlayerType:
		f = gameplayer.NewClimberPlayer
	}

	p, err := f(ctx)
	if err != nil {
		return nil, err
	}

	// Create jump skill and set up OnJump callback to publish event
	jumpSkill := skill.NewJumpSkill()
	jumpSkill.OnJump = func(body body.MovableCollidable) {
		pos := body.Position()
		// Bottom center
		jumpPos := image.Point{X: pos.Min.X + pos.Dx()/2, Y: pos.Max.Y}
		ctx.EventManager.Publish(&events.ActorJumpedEvent{
			X: float64(jumpPos.X),
			Y: float64(jumpPos.Y),
		})
	}
	p.GetCharacter().AddSkill(jumpSkill)
	p.GetCharacter().AddSkill(skill.NewHorizontalMovementSkill())

	return p, nil
}
