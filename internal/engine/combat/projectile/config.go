// internal/engine/combat/projectile/config.go
package projectile

import enginecombat "github.com/boilerplate/ebiten-template/internal/engine/combat"

// ProjectileConfig defines the basic properties for spawning a projectile.
type ProjectileConfig struct {
	Width          int
	Height         int
	Damage         int
	Faction        enginecombat.Faction `json:"faction,omitempty"`         // default 0 = Neutral
	ImpactEffect   string               `json:"impact_effect,omitempty"`   // VFX type for collision impacts
	DespawnEffect  string               `json:"despawn_effect,omitempty"`  // VFX type for lifetime expiration
	LifetimeFrames int                  `json:"lifetime_frames,omitempty"` // 0 = infinite (backward compat)
	Interceptable  bool                 `json:"interceptable,omitempty"`   // true = can be hit by other projectiles
}
