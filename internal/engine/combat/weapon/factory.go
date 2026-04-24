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
		ID                string `json:"id"`
		CooldownFrames    int    `json:"cooldown_frames"`
		ComboWindowFrames int    `json:"combo_window_frames"`
		ComboSteps        []struct {
			Damage       int    `json:"damage"`
			ActiveFrames [2]int `json:"active_frames"`
			Hitbox       *struct {
				Width   int `json:"width"`
				Height  int `json:"height"`
				OffsetX int `json:"offset_x"`
				OffsetY int `json:"offset_y"`
			} `json:"hitbox"`
		} `json:"combo_steps"`
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	if len(config.ComboSteps) < 1 || len(config.ComboSteps) > 3 {
		return nil, fmt.Errorf("combo_steps must contain 1..3 entries")
	}
	if config.ComboWindowFrames < 0 {
		return nil, fmt.Errorf("invalid combo_window_frames")
	}
	if config.CooldownFrames < 0 {
		return nil, fmt.Errorf("invalid cooldown_frames")
	}

	steps := make([]ComboStep, len(config.ComboSteps))
	for i, cs := range config.ComboSteps {
		if cs.Hitbox == nil {
			return nil, fmt.Errorf("hitbox is required for melee combo step %d", i)
		}
		if cs.ActiveFrames[0] < 0 || cs.ActiveFrames[1] < cs.ActiveFrames[0] {
			return nil, fmt.Errorf("invalid active_frames for combo step %d", i)
		}
		if cs.Hitbox.Width <= 0 || cs.Hitbox.Height <= 0 {
			return nil, fmt.Errorf("invalid hitbox dimensions for combo step %d", i)
		}
		steps[i] = ComboStep{
			Damage:          cs.Damage,
			ActiveFrames:    cs.ActiveFrames,
			HitboxW16:       fp16.To16(cs.Hitbox.Width),
			HitboxH16:       fp16.To16(cs.Hitbox.Height),
			HitboxOffsetX16: fp16.To16(cs.Hitbox.OffsetX),
			HitboxOffsetY16: fp16.To16(cs.Hitbox.OffsetY),
		}
	}

	return NewMeleeWeapon(config.ID, config.CooldownFrames, config.ComboWindowFrames, steps), nil
}
