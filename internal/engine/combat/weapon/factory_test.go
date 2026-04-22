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

func TestWeaponFactory_UnknownType_ReturnsError(t *testing.T) {
	data := []byte(`{"id": "mystery", "type": "laser"}`)

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

// --- RED-2: melee factory cases ------------------------------------------

func TestWeaponFactory_MeleeJSON(t *testing.T) {
	tests := []struct {
		name        string
		data        string
		wantErr     bool
		errContains string
		assertOK    func(t *testing.T, w *weapon.MeleeWeapon)
	}{
		{
			name: "melee weapon ok",
			data: `{
				"id": "player_melee",
				"type": "melee",
				"damage": 1,
				"cooldown_frames": 20,
				"active_frames": [4, 10],
				"hitbox": {"width": 24, "height": 16, "offset_x": 12, "offset_y": 0}
			}`,
			wantErr: false,
			assertOK: func(t *testing.T, w *weapon.MeleeWeapon) {
				if w.ID() != "player_melee" {
					t.Errorf("ID() = %q, want player_melee", w.ID())
				}
				if w.Damage() != 1 {
					t.Errorf("Damage() = %d, want 1", w.Damage())
				}
				if w.ActiveFrames() != [2]int{4, 10} {
					t.Errorf("ActiveFrames() = %v, want [4 10]", w.ActiveFrames())
				}
			},
		},
		{
			name: "melee missing hitbox",
			data: `{
				"id": "player_melee",
				"type": "melee",
				"damage": 1,
				"cooldown_frames": 20,
				"active_frames": [4, 10]
			}`,
			wantErr:     true,
			errContains: "hitbox",
		},
		{
			name: "melee inverted active_frames",
			data: `{
				"id": "player_melee",
				"type": "melee",
				"damage": 1,
				"cooldown_frames": 20,
				"active_frames": [10, 4],
				"hitbox": {"width": 24, "height": 16, "offset_x": 12, "offset_y": 0}
			}`,
			wantErr:     true,
			errContains: "active_frames",
		},
		{
			name:        "unknown type",
			data:        `{"id": "mystery", "type": "laser"}`,
			wantErr:     true,
			errContains: "unsupported weapon type",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w, err := weapon.NewWeaponFromJSON([]byte(tc.data), &mockProjectileManager{})
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil (weapon=%v)", w)
				}
				if tc.errContains != "" && !strings.Contains(err.Error(), tc.errContains) {
					t.Errorf("error message %q does not contain %q", err.Error(), tc.errContains)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			m, ok := w.(*weapon.MeleeWeapon)
			if !ok {
				t.Fatalf("factory returned %T, want *weapon.MeleeWeapon", w)
			}
			if tc.assertOK != nil {
				tc.assertOK(t, m)
			}
		})
	}
}
