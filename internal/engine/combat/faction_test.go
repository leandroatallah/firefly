package combat_test

import (
	"testing"

	enginecombat "github.com/boilerplate/ebiten-template/internal/engine/combat"
)

// TestFaction_Constants asserts the Faction type and its three canonical
// values are defined per AC3 / SPEC 2.1.
func TestFaction_Constants(t *testing.T) {
	tests := []struct {
		name string
		got  enginecombat.Faction
		want enginecombat.Faction
	}{
		{"neutral is zero value", enginecombat.Faction(0), enginecombat.FactionNeutral},
		{"player is non-neutral", enginecombat.FactionPlayer, enginecombat.FactionPlayer},
		{"enemy is non-neutral", enginecombat.FactionEnemy, enginecombat.FactionEnemy},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}

	// All three must be mutually distinct.
	if enginecombat.FactionNeutral == enginecombat.FactionPlayer ||
		enginecombat.FactionNeutral == enginecombat.FactionEnemy ||
		enginecombat.FactionPlayer == enginecombat.FactionEnemy {
		t.Error("Faction constants must be mutually distinct")
	}
}
