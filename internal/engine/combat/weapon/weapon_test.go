package weapon_test

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/combat/weapon"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/mocks"
)

func TestProjectileWeapon_Fire_SpawnsProjectileWithCorrectVelocity(t *testing.T) {
	tests := []struct {
		name     string
		faceDir  animation.FacingDirectionEnum
		shootDir body.ShootDirection
		wantVx16 int
		wantVy16 int
	}{
		{
			name:     "diagonal up forward right",
			faceDir:  animation.FaceDirectionRight,
			shootDir: body.ShootDirectionDiagonalUpForward,
			wantVx16: 70, // 100 * 707 / 1000
			wantVy16: -70,
		},
		{
			name:     "diagonal down forward right",
			faceDir:  animation.FaceDirectionRight,
			shootDir: body.ShootDirectionDiagonalDownForward,
			wantVx16: 70,
			wantVy16: 70,
		},
		{
			name:     "straight left",
			faceDir:  animation.FaceDirectionLeft,
			shootDir: body.ShootDirectionStraight,
			wantVx16: -100,
			wantVy16: 0,
		},
		{
			name:     "down",
			faceDir:  animation.FaceDirectionRight,
			shootDir: body.ShootDirectionDown,
			wantVx16: 0,
			wantVy16: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type spawnCall struct {
				projectileType string
				x16, y16       int
				vx16, vy16     int
				owner          interface{}
			}

			var got spawnCall
			mock := &mockProjectileManager{
				SpawnProjectileFunc: func(projectileType string, x16, y16, vx16, vy16 int, owner interface{}) {
					got = spawnCall{projectileType, x16, y16, vx16, vy16, owner}
				},
			}

			w := weapon.NewProjectileWeapon("test", 10, "bullet", 100, mock, "")

			w.Fire(1000, 2000, tt.faceDir, tt.shootDir)

			if got.projectileType != "bullet" {
				t.Errorf("projectileType: got %q, want %q", got.projectileType, "bullet")
			}
			if got.x16 != 1000 || got.y16 != 2000 {
				t.Errorf("position: got (%d,%d), want (1000,2000)", got.x16, got.y16)
			}
			if got.vx16 != tt.wantVx16 || got.vy16 != tt.wantVy16 {
				t.Errorf("velocity: got (%d,%d), want (%d,%d)", got.vx16, got.vy16, tt.wantVx16, tt.wantVy16)
			}
			if w.CanFire() {
				t.Error("CanFire() should return false immediately after firing")
			}
			if w.Cooldown() != 10 {
				t.Errorf("Cooldown(): got %d, want 10", w.Cooldown())
			}
		})
	}
}

func TestProjectileWeapon_Update_DecrementsCooldown(t *testing.T) {
	w := weapon.NewProjectileWeapon("test", 3, "bullet", 100, &mockProjectileManager{}, "")

	w.Fire(0, 0, animation.FaceDirectionRight, body.ShootDirectionStraight)
	w.Update()

	if w.Cooldown() != 2 {
		t.Errorf("Cooldown() after one Update: got %d, want 2", w.Cooldown())
	}
}

func TestProjectileWeapon_CanFire_TrueAfterCooldownExpires(t *testing.T) {
	w := weapon.NewProjectileWeapon("test", 2, "bullet", 100, &mockProjectileManager{}, "")

	w.Fire(0, 0, animation.FaceDirectionRight, body.ShootDirectionStraight)
	w.Update()
	w.Update()

	if !w.CanFire() {
		t.Error("CanFire() should return true after cooldown expires")
	}
}

func TestProjectileWeapon_SetCooldown(t *testing.T) {
	w := weapon.NewProjectileWeapon("test", 10, "bullet", 100, &mockProjectileManager{}, "")
	w.SetCooldown(5)

	if w.Cooldown() != 5 {
		t.Errorf("Cooldown() after SetCooldown(5): got %d, want 5", w.Cooldown())
	}
}

func TestProjectileWeapon_ID(t *testing.T) {
	w := weapon.NewProjectileWeapon("my-weapon", 10, "bullet", 100, &mockProjectileManager{}, "")
	if w.ID() != "my-weapon" {
		t.Errorf("ID(): got %q, want %q", w.ID(), "my-weapon")
	}
}

// TestProjectileWeapon_MuzzleFlashVFX_ExecutionOrder covers AC3:
// Fire() must call SpawnPuff BEFORE calling SpawnProjectile.
func TestProjectileWeapon_MuzzleFlashVFX_ExecutionOrder(t *testing.T) {
	var callOrder []string
	mockVFX := &mocks.MockVFXManager{
		SpawnPuffFunc: func(_ string, _, _ float64, _ int, _ float64) {
			callOrder = append(callOrder, "vfx")
		},
	}
	mockProj := &mockProjectileManager{
		SpawnProjectileFunc: func(_ string, _, _, _, _ int, _ interface{}) {
			callOrder = append(callOrder, "projectile")
		},
	}

	w := weapon.NewProjectileWeapon("gun", 10, "bullet", 160, mockProj, "muzzle_flash")
	w.SetVFXManager(mockVFX)

	w.Fire(0, 0, animation.FaceDirectionRight, body.ShootDirectionStraight)

	if len(callOrder) != 2 {
		t.Fatalf("expected 2 calls, got %d", len(callOrder))
	}
	if callOrder[0] != "vfx" || callOrder[1] != "projectile" {
		t.Errorf("expected vfx then projectile, got %v", callOrder)
	}
}

// TestProjectileWeapon_MuzzleFlashVFX covers AC3/AC4/AC6:
// Fire() must call SpawnPuff with the correct typeKey, fp16-converted coordinates,
// count=1, and randRange=0.0.
func TestProjectileWeapon_MuzzleFlashVFX(t *testing.T) {
	tests := []struct {
		name      string
		x16       int
		y16       int
		wantX     float64
		wantY     float64
		effectKey string
	}{
		{
			name:      "standard position 320x480",
			x16:       320,
			y16:       480,
			wantX:     20.0,
			wantY:     30.0,
			effectKey: "muzzle_flash",
		},
		{
			name:      "origin 0x0",
			x16:       0,
			y16:       0,
			wantX:     0.0,
			wantY:     0.0,
			effectKey: "muzzle_flash",
		},
		{
			name:      "arbitrary fp16 position",
			x16:       160,
			y16:       320,
			wantX:     10.0,
			wantY:     20.0,
			effectKey: "muzzle_flash",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type puffCall struct {
				typeKey   string
				x         float64
				y         float64
				count     int
				randRange float64
			}

			var got puffCall
			var spawnCount int

			mockVFX := &mocks.MockVFXManager{
				SpawnPuffFunc: func(typeKey string, x float64, y float64, count int, randRange float64) {
					spawnCount++
					got = puffCall{typeKey, x, y, count, randRange}
				},
			}

			w := weapon.NewProjectileWeapon("gun", 10, "bullet", 160, &mockProjectileManager{}, tt.effectKey)
			w.SetVFXManager(mockVFX)

			w.Fire(tt.x16, tt.y16, animation.FaceDirectionRight, body.ShootDirectionStraight)

			if spawnCount != 1 {
				t.Errorf("SpawnPuff call count: got %d, want 1", spawnCount)
			}
			if got.typeKey != tt.effectKey {
				t.Errorf("SpawnPuff typeKey: got %q, want %q", got.typeKey, tt.effectKey)
			}
			if got.x != tt.wantX {
				t.Errorf("SpawnPuff x: got %v, want %v", got.x, tt.wantX)
			}
			if got.y != tt.wantY {
				t.Errorf("SpawnPuff y: got %v, want %v", got.y, tt.wantY)
			}
			if got.count != 1 {
				t.Errorf("SpawnPuff count: got %d, want 1", got.count)
			}
			if got.randRange != 0.0 {
				t.Errorf("SpawnPuff randRange: got %v, want 0.0", got.randRange)
			}
		})
	}
}

// TestProjectileWeapon_NoVFXWhenManagerNil covers AC5:
// Fire() must not panic when vfxManager has not been set.
func TestProjectileWeapon_NoVFXWhenManagerNil(t *testing.T) {
	// No SetVFXManager call — manager stays nil.
	w := weapon.NewProjectileWeapon("gun", 10, "bullet", 160, &mockProjectileManager{}, "muzzle_flash")

	// Must not panic.
	w.Fire(320, 480, animation.FaceDirectionRight, body.ShootDirectionStraight)
}

// TestProjectileWeapon_NoVFXWhenEffectTypeEmpty covers AC5:
// Fire() must not call SpawnPuff when muzzleEffectType is empty string.
func TestProjectileWeapon_NoVFXWhenEffectTypeEmpty(t *testing.T) {
	var spawnCount int
	mockVFX := &mocks.MockVFXManager{
		SpawnPuffFunc: func(_ string, _ float64, _ float64, _ int, _ float64) {
			spawnCount++
		},
	}

	// muzzleEffectType is "" — VFX must be suppressed even when manager is set.
	w := weapon.NewProjectileWeapon("gun", 10, "bullet", 160, &mockProjectileManager{}, "")
	w.SetVFXManager(mockVFX)

	w.Fire(320, 480, animation.FaceDirectionRight, body.ShootDirectionStraight)

	if spawnCount != 0 {
		t.Errorf("SpawnPuff call count: got %d, want 0 (empty effectType must suppress VFX)", spawnCount)
	}
}

// TestProjectileWeapon_ConstructorAcceptsMuzzleEffectType covers AC1:
// NewProjectileWeapon must accept muzzleEffectType as the 6th parameter.
func TestProjectileWeapon_ConstructorAcceptsMuzzleEffectType(t *testing.T) {
	tests := []struct {
		name             string
		muzzleEffectType string
	}{
		{name: "non-empty effect type", muzzleEffectType: "muzzle_flash"},
		{name: "empty effect type disables VFX", muzzleEffectType: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// If NewProjectileWeapon does not accept 6 args, this will fail to compile.
			w := weapon.NewProjectileWeapon("gun", 10, "bullet", 160, &mockProjectileManager{}, tt.muzzleEffectType)
			if w == nil {
				t.Fatal("NewProjectileWeapon returned nil")
			}
		})
	}
}

// TestProjectileWeapon_SetVFXManager covers AC2:
// SetVFXManager must exist and accept a vfx.Manager implementation.
func TestProjectileWeapon_SetVFXManager(t *testing.T) {
	w := weapon.NewProjectileWeapon("gun", 10, "bullet", 160, &mockProjectileManager{}, "muzzle_flash")

	// If SetVFXManager does not exist this will fail to compile.
	mockVFX := &mocks.MockVFXManager{}
	w.SetVFXManager(mockVFX)

	// Verify the manager was accepted by confirming SpawnPuff is called on Fire.
	called := false
	mockVFX.SpawnPuffFunc = func(_ string, _ float64, _ float64, _ int, _ float64) {
		called = true
	}

	w.Fire(0, 0, animation.FaceDirectionRight, body.ShootDirectionStraight)

	if !called {
		t.Error("SpawnPuff was not called after SetVFXManager; manager injection may not be working")
	}
}
