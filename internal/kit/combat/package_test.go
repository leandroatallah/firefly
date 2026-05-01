package combat_test

import (
	"testing"

	_ "github.com/boilerplate/ebiten-template/internal/kit/combat"
	_ "github.com/boilerplate/ebiten-template/internal/kit/combat/inventory"
	_ "github.com/boilerplate/ebiten-template/internal/kit/combat/melee"
	_ "github.com/boilerplate/ebiten-template/internal/kit/combat/projectile"
	_ "github.com/boilerplate/ebiten-template/internal/kit/combat/weapon"
)

// TestKitCombatPackagesExist is a compile-time assertion that all kit-combat
// sub-packages (root, inventory, melee, projectile, weapon) exist and are
// importable at the new kit/combat path. The test body is intentionally empty
// — the assertion lives in the blank imports above. If any of those packages
// is missing, the file fails to build and the test cannot run.
func TestKitCombatPackagesExist(_ *testing.T) {}
