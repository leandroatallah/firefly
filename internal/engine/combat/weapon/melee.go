package weapon

import (
	"image"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
	"github.com/boilerplate/ebiten-template/internal/engine/utils"
)

// ComboStep defines the per-step hitbox and damage for a melee combo chain.
type ComboStep struct {
	Damage          int
	StartupFrames   int // frames before the swing begins (hitbox inactive)
	ActiveFrames    [2]int
	HitboxW16       int
	HitboxH16       int
	HitboxOffsetX16 int
	HitboxOffsetY16 int
}

// MeleeWeapon is a close-range swing weapon that activates a hitbox during a
// configurable active-frame window within its animation.
type MeleeWeapon struct {
	id                      string
	cooldownFrames          int
	currentCooldown         int
	comboWindowFrames       int
	postComboCooldownFrames int
	steps                   []ComboStep

	owner interface{}

	stepIndex       int
	windowRemaining int
	startup         utils.DelayTrigger
	swinging        bool
	swingFrame      int
	hitThisSwing    map[combat.Damageable]struct{}

	faceDir              animation.FacingDirectionEnum
	originX16, originY16 int
}

// NewMeleeWeapon constructs a MeleeWeapon with combo-step configuration.
func NewMeleeWeapon(id string, cooldownFrames, comboWindowFrames int, steps []ComboStep) *MeleeWeapon {
	return &MeleeWeapon{
		id:                id,
		cooldownFrames:    cooldownFrames,
		comboWindowFrames: comboWindowFrames,
		steps:             steps,
		hitThisSwing:      make(map[combat.Damageable]struct{}),
	}
}

// ID returns the weapon identifier.
func (w *MeleeWeapon) ID() string { return w.id }

// CanFire returns true when there is no active cooldown or startup.
func (w *MeleeWeapon) CanFire() bool { return w.currentCooldown == 0 && !w.startup.IsEnabled() }

// Cooldown returns the remaining cooldown frames.
func (w *MeleeWeapon) Cooldown() int { return w.currentCooldown }

// SetCooldown sets the cooldown to the given value.
func (w *MeleeWeapon) SetCooldown(frames int) { w.currentCooldown = frames }

// SetOwner sets the owner reference used for faction checks.
func (w *MeleeWeapon) SetOwner(owner interface{}) { w.owner = owner }

// SetPostComboCooldownFrames sets the cooldown applied after the final combo step finishes.
func (w *MeleeWeapon) SetPostComboCooldownFrames(frames int) { w.postComboCooldownFrames = frames }

// StepIndex returns the current combo step index.
func (w *MeleeWeapon) StepIndex() int { return w.stepIndex }

// IsSwinging returns true while a swing animation is in progress.
func (w *MeleeWeapon) IsSwinging() bool { return w.swinging }

// IsInStartup returns true while the weapon is in its pre-swing startup delay.
func (w *MeleeWeapon) IsInStartup() bool { return w.startup.IsEnabled() }

// ComboWindowRemaining returns the number of frames left in the combo window.
func (w *MeleeWeapon) ComboWindowRemaining() int { return w.windowRemaining }

// Steps returns the configured combo steps.
func (w *MeleeWeapon) Steps() []ComboStep { return w.steps }

// ResetCombo resets the combo chain to step 0 without interrupting an in-flight swing.
func (w *MeleeWeapon) ResetCombo() {
	w.stepIndex = 0
	w.windowRemaining = 0
	w.startup.Reset()
}

// AdvanceCombo advances to the next step if the combo window is open and a next step exists.
// Returns true if the combo was advanced.
func (w *MeleeWeapon) AdvanceCombo() bool {
	if w.windowRemaining == 0 {
		return false
	}
	if w.stepIndex >= len(w.steps)-1 {
		return false
	}
	w.stepIndex++
	w.windowRemaining = 0
	return true
}

// Fire begins a melee swing at the current combo step.
// If the step has startup frames, the swing is deferred until the countdown elapses.
func (w *MeleeWeapon) Fire(x16, y16 int, faceDir animation.FacingDirectionEnum, _ body.ShootDirection, _ int) {
	if !w.CanFire() {
		return
	}
	w.originX16 = x16
	w.originY16 = y16
	w.faceDir = faceDir
	w.windowRemaining = 0

	if n := w.steps[w.stepIndex].StartupFrames; n > 0 {
		w.startup.Enable(n)
		return
	}
	w.startSwing()
}

func (w *MeleeWeapon) startSwing() {
	w.swinging = true
	w.swingFrame = 0
	w.hitThisSwing = make(map[combat.Damageable]struct{})
	w.currentCooldown = w.cooldownFrames
}

// Update advances the swing frame, decrements the cooldown, and manages the combo window.
func (w *MeleeWeapon) Update() {
	if w.currentCooldown > 0 {
		w.currentCooldown--
	}
	w.startup.Update()
	if w.startup.Trigger() {
		w.startup.Reset()
		w.startSwing()
		return
	}
	if w.startup.IsEnabled() {
		return
	}
	if w.swinging {
		w.swingFrame++
		if w.swingFrame > w.steps[w.stepIndex].ActiveFrames[1] {
			w.swinging = false
			if w.stepIndex < len(w.steps)-1 {
				w.windowRemaining = w.comboWindowFrames
			} else {
				w.ResetCombo()
				if w.postComboCooldownFrames > w.currentCooldown {
					w.currentCooldown = w.postComboCooldownFrames
				}
			}
		}
	} else if w.windowRemaining > 0 {
		w.windowRemaining--
		if w.windowRemaining == 0 {
			w.ResetCombo()
		}
	}
}

// IsHitboxActive returns true when the current swing frame is within the active-frame window.
func (w *MeleeWeapon) IsHitboxActive() bool {
	if !w.swinging {
		return false
	}
	step := w.steps[w.stepIndex]
	return w.swingFrame >= step.ActiveFrames[0] && w.swingFrame <= step.ActiveFrames[1]
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
		if w.owner != nil && b == w.owner {
			continue
		}

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
		target.TakeDamage(w.steps[w.stepIndex].Damage)
	}
}

// hitboxRect computes the query rectangle in pixel space for the current step.
func (w *MeleeWeapon) hitboxRect() image.Rectangle {
	step := w.steps[w.stepIndex]
	halfW16 := step.HitboxW16 / 2

	var centerX16 int
	if w.faceDir == animation.FaceDirectionLeft {
		centerX16 = w.originX16 - step.HitboxOffsetX16
	} else {
		centerX16 = w.originX16 + step.HitboxOffsetX16
	}

	x0 := (centerX16 - halfW16) / 16
	x1 := (centerX16 + halfW16) / 16
	y0 := (w.originY16 + step.HitboxOffsetY16) / 16
	y1 := y0 + step.HitboxH16/16
	return image.Rect(x0, y0, x1, y1)
}
