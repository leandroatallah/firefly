package gamestates

import contractsbody "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"

type Bullet struct {
	movable  contractsbody.Movable
	body     contractsbody.Collidable
	space    contractsbody.BodiesSpace
	speedX16 int
	speedY16 int
}

func NewBullet(movable contractsbody.Movable, body contractsbody.Collidable, space contractsbody.BodiesSpace, speedX16, speedY16 int) *Bullet {
	return &Bullet{movable: movable, body: body, space: space, speedX16: speedX16, speedY16: speedY16}
}

func (b *Bullet) Body() contractsbody.Collidable {
	return b.body
}

func (b *Bullet) Update() {
	x, y := b.body.GetPosition16()
	x += b.speedX16
	y += b.speedY16
	b.body.SetPosition16(x, y)
	
	b.space.ResolveCollisions(b.body)

	provider := b.space.GetTilemapDimensionsProvider()
	w := provider.GetTilemapWidth()
	h := provider.GetTilemapHeight()

	if x < 0 || y < 0 || x > w<<4 || y > h<<4 {
		b.space.QueueForRemoval(b.body)
	}
}

func (b *Bullet) OnTouch(other contractsbody.Collidable) {
	if other != b.body.Owner() {
		b.space.QueueForRemoval(b.body)
	}
}

func (b *Bullet) OnBlock(other contractsbody.Collidable) {
	b.space.QueueForRemoval(b.body)
}
