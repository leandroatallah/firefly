package body

// Projectile is implemented by physics bodies that represent projectiles.
// Interceptable returns true when the projectile can be intercepted (hit) by
// other projectiles; false means it is invisible to projectile-vs-projectile
// collision resolution.
type Projectile interface {
	Interceptable() bool
}
