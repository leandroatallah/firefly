package particles

import (
	"image/color"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/render/camera"
	"github.com/hajimehoshi/ebiten/v2"
)

func TestParticleUpdate(t *testing.T) {
	config := &Config{
		FrameCount: 2,
		FrameRate:  5, // 5 ticks per animation frame
	}
	p := &Particle{
		X:          10,
		Y:          20,
		VelX:       1,
		VelY:       2,
		Duration:   10,
		Scale:      1.0,
		ScaleSpeed: 0.1,
		Config:     config,
	}

	p.Update()

	if p.X != 11 || p.Y != 22 {
		t.Errorf("expected position (11, 22), got (%f, %f)", p.X, p.Y)
	}
	if p.Duration != 9 {
		t.Errorf("expected duration 9, got %d", p.Duration)
	}
	if p.Scale != 1.1 {
		t.Errorf("expected scale 1.1, got %f", p.Scale)
	}

	// FrameRate is 5, so we need 5 updates to advance one animation frame
	// 1st update above made FrameTimer=1
	// After 4 more updates, FrameTimer=5 -> Frame=1, FrameTimer=0
	for i := 0; i < 4; i++ {
		p.Update()
	}
	if p.Frame != 1 {
		t.Errorf("expected frame 1 after 5 updates, got %d", p.Frame)
	}

	// 5 more updates to trigger frame reset (FrameCount=2, so frame wraps to 0)
	for i := 0; i < 5; i++ {
		p.Update()
	}
	if p.Frame != 0 {
		t.Errorf("expected frame reset to 0, got %d", p.Frame)
	}
}

func TestParticleUpdate_EdgeCases(t *testing.T) {
	tests := []struct {
		name           string
		initialFrame   int
		initialTimer   int
		frameCount     int
		frameRate      int
		wantFrame      int
		wantFrameTimer int
	}{
		{
			name:         "single frame no animation",
			initialFrame: 0, initialTimer: 0,
			frameCount: 1, frameRate: 1,
			wantFrame: 0, wantFrameTimer: 0,
		},
		{
			name:         "frame advance at rate boundary",
			initialFrame: 0, initialTimer: 4,
			frameCount: 3, frameRate: 5,
			wantFrame: 1, wantFrameTimer: 0,
		},
		{
			name:         "frame wrap around",
			initialFrame: 2, initialTimer: 4,
			frameCount: 3, frameRate: 5,
			wantFrame: 0, wantFrameTimer: 0,
		},
		{
			name:         "zero frame rate advances immediately",
			initialFrame: 0, initialTimer: 0,
			frameCount: 3, frameRate: 0,
			wantFrame: 1, wantFrameTimer: 0, // 0 >= 0 triggers advance
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Particle{
				Frame:      tt.initialFrame,
				FrameTimer: tt.initialTimer,
				Config: &Config{
					FrameCount: tt.frameCount,
					FrameRate:  tt.frameRate,
				},
			}

			p.Update()

			if p.Frame != tt.wantFrame {
				t.Errorf("Frame = %d, want %d", p.Frame, tt.wantFrame)
			}
			if p.FrameTimer != tt.wantFrameTimer {
				t.Errorf("FrameTimer = %d, want %d", p.FrameTimer, tt.wantFrameTimer)
			}
		})
	}
}

func TestParticleUpdate_VelocityAndScale(t *testing.T) {
	tests := []struct {
		name         string
		x, y         float64
		velX, velY   float64
		scale        float64
		scaleSpeed   float64
		wantX, wantY float64
		wantScale    float64
	}{
		{
			name: "positive velocity",
			x:    10, y: 20, velX: 5, velY: -3,
			scale: 1.0, scaleSpeed: 0.1,
			wantX: 15, wantY: 17, wantScale: 1.1,
		},
		{
			name: "zero velocity",
			x:    100, y: 200, velX: 0, velY: 0,
			scale: 2.0, scaleSpeed: 0,
			wantX: 100, wantY: 200, wantScale: 2.0,
		},
		{
			name: "negative scale speed (shrinking)",
			x:    0, y: 0, velX: 0, velY: 0,
			scale: 1.0, scaleSpeed: -0.05,
			wantX: 0, wantY: 0, wantScale: 0.95,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Particle{
				X: tt.x, Y: tt.y,
				VelX: tt.velX, VelY: tt.velY,
				Scale: tt.scale, ScaleSpeed: tt.scaleSpeed,
				Duration: 10,
				Config:   &Config{FrameCount: 1},
			}

			p.Update()

			if p.X != tt.wantX || p.Y != tt.wantY {
				t.Errorf("Position = (%f, %f), want (%f, %f)", p.X, p.Y, tt.wantX, tt.wantY)
			}
			if p.Scale != tt.wantScale {
				t.Errorf("Scale = %f, want %f", p.Scale, tt.wantScale)
			}
		})
	}
}

func TestParticleIsExpired(t *testing.T) {
	tests := []struct {
		name     string
		duration int
		wantExp  bool
	}{
		{"positive duration", 1, false},
		{"zero duration", 0, true},
		{"negative duration", -5, true},
		{"large duration", 1000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Particle{Duration: tt.duration}
			if got := p.IsExpired(); got != tt.wantExp {
				t.Errorf("IsExpired() = %v, want %v", got, tt.wantExp)
			}
		})
	}
}

func TestSystem(t *testing.T) {
	sys := NewSystem()
	p1 := &Particle{Duration: 3, Config: &Config{}}
	p2 := &Particle{Duration: 2, Config: &Config{}}

	sys.Add(p1)
	sys.Add(p2)

	sys.Update()
	// p1: Dur=2, p2: Dur=1. Both still active.
	if len(sys.particles) != 2 {
		t.Errorf("expected 2 particles, got %d", len(sys.particles))
	}

	sys.Update()
	// p1: Dur=1, p2: Dur=0. p2 should be expired and removed.
	if len(sys.particles) != 1 {
		t.Errorf("expected 1 particle (p1), got %d", len(sys.particles))
	}

	sys.Update()
	// p1: Dur=0. Removed.
	if len(sys.particles) != 0 {
		t.Errorf("expected 0 particles, got %d", len(sys.particles))
	}
}

func TestSystem_AddMultiple(t *testing.T) {
	sys := NewSystem()

	// Add 10 particles with varying durations
	for i := 0; i < 10; i++ {
		sys.Add(&Particle{Duration: i + 1, Config: &Config{}})
	}

	if len(sys.particles) != 10 {
		t.Errorf("expected 10 particles, got %d", len(sys.particles))
	}

	// After 5 updates, particles with duration <= 5 should be expired
	sys.Update() // All durations decrease by 1
	sys.Update()
	sys.Update()
	sys.Update()
	sys.Update()

	// Particles with original duration 1-5 are now expired (0 or negative)
	// Particles with original duration 6-10 remain (now 1-5)
	if len(sys.particles) != 5 {
		t.Errorf("expected 5 remaining particles, got %d", len(sys.particles))
	}
}

func TestSystem_Update(t *testing.T) {
	sys := NewSystem()

	// Add particles with different durations
	p1 := &Particle{Duration: 5, VelX: 1, Config: &Config{}}
	p2 := &Particle{Duration: 3, VelX: 2, Config: &Config{}}
	p3 := &Particle{Duration: 1, VelX: 3, Config: &Config{}}

	sys.Add(p1)
	sys.Add(p2)
	sys.Add(p3)

	// Initial positions
	p1Start := p1.X
	p2Start := p2.X
	p3Start := p3.X

	// After 1 update, all should move
	sys.Update()

	if p1.X <= p1Start {
		t.Error("p1 should have moved")
	}
	if p2.X <= p2Start {
		t.Error("p2 should have moved")
	}
	if p3.X <= p3Start {
		t.Error("p3 should have moved")
	}
}

func TestSystem_Draw(t *testing.T) {
	sys := NewSystem()
	img := ebiten.NewImage(32, 32)
	img.Fill(color.White)

	// Add a particle
	p := &Particle{
		X: 50, Y: 50,
		Duration: 10,
		Scale:    1.0,
		Config: &Config{
			Image:       img,
			FrameWidth:  32,
			FrameHeight: 32,
			FrameCount:  1,
		},
	}
	sys.Add(p)

	screen := ebiten.NewImage(100, 100)
	cam := camera.NewController(0, 0)

	// Smoke test - should not panic
	sys.Draw(screen, cam)
}

func TestSystem_DrawEmpty(t *testing.T) {
	sys := NewSystem()
	screen := ebiten.NewImage(100, 100)
	cam := camera.NewController(0, 0)

	// Should not panic with empty system
	sys.Draw(screen, cam)
}

func TestParticleDrawSmoke(t *testing.T) {
	img := ebiten.NewImage(32, 32)
	img.Fill(color.White)
	config := &Config{
		Image:       img,
		FrameWidth:  16,
		FrameHeight: 16,
		FrameCount:  2,
	}
	p := &Particle{
		X:      10,
		Y:      10,
		Scale:  1.0,
		Config: config,
		Frame:  1,
	}

	screen := ebiten.NewImage(100, 100)
	cam := camera.NewController(0, 0)

	// Smoke test
	p.Draw(screen, cam)

	// Test with nil image
	p.Config.Image = nil
	p.Draw(screen, cam)
}

func TestParticleDraw_WithColorScale(t *testing.T) {
	img := ebiten.NewImage(16, 16)
	img.Fill(color.White)
	config := &Config{
		Image:       img,
		FrameWidth:  16,
		FrameHeight: 16,
		FrameCount:  1,
	}
	p := &Particle{
		X:      50,
		Y:      50,
		Scale:  2.0,
		Config: config,
	}
	// Set a color scale (tint)
	p.ColorScale.Scale(1, 0.5, 0, 1) // Yellow tint

	screen := ebiten.NewImage(100, 100)
	cam := camera.NewController(0, 0)

	// Smoke test - should not panic with color scale
	p.Draw(screen, cam)
}

func TestParticleDraw_AnimationFrames(t *testing.T) {
	// Create a sprite sheet with 3 frames (48x16 image with 16x16 frames)
	img := ebiten.NewImage(48, 16)
	img.Fill(color.White)
	config := &Config{
		Image:       img,
		FrameWidth:  16,
		FrameHeight: 16,
		FrameCount:  3,
	}

	tests := []struct {
		name  string
		frame int
	}{
		{"frame 0", 0},
		{"frame 1", 1},
		{"frame 2", 2},
		{"frame out of bounds", 10}, // Should clamp to last frame
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Particle{
				X:      50,
				Y:      50,
				Scale:  1.0,
				Config: config,
				Frame:  tt.frame,
			}

			screen := ebiten.NewImage(100, 100)
			cam := camera.NewController(0, 0)

			// Should not panic for any frame value
			p.Draw(screen, cam)
		})
	}
}

func TestParticleDraw_ZeroFrameSize(t *testing.T) {
	img := ebiten.NewImage(32, 32)
	img.Fill(color.White)
	config := &Config{
		Image:       img,
		FrameWidth:  0, // Should use image bounds
		FrameHeight: 0,
		FrameCount:  1,
	}
	p := &Particle{
		X:      50,
		Y:      50,
		Scale:  1.0,
		Config: config,
	}

	screen := ebiten.NewImage(100, 100)
	cam := camera.NewController(0, 0)

	// Should not panic with zero frame size
	p.Draw(screen, cam)
}

func TestNewSystem(t *testing.T) {
	sys := NewSystem()
	if sys == nil {
		t.Fatal("NewSystem() returned nil")
	}
	if sys.particles == nil {
		t.Error("NewSystem() did not initialize particles slice")
	}
	if len(sys.particles) != 0 {
		t.Errorf("NewSystem() created system with %d particles, want 0", len(sys.particles))
	}
}
