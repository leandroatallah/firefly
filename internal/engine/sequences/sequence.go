package sequences

import (
	"encoding/json"
	"os"

	"github.com/leandroatallah/firefly/internal/engine/core"
)

// Command is an action to be executed in a sequence.
// It can be initialized, and it is updated every frame until it is done.
type Command interface {
	// Init is called once when the command begins.
	// It can be used to set up initial state and get resources from the app context.
	Init(appContext *core.AppContext)

	// Update is called every frame.
	// It should return true when the command is finished.
	Update() bool
}

// Sequence is a list of commands to be executed in order.
type Sequence []Command

// CommandData is a wrapper used for parsing commands from JSON.
// It holds the data for all possible command types.
type CommandData struct {
	Type string `json:"command"`

	// Fields for "dialogue"
	Lines []string `json:"lines,omitempty"`

	// Fields for "delay"
	Frames int `json:"frames,omitempty"`

	// Fields for "move_actor"
	TargetID string  `json:"target_id,omitempty"`
	EndX     float64 `json:"end_x,omitempty"`
	Speed    float64 `json:"speed,omitempty"`
}

// ToCommand converts the generic CommandData into a specific Command implementation.
func (cd *CommandData) ToCommand() Command {
	switch cd.Type {
	case "dialogue":
		return &DialogueCommand{Lines: cd.Lines}
	case "delay":
		return &DelayCommand{Frames: cd.Frames}
	case "move_actor":
		return &MoveActorCommand{
			TargetID: cd.TargetID,
			EndX:     cd.EndX,
			Speed:    cd.Speed,
		}
	}
	return nil
}

// NewSequenceFromJSON loads a sequence from a JSON file path.
func NewSequenceFromJSON(filePath string) (Sequence, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var commandData []CommandData
	if err := json.Unmarshal(data, &commandData); err != nil {
		return nil, err
	}

	var sequence Sequence
	for _, cd := range commandData {
		cmd := cd.ToCommand()
		if cmd != nil {
			sequence = append(sequence, cmd)
		}
	}

	return sequence, nil
}
