package gameplayer

import (
	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
	"github.com/leandroatallah/firefly/internal/engine/systems/sprites"
)

type CherryPlayer struct {
	actors.Player

	coinCount int
}

// TODO: Move to the right place
func GetSprites(assets map[string]string) (sprites.SpriteMap, error) {
	var s sprites.SpriteAssets
	for key, value := range assets {
		var state animation.SpriteState
		switch key {
		case "idle":
			state = actors.Idle
		case "walk":
			state = actors.Walk
		case "hurt":
			state = actors.Hurted
		default:
			continue
		}
		s = s.AddSprite(state, value)
	}
	result, err := sprites.LoadSprites(s)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewCherryPlayer(
	movementBlocker physics.PlayerMovementBlocker,
) (actors.PlayerEntity, error) {
	spriteData, statData, err := actors.ParseJsonPlayer("internal/game/actors/player/cherry.json")
	if err != nil {
		return nil, err
	}

	assets, err := GetSprites(spriteData.Assets)
	if err != nil {
		return nil, err
	}

	// Create character instance
	character := actors.NewCharacter(assets, spriteData.FrameRate)
	character.SetFaceDirection(spriteData.FacingDirection)

	// Create bodies
	body := spriteData.BodyRect
	collision := spriteData.CollisionRect
	bodyRect := physics.NewRect(body.Rect())
	collisionRect := physics.NewRect(collision.Rect())

	// TODO: (maybe) Create a builder with director to automate this process. Maybe should use a template method.
	player := &CherryPlayer{
		Player: actors.Player{Character: *character},
	}
	player.SetBody(bodyRect)
	player.SetCollisionArea(collisionRect)
	player.SetTouchable(player)
	// TODO: Create set stats method
	player.SetMaxHealth(statData.Health)
	player.SetSpeedAndMaxSpeed(
		statData.Speed, statData.MaxSpeed,
	)

	model, err := physics.NewMovementModel(physics.Platform, movementBlocker)
	if err != nil {
		return nil, err
	}
	player.SetMovementModel(model)

	return player, nil
}

func (p *CherryPlayer) AddCoinCount(amount int) {
	p.coinCount += amount
}
func (p *CherryPlayer) CoinCount() int {
	return p.coinCount
}
