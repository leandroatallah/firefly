package gamescenephases

import (
	gameplayer "github.com/leandroatallah/firefly/internal/game/entity/actors/player"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

func createPlayer() (gameentitytypes.PlatformerActorEntity, error) {
	p, err := gameplayer.NewCherryPlayer()
	if err != nil {
		return nil, err
	}

	return p, nil
}
