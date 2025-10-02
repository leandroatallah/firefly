package physics

type Alive interface {
	Health() int
	MaxHealth() int
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

func (b *AliveBody) LoseHealth(damage int)  {}
func (b *AliveBody) RestoreHealth(heal int) {}
