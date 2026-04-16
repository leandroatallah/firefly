package projectile

import (
	"testing"

	enginecombat "github.com/boilerplate/ebiten-template/internal/engine/combat"
	contractsbody "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
)

// TestProjectile_AppliesDamageOnHit covers AC2, AC3, AC7, AC8 and AC10 by
// exercising applyDamage resolution paths (direct Damageable body, Damageable
// via Owner(), non-damageable, self-owner), the faction gate (same-faction
// ignored, neutral always hurts), and the zero-damage no-op guard.
//
// The projectile must still despawn (QueueForRemoval called exactly once) on
// any non-self-owner hit, whether or not damage was applied.
func TestProjectile_AppliesDamageOnHit(t *testing.T) {
	const (
		projectileDamage = 10
	)

	tests := []struct {
		name string
		// Projectile side.
		projFaction enginecombat.Faction
		projDamage  int
		// Target factory: returns (otherArg, expectedTakeDamageTarget, isSelfOwner).
		// expectedTakeDamageTarget is nil when no TakeDamage call is expected.
		buildTarget func(projectileBody *mockCollidable) (other *mockCollidable, damageable *fakeDamageable, bodyDamageable *fakeDamageableBody, selfOwner bool)
		// Expectations.
		wantTakeDamageCalls int
		wantDamageAmount    int
		wantQueuedRemovals  int
	}{
		{
			name:        "hit on owner that is Damageable (different faction)",
			projFaction: enginecombat.FactionPlayer,
			projDamage:  projectileDamage,
			buildTarget: func(_ *mockCollidable) (*mockCollidable, *fakeDamageable, *fakeDamageableBody, bool) {
				d := &fakeDamageable{faction: enginecombat.FactionEnemy}
				return fakeCollidableWithOwner(d), d, nil, false
			},
			wantTakeDamageCalls: 1,
			wantDamageAmount:    projectileDamage,
			wantQueuedRemovals:  1,
		},
		{
			name:        "hit on body that directly implements Damageable",
			projFaction: enginecombat.FactionPlayer,
			projDamage:  projectileDamage,
			buildTarget: func(_ *mockCollidable) (*mockCollidable, *fakeDamageable, *fakeDamageableBody, bool) {
				// Body itself implements Damageable; we return it via the body
				// damageable slot. applyDamage MUST resolve the body directly
				// (step 1) and NOT consult Owner().
				_ = &fakeDamageableBody{
					id:      "direct-body",
					faction: enginecombat.FactionEnemy,
				}
				// NOTE: because OnTouch takes body.Collidable and our test uses
				// mockCollidable as the `other` argument, we simulate the
				// "body directly implements Damageable" case by attaching a
				// *fakeDamageableBody as the owner of a thin Collidable shim,
				// but the SPEC requires resolution step 1 to check the body
				// itself. To exercise step 1 we need `other` to BE the
				// Damageable. We signal this to the test harness via the
				// bodyDamageable return; the harness passes it as the `other`.
				body := &fakeDamageableBody{
					id:      "direct-body",
					faction: enginecombat.FactionEnemy,
				}
				return nil, nil, body, false
			},
			wantTakeDamageCalls: 1,
			wantDamageAmount:    projectileDamage,
			wantQueuedRemovals:  1,
		},
		{
			name:        "hit on non-damageable body (no panic, still despawns)",
			projFaction: enginecombat.FactionPlayer,
			projDamage:  projectileDamage,
			buildTarget: func(_ *mockCollidable) (*mockCollidable, *fakeDamageable, *fakeDamageableBody, bool) {
				return &mockCollidable{id: "inert"}, nil, nil, false
			},
			wantTakeDamageCalls: 0,
			wantQueuedRemovals:  1,
		},
		{
			name:        "same-faction hit ignored (no damage, still despawns)",
			projFaction: enginecombat.FactionPlayer,
			projDamage:  projectileDamage,
			buildTarget: func(_ *mockCollidable) (*mockCollidable, *fakeDamageable, *fakeDamageableBody, bool) {
				d := &fakeDamageable{faction: enginecombat.FactionPlayer}
				return fakeCollidableWithOwner(d), d, nil, false
			},
			wantTakeDamageCalls: 0,
			wantQueuedRemovals:  1,
		},
		{
			name:        "neutral projectile hurts player target",
			projFaction: enginecombat.FactionNeutral,
			projDamage:  projectileDamage,
			buildTarget: func(_ *mockCollidable) (*mockCollidable, *fakeDamageable, *fakeDamageableBody, bool) {
				d := &fakeDamageable{faction: enginecombat.FactionPlayer}
				return fakeCollidableWithOwner(d), d, nil, false
			},
			wantTakeDamageCalls: 1,
			wantDamageAmount:    projectileDamage,
			wantQueuedRemovals:  1,
		},
		{
			name:        "neutral target hurt by player projectile",
			projFaction: enginecombat.FactionPlayer,
			projDamage:  projectileDamage,
			buildTarget: func(_ *mockCollidable) (*mockCollidable, *fakeDamageable, *fakeDamageableBody, bool) {
				// Target has no explicit faction (defaults to Neutral).
				d := &fakeDamageable{faction: enginecombat.FactionNeutral}
				return fakeCollidableWithOwner(d), d, nil, false
			},
			wantTakeDamageCalls: 1,
			wantDamageAmount:    projectileDamage,
			wantQueuedRemovals:  1,
		},
		{
			name:        "zero damage no-op (no TakeDamage, still despawns)",
			projFaction: enginecombat.FactionPlayer,
			projDamage:  0,
			buildTarget: func(_ *mockCollidable) (*mockCollidable, *fakeDamageable, *fakeDamageableBody, bool) {
				d := &fakeDamageable{faction: enginecombat.FactionEnemy}
				return fakeCollidableWithOwner(d), d, nil, false
			},
			wantTakeDamageCalls: 0,
			wantQueuedRemovals:  1,
		},
		{
			name:        "self owner short-circuit (no damage, NOT queued)",
			projFaction: enginecombat.FactionPlayer,
			projDamage:  projectileDamage,
			buildTarget: func(projBody *mockCollidable) (*mockCollidable, *fakeDamageable, *fakeDamageableBody, bool) {
				// The projectile owner is a Damageable of a different faction,
				// and `other` equals the owner. Existing OnTouch path must
				// short-circuit BEFORE applyDamage runs.
				owner := &mockCollidable{id: "owner-collidable"}
				projBody.SetOwner(owner)
				return owner, nil, nil, true
			},
			wantTakeDamageCalls: 0,
			wantQueuedRemovals:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projBody := &mockCollidable{id: "projectile"}
			// Default owner: a sentinel distinct from any `other`.
			ownerSentinel := &struct{ tag string }{tag: "proj-owner"}
			projBody.SetOwner(ownerSentinel)

			other, damageable, bodyDamageable, _ := tt.buildTarget(projBody)

			space := &mockBodiesSpace{}

			p := &projectile{
				body:    projBody,
				space:   space,
				damage:  tt.projDamage,
				faction: tt.projFaction,
			}

			// Dispatch: if bodyDamageable is set, use it as `other` (tests the
			// step-1 resolution path where the body itself is Damageable).
			if bodyDamageable != nil {
				p.OnTouch(bodyDamageable)
			} else {
				p.OnTouch(other)
			}

			// Damage call assertions.
			switch {
			case bodyDamageable != nil:
				if got := len(bodyDamageable.takeDamageCalls); got != tt.wantTakeDamageCalls {
					t.Errorf("body TakeDamage call count = %d, want %d", got, tt.wantTakeDamageCalls)
				}
				if tt.wantTakeDamageCalls == 1 && len(bodyDamageable.takeDamageCalls) == 1 {
					if got := bodyDamageable.takeDamageCalls[0]; got != tt.wantDamageAmount {
						t.Errorf("body TakeDamage amount = %d, want %d", got, tt.wantDamageAmount)
					}
				}
			case damageable != nil:
				if got := len(damageable.takeDamageCalls); got != tt.wantTakeDamageCalls {
					t.Errorf("owner TakeDamage call count = %d, want %d", got, tt.wantTakeDamageCalls)
				}
				if tt.wantTakeDamageCalls == 1 && len(damageable.takeDamageCalls) == 1 {
					if got := damageable.takeDamageCalls[0]; got != tt.wantDamageAmount {
						t.Errorf("owner TakeDamage amount = %d, want %d", got, tt.wantDamageAmount)
					}
				}
			default:
				if tt.wantTakeDamageCalls != 0 {
					t.Fatalf("test misconfigured: expected TakeDamage calls but no Damageable target built")
				}
			}

			// Queue-for-removal assertion.
			if got := len(space.queuedForRemoval); got != tt.wantQueuedRemovals {
				t.Errorf("QueueForRemoval count = %d, want %d", got, tt.wantQueuedRemovals)
			}
		})
	}
}

// TestProjectile_AppliesDamageOnBlock covers AC2 for the OnBlock path: blocking
// hits must also resolve Damageable and call TakeDamage (per SPEC 2.4).
func TestProjectile_AppliesDamageOnBlock(t *testing.T) {
	projBody := &mockCollidable{id: "projectile"}
	projBody.SetOwner(&struct{}{})

	d := &fakeDamageable{faction: enginecombat.FactionEnemy}
	other := fakeCollidableWithOwner(d)

	space := &mockBodiesSpace{}
	p := &projectile{
		body:    projBody,
		space:   space,
		damage:  7,
		faction: enginecombat.FactionPlayer,
	}

	p.OnBlock(other)

	if got := len(d.takeDamageCalls); got != 1 {
		t.Fatalf("TakeDamage call count = %d, want 1", got)
	}
	if got := d.takeDamageCalls[0]; got != 7 {
		t.Errorf("TakeDamage amount = %d, want 7", got)
	}
	if got := len(space.queuedForRemoval); got != 1 {
		t.Errorf("QueueForRemoval count = %d, want 1", got)
	}
}

// TestProjectile_ResolvesDestructible covers AC7: an object implementing
// contracts/combat.Destructible (Damageable + IsDestroyed) is treated
// identically by the projectile hit path — no special-casing.
func TestProjectile_ResolvesDestructible(t *testing.T) {
	projBody := &mockCollidable{id: "projectile"}
	projBody.SetOwner(&struct{}{})

	dest := &fakeDestructible{faction: enginecombat.FactionNeutral}
	other := fakeCollidableWithOwner(dest)

	space := &mockBodiesSpace{}
	p := &projectile{
		body:    projBody,
		space:   space,
		damage:  15,
		faction: enginecombat.FactionPlayer,
	}

	p.OnTouch(other)

	if got := len(dest.takeDamageCalls); got != 1 {
		t.Fatalf("Destructible TakeDamage call count = %d, want 1", got)
	}
	if got := dest.takeDamageCalls[0]; got != 15 {
		t.Errorf("Destructible TakeDamage amount = %d, want 15", got)
	}
	if got := len(space.queuedForRemoval); got != 1 {
		t.Errorf("QueueForRemoval count = %d, want 1", got)
	}
}

// TestProjectile_NilOtherSafe covers AC2 post-condition "never panics" when the
// resolution inputs are degenerate. A nil `other` must be tolerated by the
// applyDamage helper without panicking. We exercise OnBlock which has no
// self-owner short-circuit and forwards directly.
func TestProjectile_NilOtherSafe(t *testing.T) {
	projBody := &mockCollidable{id: "projectile"}
	projBody.SetOwner(&struct{}{})

	space := &mockBodiesSpace{}
	p := &projectile{
		body:    projBody,
		space:   space,
		damage:  5,
		faction: enginecombat.FactionPlayer,
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("OnBlock(nil) panicked: %v", r)
		}
	}()
	p.OnBlock(nil)
}

// TestProjectile_PassthroughBodiesIgnored covers the requirement that projectiles
// pass through items/power-ups without being destroyed or dealing damage.
func TestProjectile_PassthroughBodiesIgnored(t *testing.T) {
	tests := []struct {
		name  string
		other func() contractsbody.Collidable
	}{
		{
			name: "passthrough body directly",
			other: func() contractsbody.Collidable {
				return &mockPassthroughBody{mockCollidable{id: "item"}}
			},
		},
		{
			name: "body whose owner is passthrough",
			other: func() contractsbody.Collidable {
				return fakeCollidableWithOwner(&mockPassthroughOwner{})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projBody := &mockCollidable{id: "projectile"}
			projBody.SetOwner(&struct{}{})

			space := &mockBodiesSpace{}
			p := &projectile{
				body:    projBody,
				space:   space,
				damage:  10,
				faction: enginecombat.FactionPlayer,
			}

			other := tt.other()

			p.OnTouch(other)
			if got := len(space.queuedForRemoval); got != 0 {
				t.Errorf("OnTouch: QueueForRemoval count = %d, want 0 (passthrough)", got)
			}

			p.OnBlock(other)
			if got := len(space.queuedForRemoval); got != 0 {
				t.Errorf("OnBlock: QueueForRemoval count = %d, want 0 (passthrough)", got)
			}
		})
	}
}
