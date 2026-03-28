package phases

import (
	"testing"

	"github.com/leandroatallah/firefly/internal/engine/contracts/navigation"
	contractseq "github.com/leandroatallah/firefly/internal/engine/contracts/sequences"
)

func TestManagerAdvanceToNextPhaseIntegration(t *testing.T) {
	m := NewManager()
	sceneType1 := navigation.SceneType(1)
	sceneType2 := navigation.SceneType(2)

	m.AddPhase(Phase{ID: 1, SceneType: sceneType1, NextPhaseID: 2})
	m.AddPhase(Phase{ID: 2, SceneType: sceneType2})

	if err := m.SetCurrentPhase(1); err != nil {
		t.Fatalf("SetCurrentPhase: %v", err)
	}

	if err := m.AdvanceToNextPhase(); err != nil {
		t.Fatalf("AdvanceToNextPhase: %v", err)
	}

	if m.CurrentPhase != 2 {
		t.Fatalf("expected CurrentPhase to be 2 after AdvanceToNextPhase, got %d", m.CurrentPhase)
	}
}

type testPlayer struct {
	playing bool
	played  []contractseq.Sequence
	updates int
	isOver  bool
}

func (p *testPlayer) IsPlaying() bool {
	return p.playing
}

func (p *testPlayer) IsOver() bool {
	return p.isOver
}

func (p *testPlayer) Play(s contractseq.Sequence) {
	p.playing = true
	p.isOver = false
	p.played = append(p.played, s)
}

func (p *testPlayer) Stop() {
	p.playing = false
}

func (p *testPlayer) PlaySequence(filePath string) {
	p.playing = true
	p.isOver = false
}

func (p *testPlayer) Update() {
	p.updates++
}

func TestSequenceGoalCompletesWhenPlayerStops(t *testing.T) {
	p := &testPlayer{
		playing: true,
	}

	called := 0
	goal := &SequenceGoal{
		Player: p,
		OnCompleteFunc: func() {
			called++
		},
	}

	if goal.IsCompleted() {
		t.Fatalf("expected goal not completed while player is playing")
	}

	p.playing = false

	if !goal.IsCompleted() {
		t.Fatalf("expected goal completed when player is not playing")
	}

	goal.OnCompletion()
	if called != 1 {
		t.Fatalf("expected OnCompleteFunc to be called once, got %d", called)
	}
}

func TestNoGoalNeverCompletes(t *testing.T) {
	g := &NoGoal{}
	if g.IsCompleted() {
		t.Fatalf("NoGoal should never be completed")
	}
	g.OnCompletion()
}
