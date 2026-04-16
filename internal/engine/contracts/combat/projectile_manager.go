package combat

// ProjectileManager handles spawning projectiles into the world.
type ProjectileManager interface {
	SpawnProjectile(projectileType string, x16, y16, vx16, vy16, damage int, owner interface{})
}
