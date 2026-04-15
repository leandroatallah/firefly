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
			x16:      100 << 4,
			y16:      50 << 4,
			vx16:     5 << 4,
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

func TestManager_SpawnProjectile(t *testing.T) {
	mockSpace := &mockBodiesSpace{}
	mgr := NewManager(mockSpace)
	mgr.SpawnProjectile("bullet", 100<<4, 50<<4, 5<<4, 0, "player")

	if len(mgr.projectiles) != 1 {
		t.Errorf("expected 1 projectile, got %d", len(mgr.projectiles))
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
			initialX16:    95 << 4,
			initialY16:    50 << 4,
			vx16:          10 << 4,
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
				mgr.Spawn(cfg, 100<<4, 50<<4, 5<<4, 0, nil)
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
	mgr.Spawn(ProjectileConfig{Width: 2, Height: 1}, 100<<4, 50<<4, 5<<4, 0, nil)
	screen := ebiten.NewImage(320, 240)

	mgr.Draw(screen)
}

func TestManager_DrawWithOffset(t *testing.T) {
	mockSpace := &mockBodiesSpace{}
	mgr := NewManager(mockSpace)
	mgr.Spawn(ProjectileConfig{Width: 2, Height: 1}, 100<<4, 50<<4, 5<<4, 0, nil)
	screen := ebiten.NewImage(320, 240)

	mgr.DrawWithOffset(screen, 10, 10)
}

func TestManager_SetVFXManager(t *testing.T) {
	tests := []struct {
		name       string
		setManager bool
	}{
		{name: "stores non-nil VFX manager", setManager: true},
		{name: "manager remains nil when setter not called", setManager: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSpace := &mockBodiesSpace{}
			mgr := NewManager(mockSpace)
			vfxMgr := &mockVFXManager{}

			if tt.setManager {
				mgr.SetVFXManager(vfxMgr)
			}

			if tt.setManager && mgr.vfxManager != vfxMgr {
				t.Error("expected vfxManager to be set")
			}
			if !tt.setManager && mgr.vfxManager != nil {
				t.Error("expected vfxManager to remain nil")
			}
		})
	}
}

func TestManager_Spawn_ForwardsVFX(t *testing.T) {
	tests := []struct {
		name              string
		setVFX            bool
		wantImpactEffect  string
		wantDespawnEffect string
	}{
		{
			name:              "projectile receives VFX manager and default effects",
			setVFX:            true,
			wantImpactEffect:  "bullet_impact",
			wantDespawnEffect: "bullet_despawn",
		},
		{
			name:              "projectile receives nil VFX manager when not set",
			setVFX:            false,
			wantImpactEffect:  "bullet_impact",
			wantDespawnEffect: "bullet_despawn",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSpace := &mockBodiesSpace{}
			mgr := NewManager(mockSpace)
			vfxMgr := &mockVFXManager{}

			if tt.setVFX {
				mgr.SetVFXManager(vfxMgr)
			}

			cfg := ProjectileConfig{Width: 2, Height: 1}
			mgr.Spawn(cfg, 100<<4, 50<<4, 5<<4, 0, nil)

			p := mgr.projectiles[0]
			if tt.setVFX && p.vfxManager != vfxMgr {
				t.Error("expected projectile to have VFX manager")
			}
			if !tt.setVFX && p.vfxManager != nil {
				t.Error("expected projectile VFX manager to be nil")
			}
			if p.impactEffect != tt.wantImpactEffect {
				t.Errorf("impactEffect = %q, want %q", p.impactEffect, tt.wantImpactEffect)
			}
			if p.despawnEffect != tt.wantDespawnEffect {
				t.Errorf("despawnEffect = %q, want %q", p.despawnEffect, tt.wantDespawnEffect)
			}
		})
	}
}

func TestProjectile_NilVFXManager_NoPanic(t *testing.T) {
	tests := []struct {
		name    string
		trigger string // "impact", "block", or "despawn"
	}{
		{name: "OnTouch with nil VFX manager", trigger: "impact"},
		{name: "OnBlock with nil VFX manager", trigger: "block"},
		{name: "Update out-of-bounds with nil VFX manager", trigger: "despawn"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSpace := &mockBodiesSpace{
				tilemapProvider: &mockTilemapDimensionsProvider{
					width: 100, height: 100,
				},
			}
			mgr := NewManager(mockSpace)
			// Do NOT call SetVFXManager -- vfxManager stays nil
			cfg := ProjectileConfig{Width: 2, Height: 1}
			mgr.Spawn(cfg, 95<<4, 50<<4, 10<<4, 0, nil)

			// Should not panic regardless of trigger
			switch tt.trigger {
			case "impact":
				mgr.projectiles[0].OnTouch(nil)
			case "block":
				mgr.projectiles[0].OnBlock(nil)
			case "despawn":
				mgr.Update() // projectile will go out of bounds
			}
		})
	}
}

func TestProjectile_VFX_SpawnPuffCalled(t *testing.T) {
	tests := []struct {
		name          string
		trigger       string
		wantTypeKey   string
		wantCallCount int
	}{
		{name: "SpawnPuff on OnTouch", trigger: "impact", wantTypeKey: "bullet_impact", wantCallCount: 1},
		{name: "SpawnPuff on OnBlock", trigger: "block", wantTypeKey: "bullet_impact", wantCallCount: 1},
		{name: "SpawnPuff on out-of-bounds", trigger: "despawn", wantTypeKey: "", wantCallCount: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vfxMgr := &mockVFXManager{}
			mockSpace := &mockBodiesSpace{
				tilemapProvider: &mockTilemapDimensionsProvider{
					width: 100, height: 100,
				},
			}
			mgr := NewManager(mockSpace)
			mgr.SetVFXManager(vfxMgr)
			cfg := ProjectileConfig{Width: 2, Height: 1}
			mgr.Spawn(cfg, 95<<4, 50<<4, 10<<4, 0, "owner")

			switch tt.trigger {
			case "impact":
				// Create a mock collidable that is not the owner
				other := &mockCollidable{id: "enemy"}
				mgr.projectiles[0].OnTouch(other)
			case "block":
				other := &mockCollidable{id: "wall"}
				mgr.projectiles[0].OnBlock(other)
			case "despawn":
				mgr.Update()
			}

			if vfxMgr.spawnPuffCallCount != tt.wantCallCount {
				t.Errorf("SpawnPuff call count = %d, want %d", vfxMgr.spawnPuffCallCount, tt.wantCallCount)
			}
			if vfxMgr.lastTypeKey != tt.wantTypeKey {
				t.Errorf("SpawnPuff typeKey = %q, want %q", vfxMgr.lastTypeKey, tt.wantTypeKey)
			}
		})
	}
}
