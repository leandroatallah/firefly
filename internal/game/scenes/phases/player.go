package gamescenephases

import (
	"image"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/vfx"
	"github.com/boilerplate/ebiten-template/internal/engine/data/schemas"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/builder"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/events"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/platformer"
	engineskill "github.com/boilerplate/ebiten-template/internal/engine/physics/skill"
	gameplayer "github.com/boilerplate/ebiten-template/internal/game/entity/actors/player"
	gameentitytypes "github.com/boilerplate/ebiten-template/internal/game/entity/types"
	"github.com/boilerplate/ebiten-template/internal/kit/combat/weapon"
)

// interface to satisfies platformer melee players
type platformerMeleeActor interface {
	body.Movable
	GetCharacter() *actors.Character
	GetSpriteData() *schemas.SpriteData
	SetMelee(w *weapon.MeleeWeapon, vfxMgr vfx.Manager)
	SetInventory(inv interface{})
}

func createPlayer(ctx *app.AppContext, playerType gameentitytypes.PlayerType) (platformer.PlatformerActorEntity, error) {
	var f func(*app.AppContext) (platformer.PlatformerActorEntity, error)

	switch playerType {
	case gameentitytypes.ClimberPlayerType:
		f = gameplayer.NewClimberPlayer
	case gameentitytypes.CodyPlayerType:
		f = gameplayer.NewCodyPlayer
	}

	p, err := f(ctx)
	if err != nil {
		return nil, err
	}

	actor, ok := p.(platformerMeleeActor)
	if !ok {
		return p, nil
	}

	spriteData := actor.GetSpriteData()
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

	actor.SetInventory(deps.Inventory)
	actor.SetMelee(gameplayer.NewPlayerMeleeWeapon(), ctx.VFX)

	if err := builder.ApplySkills(p, *spriteData, deps); err != nil {
		return nil, err
	}

	gameplayer.WireStateContributors(actor.GetCharacter(), actor)

	return p, nil
}
