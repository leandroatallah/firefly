package platformer

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/physics/space"
)

// T-I7 (negative): *PlatformerCharacter must NOT implement space.DepthLaneBody.
// Plain 2D platformer scenes must retain pure-bbox collision behaviour;
// opting them into the depth-lane gate would silently change collision
// outcomes for every existing platformer game.
func TestPlatformerCharacter_DoesNotImplementDepthLaneBody(t *testing.T) {
	if _, ok := any((*PlatformerCharacter)(nil)).(space.DepthLaneBody); ok {
		t.Fatal("*PlatformerCharacter must NOT implement space.DepthLaneBody (regression risk: 2D scenes would gain depth gating)")
	}
}
