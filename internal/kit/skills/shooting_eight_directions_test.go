package kitskills

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
	"github.com/boilerplate/ebiten-template/internal/engine/mocks"
	"github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
)

func TestShootingSkill_ShootStraight(t *testing.T) {
	var capturedDir body.ShootDirection
	handler := &mocks.MockStateTransitionHandler{
		TransitionToShootingFunc: func(direction body.ShootDirection) {
			capturedDir = direction
		},
	}

	mockWeapon := &mocks.MockWeapon{
		CanFireFunc: func() bool { return true },
		FireFunc: func(x16, y16 int, faceDir animation.FacingDirectionEnum, direction body.ShootDirection, _ int) {
			capturedDir = direction
		},
	}

	mockInventory := &mocks.MockInventory{
		ActiveWeaponFunc: func() combat.Weapon { return mockWeapon },
	}

	s := NewShootingSkill(mockInventory)
	s.SetStateTransitionHandler(handler)

	mockBody := &mockMovableCollidable{
		ObstacleRect: newMockMovableCollidable().ObstacleRect,
	}

	s.HandleInputWithDirection(mockBody, nil, nil, false, false, false, false)

	if capturedDir != body.ShootDirectionStraight {
		t.Errorf("expected ShootDirectionStraight, got %v", capturedDir)
	}
}

func TestShootingSkill_ShootUp(t *testing.T) {
	var capturedDir body.ShootDirection
	handler := &mocks.MockStateTransitionHandler{
		TransitionToShootingFunc: func(direction body.ShootDirection) {
			capturedDir = direction
		},
	}

	mockWeapon := &mocks.MockWeapon{
		CanFireFunc: func() bool { return true },
		FireFunc: func(x16, y16 int, faceDir animation.FacingDirectionEnum, direction body.ShootDirection, _ int) {
			capturedDir = direction
		},
	}

	mockInventory := &mocks.MockInventory{
		ActiveWeaponFunc: func() combat.Weapon { return mockWeapon },
	}

	s := NewShootingSkill(mockInventory)
	s.SetStateTransitionHandler(handler)

	mockBody := &mockMovableCollidable{
		ObstacleRect: newMockMovableCollidable().ObstacleRect,
	}

	s.HandleInputWithDirection(mockBody, nil, nil, true, false, false, false)

	if capturedDir != body.ShootDirectionUp {
		t.Errorf("expected ShootDirectionUp, got %v", capturedDir)
	}
}

func TestShootingSkill_ShootDownAirborne(t *testing.T) {
	var capturedDir body.ShootDirection
	handler := &mocks.MockStateTransitionHandler{
		TransitionToShootingFunc: func(direction body.ShootDirection) {
			capturedDir = direction
		},
	}

	mockWeapon := &mocks.MockWeapon{
		CanFireFunc: func() bool { return true },
		FireFunc: func(x16, y16 int, faceDir animation.FacingDirectionEnum, direction body.ShootDirection, _ int) {
			capturedDir = direction
		},
	}

	mockInventory := &mocks.MockInventory{
		ActiveWeaponFunc: func() combat.Weapon { return mockWeapon },
	}

	s := NewShootingSkill(mockInventory)
	s.SetStateTransitionHandler(handler)

	mockBody := &mockMovableCollidable{
		ObstacleRect: newMockMovableCollidable().ObstacleRect,
	}

	model := &movement.PlatformMovementModel{}
	model.SetOnGround(false)

	s.HandleInputWithDirection(mockBody, model, nil, false, true, false, false)

	if capturedDir != body.ShootDirectionDown {
		t.Errorf("expected ShootDirectionDown, got %v", capturedDir)
	}
}

func TestShootingSkill_ShootDownGrounded_Ignored(t *testing.T) {
	var capturedDir body.ShootDirection
	handler := &mocks.MockStateTransitionHandler{
		TransitionToShootingFunc: func(direction body.ShootDirection) {
			capturedDir = direction
		},
	}

	mockWeapon := &mocks.MockWeapon{
		CanFireFunc: func() bool { return true },
		FireFunc: func(x16, y16 int, faceDir animation.FacingDirectionEnum, direction body.ShootDirection, _ int) {
			capturedDir = direction
		},
	}

	mockInventory := &mocks.MockInventory{
		ActiveWeaponFunc: func() combat.Weapon { return mockWeapon },
	}

	s := NewShootingSkill(mockInventory)
	s.SetStateTransitionHandler(handler)

	mockBody := &mockMovableCollidable{
		ObstacleRect: newMockMovableCollidable().ObstacleRect,
	}

	model := &movement.PlatformMovementModel{}
	model.SetOnGround(true)

	s.HandleInputWithDirection(mockBody, model, nil, false, true, false, false)

	if capturedDir != body.ShootDirectionStraight {
		t.Errorf("expected ShootDirectionStraight (down ignored when grounded), got %v", capturedDir)
	}
}

func TestShootingSkill_DiagonalUpForward(t *testing.T) {
	var capturedDir body.ShootDirection
	handler := &mocks.MockStateTransitionHandler{
		TransitionToShootingFunc: func(direction body.ShootDirection) {
			capturedDir = direction
		},
	}

	mockWeapon := &mocks.MockWeapon{
		CanFireFunc: func() bool { return true },
		FireFunc: func(x16, y16 int, faceDir animation.FacingDirectionEnum, direction body.ShootDirection, _ int) {
			capturedDir = direction
		},
	}

	mockInventory := &mocks.MockInventory{
		ActiveWeaponFunc: func() combat.Weapon { return mockWeapon },
	}

	s := NewShootingSkill(mockInventory)
	s.SetStateTransitionHandler(handler)

	mockBody := &mockMovableCollidable{
		ObstacleRect: newMockMovableCollidable().ObstacleRect,
	}

	s.HandleInputWithDirection(mockBody, nil, nil, true, false, true, false)

	if capturedDir != body.ShootDirectionDiagonalUpForward {
		t.Errorf("expected ShootDirectionDiagonalUpForward, got %v", capturedDir)
	}
}

func TestShootingSkill_DirectionChangeMidShooting(t *testing.T) {
	var transitionCount int
	var directions []body.ShootDirection
	handler := &mocks.MockStateTransitionHandler{
		TransitionToShootingFunc: func(direction body.ShootDirection) {
			transitionCount++
			directions = append(directions, direction)
		},
	}

	mockWeapon := &mocks.MockWeapon{
		CanFireFunc: func() bool { return true },
		FireFunc: func(x16, y16 int, faceDir animation.FacingDirectionEnum, direction body.ShootDirection, _ int) {
		},
	}

	mockInventory := &mocks.MockInventory{
		ActiveWeaponFunc: func() combat.Weapon { return mockWeapon },
	}

	s := NewShootingSkill(mockInventory)
	s.SetStateTransitionHandler(handler)

	mockBody := &mockMovableCollidable{
		ObstacleRect: newMockMovableCollidable().ObstacleRect,
	}

	s.HandleInputWithDirection(mockBody, nil, nil, false, false, false, false)

	if transitionCount != 1 || directions[0] != body.ShootDirectionStraight {
		t.Fatalf("first shot should be straight, got %d transitions, direction=%v", transitionCount, directions)
	}

	s.HandleInputWithDirection(mockBody, nil, nil, true, false, false, false)

	if transitionCount != 2 || directions[1] != body.ShootDirectionUp {
		t.Errorf("direction change should trigger transition to Up, got %d transitions, direction=%v", transitionCount, directions)
	}
}

func TestShootingSkill_ReleaseDirectionalInput(t *testing.T) {
	var directions []body.ShootDirection
	handler := &mocks.MockStateTransitionHandler{
		TransitionToShootingFunc: func(direction body.ShootDirection) {
			directions = append(directions, direction)
		},
	}

	mockWeapon := &mocks.MockWeapon{
		CanFireFunc: func() bool { return true },
		FireFunc: func(x16, y16 int, faceDir animation.FacingDirectionEnum, direction body.ShootDirection, _ int) {
		},
	}

	mockInventory := &mocks.MockInventory{
		ActiveWeaponFunc: func() combat.Weapon { return mockWeapon },
	}

	s := NewShootingSkill(mockInventory)
	s.SetStateTransitionHandler(handler)

	mockBody := &mockMovableCollidable{
		ObstacleRect: newMockMovableCollidable().ObstacleRect,
	}

	s.HandleInputWithDirection(mockBody, nil, nil, true, false, false, false)

	if len(directions) != 1 || directions[0] != body.ShootDirectionUp {
		t.Fatalf("first shot should be up, got %v", directions)
	}

	s.HandleInputWithDirection(mockBody, nil, nil, false, false, false, false)

	if len(directions) != 2 || directions[1] != body.ShootDirectionStraight {
		t.Errorf("releasing up should transition to straight, got %v", directions)
	}
}

func TestShootingSkill_DuckingShooting(t *testing.T) {
	var capturedDir body.ShootDirection
	handler := &mocks.MockStateTransitionHandler{
		TransitionToShootingFunc: func(direction body.ShootDirection) {
			capturedDir = direction
		},
	}

	mockWeapon := &mocks.MockWeapon{
		CanFireFunc: func() bool { return true },
		FireFunc: func(x16, y16 int, faceDir animation.FacingDirectionEnum, direction body.ShootDirection, _ int) {
			capturedDir = direction
		},
	}

	mockInventory := &mocks.MockInventory{
		ActiveWeaponFunc: func() combat.Weapon { return mockWeapon },
	}

	s := NewShootingSkill(mockInventory)
	s.SetStateTransitionHandler(handler)

	mockBody := &mockMovableCollidable{
		ObstacleRect:  newMockMovableCollidable().ObstacleRect,
		isDuckingFunc: func() bool { return true },
	}

	model := &movement.PlatformMovementModel{}
	model.SetOnGround(true)

	s.HandleInputWithDirection(mockBody, model, nil, true, false, false, false)

	if capturedDir != body.ShootDirectionStraight {
		t.Errorf("ducking should only allow straight shooting, got %v", capturedDir)
	}
}
