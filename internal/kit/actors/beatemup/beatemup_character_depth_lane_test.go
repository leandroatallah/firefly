package beatemup_test

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/physics/space"
	"github.com/boilerplate/ebiten-template/internal/kit/actors/beatemup"
)

// T-I2: *BeatEmUpCharacter must satisfy space.DepthLaneBody at compile time
// so that game-layer subclasses (e.g. CodyPlayer) pick the methods up via
// embedding without re-declaration.
func TestBeatEmUpCharacter_ImplementsDepthLaneBody(t *testing.T) {
	var _ space.DepthLaneBody = (*beatemup.BeatEmUpCharacter)(nil)
}

// T-I5: BeatEmUpCharacter.GroundY() returns the pre-altitude ground Y
// (i.e. the body's y16>>scale value). Altitude must NOT subtract from it —
// the depth-lane gate compares depth, not screen position.
func TestBeatEmUpCharacter_GroundY_AltitudeIndependent(t *testing.T) {
	fsys, stateMap, spriteData, bodyRect := newTestFixtures()

	c, err := beatemup.NewBeatEmUpCharacter(fsys, stateMap, spriteData, bodyRect, nil)
	if err != nil {
		t.Fatalf("NewBeatEmUpCharacter returned error: %v", err)
	}

	c.SetPosition(0, 150)

	// Pre-altitude.
	gotGrounded := any(c).(space.DepthLaneBody).GroundY()
	if gotGrounded != 150 {
		t.Errorf("GroundY() with altitude=0 = %d; want 150", gotGrounded)
	}

	// Airborne — altitude must not affect GroundY.
	c.SetAltitude(40)
	gotAirborne := any(c).(space.DepthLaneBody).GroundY()
	if gotAirborne != 150 {
		t.Errorf("GroundY() with altitude=40 = %d; want 150 (altitude must not subtract)", gotAirborne)
	}
}

// T-I6: BeatEmUpCharacter.LaneHalfWidth() returns space.DefaultLaneHalfWidth.
func TestBeatEmUpCharacter_LaneHalfWidth_UsesDefault(t *testing.T) {
	fsys, stateMap, spriteData, bodyRect := newTestFixtures()

	c, err := beatemup.NewBeatEmUpCharacter(fsys, stateMap, spriteData, bodyRect, nil)
	if err != nil {
		t.Fatalf("NewBeatEmUpCharacter returned error: %v", err)
	}

	got := any(c).(space.DepthLaneBody).LaneHalfWidth()
	if got != space.DefaultLaneHalfWidth {
		t.Errorf("LaneHalfWidth() = %d; want %d (space.DefaultLaneHalfWidth)",
			got, space.DefaultLaneHalfWidth)
	}
}
