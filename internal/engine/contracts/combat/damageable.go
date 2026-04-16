package combat

// Damageable is implemented by any entity that can receive damage.
type Damageable interface {
	TakeDamage(amount int)
}

// Destructible extends Damageable with a lifecycle query: the projectile hit
// path can query whether the target has been destroyed, without requiring
// special-casing in the projectile itself.
type Destructible interface {
	Damageable
	IsDestroyed() bool
}
