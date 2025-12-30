package body

import (
	"fmt"
	"testing"

	"github.com/leandroatallah/firefly/internal/engine/utils/fp16"
)

type accVect struct {
	accelerationX int
	accelerationY int
}

func NewAccVect(x16, y16 int) accVect {
	return accVect{x16, y16}
}

func (v accVect) String() string {
	return fmt.Sprintf("{ x16: %d, y16: %d }", v.accelerationX, v.accelerationY)
}

func TestMovableBody_Movement(t *testing.T) {
	b := NewMovableBody(NewBody(NewRect(0, 0, 10, 10)))

	distance := 5
	distancex16 := fp16.To16(distance)

	tests := []struct {
		name string
		fn   func(int)
		want accVect
	}{
		{"Move Left", b.OnMoveLeft, NewAccVect(-distancex16, 0)},
		{"Move Right", b.OnMoveRight, NewAccVect(distancex16, 0)},
		{"Move Up", b.OnMoveUp, NewAccVect(0, -distancex16)},
		{"Move Down", b.OnMoveDown, NewAccVect(0, distancex16)},
		{"Move Up Left", b.OnMoveUpLeft, NewAccVect(-distancex16, -distancex16)},
		{"Move Up Right", b.OnMoveUpRight, NewAccVect(distancex16, -distancex16)},
		{"Move Down Left", b.OnMoveDownLeft, NewAccVect(-distancex16, distancex16)},
		{"Move Down Right", b.OnMoveDownRight, NewAccVect(distancex16, distancex16)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset
			b.accelerationX, b.accelerationY = 0, 0

			tt.fn(distance)
			if b.accelerationX != tt.want.accelerationX || b.accelerationY != tt.want.accelerationY {
				t.Errorf(
					"expected %v; got %v",
					tt.want, fmt.Sprintf(" { x16: %d , y16: %d }", b.accelerationX, b.accelerationY))
			}
		})
	}
}
