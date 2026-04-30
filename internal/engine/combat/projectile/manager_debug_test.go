package projectile

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
)

// TestManager_DrawCollisionBoxesWithOffset verifies the new debug-draw helper
// invokes the supplied callback once per active projectile body, never
// mutates manager state, and is robust to nil callbacks.
func TestManager_DrawCollisionBoxesWithOffset(t *testing.T) {
	tests := []struct {
		name        string
		spawnCount  int
		passNil     bool
		wantInvokes int
	}{
		{name: "no projectiles invokes callback zero times", spawnCount: 0, passNil: false, wantInvokes: 0},
		{name: "single projectile invokes callback once", spawnCount: 1, passNil: false, wantInvokes: 1},
		{name: "three projectiles invoke callback three times", spawnCount: 3, passNil: false, wantInvokes: 3},
		{name: "nil callback is a no-op (no panic)", spawnCount: 1, passNil: true, wantInvokes: 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSpace := &mockBodiesSpace{}
			mgr := NewManager(mockSpace)

			cfg := ProjectileConfig{Width: 2, Height: 1, Damage: 0}
			spawnedIDs := make(map[string]struct{}, tc.spawnCount)
			for i := 0; i < tc.spawnCount; i++ {
				mgr.Spawn(cfg, (100+i)<<4, 50<<4, 5<<4, 0, nil)
				spawnedIDs[mgr.projectiles[i].body.ID()] = struct{}{}
			}

			invokes := 0
			gotIDs := make(map[string]struct{}, tc.spawnCount)
			var cb func(b body.Collidable)
			if !tc.passNil {
				cb = func(b body.Collidable) {
					invokes++
					if b == nil {
						t.Errorf("callback received nil body.Collidable")
						return
					}
					gotIDs[b.ID()] = struct{}{}
				}
			}

			// Must not panic regardless of nil callback or empty manager.
			mgr.DrawCollisionBoxesWithOffset(cb)

			if invokes != tc.wantInvokes {
				t.Errorf("callback invoked %d times, want %d", invokes, tc.wantInvokes)
			}
			if !tc.passNil && tc.spawnCount > 0 {
				if len(gotIDs) != len(spawnedIDs) {
					t.Errorf("collected ID set size %d, want %d", len(gotIDs), len(spawnedIDs))
				}
				for id := range spawnedIDs {
					if _, ok := gotIDs[id]; !ok {
						t.Errorf("expected callback ID %q missing from collected set", id)
					}
				}
			}
		})
	}
}

// TestManager_DrawCollisionBoxesWithOffset_NoMutation guarantees that calling
// the helper does not mutate projectile state (slice length, position).
// This is the spec's no-side-effects contract for AC-5.
func TestManager_DrawCollisionBoxesWithOffset_NoMutation(t *testing.T) {
	mockSpace := &mockBodiesSpace{}
	mgr := NewManager(mockSpace)

	cfg := ProjectileConfig{Width: 2, Height: 1, Damage: 0}
	mgr.Spawn(cfg, 100<<4, 50<<4, 5<<4, 0, nil)

	preLen := len(mgr.projectiles)
	preX, preY := mgr.projectiles[0].body.GetPositionMin()
	preX16, preY16 := mgr.projectiles[0].body.GetPosition16()

	mgr.DrawCollisionBoxesWithOffset(func(_ body.Collidable) {})

	if got := len(mgr.projectiles); got != preLen {
		t.Errorf("len(projectiles) = %d after Draw helper, want %d (no mutation)", got, preLen)
	}
	x, y := mgr.projectiles[0].body.GetPositionMin()
	if x != preX || y != preY {
		t.Errorf("GetPositionMin() = (%d,%d) after helper, want (%d,%d)", x, y, preX, preY)
	}
	x16, y16 := mgr.projectiles[0].body.GetPosition16()
	if x16 != preX16 || y16 != preY16 {
		t.Errorf("GetPosition16() = (%d,%d) after helper, want (%d,%d)", x16, y16, preX16, preY16)
	}
}
