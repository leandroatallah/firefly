package projectile

import (
	"testing"

	contractsvfx "github.com/boilerplate/ebiten-template/internal/engine/contracts/vfx"
)

// TestProjectile_LifetimeDespawn verifies the frame-counted lifetime behavior.
// After N updates equal to lifetimeFrames (when > 0), the projectile must
// queue itself for removal exactly once and optionally emit a despawn VFX
// at its last fp16 position. lifetimeFrames <= 0 means infinite.
func TestProjectile_LifetimeDespawn(t *testing.T) {
	const (
		tilemapW = 1000
		tilemapH = 1000
		spawnX16 = 100 << 4
		spawnY16 = 50 << 4
	)

	tests := []struct {
		name           string
		lifetimeFrames int
		despawnEffect  string
		vfxManagerNil  bool
		updates        int
		wantQueued     bool
		wantVFXCalls   int
		wantVFXKey     string
	}{
		{
			name:           "infinite lifetime never despawns",
			lifetimeFrames: 0,
			despawnEffect:  "bullet_despawn",
			vfxManagerNil:  false,
			updates:        100,
			wantQueued:     false,
			wantVFXCalls:   0,
			wantVFXKey:     "",
		},
		{
			name:           "lifetime expires and queues removal",
			lifetimeFrames: 3,
			despawnEffect:  "bullet_despawn",
			vfxManagerNil:  false,
			updates:        3,
			wantQueued:     true,
			wantVFXCalls:   1,
			wantVFXKey:     "bullet_despawn",
		},
		{
			name:           "lifetime expires without vfx manager",
			lifetimeFrames: 2,
			despawnEffect:  "bullet_despawn",
			vfxManagerNil:  true,
			updates:        2,
			wantQueued:     true,
			wantVFXCalls:   0,
			wantVFXKey:     "",
		},
		{
			name:           "lifetime expires with empty effect",
			lifetimeFrames: 2,
			despawnEffect:  "",
			vfxManagerNil:  false,
			updates:        2,
			wantQueued:     true,
			wantVFXCalls:   0,
			wantVFXKey:     "",
		},
		{
			name:           "lifetime not yet expired",
			lifetimeFrames: 5,
			despawnEffect:  "bullet_despawn",
			vfxManagerNil:  false,
			updates:        4,
			wantQueued:     false,
			wantVFXCalls:   0,
			wantVFXKey:     "",
		},
		{
			name:           "negative lifetime treated as infinite",
			lifetimeFrames: -1,
			despawnEffect:  "bullet_despawn",
			vfxManagerNil:  false,
			updates:        10,
			wantQueued:     false,
			wantVFXCalls:   0,
			wantVFXKey:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSpace := &mockBodiesSpace{
				tilemapProvider: &mockTilemapDimensionsProvider{
					width:  tilemapW,
					height: tilemapH,
				},
			}
			mgr := NewManager(mockSpace)
			var vfxMgr *mockVFXManager
			if !tt.vfxManagerNil {
				vfxMgr = &mockVFXManager{}
				mgr.SetVFXManager(vfxMgr)
			}

			cfg := ProjectileConfig{
				Width:          2,
				Height:         1,
				DespawnEffect:  tt.despawnEffect,
				LifetimeFrames: tt.lifetimeFrames,
			}
			// Zero velocity so the projectile stays at spawnX16/spawnY16 and
			// never hits the out-of-bounds fallback.
			mgr.Spawn(cfg, spawnX16, spawnY16, 0, 0, nil)
			if len(mgr.projectiles) != 1 {
				t.Fatalf("expected 1 projectile spawned, got %d", len(mgr.projectiles))
			}
			p := mgr.projectiles[0]

			for i := 0; i < tt.updates; i++ {
				p.Update()
			}

			// Queued-for-removal assertion.
			queued := false
			for _, b := range mockSpace.queuedForRemoval {
				if b == p.body {
					queued = true
					break
				}
			}
			if queued != tt.wantQueued {
				t.Errorf("queuedForRemoval contains body = %v, want %v", queued, tt.wantQueued)
			}

			// QueueForRemoval invoked at most once invariant.
			count := 0
			for _, b := range mockSpace.queuedForRemoval {
				if b == p.body {
					count++
				}
			}
			if count > 1 {
				t.Errorf("QueueForRemoval for body invoked %d times, want at most 1", count)
			}

			// VFX assertions.
			if tt.vfxManagerNil {
				return
			}
			if got := len(vfxMgr.spawnPuffCalls); got != tt.wantVFXCalls {
				t.Errorf("spawnPuffCalls length = %d, want %d", got, tt.wantVFXCalls)
			}
			if tt.wantVFXCalls == 1 && len(vfxMgr.spawnPuffCalls) == 1 {
				call := vfxMgr.spawnPuffCalls[0]
				if call.typeKey != tt.wantVFXKey {
					t.Errorf("typeKey = %q, want %q", call.typeKey, tt.wantVFXKey)
				}
				wantX := float64(spawnX16) / 16.0
				wantY := float64(spawnY16) / 16.0
				if call.x != wantX {
					t.Errorf("x = %v, want %v", call.x, wantX)
				}
				if call.y != wantY {
					t.Errorf("y = %v, want %v", call.y, wantY)
				}
				if call.count != 1 {
					t.Errorf("count = %d, want 1", call.count)
				}
				if call.randRange != 0.0 {
					t.Errorf("randRange = %v, want 0.0", call.randRange)
				}
			}
		})
	}
}

// TestProjectile_OOBHasNoVFX verifies that the out-of-bounds safety fallback
// is silent: the body is queued for removal but no VFX is emitted.
func TestProjectile_OOBHasNoVFX(t *testing.T) {
	tests := []struct {
		name          string
		tilemapWidth  int
		tilemapHeight int
		initialX16    int
		initialY16    int
		vx16, vy16    int
	}{
		{
			name:          "horizontal out-of-bounds emits no VFX",
			tilemapWidth:  100,
			tilemapHeight: 100,
			initialX16:    95 << 4,
			initialY16:    50 << 4,
			vx16:          20 << 4,
			vy16:          0,
		},
		{
			name:          "vertical out-of-bounds emits no VFX",
			tilemapWidth:  100,
			tilemapHeight: 100,
			initialX16:    50 << 4,
			initialY16:    95 << 4,
			vx16:          0,
			vy16:          20 << 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSpace := &mockBodiesSpace{
				tilemapProvider: &mockTilemapDimensionsProvider{
					width:  tt.tilemapWidth,
					height: tt.tilemapHeight,
				},
			}
			vfxMgr := &mockVFXManager{}
			mgr := NewManager(mockSpace)
			mgr.SetVFXManager(vfxMgr)

			cfg := ProjectileConfig{
				Width:          2,
				Height:         1,
				DespawnEffect:  "bullet_despawn",
				LifetimeFrames: 0, // infinite — only OOB path can despawn
			}
			mgr.Spawn(cfg, tt.initialX16, tt.initialY16, tt.vx16, tt.vy16, nil)
			p := mgr.projectiles[0]

			p.Update()

			queued := false
			for _, b := range mockSpace.queuedForRemoval {
				if b == p.body {
					queued = true
					break
				}
			}
			if !queued {
				t.Error("expected body queued for removal after going out of bounds")
			}
			if got := len(vfxMgr.spawnPuffCalls); got != 0 {
				t.Errorf("spawnPuffCalls length = %d, want 0 (OOB must be silent)", got)
			}
		})
	}
}

// TestProjectileConfig_LifetimeFrames_Default verifies the zero-value config
// yields infinite lifetime behavior when flowed through Manager.Spawn.
func TestProjectileConfig_LifetimeFrames_Default(t *testing.T) {
	tests := []struct {
		name    string
		cfg     ProjectileConfig
		updates int
	}{
		{
			name:    "zero-value config never despawns via lifetime",
			cfg:     ProjectileConfig{Width: 2, Height: 1},
			updates: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cfg.LifetimeFrames != 0 {
				t.Fatalf("zero-value LifetimeFrames = %d, want 0", tt.cfg.LifetimeFrames)
			}

			mockSpace := &mockBodiesSpace{
				tilemapProvider: &mockTilemapDimensionsProvider{
					width:  10000,
					height: 10000,
				},
			}
			vfxMgr := &mockVFXManager{}
			mgr := NewManager(mockSpace)
			mgr.SetVFXManager(vfxMgr)
			mgr.Spawn(tt.cfg, 100<<4, 100<<4, 0, 0, nil)
			p := mgr.projectiles[0]

			for i := 0; i < tt.updates; i++ {
				p.Update()
			}

			for _, b := range mockSpace.queuedForRemoval {
				if b == p.body {
					t.Error("zero-value config projectile was queued for removal; expected infinite lifetime")
				}
			}
			if got := len(vfxMgr.spawnPuffCalls); got != 0 {
				t.Errorf("spawnPuffCalls length = %d, want 0 for infinite lifetime", got)
			}
		})
	}
}

// TestManager_Spawn_PropagatesLifetime verifies Manager.Spawn propagates the
// configured LifetimeFrames to the projectile. Observed through behavior:
// after exactly LifetimeFrames updates the body must be queued for removal.
func TestManager_Spawn_PropagatesLifetime(t *testing.T) {
	tests := []struct {
		name              string
		lifetimeFrames    int
		updatesBeforeDone int
	}{
		{
			name:              "spawn with LifetimeFrames 7 despawns after 7 updates",
			lifetimeFrames:    7,
			updatesBeforeDone: 7,
		},
		{
			name:              "spawn with LifetimeFrames 1 despawns after 1 update",
			lifetimeFrames:    1,
			updatesBeforeDone: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSpace := &mockBodiesSpace{
				tilemapProvider: &mockTilemapDimensionsProvider{
					width:  10000,
					height: 10000,
				},
			}
			mgr := NewManager(mockSpace)
			cfg := ProjectileConfig{
				Width:          2,
				Height:         1,
				LifetimeFrames: tt.lifetimeFrames,
			}
			mgr.Spawn(cfg, 100<<4, 100<<4, 0, 0, nil)
			if len(mgr.projectiles) != 1 {
				t.Fatalf("expected 1 projectile, got %d", len(mgr.projectiles))
			}
			p := mgr.projectiles[0]

			// Not yet queued before lifetime reaches zero.
			for i := 0; i < tt.updatesBeforeDone-1; i++ {
				p.Update()
				for _, b := range mockSpace.queuedForRemoval {
					if b == p.body {
						t.Fatalf("body queued for removal after %d updates, expected only after %d", i+1, tt.updatesBeforeDone)
					}
				}
			}

			// Final update should queue the body for removal exactly once.
			p.Update()
			matches := 0
			for _, b := range mockSpace.queuedForRemoval {
				if b == p.body {
					matches++
				}
			}
			if matches != 1 {
				t.Errorf("QueueForRemoval count for body = %d after %d updates, want 1", matches, tt.updatesBeforeDone)
			}
		})
	}
}

// TestProjectile_ImpactVFX verifies that a projectile spawns its impact VFX at
// the projectile's world-space position (fp16 -> float) before being queued for
// removal on OnTouch and OnBlock, with nil-manager and empty-key guards.
func TestProjectile_ImpactVFX(t *testing.T) {
	const (
		touch = "touch"
		block = "block"
	)

	type trigger struct {
		kind    string // "touch" or "block"
		isOwner bool
	}

	tests := []struct {
		name               string
		useVFXManager      bool
		impactEffect       string
		trigger            trigger
		x16, y16           int
		wantSpawnCalled    bool
		wantTypeKey        string
		wantX, wantY       float64
		wantCount          int
		wantRandRange      float64
		wantQueuedRemovals int
	}{
		{
			name:               "OnTouch spawns impact VFX",
			useVFXManager:      true,
			impactEffect:       "bullet_impact",
			trigger:            trigger{kind: touch, isOwner: false},
			x16:                32,
			y16:                16,
			wantSpawnCalled:    true,
			wantTypeKey:        "bullet_impact",
			wantX:              2.0,
			wantY:              1.0,
			wantCount:          1,
			wantRandRange:      0.0,
			wantQueuedRemovals: 1,
		},
		{
			name:               "OnTouch with owner does nothing",
			useVFXManager:      true,
			impactEffect:       "bullet_impact",
			trigger:            trigger{kind: touch, isOwner: true},
			x16:                32,
			y16:                16,
			wantSpawnCalled:    false,
			wantQueuedRemovals: 0,
		},
		{
			name:               "OnBlock spawns impact VFX",
			useVFXManager:      true,
			impactEffect:       "bullet_impact",
			trigger:            trigger{kind: block, isOwner: false},
			x16:                48,
			y16:                24,
			wantSpawnCalled:    true,
			wantTypeKey:        "bullet_impact",
			wantX:              3.0,
			wantY:              1.5,
			wantCount:          1,
			wantRandRange:      0.0,
			wantQueuedRemovals: 1,
		},
		{
			name:               "nil manager is safe on OnTouch",
			useVFXManager:      false,
			impactEffect:       "bullet_impact",
			trigger:            trigger{kind: touch, isOwner: false},
			x16:                32,
			y16:                16,
			wantSpawnCalled:    false,
			wantQueuedRemovals: 1,
		},
		{
			name:               "empty effect type is safe on OnTouch",
			useVFXManager:      true,
			impactEffect:       "",
			trigger:            trigger{kind: touch, isOwner: false},
			x16:                32,
			y16:                16,
			wantSpawnCalled:    false,
			wantQueuedRemovals: 1,
		},
		{
			name:               "nil manager is safe on OnBlock",
			useVFXManager:      false,
			impactEffect:       "bullet_impact",
			trigger:            trigger{kind: block, isOwner: false},
			x16:                48,
			y16:                24,
			wantSpawnCalled:    false,
			wantQueuedRemovals: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ownerMarker := &struct{ name string }{name: "owner"}

			projBody := &mockCollidable{id: "projectile"}
			projBody.SetOwner(ownerMarker)
			projBody.SetPosition16(tt.x16, tt.y16)

			space := &mockBodiesSpace{}

			var vfxMgr contractsvfx.Manager
			var mockMgr *mockVFXManager
			if tt.useVFXManager {
				mockMgr = &mockVFXManager{}
				vfxMgr = mockMgr
			}

			p := &projectile{
				body:         projBody,
				space:        space,
				vfxManager:   vfxMgr,
				impactEffect: tt.impactEffect,
			}

			var other *mockCollidable
			if tt.trigger.isOwner {
				// the owner is the value returned by body.Owner(); to hit the
				// "other == owner" branch, pass a Collidable whose identity
				// equals the owner. Owner() returns interface{}, so we wrap.
				owner := &mockCollidable{id: "owner-collidable"}
				projBody.SetOwner(owner)
				other = owner
			} else {
				other = &mockCollidable{id: "target"}
			}

			switch tt.trigger.kind {
			case touch:
				p.OnTouch(other)
			case block:
				p.OnBlock(other)
			default:
				t.Fatalf("unknown trigger kind %q", tt.trigger.kind)
			}

			// SpawnPuff assertions.
			if tt.wantSpawnCalled {
				if mockMgr == nil {
					t.Fatalf("test misconfigured: wantSpawnCalled without a mock manager")
				}
				if mockMgr.spawnPuffCallCount != 1 {
					t.Fatalf("SpawnPuff call count = %d, want 1", mockMgr.spawnPuffCallCount)
				}
				if mockMgr.lastTypeKey != tt.wantTypeKey {
					t.Errorf("SpawnPuff typeKey = %q, want %q", mockMgr.lastTypeKey, tt.wantTypeKey)
				}
				if mockMgr.lastX != tt.wantX {
					t.Errorf("SpawnPuff x = %v, want %v", mockMgr.lastX, tt.wantX)
				}
				if mockMgr.lastY != tt.wantY {
					t.Errorf("SpawnPuff y = %v, want %v", mockMgr.lastY, tt.wantY)
				}
				if mockMgr.lastCount != tt.wantCount {
					t.Errorf("SpawnPuff count = %d, want %d", mockMgr.lastCount, tt.wantCount)
				}
				if mockMgr.lastRandRange != tt.wantRandRange {
					t.Errorf("SpawnPuff randRange = %v, want %v", mockMgr.lastRandRange, tt.wantRandRange)
				}
			} else if mockMgr != nil && mockMgr.spawnPuffCallCount != 0 {
				t.Errorf("SpawnPuff unexpectedly called %d times", mockMgr.spawnPuffCallCount)
			}

			// QueueForRemoval assertions.
			if got := len(space.queuedForRemoval); got != tt.wantQueuedRemovals {
				t.Errorf("QueueForRemoval count = %d, want %d", got, tt.wantQueuedRemovals)
			}
			if tt.wantQueuedRemovals == 1 && len(space.queuedForRemoval) == 1 {
				if space.queuedForRemoval[0] != projBody {
					t.Errorf("QueueForRemoval queued %v, want projectile body", space.queuedForRemoval[0])
				}
			}
		})
	}
}
