package actors

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
)

const (
	frameOX = 0
	frameOY = 0
)

type PlayerEntityEnum int

const (
	TopDown PlayerEntityEnum = iota
	Platform
)

type PlayerEntity interface {
	ActorEntity
	Draw(screen *ebiten.Image)
}

func (p PlayerEntityEnum) String() string {
	PlayerEntityMap := map[PlayerEntityEnum]string{
		TopDown:  "TopDown",
		Platform: "Platform",
	}
	return PlayerEntityMap[p]
}

type Player struct {
	Character
}

func NewPlayer(playerEntity PlayerEntityEnum, playerMovementBlocker physics.PlayerMovementBlocker) (PlayerEntity, error) {
	switch playerEntity {
	case TopDown:
		p, err := NewPlayerTopDown(playerMovementBlocker)
		return p, err
	case Platform:
		p, err := NewPlayerPlatform(playerMovementBlocker)
		return p, err
	default:
		return nil, fmt.Errorf("unknown movement model type")
	}
}

// Character Methods
func (p *Player) Update(space body.BodiesSpace) error {
	return p.Character.Update(space)
}

func (p *Player) Draw(screen *ebiten.Image) {
	// TODO: Restore Draw
	p.Character.DrawCollisionBox(screen)
}
