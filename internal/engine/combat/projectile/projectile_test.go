package projectile

import (
	"testing"

	contractsvfx "github.com/boilerplate/ebiten-template/internal/engine/contracts/vfx"
)

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
