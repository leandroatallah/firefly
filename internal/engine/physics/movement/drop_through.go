package movement

import "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"

// InputSource abstracts the input signals needed for drop-through disambiguation.
type InputSource interface {
	DropHeld() bool
	DuckHeld() bool
	JumpHeld() bool
}

// tryDropThrough registers a pass-through on platform when the actor presses
// down+jump (DropHeld) while grounded on a one-way platform. Vertical velocity
// is intentionally left unchanged — no jump impulse is applied.
func tryDropThrough(actor body.MovableCollidable, platform body.OneWayPlatform, input InputSource) {
	if !input.DropHeld() || !platform.IsOneWay() {
		return
	}
	platform.SetPassThrough(actor, 2)
}
