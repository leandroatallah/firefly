package actors_test

import (
	"testing"

	contractscombat "github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
	"github.com/boilerplate/ebiten-template/internal/engine/render/sprites"
	"github.com/hajimehoshi/ebiten/v2"
)

func newTestCharacter(t *testing.T) *actors.Character {
	t.Helper()
	img := ebiten.NewImage(1, 1)
	sMap := sprites.SpriteMap{
		actors.Idle:    &sprites.Sprite{Image: img},
		actors.Walking: &sprites.Sprite{Image: img},
		actors.Hurted:  &sprites.Sprite{Image: img},
	}
	rect := bodyphysics.NewRect(0, 0, 16, 16)
	c := actors.NewCharacter(sMap, rect)
	c.SetMaxHealth(100)
	c.SetHealth(100)
	return c
}

// TestCharacter_TakeDamageDelegatesToHurt covers AC5 + AC6: Character
// satisfies the Damageable contract via a TakeDamage adapter that delegates to
// the existing Hurt path. Invulnerability guard is preserved (AC10 last
// bullet).
func TestCharacter_TakeDamageDelegatesToHurt(t *testing.T) {
	tests := []struct {
		name             string
		initialHealth    int
		invulnerable     bool
		damage           int
		wantHealthAfter  int
		wantState        actors.ActorStateEnum
		wantInvulnerable bool
	}{
		{
			name:             "damage applies when not invulnerable",
			initialHealth:    100,
			invulnerable:     false,
			damage:           25,
			wantHealthAfter:  75,
			wantState:        actors.Hurted,
			wantInvulnerable: true,
		},
		{
			name:            "damage ignored when invulnerable",
			initialHealth:   100,
			invulnerable:    true,
			damage:          25,
			wantHealthAfter: 100,
			// State is not forced to Hurted; it remains Idle.
			wantState:        actors.Idle,
			wantInvulnerable: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestCharacter(t)
			c.SetHealth(tt.initialHealth)
			c.SetInvulnerability(tt.invulnerable)

			// Assert Character satisfies the Damageable contract (AC5).
			var d contractscombat.Damageable = c
			d.TakeDamage(tt.damage)

			if got := c.Health(); got != tt.wantHealthAfter {
				t.Errorf("Health = %d, want %d", got, tt.wantHealthAfter)
			}
			if got := c.State(); got != tt.wantState {
				t.Errorf("State = %v, want %v", got, tt.wantState)
			}
			if got := c.Invulnerable(); got != tt.wantInvulnerable {
				t.Errorf("Invulnerable = %v, want %v", got, tt.wantInvulnerable)
			}
		})
	}
}

// TestCharacter_FactionAccessors covers the Faction field + accessors required
// by SPEC 2.6 / AC6. Default is FactionNeutral; SetFaction updates the value.
func TestCharacter_FactionAccessors(t *testing.T) {
	tests := []struct {
		name string
		set  *contractscombat.Faction // nil means "do not call SetFaction"
		want contractscombat.Faction
	}{
		{
			name: "default is FactionNeutral",
			set:  nil,
			want: contractscombat.FactionNeutral,
		},
		{
			name: "SetFaction(Player) round-trips",
			set:  ptrFaction(contractscombat.FactionPlayer),
			want: contractscombat.FactionPlayer,
		},
		{
			name: "SetFaction(Enemy) round-trips",
			set:  ptrFaction(contractscombat.FactionEnemy),
			want: contractscombat.FactionEnemy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestCharacter(t)
			if tt.set != nil {
				c.SetFaction(*tt.set)
			}
			if got := c.Faction(); got != tt.want {
				t.Errorf("Faction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ptrFaction(f contractscombat.Faction) *contractscombat.Faction { return &f }
