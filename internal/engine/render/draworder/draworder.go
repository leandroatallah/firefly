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
