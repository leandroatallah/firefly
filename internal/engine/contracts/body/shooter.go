package body

// Shooter is the contract ShootingSkill depends on to spawn a bullet into the world.
type Shooter interface {
	// SpawnBullet adds a bullet Body to the BodiesSpace.
	// x16, y16 are fixed-point spawn position; speedX16 is signed fixed-point horizontal velocity.
	SpawnBullet(x16, y16, speedX16 int, owner interface{})
}
