package gamestates

import contractsbody "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"

type Bullet struct {
	body     contractsbody.MovableCollidable
	space    contractsbody.BodiesSpace
	speedX16 int
}

func NewBullet(body contractsbody.MovableCollidable, space contractsbody.BodiesSpace, speedX16 int) *Bullet {
	return &Bullet{body: body, space: space, speedX16: speedX16}
}

func (b *Bullet) Update() {
	b.body.SetVelocity(b.speedX16, 0)
	b.space.ResolveCollisions(b.body)

	x, y := b.body.GetPosition16()
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
