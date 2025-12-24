package gameenemies

import (
	"log"

	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/actors/enemies"
	"github.com/leandroatallah/firefly/internal/engine/actors/movement"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
)

type BlueEnemy struct {
	actors.Character
	count int
}

func NewBlueEnemy(x, y int, id string) (*BlueEnemy, error) {
	spriteData, statData, err := enemies.ParseJsonEnemy("internal/game/actors/enemies/blue_enemy.json")
	if err != nil {
		log.Fatal(err)
	}

	character, err := CreateAnimatedCharacter(spriteData)
	if err != nil {
		log.Fatal(err)
	}

	character.SetPosition(x, y)
	enemy := &BlueEnemy{Character: *character}
	enemy.SetID(id)

	if err = SetEnemyStats(enemy, statData); err != nil {
		return nil, err
	}
	if err = SetEnemyBodies(enemy, spriteData); err != nil {
		return nil, err
	}
	enemy.SetTouchable(enemy)

	return enemy, nil
}

func (e *BlueEnemy) SetTarget(target body.MovableCollidable) {
	e.Character.SetMovementState(movement.DumbChase, target)
}

// Character Methods
func (e *BlueEnemy) Update(space body.BodiesSpace) error {
	e.count++
	return e.Character.Update(space)
}

func (e *BlueEnemy) OnTouch(other body.Collidable) {
	player := e.MovementState().Target()
	if other.ID() == player.ID() {
		player.(*actors.Character).Hurt(1)
	}
}
