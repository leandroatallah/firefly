package gameplayer

import (
	"log"

	"github.com/boilerplate/ebiten-template/internal/engine/combat/inventory"
	"github.com/boilerplate/ebiten-template/internal/engine/combat/weapon"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/vfx"
	actors "github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/utils/fp16"
)

// StateOffsetEntry holds pixel-space spawn offset values for a single actor state.
type StateOffsetEntry struct {
	X int
	Y int
}

// BuildStateSpawnOffsets converts a map of state-name → pixel offset into a
// map of ActorStateEnum (as int) → fp16 offset pairs. Unknown state names are
// skipped with a log warning. A nil input returns nil.
func BuildStateSpawnOffsets(input map[string]StateOffsetEntry) map[int][2]int {
	if len(input) == 0 {
		return nil
	}
	result := make(map[int][2]int, len(input))
	for name, entry := range input {
		enum, ok := actors.GetStateEnum(name)
		if !ok {
			log.Printf("US-037: unknown state %q in state_spawn_offsets, skipping", name)
			continue
		}
		result[int(enum)] = [2]int{fp16.To16(entry.X), fp16.To16(entry.Y)}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

// NewClimberInventory creates a new inventory with light_blaster and heavy_cannon.
func NewClimberInventory(projectileManager combat.ProjectileManager, vfxManager vfx.Manager) *inventory.Inventory {
	inv := inventory.New()
	// Speeds in fp16 (scale factor 16): 6 pixels/frame = 96
	light := weapon.NewProjectileWeapon("light_blaster", 8, "bullet_small", 96, projectileManager, "muzzle_flash", fp16.To16(5), fp16.To16(10))
	light.SetDamage(1)
	light.SetVFXManager(vfxManager)
	inv.AddWeapon(light)
	inv.AddWeapon(NewHeavyCannonWeapon(projectileManager, vfxManager))
	return inv
}

// NewPlayerMeleeWeapon creates the player's melee weapon with a 3-step combo chain.
func NewPlayerMeleeWeapon() *weapon.MeleeWeapon {
	steps := []weapon.ComboStep{
		{Damage: 1, ActiveFrames: [2]int{4, 10}, HitboxW16: fp16.To16(24), HitboxH16: fp16.To16(16), HitboxOffsetX16: fp16.To16(12), HitboxOffsetY16: fp16.To16(0)},
		{Damage: 1, ActiveFrames: [2]int{3, 8}, HitboxW16: fp16.To16(28), HitboxH16: fp16.To16(16), HitboxOffsetX16: fp16.To16(14), HitboxOffsetY16: fp16.To16(-4)},
		{Damage: 2, ActiveFrames: [2]int{5, 12}, HitboxW16: fp16.To16(32), HitboxH16: fp16.To16(20), HitboxOffsetX16: fp16.To16(16), HitboxOffsetY16: fp16.To16(0)},
	}
	return weapon.NewMeleeWeapon("player_melee", 20, 15, steps)
}

// NewHeavyCannonWeapon creates the heavy_cannon weapon.
func NewHeavyCannonWeapon(projectileManager combat.ProjectileManager, vfxManager vfx.Manager) combat.Weapon {
	// Speed in fp16: 9 pixels/frame = 144
	w := weapon.NewProjectileWeapon("heavy_cannon", 30, "bullet_large", 144, projectileManager, "muzzle_flash", fp16.To16(5), fp16.To16(10))
	w.SetDamage(3)
	w.SetVFXManager(vfxManager)
	return w
}
