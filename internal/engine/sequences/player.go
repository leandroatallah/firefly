package sequences

import (
	"github.com/leandroatallah/firefly/internal/engine/core"
)

// SequencePlayer manages the execution of a sequence.
type SequencePlayer struct {
	appContext          *core.AppContext
	currentSequence     Sequence
	currentCommandIndex int
	isPlaying           bool
}

// NewSequencePlayer creates a new player.
func NewSequencePlayer(appContext *core.AppContext) *SequencePlayer {
	return &SequencePlayer{
		appContext: appContext,
	}
}

// Play starts executing a sequence.
func (p *SequencePlayer) Play(sequence Sequence) {
	if p.isPlaying {
		return // Do not play if another sequence is already in progress
	}
	p.currentSequence = sequence
	p.currentCommandIndex = -1 // Will be incremented to 0 by advanceToNextCommand
	p.isPlaying = true

	if sequence.BlockPlayerMovement {
		if player, found := p.appContext.ActorManager.GetPlayer(); found {
			player.BlockMovement()
		}
	}
	p.advanceToNextCommand()
}

// IsPlaying returns true if a sequence is currently being played.
func (p *SequencePlayer) IsPlaying() bool {
	return p.isPlaying
}

// Update should be called every frame. It updates the current command.
func (p *SequencePlayer) Update() {
	if !p.isPlaying {
		return
	}

	if p.currentCommandIndex >= len(p.currentSequence.Commands) {
		p.endSequence()
		return
	}

	currentCommand := p.currentSequence.Commands[p.currentCommandIndex]
	if currentCommand.Update() {
		p.advanceToNextCommand()
	}
}

// advanceToNextCommand moves to the next command in the queue and initializes it.
func (p *SequencePlayer) advanceToNextCommand() {
	p.currentCommandIndex++
	if p.currentCommandIndex >= len(p.currentSequence.Commands) {
		p.endSequence()
		return
	}

	nextCommand := p.currentSequence.Commands[p.currentCommandIndex]
	nextCommand.Init(p.appContext)
}

func (p *SequencePlayer) endSequence() {
	p.isPlaying = false
	if p.currentSequence.BlockPlayerMovement {
		if player, found := p.appContext.ActorManager.GetPlayer(); found {
			player.UnblockMovement()
		}
	}
}
