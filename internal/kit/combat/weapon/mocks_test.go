package weapon_test

import (
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/hajimehoshi/ebiten/v2"
)

type mockProjectileManager struct {
	SpawnProjectileFunc func(projectileType string, x16, y16, vx16, vy16, damage int, owner interface{})
}

func (m *mockProjectileManager) SpawnProjectile(projectileType string, x16, y16, vx16, vy16, damage int, owner interface{}) {
	if m.SpawnProjectileFunc != nil {
		m.SpawnProjectileFunc(projectileType, x16, y16, vx16, vy16, damage, owner)
	}
}

func (m *mockProjectileManager) Update() {}

func (m *mockProjectileManager) Draw(screen *ebiten.Image) {}

func (m *mockProjectileManager) DrawWithOffset(screen *ebiten.Image, camX, camY float64) {}

func (m *mockProjectileManager) DrawCollisionBoxesWithOffset(draw func(b body.Collidable)) {}

func (m *mockProjectileManager) Clear() {}
