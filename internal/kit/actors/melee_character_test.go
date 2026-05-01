package kitactors

import (
	"testing"

	meleeengine "github.com/boilerplate/ebiten-template/internal/engine/combat/melee"
)

func TestMeleeCharacter_MeleeController_NilBeforeSet(t *testing.T) {
	mc := NewMeleeCharacter()
	if mc.MeleeController() != nil {
		t.Error("expected MeleeController() to be nil before SetMeleeController")
	}
}

func TestMeleeCharacter_SetMeleeController_RoundTrip(t *testing.T) {
	mc := NewMeleeCharacter()
	ctrl := &meleeengine.Controller{}
	mc.SetMeleeController(ctrl)
	if mc.MeleeController() != ctrl {
		t.Error("MeleeController() did not return the value set via SetMeleeController")
	}
}
