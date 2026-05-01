package combat

import (
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/hajimehoshi/ebiten/v2"
)

// ProjectileManager handles spawning projectiles into the world.
type ProjectileManager interface {
	SpawnProjectile(projectileType string, x16, y16, vx16, vy16, damage int, owner interface{})
	Update()
	Draw(screen *ebiten.Image)
	DrawWithOffset(screen *ebiten.Image, camX, camY float64)
	DrawCollisionBoxesWithOffset(draw func(b body.Collidable))
	Clear()
}
