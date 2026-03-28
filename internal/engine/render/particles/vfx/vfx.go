package vfx

import (
	"bytes"
	"encoding/json"
	"image/color"
	"io/fs"
	"log"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
	"github.com/leandroatallah/firefly/internal/engine/data/schemas"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/render/camera"
	"github.com/leandroatallah/firefly/internal/engine/render/particles"
	"github.com/leandroatallah/firefly/internal/engine/render/vfx/text"
)

type VFXConfig struct {
	Type string `json:"type"`
	schemas.ParticleData
}

// Manager handles all visual effects for the game (particles + floating text).
type Manager struct {
	app.AppContextHolder

	system      *particles.System
	configs     map[string]*particles.Config
	textManager *text.Manager
	pixelConfig *particles.Config
}

func NewManager(fsys fs.FS, path string) *Manager {
	configs := make(map[string]*particles.Config)

	// Load vfx.json
	data, err := fs.ReadFile(fsys, path)
	if err != nil {
		log.Printf("failed to load vfx config: %v", err)
	} else {
		var vfxList []VFXConfig
		if err := json.Unmarshal(data, &vfxList); err != nil {
			log.Printf("failed to parse vfx config: %v", err)
		}

		for _, vfx := range vfxList {
			imgData, err := fs.ReadFile(fsys, vfx.Image)
			if err != nil {
				log.Printf("failed to load particle image %s: %v", vfx.Image, err)
				// Fallback to white pixel
				img := ebiten.NewImage(1, 1)
				img.Fill(color.White)
				configs[vfx.Type] = createConfigFromImage(img, vfx)
				continue
			}
			img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(imgData))
			if err != nil {
				log.Printf("failed to parse particle image %s: %v", vfx.Image, err)
				img = ebiten.NewImage(1, 1)
				img.Fill(color.White)
			}
			configs[vfx.Type] = createConfigFromImage(img, vfx)
		}
	}

	pixelImg := ebiten.NewImage(1, 1)
	pixelImg.Fill(color.White)

	return &Manager{
		system:      particles.NewSystem(),
		configs:     configs,
		textManager: text.NewManager(),
		pixelConfig: &particles.Config{
			Image:       pixelImg,
			FrameWidth:  1,
			FrameHeight: 1,
			FrameCount:  1,
		},
	}
}

// NewManagerFromPath loads a VFX manager from an OS filesystem path.
func NewManagerFromPath(path string) *Manager {
	dir, file := filepath.Split(path)
	if dir == "" {
		dir = "."
	}
	return NewManager(os.DirFS(dir), file)
}

func createConfigFromImage(img *ebiten.Image, vfx VFXConfig) *particles.Config {
	frameCount := 1
	if vfx.FrameWidth > 0 {
		frameCount = img.Bounds().Dx() / vfx.FrameWidth
	}

	return &particles.Config{
		Image:       img,
		FrameWidth:  vfx.FrameWidth,
		FrameHeight: vfx.FrameHeight,
		FrameCount:  frameCount,
		FrameRate:   vfx.FrameRate,
	}
}

// SetDefaultFont sets the default font for floating text effects.
func (m *Manager) SetDefaultFont(f *font.FontText) {
	m.textManager.SetDefaultFont(f)
}

// SpawnPuff creates a puff of particles of the specified type at the given location.
func (m *Manager) SpawnPuff(typeKey string, x, y float64, count int, randRange float64) {
	config, ok := m.configs[typeKey]
	if !ok {
		return
	}

	for i := 0; i < count; i++ {
		p := &particles.Particle{
			X:           x,
			Y:           y,
			VelX:        (rand.Float64() - 0.5) * randRange,
			VelY:        (rand.Float64() - 0.5) * randRange,
			Duration:    config.FrameCount * config.FrameRate,
			MaxDuration: config.FrameCount * config.FrameRate,
			Scale:       1.0,
			ScaleSpeed:  0,
			Config:      config,
		}
		m.system.Add(p)
	}
}

// SpawnJumpPuff creates a jump dust effect at the specified location.
// The randRange parameter controls the randomness of the particle velocities.
func (m *Manager) SpawnJumpPuff(x, y float64, count int) {
	m.SpawnPuff("jump", x, y, count, 0.1)
}

// SpawnLandingPuff creates a landing dust effect at the specified location.
// The randRange parameter controls the randomness of the particle velocities.
func (m *Manager) SpawnLandingPuff(x, y float64, count int) {
	m.SpawnPuff("landing", x, y, count, 0.1)
}

// SpawnFallingRocks spawns falling pixel particles across a specified area.
func (m *Manager) SpawnFallingRocks(x, y, width float64, count int) {
	for i := 0; i < count; i++ {
		rx := x + rand.Float64()*width
		ry := y + (rand.Float64()-0.5)*10

		// Random gray color for dust/rocks (mostly opaque)
		gray := uint8(150 + rand.Intn(80))
		c := color.RGBA{gray, gray, gray, 255}

		scale := 1.0
		if rand.Float64() > 0.7 {
			scale = 2.0
		}

		p := &particles.Particle{
			X:           rx,
			Y:           ry,
			VelX:        (rand.Float64() - 0.5) * 0.2,
			VelY:        rand.Float64() * 0.5,
			AccY:        0.03 + rand.Float64()*0.05, // Lighter gravity for dust
			Duration:    90 + rand.Intn(60),
			MaxDuration: 150,
			Scale:       scale,
			Config:      m.pixelConfig,
		}
		p.ColorScale.ScaleWithColor(c)
		m.system.Add(p)
	}
}

// SpawnDeathExplosion spawns an explosion effect at the specified location.
// Particles explode outward and upward, then fall with gravity.
func (m *Manager) SpawnDeathExplosion(x, y float64, count int) {
	for i := 0; i < count; i++ {
		// Spread particles around the death location
		rx := x + (rand.Float64()-0.5)*24
		ry := y + (rand.Float64()-0.5)*24

		// Red color for 1-bit aesthetic
		c := color.RGBA{255, 0, 0, 255}

		// Explosion velocity: outward and upward
		velX := (rand.Float64() - 0.5) * 4.0
		velY := (rand.Float64() - 0.8) * 5.0 // Bias upward

		scale := 1.0 + rand.Float64()*1.5

		p := &particles.Particle{
			X:           rx,
			Y:           ry,
			VelX:        velX,
			VelY:        velY,
			AccY:        0.15, // Gravity pulls particles down
			Duration:    30 + rand.Intn(20),
			MaxDuration: 60,
			Scale:       scale,
			Config:      m.pixelConfig,
		}
		p.ColorScale.ScaleWithColor(c)
		m.system.Add(p)
	}
}

// AddParticle adds a custom particle to the VFX system.
func (m *Manager) AddParticle(p *particles.Particle) {
	m.system.Add(p)
}

// PixelConfig returns the default 1x1 pixel particle configuration.
func (m *Manager) PixelConfig() *particles.Config {
	return m.pixelConfig
}

// SpawnFloatingText spawns floating text at the specified location.
func (m *Manager) SpawnFloatingText(msg string, x, y float64, duration int) {
	ft := text.NewFloatingText(msg, x, y, duration)
	m.textManager.Add(ft)
}

// SpawnFloatingTextAbove spawns floating text above an actor.
func (m *Manager) SpawnFloatingTextAbove(actor actors.ActorEntity, msg string, duration int) {
	pos := actor.Position()
	x := float64(pos.Min.X + pos.Dx()/2) // Center horizontally
	y := float64(pos.Min.Y)              // Top of actor
	m.SpawnFloatingText(msg, x, y, duration)
}

// AddTrauma adds trauma to the camera to cause a screen shake and spawns falling rocks.
func (m *Manager) AddTrauma(cam *camera.Controller, amount float64) {
	if cam != nil {
		cam.AddTrauma(amount)

		// Spawn rocks proportional to trauma
		k := cam.Kamera()
		vw := cam.Width() / k.ZoomFactor
		// 1 to 5 rocks depending on intensity
		count := 1 + int(amount*4.0)
		m.SpawnFallingRocks(k.X, k.Y-10, vw, count)
	}
}

func (m *Manager) Update() {
	m.system.Update()
	m.textManager.Update()
}

func (m *Manager) Draw(screen *ebiten.Image, cam *camera.Controller) {
	m.system.Draw(screen, cam)
	m.textManager.Draw(screen, cam)
}

// TriggerScreenFlash triggers a screen flash effect on the current scene if it supports it.
func (m *Manager) TriggerScreenFlash() {
	if m.AppContext().SceneManager == nil {
		return
	}
	// Type assertion to access PhasesScene-specific method
	type screenFlasher interface {
		TriggerScreenFlash()
	}
	if scene, ok := m.AppContext().SceneManager.CurrentScene().(screenFlasher); ok {
		scene.TriggerScreenFlash()
	}
}
