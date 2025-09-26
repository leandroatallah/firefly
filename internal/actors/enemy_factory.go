package actors

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/systems/physics"
)

// TODO: Should it be a Actor factory?
type EnemyFactory interface {
	Create(enemyType EnemyType) (ActorEntity, error)
}

type EnemyType int

const (
	BlueEnemyType EnemyType = iota
)

type DefaultEnemyFactory struct{}

func NewDefaultEnemyFactory() *DefaultEnemyFactory {
	return &DefaultEnemyFactory{}
}

func (f *DefaultEnemyFactory) Create(enemyType EnemyType) (ActorEntity, error) {
	switch enemyType {
	case BlueEnemyType:
		return NewBlueEnemy(), nil
	default:
		return nil, fmt.Errorf("unknown enemy type")
	}
}

// Blue Enemy
type BlueEnemy struct {
	Character
}

func NewBlueEnemy() *BlueEnemy {
	const (
		frameWidth  = 32
		frameHeight = 32
	)

	var assets spriteAssets
	assets = assets.addSprite(Idle, "assets/blue-enemy.png").
		addSprite(Walk, "assets/blue-enemy.png")

	sprites, err := loadSprites(assets)
	if err != nil {
		log.Fatal(err)
	}

	character := NewCharacter(sprites)

	// TODO: How to define the initial position?
	x, y := 30, 30
	bodyRect := physics.NewRect(x, y, frameWidth, frameHeight)
	collisionRect := physics.NewRect(x, y, frameWidth, frameHeight)

	// TODO: Create a builder with director to automate this process
	enemy := &BlueEnemy{Character: character}
	enemy.SetBody(bodyRect)
	enemy.SetCollisionArea(collisionRect)
	enemy.SetMovementFunc(enemy.HandleMovement)

	return enemy
}

// Blue Enemy - Character Methods
func (e *BlueEnemy) SetBody(rect *physics.Rect) ActorEntity {
	return e.Character.SetBody(rect)
}

func (e *BlueEnemy) SetCollisionArea(rect *physics.Rect) ActorEntity {
	return e.Character.SetCollisionArea(rect)
}

func (e *BlueEnemy) Update(boundaries []physics.Body) error {
	return e.Character.Update(boundaries)
}

func (e *BlueEnemy) Draw(screen *ebiten.Image) {
	e.Character.Draw(screen)
}

func (e *BlueEnemy) HandleMovement() {
	e.OnMoveDown()
}
