package body_test

import (
	"testing"

	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
	"github.com/boilerplate/ebiten-template/internal/engine/physics/space"
)

// T-I1: *ObstacleRect must satisfy space.DepthLaneBody at compile time.
// Story 069 — without this, the depth-lane gate inside space.HasCollision
// never fires for any obstacle pair and falls back to 2D bbox overlap.
func TestObstacleRect_ImplementsDepthLaneBody(t *testing.T) {
	var _ space.DepthLaneBody = (*bodyphysics.ObstacleRect)(nil)
}

// T-I3: ObstacleRect.GroundY() returns the bottom edge in world coords
// (Position().Max.Y). Obstacles are not airborne, so this is altitude-
// independent by construction.
func TestObstacleRect_GroundY_ReturnsBottomEdge(t *testing.T) {
	obs := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 32, 16))
	obs.SetID("wall")
	obs.SetPosition(10, 20)

	got := any(obs).(space.DepthLaneBody).GroundY()
	want := 36 // 20 (top) + 16 (height)
	if got != want {
		t.Errorf("GroundY() = %d; want %d (SetPosition(10,20) SetSize(32,16) → Position().Max.Y)", got, want)
	}
}

// T-I4: ObstacleRect.LaneHalfWidth() = max(height, space.DefaultLaneHalfWidth).
// A short obstacle must NOT be pinned to its tiny half-height; it falls back
// to the engine default so character feet within the default tolerance still
// collide.
func TestObstacleRect_LaneHalfWidth_MaxOfHeightAndDefault(t *testing.T) {
	cases := []struct {
		name   string
		width  int
		height int
		want   int
	}{
		{"height below default falls back to DefaultLaneHalfWidth", 32, 4, space.DefaultLaneHalfWidth},
		{"height above default uses full height", 32, 32, 32},
		{"height equal to default uses default", 32, space.DefaultLaneHalfWidth, space.DefaultLaneHalfWidth},
		{"zero height falls back to default", 32, 0, space.DefaultLaneHalfWidth},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			obs := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, tc.width, tc.height))
			obs.SetID("wall")

			got := any(obs).(space.DepthLaneBody).LaneHalfWidth()
			if got != tc.want {
				t.Errorf("LaneHalfWidth() = %d; want %d (width=%d height=%d)",
					got, tc.want, tc.width, tc.height)
			}
		})
	}
}
