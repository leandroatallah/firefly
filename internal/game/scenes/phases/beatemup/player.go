package gamebeatemupphase

import (
	"fmt"
	"image"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/vfx"
	"github.com/boilerplate/ebiten-template/internal/engine/data/schemas"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/builder"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/events"
	gameplayer "github.com/boilerplate/ebiten-template/internal/game/entity/actors/player"
	beatemupkit "github.com/boilerplate/ebiten-template/internal/kit/actors/beatemup"
	"github.com/boilerplate/ebiten-template/internal/kit/combat/weapon"
	kitskills "github.com/boilerplate/ebiten-template/internal/kit/skills"
)

type beatemupMeleeActor interface {
	body.Movable
	GetCharacter() *actors.Character
	GetSpriteData() *schemas.SpriteData
	SetMelee(w *weapon.MeleeWeapon, vfxMgr vfx.Manager)
	SetInventory(inv interface{})
}

func createPlayer(ctx *app.AppContext) (beatemupkit.BeatEmUpActorEntity, error) {
	p, err := gameplayer.NewCodyPlayer(ctx)
	if err != nil {
		return nil, err
	}

	be, ok := p.(beatemupkit.BeatEmUpActorEntity)
	if !ok {
		return nil, fmt.Errorf("CodyPlayer does not satisfy BeatEmUpActorEntity")
	}

	actor, ok := p.(beatemupMeleeActor)
	if !ok {
		return be, nil
	}

	spriteData := actor.GetSpriteData()
	if spriteData == nil {
		return be, nil
	}

	deps := kitskills.SkillDeps{
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

	skills := kitskills.FromConfig(spriteData.Skills, deps)
	if err := builder.ApplySkills(p, skills); err != nil {
		return nil, err
	}

	gameplayer.WireStateContributors(actor.GetCharacter(), actor)

	return be, nil
}
