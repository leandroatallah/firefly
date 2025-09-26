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
func (e *BaseEnemy) SetBody(rect *physics.Rect) actors.ActorEntity {
	return e.Character.SetBody(rect)
}

func (e *BaseEnemy) SetCollisionArea(rect *physics.Rect) actors.ActorEntity {
	return e.Character.SetCollisionArea(rect)
}

func (e *BaseEnemy) Update(boundaries []physics.Body) error {
	return e.Character.Update(boundaries)
}

func (e *BaseEnemy) Draw(screen *ebiten.Image) {
	e.Character.Draw(screen)
}

func (e *BaseEnemy) HandleMovement() {}
