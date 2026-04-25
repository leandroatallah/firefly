package gameplayer_test

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	gameplayer "github.com/boilerplate/ebiten-template/internal/game/entity/actors/player"
)

// stubMeleeActive implements the meleeActiveIface surface needed by the
// meleeContributor (IsSwinging + StepIndex).
type stubMeleeActive struct {
	swinging  bool
	inStartup bool
	step      int
}

func (s *stubMeleeActive) IsSwinging() bool  { return s.swinging }
func (s *stubMeleeActive) IsInStartup() bool { return s.inStartup }
func (s *stubMeleeActive) StepIndex() int    { return s.step }

// TestMeleeContributor_StateMapping verifies the contributor's behavior across
// the three observable conditions:
//  1. Not swinging → defer (0, false).
//  2. Swinging with a valid step index → that step's enum + true.
//  3. Swinging with an out-of-range step index → safely defer (0, false).
func TestMeleeContributor_StateMapping(t *testing.T) {
	// Three distinct step state enums — produced by the helper under test.
	stepStates := []actors.ActorStateEnum{101, 102, 103}

	tests := []struct {
		name      string
		swinging  bool
		step      int
		wantState actors.ActorStateEnum
		wantOK    bool
	}{
		{name: "not swinging defers", swinging: false, step: 0, wantState: 0, wantOK: false},
		{name: "not swinging with nonzero step still defers", swinging: false, step: 2, wantState: 0, wantOK: false},
		{name: "swinging step 0", swinging: true, step: 0, wantState: stepStates[0], wantOK: true},
		{name: "swinging step 1", swinging: true, step: 1, wantState: stepStates[1], wantOK: true},
		{name: "swinging step 2 (last)", swinging: true, step: 2, wantState: stepStates[2], wantOK: true},
		{name: "swinging step out of range high defers", swinging: true, step: 3, wantState: 0, wantOK: false},
		{name: "swinging step negative defers", swinging: true, step: -1, wantState: 0, wantOK: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stub := &stubMeleeActive{swinging: tc.swinging, step: tc.step}
			contrib := gameplayer.NewMeleeContributorForTest(stub, stepStates)

			got, ok := contrib.ContributeState(actors.Idle)
			if ok != tc.wantOK {
				t.Errorf("ok = %v, want %v", ok, tc.wantOK)
			}
			if ok && got != tc.wantState {
				t.Errorf("state = %v, want %v", got, tc.wantState)
			}
		})
	}
}
