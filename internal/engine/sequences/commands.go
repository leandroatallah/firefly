package sequences

import (
	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/dialogue"
	"github.com/boilerplate/ebiten-template/internal/engine/event"
)

// EventCommand publishes an event to the global event manager.
type EventCommand struct {
	EventType string
	Payload   map[string]interface{}

	eventManager *event.Manager
}

func (c *EventCommand) Init(appContext any) {
	c.eventManager = appContext.(*app.AppContext).EventManager
	if c.eventManager != nil {
		evt := event.GenericEvent{
			EventType: c.EventType,
			Payload:   c.Payload,
		}
		c.eventManager.Publish(evt)
	}
}

func (c *EventCommand) Update() bool {
	return true
}

// DialogueCommand displays one or more lines of text and waits for player input.
type DialogueCommand struct {
	Lines               []string
	Position            string
	Speed               int
	SpeechID            string
	SpeechAudio         []string
	EnableSpeechSkip    *bool
	EnablePlayerAdvance *bool
	Accumulative        *bool
	HideIndicator       *bool
	dialogueManager     dialogue.Manager
}

func (c *DialogueCommand) Init(appContext any) {
	ctx := appContext.(*app.AppContext)
	c.dialogueManager = ctx.DialogueManager
	speechID := c.SpeechID
	if speechID == "" {
		speechID = dialogue.BubbleSpeechID
	}
	c.dialogueManager.SetActiveSpeech(speechID)

	s := c.dialogueManager.GetActiveSpeech()

	skipEnabled := c.EnableSpeechSkip != nil && *c.EnableSpeechSkip
	c.dialogueManager.SetSpeechSkipEnabled(skipEnabled)
	playerAdvance := c.EnablePlayerAdvance == nil || *c.EnablePlayerAdvance
	c.dialogueManager.SetPlayerAdvanceEnabled(playerAdvance)
	if len(c.SpeechAudio) > 0 {
		c.dialogueManager.SetSpeechAudioQueue(c.SpeechAudio)
	} else {
		c.dialogueManager.ApplyDefaultSpeechAudio(len(c.Lines))
	}

	if c.Accumulative != nil {
		s.SetAccumulative(*c.Accumulative)
	}

	if c.HideIndicator != nil && *c.HideIndicator {
		s.SetIndicator(nil)
	}

	c.dialogueManager.ShowMessages(c.Lines, c.Position, c.Speed)
}

func (c *DialogueCommand) Update() bool {
	// The command is done when the dialogue manager is no longer speaking.
	return !c.dialogueManager.IsSpeaking()
}

type DialogueResetCommand struct {
	dialogueManager dialogue.Manager
}

func (c *DialogueResetCommand) Init(appContext any) {
	c.dialogueManager = appContext.(*app.AppContext).DialogueManager
	c.dialogueManager.ClearSpeechAudioQueue()
	c.dialogueManager.Stop()
}

func (c *DialogueResetCommand) Update() bool {
	return true
}

// DelayCommand waits for a specified number of frames.
type DelayCommand struct {
	Frames int
	timer  int
}

func (c *DelayCommand) Init(appContext any) {
	c.timer = 0
}

func (c *DelayCommand) Update() bool {
	c.timer++
	return c.timer >= c.Frames
}
