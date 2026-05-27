package schemas

import (
	"encoding/json"
	"testing"
)

// T-S1 (story 070): SpriteOffset unmarshals with x_flipped present → XFlipped
// is a non-nil pointer carrying the JSON value. Pins the schema contract that
// distinguishes per-facing-direction X overrides from the single X used today.
func TestSpriteOffset_UnmarshalsXFlipped_Present(t *testing.T) {
	raw := []byte(`{"x":-4,"y":2,"x_flipped":6}`)

	var o SpriteOffset
	if err := json.Unmarshal(raw, &o); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if o.X != -4 {
		t.Errorf("X = %d, want -4", o.X)
	}
	if o.Y != 2 {
		t.Errorf("Y = %d, want 2", o.Y)
	}
	if o.XFlipped == nil {
		t.Fatal("XFlipped = nil; want non-nil pointer when x_flipped key present")
	}
	if *o.XFlipped != 6 {
		t.Errorf("*XFlipped = %d, want 6", *o.XFlipped)
	}
}

// T-S2 (story 070): SpriteOffset unmarshals when x_flipped is absent → XFlipped
// is nil. Guarantees zero regression for all existing actor JSON files that
// only declare {x, y} (story 068 schema).
func TestSpriteOffset_UnmarshalsXFlipped_Absent(t *testing.T) {
	raw := []byte(`{"x":-4,"y":2}`)

	var o SpriteOffset
	if err := json.Unmarshal(raw, &o); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if o.X != -4 || o.Y != 2 {
		t.Errorf("(X, Y) = (%d, %d), want (-4, 2)", o.X, o.Y)
	}
	if o.XFlipped != nil {
		t.Errorf("XFlipped = %v, want nil when x_flipped key omitted", *o.XFlipped)
	}
}

// T-S3 (story 070): SpriteOffset unmarshals with x_flipped explicitly 0 →
// XFlipped is a non-nil pointer to 0. The pointer-based encoding must
// distinguish "explicit zero override" from "unset" — otherwise authors lose
// the ability to declare a zero left-facing nudge.
func TestSpriteOffset_UnmarshalsXFlipped_ExplicitZero(t *testing.T) {
	raw := []byte(`{"x":-4,"y":2,"x_flipped":0}`)

	var o SpriteOffset
	if err := json.Unmarshal(raw, &o); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if o.XFlipped == nil {
		t.Fatal("XFlipped = nil; want non-nil pointer to 0 for explicit zero override")
	}
	if *o.XFlipped != 0 {
		t.Errorf("*XFlipped = %d, want 0 (explicit override)", *o.XFlipped)
	}
}
