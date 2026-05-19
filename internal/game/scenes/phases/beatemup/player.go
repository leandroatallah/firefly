package gamebeatemupphase

import (
	"fmt"
	"image"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/events"
	gameplayer "github.com/boilerplate/ebiten-template/internal/game/entity/actors/player"
	beatemupkit "github.com/boilerplate/ebiten-template/internal/kit/actors/beatemup"
	kitbuilder "github.com/boilerplate/ebiten-template/internal/kit/actors/builder"
	beatemupphasescene "github.com/boilerplate/ebiten-template/internal/kit/scenes/phases/beatemup"
	kitskills "github.com/boilerplate/ebiten-template/internal/kit/skills"
)

func newCodyPlayer(ctx *app.AppContext) (beatemupphasescene.Player, error) {
	p, err := gameplayer.NewCodyPlayer(ctx)
	if err != nil {
		return nil, err
	}
	be, ok := p.(beatemupkit.BeatEmUpActorEntity)
	if !ok {
		return nil, fmt.Errorf("CodyPlayer does not satisfy BeatEmUpActorEntity")
	}
	cp, ok := p.(*gameplayer.CodyPlayer)
	if !ok {
		return be, nil
	}
	built, err := kitbuilder.BuildPlayer(cp, kitbuilder.PlayerDeps{
		SpriteData:  cp.GetSpriteData(),
		Inventory:   gameplayer.NewClimberInventory(ctx.ProjectileManager, ctx.VFX),
		MeleeWeapon: gameplayer.NewPlayerMeleeWeapon(),
		VFXManager:  ctx.VFX,
		SkillDeps: kitskills.SkillDeps{
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
		},
		WireState: func(c *actors.Character) {
			gameplayer.WireStateContributors(c, cp)
		},
	})
	if err != nil {
		return nil, err
	}
	return built, nil
}
