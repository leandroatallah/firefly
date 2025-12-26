package gameenemies

import (
	"log"

	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/actors/enemies"
	"github.com/leandroatallah/firefly/internal/engine/actors/movement"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
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

	if err = SetEnemyStats(enemy, statData); err != nil {
		return nil, err
	}
	if err = SetEnemyBodies(enemy, spriteData, id); err != nil {
		return nil, err
	}

	model, err := physics.NewMovementModel(physics.Platform, nil)
	if err != nil {
		return nil, err
	}
	enemy.SetMovementModel(model)
	enemy.SetTouchable(enemy)

	return enemy, nil
}

func (e *BlueEnemy) SetTarget(target body.MovableCollidable) {
	e.Character.SetMovementState(movement.Idle, target)
}

// Character Methods
func (e *BlueEnemy) Update(space body.BodiesSpace) error {
	e.count++
	return e.Character.Update(space)
}

func (e *BlueEnemy) GetCharacter() *actors.Character {
	return &e.Character
}

func (e *BlueEnemy) OnTouch(other body.Collidable) {
	player := e.MovementState().Target()
	if other.ID() == player.ID() {
		player.(actors.ActorEntity).GetCharacter().Hurt(1)
	}
}
