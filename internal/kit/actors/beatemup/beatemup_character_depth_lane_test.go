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

// T-I5: BeatEmUpCharacter.GroundY() returns the floor-projected bottom edge
// of the body (position.Y + height), matching the convention used by
// ObstacleRect.GroundY(). Altitude must NOT affect the result — the
// depth-lane gate compares floor positions, not screen positions.
//
// newTestFixtures creates an 8×8 body, so GroundY = SetPosition Y + 8.
func TestBeatEmUpCharacter_GroundY_AltitudeIndependent(t *testing.T) {
	fsys, stateMap, spriteData, bodyRect := newTestFixtures()

	c, err := beatemup.NewBeatEmUpCharacter(fsys, stateMap, spriteData, bodyRect, nil)
	if err != nil {
		t.Fatalf("NewBeatEmUpCharacter returned error: %v", err)
	}

	c.SetPosition(0, 150)

	// Pre-altitude: GroundY = 150 + 8 = 158 (floor bottom).
	gotGrounded := any(c).(space.DepthLaneBody).GroundY()
	if gotGrounded != 158 {
		t.Errorf("GroundY() with altitude=0 = %d; want 158", gotGrounded)
	}

	// Airborne — altitude must not affect GroundY.
	c.SetAltitude(40)
	gotAirborne := any(c).(space.DepthLaneBody).GroundY()
	if gotAirborne != 158 {
		t.Errorf("GroundY() with altitude=40 = %d; want 158 (altitude must not subtract)", gotAirborne)
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
