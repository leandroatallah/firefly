package gamescenephases

import (
	"image"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/builder"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/events"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/platformer"
	engineskill "github.com/boilerplate/ebiten-template/internal/engine/physics/skill"
	gameplayer "github.com/boilerplate/ebiten-template/internal/game/entity/actors/player"
	gameentitytypes "github.com/boilerplate/ebiten-template/internal/game/entity/types"
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

	climber, ok := p.(*gameplayer.ClimberPlayer)
	if !ok {
		return p, nil
	}

	spriteData := climber.GetSpriteData()
	if spriteData == nil {
		return p, nil
	}

	deps := engineskill.SkillDeps{
		Inventory:         gameplayer.NewClimberInventory(ctx.ProjectileManager, ctx.VFX),
		ProjectileManager: ctx.ProjectileManager,
		OnJump: func(b interface{}) {
			bodyObj, ok := b.(body.MovableCollidable)
			if !ok {
				return
			}
			pos := bodyObj.Position()
			jumpPos := image.Point{X: pos.Min.X + pos.Dx()/2, Y: pos.Max.Y}
			ctx.EventManager.Publish(&events.ActorJumpedEvent{
				X: float64(jumpPos.X),
				Y: float64(jumpPos.Y),
			})
		},
		EventManager: ctx.EventManager,
	}

	climber.SetInventory(deps.Inventory)

	if err := builder.ApplySkills(p, *spriteData, deps); err != nil {
		return nil, err
	}

	gameplayer.WireStateContributors(climber)

	return p, nil
}
