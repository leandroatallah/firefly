package mocks

import (
	"github.com/boilerplate/ebiten-template/internal/engine/assets/font"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/render/camera"
	"github.com/boilerplate/ebiten-template/internal/engine/render/particles"
	"github.com/hajimehoshi/ebiten/v2"
)

// MockVFXManager is a shared mock implementation of vfx.Manager.
// Configure SpawnPuffFunc (and other Func fields) to capture or assert calls in tests.
type MockVFXManager struct {
	SetAppContextFunc          func(appContext any)
	UpdateFunc                 func()
	DrawFunc                   func(screen *ebiten.Image, cam *camera.Controller)
	AddParticleFunc            func(p *particles.Particle)
	AddTraumaFunc              func(cam *camera.Controller, amount float64)
	PixelConfigFunc            func() *particles.Config
	SetDefaultFontFunc         func(f *font.FontText)
	SpawnDeathExplosionFunc    func(x float64, y float64, count int)
	SpawnFallingRocksFunc      func(x float64, y float64, width float64, count int)
	SpawnFloatingTextFunc      func(msg string, x float64, y float64, duration int)
	SpawnFloatingTextAboveFunc func(target body.Body, msg string, duration int)
	SpawnJumpPuffFunc          func(x float64, y float64, count int)
	SpawnLandingPuffFunc       func(x float64, y float64, count int)
	SpawnPuffFunc              func(typeKey string, x float64, y float64, count int, randRange float64)
	SpawnDirectionalPuffFunc   func(typeKey string, x float64, y float64, faceRight bool, count int, randRange float64)
	TriggerScreenFlashFunc     func()
}

func (m *MockVFXManager) SetAppContext(appContext any) {
	if m.SetAppContextFunc != nil {
		m.SetAppContextFunc(appContext)
	}
}

func (m *MockVFXManager) Update() {
	if m.UpdateFunc != nil {
		m.UpdateFunc()
	}
}

func (m *MockVFXManager) Draw(screen *ebiten.Image, cam *camera.Controller) {
	if m.DrawFunc != nil {
		m.DrawFunc(screen, cam)
	}
}

func (m *MockVFXManager) AddParticle(p *particles.Particle) {
	if m.AddParticleFunc != nil {
		m.AddParticleFunc(p)
	}
}

func (m *MockVFXManager) AddTrauma(cam *camera.Controller, amount float64) {
	if m.AddTraumaFunc != nil {
		m.AddTraumaFunc(cam, amount)
	}
}

func (m *MockVFXManager) PixelConfig() *particles.Config {
	if m.PixelConfigFunc != nil {
		return m.PixelConfigFunc()
	}
	return nil
}

func (m *MockVFXManager) SetDefaultFont(f *font.FontText) {
	if m.SetDefaultFontFunc != nil {
		m.SetDefaultFontFunc(f)
	}
}

func (m *MockVFXManager) SpawnDeathExplosion(x float64, y float64, count int) {
	if m.SpawnDeathExplosionFunc != nil {
		m.SpawnDeathExplosionFunc(x, y, count)
	}
}

func (m *MockVFXManager) SpawnFallingRocks(x float64, y float64, width float64, count int) {
	if m.SpawnFallingRocksFunc != nil {
		m.SpawnFallingRocksFunc(x, y, width, count)
	}
}

func (m *MockVFXManager) SpawnFloatingText(msg string, x float64, y float64, duration int) {
	if m.SpawnFloatingTextFunc != nil {
		m.SpawnFloatingTextFunc(msg, x, y, duration)
	}
}

func (m *MockVFXManager) SpawnFloatingTextAbove(target body.Body, msg string, duration int) {
	if m.SpawnFloatingTextAboveFunc != nil {
		m.SpawnFloatingTextAboveFunc(target, msg, duration)
	}
}

func (m *MockVFXManager) SpawnJumpPuff(x float64, y float64, count int) {
	if m.SpawnJumpPuffFunc != nil {
		m.SpawnJumpPuffFunc(x, y, count)
	}
}

func (m *MockVFXManager) SpawnLandingPuff(x float64, y float64, count int) {
	if m.SpawnLandingPuffFunc != nil {
		m.SpawnLandingPuffFunc(x, y, count)
	}
}

func (m *MockVFXManager) SpawnPuff(typeKey string, x float64, y float64, count int, randRange float64) {
	if m.SpawnPuffFunc != nil {
		m.SpawnPuffFunc(typeKey, x, y, count, randRange)
	}
}

func (m *MockVFXManager) SpawnDirectionalPuff(typeKey string, x float64, y float64, faceRight bool, count int, randRange float64) {
	if m.SpawnDirectionalPuffFunc != nil {
		m.SpawnDirectionalPuffFunc(typeKey, x, y, faceRight, count, randRange)
	}
}

func (m *MockVFXManager) TriggerScreenFlash() {
	if m.TriggerScreenFlashFunc != nil {
		m.TriggerScreenFlashFunc()
	}
}
