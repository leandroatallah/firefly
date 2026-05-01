package kitstates_test

import (
	"testing"

	kitstates "github.com/boilerplate/ebiten-template/internal/kit/states"
)

// subStateEnum is the local enum supplied as the type parameter E to
// IdleSubState. It mirrors the shape of gamestates.GroundedSubStateEnum
// without coupling kit tests to the game package.
type subStateEnum int

const (
	idle subStateEnum = iota
	walking
	ducking
	aimLock
)

// fakeInput satisfies kitstates.GroundedInputLike via func fields so each
// table row can declare its own input surface inline without a mock
// framework.
type fakeInput struct {
	aimLockHeld     bool
	duckHeld        bool
	horizontalInput int
}

func (f *fakeInput) AimLockHeld() bool    { return f.aimLockHeld }
func (f *fakeInput) DuckHeld() bool       { return f.duckHeld }
func (f *fakeInput) HorizontalInput() int { return f.horizontalInput }

func newIdleSubState() *kitstates.IdleSubState[subStateEnum, *fakeInput] {
	return &kitstates.IdleSubState[subStateEnum, *fakeInput]{
		Idle:    idle,
		Walking: walking,
		Ducking: ducking,
		AimLock: aimLock,
	}
}

func TestIdleSubStateTransitionTo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *fakeInput
		want  subStateEnum
	}{
		{
			name:  "all zero stays idle",
			input: &fakeInput{},
			want:  idle,
		},
		{
			name:  "horizontal input right transitions to walking",
			input: &fakeInput{horizontalInput: 1},
			want:  walking,
		},
		{
			name:  "horizontal input left transitions to walking",
			input: &fakeInput{horizontalInput: -1},
			want:  walking,
		},
		{
			name:  "duck held transitions to ducking",
			input: &fakeInput{duckHeld: true},
			want:  ducking,
		},
		{
			name:  "aim lock held transitions to aim lock",
			input: &fakeInput{aimLockHeld: true},
			want:  aimLock,
		},
		{
			name:  "aim lock takes precedence over duck",
			input: &fakeInput{aimLockHeld: true, duckHeld: true},
			want:  aimLock,
		},
		{
			name:  "duck takes precedence over horizontal input",
			input: &fakeInput{duckHeld: true, horizontalInput: 1},
			want:  ducking,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := newIdleSubState()
			got := s.TransitionTo(tc.input)
			if got != tc.want {
				t.Fatalf("TransitionTo(%+v) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestIdleSubStateOnStartIsNoOp(t *testing.T) {
	t.Parallel()

	s := newIdleSubState()
	before := *s

	// Should not panic and must not mutate observable struct fields.
	s.OnStart(7)

	after := *s
	if before != after {
		t.Fatalf("OnStart mutated state: before=%+v after=%+v", before, after)
	}

	// After a no-op OnStart, transition behaviour is unchanged.
	if got := s.TransitionTo(&fakeInput{}); got != idle {
		t.Fatalf("TransitionTo after OnStart = %v, want %v", got, idle)
	}
}

func TestIdleSubStateOnFinishIsNoOp(t *testing.T) {
	t.Parallel()

	s := newIdleSubState()
	before := *s

	// Should not panic.
	s.OnFinish()

	after := *s
	if before != after {
		t.Fatalf("OnFinish mutated state: before=%+v after=%+v", before, after)
	}
}
