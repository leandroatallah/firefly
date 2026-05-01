package projectile

import (
	"testing"

	spacepkg "github.com/boilerplate/ebiten-template/internal/engine/physics/space"
)

// AC1 integration — two default projectiles overlapping in a real Space must
// both survive a frame of Manager.Update. Their bodies must still be findable
// by ID after the update.
func TestSpace_TwoDefaultProjectilesOverlap_BothSurvive(t *testing.T) {
	tests := []struct {
		name       string
		x16, y16   int
		v1x16      int
		v2x16      int
		tilemapW   int
		tilemapH   int
		updateRuns int
	}{
		{
			name:       "two bullets at the same position with opposite velocities",
			x16:        100 << 4,
			y16:        100 << 4,
			v1x16:      1 << 4,
			v2x16:      -1 << 4,
			tilemapW:   1000,
			tilemapH:   1000,
			updateRuns: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			space := spacepkg.NewSpace()
			space.SetTilemapDimensionsProvider(&mockTilemapDimensionsProvider{
				width:  tt.tilemapW,
				height: tt.tilemapH,
			})

			mgr := NewManager(space)

			cfg := ProjectileConfig{Width: 4, Height: 4}
			mgr.Spawn(cfg, tt.x16, tt.y16, tt.v1x16, 0, "ownerA")
			mgr.Spawn(cfg, tt.x16, tt.y16, tt.v2x16, 0, "ownerB")

			if len(mgr.projectiles) != 2 {
				t.Fatalf("expected 2 projectiles spawned, got %d", len(mgr.projectiles))
			}

			ids := make([]string, 0, 2)
			for _, p := range mgr.projectiles {
				ids = append(ids, p.body.ID())
			}

			for i := 0; i < tt.updateRuns; i++ {
				mgr.Update()
			}

			for _, id := range ids {
				if got := space.Find(id); got == nil {
					t.Errorf("projectile body %q was removed; expected to survive projectile-vs-projectile contact", id)
				}
			}
			if got := len(mgr.projectiles); got != 2 {
				t.Errorf("Manager.projectiles len = %d, want 2 after frame", got)
			}
		})
	}
}
