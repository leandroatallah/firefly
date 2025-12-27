package gamescenelevels

import (
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	gameplayer "github.com/leandroatallah/firefly/internal/game/entity/actors/player"
)

func createPlayer() (actors.ActorEntity, error) {
	p, err := gameplayer.NewCherryPlayer()
	if err != nil {
		return nil, err
	}

	return p, nil
}
