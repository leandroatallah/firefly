package kitactors

import "testing"

func TestPlayerDeathBehavior_OnDie_SetsHealthToZero(t *testing.T) {
	tests := []struct {
		name          string
		initialHealth int
	}{
		{"health already zero", 0},
		{"health positive", 100},
		{"health negative", -5},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actor := &mockPlatformerActor{health: tc.initialHealth}
			b := NewPlayerDeathBehavior(actor)
			b.OnDie()
			if actor.health != 0 {
				t.Errorf("expected health 0 after OnDie, got %d", actor.health)
			}
		})
	}
}
