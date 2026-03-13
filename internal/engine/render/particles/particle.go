package particles

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/render/camera"
)

// Config defines the properties of a particle type.
type Config struct {
	Image       *ebiten.Image
	FrameWidth  int
	FrameHeight int
	FrameCount  int
	FrameRate   int // Ticks per frame
}

// Particle represents an active particle instance.
type Particle struct {
	X, Y       float64
	VelX, VelY float64
	AccX, AccY float64

	Duration    int // Current remaining ticks
	MaxDuration int // Initial total duration

	Scale      float64
	ScaleSpeed float64

	ColorScale ebiten.ColorScale

	// Animation state
	Frame      int
	FrameTimer int

	Config *Config
}

// Update advances the particle state.
func (p *Particle) Update() {
	p.VelX += p.AccX
	p.VelY += p.AccY
	p.X += p.VelX
	p.Y += p.VelY
	p.Duration--
	p.Scale += p.ScaleSpeed

	if p.Config.FrameCount > 1 {
		p.FrameTimer++
		if p.FrameTimer >= p.Config.FrameRate {
			p.FrameTimer = 0
			p.Frame++
			if p.Frame >= p.Config.FrameCount {
				p.Frame = 0
			}
		}
	}
}

func (p *Particle) IsExpired() bool {
	return p.Duration <= 0
}

func (p *Particle) Draw(screen *ebiten.Image, cam *camera.Controller) {
	if p.Config.Image == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(p.Scale, p.Scale)
	op.GeoM.Translate(p.X, p.Y)

	// Center logic
	w, h := p.Config.FrameWidth, p.Config.FrameHeight
	if w == 0 {
		w = p.Config.Image.Bounds().Dx()
	}
	if h == 0 {
		h = p.Config.Image.Bounds().Dy()
	}

	// Center the sprite
	op.GeoM.Translate(-float64(w)*p.Scale/2, -float64(h)*p.Scale)

	op.ColorScale = p.ColorScale

	var subImg *ebiten.Image
	if p.Config.FrameCount > 1 {
		sx := p.Frame * w
		// Clamp
		if sx >= p.Config.Image.Bounds().Dx() {
			sx = p.Config.Image.Bounds().Dx() - w
		}
		subImg = p.Config.Image.SubImage(image.Rect(sx, 0, sx+w, h)).(*ebiten.Image)
	} else {
		subImg = p.Config.Image
	}

	cam.Draw(subImg, op, screen)
}

// System manages a collection of particles.
type System struct {
	particles []*Particle
}

func NewSystem() *System {
	return &System{
		particles: make([]*Particle, 0, 100),
	}
}

func (s *System) Add(p *Particle) {
	s.particles = append(s.particles, p)
}

func (s *System) Update() {
	activeParticles := s.particles[:0]
	for _, p := range s.particles {
		p.Update()
		if !p.IsExpired() {
			activeParticles = append(activeParticles, p)
		}
	}
	s.particles = activeParticles
}

func (s *System) Draw(screen *ebiten.Image, cam *camera.Controller) {
	for _, p := range s.particles {
		p.Draw(screen, cam)
	}
}
