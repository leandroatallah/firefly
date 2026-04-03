package gamescenephases

import (
	"image"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/events"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/platformer"
	"github.com/boilerplate/ebiten-template/internal/engine/physics/skill"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/fp16"
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

	addJumpSkill(ctx, p)
	p.GetCharacter().AddSkill(skill.NewDashSkill())
	p.GetCharacter().AddSkill(skill.NewHorizontalMovementSkill())

	// Add shooting skill - get shooter from current scene
	if currentScene := ctx.SceneManager.CurrentScene(); currentScene != nil {
		if shooter, ok := currentScene.(body.Shooter); ok {
			addShootingSkill(p, shooter)
		}
	}

	return p, nil
}

func addJumpSkill(ctx *app.AppContext, p platformer.PlatformerActorEntity) {
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
	jumpSkill.SetJumpCutMultiplier(0.4)
	p.GetCharacter().AddSkill(jumpSkill)
}

func addShootingSkill(p platformer.PlatformerActorEntity, shooter body.Shooter) {
	shootingSkill := skill.NewShootingSkill(shooter, 15, 4, fp16.To16(6), 4)
	p.GetCharacter().AddSkill(shootingSkill)
}
