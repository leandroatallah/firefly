package weapon

import (
	"encoding/json"
	"fmt"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
)

// NewWeaponFromJSON creates a Weapon from JSON configuration.
func NewWeaponFromJSON(data []byte, manager combat.ProjectileManager) (combat.Weapon, error) {
	var config struct {
		ID               string `json:"id"`
		Type             string `json:"type"`
		CooldownFrames   int    `json:"cooldown_frames"`
		MuzzleEffectType string `json:"muzzle_effect_type"`
		Projectile       *struct {
			Type   string `json:"type"`
			Speed  int    `json:"speed"`
			Damage int    `json:"damage"`
		} `json:"projectile"`
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if config.Type != "projectile" {
		return nil, fmt.Errorf("unsupported weapon type: %s", config.Type)
	}

	if config.Projectile == nil {
		return nil, fmt.Errorf("projectile object is required for projectile weapons")
	}

	w := NewProjectileWeapon(config.ID, config.CooldownFrames, config.Projectile.Type, config.Projectile.Speed, manager, config.MuzzleEffectType, 0, 0)
	w.SetDamage(config.Projectile.Damage)
	return w, nil
}
