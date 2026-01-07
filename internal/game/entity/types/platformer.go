package gameentitytypes

import (
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/physics/body"
	"github.com/leandroatallah/firefly/internal/engine/render/sprites"
)

type CoinCollector interface {
	AddCoinCount(amount int)
	CoinCount() int
}

type PlatformerActorEntity interface {
	actors.ActorEntity
	CoinCollector
}

type PlatformerCharacter struct {
	actors.Character

	coinCount        int
	movementBlockers int
}

func NewPlatformerCharacter(s sprites.SpriteMap, bodyRect *body.Rect) *PlatformerCharacter {
	c := actors.NewCharacter(s, bodyRect)
	return &PlatformerCharacter{
		Character: *c,
	}
}

func (p *PlatformerCharacter) AddCoinCount(amount int) {
	p.coinCount += amount
}
func (p *PlatformerCharacter) CoinCount() int {
	return p.coinCount
}
