package gamescenephases

import (
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	gameenemies "github.com/leandroatallah/firefly/internal/game/entity/actors/enemies"
)

type BodyCounter struct {
	wolf       int
	wolfKilled int
}

func (b *BodyCounter) setBodyCounter(space body.BodiesSpace) {
	b.wolf = 0
	b.wolfKilled = 0
	for _, sb := range space.Bodies() {
		switch sb.(type) {
		case *gameenemies.WolfEnemy:
			b.wolf++
		}
	}
}
