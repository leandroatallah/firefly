package gameitems

import (
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	gamevfx "github.com/leandroatallah/firefly/internal/game/render/vfx"
)

// GrowPowerItem is a collectible power-up item that activates the grow skill.
type GrowPowerItem struct {
	PowerUpItem
}

// NewGrowPowerItem creates a new grow power-up item.
func NewGrowPowerItem(ctx *app.AppContext, x, y int, id string) (*GrowPowerItem, error) {
	var powerItem *GrowPowerItem
	
	powerItem = &GrowPowerItem{
		PowerUpItem: PowerUpItem{
			activateSkill: func() {
				// Activate grow skill when collected
				player, found := ctx.ActorManager.GetPlayer()
				if !found {
					return
				}
				// Pass self as item reference for respawn tracking
				if skillUser, ok := player.(interface {
					ActivateGrowSkillWithItem(item interface{})
				}); ok {
					skillUser.ActivateGrowSkillWithItem(powerItem)
				}
			},
		},
	}
	
	// Initialize base item
	base, err := createPowerUpBase(ctx, x, y, id, "assets/entities/items/item-power-grow.json", powerItem.activateSkill)
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

// Update spawns orange aura particles around the item.
func (g *GrowPowerItem) Update(space body.BodiesSpace) error {
	if !g.IsRemoved() {
		if ctx := g.AppContext(); ctx != nil && ctx.VFX != nil && ctx.FrameCount%5 == 0 {
			pos := g.Position()
			centerX := float64(pos.Min.X + pos.Dx()/2)
			centerY := float64(pos.Min.Y + pos.Dy()/2)
			gamevfx.SpawnGrowAuraParticles(ctx.VFX, centerX, centerY, 4)
		}
	}
	return g.BaseItem.Update(space)
}
