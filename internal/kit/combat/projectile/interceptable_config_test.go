package projectile

import (
	"testing"

	contractsbody "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
)

// AC5 — ProjectileConfig.Interceptable defaults to false, and Manager.Spawn
// produces a body that satisfies body.Projectile and reports Interceptable()==false.
func TestProjectileConfig_DefaultInterceptableIsFalse(t *testing.T) {
	tests := []struct {
		name              string
		cfg               ProjectileConfig
		wantInterceptable bool
	}{
		{
			name:              "zero-value config defaults Interceptable to false",
			cfg:               ProjectileConfig{Width: 2, Height: 1},
			wantInterceptable: false,
		},
		{
			name:              "explicit Interceptable=true propagates to spawned body",
			cfg:               ProjectileConfig{Width: 2, Height: 1, Interceptable: true},
			wantInterceptable: true,
		},
		{
			name:              "explicit Interceptable=false propagates to spawned body",
			cfg:               ProjectileConfig{Width: 2, Height: 1, Interceptable: false},
			wantInterceptable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "zero-value config defaults Interceptable to false" {
				if (ProjectileConfig{}).Interceptable != false {
					t.Errorf("zero-value Interceptable = true, want false")
				}
			}

			mockSpace := &mockBodiesSpace{}
			mgr := NewManager(mockSpace)
			mgr.Spawn(tt.cfg, 100<<4, 50<<4, 0, 0, nil)

			if len(mgr.projectiles) != 1 {
				t.Fatalf("expected 1 projectile, got %d", len(mgr.projectiles))
			}

			p := mgr.projectiles[0]
			proj, ok := p.body.(contractsbody.Projectile)
			if !ok {
				t.Fatalf("spawned body does not satisfy contractsbody.Projectile; the manager must wrap the body so the trait is discoverable")
			}
			if got := proj.Interceptable(); got != tt.wantInterceptable {
				t.Errorf("Interceptable() = %v, want %v", got, tt.wantInterceptable)
			}
		})
	}
}
