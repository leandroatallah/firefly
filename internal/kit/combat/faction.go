package combat

import contractscombat "github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"

// Faction aliases the canonical contracts/combat.Faction type so that kit
// callers can refer to `kitcombat.Faction` without importing the contracts
// package directly.
type Faction = contractscombat.Faction

const (
	FactionNeutral = contractscombat.FactionNeutral
	FactionPlayer  = contractscombat.FactionPlayer
	FactionEnemy   = contractscombat.FactionEnemy
)
