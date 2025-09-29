package movement

import "math/rand/v2"

type RandMovementState struct {
	BaseMovementState
}

func (s *RandMovementState) Move() {
	a := []func(){
		func() { s.actor.OnMoveLeft(s.actor.Speed()) },
		func() { s.actor.OnMoveRight(s.actor.Speed()) },
		func() { s.actor.OnMoveUp(s.actor.Speed()) },
		func() { s.actor.OnMoveDown(s.actor.Speed()) },
	}

	i := rand.IntN(4)
	a[i]()
}
