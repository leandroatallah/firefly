package gameplatformerphase

import (
	"image"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/events"
	gameplayer "github.com/boilerplate/ebiten-template/internal/game/entity/actors/player"
	kitbuilder "github.com/boilerplate/ebiten-template/internal/kit/actors/builder"
	platformerphasescene "github.com/boilerplate/ebiten-template/internal/kit/scenes/phases/platformer"
	kitskills "github.com/boilerplate/ebiten-template/internal/kit/skills"
	"github.com/hajimehoshi/ebiten/v2"
)

func newClimberPlayer(ctx *app.AppContext) (platformerphasescene.Player, error) {
	p, err := gameplayer.NewClimberPlayer(ctx)
	if err != nil {
		return nil, err
	}
	cp, ok := p.(*gameplayer.ClimberPlayer)
	if !ok {
		return p, nil
	}
	return kitbuilder.BuildPlayer(cp, kitbuilder.PlayerDeps{
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
}

// makeClimberDebugHook creates a debug draw hook for the climber player.
func makeClimberDebugHook(ctx *app.AppContext) func(*ebiten.Image) {
	return func(screen *ebiten.Image) {
		if !config.Get().CollisionBox {
			return
		}
		p, ok := ctx.ActorManager.GetPlayer()
		if !ok {
			return
		}
		cp, ok := p.(*gameplayer.ClimberPlayer)
		if !ok {
			return
		}
		mc := cp.MeleeController()
		if mc == nil {
			return
		}
		// Melee hitbox draw requires camera which is managed by the kit scene.
		// Access camera via ActorManager is not possible; camera offset is skipped.
		if _, active := mc.Weapon().ActiveHitboxRect(); !active {
			return
		}
		_ = screen
	}
}
