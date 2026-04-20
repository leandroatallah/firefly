package gameplayer_test

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	gameplayer "github.com/boilerplate/ebiten-template/internal/game/entity/actors/player"
	gamestates "github.com/boilerplate/ebiten-template/internal/game/entity/actors/states"
)

// stubMovement implements the movement-query interface used by shootingContributor.
type stubMovement struct {
	walking, goingUp, falling bool
}

func (s *stubMovement) IsWalking() bool { return s.walking }
func (s *stubMovement) IsGoingUp() bool { return s.goingUp }
func (s *stubMovement) IsFalling() bool { return s.falling }

func TestShootingContributor_StateMapping(t *testing.T) {
	tests := []struct {
		name          string
		shootHeld     bool
		walking       bool
		goingUp       bool
		falling       bool
		expectedState actors.ActorStateEnum
		expectedOK    bool
	}{
		{"inactive — defer", false, false, false, false, 0, false},
		{"active idle", true, false, false, false, actors.IdleShooting, true},
		{"active walking", true, true, false, false, actors.WalkingShooting, true},
		{"active going up", true, false, true, false, actors.JumpingShooting, true},
		{"active falling", true, false, false, true, actors.FallingShooting, true},
		{"active going up beats walking", true, true, true, false, actors.JumpingShooting, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contrib := gameplayer.NewShootingContributorForTest(tt.shootHeld, &stubMovement{
				walking: tt.walking,
				goingUp: tt.goingUp,
				falling: tt.falling,
			})
			got, ok := contrib.ContributeState(actors.Idle)
			if ok != tt.expectedOK {
				t.Errorf("ok: want %v got %v", tt.expectedOK, ok)
			}
			if ok && got != tt.expectedState {
				t.Errorf("state: want %v got %v", tt.expectedState, got)
			}
		})
	}
}

func TestDashContributor_StateMapping(t *testing.T) {
	tests := []struct {
		name          string
		active        bool
		expectedState actors.ActorStateEnum
		expectedOK    bool
	}{
		{"inactive — defer", false, 0, false},
		{"active — StateDashing", true, gamestates.StateDashing, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contrib := gameplayer.NewDashContributorForTest(tt.active)
			got, ok := contrib.ContributeState(actors.Idle)
			if ok != tt.expectedOK {
				t.Errorf("ok: want %v got %v", tt.expectedOK, ok)
			}
			if ok && got != tt.expectedState {
				t.Errorf("state: want %v got %v", tt.expectedState, got)
			}
		})
	}
}
