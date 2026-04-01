package skill

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/fp16"
)

func TestJumpSkill_SetJumpCutMultiplier(t *testing.T) {
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
		s := NewJumpSkill()
		s.SetJumpCutMultiplier(tc.input)
		if s.jumpCutMultiplier != tc.want {
			t.Errorf("SetJumpCutMultiplier(%v): got %v, want %v", tc.input, s.jumpCutMultiplier, tc.want)
		}
	}
}

func TestJumpSkill_JumpCut(t *testing.T) {
	cfg := &config.AppConfig{Physics: config.PhysicsConfig{
		JumpForce:        8,
		CoyoteTimeFrames: 5,
		JumpBufferFrames: 5,
		DownwardGravity:  4,
	}}
	config.Set(cfg)

	// vy16Before values: negative = going up, positive = falling
	tests := []struct {
		name              string
		jumpCutMultiplier float64
		vy16Before        int
		applyRelease      bool
		wantVy16After     int
		wantPending       bool
	}{
		{
			name:              "full hold — no cut",
			jumpCutMultiplier: 1.0,
			vy16Before:        -320,
			applyRelease:      false,
			wantVy16After:     -320,
			wantPending:       true, // pending stays true; no release
		},
		{
			name:              "short press — cut applied",
			jumpCutMultiplier: 0.5,
			vy16Before:        -320,
			applyRelease:      true,
			wantVy16After:     -160,
			wantPending:       false,
		},
		{
			name:              "release while falling — no cut",
			jumpCutMultiplier: 0.5,
			vy16Before:        fp16.To16(5), // positive = falling
			applyRelease:      true,
			wantVy16After:     fp16.To16(5),
			wantPending:       false,
		},
		{
			name:              "multiplier clamped below zero",
			jumpCutMultiplier: -1.0,
			vy16Before:        -320,
			applyRelease:      true,
			wantVy16After:     int(float64(-320) * 0.1),
			wantPending:       false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actor := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 10, 10))
			actor.SetID("actor")
			actor.SetJumpForceMultiplier(1.0)

			s := NewJumpSkill()
			s.SetJumpCutMultiplier(tc.jumpCutMultiplier)
			s.jumpCutPending = true
			actor.SetVelocity(0, tc.vy16Before)

			if tc.applyRelease {
				s.applyJumpCut(actor)
			}

			_, gotVy := actor.Velocity()
			if gotVy != tc.wantVy16After {
				t.Errorf("vy16: got %d, want %d", gotVy, tc.wantVy16After)
			}
			if s.jumpCutPending != tc.wantPending {
				t.Errorf("jumpCutPending: got %v, want %v", s.jumpCutPending, tc.wantPending)
			}
		})
	}
}

func TestJumpSkill_JumpCut_AppliedOnlyOnce(t *testing.T) {
	cfg := &config.AppConfig{Physics: config.PhysicsConfig{DownwardGravity: 4}}
	config.Set(cfg)

	actor := bodyphysics.NewObstacleRect(bodyphysics.NewRect(0, 0, 10, 10))
	actor.SetID("actor")
	actor.SetVelocity(0, -320)

	s := NewJumpSkill()
	s.SetJumpCutMultiplier(0.5)
	s.jumpCutPending = true

	// First release: cut applied, pending cleared.
	s.applyJumpCut(actor)
	_, vy := actor.Velocity()
	if vy != -160 {
		t.Fatalf("after first cut: got %d, want -160", vy)
	}
	if s.jumpCutPending {
		t.Fatal("jumpCutPending should be false after cut")
	}

	// HandleInput guard: pending is false, so a second release is a no-op.
	// Simulate what HandleInput does: only call applyJumpCut when pending.
	if s.jumpCutPending {
		s.applyJumpCut(actor)
	}
	_, vy = actor.Velocity()
	if vy != -160 {
		t.Fatalf("vy changed on second release: got %d, want -160", vy)
	}
}
