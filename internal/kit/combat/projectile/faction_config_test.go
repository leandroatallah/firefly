package projectile

import (
	"encoding/json"
	"testing"

	enginecombat "github.com/boilerplate/ebiten-template/internal/kit/combat"
)

// TestProjectileConfig_FactionField covers AC4: ProjectileConfig has a new
// Faction field of type combat.Faction with zero-value == FactionNeutral.
// JSON round-trips via the `faction` tag.
func TestProjectileConfig_FactionField(t *testing.T) {
	tests := []struct {
		name        string
		cfg         ProjectileConfig
		wantFaction enginecombat.Faction
	}{
		{
			name:        "zero-value defaults to FactionNeutral",
			cfg:         ProjectileConfig{},
			wantFaction: enginecombat.FactionNeutral,
		},
		{
			name:        "faction player stored",
			cfg:         ProjectileConfig{Faction: enginecombat.FactionPlayer},
			wantFaction: enginecombat.FactionPlayer,
		},
		{
			name:        "faction enemy stored",
			cfg:         ProjectileConfig{Faction: enginecombat.FactionEnemy},
			wantFaction: enginecombat.FactionEnemy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cfg.Faction != tt.wantFaction {
				t.Errorf("Faction = %v, want %v", tt.cfg.Faction, tt.wantFaction)
			}
		})
	}
}

// TestProjectileConfig_FactionJSONRoundTrip verifies JSON round-trip for the
// new Faction field under the documented `faction` tag.
func TestProjectileConfig_FactionJSONRoundTrip(t *testing.T) {
	orig := ProjectileConfig{
		Width:   2,
		Height:  1,
		Damage:  10,
		Faction: enginecombat.FactionEnemy,
	}
	data, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}
	var got ProjectileConfig
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}
	if got.Faction != orig.Faction {
		t.Errorf("round-trip Faction = %v, want %v", got.Faction, orig.Faction)
	}
}

// TestManager_Spawn_PropagatesFaction covers AC4 post-condition: Manager.Spawn
// copies ProjectileConfig.Faction into the internal projectile.faction field.
// Observed by triggering a hit against a Damageable of the same faction and
// asserting it is ignored (faction gate triggered).
func TestManager_Spawn_PropagatesFaction(t *testing.T) {
	mockSpace := &mockBodiesSpace{
		tilemapProvider: &mockTilemapDimensionsProvider{width: 10000, height: 10000},
	}
	mgr := NewManager(mockSpace)

	cfg := ProjectileConfig{
		Width:   2,
		Height:  1,
		Damage:  10,
		Faction: enginecombat.FactionPlayer,
	}
	mgr.Spawn(cfg, 100<<4, 100<<4, 0, 0, nil)
	if len(mgr.projectiles) != 1 {
		t.Fatalf("expected 1 projectile, got %d", len(mgr.projectiles))
	}
	p := mgr.projectiles[0]

	// Same-faction target: if faction was propagated, TakeDamage must NOT fire.
	sameFactionTarget := &fakeDamageable{faction: enginecombat.FactionPlayer}
	other := fakeCollidableWithOwner(sameFactionTarget)

	p.OnTouch(other)

	if got := len(sameFactionTarget.takeDamageCalls); got != 0 {
		t.Errorf("same-faction TakeDamage count = %d, want 0 (proves Spawn propagated Faction)", got)
	}
}

// TestManager_Spawn_PropagatesDamage covers AC4 post-condition for Damage:
// Manager.Spawn copies ProjectileConfig.Damage into projectile.damage.
func TestManager_Spawn_PropagatesDamage(t *testing.T) {
	mockSpace := &mockBodiesSpace{
		tilemapProvider: &mockTilemapDimensionsProvider{width: 10000, height: 10000},
	}
	mgr := NewManager(mockSpace)

	cfg := ProjectileConfig{
		Width:   2,
		Height:  1,
		Damage:  23,
		Faction: enginecombat.FactionPlayer,
	}
	mgr.Spawn(cfg, 100<<4, 100<<4, 0, 0, nil)
	p := mgr.projectiles[0]

	target := &fakeDamageable{faction: enginecombat.FactionEnemy}
	other := fakeCollidableWithOwner(target)

	p.OnTouch(other)

	if got := len(target.takeDamageCalls); got != 1 {
		t.Fatalf("TakeDamage count = %d, want 1", got)
	}
	if got := target.takeDamageCalls[0]; got != 23 {
		t.Errorf("TakeDamage amount = %d, want 23 (proves Spawn propagated Damage)", got)
	}
}
