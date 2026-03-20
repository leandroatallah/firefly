package gameitems

import (
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	gamevfx "github.com/leandroatallah/firefly/internal/game/render/vfx"
)

// FreezePowerItem is a collectible power-up item that activates the freeze skill.
type FreezePowerItem struct {
	PowerUpItem
}

// NewFreezePowerItem creates a new freeze power-up item.
func NewFreezePowerItem(ctx *app.AppContext, x, y int, id string) (*FreezePowerItem, error) {
	var powerItem *FreezePowerItem

	powerItem = &FreezePowerItem{
		PowerUpItem: PowerUpItem{
			activateSkill: func() {
				// Activate freeze skill when collected
				player, found := ctx.ActorManager.GetPlayer()
				if !found {
					return
				}
				// Pass self as item reference for respawn tracking
				if skillUser, ok := player.(interface {
					ActivateFreezeSkillWithItem(item interface{})
				}); ok {
					skillUser.ActivateFreezeSkillWithItem(powerItem)
				}
			},
		},
	}

	// Initialize base item
	base, err := createPowerUpBase(ctx, x, y, id, "internal/game/entity/items/power.json", powerItem.activateSkill)
	if err != nil {
		return nil, err
	}
	powerItem.BaseItem = *base

	// Set collection feedback callback
	powerItem.SetOnCollect(func() {
		// Play sound effect at reduced volume
		if ctx.AudioManager != nil {
			ctx.AudioManager.PlaySoundAtVolume("assets/audio/Booster.ogg", 0.3)
		}
		// Trigger screen flash via AppContext
		ctx.VFX.TriggerScreenFlash()
	})

	return powerItem, nil
}

// Update spawns blue aura particles around the item.
func (f *FreezePowerItem) Update(space body.BodiesSpace) error {
	if !f.IsRemoved() {
		if ctx := f.AppContext(); ctx != nil && ctx.VFX != nil && ctx.FrameCount%5 == 0 {
			pos := f.Position()
			centerX := float64(pos.Min.X + pos.Dx()/2)
			centerY := float64(pos.Min.Y + pos.Dy()/2)
			gamevfx.SpawnFreezeAuraParticles(ctx.VFX, centerX, centerY, 4)
		}
	}
	return f.BaseItem.Update(space)
}
