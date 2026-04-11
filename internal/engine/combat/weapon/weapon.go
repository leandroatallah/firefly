package weapon

import (
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/vfx"
)

// ProjectileWeapon is a weapon that spawns projectiles.
type ProjectileWeapon struct {
	id               string
	cooldownFrames   int
	currentCooldown  int
	projectileType   string
	projectileSpeed  int
	manager          combat.ProjectileManager
	owner            interface{}
	muzzleEffectType string
	vfxManager       vfx.Manager
}

// NewProjectileWeapon creates a new projectile weapon.
func NewProjectileWeapon(id string, cooldownFrames int, projectileType string, projectileSpeed int, manager combat.ProjectileManager, muzzleEffectType string) *ProjectileWeapon {
	return &ProjectileWeapon{
		id:               id,
		cooldownFrames:   cooldownFrames,
		currentCooldown:  0,
		projectileType:   projectileType,
		projectileSpeed:  projectileSpeed,
		manager:          manager,
		owner:            nil,
		muzzleEffectType: muzzleEffectType,
	}
}

// SetOwner sets the owner of projectiles fired by this weapon.
func (w *ProjectileWeapon) SetOwner(owner interface{}) {
	w.owner = owner
}

// SetVFXManager sets the visual effects manager for this weapon.
func (w *ProjectileWeapon) SetVFXManager(manager vfx.Manager) {
	w.vfxManager = manager
}

// ID returns the weapon's unique identifier.
func (w *ProjectileWeapon) ID() string {
	return w.id
}

// Fire spawns a projectile if the weapon can fire.
func (w *ProjectileWeapon) Fire(x16, y16 int, faceDir animation.FacingDirectionEnum, direction body.ShootDirection) {
	// Spawn muzzle flash VFX
	if w.vfxManager != nil && w.muzzleEffectType != "" {
		x := float64(x16) / 16.0
		y := float64(y16) / 16.0
		w.vfxManager.SpawnPuff(w.muzzleEffectType, x, y, 1, 0.0)
	}

	vx16, vy16 := w.calculateVelocity(direction, faceDir)
	w.manager.SpawnProjectile(w.projectileType, x16, y16, vx16, vy16, w.owner)
	w.currentCooldown = w.cooldownFrames
}

// CanFire returns true if the weapon is ready to fire.
func (w *ProjectileWeapon) CanFire() bool {
	return w.currentCooldown == 0
}

// Update decrements the cooldown timer.
func (w *ProjectileWeapon) Update() {
	if w.currentCooldown > 0 {
		w.currentCooldown--
	}
}

// Cooldown returns the current cooldown value.
func (w *ProjectileWeapon) Cooldown() int {
	return w.currentCooldown
}

// SetCooldown sets the cooldown to the given value.
func (w *ProjectileWeapon) SetCooldown(frames int) {
	w.currentCooldown = frames
}

func (w *ProjectileWeapon) calculateVelocity(direction body.ShootDirection, faceDir animation.FacingDirectionEnum) (vx16, vy16 int) {
	speed := w.projectileSpeed
	sign := 1
	if faceDir == animation.FaceDirectionLeft {
		sign = -1
	}

	switch direction {
	case body.ShootDirectionStraight:
		return sign * speed, 0
	case body.ShootDirectionUp:
		return 0, -speed
	case body.ShootDirectionDown:
		return 0, speed
	case body.ShootDirectionDiagonalUpForward:
		return sign * speed * 707 / 1000, -speed * 707 / 1000
	case body.ShootDirectionDiagonalDownForward:
		return sign * speed * 707 / 1000, speed * 707 / 1000
	}
	return sign * speed, 0
}
