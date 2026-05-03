package kitskills

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
	"github.com/boilerplate/ebiten-template/internal/engine/input"
	"github.com/boilerplate/ebiten-template/internal/engine/mocks"
)

func TestShootingSkill_FireDelegatesToActiveWeapon(t *testing.T) {
	fireCalled := false
	mockWeapon := &mocks.MockWeapon{
		CanFireFunc: func() bool { return true },
		FireFunc: func(x16, y16 int, faceDir animation.FacingDirectionEnum, direction body.ShootDirection, _ int) {
			fireCalled = true
		},
	}

	mockInventory := &mocks.MockInventory{
		ActiveWeaponFunc: func() combat.Weapon { return mockWeapon },
	}

	s := NewShootingSkill(mockInventory)
	actor := &mockMovableCollidable{
		ObstacleRect: newMockMovableCollidable().ObstacleRect,
	}

	// Inject shoot command
	oldReader := input.CommandsReader
	defer func() { input.CommandsReader = oldReader }()
	input.CommandsReader = func() input.PlayerCommands {
		return input.PlayerCommands{Shoot: true}
	}

	s.HandleInput(actor, nil, nil)

	if !fireCalled {
		t.Error("expected weapon.Fire() to be called, but it was not")
	}
}

func TestShootingSkill_NoFireWhenWeaponOnCooldown(t *testing.T) {
	fireCalled := false
	mockWeapon := &mocks.MockWeapon{
		CanFireFunc: func() bool { return false },
		FireFunc: func(x16, y16 int, faceDir animation.FacingDirectionEnum, direction body.ShootDirection, _ int) {
			fireCalled = true
		},
	}

	mockInventory := &mocks.MockInventory{
		ActiveWeaponFunc: func() combat.Weapon { return mockWeapon },
	}

	s := NewShootingSkill(mockInventory)
	actor := &mockMovableCollidable{
		ObstacleRect: newMockMovableCollidable().ObstacleRect,
	}

	// Inject shoot command
	oldReader := input.CommandsReader
	defer func() { input.CommandsReader = oldReader }()
	input.CommandsReader = func() input.PlayerCommands {
		return input.PlayerCommands{Shoot: true}
	}

	s.HandleInput(actor, nil, nil)

	if fireCalled {
		t.Error("expected weapon.Fire() NOT to be called when CanFire() is false, but it was")
	}
}

func TestShootingSkill_NoFireWhenInventoryEmpty(t *testing.T) {
	mockInventory := &mocks.MockInventory{
		ActiveWeaponFunc: func() combat.Weapon { return nil },
	}

	s := NewShootingSkill(mockInventory)
	actor := &mockMovableCollidable{
		ObstacleRect: newMockMovableCollidable().ObstacleRect,
	}

	// Inject shoot command
	oldReader := input.CommandsReader
	defer func() { input.CommandsReader = oldReader }()
	input.CommandsReader = func() input.PlayerCommands {
		return input.PlayerCommands{Shoot: true}
	}

	// Should not panic
	s.HandleInput(actor, nil, nil)
}

func TestShootingSkill_WeaponSwitchingOnInput(t *testing.T) {
	switchNextCalled := false
	switchPrevCalled := false

	mockInventory := &mocks.MockInventory{
		ActiveWeaponFunc: func() combat.Weapon { return nil },
		SwitchNextFunc: func() {
			switchNextCalled = true
		},
		SwitchPrevFunc: func() {
			switchPrevCalled = true
		},
	}

	s := NewShootingSkill(mockInventory)
	actor := &mockMovableCollidable{
		ObstacleRect: newMockMovableCollidable().ObstacleRect,
	}

	// Test WeaponNext
	oldReader := input.CommandsReader
	defer func() { input.CommandsReader = oldReader }()
	input.CommandsReader = func() input.PlayerCommands {
		return input.PlayerCommands{WeaponNext: true}
	}

	s.HandleInput(actor, nil, nil)

	if !switchNextCalled {
		t.Error("expected inv.SwitchNext() to be called on WeaponNext input, but it was not")
	}

	// Test WeaponPrev
	switchNextCalled = false
	input.CommandsReader = func() input.PlayerCommands {
		return input.PlayerCommands{WeaponPrev: true}
	}

	s.HandleInput(actor, nil, nil)

	if !switchPrevCalled {
		t.Error("expected inv.SwitchPrev() to be called on WeaponPrev input, but it was not")
	}
}

func TestShootingSkill_UpdateHandlesShootRelease(t *testing.T) {
	transitionFromShootingCalled := false
	handler := &mocks.MockStateTransitionHandler{
		TransitionFromShootingFunc: func() {
			transitionFromShootingCalled = true
		},
	}

	mockInventory := &mocks.MockInventory{
		ActiveWeaponFunc: func() combat.Weapon { return nil },
	}

	s := NewShootingSkill(mockInventory)
	s.SetStateTransitionHandler(handler)

	actor := &mockMovableCollidable{
		ObstacleRect: newMockMovableCollidable().ObstacleRect,
	}

	oldReader := input.CommandsReader
	defer func() { input.CommandsReader = oldReader }()

	// Simulate shoot held
	input.CommandsReader = func() input.PlayerCommands {
		return input.PlayerCommands{Shoot: true}
	}
	s.Update(actor, nil)

	// Simulate shoot released
	input.CommandsReader = func() input.PlayerCommands {
		return input.PlayerCommands{Shoot: false}
	}
	s.Update(actor, nil)

	if !transitionFromShootingCalled {
		t.Error("expected TransitionFromShooting() to be called on shoot release, but it was not")
	}
}

func TestShootingSkill_IsActive_ReflectsShootHeld(t *testing.T) {
	mockInventory := &mocks.MockInventory{
		ActiveWeaponFunc: func() combat.Weapon { return nil },
	}
	s := NewShootingSkill(mockInventory)

	oldReader := input.CommandsReader
	defer func() { input.CommandsReader = oldReader }()

	shoot := false
	input.CommandsReader = func() input.PlayerCommands {
		return input.PlayerCommands{Shoot: shoot}
	}

	if s.IsActive() {
		t.Error("expected IsActive() false before any Update")
	}

	shoot = true
	s.Update(nil, nil)
	if !s.IsActive() {
		t.Error("expected IsActive() true when shoot is held")
	}

	shoot = false
	s.Update(nil, nil)
	if s.IsActive() {
		t.Error("expected IsActive() false after shoot released")
	}
}
