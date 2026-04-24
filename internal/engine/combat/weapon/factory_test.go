package weapon_test

import (
	"strings"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/combat/weapon"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/fp16"
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

// ---------------------------------------------------------------------------
// §4 RED-2 — Melee combo-schema factory cases
// ---------------------------------------------------------------------------

const comboOKJSON = `{
	"id": "player_melee",
	"type": "melee",
	"cooldown_frames": 20,
	"combo_window_frames": 15,
	"combo_steps": [
		{ "damage": 1, "active_frames": [4, 10], "hitbox": { "width": 24, "height": 16, "offset_x": 12, "offset_y": 0 } },
		{ "damage": 1, "active_frames": [3, 8],  "hitbox": { "width": 28, "height": 16, "offset_x": 14, "offset_y": -4 } },
		{ "damage": 2, "active_frames": [5, 12], "hitbox": { "width": 32, "height": 20, "offset_x": 16, "offset_y": 0 } }
	]
}`

func TestWeaponFactory_MeleeJSON(t *testing.T) {
	tests := []struct {
		name        string
		data        string
		wantErr     bool
		errContains string
		assertOK    func(t *testing.T, w *weapon.MeleeWeapon)
	}{
		{
			name:    "melee combo ok",
			data:    comboOKJSON,
			wantErr: false,
			assertOK: func(t *testing.T, w *weapon.MeleeWeapon) {
				if w.ID() != "player_melee" {
					t.Errorf("ID() = %q, want player_melee", w.ID())
				}
				steps := w.Steps()
				if len(steps) != 3 {
					t.Fatalf("len(Steps()) = %d, want 3", len(steps))
				}
				if steps[0].Damage != 1 {
					t.Errorf("Steps()[0].Damage = %d, want 1", steps[0].Damage)
				}
				if steps[2].Damage != 2 {
					t.Errorf("Steps()[2].Damage = %d, want 2", steps[2].Damage)
				}
				if steps[1].HitboxOffsetY16 != fp16.To16(-4) {
					t.Errorf("Steps()[1].HitboxOffsetY16 = %d, want %d (fp16 of -4)", steps[1].HitboxOffsetY16, fp16.To16(-4))
				}
				if steps[0].ActiveFrames != [2]int{4, 10} {
					t.Errorf("Steps()[0].ActiveFrames = %v, want [4 10]", steps[0].ActiveFrames)
				}
				if steps[2].HitboxW16 != fp16.To16(32) {
					t.Errorf("Steps()[2].HitboxW16 = %d, want %d (fp16 of 32)", steps[2].HitboxW16, fp16.To16(32))
				}
			},
		},
		{
			name: "melee combo missing combo_steps",
			data: `{
				"id": "player_melee",
				"type": "melee",
				"cooldown_frames": 20,
				"combo_window_frames": 15
			}`,
			wantErr:     true,
			errContains: "combo_steps",
		},
		{
			name: "melee combo too many steps",
			data: `{
				"id": "player_melee",
				"type": "melee",
				"cooldown_frames": 20,
				"combo_window_frames": 15,
				"combo_steps": [
					{ "damage": 1, "active_frames": [4, 10], "hitbox": { "width": 24, "height": 16, "offset_x": 12, "offset_y": 0 } },
					{ "damage": 1, "active_frames": [3, 8],  "hitbox": { "width": 28, "height": 16, "offset_x": 14, "offset_y": -4 } },
					{ "damage": 2, "active_frames": [5, 12], "hitbox": { "width": 32, "height": 20, "offset_x": 16, "offset_y": 0 } },
					{ "damage": 3, "active_frames": [5, 12], "hitbox": { "width": 32, "height": 20, "offset_x": 16, "offset_y": 0 } }
				]
			}`,
			wantErr:     true,
			errContains: "combo_steps",
		},
		{
			name: "melee combo step missing hitbox",
			data: `{
				"id": "player_melee",
				"type": "melee",
				"cooldown_frames": 20,
				"combo_window_frames": 15,
				"combo_steps": [
					{ "damage": 1, "active_frames": [4, 10], "hitbox": { "width": 24, "height": 16, "offset_x": 12, "offset_y": 0 } },
					{ "damage": 1, "active_frames": [3, 8] }
				]
			}`,
			wantErr:     true,
			errContains: "hitbox",
		},
		{
			name: "melee combo negative window",
			data: `{
				"id": "player_melee",
				"type": "melee",
				"cooldown_frames": 20,
				"combo_window_frames": -1,
				"combo_steps": [
					{ "damage": 1, "active_frames": [4, 10], "hitbox": { "width": 24, "height": 16, "offset_x": 12, "offset_y": 0 } }
				]
			}`,
			wantErr:     true,
			errContains: "combo_window_frames",
		},
		{
			name: "melee combo inverted active_frames for step 2",
			data: `{
				"id": "player_melee",
				"type": "melee",
				"cooldown_frames": 20,
				"combo_window_frames": 15,
				"combo_steps": [
					{ "damage": 1, "active_frames": [4, 10], "hitbox": { "width": 24, "height": 16, "offset_x": 12, "offset_y": 0 } },
					{ "damage": 1, "active_frames": [8, 3], "hitbox": { "width": 28, "height": 16, "offset_x": 14, "offset_y": -4 } }
				]
			}`,
			wantErr:     true,
			errContains: "active_frames",
		},
		{
			name: "melee combo step invalid hitbox dimensions",
			data: `{
				"id": "player_melee",
				"type": "melee",
				"cooldown_frames": 20,
				"combo_window_frames": 15,
				"combo_steps": [
					{ "damage": 1, "active_frames": [4, 10], "hitbox": { "width": 0, "height": 16, "offset_x": 12, "offset_y": 0 } }
				]
			}`,
			wantErr:     true,
			errContains: "hitbox",
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

// §4 RED-4 — JSON-driven end-to-end combo chain smoke test.
// Builds the weapon from the AC6 JSON and runs a full 3-hit chain, asserting
// each step's damage value is pulled from the JSON (AC2 + AC6 integration).
func TestWeaponFactory_MeleeJSON_DrivesFullComboChain(t *testing.T) {
	w, err := weapon.NewWeaponFromJSON([]byte(comboOKJSON), &mockProjectileManager{})
	if err != nil {
		t.Fatalf("unexpected error building weapon: %v", err)
	}
	m, ok := w.(*weapon.MeleeWeapon)
	if !ok {
		t.Fatalf("factory returned %T, want *weapon.MeleeWeapon", w)
	}
	owner := newMeleeOwner(100, 100, combat.FactionPlayer, animation.FaceDirectionRight)
	m.SetOwner(owner)

	wantDamagePerStep := []int{1, 1, 2}
	wantFirstActive := []int{4, 3, 5} // from the JSON

	for stepIdx, wantDmg := range wantDamagePerStep {
		if m.StepIndex() != stepIdx {
			t.Fatalf("step %d: StepIndex() = %d, want %d", stepIdx, m.StepIndex(), stepIdx)
		}

		enemy := newMeleeTarget("enemy", 110, 100, 8, 8, combat.FactionEnemy)
		sp := &fakeSpace{}
		sp.AddBody(enemy)

		// Reset cooldown for the test — the JSON sets cooldown=20 which would
		// otherwise block the next Fire. This test exercises combo mechanics,
		// not cooldown gating (cooldown gating is covered by the parity tests).
		m.SetCooldown(0)

		m.Fire(owner.x16, owner.y16, owner.face, body.ShootDirectionStraight, 0)
		// Advance to this step's first active frame.
		for i := 0; i < wantFirstActive[stepIdx]; i++ {
			m.Update()
		}
		if !m.IsHitboxActive() {
			t.Fatalf("step %d: expected hitbox active at frame %d", stepIdx, wantFirstActive[stepIdx])
		}
		m.ApplyHitbox(sp)
		if len(enemy.damageCalls) != 1 || enemy.damageCalls[0] != wantDmg {
			t.Errorf("step %d: damageCalls = %v, want [%d]", stepIdx, enemy.damageCalls, wantDmg)
		}

		// Finish the step's active window and pass frame[1]+1 to open the window (or reset).
		// Step active-frame last values: 10, 8, 12.
		lastActive := []int{10, 8, 12}[stepIdx]
		for i := wantFirstActive[stepIdx]; i <= lastActive+1; i++ {
			m.Update()
		}
		if stepIdx < 2 {
			if !m.AdvanceCombo() {
				t.Fatalf("step %d: AdvanceCombo() = false; want true", stepIdx)
			}
		}
	}

	// After step 3, chain has wrapped (AC4).
	if m.StepIndex() != 0 {
		t.Errorf("after full chain, StepIndex() = %d, want 0 (AC4 auto-reset)", m.StepIndex())
	}
}
