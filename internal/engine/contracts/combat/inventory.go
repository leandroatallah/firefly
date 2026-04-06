package combat

// Inventory manages a collection of weapons with ammo tracking.
type Inventory interface {
	AddWeapon(weapon Weapon)
	ActiveWeapon() Weapon
	SwitchNext()
	SwitchPrev()
	SwitchTo(index int)
	HasAmmo(weaponID string) bool
	ConsumeAmmo(weaponID string, amount int)
	SetAmmo(weaponID string, amount int)
}
