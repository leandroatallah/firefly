package enemies

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/actors"
	"github.com/leandroatallah/firefly/internal/systems/physics"
)

type BaseEnemy struct {
	actors.Character
}

func NewBaseEnemy() *BaseEnemy {
	return &BaseEnemy{}
}

// Character Methods
func (e *BaseEnemy) Update(space *physics.Space) error {
	return e.Character.Update(space)
}

func (e *BaseEnemy) Draw(screen *ebiten.Image) {
	e.Character.Draw(screen)
}
