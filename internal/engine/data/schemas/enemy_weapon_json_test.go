package schemas

import (
	"encoding/json"
	"testing"
)

// TestSpriteData_WeaponBlock verifies that SpriteData.Weapon populates with
// all enemy weapon fields (projectile_type, speed, cooldown, damage, range,
// shoot_mode, shoot_direction, shoot_state) when present in JSON.
func TestSpriteData_WeaponBlock(t *testing.T) {
	raw := []byte(`{
		"body_rect": {"x": 0, "y": 0, "width": 16, "height": 16},
		"assets": {},
		"frame_rate": 8,
		"facing_direction": 0,
		"weapon": {
			"projectile_type": "bullet_small",
			"speed": 6,
			"cooldown": 90,
			"damage": 1,
			"range": 160,
			"shoot_mode": "on_sight",
			"shoot_direction": "horizontal",
			"shoot_state": "walk"
		}
	}`)

	var sd SpriteData
	if err := json.Unmarshal(raw, &sd); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if sd.Weapon == nil {
		t.Fatal("SpriteData.Weapon is nil; expected populated EnemyWeaponConfig")
	}
	w := sd.Weapon
	if w.ProjectileType != "bullet_small" {
		t.Errorf("ProjectileType = %q, want %q", w.ProjectileType, "bullet_small")
	}
	if w.Speed != 6 {
		t.Errorf("Speed = %d, want 6", w.Speed)
	}
	if w.Cooldown != 90 {
		t.Errorf("Cooldown = %d, want 90", w.Cooldown)
	}
	if w.Damage != 1 {
		t.Errorf("Damage = %d, want 1", w.Damage)
	}
	if w.Range != 160 {
		t.Errorf("Range = %d, want 160", w.Range)
	}
	if w.ShootMode != "on_sight" {
		t.Errorf("ShootMode = %q, want %q", w.ShootMode, "on_sight")
	}
	if w.ShootDirection != "horizontal" {
		t.Errorf("ShootDirection = %q, want %q", w.ShootDirection, "horizontal")
	}
	if w.ShootState != "walk" {
		t.Errorf("ShootState = %q, want %q", w.ShootState, "walk")
	}
}

// TestSpriteData_NoWeaponBlock preserves backward compatibility: SpriteData.Weapon
// is nil when the JSON omits the weapon block entirely.
func TestSpriteData_NoWeaponBlock(t *testing.T) {
	raw := []byte(`{
		"body_rect": {"x": 0, "y": 0, "width": 16, "height": 16},
		"assets": {},
		"frame_rate": 8,
		"facing_direction": 0
	}`)

	var sd SpriteData
	if err := json.Unmarshal(raw, &sd); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if sd.Weapon != nil {
		t.Errorf("SpriteData.Weapon should be nil when block is absent, got %+v", sd.Weapon)
	}
}

// TestSpriteData_WeaponOptionalKeysDefault verifies that optional string keys
// (shoot_mode, shoot_direction, shoot_state) default to the empty string when
// omitted from within the weapon block.
func TestSpriteData_WeaponOptionalKeysDefault(t *testing.T) {
	raw := []byte(`{
		"body_rect": {"x": 0, "y": 0, "width": 16, "height": 16},
		"assets": {},
		"frame_rate": 8,
		"facing_direction": 0,
		"weapon": {
			"projectile_type": "bullet_small",
			"speed": 5,
			"cooldown": 60,
			"damage": 1,
			"range": 0
		}
	}`)

	var sd SpriteData
	if err := json.Unmarshal(raw, &sd); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if sd.Weapon == nil {
		t.Fatal("SpriteData.Weapon is nil; expected populated EnemyWeaponConfig")
	}
	if sd.Weapon.ShootMode != "" {
		t.Errorf("ShootMode default: got %q, want \"\"", sd.Weapon.ShootMode)
	}
	if sd.Weapon.ShootDirection != "" {
		t.Errorf("ShootDirection default: got %q, want \"\"", sd.Weapon.ShootDirection)
	}
	if sd.Weapon.ShootState != "" {
		t.Errorf("ShootState default: got %q, want \"\"", sd.Weapon.ShootState)
	}
}
