package sequences

import (
	"encoding/json"
	"os"

	"github.com/leandroatallah/firefly/internal/engine/contracts/sequences"
)

type Sequence struct {
	commands            []sequences.Command
	BlockPlayerMovement bool
	interruptible       bool
	oneTime             bool
	Path                string
	blockSequenceFlags  []bool
}

func (s *Sequence) Commands() []sequences.Command {
	return s.commands
}

// Interruptible returns whether this sequence can be interrupted by other sequences.
func (s *Sequence) Interruptible() bool {
	return s.interruptible
}

// OneTime returns whether this sequence can only be played once.
func (s *Sequence) OneTime() bool {
	return s.oneTime
}

// GetPath returns the path of this sequence.
func (s *Sequence) GetPath() string {
	return s.Path
}

// CommandData is a wrapper used for parsing commands from JSON.
// It holds the data for all possible command types.
type CommandData struct {
	Type string `json:"command"`

	// Fields for "dialogue"
	Lines            []string `json:"lines,omitempty"`
	Position         string   `json:"position,omitempty"`
	SpeechSpeed      int      `json:"speech_speed,omitempty"`
	SpeechID         string   `json:"speech_id,omitempty"`
	SpeechAudio      []string `json:"speech_audio,omitempty"`
	EnableSpeechSkip *bool    `json:"enable_speech_skip,omitempty"`
	Accumulative     *bool    `json:"accumulative,omitempty"`

	// Fields for "delay"
	Frames int `json:"frames,omitempty"`

	// Fields for "move_actor"
	TargetID string  `json:"target_id,omitempty"`
	EndX     float64 `json:"end_x,omitempty"`
	Speed    float64 `json:"speed,omitempty"`

	// Fields for "follow_player"
	StayOnPlatform bool `json:"stay_on_platform,omitempty"`

	// Per-command control over whether this command blocks the sequence timeline.
	// If omitted, commands are treated as blocking (current default behavior).
	BlockSequence *bool `json:"block_sequence,omitempty"`

	// Fields for "event"
	EventType string                 `json:"event_type,omitempty"`
	Payload   map[string]interface{} `json:"payload,omitempty"`

	// Fields for "camera_zoom"
	Zoom        float64 `json:"zoom,omitempty"`
	Duration    int     `json:"duration,omitempty"`
	Delay       int     `json:"delay,omitempty"`
	OutDuration int     `json:"out_duration,omitempty"`

	// Fields for "camera_move"
	X      float64 `json:"x,omitempty"`
	Y      float64 `json:"y,omitempty"`
	Smooth bool    `json:"smooth,omitempty"`

	// Fields for "camera_reset"
	DefaultZoom float64 `json:"default_zoom,omitempty"`

	// Fields for "call_sequence" and "play_music"
	Path string `json:"path,omitempty"`

	// Fields for "play_music"
	MusicRewind bool    `json:"rewind,omitempty"`
	Volume      float64 `json:"volume,omitempty"`
	Loop        bool    `json:"loop,omitempty"`

	// Fields for "spawn_text"
	Text     string `json:"text,omitempty"`
	TextType string `json:"type,omitempty"` // For spawn_text: "overhead" or "screen"

	// Fields for "camera_shake"
	Trauma float64 `json:"trauma,omitempty"`
}

// SequenceData is a wrapper used for parsing a full sequence from JSON.
type SequenceData struct {
	Commands            []CommandData `json:"commands"`
	BlockPlayerMovement bool          `json:"block_player_movement,omitempty"`
	Interruptible       *bool         `json:"interruptible,omitempty"`
	OneTime             *bool         `json:"one_time,omitempty"`
}

// ToCommand converts the generic CommandData into a specific Command implementation.
func (cd *CommandData) ToCommand() sequences.Command {
	switch cd.Type {
	case "dialogue":
		speed := cd.SpeechSpeed
		if speed == 0 && cd.Speed > 0 {
			speed = int(cd.Speed)
		}
		return &DialogueCommand{Lines: cd.Lines, Position: cd.Position, Speed: speed, SpeechID: cd.SpeechID, SpeechAudio: cd.SpeechAudio, EnableSpeechSkip: cd.EnableSpeechSkip, Accumulative: cd.Accumulative}
	case "delay":
		return &DelayCommand{Frames: cd.Frames}
	case "move_actor":
		return &MoveActorCommand{
			TargetID: cd.TargetID,
			EndX:     cd.EndX,
			Speed:    cd.Speed,
		}
	case "set_speed":
		return &SetSpeedCommand{TargetID: cd.TargetID, Speed: cd.Speed}
	case "follow_player":
		return &FollowPlayerCommand{
			TargetID:       cd.TargetID,
			StayOnPlatform: cd.StayOnPlatform,
		}
	case "stop_following":
		return &StopFollowingCommand{
			TargetID: cd.TargetID,
		}
	case "remove_actor":
		return &RemoveActorCommand{
			TargetID: cd.TargetID,
		}
	case "event":
		return &EventCommand{
			EventType: cd.EventType,
			Payload:   cd.Payload,
		}
	case "camera_zoom":
		return &CameraZoomCommand{
			Zoom:        cd.Zoom,
			Duration:    cd.Duration,
			Delay:       cd.Delay,
			OutDuration: cd.OutDuration,
			TargetID:    cd.TargetID,
		}
	case "camera_move":
		return &CameraMoveCommand{
			X:        cd.X,
			Y:        cd.Y,
			Duration: cd.Duration,
			Smooth:   cd.Smooth,
		}
	case "camera_reset":
		return &CameraResetCommand{
			DefaultZoom: cd.DefaultZoom,
			Duration:    cd.Duration,
		}
	case "camera_set_target":
		return &CameraSetTargetCommand{
			TargetID: cd.TargetID,
			Duration: cd.Duration,
		}
	case "camera_shake":
		return &CameraShakeCommand{
			Trauma: cd.Trauma,
		}
	case "quake":
		return &QuakeCommand{
			Trauma:   cd.Trauma,
			Duration: cd.Duration,
		}
	case "call_sequence":
		return &CallSequenceCommand{
			Path: cd.Path,
		}
	case "play_music":
		return &PlayMusicCommand{
			Path:   cd.Path,
			Rewind: cd.MusicRewind,
			Volume: cd.Volume,
			Loop:   cd.Loop,
		}
	case "pause_all_music":
		return &PauseAllMusicCommand{}
	case "fadeout_all_music":
		return &FadeOutAllMusicCommand{
			Duration: cd.Duration,
		}
	case "spawn_text":
		return &SpawnTextCommand{
			TargetID: cd.TargetID,
			Text:     cd.Text,
			Duration: cd.Duration,
			Type:     cd.TextType,
			X:        cd.X,
			Y:        cd.Y,
		}
	}
	return nil
}

// NewSequenceFromJSON loads a sequence from a JSON file path.
func NewSequenceFromJSON(filePath string) (*Sequence, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return &Sequence{}, err
	}

	var sequenceData SequenceData
	if err := json.Unmarshal(data, &sequenceData); err != nil {
		return &Sequence{}, err
	}

	var commands []sequences.Command
	var flags []bool
	for _, cd := range sequenceData.Commands {
		cmd := cd.ToCommand()
		if cmd != nil {
			commands = append(commands, cmd)
			block := true
			if cd.BlockSequence != nil {
				block = *cd.BlockSequence
			}
			flags = append(flags, block)
		}
	}

	// Default values for optional fields
	interruptible := true
	if sequenceData.Interruptible != nil {
		interruptible = *sequenceData.Interruptible
	}
	oneTime := false
	if sequenceData.OneTime != nil {
		oneTime = *sequenceData.OneTime
	}

	return &Sequence{
		commands:            commands,
		BlockPlayerMovement: sequenceData.BlockPlayerMovement,
		interruptible:       interruptible,
		oneTime:             oneTime,
		Path:                filePath,
		blockSequenceFlags:  flags,
	}, nil
}
