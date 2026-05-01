package projectile

import (
	"image"
	"testing"

	enginecombat "github.com/boilerplate/ebiten-template/internal/engine/combat"
	body "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/hajimehoshi/ebiten/v2"
)

// mockProjectileBody is a Collidable that ALSO implements body.Projectile.
// It models the wrapper described in SPEC §"Implementation Touchpoints".
type mockProjectileBody struct {
	id            string
	owner         interface{}
	x16, y16      int
	interceptable bool
}

func (m *mockProjectileBody) ID() string                 { return m.id }
func (m *mockProjectileBody) SetID(id string)            { m.id = id }
func (m *mockProjectileBody) Position() image.Rectangle  { return image.Rectangle{} }
func (m *mockProjectileBody) SetPosition(_, _ int)       {}
func (m *mockProjectileBody) SetPosition16(x, y int)     { m.x16 = x; m.y16 = y }
func (m *mockProjectileBody) SetSize(_, _ int)           {}
func (m *mockProjectileBody) Scale() float64             { return 1.0 }
func (m *mockProjectileBody) SetScale(_ float64)         {}
func (m *mockProjectileBody) GetPosition16() (int, int)  { return m.x16, m.y16 }
func (m *mockProjectileBody) GetPositionMin() (int, int) { return 0, 0 }
func (m *mockProjectileBody) GetShape() body.Shape       { return nil }
func (m *mockProjectileBody) Owner() interface{}         { return m.owner }
func (m *mockProjectileBody) SetOwner(o interface{})     { m.owner = o }
func (m *mockProjectileBody) LastOwner() interface{}     { return nil }
func (m *mockProjectileBody) OnTouch(_ body.Collidable)  {}
func (m *mockProjectileBody) OnBlock(_ body.Collidable)  {}
func (m *mockProjectileBody) GetTouchable() body.Touchable {
	return m
}
func (m *mockProjectileBody) DrawCollisionBox(_ *ebiten.Image, _ image.Rectangle) {}
func (m *mockProjectileBody) CollisionPosition() []image.Rectangle                { return nil }
func (m *mockProjectileBody) CollisionShapes() []body.Collidable                  { return nil }
func (m *mockProjectileBody) IsObstructive() bool                                 { return false }
func (m *mockProjectileBody) SetIsObstructive(_ bool)                             {}
func (m *mockProjectileBody) AddCollision(_ ...body.Collidable)                   {}
func (m *mockProjectileBody) ClearCollisions()                                    {}
func (m *mockProjectileBody) SetTouchable(_ body.Touchable)                       {}
func (m *mockProjectileBody) ApplyValidPosition(_ int, _ bool, _ body.BodiesSpace) (int, int, bool) {
	return 0, 0, false
}
func (m *mockProjectileBody) Interceptable() bool { return m.interceptable }

// Compile-time assertion: mockProjectileBody must implement body.Projectile.
var _ body.Projectile = (*mockProjectileBody)(nil)

// newProjectileForTest builds a *projectile wired for OnTouch/OnBlock test calls.
func newProjectileForTest(projBody body.Collidable, space body.BodiesSpace, vfxMgr *mockVFXManager, damage int, faction enginecombat.Faction) *projectile {
	p := &projectile{
		body:         projBody,
		space:        space,
		damage:       damage,
		faction:      faction,
		impactEffect: "bullet_impact",
	}
	if vfxMgr != nil {
		p.vfxManager = vfxMgr
	}
	return p
}

// AC1 — projectile-vs-default-projectile no-op on OnTouch and OnBlock.
func TestProjectile_OnTouch_IgnoresOtherDefaultProjectile(t *testing.T) {
	tests := []struct {
		name    string
		trigger string // "touch" or "block"
	}{
		{name: "OnTouch ignores other default projectile", trigger: "touch"},
		{name: "OnBlock ignores other default projectile", trigger: "block"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			space := &mockBodiesSpace{}
			vfx := &mockVFXManager{}

			selfBody := &mockProjectileBody{id: "p_self", interceptable: false}
			otherBody := &mockProjectileBody{id: "p_other", interceptable: false}

			p := newProjectileForTest(selfBody, space, vfx, 5, enginecombat.FactionPlayer)

			switch tt.trigger {
			case "touch":
				p.OnTouch(otherBody)
			case "block":
				p.OnBlock(otherBody)
			}

			if got := len(space.queuedForRemoval); got != 0 {
				t.Errorf("QueueForRemoval count = %d, want 0 (no projectile-vs-projectile despawn)", got)
			}
			if got := vfx.spawnPuffCallCount; got != 0 {
				t.Errorf("SpawnPuff call count = %d, want 0 (no impact VFX)", got)
			}
		})
	}
}

// AC2 — projectile vs. Actor body still applies damage and despawns.
func TestProjectile_OnTouch_HitsActorBody_AppliesDamageAndDespawns(t *testing.T) {
	space := &mockBodiesSpace{}
	vfx := &mockVFXManager{}

	selfBody := &mockProjectileBody{id: "p_self", interceptable: false}

	enemyDamageable := &fakeDamageable{faction: enginecombat.FactionEnemy}
	actorBody := fakeCollidableWithOwner(enemyDamageable)
	actorBody.SetPosition16(64, 32)
	selfBody.SetPosition16(64, 32)

	p := newProjectileForTest(selfBody, space, vfx, 5, enginecombat.FactionPlayer)

	p.OnTouch(actorBody)

	if got := len(enemyDamageable.takeDamageCalls); got != 1 {
		t.Fatalf("TakeDamage call count = %d, want 1", got)
	}
	if enemyDamageable.takeDamageCalls[0] != 5 {
		t.Errorf("TakeDamage amount = %d, want 5", enemyDamageable.takeDamageCalls[0])
	}
	if got := len(space.queuedForRemoval); got != 1 {
		t.Errorf("QueueForRemoval count = %d, want 1", got)
	}
	if got := vfx.spawnPuffCallCount; got != 1 {
		t.Errorf("SpawnPuff count = %d, want 1", got)
	}
}

// AC3 — non-projectile body with Damageable owner still triggers damage and despawn.
func TestProjectile_OnTouch_DoesNotShortCircuitOnNonProjectileBody(t *testing.T) {
	space := &mockBodiesSpace{}
	vfx := &mockVFXManager{}

	selfBody := &mockProjectileBody{id: "p_self", interceptable: false}

	dmg := &fakeDamageable{faction: enginecombat.FactionEnemy}
	meleeBody := fakeCollidableWithOwner(dmg) // *mockCollidable: NOT a body.Projectile

	p := newProjectileForTest(selfBody, space, vfx, 7, enginecombat.FactionPlayer)

	p.OnTouch(meleeBody)

	if got := len(dmg.takeDamageCalls); got != 1 {
		t.Errorf("TakeDamage call count = %d, want 1 (non-projectile must not be filtered)", got)
	}
	if got := len(space.queuedForRemoval); got != 1 {
		t.Errorf("QueueForRemoval count = %d, want 1", got)
	}
	if got := vfx.spawnPuffCallCount; got != 1 {
		t.Errorf("SpawnPuff count = %d, want 1", got)
	}
}

// AC4 — deterministic across ordering: invoking either side or both in any
// order produces identical (zero-effect) outcomes.
func TestProjectile_PvP_DeterministicAcrossOrdering(t *testing.T) {
	tests := []struct {
		name  string
		order []int // sequence of which projectile invokes OnTouch on the other (1->2 or 2->1)
	}{
		{name: "p1 then p2", order: []int{12, 21}},
		{name: "p2 then p1", order: []int{21, 12}},
		{name: "only p1", order: []int{12}},
		{name: "only p2", order: []int{21}},
		{name: "p1 twice", order: []int{12, 12}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			space := &mockBodiesSpace{}
			vfx := &mockVFXManager{}

			body1 := &mockProjectileBody{id: "p1", interceptable: false}
			body2 := &mockProjectileBody{id: "p2", interceptable: false}

			p1 := newProjectileForTest(body1, space, vfx, 5, enginecombat.FactionPlayer)
			p2 := newProjectileForTest(body2, space, vfx, 5, enginecombat.FactionEnemy)

			for _, step := range tt.order {
				switch step {
				case 12:
					p1.OnTouch(body2)
				case 21:
					p2.OnTouch(body1)
				}
			}

			if got := len(space.queuedForRemoval); got != 0 {
				t.Errorf("QueueForRemoval count = %d, want 0", got)
			}
			if got := vfx.spawnPuffCallCount; got != 0 {
				t.Errorf("SpawnPuff count = %d, want 0", got)
			}
		})
	}
}

// AC5 — opt-in interceptability: an interceptable target IS hit by a default
// projectile (target-property semantic).
func TestProjectile_OnTouch_InterceptableTargetIsHit(t *testing.T) {
	space := &mockBodiesSpace{}
	vfx := &mockVFXManager{}

	selfBody := &mockProjectileBody{id: "shooter", interceptable: false}

	rocketDamageable := &fakeDamageable{faction: enginecombat.FactionEnemy}
	targetBody := &mockProjectileBody{id: "rocket", interceptable: true}
	targetBody.SetOwner(rocketDamageable)

	p := newProjectileForTest(selfBody, space, vfx, 3, enginecombat.FactionPlayer)

	p.OnTouch(targetBody)

	if got := len(rocketDamageable.takeDamageCalls); got != 1 {
		t.Errorf("TakeDamage call count = %d, want 1 (interceptable target must take damage)", got)
	}
	if got := len(space.queuedForRemoval); got != 1 {
		t.Errorf("QueueForRemoval count = %d, want 1 (shooter despawns on interceptable hit)", got)
	}
	if space.queuedForRemoval[0] != selfBody {
		t.Errorf("QueueForRemoval queued %v, want shooter body %v", space.queuedForRemoval[0], selfBody)
	}
}

// AC5 — symmetry: an interceptable shooter still ignores default-projectile
// targets (interceptability is a property of the *target*, not the shooter).
func TestProjectile_OnTouch_InterceptableShooterStillIgnoresDefaultBullet(t *testing.T) {
	space := &mockBodiesSpace{}
	vfx := &mockVFXManager{}

	selfBody := &mockProjectileBody{id: "rocket", interceptable: true}

	otherDamageable := &fakeDamageable{faction: enginecombat.FactionEnemy}
	defaultBullet := &mockProjectileBody{id: "bullet", interceptable: false}
	defaultBullet.SetOwner(otherDamageable)

	p := newProjectileForTest(selfBody, space, vfx, 9, enginecombat.FactionPlayer)

	p.OnTouch(defaultBullet)

	if got := len(otherDamageable.takeDamageCalls); got != 0 {
		t.Errorf("TakeDamage call count = %d, want 0 (default bullet must not be hit even by interceptable shooter)", got)
	}
	if got := len(space.queuedForRemoval); got != 0 {
		t.Errorf("QueueForRemoval count = %d, want 0 (interceptable shooter must survive)", got)
	}
}
