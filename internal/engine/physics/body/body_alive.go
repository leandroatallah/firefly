package bodyphysics

import "github.com/leandroatallah/firefly/internal/engine/contracts/body"

// Implements Alive
type AliveBody struct {
	body.Alive

	*Body

	health       int
	maxHealth    int
	invulnerable bool
}

func NewAliveBody(body *Body) *AliveBody {
	if body == nil {
		panic("NewAliveBody: body must not be nil")
	}
	return &AliveBody{Body: body}
}

func (b *AliveBody) Health() int {
	return b.health
}

func (b *AliveBody) MaxHealth() int {
	return b.maxHealth
}

func (b *AliveBody) SetHealth(health int) {
	b.health = health
}

func (b *AliveBody) SetMaxHealth(health int) {
	b.health = health
	b.maxHealth = health
}

func (b *AliveBody) LoseHealth(damage int) {
	b.health = max(b.health-damage, 0)
}
func (b *AliveBody) RestoreHealth(heal int) {
	b.health = min(b.health+heal, b.maxHealth)
}

func (b *AliveBody) Invulnerable() bool {
	return b.invulnerable
}

func (b *AliveBody) SetInvulnerability(value bool) {
	b.invulnerable = value
}
