package draworder

import (
	"sort"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
)

func SortByGroundY(in []body.Collidable) []body.Collidable {
	out := make([]body.Collidable, len(in))
	copy(out, in)
	sort.SliceStable(out, func(i, j int) bool {
		_, yi16 := out[i].GetPosition16()
		_, yj16 := out[j].GetPosition16()
		return yi16 < yj16
	})
	return out
}

// Altitudable is satisfied by bodies that carry a vertical altitude axis
// (used in beat-em-up genres where entities can move on a Z plane).
type Altitudable interface {
	Altitude16() int
}

// SortByGroundYAltitude sorts collidable bodies by their effective ground depth,
// computed as y16 - altitude16. Bodies with lower effective depth are drawn first
// (they appear behind bodies with higher depth). Non-Altitudable bodies are treated
// as altitude=0, making their effective depth equal to y16.
func SortByGroundYAltitude(in []body.Collidable) []body.Collidable {
	out := make([]body.Collidable, len(in))
	copy(out, in)
	sort.SliceStable(out, func(i, j int) bool {
		_, yi16 := out[i].GetPosition16()
		_, yj16 := out[j].GetPosition16()

		var altI, altJ int
		if a, ok := out[i].(Altitudable); ok {
			altI = a.Altitude16()
		}
		if a, ok := out[j].(Altitudable); ok {
			altJ = a.Altitude16()
		}

		effI := yi16 - altI
		effJ := yj16 - altJ
		return effI < effJ
	})
	return out
}
