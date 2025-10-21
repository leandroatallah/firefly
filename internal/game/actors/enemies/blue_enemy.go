package gameenemies

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
	"github.com/leandroatallah/firefly/internal/engine/systems/sprites"
)

type BlueEnemy struct {
	actors.Character
	count int
}

func NewBlueEnemy() *BlueEnemy {
	// TODO: Set logic to initial position
	x, y := 0, 0
	const (
		frameWidth  = 32
		frameHeight = 32
	)

	var assets sprites.SpriteAssets
	assets = assets.AddSprite(actors.Idle, "assets/images/blue-enemy.png").
		AddSprite(actors.Walk, "assets/images/blue-enemy.png")

	sprites, err := sprites.LoadSprites(assets)
	if err != nil {
		log.Fatal(err)
	}

	character := actors.NewCharacter(sprites, 8)

	bodyRect := physics.NewRect(x, y, frameWidth, frameHeight)
	collisionRect := physics.NewRect(x, y, frameWidth, frameHeight)

	// TODO: Create a builder with director to automate this process
	enemy := &BlueEnemy{Character: *character}
	enemy.SetBody(bodyRect)
	// TODO: Move it to the right place (builder)
	enemy.SetSpeedAndMaxSpeed(2, 2)
	enemy.SetCollisionArea(collisionRect)
	enemy.PhysicsBody.SetTouchable(enemy)

	return enemy
}

// Character Methods
func (e *BlueEnemy) Update(space body.BodiesSpace) error {
	e.count++
	return e.Character.Update(space)
}

func (e *BlueEnemy) Draw(screen *ebiten.Image) {
	e.Character.Draw(screen)
}

func (e *BlueEnemy) OnTouch(other body.Body) {
	player := e.MovementState().Target()
	if other.ID() == player.ID() {
		// TODO: Replace the condition of hurting
		player.(*actors.Player).Hurt(1)
	}
}
