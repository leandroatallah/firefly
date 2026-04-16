package input

import "github.com/hajimehoshi/ebiten/v2"

type PlayerCommands struct {
	Up         bool
	Down       bool
	Left       bool
	Right      bool
	Shoot      bool
	Melee      bool
	Jump       bool
	Dash       bool
	Confirm    bool
	Cancel     bool
	WeaponNext bool
	WeaponPrev bool
}

// ReadPlayerCommands returns the default keyboard mapping.
// Swappable via CommandsReader for game-layer overrides.
func ReadPlayerCommands() PlayerCommands {
	return PlayerCommands{
		Up:         isKeyPressed(ebiten.KeyUp) || isKeyPressed(ebiten.KeyW),
		Down:       isKeyPressed(ebiten.KeyDown) || isKeyPressed(ebiten.KeyS),
		Left:       isKeyPressed(ebiten.KeyLeft) || isKeyPressed(ebiten.KeyA),
		Right:      isKeyPressed(ebiten.KeyRight) || isKeyPressed(ebiten.KeyD),
		Shoot:      isKeyPressed(ebiten.KeyX),
		Melee:      isKeyPressed(ebiten.KeyZ),
		Jump:       isKeyPressed(ebiten.KeySpace),
		Dash:       isKeyPressed(ebiten.KeyShift),
		Confirm:    isKeyPressed(ebiten.KeyEnter),
		Cancel:     isKeyPressed(ebiten.KeyEscape),
		WeaponNext: isKeyPressed(ebiten.KeyE),
		WeaponPrev: isKeyPressed(ebiten.KeyQ),
	}
}

// Swappable function var: allows injection in tests and game-layer overrides
//
//nolint:gochecknoglobals
var CommandsReader func() PlayerCommands = ReadPlayerCommands
