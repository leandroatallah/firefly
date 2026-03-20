package sequences

import (
	"path/filepath"

	"github.com/leandroatallah/firefly/internal/engine/app"
)

// CallSequenceCommand calls a nested sequence from a JSON file path.
// The sequence runs to completion before this command finishes (unless block_sequence is false).
type CallSequenceCommand struct {
	Path string `json:"path"` // Path to the nested sequence JSON file

	nestedSequence *Sequence
	sequencePlayer *SequencePlayer
	isComplete     bool
	blockSequence  *bool
}

func (c *CallSequenceCommand) Init(appContext any) {
	ctx := appContext.(*app.AppContext)

	// Resolve the path relative to the working directory or use absolute path
	sequencePath := c.Path
	if !filepath.IsAbs(sequencePath) {
		// Assume path is relative to assets/sequences/
		sequencePath = filepath.Join("assets/sequences", c.Path)
	}

	var err error
	c.nestedSequence, err = NewSequenceFromJSON(ctx.Assets, sequencePath)
	if err != nil {
		// If loading fails, mark as complete so the sequence continues
		c.isComplete = true
		return
	}

	// Create a sequence player if not already available
	if c.sequencePlayer == nil {
		c.sequencePlayer = NewSequencePlayer(ctx)
	}

	// Mark that this player is blocked by parent sequence
	c.sequencePlayer.blockedByParent = true

	c.sequencePlayer.Play(c.nestedSequence)
	c.isComplete = false
}

func (c *CallSequenceCommand) Update() bool {
	if c.isComplete {
		return true
	}

	if c.sequencePlayer == nil {
		c.isComplete = true
		return true
	}

	// Check if the nested sequence has finished
	if c.sequencePlayer.IsOver() || !c.sequencePlayer.IsPlaying() {
		c.isComplete = true
		return true
	}

	// Update the nested sequence player
	c.sequencePlayer.Update()

	return false
}
