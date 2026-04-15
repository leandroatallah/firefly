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

			w := weapon.NewProjectileWeapon("test", 10, "bullet", 100, mock, "", 0, 0)

			w.Fire(1000, 2000, tt.faceDir, tt.shootDir, 0)

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
	w := weapon.NewProjectileWeapon("test", 3, "bullet", 100, &mockProjectileManager{}, "", 0, 0)

	w.Fire(0, 0, animation.FaceDirectionRight, body.ShootDirectionStraight, 0)
	w.Update()

	if w.Cooldown() != 2 {
		t.Errorf("Cooldown() after one Update: got %d, want 2", w.Cooldown())
	}
}

func TestProjectileWeapon_CanFire_TrueAfterCooldownExpires(t *testing.T) {
	w := weapon.NewProjectileWeapon("test", 2, "bullet", 100, &mockProjectileManager{}, "", 0, 0)

	w.Fire(0, 0, animation.FaceDirectionRight, body.ShootDirectionStraight, 0)
	w.Update()
	w.Update()

	if !w.CanFire() {
		t.Error("CanFire() should return true after cooldown expires")
	}
}

func TestProjectileWeapon_SetCooldown(t *testing.T) {
	w := weapon.NewProjectileWeapon("test", 10, "bullet", 100, &mockProjectileManager{}, "", 0, 0)
	w.SetCooldown(5)

	if w.Cooldown() != 5 {
		t.Errorf("Cooldown() after SetCooldown(5): got %d, want 5", w.Cooldown())
	}
}

func TestProjectileWeapon_ID(t *testing.T) {
	w := weapon.NewProjectileWeapon("my-weapon", 10, "bullet", 100, &mockProjectileManager{}, "", 0, 0)
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

	w := weapon.NewProjectileWeapon("gun", 10, "bullet", 160, mockProj, "muzzle_flash", 0, 0)
	w.SetVFXManager(mockVFX)

	w.Fire(0, 0, animation.FaceDirectionRight, body.ShootDirectionStraight, 0)

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

			w := weapon.NewProjectileWeapon("gun", 10, "bullet", 160, &mockProjectileManager{}, tt.effectKey, 0, 0)
			w.SetVFXManager(mockVFX)

			w.Fire(tt.x16, tt.y16, animation.FaceDirectionRight, body.ShootDirectionStraight, 0)

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
	w := weapon.NewProjectileWeapon("gun", 10, "bullet", 160, &mockProjectileManager{}, "muzzle_flash", 0, 0)

	// Must not panic.
	w.Fire(320, 480, animation.FaceDirectionRight, body.ShootDirectionStraight, 0)
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
	w := weapon.NewProjectileWeapon("gun", 10, "bullet", 160, &mockProjectileManager{}, "", 0, 0)
	w.SetVFXManager(mockVFX)

	w.Fire(320, 480, animation.FaceDirectionRight, body.ShootDirectionStraight, 0)

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
			// If NewProjectileWeapon does not accept 8 args, this will fail to compile.
			w := weapon.NewProjectileWeapon("gun", 10, "bullet", 160, &mockProjectileManager{}, tt.muzzleEffectType, 0, 0)
			if w == nil {
				t.Fatal("NewProjectileWeapon returned nil")
			}
		})
	}
}

// TestProjectileWeapon_SetVFXManager covers AC2:
// SetVFXManager must exist and accept a vfx.Manager implementation.
func TestProjectileWeapon_SetVFXManager(t *testing.T) {
	w := weapon.NewProjectileWeapon("gun", 10, "bullet", 160, &mockProjectileManager{}, "muzzle_flash", 0, 0)

	// If SetVFXManager does not exist this will fail to compile.
	mockVFX := &mocks.MockVFXManager{}
	w.SetVFXManager(mockVFX)

	// Verify the manager was accepted by confirming SpawnPuff is called on Fire.
	called := false
	mockVFX.SpawnPuffFunc = func(_ string, _ float64, _ float64, _ int, _ float64) {
		called = true
	}

	w.Fire(0, 0, animation.FaceDirectionRight, body.ShootDirectionStraight, 0)

	if !called {
		t.Error("SpawnPuff was not called after SetVFXManager; manager injection may not be working")
	}
}

// TestProjectileWeapon_Fire_SpawnOffset covers AC1-AC6 for US-034:
// Configurable spawn offset applied to both projectile and VFX spawn positions.
func TestProjectileWeapon_Fire_SpawnOffset(t *testing.T) {
	tests := []struct {
		name           string
		spawnOffsetX16 int
		spawnOffsetY16 int
		fireX16        int
		fireY16        int
		faceDir        animation.FacingDirectionEnum
		wantSpawnX16   int
		wantSpawnY16   int
		wantVFXX       float64
		wantVFXY       float64
		hasVFX         bool
	}{
		{
			name:           "facing right applies positive X offset",
			spawnOffsetX16: 128,
			spawnOffsetY16: 64,
			fireX16:        320,
			fireY16:        480,
			faceDir:        animation.FaceDirectionRight,
			wantSpawnX16:   448,  // 320 + 128
			wantSpawnY16:   544,  // 480 + 64
			wantVFXX:       28.0, // 448 / 16
			wantVFXY:       34.0, // 544 / 16
			hasVFX:         true,
		},
		{
			name:           "facing left negates X offset",
			spawnOffsetX16: 128,
			spawnOffsetY16: 64,
			fireX16:        320,
			fireY16:        480,
			faceDir:        animation.FaceDirectionLeft,
			wantSpawnX16:   192,  // 320 - 128
			wantSpawnY16:   544,  // 480 + 64
			wantVFXX:       12.0, // 192 / 16
			wantVFXY:       34.0, // 544 / 16
			hasVFX:         true,
		},
		{
			name:           "zero offset preserves current behavior",
			spawnOffsetX16: 0,
			spawnOffsetY16: 0,
			fireX16:        320,
			fireY16:        480,
			faceDir:        animation.FaceDirectionRight,
			wantSpawnX16:   320,
			wantSpawnY16:   480,
			wantVFXX:       20.0,
			wantVFXY:       30.0,
			hasVFX:         true,
		},
		{
			name:           "zero offset without VFX manager",
			spawnOffsetX16: 0,
			spawnOffsetY16: 0,
			fireX16:        320,
			fireY16:        480,
			faceDir:        animation.FaceDirectionRight,
			wantSpawnX16:   320,
			wantSpawnY16:   480,
			hasVFX:         false,
		},
		{
			name:           "negative Y offset shifts spawn upward",
			spawnOffsetX16: 0,
			spawnOffsetY16: -32,
			fireX16:        160,
			fireY16:        320,
			faceDir:        animation.FaceDirectionRight,
			wantSpawnX16:   160,
			wantSpawnY16:   288,  // 320 - 32
			wantVFXX:       10.0, // 160 / 16
			wantVFXY:       18.0, // 288 / 16
			hasVFX:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotProjX16, gotProjY16 int
			mockProj := &mockProjectileManager{
				SpawnProjectileFunc: func(_ string, x16, y16, _, _ int, _ interface{}) {
					gotProjX16 = x16
					gotProjY16 = y16
				},
			}

			w := weapon.NewProjectileWeapon("gun", 10, "bullet", 160, mockProj, "muzzle_flash", tt.spawnOffsetX16, tt.spawnOffsetY16)

			var gotVFXX, gotVFXY float64
			var vfxCalled bool
			if tt.hasVFX {
				mockVFX := &mocks.MockVFXManager{
					SpawnPuffFunc: func(_ string, x, y float64, _ int, _ float64) {
						vfxCalled = true
						gotVFXX = x
						gotVFXY = y
					},
				}
				w.SetVFXManager(mockVFX)
			}

			w.Fire(tt.fireX16, tt.fireY16, tt.faceDir, body.ShootDirectionStraight, 0)

			if gotProjX16 != tt.wantSpawnX16 {
				t.Errorf("projectile x16: got %d, want %d", gotProjX16, tt.wantSpawnX16)
			}
			if gotProjY16 != tt.wantSpawnY16 {
				t.Errorf("projectile y16: got %d, want %d", gotProjY16, tt.wantSpawnY16)
			}

			if tt.hasVFX {
				if !vfxCalled {
					t.Fatal("expected SpawnPuff to be called")
				}
				if gotVFXX != tt.wantVFXX {
					t.Errorf("VFX x: got %v, want %v", gotVFXX, tt.wantVFXX)
				}
				if gotVFXY != tt.wantVFXY {
					t.Errorf("VFX y: got %v, want %v", gotVFXY, tt.wantVFXY)
				}
			}
		})
	}
}

// TestProjectileWeapon_Fire_StateSpawnOffset covers US-037 AC3-AC6:
// Per-state offset table overrides default offset when Fire's state arg matches
// an entry; fallback to default when no match; facing-left negates resolved X;
// nil table or (0,0) default reproduces pre-US-037 behavior.
func TestProjectileWeapon_Fire_StateSpawnOffset(t *testing.T) {
	tests := []struct {
		name           string
		spawnOffsetX16 int
		spawnOffsetY16 int
		stateOffsets   map[int][2]int // nil means do not call SetStateSpawnOffsets
		state          int
		faceDir        animation.FacingDirectionEnum
		fireX16        int
		fireY16        int
		wantSpawnX16   int
		wantSpawnY16   int
	}{
		{
			name:           "state with matching offset uses it",
			spawnOffsetX16: 80,
			spawnOffsetY16: 160,
			stateOffsets:   map[int][2]int{42: {96, 192}},
			state:          42,
			faceDir:        animation.FaceDirectionRight,
			fireX16:        320,
			fireY16:        480,
			wantSpawnX16:   416, // 320 + 96
			wantSpawnY16:   672, // 480 + 192
		},
		{
			name:           "state without entry falls back to default",
			spawnOffsetX16: 80,
			spawnOffsetY16: 160,
			stateOffsets:   map[int][2]int{42: {96, 192}},
			state:          99,
			faceDir:        animation.FaceDirectionRight,
			fireX16:        320,
			fireY16:        480,
			wantSpawnX16:   400, // 320 + 80
			wantSpawnY16:   640, // 480 + 160
		},
		{
			name:           "facing left negates per-state X",
			spawnOffsetX16: 80,
			spawnOffsetY16: 160,
			stateOffsets:   map[int][2]int{42: {96, 192}},
			state:          42,
			faceDir:        animation.FaceDirectionLeft,
			fireX16:        320,
			fireY16:        480,
			wantSpawnX16:   224, // 320 - 96
			wantSpawnY16:   672, // 480 + 192
		},
		{
			name:           "nil state table uses default offset",
			spawnOffsetX16: 80,
			spawnOffsetY16: 160,
			stateOffsets:   nil,
			state:          0,
			faceDir:        animation.FaceDirectionRight,
			fireX16:        320,
			fireY16:        480,
			wantSpawnX16:   400, // 320 + 80
			wantSpawnY16:   640, // 480 + 160
		},
		{
			name:           "(0,0) default + no table = no offset",
			spawnOffsetX16: 0,
			spawnOffsetY16: 0,
			stateOffsets:   nil,
			state:          0,
			faceDir:        animation.FaceDirectionRight,
			fireX16:        320,
			fireY16:        480,
			wantSpawnX16:   320,
			wantSpawnY16:   480,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotX16, gotY16 int
			mock := &mockProjectileManager{
				SpawnProjectileFunc: func(_ string, x16, y16, _, _ int, _ interface{}) {
					gotX16 = x16
					gotY16 = y16
				},
			}

			w := weapon.NewProjectileWeapon("test", 10, "bullet", 100, mock, "", tt.spawnOffsetX16, tt.spawnOffsetY16)
			if tt.stateOffsets != nil {
				w.SetStateSpawnOffsets(tt.stateOffsets)
			}

			w.Fire(tt.fireX16, tt.fireY16, tt.faceDir, body.ShootDirectionStraight, tt.state)

			if gotX16 != tt.wantSpawnX16 {
				t.Errorf("spawnX16: got %d, want %d", gotX16, tt.wantSpawnX16)
			}
			if gotY16 != tt.wantSpawnY16 {
				t.Errorf("spawnY16: got %d, want %d", gotY16, tt.wantSpawnY16)
			}
		})
	}
}
