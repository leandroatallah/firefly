package gamescenelevels

import (
	"github.com/leandroatallah/firefly/internal/engine/actors"
	gameplayer "github.com/leandroatallah/firefly/internal/game/actors/player"
)

func createPlayer() (actors.PlayerEntity, error) {
	p, err := gameplayer.NewCherryPlayer()
	if err != nil {
		return nil, err
	}

	return p, nil
}
