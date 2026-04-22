package weapon

import (
	"encoding/json"
	"fmt"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/fp16"
)

// NewWeaponFromJSON creates a Weapon from JSON configuration.
func NewWeaponFromJSON(data []byte, manager combat.ProjectileManager) (combat.Weapon, error) {
	var base struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &base); err != nil {
		return nil, err
	}

	switch base.Type {
	case "projectile":
		return parseProjectileWeapon(data, manager)
	case "melee":
		return parseMeleeWeapon(data)
	default:
		return nil, fmt.Errorf("unsupported weapon type: %s", base.Type)
	}
}

func parseProjectileWeapon(data []byte, manager combat.ProjectileManager) (combat.Weapon, error) {
	var config struct {
		ID               string `json:"id"`
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
	if config.Projectile == nil {
		return nil, fmt.Errorf("projectile object is required for projectile weapons")
	}
	w := NewProjectileWeapon(config.ID, config.CooldownFrames, config.Projectile.Type, config.Projectile.Speed, manager, config.MuzzleEffectType, 0, 0)
	w.SetDamage(config.Projectile.Damage)
	return w, nil
}

func parseMeleeWeapon(data []byte) (*MeleeWeapon, error) {
	var config struct {
		ID             string `json:"id"`
		Damage         int    `json:"damage"`
		CooldownFrames int    `json:"cooldown_frames"`
		ActiveFrames   [2]int `json:"active_frames"`
		Hitbox         *struct {
			Width   int `json:"width"`
			Height  int `json:"height"`
			OffsetX int `json:"offset_x"`
			OffsetY int `json:"offset_y"`
		} `json:"hitbox"`
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	if config.Hitbox == nil {
		return nil, fmt.Errorf("hitbox is required for melee weapons")
	}
	if config.ActiveFrames[0] < 0 || config.ActiveFrames[1] < config.ActiveFrames[0] {
		return nil, fmt.Errorf("invalid active_frames")
	}
	if config.Hitbox.Width <= 0 || config.Hitbox.Height <= 0 {
		return nil, fmt.Errorf("invalid hitbox dimensions")
	}
	if config.CooldownFrames < 0 {
		return nil, fmt.Errorf("invalid cooldown_frames")
	}
	w := NewMeleeWeapon(
		config.ID,
		config.Damage,
		config.CooldownFrames,
		config.ActiveFrames,
		fp16.To16(config.Hitbox.Width),
		fp16.To16(config.Hitbox.Height),
		fp16.To16(config.Hitbox.OffsetX),
		fp16.To16(config.Hitbox.OffsetY),
	)
	return w, nil
}
