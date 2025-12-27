package sequences

import (
	"fmt"
	"math"

	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/systems/speech"
)

// DialogueCommand displays one or more lines of text and waits for player input.
type DialogueCommand struct {
	Lines           []string
	dialogueManager *speech.Manager
}

func (c *DialogueCommand) Init(appContext *app.AppContext) {
	c.dialogueManager = appContext.DialogueManager
	c.dialogueManager.ShowMessages(c.Lines)
}

func (c *DialogueCommand) Update() bool {
	// The command is done when the dialogue manager is no longer speaking.
	return !c.dialogueManager.IsSpeaking()
}

// DelayCommand waits for a specified number of frames.
type DelayCommand struct {
	Frames int
	timer  int
}

func (c *DelayCommand) Init(appContext *app.AppContext) {
	c.timer = 0
}

func (c *DelayCommand) Update() bool {
	c.timer++
	return c.timer >= c.Frames
}

// MoveActorCommand moves a target actor to a specified X position.
type MoveActorCommand struct {
	TargetID string
	EndX     float64
	Speed    float64

	targetActor actors.ActorEntity
	isDone      bool
}

func (c *MoveActorCommand) Init(appContext *app.AppContext) {
	actor, found := appContext.ActorManager.Find(c.TargetID)
	if !found {
		fmt.Printf("MoveActorCommand: Actor with ID '%s' not found.\n", c.TargetID)
		c.isDone = true
		return
	}
	c.targetActor = actor

	if model := actor.MovementModel(); model != nil {
		model.SetIsScripted(true)
	}
}

func (c *MoveActorCommand) Update() bool {
	if c.isDone || c.targetActor == nil {
		return true
	}

	currentX := float64(c.targetActor.Position().Min.X)
	distance := c.EndX - currentX

	const arrivalThreshold = 20.0
	const brakingDistance = 10.0 // This value may need tuning depending on friction and speed

	if math.Abs(distance) < arrivalThreshold {
		c.isDone = true
		// Restore player control before finishing the command.
		if model := c.targetActor.MovementModel(); model != nil {
			model.SetIsScripted(false)
		}
		return true
	}

	// When we are close, stop applying force and let friction do the work.
	if math.Abs(distance) < brakingDistance {
		// No-op, just wait for the actor to glide to a stop.
	} else {
		// Apply force to move towards the target.
		if distance > 0 {
			c.targetActor.OnMoveRight(int(c.Speed))
		} else {
			c.targetActor.OnMoveLeft(int(c.Speed))
		}
	}

	return false
}
