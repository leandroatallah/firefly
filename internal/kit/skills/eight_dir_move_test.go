package kitskills

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/input"
	"github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
	"github.com/boilerplate/ebiten-template/internal/engine/skill"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/fp16"
	"github.com/hajimehoshi/ebiten/v2"
)

// mockEightDirBody wraps mockMovableCollidable to record OnMove* calls
// for all four cardinal directions.
type mockEightDirBody struct {
	*mockMovableCollidable
	onMoveLeftCalls  int
	onMoveRightCalls int
	onMoveUpCalls    int
	onMoveDownCalls  int
	lastLeftArg      int
	lastRightArg     int
	lastUpArg        int
	lastDownArg      int
}

func newMockEightDirBody() *mockEightDirBody {
	return &mockEightDirBody{mockMovableCollidable: newMockMovableCollidable()}
}

func (m *mockEightDirBody) OnMoveLeft(speed int) {
	m.onMoveLeftCalls++
	m.lastLeftArg = speed
	m.mockMovableCollidable.OnMoveLeft(speed)
}

func (m *mockEightDirBody) OnMoveRight(speed int) {
	m.onMoveRightCalls++
	m.lastRightArg = speed
	m.mockMovableCollidable.OnMoveRight(speed)
}

func (m *mockEightDirBody) OnMoveUp(speed int) {
	m.onMoveUpCalls++
	m.lastUpArg = speed
	m.mockMovableCollidable.OnMoveUp(speed)
}

func (m *mockEightDirBody) OnMoveDown(speed int) {
	m.onMoveDownCalls++
	m.lastDownArg = speed
	m.mockMovableCollidable.OnMoveDown(speed)
}

func TestEightDirectionalMovementSkill_HandleInput(t *testing.T) {
	const speed = 200

	tests := []struct {
		name     string
		cmds     input.PlayerCommands
		immobile bool
		blocked  bool

		wantLeft, wantRight, wantUp, wantDown int

		preVX, preVY     int
		preAccX, preAccY int

		// expected post-state for guard cases
		checkUnchangedVel  bool
		expectVX, expectVY int
		checkZeroedVelAcc  bool
	}{
		{
			name:     "move_left_only",
			cmds:     input.PlayerCommands{Left: true},
			wantLeft: 1,
		},
		{
			name:      "move_right_only",
			cmds:      input.PlayerCommands{Right: true},
			wantRight: 1,
		},
		{
			name:   "move_up_only",
			cmds:   input.PlayerCommands{Up: true},
			wantUp: 1,
		},
		{
			name:     "move_down_only",
			cmds:     input.PlayerCommands{Down: true},
			wantDown: 1,
		},
		{
			name:     "diagonal_left_up",
			cmds:     input.PlayerCommands{Left: true, Up: true},
			wantLeft: 1,
			wantUp:   1,
		},
		{
			name:              "immobile_guard",
			cmds:              input.PlayerCommands{Left: true},
			immobile:          true,
			preVX:             fp16.To16(5),
			preVY:             fp16.To16(3),
			preAccX:           fp16.To16(2),
			preAccY:           fp16.To16(1),
			checkZeroedVelAcc: true,
		},
		{
			name:              "input_blocked_guard",
			cmds:              input.PlayerCommands{Left: true},
			blocked:           true,
			preVX:             fp16.To16(7),
			preVY:             fp16.To16(4),
			checkUnchangedVel: true,
			expectVX:          fp16.To16(7),
			expectVY:          fp16.To16(4),
		},
		{
			name: "no_input",
			cmds: input.PlayerCommands{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actor := newMockEightDirBody()
			actor.SetSpeed(speed)
			if tc.immobile {
				actor.SetImmobile(true)
			}
			if tc.preVX != 0 || tc.preVY != 0 {
				actor.SetVelocity(tc.preVX, tc.preVY)
			}
			if tc.preAccX != 0 || tc.preAccY != 0 {
				actor.SetAcceleration(tc.preAccX, tc.preAccY)
			}

			model := movement.NewPlatformMovementModel(&mockPlayerMovementBlocker{blocked: tc.blocked})

			orig := input.CommandsReader
			defer func() { input.CommandsReader = orig }()
			input.CommandsReader = func() input.PlayerCommands { return tc.cmds }

			s := NewEightDirectionalMovementSkill()
			s.HandleInput(actor, model, nil)

			if actor.onMoveLeftCalls != tc.wantLeft {
				t.Errorf("OnMoveLeft calls: got %d, want %d", actor.onMoveLeftCalls, tc.wantLeft)
			}
			if actor.onMoveRightCalls != tc.wantRight {
				t.Errorf("OnMoveRight calls: got %d, want %d", actor.onMoveRightCalls, tc.wantRight)
			}
			if actor.onMoveUpCalls != tc.wantUp {
				t.Errorf("OnMoveUp calls: got %d, want %d", actor.onMoveUpCalls, tc.wantUp)
			}
			if actor.onMoveDownCalls != tc.wantDown {
				t.Errorf("OnMoveDown calls: got %d, want %d", actor.onMoveDownCalls, tc.wantDown)
			}

			if tc.wantLeft > 0 && actor.lastLeftArg != speed {
				t.Errorf("OnMoveLeft arg: got %d, want %d", actor.lastLeftArg, speed)
			}
			if tc.wantRight > 0 && actor.lastRightArg != speed {
				t.Errorf("OnMoveRight arg: got %d, want %d", actor.lastRightArg, speed)
			}
			if tc.wantUp > 0 && actor.lastUpArg != speed {
				t.Errorf("OnMoveUp arg: got %d, want %d", actor.lastUpArg, speed)
			}
			if tc.wantDown > 0 && actor.lastDownArg != speed {
				t.Errorf("OnMoveDown arg: got %d, want %d", actor.lastDownArg, speed)
			}

			if tc.checkZeroedVelAcc {
				vx, vy := actor.Velocity()
				if vx != 0 || vy != 0 {
					t.Errorf("expected zeroed velocity when immobile; got (%d, %d)", vx, vy)
				}
				accX, accY := actor.Acceleration()
				if accX != 0 || accY != 0 {
					t.Errorf("expected zeroed acceleration when immobile; got (%d, %d)", accX, accY)
				}
			}

			if tc.checkUnchangedVel {
				vx, vy := actor.Velocity()
				if vx != tc.expectVX || vy != tc.expectVY {
					t.Errorf("expected velocity unchanged (%d, %d); got (%d, %d)", tc.expectVX, tc.expectVY, vx, vy)
				}
			}
		})
	}
}

func TestEightDirectionalMovementSkill_New(t *testing.T) {
	s := NewEightDirectionalMovementSkill()
	if s == nil {
		t.Fatal("NewEightDirectionalMovementSkill returned nil")
	}
	if s.State() != skill.StateReady {
		t.Errorf("expected state Ready; got %s", s.State())
	}
}

func TestEightDirectionalMovementSkill_Update_NoOp(t *testing.T) {
	actor := newMockEightDirBody()
	model := movement.NewPlatformMovementModel(nil)

	s := NewEightDirectionalMovementSkill()
	preState := s.State()

	s.Update(actor, model)

	if s.State() != preState {
		t.Errorf("expected state unchanged (%s); got %s", preState, s.State())
	}
	if actor.onMoveLeftCalls != 0 || actor.onMoveRightCalls != 0 || actor.onMoveUpCalls != 0 || actor.onMoveDownCalls != 0 {
		t.Errorf("expected no OnMove* calls during Update; got L=%d R=%d U=%d D=%d",
			actor.onMoveLeftCalls, actor.onMoveRightCalls, actor.onMoveUpCalls, actor.onMoveDownCalls)
	}
}

func TestEightDirectionalMovementSkill_ActivationKey(t *testing.T) {
	s := NewEightDirectionalMovementSkill()
	if s.ActivationKey() != ebiten.Key(0) {
		t.Errorf("expected zero ebiten.Key; got %v", s.ActivationKey())
	}
}
