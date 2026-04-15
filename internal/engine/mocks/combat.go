package mocks

import (
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
)

// MockWeapon is a shared mock for the Weapon interface.
type MockWeapon struct {
	IDFunc          func() string
	FireFunc        func(x16, y16 int, faceDir animation.FacingDirectionEnum, direction body.ShootDirection, state int)
	CanFireFunc     func() bool
	UpdateFunc      func()
	CooldownFunc    func() int
	SetCooldownFunc func(frames int)
}

func (m *MockWeapon) ID() string {
	if m.IDFunc != nil {
		return m.IDFunc()
	}
	return ""
}

func (m *MockWeapon) Fire(x16, y16 int, faceDir animation.FacingDirectionEnum, direction body.ShootDirection, state int) {
	if m.FireFunc != nil {
		m.FireFunc(x16, y16, faceDir, direction, state)
	}
}

func (m *MockWeapon) CanFire() bool {
	if m.CanFireFunc != nil {
		return m.CanFireFunc()
	}
	return false
}

func (m *MockWeapon) Update() {
	if m.UpdateFunc != nil {
		m.UpdateFunc()
	}
}

func (m *MockWeapon) Cooldown() int {
	if m.CooldownFunc != nil {
		return m.CooldownFunc()
	}
	return 0
}

func (m *MockWeapon) SetCooldown(frames int) {
	if m.SetCooldownFunc != nil {
		m.SetCooldownFunc(frames)
	}
}

// MockInventory is a shared mock for the Inventory interface.
type MockInventory struct {
	AddWeaponFunc    func(weapon combat.Weapon)
	ActiveWeaponFunc func() combat.Weapon
	SwitchNextFunc   func()
	SwitchPrevFunc   func()
	SwitchToFunc     func(index int)
	HasAmmoFunc      func(weaponID string) bool
	ConsumeAmmoFunc  func(weaponID string, amount int)
	SetAmmoFunc      func(weaponID string, amount int)
	UpdateFunc       func()
}

func (m *MockInventory) AddWeapon(weapon combat.Weapon) {
	if m.AddWeaponFunc != nil {
		m.AddWeaponFunc(weapon)
	}
}

func (m *MockInventory) ActiveWeapon() combat.Weapon {
	if m.ActiveWeaponFunc != nil {
		return m.ActiveWeaponFunc()
	}
	return nil
}

func (m *MockInventory) SwitchNext() {
	if m.SwitchNextFunc != nil {
		m.SwitchNextFunc()
	}
}

func (m *MockInventory) SwitchPrev() {
	if m.SwitchPrevFunc != nil {
		m.SwitchPrevFunc()
	}
}

func (m *MockInventory) SwitchTo(index int) {
	if m.SwitchToFunc != nil {
		m.SwitchToFunc(index)
	}
}

func (m *MockInventory) HasAmmo(weaponID string) bool {
	if m.HasAmmoFunc != nil {
		return m.HasAmmoFunc(weaponID)
	}
	return false
}

func (m *MockInventory) ConsumeAmmo(weaponID string, amount int) {
	if m.ConsumeAmmoFunc != nil {
		m.ConsumeAmmoFunc(weaponID, amount)
	}
}

func (m *MockInventory) SetAmmo(weaponID string, amount int) {
	if m.SetAmmoFunc != nil {
		m.SetAmmoFunc(weaponID, amount)
	}
}

func (m *MockInventory) Update() {
	if m.UpdateFunc != nil {
		m.UpdateFunc()
	}
}
