// internal/engine/combat/projectile/config.go
package projectile

// ProjectileConfig defines the basic properties for spawning a projectile.
type ProjectileConfig struct {
	Width          int
	Height         int
	Damage         int
	ImpactEffect   string `json:"impact_effect,omitempty"`   // VFX type for collision impacts
	DespawnEffect  string `json:"despawn_effect,omitempty"`  // VFX type for lifetime expiration
	LifetimeFrames int    `json:"lifetime_frames,omitempty"` // 0 = infinite (backward compat)
}
