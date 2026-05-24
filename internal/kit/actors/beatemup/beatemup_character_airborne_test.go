package beatemup_test

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/kit/actors/beatemup"
)

// Ensure DownwardGravity is set: the handler uses it as the ground-plane
// motion threshold. A value of 4 matches the platformer test default and the
// existing beatemup_jump test default.
func init() {
	config.Set(&config.AppConfig{
		Physics: config.PhysicsConfig{
			DownwardGravity: 4,
			UpwardGravity:   2,
		},
	})
}

// fakeAirborneState is a minimal package-local ActorState used to control
// IsAnimationFinished() for Landing-lock tests. The handler reads the current
// state instance's IsAnimationFinished(), so wiring a fake here is the
// lowest-friction way to pin AC-4 behaviour deterministically.
type fakeAirborneState struct {
	enum     actors.ActorStateEnum
	finished bool
}

func (f *fakeAirborneState) State() actors.ActorStateEnum { return f.enum }
func (f *fakeAirborneState) OnStart(_ int)                {}
func (f *fakeAirborneState) OnFinish()                    {}
func (f *fakeAirborneState) GetAnimationCount(_ int) int  { return 0 }
func (f *fakeAirborneState) IsAnimationFinished() bool    { return f.finished }

// newAirborneCharacter builds a freshly-constructed BeatEmUpCharacter wired
// with the production movement transition handler. It registers fake
// Jumping/Falling/Landing states so the seeded initial state has a
// controllable IsAnimationFinished() value.
func newAirborneCharacter(t *testing.T, seed actors.ActorStateEnum, landingFinished bool) *beatemup.BeatEmUpCharacter {
	t.Helper()
	fsys, stateMap, spriteData, bodyRect := newTestFixtures()
	c, err := beatemup.NewBeatEmUpCharacter(fsys, stateMap, spriteData, bodyRect, nil)
	if err != nil {
		t.Fatalf("NewBeatEmUpCharacter: %v", err)
	}
	// Register controllable fakes for every airborne state we may transition
	// to or seed from. The Landing fake's IsAnimationFinished() is the
	// switch the handler reads to release the landing lock.
	c.SetStateInstance(actors.Idle, &fakeAirborneState{enum: actors.Idle})
	c.SetStateInstance(actors.Walking, &fakeAirborneState{enum: actors.Walking})
	c.SetStateInstance(actors.Jumping, &fakeAirborneState{enum: actors.Jumping})
	c.SetStateInstance(actors.Falling, &fakeAirborneState{enum: actors.Falling})
	c.SetStateInstance(actors.Landing, &fakeAirborneState{enum: actors.Landing, finished: landingFinished})
	c.SetNewStateFatal(seed)
	return c
}

// TestBeatemupMovementTransitions_AirborneStates is the table-driven Red
// suite for story 066. It exercises beatemupMovementTransitions through the
// public MovementTransitionHandler field installed by the constructor.
//
// Each row mirrors a row in SPEC §6 (T-A1..T-A15) and asserts the resulting
// observable ActorState. The handler currently knows only about
// Idle/Walking, so every airborne row is expected to FAIL until the handler
// is extended.
func TestBeatemupMovementTransitions_AirborneStates(t *testing.T) {
	type setup struct {
		seedState       actors.ActorStateEnum
		vAlt16          int
		altitude        int
		vx16            int
		vy16            int
		landingFinished bool
	}
	tests := []struct {
		name  string
		ac    string
		setup setup
		want  actors.ActorStateEnum
	}{
		{
			name:  "T-A1 [AC-1] ascend from Idle -> Jumping",
			ac:    "AC-1",
			setup: setup{seedState: actors.Idle, vAlt16: -1000, altitude: 0},
			want:  actors.Jumping,
		},
		{
			name:  "T-A2 [AC-1,AC-5] ascend with ground velocity stays Jumping (no Walking)",
			ac:    "AC-1/AC-5",
			setup: setup{seedState: actors.Idle, vAlt16: -1000, altitude: 10, vx16: 2000},
			want:  actors.Jumping,
		},
		{
			name:  "T-A3 [AC-6] apex during Jumping -> Falling",
			ac:    "AC-6",
			setup: setup{seedState: actors.Jumping, vAlt16: 0, altitude: 20},
			want:  actors.Falling,
		},
		{
			name:  "T-A4 [AC-6] descent during Jumping -> Falling",
			ac:    "AC-6",
			setup: setup{seedState: actors.Jumping, vAlt16: 500, altitude: 20},
			want:  actors.Falling,
		},
		{
			name:  "T-A5 [AC-2] descent airborne from Idle -> Falling",
			ac:    "AC-2",
			setup: setup{seedState: actors.Idle, vAlt16: 500, altitude: 20},
			want:  actors.Falling,
		},
		{
			name:  "T-A6 [AC-3] touchdown from Falling -> Landing",
			ac:    "AC-3",
			setup: setup{seedState: actors.Falling, vAlt16: 0, altitude: 0},
			want:  actors.Landing,
		},
		{
			name:  "T-A7 [AC-4] Landing locked while animation playing",
			ac:    "AC-4",
			setup: setup{seedState: actors.Landing, landingFinished: false},
			want:  actors.Landing,
		},
		{
			name:  "T-A8 [AC-4] Landing -> Idle when finished and still",
			ac:    "AC-4",
			setup: setup{seedState: actors.Landing, landingFinished: true},
			want:  actors.Idle,
		},
		{
			name:  "T-A9 [AC-4] Landing -> Walking when finished and ground-moving",
			ac:    "AC-4",
			setup: setup{seedState: actors.Landing, landingFinished: true, vx16: 2000},
			want:  actors.Walking,
		},
		{
			name:  "T-A10 [AC-5] Jumping with ground velocity stays Jumping",
			ac:    "AC-5",
			setup: setup{seedState: actors.Jumping, vAlt16: -500, altitude: 20, vx16: 2000},
			want:  actors.Jumping,
		},
		{
			name:  "T-A11 [AC-5] Falling with ground velocity stays Falling",
			ac:    "AC-5",
			setup: setup{seedState: actors.Falling, vAlt16: 500, altitude: 20, vx16: 2000},
			want:  actors.Falling,
		},
		{
			name:  "T-A12 (edge) jump impulse with altitude==0 -> Jumping (no spurious Landing)",
			ac:    "edge",
			setup: setup{seedState: actors.Idle, vAlt16: -1000, altitude: 0},
			want:  actors.Jumping,
		},
		{
			name:  "T-A13 (edge) buffered jump on landing frame: Falling with alt==0 -> Landing first",
			ac:    "edge",
			setup: setup{seedState: actors.Falling, vAlt16: 0, altitude: 0},
			want:  actors.Landing,
		},
		{
			name:  "T-A14 (regression) idle on ground, no motion -> Idle",
			ac:    "regression",
			setup: setup{seedState: actors.Idle, vAlt16: 0, altitude: 0},
			want:  actors.Idle,
		},
		{
			name:  "T-A15 (regression) ground-moving on ground -> Walking",
			ac:    "regression",
			setup: setup{seedState: actors.Idle, vAlt16: 0, altitude: 0, vx16: 2000},
			want:  actors.Walking,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newAirborneCharacter(t, tt.setup.seedState, tt.setup.landingFinished)
			c.SetAltitude(tt.setup.altitude)
			c.SetVAltitude16(tt.setup.vAlt16)
			c.SetVelocity(tt.setup.vx16, tt.setup.vy16)

			if c.MovementTransitionHandler == nil {
				t.Fatal("MovementTransitionHandler is nil; constructor must wire beatemupMovementTransitions")
			}
			c.MovementTransitionHandler(c.Character)

			if got := c.State(); got != tt.want {
				t.Errorf("[%s] state after handler = %v, want %v (seed=%v vAlt16=%d alt=%d vx=%d vy=%d landingFinished=%v)",
					tt.ac, got, tt.want,
					tt.setup.seedState, tt.setup.vAlt16, tt.setup.altitude,
					tt.setup.vx16, tt.setup.vy16, tt.setup.landingFinished)
			}
		})
	}
}
