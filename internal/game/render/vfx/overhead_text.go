package vfx

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/render/camera"
	"github.com/leandroatallah/firefly/internal/engine/render/vfx/text"
)

// OverheadText displays text above an actor, following its position.
type OverheadText struct {
	*text.SimpleFloatingText
	actor   actors.ActorEntity
	offsetX float64
	offsetY float64
}

// NewOverheadText creates text that follows an actor.
func NewOverheadText(actor actors.ActorEntity, msg string, duration int) *OverheadText {
	return &OverheadText{
		SimpleFloatingText: text.NewFloatingText(msg, 0, 0, duration),
		actor:              actor,
		offsetX:            0,
		offsetY:            -10,
	}
}

func (ot *OverheadText) Draw(screen *ebiten.Image, cam *camera.Controller) {
	pos := ot.actor.Position()
	ot.X = float64(pos.Min.X+pos.Dx()/2) + ot.offsetX
	ot.Y = float64(pos.Min.Y) + ot.offsetY
	ot.SimpleFloatingText.Draw(screen, cam)
}
