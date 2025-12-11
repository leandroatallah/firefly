package gameenemies

import (
	"log"

	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
	"github.com/leandroatallah/firefly/internal/engine/systems/sprites"
)

type BlueEnemy struct {
	actors.Character
	count int
}

func NewBlueEnemy() (*BlueEnemy, error) {
	// NOTE: Ignore for now. Set logic to initial position
	x, y := 0, 0
	const (
		frameWidth  = 32
		frameHeight = 32
	)

	var assets sprites.SpriteAssets
	assets = assets.AddSprite(actors.Idle, "assets/images/blue-enemy.png").
		AddSprite(actors.Walking, "assets/images/blue-enemy.png")

	sprites, err := sprites.LoadSprites(assets)
	if err != nil {
		log.Fatal(err)
	}

	rect := physics.NewRect(x, y, frameWidth, frameHeight)
	character := actors.NewCharacter(sprites, rect)

	collisionRect := physics.NewRect(x, y, frameWidth, frameHeight)

	enemy := &BlueEnemy{Character: *character}
	// TODO: Handle repeated IDs
	enemy.SetID("BLUEENEMY")
	err = enemy.SetSpeed(2)
	if err != nil {
		return nil, err
	}
	err = enemy.SetMaxSpeed(2)
	if err != nil {
		return nil, err
	}
	enemy.AddCollision(physics.NewCollidableBodyFromRect(collisionRect))
	enemy.SetTouchable(enemy)

	return enemy, nil
}

// Character Methods
func (e *BlueEnemy) Update(space body.BodiesSpace) error {
	e.count++
	return e.Character.Update(space)
}

func (e *BlueEnemy) OnTouch(other body.Collidable) {
	player := e.MovementState().Target()
	if other.ID() == player.ID() {
		// TODO: Replace the condition of hurting
		player.(*actors.Character).Hurt(1)
	}
}
