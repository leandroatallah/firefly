// internal/engine/contracts/projectile/projectile.go
package projectile

import "github.com/hajimehoshi/ebiten/v2"

// Manager defines the interface for projectile management to avoid circular imports.
type Manager interface {
	Spawn(cfg interface{}, x16, y16, vx16, vy16 int, owner interface{})
	Update()
	Draw(screen *ebiten.Image)
	Clear()
}
