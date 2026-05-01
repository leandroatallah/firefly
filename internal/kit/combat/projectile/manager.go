// internal/engine/combat/projectile/manager.go
package projectile

import (
	"fmt"
	"image/color"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	contractsvfx "github.com/boilerplate/ebiten-template/internal/engine/contracts/vfx"
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
	enginecombat "github.com/boilerplate/ebiten-template/internal/kit/combat"
	"github.com/hajimehoshi/ebiten/v2"
)

// getDefaultBulletImg creates and returns a white 2x1 image.
// This is called per Draw to avoid global state.
func getDefaultBulletImg() *ebiten.Image {
	img := ebiten.NewImage(2, 1)
	img.Fill(color.White)
	return img
}

// Manager handles the lifecycle of all projectiles in the game.
type Manager struct {
	projectiles   []*projectile
	space         body.BodiesSpace
	counter       int
	vfxManager    contractsvfx.Manager
	impactEffect  string
	despawnEffect string
}

// NewManager creates a new projectile manager.
func NewManager(space body.BodiesSpace) *Manager {
	return &Manager{
		projectiles:   make([]*projectile, 0),
		space:         space,
		impactEffect:  "bullet_impact",
		despawnEffect: "bullet_despawn",
	}
}

// SetVFXManager sets the VFX manager to be used for projectile effects.
func (m *Manager) SetVFXManager(v contractsvfx.Manager) {
	m.vfxManager = v
}

// Spawn creates a new projectile and registers it in the physics space.
func (m *Manager) Spawn(cfg interface{}, x16, y16, vx16, vy16 int, owner interface{}) {
	config, ok := cfg.(ProjectileConfig)
	if !ok {
		return
	}

	m.counter++
	id := fmt.Sprintf("bullet_%d", m.counter)

	baseBody := bodyphysics.NewBody(bodyphysics.NewRect(0, 0, config.Width, config.Height))
	baseBody.SetID(id)
	baseBody.SetPosition16(x16, y16)

	movableBody := bodyphysics.NewMovableBody(baseBody)
	collidableBody := bodyphysics.NewCollidableBody(baseBody)
	collidableBody.SetOwner(owner)

	b := bodyphysics.NewCollidableBodyFromRect(baseBody.GetShape())
	x, y := baseBody.GetPositionMin()
	b.SetPosition(x, y)
	b.SetID(fmt.Sprintf("%v_COLLISION_0", id))
	collidableBody.AddCollision(b)

	lifetime := max(config.LifetimeFrames, 0)

	impactEffect := config.ImpactEffect
	if impactEffect == "" && config.LifetimeFrames == 0 {
		impactEffect = m.impactEffect
	}
	despawnEffect := config.DespawnEffect
	if despawnEffect == "" && config.LifetimeFrames == 0 {
		despawnEffect = m.despawnEffect
	}

	faction := config.Faction
	if faction == enginecombat.FactionNeutral && owner != nil {
		if f, ok := owner.(factioned); ok {
			faction = f.Faction()
		}
	}

	// Wrap the collidable body so the Projectile trait is discoverable on the body itself.
	wrappedBody := &projectileBody{
		Collidable:    collidableBody,
		interceptable: config.Interceptable,
	}

	p := &projectile{
		movable:         movableBody,
		body:            wrappedBody,
		space:           m.space,
		speedX16:        vx16,
		speedY16:        vy16,
		vfxManager:      m.vfxManager,
		impactEffect:    impactEffect,
		despawnEffect:   despawnEffect,
		lifetimeFrames:  lifetime,
		currentLifetime: lifetime,
		damage:          config.Damage,
		faction:         faction,
		interceptable:   config.Interceptable,
	}

	// Register collision callbacks
	wrappedBody.SetTouchable(p)

	m.projectiles = append(m.projectiles, p)
	m.space.AddBody(wrappedBody)
}

// SpawnProjectile implements the ProjectileManager interface.
func (m *Manager) SpawnProjectile(projectileType string, x16, y16, vx16, vy16, damage int, owner interface{}) {
	cfg := ProjectileConfig{Width: 2, Height: 1, Damage: damage}
	m.Spawn(cfg, x16, y16, vx16, vy16, owner)
}

// Update advances all active projectiles and removes those that are despawned.
func (m *Manager) Update() {
	// First, update all projectiles
	for _, p := range m.projectiles {
		p.Update()
	}

	// Process any queued removals
	m.space.ProcessRemovals()

	// Then, remove projectiles whose bodies are no longer in the space
	activeProjectiles := m.projectiles[:0]
	for _, p := range m.projectiles {
		// Check if the body is still managed by the space
		if m.isBodyInSpace(p.body) {
			activeProjectiles = append(activeProjectiles, p)
		}
	}
	m.projectiles = activeProjectiles
}

// isBodyInSpace checks if the given body is still present in the BodiesSpace.
func (m *Manager) isBodyInSpace(b body.Collidable) bool {
	return m.space.Find(b.ID()) != nil
}

// Draw renders all active projectiles to the screen.
// Note: As specified in the interface, this does not take a camera.
// For world-space rendering, the caller is expected to provide a translated screen
// or the interface should be updated to include a camera.
func (m *Manager) Draw(screen *ebiten.Image) {
	for _, p := range m.projectiles {
		opts := &ebiten.DrawImageOptions{}
		x, y := p.body.GetPositionMin()
		opts.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(getDefaultBulletImg(), opts)
	}
}

// DrawWithOffset renders all active projectiles with camera offset applied.
func (m *Manager) DrawWithOffset(screen *ebiten.Image, camX, camY float64) {
	for _, p := range m.projectiles {
		opts := &ebiten.DrawImageOptions{}
		x, y := p.body.GetPositionMin()
		opts.GeoM.Translate(float64(x)-camX, float64(y)-camY)
		screen.DrawImage(getDefaultBulletImg(), opts)
	}
}

// DrawCollisionBoxesWithOffset renders each active projectile's collision box
// using the given camera-space draw helper. The helper is invoked once per
// active projectile body.
func (m *Manager) DrawCollisionBoxesWithOffset(draw func(b body.Collidable)) {
	if draw == nil {
		return
	}
	for _, p := range m.projectiles {
		draw(p.body)
	}
}

// Clear removes all projectiles and their bodies from the physics space.
func (m *Manager) Clear() {
	for _, p := range m.projectiles {
		m.space.RemoveBody(p.body)
	}
	m.projectiles = m.projectiles[:0]
}
