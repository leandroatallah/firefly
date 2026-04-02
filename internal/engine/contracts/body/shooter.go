package body

// Shooter is the contract ShootingSkill depends on to spawn a bullet into the world.
type Shooter interface {
	// SpawnBullet adds a bullet Body to the BodiesSpace.
	// x16, y16 are fixed-point spawn position; vx16, vy16 are signed fixed-point velocities.
	SpawnBullet(x16, y16, vx16, vy16 int, owner interface{})
}
