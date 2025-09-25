package physics

import "github.com/hajimehoshi/ebiten/v2"

type BaseEnemy struct {
	Character
}

func NewBaseEnemy() *BaseEnemy {
	return &BaseEnemy{}
}

// Character Methods
func (e *BaseEnemy) SetBody(rect *Rect) ActorEntity {
	return e.Character.SetBody(rect)
}

func (e *BaseEnemy) SetCollisionArea(rect *Rect) ActorEntity {
	return e.Character.SetCollisionArea(rect)
}

func (e *BaseEnemy) Update(boundaries []Body) error {
	return e.Character.Update(boundaries, e.HandleMovement)
}

func (e *BaseEnemy) Draw(screen *ebiten.Image) {
	e.Character.Draw(screen)
}

func (e *BaseEnemy) HandleMovement() {
	panic("Implement me")
}
