package weapon_test

type mockProjectileManager struct {
	SpawnProjectileFunc func(projectileType string, x16, y16, vx16, vy16 int, owner interface{})
}

func (m *mockProjectileManager) SpawnProjectile(projectileType string, x16, y16, vx16, vy16 int, owner interface{}) {
	if m.SpawnProjectileFunc != nil {
		m.SpawnProjectileFunc(projectileType, x16, y16, vx16, vy16, owner)
	}
}
