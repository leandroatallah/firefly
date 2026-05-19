// Red-Phase tests for story 059-thin-game-phase-scenes [AC-7].
// SPEC.md §3 introduces the generic player builder
//
//	BuildPlayer[T actors.ActorEntity](p T, deps PlayerDeps) (T, error)
//
// which applies skills (when SpriteData != nil), then optionally injects
// Inventory and MeleeWeapon when non-nil.
//
// These tests fail until the kitbuilder package exposes BuildPlayer and
// PlayerDeps with the documented semantics.
package kitbuilder

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/kit/combat/weapon"
)

func TestBuildPlayer_NilInventoryAndMelee_SkipsBothWires(t *testing.T) {
	p := newMockPlayerWithWiring()

	got, err := BuildPlayer(p, PlayerDeps{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != p {
		t.Fatal("expected BuildPlayer to return the same player instance")
	}
	if p.setInventoryCalls != 0 {
		t.Fatalf("expected SetInventory not called when Inventory==nil, got %d", p.setInventoryCalls)
	}
	if p.setMeleeCalls != 0 {
		t.Fatalf("expected SetMelee not called when MeleeWeapon==nil, got %d", p.setMeleeCalls)
	}
}

func TestBuildPlayer_NilMelee_AppliesOnlyInventory(t *testing.T) {
	p := newMockPlayerWithWiring()
	inv := &stubInventory{id: "inv-1"}

	if _, err := BuildPlayer(p, PlayerDeps{Inventory: inv}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.setInventoryCalls != 1 {
		t.Fatalf("expected SetInventory called once, got %d", p.setInventoryCalls)
	}
	if p.lastInventory != inv {
		t.Fatalf("expected SetInventory invoked with the same inventory pointer")
	}
	if p.setMeleeCalls != 0 {
		t.Fatalf("expected SetMelee not called when MeleeWeapon==nil, got %d", p.setMeleeCalls)
	}
}

func TestBuildPlayer_NilInventory_AppliesOnlyMelee(t *testing.T) {
	p := newMockPlayerWithWiring()
	wpn := weapon.NewMeleeWeapon("test", 0, 0, nil)

	if _, err := BuildPlayer(p, PlayerDeps{MeleeWeapon: wpn}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.setInventoryCalls != 0 {
		t.Fatalf("expected SetInventory not called, got %d", p.setInventoryCalls)
	}
	if p.setMeleeCalls != 1 {
		t.Fatalf("expected SetMelee called once, got %d", p.setMeleeCalls)
	}
	if p.lastMelee != wpn {
		t.Fatal("expected SetMelee invoked with the same weapon pointer")
	}
}

func TestBuildPlayer_BothNonNil_AppliesBoth(t *testing.T) {
	p := newMockPlayerWithWiring()
	inv := &stubInventory{id: "inv-2"}
	wpn := weapon.NewMeleeWeapon("test", 0, 0, nil)

	if _, err := BuildPlayer(p, PlayerDeps{Inventory: inv, MeleeWeapon: wpn}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.setInventoryCalls != 1 {
		t.Fatalf("expected SetInventory called once, got %d", p.setInventoryCalls)
	}
	if p.setMeleeCalls != 1 {
		t.Fatalf("expected SetMelee called once, got %d", p.setMeleeCalls)
	}
}

func TestBuildPlayer_PlayerWithoutWiring_NoOp(t *testing.T) {
	// A player that does NOT implement the internal playerWiring interface
	// must be returned unchanged with no error and no panic, even when
	// inventory and melee are non-nil.
	p := newMockPlayerNoWiring()
	inv := &stubInventory{id: "should-not-be-used"}
	wpn := weapon.NewMeleeWeapon("test", 0, 0, nil)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("expected no panic, got: %v", r)
		}
	}()

	got, err := BuildPlayer(p, PlayerDeps{Inventory: inv, MeleeWeapon: wpn})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != p {
		t.Fatal("expected BuildPlayer to return the same player instance")
	}
}

func TestBuildPlayer_WireStateInvokedWhenSet(t *testing.T) {
	p := newMockPlayerWithWiring()

	calls := 0
	deps := PlayerDeps{
		WireState: func(*actors.Character) { calls++ },
	}

	if _, err := BuildPlayer(p, deps); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected WireState callback invoked once, got %d", calls)
	}
}
