package combat_test

import (
	"testing"

	kitcombat "github.com/boilerplate/ebiten-template/internal/kit/combat"
)

// TestFaction_Constants asserts the Faction type and its three canonical
// values are defined per AC3 / SPEC 2.1.
func TestFaction_Constants(t *testing.T) {
	tests := []struct {
		name string
		got  kitcombat.Faction
		want kitcombat.Faction
	}{
		{"neutral is zero value", kitcombat.Faction(0), kitcombat.FactionNeutral},
		{"player is non-neutral", kitcombat.FactionPlayer, kitcombat.FactionPlayer},
		{"enemy is non-neutral", kitcombat.FactionEnemy, kitcombat.FactionEnemy},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}

	// All three must be mutually distinct.
	if kitcombat.FactionNeutral == kitcombat.FactionPlayer ||
		kitcombat.FactionNeutral == kitcombat.FactionEnemy ||
		kitcombat.FactionPlayer == kitcombat.FactionEnemy {
		t.Error("Faction constants must be mutually distinct")
	}
}
