package mocks

// MockShooter is a shared mock for the Shooter interface.
type MockShooter struct {
	SpawnBulletFunc func(x16, y16, vx16, vy16 int, owner interface{})
}

func (m *MockShooter) SpawnBullet(x16, y16, vx16, vy16 int, owner interface{}) {
	if m.SpawnBulletFunc != nil {
		m.SpawnBulletFunc(x16, y16, vx16, vy16, owner)
	}
}
