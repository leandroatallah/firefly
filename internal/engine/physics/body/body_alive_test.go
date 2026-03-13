package body

import (
	"testing"
)

func TestNewAliveBody(t *testing.T) {
	b := NewBody(NewRect(0, 0, 10, 10))
	ab := NewAliveBody(b)

	if ab == nil {
		t.Fatal("NewAliveBody returned nil")
	}
	if ab.Body != b {
		t.Errorf("expected Body to be set; got %v", ab.Body)
	}
}

func TestNewAliveBody_NilBody(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("NewAliveBody did not panic with nil body")
		}
	}()
	NewAliveBody(nil)
}

func TestAliveBody_Health(t *testing.T) {
	ab := NewAliveBody(NewBody(NewRect(0, 0, 10, 10)))
	ab.SetHealth(50)

	if ab.Health() != 50 {
		t.Errorf("expected health 50; got %d", ab.Health())
	}
}

func TestAliveBody_MaxHealth(t *testing.T) {
	ab := NewAliveBody(NewBody(NewRect(0, 0, 10, 10)))
	ab.SetMaxHealth(100)

	if ab.MaxHealth() != 100 {
		t.Errorf("expected maxHealth 100; got %d", ab.MaxHealth())
	}
}

func TestAliveBody_SetHealth(t *testing.T) {
	ab := NewAliveBody(NewBody(NewRect(0, 0, 10, 10)))

	tests := []struct {
		name   string
		health int
		want   int
	}{
		{"positive", 50, 50},
		{"zero", 0, 0},
		{"negative", -10, -10},
		{"above max", 150, 150},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab.SetHealth(tt.health)
			if ab.Health() != tt.want {
				t.Errorf("expected health %d; got %d", tt.want, ab.Health())
			}
		})
	}
}

func TestAliveBody_SetMaxHealth(t *testing.T) {
	ab := NewAliveBody(NewBody(NewRect(0, 0, 10, 10)))

	tests := []struct {
		name      string
		maxHealth int
		wantHealth int
		wantMax   int
	}{
		{"positive", 100, 100, 100},
		{"zero", 0, 0, 0},
		{"large", 1000, 1000, 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab.SetMaxHealth(tt.maxHealth)
			if ab.Health() != tt.wantHealth || ab.MaxHealth() != tt.wantMax {
				t.Errorf("expected health %d, maxHealth %d; got health %d, maxHealth %d",
					tt.wantHealth, tt.wantMax, ab.Health(), ab.MaxHealth())
			}
		})
	}
}

func TestAliveBody_LoseHealth(t *testing.T) {
	tests := []struct {
		name       string
		startHealth int
		damage     int
		wantHealth int
	}{
		{"normal damage", 100, 30, 70},
		{"damage to zero", 50, 50, 0},
		{"excess damage", 30, 50, 0},
		{"no damage", 100, 0, 100},
		{"small damage", 100, 1, 99},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAliveBody(NewBody(NewRect(0, 0, 10, 10)))
			ab.SetHealth(tt.startHealth)

			ab.LoseHealth(tt.damage)

			if ab.Health() != tt.wantHealth {
				t.Errorf("expected health %d after losing %d; got %d",
					tt.wantHealth, tt.damage, ab.Health())
			}
		})
	}
}

func TestAliveBody_RestoreHealth(t *testing.T) {
	tests := []struct {
		name       string
		startHealth int
		maxHealth  int
		heal       int
		wantHealth int
	}{
		{"normal heal", 50, 100, 30, 80},
		{"heal to max", 50, 100, 50, 100},
		{"excess heal", 80, 100, 50, 100},
		{"no heal", 50, 100, 0, 50},
		{"small heal", 50, 100, 1, 51},
		{"heal from zero", 0, 100, 25, 25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAliveBody(NewBody(NewRect(0, 0, 10, 10)))
			ab.SetMaxHealth(tt.maxHealth)
			ab.SetHealth(tt.startHealth)

			ab.RestoreHealth(tt.heal)

			if ab.Health() != tt.wantHealth {
				t.Errorf("expected health %d after healing %d; got %d",
					tt.wantHealth, tt.heal, ab.Health())
			}
		})
	}
}

func TestAliveBody_Invulnerable(t *testing.T) {
	ab := NewAliveBody(NewBody(NewRect(0, 0, 10, 10)))

	if ab.Invulnerable() {
		t.Error("expected invulnerable to be false by default")
	}

	ab.SetInvulnerability(true)
	if !ab.Invulnerable() {
		t.Error("expected invulnerable to be true after setting")
	}

	ab.SetInvulnerability(false)
	if ab.Invulnerable() {
		t.Error("expected invulnerable to be false after clearing")
	}
}

func TestAliveBody_HealthBoundary(t *testing.T) {
	ab := NewAliveBody(NewBody(NewRect(0, 0, 10, 10)))
	ab.SetMaxHealth(100)
	ab.SetHealth(100)

	// Lose all health
	ab.LoseHealth(100)
	if ab.Health() != 0 {
		t.Errorf("expected health 0; got %d", ab.Health())
	}

	// Try to lose more health
	ab.LoseHealth(50)
	if ab.Health() != 0 {
		t.Errorf("expected health to stay at 0; got %d", ab.Health())
	}

	// Restore to max
	ab.RestoreHealth(200)
	if ab.Health() != 100 {
		t.Errorf("expected health capped at max 100; got %d", ab.Health())
	}
}

func TestAliveBody_ChainOperations(t *testing.T) {
	ab := NewAliveBody(NewBody(NewRect(0, 0, 10, 10)))
	ab.SetMaxHealth(100)
	ab.SetHealth(100)

	ab.LoseHealth(30)
	ab.RestoreHealth(10)
	ab.LoseHealth(20)

	if ab.Health() != 60 {
		t.Errorf("expected health 60 after chain operations; got %d", ab.Health())
	}
}
