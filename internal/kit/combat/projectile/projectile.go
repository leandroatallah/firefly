// internal/engine/combat/projectile/projectile.go
package projectile

import (
	contractsbody "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	contractscombat "github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
	contractsvfx "github.com/boilerplate/ebiten-template/internal/engine/contracts/vfx"
	enginecombat "github.com/boilerplate/ebiten-template/internal/kit/combat"
)

// factioned is a file-local interface for entities that expose their faction.
type factioned interface {
	Faction() enginecombat.Faction
}

// idable is a file-local interface for entities that have an ID.
type idable interface {
	ID() string
}

// isPassthrough returns true if the body or its owner implements Passthrough.
func isPassthrough(other contractsbody.Collidable) bool {
	if other == nil {
		return false
	}
	if pt, ok := other.(contractsbody.Passthrough); ok && pt.IsPassthrough() {
		return true
	}
	if owner := other.Owner(); owner != nil {
		if pt, ok := owner.(contractsbody.Passthrough); ok && pt.IsPassthrough() {
			return true
		}
	}
	return false
}

// isProjectile checks whether other (or its owner) implements body.Projectile.
// Returns (true, interceptable) when it does, (false, false) otherwise.
func isProjectile(other contractsbody.Collidable) (bool, bool) {
	if other == nil {
		return false, false
	}
	if proj, ok := other.(contractsbody.Projectile); ok {
		return true, proj.Interceptable()
	}
	if owner := other.Owner(); owner != nil {
		if proj, ok := owner.(contractsbody.Projectile); ok {
			return true, proj.Interceptable()
		}
	}
	return false, false
}

// projectileBody wraps a contractsbody.Collidable and adds Interceptable so the
// trait is discoverable directly on the body registered in the physics space.
type projectileBody struct {
	contractsbody.Collidable
	interceptable bool
}

func (pb *projectileBody) Interceptable() bool { return pb.interceptable }

// projectile is the internal state of a spawned projectile.
type projectile struct {
	movable         contractsbody.Movable
	body            contractsbody.Collidable
	space           contractsbody.BodiesSpace
	speedX16        int
	speedY16        int
	vfxManager      contractsvfx.Manager
	impactEffect    string
	despawnEffect   string
	lifetimeFrames  int // configured total lifetime (0 = infinite)
	currentLifetime int // frames remaining; only meaningful when lifetimeFrames > 0
	damage          int
	faction         enginecombat.Faction
	interceptable   bool
}

func (p *projectile) Interceptable() bool { return p.interceptable }

func (p *projectile) Update() {
	x, y := p.body.GetPosition16()
	x += p.speedX16
	y += p.speedY16
	p.body.SetPosition16(x, y)

	p.space.ResolveCollisions(p.body)

	// Lifetime tick: if lifetimeFrames > 0, decrement and despawn when expired.
	if p.lifetimeFrames > 0 {
		p.currentLifetime--
		if p.currentLifetime <= 0 {
			p.spawnVFX(p.despawnEffect)
			p.space.QueueForRemoval(p.body)
			return
		}
	}

	provider := p.space.GetTilemapDimensionsProvider()
	if provider == nil {
		return
	}
	w := provider.GetTilemapWidth()
	h := provider.GetTilemapHeight()

	// Bounds are in fp16 units (scale factor 16: 1 pixel = 16 units)
	if x < 0 || y < 0 || x > w<<4 || y > h<<4 {
		p.space.QueueForRemoval(p.body)
	}
}

func (p *projectile) OnTouch(other contractsbody.Collidable) {
	if p.isOwner(other) {
		return
	}
	if isPassthrough(other) {
		return
	}
	if isProj, interceptable := isProjectile(other); isProj && !interceptable {
		return
	}
	p.applyDamage(other)
	p.spawnVFX(p.impactEffect)
	p.space.QueueForRemoval(p.body)
}

func (p *projectile) OnBlock(other contractsbody.Collidable) {
	if p.isOwner(other) {
		return
	}
	if isPassthrough(other) {
		return
	}
	if isProj, interceptable := isProjectile(other); isProj && !interceptable {
		return
	}
	p.applyDamage(other)
	p.spawnVFX(p.impactEffect)
	p.space.QueueForRemoval(p.body)
}

// isOwner returns true if the other body is the projectile's owner or belongs to it.
func (p *projectile) isOwner(other contractsbody.Collidable) bool {
	if other == nil {
		return false
	}
	owner := p.body.Owner()
	if owner == nil {
		return false
	}

	// Direct equality check
	if other == owner {
		return true
	}

	// Check if other's owner is our owner
	otherOwner := other.Owner()
	if otherOwner != nil && otherOwner == owner {
		return true
	}

	// Robust check via ID if available
	if ownerID, ok := owner.(idable); ok {
		targetID := ""
		if otherID, ok := other.(idable); ok {
			targetID = otherID.ID()
		} else if otherOwner != nil {
			if ooID, ok := otherOwner.(idable); ok {
				targetID = ooID.ID()
			}
		}

		if targetID != "" && targetID == ownerID.ID() {
			return true
		}
	}

	return false
}

// applyDamage resolves a Damageable from the hit body and calls TakeDamage,
// honouring faction and zero-damage guards. Safe on nil / non-damageable others.
func (p *projectile) applyDamage(other contractsbody.Collidable) {
	if p.damage == 0 {
		return
	}
	if other == nil {
		return
	}

	target, tFaction, ok := p.resolveDamageable(other)
	if !ok {
		return
	}

	// Faction gate: skip only when both sides are non-neutral AND equal.
	if p.faction != enginecombat.FactionNeutral &&
		tFaction != enginecombat.FactionNeutral &&
		p.faction == tFaction {
		return
	}

	target.TakeDamage(p.damage)
}

// resolveDamageable tries (1) the body itself, then (2) body.Owner().
// Returns the Damageable, the target's faction (FactionNeutral when not factioned),
// and whether resolution succeeded.
func (p *projectile) resolveDamageable(other contractsbody.Collidable) (contractscombat.Damageable, enginecombat.Faction, bool) {
	// Step 1: body itself.
	if d, ok := other.(contractscombat.Damageable); ok {
		f := enginecombat.FactionNeutral
		if fac, ok := other.(factioned); ok {
			f = fac.Faction()
		}
		return d, f, true
	}

	// Step 2: body.Owner().
	owner := other.Owner()
	if owner == nil {
		return nil, enginecombat.FactionNeutral, false
	}
	if d, ok := owner.(contractscombat.Damageable); ok {
		f := enginecombat.FactionNeutral
		if fac, ok := owner.(factioned); ok {
			f = fac.Faction()
		}
		return d, f, true
	}

	return nil, enginecombat.FactionNeutral, false
}

func (p *projectile) spawnVFX(typeKey string) {
	if p.vfxManager == nil || typeKey == "" {
		return
	}
	x16, y16 := p.body.GetPosition16()
	p.vfxManager.SpawnPuff(typeKey, float64(x16)/16.0, float64(y16)/16.0, 1, 0.0)
}
