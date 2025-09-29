package enemies

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/actors"
	"github.com/leandroatallah/firefly/internal/systems/physics"
)

type BlueEnemy struct {
	actors.Character
	count int
}

func NewBlueEnemy(x, y int) *BlueEnemy {
	const (
		frameWidth  = 32
		frameHeight = 32
	)

	var assets actors.SpriteAssets
	assets = assets.AddSprite(actors.Idle, "assets/blue-enemy.png").
		AddSprite(actors.Walk, "assets/blue-enemy.png")

	sprites, err := actors.LoadSprites(assets)
	if err != nil {
		log.Fatal(err)
	}

	character := actors.NewCharacter(sprites)

	bodyRect := physics.NewRect(x, y, frameWidth, frameHeight)
	collisionRect := physics.NewRect(x, y, frameWidth, frameHeight)

	// TODO: Create a builder with director to automate this process
	enemy := &BlueEnemy{Character: *character}
	enemy.SetBody(bodyRect)
	// TODO: Move it to the right place (builder)
	enemy.SetSpeedAndMaxSpeed(2, 2)
	enemy.SetCollisionArea(collisionRect)

	return enemy
}

// Character Methods
func (e *BlueEnemy) Update(boundaries []physics.Body) error {
	e.count++

	// Example of movement state change
	if e.count > 200 {
		e.SwitchMovementState(actors.Rand)
	}
	if e.count > 400 {
		e.SwitchMovementState(actors.DumbChase)
	}

	return e.Character.Update(boundaries)
}

func (e *BlueEnemy) Draw(screen *ebiten.Image) {
	e.Character.Draw(screen)
}
