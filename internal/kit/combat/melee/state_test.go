package melee_test

import (
	"image"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/tilemaplayer"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/mocks"
	"github.com/boilerplate/ebiten-template/internal/kit/combat/melee"
)

type mockOwner struct {
	mocks.MockActor
	falling bool
	goingUp bool
	ducking bool
	faceDir animation.FacingDirectionEnum
}

func (m *mockOwner) IsFalling() bool                              { return m.falling }
func (m *mockOwner) IsGoingUp() bool                              { return m.goingUp }
func (m *mockOwner) IsDucking() bool                              { return m.ducking }
func (m *mockOwner) FaceDirection() animation.FacingDirectionEnum { return m.faceDir }
func (m *mockOwner) GetPosition16() (int, int)                    { return m.MockActor.GetPosition16() }

type mockSpace struct{}

func (s *mockSpace) AddBody(_ body.Collidable)                                             {}
func (s *mockSpace) Bodies() []body.Collidable                                             { return nil }
func (s *mockSpace) RemoveBody(_ body.Collidable)                                          {}
func (s *mockSpace) QueueForRemoval(_ body.Collidable)                                     {}
func (s *mockSpace) ProcessRemovals()                                                      {}
func (s *mockSpace) Clear()                                                                {}
func (s *mockSpace) ResolveCollisions(_ body.Collidable) (bool, bool)                      { return false, false }
func (s *mockSpace) SetTilemapDimensionsProvider(_ tilemaplayer.TilemapDimensionsProvider) {}
func (s *mockSpace) GetTilemapDimensionsProvider() tilemaplayer.TilemapDimensionsProvider  { return nil }
func (s *mockSpace) Find(_ string) body.Collidable                                         { return nil }
func (s *mockSpace) Query(_ image.Rectangle) []body.Collidable                             { return nil }

func TestState_GetAnimationCount_ProgressesWithTime(t *testing.T) {
	w := threeStepWeapon()
	owner := &mockOwner{faceDir: animation.FaceDirectionRight}
	owner.SetPosition16(100, 200)
	space := &mockSpace{}

	// Create state
	st := melee.NewState(owner, space, w, nil, meleeAttackEnum, actors.Idle, actors.Falling)
	st.SetAnimationFrames(10)

	// Start state at count 100
	st.OnStart(100)

	// Immediately after OnStart, animation count should be 0
	if got := st.GetAnimationCount(100); got != 0 {
		t.Errorf("GetAnimationCount(100) = %d, want 0", got)
	}

	// After 1 frame (c.count = 101)
	if got := st.GetAnimationCount(101); got != 1 {
		t.Errorf("GetAnimationCount(101) = %d, want 1", got)
	}

	// After 5 frames (c.count = 105)
	if got := st.GetAnimationCount(105); got != 5 {
		t.Errorf("GetAnimationCount(105) = %d, want 5", got)
	}

	// Verify it still works if Update() was called (which increments s.frame but shouldn't affect GetAnimationCount)
	st.Update() // s.frame becomes 1
	if got := st.GetAnimationCount(101); got != 1 {
		t.Errorf("After Update, GetAnimationCount(101) = %d, want 1", got)
	}
}

func TestState_AnimationFinished(t *testing.T) {
	w := threeStepWeapon()
	owner := &mockOwner{}
	space := &mockSpace{}
	st := melee.NewState(owner, space, w, nil, meleeAttackEnum, actors.Idle, actors.Falling)
	st.SetAnimationFrames(3)
	st.OnStart(100)

	if st.IsAnimationFinished() {
		t.Error("IsAnimationFinished() = true, want false at start")
	}

	st.Update() // frame 1
	if st.IsAnimationFinished() {
		t.Error("IsAnimationFinished() = true, want false at frame 1")
	}

	st.Update() // frame 2
	if st.IsAnimationFinished() {
		t.Error("IsAnimationFinished() = true, want false at frame 2")
	}

	st.Update() // frame 3
	if !st.IsAnimationFinished() {
		t.Error("IsAnimationFinished() = false, want true at frame 3")
	}
}
