package combat_test

import (
	"testing"

	contractscombat "github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
)

// stubDamageable implements contractscombat.Damageable for the interface
// compile-time / identity checks below. Covers AC1.
type stubDamageable struct {
	got int
}

func (s *stubDamageable) TakeDamage(amount int) { s.got = amount }

// stubDestructible implements contractscombat.Destructible. Covers AC7.
type stubDestructible struct {
	got       int
	destroyed bool
}

func (s *stubDestructible) TakeDamage(amount int) { s.got = amount }
func (s *stubDestructible) IsDestroyed() bool     { return s.destroyed }

// TestDamageable_Interface asserts the Damageable interface exists with the
// exact shape required by AC1: a single method TakeDamage(int).
func TestDamageable_Interface(t *testing.T) {
	var d contractscombat.Damageable = &stubDamageable{}
	d.TakeDamage(7)
	got := d.(*stubDamageable).got
	if got != 7 {
		t.Errorf("TakeDamage stored = %d, want 7", got)
	}
}

// TestDestructible_Interface asserts the Destructible interface is exactly
// Damageable + IsDestroyed() bool (AC7). A Destructible is assignable to
// Damageable (embedding relationship).
func TestDestructible_Interface(t *testing.T) {
	s := &stubDestructible{destroyed: false}
	var d contractscombat.Destructible = s
	d.TakeDamage(4)
	if s.got != 4 {
		t.Errorf("TakeDamage stored = %d, want 4", s.got)
	}
	if d.IsDestroyed() {
		t.Error("IsDestroyed = true, want false")
	}
	s.destroyed = true
	if !d.IsDestroyed() {
		t.Error("IsDestroyed = false after flip, want true")
	}

	// Destructible is a Damageable.
	var asDamageable contractscombat.Damageable = d
	asDamageable.TakeDamage(1)
}
