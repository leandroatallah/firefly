// internal/engine/combat/projectile/projectile.go
package projectile

import (
	contractsbody "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	contractsvfx "github.com/boilerplate/ebiten-template/internal/engine/contracts/vfx"
)

// projectile is the internal state of a spawned projectile.
type projectile struct {
	movable       contractsbody.Movable
	body          contractsbody.Collidable
	space         contractsbody.BodiesSpace
	speedX16      int
	speedY16      int
	vfxManager    contractsvfx.Manager
	impactEffect  string
	despawnEffect string
}

func (p *projectile) Update() {
	x, y := p.body.GetPosition16()
	x += p.speedX16
	y += p.speedY16
	p.body.SetPosition16(x, y)

	p.space.ResolveCollisions(p.body)

	provider := p.space.GetTilemapDimensionsProvider()
	if provider == nil {
		return
	}
	w := provider.GetTilemapWidth()
	h := provider.GetTilemapHeight()

	// Bounds are in fp16 units (scale factor 16: 1 pixel = 16 units)
	if x < 0 || y < 0 || x > w<<4 || y > h<<4 {
		p.spawnVFX(p.despawnEffect)
		p.space.QueueForRemoval(p.body)
	}
}

func (p *projectile) OnTouch(other contractsbody.Collidable) {
	if other != p.body.Owner() {
		p.spawnVFX(p.impactEffect)
		p.space.QueueForRemoval(p.body)
	}
}

func (p *projectile) OnBlock(other contractsbody.Collidable) {
	p.spawnVFX(p.impactEffect)
	p.space.QueueForRemoval(p.body)
}

func (p *projectile) spawnVFX(typeKey string) {
	if p.vfxManager == nil || typeKey == "" {
		return
	}
	x16, y16 := p.body.GetPosition16()
	p.vfxManager.SpawnPuff(typeKey, float64(x16)/16.0, float64(y16)/16.0, 1, 0.0)
}
