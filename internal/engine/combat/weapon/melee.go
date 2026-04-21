package weapon

import (
	"image"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
)

// MeleeWeapon is a close-range swing weapon that activates a hitbox during a
// configurable active-frame window within its animation.
type MeleeWeapon struct {
	id              string
	damage          int
	cooldownFrames  int
	currentCooldown int
	activeFrames    [2]int // [first, last] inclusive
	hitboxW16       int
	hitboxH16       int
	hitboxOffsetX16 int // offset to hitbox center in fp16 (mirrored when facing left)
	hitboxOffsetY16 int

	owner interface{}

	// swing state
	swinging     bool
	swingFrame   int
	hitThisSwing map[combat.Damageable]struct{}

	// stored when Fire() is called, used by ApplyHitbox
	faceDir              animation.FacingDirectionEnum
	originX16, originY16 int
}

// NewMeleeWeapon constructs a MeleeWeapon.
// hitboxW16/H16 and offsetX16/Y16 are in fp16 units (offset points to center of hitbox).
func NewMeleeWeapon(id string, damage, cooldownFrames int, activeFrames [2]int, hitboxW16, hitboxH16, hitboxOffsetX16, hitboxOffsetY16 int) *MeleeWeapon {
	return &MeleeWeapon{
		id:              id,
		damage:          damage,
		cooldownFrames:  cooldownFrames,
		activeFrames:    activeFrames,
		hitboxW16:       hitboxW16,
		hitboxH16:       hitboxH16,
		hitboxOffsetX16: hitboxOffsetX16,
		hitboxOffsetY16: hitboxOffsetY16,
		hitThisSwing:    make(map[combat.Damageable]struct{}),
	}
}

// ID returns the weapon identifier.
func (w *MeleeWeapon) ID() string { return w.id }

// Damage returns the damage value per hit.
func (w *MeleeWeapon) Damage() int { return w.damage }

// ActiveFrames returns the [first, last] active-frame window.
func (w *MeleeWeapon) ActiveFrames() [2]int { return w.activeFrames }

// CanFire returns true when there is no active cooldown.
func (w *MeleeWeapon) CanFire() bool { return w.currentCooldown == 0 }

// Cooldown returns the remaining cooldown frames.
func (w *MeleeWeapon) Cooldown() int { return w.currentCooldown }

// SetCooldown sets the cooldown to the given value.
func (w *MeleeWeapon) SetCooldown(frames int) { w.currentCooldown = frames }

// SetOwner sets the owner reference used for faction checks.
func (w *MeleeWeapon) SetOwner(owner interface{}) { w.owner = owner }

// Fire begins a melee swing, resetting swing state and starting the cooldown.
func (w *MeleeWeapon) Fire(x16, y16 int, faceDir animation.FacingDirectionEnum, _ body.ShootDirection, _ int) {
	if !w.CanFire() {
		return
	}
	w.originX16 = x16
	w.originY16 = y16
	w.faceDir = faceDir
	w.swinging = true
	w.swingFrame = 0
	w.hitThisSwing = make(map[combat.Damageable]struct{})
	w.currentCooldown = w.cooldownFrames
}

// Update advances the swing frame and decrements the cooldown.
func (w *MeleeWeapon) Update() {
	if w.currentCooldown > 0 {
		w.currentCooldown--
	}
	if w.swinging {
		w.swingFrame++
		if w.swingFrame > w.activeFrames[1] {
			w.swinging = false
		}
	}
}

// IsHitboxActive returns true when the current swing frame is within the
// active-frame window.
func (w *MeleeWeapon) IsHitboxActive() bool {
	if !w.swinging {
		return false
	}
	return w.swingFrame >= w.activeFrames[0] && w.swingFrame <= w.activeFrames[1]
}

// ApplyHitbox queries space for targets in the hitbox rect and applies
// faction-gated, single-hit-per-swing damage.
func (w *MeleeWeapon) ApplyHitbox(space body.BodiesSpace) {
	if !w.IsHitboxActive() {
		return
	}

	rect := w.hitboxRect()
	hits := space.Query(rect)

	for _, b := range hits {
		// Self-damage guard.
		if w.owner != nil && b == w.owner {
			continue
		}

		// Same-faction gate (only if both are factioned).
		if ownerFactioned, ok := w.owner.(combat.Factioned); ok {
			if targetFactioned, ok := b.(combat.Factioned); ok {
				if targetFactioned.Faction() == ownerFactioned.Faction() {
					continue
				}
			}
		}

		target, ok := b.(combat.Damageable)
		if !ok {
			continue
		}

		if _, already := w.hitThisSwing[target]; already {
			continue
		}
		w.hitThisSwing[target] = struct{}{}
		target.TakeDamage(w.damage)
	}
}

// hitboxRect computes the query rectangle in pixel space.
func (w *MeleeWeapon) hitboxRect() image.Rectangle {
	halfW16 := w.hitboxW16 / 2

	var centerX16 int
	if w.faceDir == animation.FaceDirectionLeft {
		centerX16 = w.originX16 - w.hitboxOffsetX16
	} else {
		centerX16 = w.originX16 + w.hitboxOffsetX16
	}

	x0 := (centerX16 - halfW16) / 16
	x1 := (centerX16 + halfW16) / 16
	y0 := (w.originY16 + w.hitboxOffsetY16) / 16
	y1 := y0 + w.hitboxH16/16
	return image.Rect(x0, y0, x1, y1)
}
