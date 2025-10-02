package physics

type Alive interface {
	Health() int
	MaxHealth() int
	SetHealth(health int)
	SetMaxHealth(health int)
	LoseHealth(damage int)
	RestoreHealth(heal int)
	Invulnerable() bool
	SetInvulnerable(value bool)
}

type AliveBody struct {
	health    int
	maxHealth int
}

func NewAliveBody() *AliveBody {
	return &AliveBody{}
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
