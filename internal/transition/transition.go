package transition

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Transition interface {
	Update()
	Draw(screen *ebiten.Image)
	StartTransition(func())
	EndTransition(func())
}

type BaseTransition struct {
	active   bool
	starting bool
	exiting  bool
	onExitCb func()
}

func NewBaseTransition() *BaseTransition {
	return &BaseTransition{}
}

func (t *BaseTransition) Update() {
	panic("You should implement this method in derivated structs")
}

func (t *BaseTransition) Draw(screen *ebiten.Image) {
	panic("You should implement this method in derivated structs")
}

func (t *BaseTransition) StartTransition(cb func()) {
	panic("You should implement this method in derivated structs")
}

func (t *BaseTransition) EndTransition(cb func()) {
	panic("You should implement this method in derivated structs")
}
