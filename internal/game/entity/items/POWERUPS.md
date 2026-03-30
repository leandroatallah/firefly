# Power-Up Items Guide

This boilerplate includes `FallingPlatformType` and the `PowerUpItem` base as examples.
To add a new power-up (e.g., a "Grow" or "Freeze" skill):

## 1. Define the skill in `internal/game/physics/skill/`

```go
type GrowSkill struct { /* ... */ }
func (s *GrowSkill) RequestActivation() { /* ... */ }
func (s *GrowSkill) IsActive() bool { /* ... */ }
func (s *GrowSkill) Reset(owner interface{}) { /* ... */ }
```

## 2. Add skill fields to the player (`ClimberPlayer`)

```go
type ClimberPlayer struct {
    *platformer.PlatformerCharacter
    baseSpeed int
    growSkill *skill.GrowSkill
    *gameplayermethods.PlayerDeathBehavior
}
```

Wire it in `NewClimberPlayer`:
```go
character.AddSkill(player.growSkill)
```

## 3. Create the item in `internal/game/entity/items/`

```go
type GrowPowerItem struct { PowerUpItem }

func NewGrowPowerItem(ctx *app.AppContext, x, y int, id string) (*GrowPowerItem, error) {
    // Use createPowerUpBase + activateSkill callback
}
```

## 4. Register in `init_items.go`

```go
const GrowPowerUpType items.ItemType = "GROW_POWER_UP"

// in InitItemMap:
GrowPowerUpType: func(x, y int, id string) items.Item {
    return itemFactoryOrFatal(NewGrowPowerItem(ctx, x, y, id))
},
```

## 5. Add invincibility guard in enemies (optional)

In `bat.go` / `wolf.go` `OnTouch`:
```go
if invincible, ok := owner.(interface{ IsGrowActive() bool }); ok {
    if invincible.IsGrowActive() {
        return
    }
}
```

## 6. Reset on player death / phase restart

In `internal/game/scenes/phases/scene.go`, call `player.ResetSkills()` after death
and in the phase restart block.
