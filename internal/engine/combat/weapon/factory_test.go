package weapon_test

import (
	"strings"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/combat/weapon"
)

func TestWeaponFactory_ValidJSON_ReturnsWeapon(t *testing.T) {
	data := []byte(`{
		"id": "basic_blaster",
		"type": "projectile",
		"cooldown_frames": 15,
		"projectile": {
			"type": "bullet",
			"speed": 327680,
			"damage": 1
		}
	}`)

	w, err := weapon.NewWeaponFromJSON(data, &mockProjectileManager{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.ID() != "basic_blaster" {
		t.Errorf("ID(): got %q, want %q", w.ID(), "basic_blaster")
	}
	if w.Cooldown() != 0 {
		t.Errorf("initial Cooldown(): got %d, want 0", w.Cooldown())
	}
}

func TestWeaponFactory_InvalidType_ReturnsError(t *testing.T) {
	data := []byte(`{"id": "sword", "type": "melee"}`)

	_, err := weapon.NewWeaponFromJSON(data, &mockProjectileManager{})
	if err == nil {
		t.Fatal("expected error for unsupported weapon type, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported weapon type") {
		t.Errorf("error message: got %q, want it to contain %q", err.Error(), "unsupported weapon type")
	}
}

func TestWeaponFactory_MissingProjectileObject_ReturnsError(t *testing.T) {
	data := []byte(`{"id": "blaster", "type": "projectile", "cooldown_frames": 10}`)

	_, err := weapon.NewWeaponFromJSON(data, &mockProjectileManager{})
	if err == nil {
		t.Fatal("expected error for missing projectile object, got nil")
	}
}
