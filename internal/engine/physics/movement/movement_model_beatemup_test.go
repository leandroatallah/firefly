package movement

import (
	"testing"

	contractsbody "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
	"github.com/boilerplate/ebiten-template/internal/engine/physics/space"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/fp16"
)

// altitudeSpyBodiesSpace forwards every BodiesSpace call to an inner
// BodiesSpace (real *space.Space) but intercepts ResolveCollisions to record
// the moving body's Altitude16 at the instant of the call. This is how the
// test verifies that Block 1 (the zero-altitude wrap around
// ApplyValidPosition) has been removed — story 069 AC-6.
//
// With Block 1 in place, BeatEmUpMovementModel.Update temporarily zeroes
// b.Altitude16 before ApplyValidPosition, so the spy observes 0 for every
// frame in which the body started airborne. After AC-6, the spy must NEVER
// observe a transition to 0 when the player started with a non-zero altitude.
type altitudeSpyBodiesSpace struct {
	contractsbody.BodiesSpace
	observed    []int
	trackBodyID string
}

func newAltitudeSpy(inner contractsbody.BodiesSpace, trackID string) *altitudeSpyBodiesSpace {
	return &altitudeSpyBodiesSpace{BodiesSpace: inner, trackBodyID: trackID}
}

func (s *altitudeSpyBodiesSpace) ResolveCollisions(b contractsbody.Collidable) (bool, bool) {
	if b != nil && b.ID() == s.trackBodyID {
		if alt, ok := b.(interface{ Altitude16() int }); ok {
			s.observed = append(s.observed, alt.Altitude16())
		}
	}
	return s.BodiesSpace.ResolveCollisions(b)
}

// T-M3: Block 1 (zero-altitude wrap) removed [AC-6].
// When the player starts a frame airborne (Altitude16 > 0), the movement
// model must NOT zero altitude before calling ApplyValidPosition. The spy
// captures Altitude16 at the moment of each ResolveCollisions call.
func TestBeatEmUpMovementModel_NoZeroAltitudeWrap_DuringApplyValidPosition(t *testing.T) {
	config.Set(&config.AppConfig{
		ScreenWidth:  320,
		ScreenHeight: 240,
		Physics:      config.PhysicsConfig{SpeedMultiplier: 1.0},
	})

	inner := space.NewSpace()
	sp := newAltitudeSpy(inner, "player")

	player := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 16, 16))
	player.SetID("player")
	player.SetPosition(100, 100)
	player.SetMaxSpeed(30)
	player.SetVelocity(fp16.To16(5), 0) // X-axis movement → ApplyValidPosition fires
	startAlt16 := fp16.To16(40)
	player.SetAltitude16(startAlt16)
	inner.AddBody(player)

	model := NewBeatEmUpMovementModel(nil)
	if err := model.Update(player, sp); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	if len(sp.observed) == 0 {
		t.Fatalf("expected ResolveCollisions to be called at least once for player")
	}
	for i, alt16 := range sp.observed {
		if alt16 == 0 {
			t.Fatalf("call #%d: observed Altitude16=0 during ApplyValidPosition; want non-zero (Block 1 must be removed — AC-6)", i)
		}
	}
}

// T-M2: airborne player still blocked by a same-depth obstacle when the
// depth-lane gate admits the pair [AC-3, AC-4 — integration with movement].
//
// Requires *ObstacleRect to implement space.DepthLaneBody (story 069 AC-1)
// AND Block 1 to be removed (AC-6). After implementation, the player's
// screen rect overlaps the wall's screen rect (both visible), their ground
// depths match, the gate admits the pair, and the player must be blocked.
func TestBeatEmUpMovementModel_AirbornePlayer_BlockedBySameDepthWall(t *testing.T) {
	config.Set(&config.AppConfig{
		ScreenWidth:  320,
		ScreenHeight: 240,
		Physics:      config.PhysicsConfig{SpeedMultiplier: 1.0},
	})

	sp := space.NewSpace()

	player := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 16, 16))
	player.SetID("player")
	player.SetPosition(100, 100)
	player.SetMaxSpeed(30)
	player.SetVelocity(fp16.To16(10), 0)
	player.SetAltitude16(fp16.To16(40)) // airborne; screen Y = 60
	sp.AddBody(player)

	// Wall ground depth equals player ground depth (both at Y=116 — player
	// bottom edge pre-altitude: groundY = position.Y + height = 100 + 16 = 116).
	// Wall position is chosen so its screen rect overlaps the airborne player's
	// screen rect: player screen rect = (100, 60, 116, 76); wall top = 68 < 76.
	// wallGroundY = wallY + wallH = 68 + 48 = 116, matching the player.
	wall := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 16, 48))
	wall.SetID("wall")
	wall.SetIsObstructive(true)
	wall.SetPosition(120, 68) // bottom edge = 68 + 48 = 116, same as player ground depth
	wall.AddCollisionBodies()
	sp.AddBody(wall)

	model := NewBeatEmUpMovementModel(nil)
	if err := model.Update(player, sp); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	pos := player.Position().Min
	// Player must be blocked before crossing into the wall (wall.Min.X=120;
	// player width=16, so player.Min.X must remain <= 104 to not overlap).
	if pos.X > 104 {
		t.Errorf("airborne same-depth: expected player blocked before crossing into wall (pos.X <= 104); got pos.X=%d", pos.X)
	}
	if pos.X < 100 {
		t.Errorf("airborne same-depth: expected player to advance (pos.X >= 100); got pos.X=%d", pos.X)
	}
}

// T-M1: airborne player NOT blocked by a depth-mismatched wall [AC-3].
//
// Direct gate behavioural assertion — independent of movement-model
// integration. Post-implementation: gate denies; HasCollision must be
// false because |player.GroundY - wall.GroundY| far exceeds the max
// LaneHalfWidth. Today the obstacle does not implement DepthLaneBody,
// so the gate falls through and HasCollision returns true on screen-bbox
// overlap — failing this assertion.
func TestBeatEmUpMovementModel_AirbornePlayer_NotBlockedByDifferentDepthWall(t *testing.T) {
	config.Set(&config.AppConfig{
		ScreenWidth:  320,
		ScreenHeight: 240,
		Physics:      config.PhysicsConfig{SpeedMultiplier: 1.0},
	})

	sp := space.NewSpace()

	player := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 16, 16))
	player.SetID("player")
	player.SetPosition(100, 100)
	player.SetMaxSpeed(30)
	player.SetAltitude16(fp16.To16(40))
	sp.AddBody(player)

	// Background wall whose screen rect overlaps the airborne player's
	// screen rect but whose ground projection is far below. Player ground
	// bottom = 100 + 16 = 116; wall ground bottom = 60 + 16 = 76 (depth
	// diff = 40, far greater than DefaultLaneHalfWidth=8 and the wall's
	// height-derived half-width = 16). Player's screen Y after altitude
	// subtraction = 100 - 40 = 60, so the player rect (60..76) overlaps
	// the wall rect (60..76) in screen space — pure 2D bbox would collide.
	wall := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 16, 16))
	wall.SetID("wall")
	wall.SetIsObstructive(true)
	wall.SetPosition(108, 60)
	wall.AddCollisionBodies()
	sp.AddBody(wall)

	if got := space.HasCollision(player, wall); got {
		t.Errorf("HasCollision(airborne player @ ground 116, wall @ ground 76) = true; want false (depth-lane gate must deny — story 069 AC-3)")
	}
}
