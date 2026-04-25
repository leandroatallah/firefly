package gamestates_test

import (
	"fmt"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	gamestates "github.com/boilerplate/ebiten-template/internal/game/entity/actors/states"
)

// TestMeleeAttackStepStates_LengthAndNonZero verifies that the helper returns
// a slice of the requested length and every entry is a registered (non-zero)
// state enum.
func TestMeleeAttackStepStates_LengthAndNonZero(t *testing.T) {
	tests := []struct {
		name string
		n    int
	}{
		{name: "single step", n: 1},
		{name: "three combo steps", n: 3},
		{name: "five combo steps", n: 5},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := gamestates.MeleeAttackStepStates(tc.n)
			if len(got) != tc.n {
				t.Fatalf("MeleeAttackStepStates(%d) len = %d, want %d", tc.n, len(got), tc.n)
			}
			for i, enum := range got {
				if enum == 0 {
					t.Errorf("index %d: state enum = 0, want non-zero (should be registered)", i)
				}
			}
		})
	}
}

// TestMeleeAttackStepStates_Idempotent verifies that calling the helper
// repeatedly with the same N returns the same enum values — the underlying
// state registry must not allocate new enums for already-registered names.
func TestMeleeAttackStepStates_Idempotent(t *testing.T) {
	first := gamestates.MeleeAttackStepStates(3)
	second := gamestates.MeleeAttackStepStates(3)

	if len(first) != len(second) {
		t.Fatalf("len mismatch: first=%d, second=%d", len(first), len(second))
	}
	for i := range first {
		if first[i] != second[i] {
			t.Errorf("index %d: first=%v, second=%v (should be idempotent)", i, first[i], second[i])
		}
	}
}

// TestMeleeAttackStepStates_RegistersNamedPattern verifies that the helper
// registers states under the name pattern "melee_attack_step_<i>" and that
// GetStateEnum resolves to the same enum returned by the helper.
func TestMeleeAttackStepStates_RegistersNamedPattern(t *testing.T) {
	const n = 3
	got := gamestates.MeleeAttackStepStates(n)

	for i := 0; i < n; i++ {
		name := fmt.Sprintf("melee_attack_step_%d", i)
		enum, ok := actors.GetStateEnum(name)
		if !ok {
			t.Errorf("GetStateEnum(%q) not registered; MeleeAttackStepStates should register the per-step name", name)
			continue
		}
		if enum != got[i] {
			t.Errorf("index %d (%q): GetStateEnum = %v, helper returned %v", i, name, enum, got[i])
		}
	}
}

// TestMeleeAttackStepStates_DistinctPerStep verifies that different step
// indices map to distinct enum values (so the animation layer can dispatch
// on step without collisions).
func TestMeleeAttackStepStates_DistinctPerStep(t *testing.T) {
	got := gamestates.MeleeAttackStepStates(3)
	seen := make(map[actors.ActorStateEnum]int)
	for i, enum := range got {
		if prev, ok := seen[enum]; ok {
			t.Errorf("duplicate enum %v at indices %d and %d", enum, prev, i)
		}
		seen[enum] = i
	}
}
