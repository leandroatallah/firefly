package actors

import (
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
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

func (p PlayerEntityEnum) String() string {
	PlayerEntityMap := map[PlayerEntityEnum]string{
		TopDown:  "TopDown",
		Platform: "Platform",
	}
	return PlayerEntityMap[p]
}

// TODO: Should remove and use only Character?
type Player struct {
	Character
}

func (p *Player) Update(space body.BodiesSpace) error {
	return p.Character.Update(space)
}
