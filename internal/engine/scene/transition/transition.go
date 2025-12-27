package transition

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type BaseTransition struct {
	active   bool
	starting bool
	exiting  bool
	onExitCb func()
}

func NewBaseTransition() *BaseTransition {
	return &BaseTransition{}
}

func (t *BaseTransition) Update() {}

func (t *BaseTransition) Draw(screen *ebiten.Image) {}

func (t *BaseTransition) StartTransition(cb func()) {}

func (t *BaseTransition) EndTransition(cb func()) {}
