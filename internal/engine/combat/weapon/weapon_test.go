package weapon_test

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/combat/weapon"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
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

			w := weapon.NewProjectileWeapon("test", 10, "bullet", 100, mock)

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
	w := weapon.NewProjectileWeapon("test", 3, "bullet", 100, &mockProjectileManager{})

	w.Fire(0, 0, animation.FaceDirectionRight, body.ShootDirectionStraight)
	w.Update()

	if w.Cooldown() != 2 {
		t.Errorf("Cooldown() after one Update: got %d, want 2", w.Cooldown())
	}
}

func TestProjectileWeapon_CanFire_TrueAfterCooldownExpires(t *testing.T) {
	w := weapon.NewProjectileWeapon("test", 2, "bullet", 100, &mockProjectileManager{})

	w.Fire(0, 0, animation.FaceDirectionRight, body.ShootDirectionStraight)
	w.Update()
	w.Update()

	if !w.CanFire() {
		t.Error("CanFire() should return true after cooldown expires")
	}
}

func TestProjectileWeapon_SetCooldown(t *testing.T) {
	w := weapon.NewProjectileWeapon("test", 10, "bullet", 100, &mockProjectileManager{})
	w.SetCooldown(5)

	if w.Cooldown() != 5 {
		t.Errorf("Cooldown() after SetCooldown(5): got %d, want 5", w.Cooldown())
	}
}

func TestProjectileWeapon_ID(t *testing.T) {
	w := weapon.NewProjectileWeapon("my-weapon", 10, "bullet", 100, &mockProjectileManager{})
	if w.ID() != "my-weapon" {
		t.Errorf("ID(): got %q, want %q", w.ID(), "my-weapon")
	}
}
