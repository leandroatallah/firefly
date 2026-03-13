package sequences

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSequence_Interruptible(t *testing.T) {
	tests := []struct {
		name          string
		interruptible bool
	}{
		{"Interruptible", true},
		{"NonInterruptible", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Sequence{
				interruptible: tt.interruptible,
			}

			if s.Interruptible() != tt.interruptible {
				t.Errorf("Interruptible() = %v, want %v", s.Interruptible(), tt.interruptible)
			}
		})
	}
}

func TestSequence_OneTime(t *testing.T) {
	tests := []struct {
		name    string
		oneTime bool
	}{
		{"OneTime", true},
		{"NotOneTime", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Sequence{
				oneTime: tt.oneTime,
			}

			if s.OneTime() != tt.oneTime {
				t.Errorf("OneTime() = %v, want %v", s.OneTime(), tt.oneTime)
			}
		})
	}
}

func TestSequence_GetPath(t *testing.T) {
	s := &Sequence{
		Path: "assets/sequences/test.json",
	}

	if s.GetPath() != "assets/sequences/test.json" {
		t.Errorf("GetPath() = %q, want %q", s.GetPath(), "assets/sequences/test.json")
	}
}

func TestSequence_ToCommand_Dialogue(t *testing.T) {
	tests := []struct {
		name     string
		cmdData  CommandData
		wantType string
	}{
		{
			name: "dialogue",
			cmdData: CommandData{
				Type:  "dialogue",
				Lines: []string{"Hello"},
			},
			wantType: "*sequences.DialogueCommand",
		},
		{
			name: "delay",
			cmdData: CommandData{
				Type:   "delay",
				Frames: 30,
			},
			wantType: "*sequences.DelayCommand",
		},
		{
			name: "move_actor",
			cmdData: CommandData{
				Type:   "move_actor",
				TargetID: "actor1",
				EndX:   100.0,
			},
			wantType: "*sequences.MoveActorCommand",
		},
		{
			name: "set_speed",
			cmdData: CommandData{
				Type:   "set_speed",
				TargetID: "actor1",
				Speed:  5.0,
			},
			wantType: "*sequences.SetSpeedCommand",
		},
		{
			name: "follow_player",
			cmdData: CommandData{
				Type:           "follow_player",
				TargetID:       "npc1",
				StayOnPlatform: true,
			},
			wantType: "*sequences.FollowPlayerCommand",
		},
		{
			name: "stop_following",
			cmdData: CommandData{
				Type:     "stop_following",
				TargetID: "npc1",
			},
			wantType: "*sequences.StopFollowingCommand",
		},
		{
			name: "remove_actor",
			cmdData: CommandData{
				Type:     "remove_actor",
				TargetID: "actor1",
			},
			wantType: "*sequences.RemoveActorCommand",
		},
		{
			name: "event",
			cmdData: CommandData{
				Type:      "event",
				EventType: "custom_event",
				Payload:   map[string]interface{}{"key": "value"},
			},
			wantType: "*sequences.EventCommand",
		},
		{
			name: "camera_zoom",
			cmdData: CommandData{
				Type:   "camera_zoom",
				Zoom:   2.0,
				Duration: 30,
			},
			wantType: "*sequences.CameraZoomCommand",
		},
		{
			name: "camera_move",
			cmdData: CommandData{
				Type:     "camera_move",
				X:        100.0,
				Y:        200.0,
				Duration: 30,
			},
			wantType: "*sequences.CameraMoveCommand",
		},
		{
			name: "camera_reset",
			cmdData: CommandData{
				Type:        "camera_reset",
				DefaultZoom: 1.0,
			},
			wantType: "*sequences.CameraResetCommand",
		},
		{
			name: "camera_set_target",
			cmdData: CommandData{
				Type:     "camera_set_target",
				TargetID: "actor1",
			},
			wantType: "*sequences.CameraSetTargetCommand",
		},
		{
			name: "call_sequence",
			cmdData: CommandData{
				Type: "call_sequence",
				Path: "nested.json",
			},
			wantType: "*sequences.CallSequenceCommand",
		},
		{
			name: "play_music",
			cmdData: CommandData{
				Type: "play_music",
				Path: "music.mp3",
			},
			wantType: "*sequences.PlayMusicCommand",
		},
		{
			name: "pause_all_music",
			cmdData: CommandData{
				Type: "pause_all_music",
			},
			wantType: "*sequences.PauseAllMusicCommand",
		},
		{
			name: "fadeout_all_music",
			cmdData: CommandData{
				Type:     "fadeout_all_music",
				Duration: 60,
			},
			wantType: "*sequences.FadeOutAllMusicCommand",
		},
		{
			name: "spawn_text",
			cmdData: CommandData{
				Type:     "spawn_text",
				TargetID: "actor1",
				Text:     "Hello",
				Duration: 60,
			},
			wantType: "*sequences.SpawnTextCommand",
		},
		{
			name: "camera_shake",
			cmdData: CommandData{
				Type:   "camera_shake",
				Trauma: 0.5,
			},
			wantType: "*sequences.CameraShakeCommand",
		},
		{
			name: "quake",
			cmdData: CommandData{
				Type:     "quake",
				Trauma:   0.5,
				Duration: 60,
			},
			wantType: "*sequences.QuakeCommand",
		},
		{
			name: "unknown",
			cmdData: CommandData{
				Type: "unknown_command",
			},
			wantType: "nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.cmdData.ToCommand()

			if tt.wantType == "nil" {
				if cmd != nil {
					t.Errorf("ToCommand() for unknown type should return nil, got %T", cmd)
				}
			} else {
				cmdType := getTypeName(cmd)
				if cmdType != tt.wantType {
					t.Errorf("ToCommand() type = %q, want %q", cmdType, tt.wantType)
				}
			}
		})
	}
}

func TestSequence_ToCommand_Dialogue_SpeedFallback(t *testing.T) {
	// Test that Speed field is used as fallback for SpeechSpeed
	cd := CommandData{
		Type:  "dialogue",
		Lines: []string{"Hello"},
		Speed: 10, // Old field name
	}

	cmd := cd.ToCommand()
	dialogueCmd, ok := cmd.(*DialogueCommand)
	if !ok {
		t.Fatalf("expected DialogueCommand, got %T", cmd)
	}

	if dialogueCmd.Speed != 10 {
		t.Errorf("expected Speed to be 10 (from Speed field), got %d", dialogueCmd.Speed)
	}
}

func TestSequence_ToCommand_Dialogue_SpeedPriority(t *testing.T) {
	// Test that SpeechSpeed takes priority over Speed
	cd := CommandData{
		Type:        "dialogue",
		Lines:       []string{"Hello"},
		Speed:       10,
		SpeechSpeed: 5, // New field name (should take priority)
	}

	cmd := cd.ToCommand()
	dialogueCmd, ok := cmd.(*DialogueCommand)
	if !ok {
		t.Fatalf("expected DialogueCommand, got %T", cmd)
	}

	if dialogueCmd.Speed != 5 {
		t.Errorf("expected Speed to be 5 (from SpeechSpeed field), got %d", dialogueCmd.Speed)
	}
}

func TestNewSequenceFromJSON_Valid(t *testing.T) {
	tmpDir := t.TempDir()
	seqPath := filepath.Join(tmpDir, "valid.json")
	seqContent := `{
		"commands": [
			{ "command": "delay", "frames": 30 },
			{ "command": "delay", "frames": 60, "block_sequence": false }
		],
		"block_player_movement": true,
		"interruptible": false,
		"one_time": true
	}`

	if err := os.WriteFile(seqPath, []byte(seqContent), 0644); err != nil {
		t.Fatalf("failed to create sequence file: %v", err)
	}

	seq, err := NewSequenceFromJSON(seqPath)
	if err != nil {
		t.Fatalf("NewSequenceFromJSON() error = %v", err)
	}

	if len(seq.commands) != 2 {
		t.Errorf("expected 2 commands, got %d", len(seq.commands))
	}

	if !seq.BlockPlayerMovement {
		t.Error("expected BlockPlayerMovement to be true")
	}

	if seq.Interruptible() {
		t.Error("expected Interruptible() to be false")
	}

	if !seq.OneTime() {
		t.Error("expected OneTime() to be true")
	}

	if seq.GetPath() != seqPath {
		t.Errorf("expected Path %q, got %q", seqPath, seq.GetPath())
	}

	// Check block_sequence flags
	if len(seq.blockSequenceFlags) != 2 {
		t.Errorf("expected 2 block flags, got %d", len(seq.blockSequenceFlags))
	}
	if !seq.blockSequenceFlags[0] {
		t.Error("expected first command to be blocking")
	}
	if seq.blockSequenceFlags[1] {
		t.Error("expected second command to be non-blocking")
	}
}

func TestNewSequenceFromJSON_Defaults(t *testing.T) {
	tmpDir := t.TempDir()
	seqPath := filepath.Join(tmpDir, "defaults.json")
	seqContent := `{
		"commands": [
			{ "command": "delay", "frames": 30 }
		]
	}`

	if err := os.WriteFile(seqPath, []byte(seqContent), 0644); err != nil {
		t.Fatalf("failed to create sequence file: %v", err)
	}

	seq, err := NewSequenceFromJSON(seqPath)
	if err != nil {
		t.Fatalf("NewSequenceFromJSON() error = %v", err)
	}

	// Default values
	if !seq.Interruptible() {
		t.Error("expected default Interruptible() to be true")
	}

	if seq.OneTime() {
		t.Error("expected default OneTime() to be false")
	}
}

func TestNewSequenceFromJSON_InvalidPath(t *testing.T) {
	seq, err := NewSequenceFromJSON("nonexistent.json")
	if err == nil {
		t.Error("NewSequenceFromJSON() with invalid path should return error")
	}
	if seq == nil {
		t.Error("NewSequenceFromJSON() should return empty Sequence on error, got nil")
	}
}

func TestNewSequenceFromJSON_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	seqPath := filepath.Join(tmpDir, "invalid.json")
	seqContent := `{ invalid json }`

	if err := os.WriteFile(seqPath, []byte(seqContent), 0644); err != nil {
		t.Fatalf("failed to create sequence file: %v", err)
	}

	seq, err := NewSequenceFromJSON(seqPath)
	if err == nil {
		t.Error("NewSequenceFromJSON() with invalid JSON should return error")
	}
	if seq == nil {
		t.Error("NewSequenceFromJSON() should return empty Sequence on error, got nil")
	}
}

func getTypeName(v interface{}) string {
	if v == nil {
		return "nil"
	}
	switch v.(type) {
	case *DialogueCommand:
		return "*sequences.DialogueCommand"
	case *DelayCommand:
		return "*sequences.DelayCommand"
	case *MoveActorCommand:
		return "*sequences.MoveActorCommand"
	case *SetSpeedCommand:
		return "*sequences.SetSpeedCommand"
	case *FollowPlayerCommand:
		return "*sequences.FollowPlayerCommand"
	case *StopFollowingCommand:
		return "*sequences.StopFollowingCommand"
	case *RemoveActorCommand:
		return "*sequences.RemoveActorCommand"
	case *EventCommand:
		return "*sequences.EventCommand"
	case *CameraZoomCommand:
		return "*sequences.CameraZoomCommand"
	case *CameraMoveCommand:
		return "*sequences.CameraMoveCommand"
	case *CameraResetCommand:
		return "*sequences.CameraResetCommand"
	case *CameraSetTargetCommand:
		return "*sequences.CameraSetTargetCommand"
	case *CallSequenceCommand:
		return "*sequences.CallSequenceCommand"
	case *PlayMusicCommand:
		return "*sequences.PlayMusicCommand"
	case *PauseAllMusicCommand:
		return "*sequences.PauseAllMusicCommand"
	case *FadeOutAllMusicCommand:
		return "*sequences.FadeOutAllMusicCommand"
	case *CameraShakeCommand:
		return "*sequences.CameraShakeCommand"
	case *QuakeCommand:
		return "*sequences.QuakeCommand"
	case *SpawnTextCommand:
		return "*sequences.SpawnTextCommand"
	default:
		return "unknown"
	}
}
