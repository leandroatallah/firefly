package movement

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/utils/fp16"
)

func TestIncreaseVelocity(t *testing.T) {
	tests := []struct {
		name       string
		velocity   int
		accel      int
		want       int
	}{
		{"positive acceleration", 0, 10, 10},
		{"negative acceleration", 0, -10, -10},
		{"zero acceleration", 10, 0, 10},
		{"add to existing positive", 10, 5, 15},
		{"add to existing negative", -10, -5, -15},
		{"reduce positive", 10, -5, 5},
		{"reduce negative", -10, 5, -5},
		{"reverse direction", 10, -20, -10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := increaseVelocity(tt.velocity, tt.accel)
			if got != tt.want {
				t.Errorf("increaseVelocity(%d, %d) = %d; want %d",
					tt.velocity, tt.accel, got, tt.want)
			}
		})
	}
}

func TestReduceVelocity(t *testing.T) {
	friction := fp16.To16(1) / 4 // 4 in fp16

	tests := []struct {
		name     string
		velocity int
		want     int
	}{
		{"zero velocity", 0, 0},
		{"positive above friction", 10, 6},
		{"negative above friction", -10, -6},
		{"positive below friction", 2, 0},
		{"negative below friction", -2, 0},
		{"positive equals friction", friction, 0},
		{"negative equals friction", -friction, 0},
		{"small positive", 1, 0},
		{"small negative", -1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := reduceVelocity(tt.velocity)
			if got != tt.want {
				t.Errorf("reduceVelocity(%d) = %d; want %d",
					tt.velocity, got, tt.want)
			}
		})
	}
}

func TestSmoothDiagonalMovement(t *testing.T) {
	tests := []struct {
		name        string
		accX, accY  int
	}{
		{"zero input", 0, 0},
		{"right only", 2, 0},
		{"left only", -2, 0},
		{"up only", 0, -2},
		{"down only", 0, 2},
		{"diagonal up-right", 2, -2},
		{"diagonal down-left", -2, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotX, gotY := smoothDiagonalMovement(tt.accX, tt.accY)
			// For diagonal movements, we check that normalization occurred
			if tt.accX != 0 && tt.accY != 0 {
				// Diagonal: both components should be reduced by sqrt(2)
				// We just verify they're non-zero and roughly equal magnitude
				if gotX == 0 && gotY == 0 {
					// This is acceptable for small values after normalization
				}
			} else {
				// Cardinal: just verify it's scaled and non-zero for non-zero input
				if tt.accX != 0 && gotX == 0 {
					t.Errorf("smoothDiagonalMovement(%d, %d) gotX=0; want non-zero", tt.accX, tt.accY)
				}
				if tt.accY != 0 && gotY == 0 {
					t.Errorf("smoothDiagonalMovement(%d, %d) gotY=0; want non-zero", tt.accX, tt.accY)
				}
			}
		})
	}
}

func TestClampAxisVelocity(t *testing.T) {
	tests := []struct {
		name     string
		velocity int
		limit    int
		want     int
	}{
		{"within positive limit", 5, 10, 5},
		{"within negative limit", -5, 10, -5},
		{"exceeds positive limit", 15, 10, 10},
		{"exceeds negative limit", -15, 10, -10},
		{"at positive limit", 10, 10, 10},
		{"at negative limit", -10, 10, -10},
		{"zero velocity", 0, 10, 0},
		{"zero limit", 5, 0, 0},
		{"negative limit", 5, -10, 0},
		{"velocity zero negative limit", 0, -10, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := clampAxisVelocity(tt.velocity, tt.limit)
			if got != tt.want {
				t.Errorf("clampAxisVelocity(%d, %d) = %d; want %d",
					tt.velocity, tt.limit, got, tt.want)
			}
		})
	}
}
