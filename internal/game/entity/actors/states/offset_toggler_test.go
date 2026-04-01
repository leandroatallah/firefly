package gamestates

import "testing"

func TestOffsetTogglerSequence(t *testing.T) {
	o := NewOffsetToggler(4)
	expected := []int{4, -4, 4, -4}
	for i, want := range expected {
		if got := o.Next(); got != want {
			t.Errorf("call %d: got %d, want %d", i+1, got, want)
		}
	}
}

func TestOffsetTogglerZero(t *testing.T) {
	o := NewOffsetToggler(0)
	for i := range 4 {
		if got := o.Next(); got != 0 {
			t.Errorf("call %d: got %d, want 0", i+1, got)
		}
	}
}

func TestOffsetTogglerNegativeInit(t *testing.T) {
	o := NewOffsetToggler(-3)
	expected := []int{-3, 3, -3, 3}
	for i, want := range expected {
		if got := o.Next(); got != want {
			t.Errorf("call %d: got %d, want %d", i+1, got, want)
		}
	}
}
