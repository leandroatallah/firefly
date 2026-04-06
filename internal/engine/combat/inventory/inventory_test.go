package inventory

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/mocks"
)

func TestActiveWeaponEmptyInventory(t *testing.T) {
	inv := New()
	if inv.ActiveWeapon() != nil {
		t.Error("expected nil for empty inventory, got weapon")
	}
}

func TestAddAndRetrieveWeapon(t *testing.T) {
	inv := New()
	weapon := &mocks.MockWeapon{
		IDFunc: func() string { return "laser" },
	}
	inv.AddWeapon(weapon)

	if inv.ActiveWeapon() != weapon {
		t.Error("expected added weapon to be active")
	}
}

func TestSwitchNextWrapAround(t *testing.T) {
	inv := New()
	w1 := &mocks.MockWeapon{IDFunc: func() string { return "w1" }}
	w2 := &mocks.MockWeapon{IDFunc: func() string { return "w2" }}
	w3 := &mocks.MockWeapon{IDFunc: func() string { return "w3" }}

	inv.AddWeapon(w1)
	inv.AddWeapon(w2)
	inv.AddWeapon(w3)

	// Move to last weapon
	inv.SwitchNext()
	inv.SwitchNext()
	if inv.ActiveWeapon() != w3 {
		t.Error("expected w3 to be active")
	}

	// Wrap around to first
	inv.SwitchNext()
	if inv.ActiveWeapon() != w1 {
		t.Error("expected wrap-around to w1")
	}
}

func TestSwitchPrevWrapAround(t *testing.T) {
	inv := New()
	w1 := &mocks.MockWeapon{IDFunc: func() string { return "w1" }}
	w2 := &mocks.MockWeapon{IDFunc: func() string { return "w2" }}
	w3 := &mocks.MockWeapon{IDFunc: func() string { return "w3" }}

	inv.AddWeapon(w1)
	inv.AddWeapon(w2)
	inv.AddWeapon(w3)

	// At first weapon, go prev
	if inv.ActiveWeapon() != w1 {
		t.Error("expected w1 to be active initially")
	}

	inv.SwitchPrev()
	if inv.ActiveWeapon() != w3 {
		t.Error("expected wrap-around to w3")
	}
}

func TestUnlimitedAmmo(t *testing.T) {
	inv := New()
	weapon := &mocks.MockWeapon{
		IDFunc: func() string { return "laser" },
	}
	inv.AddWeapon(weapon)

	// Unlimited ammo (-1) should return true
	if !inv.HasAmmo("laser") {
		t.Error("expected HasAmmo to return true for unlimited ammo")
	}

	// Consuming ammo should not change unlimited status
	inv.ConsumeAmmo("laser", 10)
	if !inv.HasAmmo("laser") {
		t.Error("expected HasAmmo to still return true after consuming from unlimited")
	}
}

func TestLimitedAmmoConsumption(t *testing.T) {
	inv := New()
	weapon := &mocks.MockWeapon{
		IDFunc: func() string { return "plasma" },
	}
	inv.AddWeapon(weapon)
	inv.SetAmmo("plasma", 5)

	// Consume 2, should have 3 left
	inv.ConsumeAmmo("plasma", 2)
	if !inv.HasAmmo("plasma") {
		t.Error("expected HasAmmo to return true with 3 ammo remaining")
	}

	// Consume 3 more, should have 0 left
	inv.ConsumeAmmo("plasma", 3)
	if inv.HasAmmo("plasma") {
		t.Error("expected HasAmmo to return false with 0 ammo")
	}
}

func TestSwitchToBoundsCheck(t *testing.T) {
	inv := New()
	w1 := &mocks.MockWeapon{IDFunc: func() string { return "w1" }}
	w2 := &mocks.MockWeapon{IDFunc: func() string { return "w2" }}
	w3 := &mocks.MockWeapon{IDFunc: func() string { return "w3" }}

	inv.AddWeapon(w1)
	inv.AddWeapon(w2)
	inv.AddWeapon(w3)

	// Switch to valid index
	inv.SwitchTo(1)
	if inv.ActiveWeapon() != w2 {
		t.Error("expected w2 to be active after SwitchTo(1)")
	}

	// Try to switch to out-of-bounds index (should be no-op)
	inv.SwitchTo(5)
	if inv.ActiveWeapon() != w2 {
		t.Error("expected activeIndex to remain 1 after invalid SwitchTo(5)")
	}
}
