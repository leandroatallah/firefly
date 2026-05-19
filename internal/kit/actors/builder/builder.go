// Package kitbuilder provides a generic player-builder utility that applies
// skills, inventory, and melee-weapon wiring in a genre-agnostic way.
package kitbuilder

import (
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/vfx"
	"github.com/boilerplate/ebiten-template/internal/engine/data/schemas"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	enginebuilder "github.com/boilerplate/ebiten-template/internal/engine/entity/actors/builder"
	"github.com/boilerplate/ebiten-template/internal/kit/combat/weapon"
	kitskills "github.com/boilerplate/ebiten-template/internal/kit/skills"
)

// PlayerDeps configures optional wiring applied by BuildPlayer.
type PlayerDeps struct {
	SkillDeps   kitskills.SkillDeps     // required if SpriteData has skills
	Inventory   interface{}             // optional; applied via SetInventory if non-nil
	MeleeWeapon *weapon.MeleeWeapon     // optional; applied via SetMelee if non-nil
	VFXManager  vfx.Manager             // passed to SetMelee
	SpriteData  *schemas.SpriteData     // required if applying skills; nil => skip skills
	WireState   func(*actors.Character) // optional; e.g., WireStateContributors
}

// playerWiring is the optional interface used to inject inventory/melee.
// Players that implement this interface receive the optional wiring;
// players that do not are returned unchanged (no-op).
type playerWiring interface {
	SetInventory(interface{})
	SetMelee(w *weapon.MeleeWeapon, vfxMgr vfx.Manager)
	GetCharacter() *actors.Character
}

// BuildPlayer applies skills (when SpriteData non-nil), then optionally injects
// Inventory and MeleeWeapon. Returns p untouched on success.
// If p does not implement the internal playerWiring interface, the function is a no-op.
func BuildPlayer[T actors.ActorEntity](p T, deps PlayerDeps) (T, error) {
	pw, ok := any(p).(playerWiring)
	if !ok {
		return p, nil // p does not opt-in to inventory/melee wiring
	}
	if deps.Inventory != nil {
		pw.SetInventory(deps.Inventory)
	}
	if deps.MeleeWeapon != nil {
		pw.SetMelee(deps.MeleeWeapon, deps.VFXManager)
	}
	if deps.SpriteData != nil {
		skills := kitskills.FromConfig(deps.SpriteData.Skills, deps.SkillDeps)
		if err := enginebuilder.ApplySkills(p, skills); err != nil {
			return p, err
		}
	}
	if deps.WireState != nil {
		deps.WireState(pw.GetCharacter())
	}
	return p, nil
}
