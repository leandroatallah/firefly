package kitactors

import meleeengine "github.com/boilerplate/ebiten-template/internal/engine/combat/melee"

// MeleeCharacter is a reusable trait that holds a melee Controller and provides accessor methods.
type MeleeCharacter struct {
	melee *meleeengine.Controller
}

// NewMeleeCharacter creates a new MeleeCharacter with melee initialized to nil.
func NewMeleeCharacter() *MeleeCharacter {
	return &MeleeCharacter{melee: nil}
}

// MeleeController returns the melee controller (may be nil).
func (m *MeleeCharacter) MeleeController() *meleeengine.Controller {
	return m.melee
}

// SetMeleeController assigns the controller field.
func (m *MeleeCharacter) SetMeleeController(c *meleeengine.Controller) {
	m.melee = c
}
