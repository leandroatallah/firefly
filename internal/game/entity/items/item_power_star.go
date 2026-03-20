package gameitems

import (
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	gamevfx "github.com/leandroatallah/firefly/internal/game/render/vfx"
)

// StarPowerItem is a collectible power-up item that activates the star skill.
type StarPowerItem struct {
	PowerUpItem
}

// NewStarPowerItem creates a new star power-up item.
func NewStarPowerItem(ctx *app.AppContext, x, y int, id string) (*StarPowerItem, error) {
	var powerItem *StarPowerItem

	powerItem = &StarPowerItem{
		PowerUpItem: PowerUpItem{
			activateSkill: func() {
				// Activate star skill when collected
				player, found := ctx.ActorManager.GetPlayer()
				if !found {
					return
				}
				// Pass self as item reference for respawn tracking
				if skillUser, ok := player.(interface {
					ActivateStarSkillWithItem(item interface{})
				}); ok {
					skillUser.ActivateStarSkillWithItem(powerItem)
				}
			},
		},
	}

	// Initialize base item
	base, err := createPowerUpBase(ctx, x, y, id, "internal/game/entity/items/star.json", powerItem.activateSkill)
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

// Update spawns rainbow aura particles around the item.
func (s *StarPowerItem) Update(space body.BodiesSpace) error {
	if !s.IsRemoved() {
		if ctx := s.AppContext(); ctx != nil && ctx.VFX != nil && ctx.FrameCount%5 == 0 {
			pos := s.Position()
			centerX := float64(pos.Min.X + pos.Dx()/2)
			centerY := float64(pos.Min.Y + pos.Dy()/2)
			gamevfx.SpawnStarParticles(ctx.VFX, centerX, centerY, 4)
		}
	}
	return s.BaseItem.Update(space)
}
