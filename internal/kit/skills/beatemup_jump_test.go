package kitskills

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/data/schemas"
	"github.com/boilerplate/ebiten-template/internal/engine/input"
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
	"github.com/boilerplate/ebiten-template/internal/engine/physics/space"
)

// -- Test helpers --------------------------------------------------------

// stubBlocker satisfies physicsmovement.PlayerMovementBlocker.
type stubBlocker struct{ blocked bool }

func (s *stubBlocker) IsMovementBlocked() bool { return s.blocked }

// beatEmUpTestCfg returns a deterministic physics config for the suite.
func beatEmUpTestCfg() *config.AppConfig {
	return &config.AppConfig{Physics: config.PhysicsConfig{
		JumpForce:        8,
		CoyoteTimeFrames: 5,
		JumpBufferFrames: 5,
		UpwardGravity:    2,
		DownwardGravity:  4,
	}}
}

// newBeatEmUpActor builds a body at altitude 0 with multiplier 1.0.
func newBeatEmUpActor() *bodyphysics.ObstacleRect {
	a := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 10, 10))
	a.SetID("actor")
	a.SetJumpForceMultiplier(1.0)
	return a
}

// withCommands installs a CommandsReader stub for the duration of fn.
func withCommands(cmds input.PlayerCommands, fn func()) {
	old := input.CommandsReader
	defer func() { input.CommandsReader = old }()
	input.CommandsReader = func() input.PlayerCommands { return cmds }
	fn()
}

// -- T1: grounded jump fires [AC-2] -------------------------------------

func TestBeatEmUpJumpSkill_GroundedJumpFires(t *testing.T) {
	cfg := beatEmUpTestCfg()
	config.Set(cfg)

	actor := newBeatEmUpActor()
	actor.SetAltitude(0)
	actor.SetVAltitude16(0)

	model := physicsmovement.NewBeatEmUpMovementModel(nil)
	sp := space.NewSpace()

	s := NewBeatEmUpJumpSkill()
	jumpCalls := 0
	s.OnJump = func(b body.MovableCollidable) { jumpCalls++ }

	withCommands(input.PlayerCommands{Jump: true}, func() {
		s.HandleInput(actor, model, sp)
	})

	wantV := -cfg.Physics.JumpForce
	if got := actor.VAltitude16(); got != wantV {
		t.Errorf("VAltitude16 after grounded jump: got %d, want %d", got, wantV)
	}
	if !s.jumpCutPending {
		t.Errorf("jumpCutPending: got false, want true")
	}
	if jumpCalls != 1 {
		t.Errorf("OnJump invocations: got %d, want 1", jumpCalls)
	}
}

// -- T2: no double-jump while airborne [AC-3] ---------------------------

func TestBeatEmUpJumpSkill_NoDoubleJumpWhileAirborne(t *testing.T) {
	cfg := beatEmUpTestCfg()
	config.Set(cfg)

	actor := newBeatEmUpActor()
	actor.SetAltitude(100)
	actor.SetVAltitude16(-50)

	model := physicsmovement.NewBeatEmUpMovementModel(nil)
	sp := space.NewSpace()

	s := NewBeatEmUpJumpSkill()

	withCommands(input.PlayerCommands{Jump: true}, func() {
		s.HandleInput(actor, model, sp)
	})

	if got := actor.VAltitude16(); got != -50 {
		t.Errorf("VAltitude16: got %d, want -50 (unchanged)", got)
	}
	if got := s.jumpBufferCounter; got != cfg.Physics.JumpBufferFrames {
		t.Errorf("jumpBufferCounter: got %d, want %d", got, cfg.Physics.JumpBufferFrames)
	}
}

// -- T3: coyote jump [AC-5] ---------------------------------------------

func TestBeatEmUpJumpSkill_CoyoteJumpFires(t *testing.T) {
	cfg := beatEmUpTestCfg()
	config.Set(cfg)

	actor := newBeatEmUpActor()
	actor.SetAltitude(20)
	actor.SetVAltitude16(0)

	model := physicsmovement.NewBeatEmUpMovementModel(nil)
	sp := space.NewSpace()

	s := NewBeatEmUpJumpSkill()
	s.coyoteTimeCounter = 3

	withCommands(input.PlayerCommands{Jump: true}, func() {
		s.HandleInput(actor, model, sp)
	})

	wantV := -cfg.Physics.JumpForce
	if got := actor.VAltitude16(); got != wantV {
		t.Errorf("VAltitude16 after coyote jump: got %d, want %d", got, wantV)
	}
	if s.coyoteTimeCounter != 0 {
		t.Errorf("coyoteTimeCounter after jump: got %d, want 0", s.coyoteTimeCounter)
	}
}

// -- T4: jump buffered → fires on landing [AC-4] -------------------------

func TestBeatEmUpJumpSkill_BufferedJumpFiresOnLanding(t *testing.T) {
	cfg := beatEmUpTestCfg()
	config.Set(cfg)

	actor := newBeatEmUpActor()
	actor.SetAltitude(0) // just landed
	actor.SetVAltitude16(0)

	model := physicsmovement.NewBeatEmUpMovementModel(nil)

	s := NewBeatEmUpJumpSkill()
	s.jumpBufferCounter = 5
	jumpCalls := 0
	s.OnJump = func(b body.MovableCollidable) { jumpCalls++ }

	s.Update(actor, model)

	wantV := -cfg.Physics.JumpForce
	if got := actor.VAltitude16(); got != wantV {
		t.Errorf("VAltitude16 after buffered fire: got %d, want %d", got, wantV)
	}
	if s.jumpBufferCounter != 0 {
		t.Errorf("jumpBufferCounter after fire: got %d, want 0", s.jumpBufferCounter)
	}
	if jumpCalls != 1 {
		t.Errorf("OnJump invocations: got %d, want 1", jumpCalls)
	}
}

// -- T5: jump-cut applies multiplier [AC-6] ------------------------------

func TestBeatEmUpJumpSkill_JumpCutAppliesMultiplier(t *testing.T) {
	cfg := beatEmUpTestCfg()
	config.Set(cfg)

	actor := newBeatEmUpActor()
	actor.SetAltitude(50)
	actor.SetVAltitude16(-320)

	model := physicsmovement.NewBeatEmUpMovementModel(nil)
	sp := space.NewSpace()

	s := NewBeatEmUpJumpSkill()
	s.SetJumpCutMultiplier(0.5)
	s.jumpCutPending = true
	s.jumpPressed = true // previous frame was holding

	// Trailing edge: release Jump.
	withCommands(input.PlayerCommands{Jump: false}, func() {
		s.HandleInput(actor, model, sp)
	})

	if got := actor.VAltitude16(); got != -160 {
		t.Errorf("VAltitude16 after jump-cut: got %d, want -160", got)
	}
	if s.jumpCutPending {
		t.Errorf("jumpCutPending after cut: got true, want false")
	}
}

// -- T6: no-op when model is *PlatformMovementModel [AC-1] ---------------

func TestBeatEmUpJumpSkill_NoOpOnPlatformModel(t *testing.T) {
	cfg := beatEmUpTestCfg()
	cfg.Physics.MaxFallSpeed = 100
	config.Set(cfg)

	actor := newBeatEmUpActor()
	actor.SetAltitude(0)
	actor.SetVAltitude16(0)

	model := physicsmovement.NewPlatformMovementModel(nil)
	sp := space.NewSpace()

	s := NewBeatEmUpJumpSkill()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("HandleInput panicked with PlatformMovementModel: %v", r)
		}
	}()

	withCommands(input.PlayerCommands{Jump: true}, func() {
		s.HandleInput(actor, model, sp)
	})

	if got := actor.VAltitude16(); got != 0 {
		t.Errorf("VAltitude16 must be unchanged for non-BeatEmUp model: got %d, want 0", got)
	}
}

// -- T7: input-blocked guard [AC-8] --------------------------------------

func TestBeatEmUpJumpSkill_InputBlockedGuard(t *testing.T) {
	cfg := beatEmUpTestCfg()
	config.Set(cfg)

	actor := newBeatEmUpActor()
	actor.SetAltitude(0)
	actor.SetVAltitude16(0)

	blocker := &stubBlocker{blocked: true}
	model := physicsmovement.NewBeatEmUpMovementModel(blocker)
	sp := space.NewSpace()

	s := NewBeatEmUpJumpSkill()

	withCommands(input.PlayerCommands{Jump: true}, func() {
		s.HandleInput(actor, model, sp)
	})

	if got := actor.VAltitude16(); got != 0 {
		t.Errorf("VAltitude16 must be untouched when input blocked: got %d, want 0", got)
	}
	if s.jumpPressed {
		t.Errorf("jumpPressed must not advance when input blocked: got true, want false")
	}
	if s.jumpCutPending {
		t.Errorf("jumpCutPending must not be set when input blocked")
	}
}

// -- T8: SetJumpCutMultiplier clamp [AC-7] -------------------------------

func TestBeatEmUpJumpSkill_SetJumpCutMultiplier(t *testing.T) {
	tests := []struct {
		input float64
		want  float64
	}{
		{0.5, 0.5},
		{1.0, 1.0},
		{0.0, 0.1},
		{-1.0, 0.1},
		{1.5, 1.0},
	}
	for _, tc := range tests {
		s := NewBeatEmUpJumpSkill()
		s.SetJumpCutMultiplier(tc.input)
		if s.jumpCutMultiplier != tc.want {
			t.Errorf("SetJumpCutMultiplier(%v): got %v, want %v", tc.input, s.jumpCutMultiplier, tc.want)
		}
	}
}

// -- T9: coyote decrement while airborne [AC-5] --------------------------

func TestBeatEmUpJumpSkill_CoyoteDecrementsAirborne(t *testing.T) {
	cfg := beatEmUpTestCfg()
	config.Set(cfg)

	actor := newBeatEmUpActor()
	actor.SetAltitude(10)
	actor.SetVAltitude16(0)

	model := physicsmovement.NewBeatEmUpMovementModel(nil)

	s := NewBeatEmUpJumpSkill()
	s.coyoteTimeCounter = 2

	s.Update(actor, model)

	if s.coyoteTimeCounter != 1 {
		t.Errorf("coyoteTimeCounter after airborne Update: got %d, want 1", s.coyoteTimeCounter)
	}
}

// -- T10: coyote reset while grounded [AC-5] -----------------------------

func TestBeatEmUpJumpSkill_CoyoteResetsGrounded(t *testing.T) {
	cfg := beatEmUpTestCfg()
	config.Set(cfg)

	actor := newBeatEmUpActor()
	actor.SetAltitude(0)
	actor.SetVAltitude16(0)

	model := physicsmovement.NewBeatEmUpMovementModel(nil)

	s := NewBeatEmUpJumpSkill()
	s.coyoteTimeCounter = 0

	s.Update(actor, model)

	if s.coyoteTimeCounter != cfg.Physics.CoyoteTimeFrames {
		t.Errorf("coyoteTimeCounter on grounded Update: got %d, want %d", s.coyoteTimeCounter, cfg.Physics.CoyoteTimeFrames)
	}
}

// -- T11: Freeze blocks Update mutation [AC-9] ---------------------------

func TestBeatEmUpJumpSkill_FreezeBlocksUpdate(t *testing.T) {
	cfg := beatEmUpTestCfg()
	config.Set(cfg)

	actor := newBeatEmUpActor()
	actor.SetAltitude(10)
	actor.SetVAltitude16(-100)
	actor.SetFreeze(true)

	model := physicsmovement.NewBeatEmUpMovementModel(nil)

	s := NewBeatEmUpJumpSkill()
	s.coyoteTimeCounter = 2
	s.jumpBufferCounter = 2

	s.Update(actor, model)

	if s.coyoteTimeCounter != 2 {
		t.Errorf("coyoteTimeCounter must not advance while frozen: got %d, want 2", s.coyoteTimeCounter)
	}
	if s.jumpBufferCounter != 2 {
		t.Errorf("jumpBufferCounter must not advance while frozen: got %d, want 2", s.jumpBufferCounter)
	}
	if got := actor.VAltitude16(); got != -100 {
		t.Errorf("VAltitude16 must not change while frozen: got %d, want -100", got)
	}
}

// -- T12: force <= 0 skips silently --------------------------------------

func TestBeatEmUpJumpSkill_ZeroForceSkipsSilently(t *testing.T) {
	cfg := beatEmUpTestCfg()
	config.Set(cfg)

	actor := newBeatEmUpActor()
	actor.SetAltitude(0)
	actor.SetVAltitude16(0)
	actor.SetJumpForceMultiplier(0)

	model := physicsmovement.NewBeatEmUpMovementModel(nil)
	sp := space.NewSpace()

	s := NewBeatEmUpJumpSkill()
	jumpCalls := 0
	s.OnJump = func(b body.MovableCollidable) { jumpCalls++ }

	withCommands(input.PlayerCommands{Jump: true}, func() {
		s.HandleInput(actor, model, sp)
	})

	if got := actor.VAltitude16(); got != 0 {
		t.Errorf("VAltitude16 must remain 0 when force <= 0: got %d", got)
	}
	if s.jumpCutPending {
		t.Errorf("jumpCutPending must not be set when force <= 0")
	}
	if jumpCalls != 0 {
		t.Errorf("OnJump must not fire when force <= 0: got %d calls", jumpCalls)
	}
}

// -- T13: factory selects BeatEmUpJumpSkill for eight_dir [AC-11] --------

func TestFromConfig_BeatEmUpJumpForEightDir(t *testing.T) {
	cfg := &schemas.SkillsConfig{
		Movement: &schemas.MovementConfig{Enabled: ptrBool(true), Mode: schemas.MovementModeEightDir},
		Jump:     &schemas.JumpConfig{Enabled: ptrBool(true)},
	}

	skills := FromConfig(cfg, SkillDeps{})

	var foundBeatEmUp, foundPlatform bool
	for _, sk := range skills {
		if _, ok := sk.(*BeatEmUpJumpSkill); ok {
			foundBeatEmUp = true
		}
		if _, ok := sk.(*JumpSkill); ok {
			foundPlatform = true
		}
	}
	if !foundBeatEmUp {
		t.Errorf("expected *BeatEmUpJumpSkill in skills for mode=eight_dir; got skills=%v", skills)
	}
	if foundPlatform {
		t.Errorf("must NOT instantiate *JumpSkill for mode=eight_dir")
	}
}

// -- T14: factory selects JumpSkill for default/horizontal [AC-11] -------

func TestFromConfig_JumpSkillForHorizontal(t *testing.T) {
	tests := []struct {
		name string
		mode string
	}{
		{"empty mode (default)", ""},
		{"horizontal mode", schemas.MovementModeHorizontal},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &schemas.SkillsConfig{
				Movement: &schemas.MovementConfig{Enabled: ptrBool(true), Mode: tc.mode},
				Jump:     &schemas.JumpConfig{Enabled: ptrBool(true)},
			}

			skills := FromConfig(cfg, SkillDeps{})

			var foundBeatEmUp, foundPlatform bool
			for _, sk := range skills {
				if _, ok := sk.(*BeatEmUpJumpSkill); ok {
					foundBeatEmUp = true
				}
				if _, ok := sk.(*JumpSkill); ok {
					foundPlatform = true
				}
			}
			if foundBeatEmUp {
				t.Errorf("must NOT instantiate *BeatEmUpJumpSkill for mode=%q", tc.mode)
			}
			if !foundPlatform {
				t.Errorf("expected *JumpSkill in skills for mode=%q", tc.mode)
			}
		})
	}
}
