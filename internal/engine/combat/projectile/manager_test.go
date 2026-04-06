package projectile

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/hajimehoshi/ebiten/v2"
)

func TestManager_Spawn(t *testing.T) {
	tests := []struct {
		name       string
		spawnCfg   ProjectileConfig
		x16, y16   int
		vx16, vy16 int
		owner      interface{}
	}{
		{
			name:     "spawn increases projectile count",
			spawnCfg: ProjectileConfig{Width: 2, Height: 1, Damage: 10},
			x16:      100 << 16,
			y16:      50 << 16,
			vx16:     5 << 16,
			vy16:     0,
			owner:    "player",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addBodyCalled := 0
			mockSpace := &mockBodiesSpace{
				AddBodyFunc: func(_ body.Collidable) {
					addBodyCalled++
				},
			}

			mgr := NewManager(mockSpace)
			mgr.Spawn(tt.spawnCfg, tt.x16, tt.y16, tt.vx16, tt.vy16, tt.owner)

			if len(mgr.projectiles) != 1 {
				t.Errorf("expected 1 projectile, got %d", len(mgr.projectiles))
			}
			if addBodyCalled != 1 {
				t.Errorf("expected AddBody called once, got %d", addBodyCalled)
			}
		})
	}
}

func TestManager_Update_OutOfBounds(t *testing.T) {
	tests := []struct {
		name          string
		tilemapWidth  int
		tilemapHeight int
		initialX16    int
		initialY16    int
		vx16, vy16    int
		expectRemoved bool
	}{
		{
			name:          "projectile removed when out of bounds",
			tilemapWidth:  100,
			tilemapHeight: 100,
			initialX16:    95 << 16,
			initialY16:    50 << 16,
			vx16:          10 << 16,
			vy16:          0,
			expectRemoved: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queueForRemovalCalled := 0
			mockSpace := &mockBodiesSpace{
				QueueForRemovalFunc: func(_ body.Collidable) {
					queueForRemovalCalled++
				},
				tilemapProvider: &mockTilemapDimensionsProvider{
					width:  tt.tilemapWidth,
					height: tt.tilemapHeight,
				},
			}

			mgr := NewManager(mockSpace)
			cfg := ProjectileConfig{Width: 2, Height: 1, Damage: 10}
			mgr.Spawn(cfg, tt.initialX16, tt.initialY16, tt.vx16, tt.vy16, nil)

			mgr.Update()

			if tt.expectRemoved && len(mgr.projectiles) != 0 {
				t.Errorf("expected projectile removed, got %d projectiles", len(mgr.projectiles))
			}
			if tt.expectRemoved && queueForRemovalCalled == 0 {
				t.Error("expected QueueForRemoval called")
			}
		})
	}
}

func TestManager_Clear(t *testing.T) {
	tests := []struct {
		name            string
		projectileCount int
	}{
		{
			name:            "clear removes all projectiles",
			projectileCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			removeBodyCalled := 0
			mockSpace := &mockBodiesSpace{
				RemoveBodyFunc: func(_ body.Collidable) {
					removeBodyCalled++
				},
			}

			mgr := NewManager(mockSpace)
			cfg := ProjectileConfig{Width: 2, Height: 1, Damage: 10}
			for i := 0; i < tt.projectileCount; i++ {
				mgr.Spawn(cfg, 100<<16, 50<<16, 5<<16, 0, nil)
			}

			mgr.Clear()

			if len(mgr.projectiles) != 0 {
				t.Errorf("expected 0 projectiles after Clear, got %d", len(mgr.projectiles))
			}
			if removeBodyCalled != tt.projectileCount {
				t.Errorf("expected RemoveBody called %d times, got %d", tt.projectileCount, removeBodyCalled)
			}
		})
	}
}

func TestManager_Draw(t *testing.T) {
	mockSpace := &mockBodiesSpace{}
	mgr := NewManager(mockSpace)
	screen := ebiten.NewImage(320, 240)

	mgr.Draw(screen)
}
