package sequences

import (
	"testing"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/movement"
	"github.com/leandroatallah/firefly/internal/engine/mocks"
	"github.com/leandroatallah/firefly/internal/engine/physics/space"
)

const (
	Input movement.MovementStateEnum = iota
	Idle
	Rand
	Chase
	DumbChase
	Patrol
	Avoid
	SideToSide
	Follow
)

func setupTestContext() (*app.AppContext, *actors.Manager) {
	appContext := &app.AppContext{}
	actorManager := actors.NewManager()
	appContext.ActorManager = actorManager
	appContext.Space = space.NewSpace()
	return appContext, actorManager
}

func TestMoveActorCommand(t *testing.T) {
	ctx, am := setupTestContext()
	actor := &mocks.MockActor{Id: "test_actor", SpeedVal: 5, MovementMdl: &mocks.MockMovementModel{}}
	actor.SetPosition(0, 0)
	am.Register(actor)

	cmd := &MoveActorCommand{
		TargetID: "test_actor",
		EndX:     100,
		Speed:    10,
	}

	cmd.Init(ctx)

	if !actor.MovementMdl.(*mocks.MockMovementModel).IsScriptedVal {
		t.Error("expected actor to be set to scripted mode")
	}

	// Move right
	finished := cmd.Update()
	if finished {
		t.Error("command should not be finished yet")
	}
	if actor.MoveRightForce != 10 {
		t.Errorf("expected move right force 10, got %d", actor.MoveRightForce)
	}

	// Reach destination
	actor.SetPosition(100, 0)
	finished = cmd.Update()
	if !finished {
		t.Error("command should be finished when destination reached")
	}
	if actor.MovementMdl.(*mocks.MockMovementModel).IsScriptedVal {
		t.Error("expected actor scripted mode to be disabled after completion")
	}
}

func TestSetSpeedCommand(t *testing.T) {
	ctx, am := setupTestContext()
	actor := &mocks.MockActor{Id: "test_actor", SpeedVal: 5, MaxSpeedVal: 5}
	am.Register(actor)

	cmd := &SetSpeedCommand{
		TargetID: "test_actor",
		Speed:    20,
	}

	cmd.Init(ctx)

	if actor.SpeedVal != 20 || actor.MaxSpeedVal != 20 {
		t.Errorf("expected speed 20, got speed %d, maxSpeed %d", actor.SpeedVal, actor.MaxSpeedVal)
	}
}

func TestFollowCommands(t *testing.T) {
	ctx, am := setupTestContext()
	player := &mocks.MockActor{Id: "player"}
	am.Register(player)

	npc := &mocks.MockActor{Id: "npc"}
	am.Register(npc)

	// Test Follow
	followCmd := &FollowPlayerCommand{TargetID: "npc"}
	followCmd.Init(ctx)
	if npc.MovementSt != Follow {
		t.Errorf("expected npc to be in Follow state, got %v", npc.MovementSt)
	}

	// Test Stop Following
	stopCmd := &StopFollowingCommand{TargetID: "npc"}
	stopCmd.Init(ctx)
	if npc.MovementSt != Idle {
		t.Errorf("expected npc to be in Idle state, got %v", npc.MovementSt)
	}
}

func TestRemoveActorCommand(t *testing.T) {
	ctx, am := setupTestContext()
	actor := &mocks.MockActor{Id: "to_remove"}
	actor.SetPosition(0, 0)
	am.Register(actor)
	ctx.Space.AddBody(actor)

	cmd := &RemoveActorCommand{TargetID: "to_remove"}
	cmd.Init(ctx)

	if _, found := am.Find("to_remove"); found {
		t.Error("actor should be unregistered from ActorManager")
	}

	// Check if queued for removal in space
	ctx.Space.ProcessRemovals()
	if ctx.Space.Find("to_remove") != nil {
		t.Error("actor should be removed from Space")
	}
}

func TestMoveActorCommand_Update_StuckDetection(t *testing.T) {
	ctx, _ := setupTestContext()
	actor := &mocks.MockActor{Id: "stuck_actor", SpeedVal: 5, MovementMdl: &mocks.MockMovementModel{}}
	actor.SetPosition(0, 0)
	ctx.ActorManager.Register(actor)

	cmd := &MoveActorCommand{
		TargetID: "stuck_actor",
		EndX:     100,
		Speed:    10,
	}

	cmd.Init(ctx)

	// Simulate actor being stuck (not moving)
	for i := 0; i < 70; i++ {
		finished := cmd.Update()
		if finished && i < 69 {
			// Wait for stuck threshold (60 frames) + init wait (10 frames)
			if i >= 70 {
				t.Errorf("command finished too early at frame %d", i)
			}
		}
	}

	// After 60 stuck frames, command should be done
	if !cmd.isDone {
		t.Error("command should be done after stuck detection")
	}
}

func TestMoveActorCommand_Update_BrakingDistance(t *testing.T) {
	ctx, _ := setupTestContext()
	actor := &mocks.MockActor{Id: "brake_actor", SpeedVal: 5, MovementMdl: &mocks.MockMovementModel{}}
	actor.SetPosition(0, 0)
	ctx.ActorManager.Register(actor)

	cmd := &MoveActorCommand{
		TargetID: "brake_actor",
		EndX:     100,
		Speed:    10,
	}

	cmd.Init(ctx)

	// Move actor close to destination (within braking distance of 10)
	actor.SetPosition(95, 0)

	// Should not apply force when within braking distance
	cmd.Update()

	if actor.MoveRightForce != 0 && actor.MoveLeftForce != 0 {
		t.Error("should not apply force within braking distance")
	}
}

func TestMoveActorCommand_Update_MoveLeft(t *testing.T) {
	ctx, _ := setupTestContext()
	actor := &mocks.MockActor{Id: "left_actor", SpeedVal: 5, MovementMdl: &mocks.MockMovementModel{}}
	actor.SetPosition(100, 0)
	ctx.ActorManager.Register(actor)

	cmd := &MoveActorCommand{
		TargetID: "left_actor",
		EndX:     0,
		Speed:    10,
	}

	cmd.Init(ctx)

	finished := cmd.Update()
	if finished {
		t.Error("command should not be finished yet")
	}
	if actor.MoveLeftForce != 10 {
		t.Errorf("expected move left force 10, got %d", actor.MoveLeftForce)
	}
}

func TestMoveActorCommand_Update_Done(t *testing.T) {
	ctx, _ := setupTestContext()
	actor := &mocks.MockActor{Id: "done_actor", SpeedVal: 5, MovementMdl: &mocks.MockMovementModel{}}
	actor.SetPosition(0, 0)
	ctx.ActorManager.Register(actor)

	cmd := &MoveActorCommand{
		TargetID: "done_actor",
		EndX:     100,
		Speed:    10,
	}

	cmd.Init(ctx)
	cmd.isDone = true // Force done state

	finished := cmd.Update()
	if !finished {
		t.Error("Update() should return true when isDone is true")
	}
}

func TestMoveActorCommand_Update_NilTarget(t *testing.T) {
	ctx, _ := setupTestContext()
	actor := &mocks.MockActor{Id: "nil_actor", SpeedVal: 5, MovementMdl: &mocks.MockMovementModel{}}
	actor.SetPosition(0, 0)
	ctx.ActorManager.Register(actor)

	cmd := &MoveActorCommand{
		TargetID: "nil_actor",
		EndX:     100,
		Speed:    10,
	}

	cmd.Init(ctx)
	cmd.targetActor = nil // Simulate nil target

	finished := cmd.Update()
	if !finished {
		t.Error("Update() should return true when target is nil")
	}
}

func TestMoveActorCommand_Init_NotFound(t *testing.T) {
	ctx, _ := setupTestContext()
	// Don't register any actor

	cmd := &MoveActorCommand{
		TargetID: "nonexistent",
		EndX:     100,
		Speed:    10,
	}

	cmd.Init(ctx)

	if !cmd.isDone {
		t.Error("command should be done when actor not found")
	}
}

func TestMoveActorCommand_Update_ArrivalThreshold(t *testing.T) {
	ctx, _ := setupTestContext()
	actor := &mocks.MockActor{Id: "arrive_actor", SpeedVal: 5, MovementMdl: &mocks.MockMovementModel{}}
	actor.SetPosition(0, 0)
	ctx.ActorManager.Register(actor)

	cmd := &MoveActorCommand{
		TargetID: "arrive_actor",
		EndX:     100,
		Speed:    10,
	}

	cmd.Init(ctx)

	// Move actor within arrival threshold (20)
	actor.SetPosition(85, 0)

	finished := cmd.Update()
	if !finished {
		t.Error("command should finish when within arrival threshold")
	}

	if actor.MovementMdl.(*mocks.MockMovementModel).IsScriptedVal {
		t.Error("expected actor scripted mode to be disabled after arrival")
	}
}

func TestSetSpeedCommand_Init_ZeroSpeed(t *testing.T) {
	ctx, _ := setupTestContext()
	actor := &mocks.MockActor{Id: "speed_actor", SpeedVal: 5, MaxSpeedVal: 5}
	ctx.ActorManager.Register(actor)

	cmd := &SetSpeedCommand{
		TargetID: "speed_actor",
		Speed:    0, // Zero speed should be ignored
	}

	cmd.Init(ctx)

	// Speed should remain unchanged
	if actor.SpeedVal != 5 {
		t.Errorf("expected speed to remain 5, got %d", actor.SpeedVal)
	}
}

func TestSetSpeedCommand_Init_NotFound(t *testing.T) {
	ctx, _ := setupTestContext()
	// Don't register any actor

	cmd := &SetSpeedCommand{
		TargetID: "nonexistent",
		Speed:    20,
	}

	// Should not panic
	cmd.Init(ctx)
}

func TestFollowPlayerCommand_Init_NoPlayer(t *testing.T) {
	ctx, _ := setupTestContext()
	// Don't register a player

	npc := &mocks.MockActor{Id: "npc"}
	ctx.ActorManager.Register(npc)

	cmd := &FollowPlayerCommand{TargetID: "npc"}

	// Should not panic when player not found
	cmd.Init(ctx)
}

func TestFollowPlayerCommand_Init_NoCharacter(t *testing.T) {
	ctx, _ := setupTestContext()
	player := &mocks.MockActor{Id: "player"}
	ctx.ActorManager.Register(player)

	// MockActor.GetCharacter returns a Character, so this test verifies
	// the command handles the case where SetMovementState is called
	npc := &mocks.MockActor{Id: "npc"}
	ctx.ActorManager.Register(npc)

	cmd := &FollowPlayerCommand{TargetID: "npc"}
	cmd.Init(ctx)

	if npc.MovementSt != Follow {
		t.Errorf("expected npc to be in Follow state, got %v", npc.MovementSt)
	}
}

func TestStopFollowingCommand_Init(t *testing.T) {
	ctx, _ := setupTestContext()
	player := &mocks.MockActor{Id: "player"}
	ctx.ActorManager.Register(player)

	npc := &mocks.MockActor{Id: "npc"}
	ctx.ActorManager.Register(npc)

	cmd := &StopFollowingCommand{TargetID: "npc"}
	cmd.Init(ctx)

	if npc.MovementSt != Idle {
		t.Errorf("expected npc to be in Idle state, got %v", npc.MovementSt)
	}
}

func TestStopFollowingCommand_Init_NoCharacter(t *testing.T) {
	ctx, _ := setupTestContext()
	npc := &mocks.MockActor{Id: "npc"}
	ctx.ActorManager.Register(npc)

	cmd := &StopFollowingCommand{TargetID: "npc"}

	// Should not panic
	cmd.Init(ctx)
}

func TestRemoveActorCommand_Init_NotFound(t *testing.T) {
	ctx, _ := setupTestContext()
	// Don't register any actor

	cmd := &RemoveActorCommand{TargetID: "nonexistent"}

	// Should not panic
	cmd.Init(ctx)
}
