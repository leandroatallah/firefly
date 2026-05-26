package space

import (
	"image"
	"testing"
)

// TestHasCollisionDepthLane_AirbornePlayerScenarios exercises the depth-lane
// gate against the two scenarios that story 069 targets: an airborne player
// near a same-depth wall (should still collide on bbox overlap) and an
// airborne player near a different-depth background wall (should be filtered
// out by the gate even though screen-bbox overlaps).
//
// These are the *behavioural* contracts the new ObstacleRect / BeatEmUpCharacter
// implementations must satisfy through HasCollision.
func TestHasCollisionDepthLane_AirbornePlayerScenarios(t *testing.T) {
	// Player rect is the SCREEN rect (Y already had altitude subtracted by
	// the body layer). The depth-lane gate compares GroundY (pre-altitude),
	// which is what DepthLaneBody.GroundY() returns.
	cases := []struct {
		name        string
		playerRect  image.Rectangle
		playerGY    int
		playerHalfW int
		wallRect    image.Rectangle
		wallGY      int
		wallHalfW   int
		want        bool
	}{
		{
			// T-S3: player jumping (alt=40, GroundY=100, screen Y=60) into
			// a wall at the same ground depth — must still block.
			name:        "airborne player same-depth wall blocks",
			playerRect:  image.Rect(100, 60, 116, 76),
			playerGY:    100,
			playerHalfW: 8,
			wallRect:    image.Rect(108, 60, 140, 110),
			wallGY:      110,
			wallHalfW:   16, // height-derived; |100-110|=10 <= 16
			want:        true,
		},
		{
			// T-S4: player jumping into a background wall whose ground
			// projection is far below — bbox can overlap because the wall is
			// tall, but the gate must filter it out.
			name:        "airborne player different-depth wall no block",
			playerRect:  image.Rect(100, 60, 116, 76),
			playerGY:    100,
			playerHalfW: 8,
			wallRect:    image.Rect(108, 40, 140, 168),
			wallGY:      168,
			wallHalfW:   16, // |100-168|=68 > 16
			want:        false,
		},
		{
			// T-S6 boundary: diff == max(half) → inclusive, returns true.
			name:        "boundary inclusive diff equals half",
			playerRect:  image.Rect(0, 0, 16, 16),
			playerGY:    100,
			playerHalfW: 8,
			wallRect:    image.Rect(8, 8, 24, 24),
			wallGY:      108,
			wallHalfW:   8,
			want:        true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			player := &depthLaneCollidable{
				testCollidable: newTestCollidable("player", tc.playerRect, false),
				groundY:        tc.playerGY,
				laneHalfWidth:  tc.playerHalfW,
			}
			wall := &depthLaneCollidable{
				testCollidable: newTestCollidable("wall", tc.wallRect, true),
				groundY:        tc.wallGY,
				laneHalfWidth:  tc.wallHalfW,
			}
			if got := HasCollision(player, wall); got != tc.want {
				t.Fatalf("HasCollision(player, wall) = %v, want %v (player GY=%d half=%d; wall GY=%d half=%d)",
					got, tc.want, tc.playerGY, tc.playerHalfW, tc.wallGY, tc.wallHalfW)
			}
		})
	}
}
