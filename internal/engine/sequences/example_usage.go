//go:build ignore

// This file provides an example of how to use the sequences package.
// It is not part of the build.

package sequences

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/leandroatallah/firefly/internal/engine/core"
)

// The following is an example of how to integrate the SequencePlayer into a scene.

// 1. Add SequencePlayer to your scene struct.
// Make sure your scene has access to *app.AppContext
type SceneWithSequence struct {
	// ... other scene fields
	sequencePlayer *SequencePlayer
	AppContext     *app.AppContext
}

// 2. Initialize the SequencePlayer in your scene's startup method.
func (s *SceneWithSequence) OnStart() {
	// ... other startup logic
	s.sequencePlayer = NewSequencePlayer(s.AppContext)
}

// 3. Update the SequencePlayer in your scene's update loop.
func (s *SceneWithSequence) Update() error {
	// ... other update logic

	s.sequencePlayer.Update()

	// 4. Trigger a sequence.
	// For example, on a key press.
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		// Make sure the JSON file path is correct.
		sequence, err := NewSequenceFromJSON("assets/sequences/sample.json")
		if err != nil {
			log.Fatal(err)
		}
		s.sequencePlayer.Play(sequence)
	}

	return nil
}
