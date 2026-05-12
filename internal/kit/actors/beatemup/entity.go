package beatemup

import (
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
)

// BeatEmUpActorEntity is the genre contract for beat-em-up actor entities.
// It extends ActorEntity with altitude accessors required by the beat-em-up
// draw-order and physics systems.
type BeatEmUpActorEntity interface {
	actors.ActorEntity

	Altitude16() int
	SetAltitude16(int)
}
