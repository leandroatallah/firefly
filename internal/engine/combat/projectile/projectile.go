// internal/engine/combat/projectile/projectile.go
package projectile

import (
	contractsbody "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
)

// projectile is the internal state of a spawned projectile.
type projectile struct {
	movable  contractsbody.Movable
	body     contractsbody.Collidable
	space    contractsbody.BodiesSpace
	speedX16 int
	speedY16 int
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

	if x < 0 || y < 0 || x > w<<4 || y > h<<4 {
		p.space.QueueForRemoval(p.body)
	}
}

func (p *projectile) OnTouch(other contractsbody.Collidable) {
	if other != p.body.Owner() {
		p.space.QueueForRemoval(p.body)
	}
}

func (p *projectile) OnBlock(_ contractsbody.Collidable) {
	p.space.QueueForRemoval(p.body)
}
