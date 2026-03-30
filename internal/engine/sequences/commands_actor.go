package sequences

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/movement"
)

// resolveActorTargets returns actors matching a target ID or a @query: regex pattern
func resolveActorTargets(ctx *app.AppContext, targetID string) []actors.ActorEntity {
	targets := []actors.ActorEntity{}
	if strings.HasPrefix(targetID, "@query:") {
		pattern := strings.TrimPrefix(targetID, "@query:")
		re, err := regexp.Compile(pattern)
		if err != nil {
			return targets
		}
		for _, b := range ctx.Space.Bodies() {
			if re.MatchString(b.ID()) {
				if a, ok := ctx.ActorManager.Find(b.ID()); ok {
					targets = append(targets, a)
				}
			}
		}
		return targets
	}
	if a, ok := ctx.ActorManager.Find(targetID); ok {
		targets = append(targets, a)
	}
	return targets
}

// MoveActorCommand moves a target actor to a specified X position.
type MoveActorCommand struct {
	TargetID string
	EndX     float64
	Speed    float64

	targetActor actors.ActorEntity
	isDone      bool

	lastX         float64
	stuckFrames   int
	initWaitCount int
}

func (c *MoveActorCommand) Init(appContext any) {
	actor, found := appContext.(*app.AppContext).ActorManager.Find(c.TargetID)
	if !found {
		fmt.Printf("MoveActorCommand: Actor with ID '%s' not found.\n", c.TargetID)
		c.isDone = true
		return
	}
	c.targetActor = actor

	if model := actor.MovementModel(); model != nil {
		model.SetIsScripted(true)
	}

	c.lastX = float64(actor.Position().Min.X)
	c.stuckFrames = 0
	c.initWaitCount = 0
}

func (c *MoveActorCommand) Update() bool {
	if c.isDone || c.targetActor == nil {
		return true
	}

	currentX := float64(c.targetActor.Position().Min.X)

	// Stuck detection: if we haven't moved significantly for a while, finish the command.
	// We wait a few frames (initWaitCount) to allow physics to kick in.
	if c.initWaitCount < 10 {
		c.initWaitCount++
	} else {
		if math.Abs(currentX-c.lastX) < 0.1 {
			c.stuckFrames++
		} else {
			c.stuckFrames = 0
		}
	}
	c.lastX = currentX

	const stuckThreshold = 60 // 1 second at 60fps
	if c.stuckFrames >= stuckThreshold {
		c.isDone = true
	}

	speed := c.Speed
	if speed == 0 {
		speed = float64(c.targetActor.Speed())
	}

	distance := c.EndX - currentX

	const arrivalThreshold = 20.0
	const brakingDistance = 10.0 // This value may need tuning depending on friction and speed

	if c.isDone || math.Abs(distance) < arrivalThreshold {
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
			c.targetActor.OnMoveRight(int(speed))
		} else {
			c.targetActor.OnMoveLeft(int(speed))
		}
	}

	return false
}

// SetSpeedCommand sets speed and max speed for one or more targets
type SetSpeedCommand struct {
	TargetID string
	Speed    float64
}

func (c *SetSpeedCommand) Init(appContext any) {
	ctx := appContext.(*app.AppContext)
	if c.Speed <= 0 {
		return
	}
	for _, a := range resolveActorTargets(ctx, c.TargetID) {
		if m, ok := any(a).(body.Movable); ok {
			spd := int(c.Speed)
			m.SetSpeed(spd)
			m.SetMaxSpeed(spd)
		}
	}
}

func (c *SetSpeedCommand) Update() bool { return true }

// FollowActorCommand sets movement state to Follow targeting a specific actor by ID.
type FollowActorCommand struct {
	TargetID       string
	SubjectID      string
	StayOnPlatform bool
}

func (c *FollowActorCommand) Init(appContext any) {
	ctx := appContext.(*app.AppContext)
	subject, ok := ctx.ActorManager.Find(c.SubjectID)
	if !ok {
		return
	}
	for _, a := range resolveActorTargets(ctx, c.TargetID) {
		ch := a.GetCharacter()
		if ch != nil {
			ch.ClearSkills()
			if c.StayOnPlatform {
				a.SetMovementState(movement.Follow, subject, movement.WithPlatformFollow())
			} else {
				a.SetMovementState(movement.Follow, subject)
			}
		}
	}
}

func (c *FollowActorCommand) Update() bool { return true }

// FollowPlayerCommand sets movement state to Follow targeting the primary player.
// It is a convenience wrapper around FollowActorCommand for single-player games.
type FollowPlayerCommand struct {
	TargetID       string
	StayOnPlatform bool
}

func (c *FollowPlayerCommand) Init(appContext any) {
	ctx := appContext.(*app.AppContext)
	player, ok := ctx.ActorManager.GetPlayer()
	if !ok {
		return
	}
	cmd := &FollowActorCommand{
		TargetID:       c.TargetID,
		SubjectID:      player.ID(),
		StayOnPlatform: c.StayOnPlatform,
	}
	cmd.Init(appContext)
}

func (c *FollowPlayerCommand) Update() bool { return true }

type StopFollowingCommand struct {
	TargetID       string
	StayOnPlatform bool
}

func (c *StopFollowingCommand) Init(appContext any) {
	ctx := appContext.(*app.AppContext)
	for _, a := range resolveActorTargets(ctx, c.TargetID) {
		ch := a.GetCharacter()
		if ch != nil {
			a.SetMovementState(movement.Idle, nil)
		}
	}
}

func (c *StopFollowingCommand) Update() bool { return true }

// RemoveActorCommand removes an actor from the actor manager and the physics space.
type RemoveActorCommand struct {
	TargetID string
}

func (c *RemoveActorCommand) Init(appContext any) {
	ctx := appContext.(*app.AppContext)
	for _, a := range resolveActorTargets(ctx, c.TargetID) {
		ctx.ActorManager.Unregister(a)
		ctx.Space.QueueForRemoval(a)
	}
}

func (c *RemoveActorCommand) Update() bool { return true }
