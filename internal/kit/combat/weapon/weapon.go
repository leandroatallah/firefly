package weapon

import (
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/vfx"
	"github.com/boilerplate/ebiten-template/internal/engine/utils"
)

// ProjectileWeapon is a weapon that spawns projectiles.
type ProjectileWeapon struct {
	id               string
	cooldownFrames   int
	currentCooldown  int
	startupFrames    int
	startup          utils.DelayTrigger
	projectileType   string
	projectileSpeed  int
	manager          combat.ProjectileManager
	owner            interface{}
	muzzleEffectType string
	vfxManager       vfx.Manager
	spawnOffsetX16   int
	spawnOffsetY16   int
	stateOffsets     map[int][2]int
	damage           int

	// pending fire params stored during startup
	pendingX16, pendingY16 int
	pendingFaceDir         animation.FacingDirectionEnum
	pendingDirection       body.ShootDirection
	pendingState           int
}

// NewProjectileWeapon creates a new projectile weapon.
func NewProjectileWeapon(id string, cooldownFrames int, projectileType string, projectileSpeed int, manager combat.ProjectileManager, muzzleEffectType string, spawnOffsetX16 int, spawnOffsetY16 int) *ProjectileWeapon {
	return &ProjectileWeapon{
		id:               id,
		cooldownFrames:   cooldownFrames,
		currentCooldown:  0,
		projectileType:   projectileType,
		projectileSpeed:  projectileSpeed,
		manager:          manager,
		owner:            nil,
		muzzleEffectType: muzzleEffectType,
		spawnOffsetX16:   spawnOffsetX16,
		spawnOffsetY16:   spawnOffsetY16,
	}
}

// SetOwner sets the owner of projectiles fired by this weapon.
func (w *ProjectileWeapon) SetOwner(owner interface{}) {
	w.owner = owner
}

// SetStateSpawnOffsets registers per-state spawn offsets. Values are fp16 (x16, y16).
// Passing a nil or empty map clears all per-state overrides.
func (w *ProjectileWeapon) SetStateSpawnOffsets(offsets map[int][2]int) {
	w.stateOffsets = offsets
}

// SetDamage sets the damage dealt by each projectile fired by this weapon.
func (w *ProjectileWeapon) SetDamage(d int) {
	w.damage = d
}

// SetVFXManager sets the visual effects manager for this weapon.
func (w *ProjectileWeapon) SetVFXManager(manager vfx.Manager) {
	w.vfxManager = manager
}

// SetStartupFrames sets the number of frames to wait before spawning the projectile.
func (w *ProjectileWeapon) SetStartupFrames(frames int) { w.startupFrames = frames }

// ID returns the weapon's unique identifier.
func (w *ProjectileWeapon) ID() string {
	return w.id
}

// Fire spawns a projectile if the weapon can fire.
// If startup_frames > 0, the spawn is deferred until the countdown elapses.
func (w *ProjectileWeapon) Fire(x16, y16 int, faceDir animation.FacingDirectionEnum, direction body.ShootDirection, state int) {
	if w.startupFrames > 0 {
		w.startup.Enable(w.startupFrames)
		w.pendingX16 = x16
		w.pendingY16 = y16
		w.pendingFaceDir = faceDir
		w.pendingDirection = direction
		w.pendingState = state
		return
	}
	w.doFire(x16, y16, faceDir, direction, state)
}

func (w *ProjectileWeapon) doFire(x16, y16 int, faceDir animation.FacingDirectionEnum, direction body.ShootDirection, state int) {
	offsetX16 := w.spawnOffsetX16
	offsetY16 := w.spawnOffsetY16

	if w.stateOffsets != nil {
		if override, ok := w.stateOffsets[state]; ok {
			offsetX16 = override[0]
			offsetY16 = override[1]
		}
	}

	if faceDir == animation.FaceDirectionLeft {
		offsetX16 = -offsetX16
	}

	spawnX16 := x16 + offsetX16
	spawnY16 := y16 + offsetY16

	if w.vfxManager != nil && w.muzzleEffectType != "" {
		x := float64(spawnX16) / 16.0
		y := float64(spawnY16) / 16.0
		w.vfxManager.SpawnDirectionalPuff(w.muzzleEffectType, x, y, faceDir == animation.FaceDirectionRight, 1, 0.0)
	}

	vx16, vy16 := w.calculateVelocity(direction, faceDir)
	w.manager.SpawnProjectile(w.projectileType, spawnX16, spawnY16, vx16, vy16, w.damage, w.owner)
	w.currentCooldown = w.cooldownFrames
}

// CanFire returns true if the weapon is ready to fire.
func (w *ProjectileWeapon) CanFire() bool {
	return w.currentCooldown == 0 && !w.startup.IsEnabled()
}

// Update decrements the startup and cooldown timers.
func (w *ProjectileWeapon) Update() {
	w.startup.Update()
	if w.startup.Trigger() {
		w.startup.Reset()
		w.doFire(w.pendingX16, w.pendingY16, w.pendingFaceDir, w.pendingDirection, w.pendingState)
		return
	}
	if w.startup.IsEnabled() {
		return
	}
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
