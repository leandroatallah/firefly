package gameplayer_test

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	actors "github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/mocks"
	gameplayer "github.com/boilerplate/ebiten-template/internal/game/entity/actors/player"
)

type mockProjectileManager struct{}

func (m *mockProjectileManager) SpawnProjectile(projectileType string, x16, y16, vx16, vy16, damage int, owner interface{}) {
}

func TestNewClimberInventory(t *testing.T) {
	spawnPuffCalled := false
	vfxMock := &mocks.MockVFXManager{
		SpawnPuffFunc: func(typeKey string, x float64, y float64, count int, randRange float64) {
			if typeKey == "muzzle_flash" {
				spawnPuffCalled = true
			}
		},
	}

	inv := gameplayer.NewClimberInventory(&mockProjectileManager{}, vfxMock)

	if inv.ActiveWeapon().ID() != "light_blaster" {
		t.Fatalf("expected light_blaster, got %s", inv.ActiveWeapon().ID())
	}
	if !inv.ActiveWeapon().CanFire() {
		t.Fatal("expected CanFire() == true")
	}

	// Test VFX trigger on fire
	inv.ActiveWeapon().Fire(160, 160, animation.FaceDirectionRight, body.ShootDirectionStraight, 0)
	if !spawnPuffCalled {
		t.Error("expected SpawnPuff to be called for light_blaster")
	}

	inv.SwitchNext()
	if inv.ActiveWeapon().ID() != "heavy_cannon" {
		t.Fatalf("expected heavy_cannon, got %s", inv.ActiveWeapon().ID())
	}

	spawnPuffCalled = false
	inv.ActiveWeapon().Fire(160, 160, animation.FaceDirectionRight, body.ShootDirectionStraight, 0)
	if !spawnPuffCalled {
		t.Error("expected SpawnPuff to be called for heavy_cannon")
	}
}

// TestLoadShootingSkill_StateSpawnOffsets covers US-037 AC1-AC3:
// BuildStateSpawnOffsets converts a JSON-parsed pixel-int offset map keyed by
// state name into an fp16 offset map keyed by actor state enum. Unknown state
// names are skipped with a log warning (no error returned). A nil input yields
// a nil/empty output.
func TestLoadShootingSkill_StateSpawnOffsets(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string]gameplayer.StateOffsetEntry
		want    map[int][2]int
		wantLen int
	}{
		{
			name:    "absent/nil map produces nil table",
			input:   nil,
			want:    nil,
			wantLen: 0,
		},
		{
			name: "valid entries are registered with fp16 conversion",
			input: map[string]gameplayer.StateOffsetEntry{
				"duck": {X: 6, Y: 12},
			},
			want: map[int][2]int{
				int(actors.Ducking): {96, 192},
			},
			wantLen: 1,
		},
		{
			name: "unknown state name is skipped",
			input: map[string]gameplayer.StateOffsetEntry{
				"bogus": {X: 1, Y: 1},
				"duck":  {X: 2, Y: 3},
			},
			want: map[int][2]int{
				int(actors.Ducking): {32, 48},
			},
			wantLen: 1,
		},
		{
			name: "multiple states load independently",
			input: map[string]gameplayer.StateOffsetEntry{
				"duck":       {X: 1, Y: 1},
				"jump_shoot": {X: 2, Y: 2},
			},
			want: map[int][2]int{
				int(actors.Ducking):         {16, 16},
				int(actors.JumpingShooting): {32, 32},
			},
			wantLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gameplayer.BuildStateSpawnOffsets(tt.input)

			if len(got) != tt.wantLen {
				t.Errorf("len: got %d, want %d (map=%v)", len(got), tt.wantLen, got)
			}

			for k, v := range tt.want {
				gv, ok := got[k]
				if !ok {
					t.Errorf("missing key %d in output", k)
					continue
				}
				if gv != v {
					t.Errorf("key %d: got %v, want %v", k, gv, v)
				}
			}
		})
	}
}
