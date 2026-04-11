package gameplayer_test

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/mocks"
	gameplayer "github.com/boilerplate/ebiten-template/internal/game/entity/actors/player"
)

type mockProjectileManager struct{}

func (m *mockProjectileManager) SpawnProjectile(projectileType string, x16, y16, vx16, vy16 int, owner interface{}) {
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
	inv.ActiveWeapon().Fire(160, 160, animation.FaceDirectionRight, body.ShootDirectionStraight)
	if !spawnPuffCalled {
		t.Error("expected SpawnPuff to be called for light_blaster")
	}

	inv.SwitchNext()
	if inv.ActiveWeapon().ID() != "heavy_cannon" {
		t.Fatalf("expected heavy_cannon, got %s", inv.ActiveWeapon().ID())
	}

	spawnPuffCalled = false
	inv.ActiveWeapon().Fire(160, 160, animation.FaceDirectionRight, body.ShootDirectionStraight)
	if !spawnPuffCalled {
		t.Error("expected SpawnPuff to be called for heavy_cannon")
	}
}
