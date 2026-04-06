package gameplayer_test

import (
	"testing"

	gameplayer "github.com/boilerplate/ebiten-template/internal/game/entity/actors/player"
)

type mockProjectileManager struct{}

func (m *mockProjectileManager) SpawnProjectile(projectileType string, x16, y16, vx16, vy16 int, owner interface{}) {
}

func TestNewClimberInventory(t *testing.T) {
	inv := gameplayer.NewClimberInventory(&mockProjectileManager{})

	if inv.ActiveWeapon().ID() != "light_blaster" {
		t.Fatalf("expected light_blaster, got %s", inv.ActiveWeapon().ID())
	}
	if !inv.ActiveWeapon().CanFire() {
		t.Fatal("expected CanFire() == true")
	}

	inv.SwitchNext()
	if inv.ActiveWeapon().ID() != "heavy_cannon" {
		t.Fatalf("expected heavy_cannon, got %s", inv.ActiveWeapon().ID())
	}
}
