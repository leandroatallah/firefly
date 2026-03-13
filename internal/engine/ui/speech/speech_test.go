package speech

import (
	"image/color"
	"testing"
)

func TestSpeechBase(t *testing.T) {
	sb := NewSpeechBase(nil)
	sb.SetID("test-id")

	if sb.ID() != "test-id" {
		t.Errorf("Expected ID to be test-id, got %s", sb.ID())
	}

	if sb.Visible() {
		t.Error("Expected Visible to be false initially")
	}

	sb.Show()
	if !sb.Visible() {
		t.Error("Expected Visible to be true after Show()")
	}

	sb.Hide()
	if sb.Visible() {
		t.Error("Expected Visible to be false after Hide()")
	}

	sb.SetColor(color.White)
	if sb.Color() != color.White {
		t.Error("Expected color to be White")
	}

	if sb.Image(nil) != nil {
		t.Error("Expected Image to return nil")
	}
}

func TestSpeechBase_Spelling(t *testing.T) {
	sb := NewSpeechBase(nil)
	sb.SetSpellingDelay(10)
	sb.SetSpeed(0) // Test default speed 4

	msg := "Hello"
	// Before spelling delay
	if text := sb.Text(msg, 0); text != "" {
		t.Errorf("Expected empty text before delay, got %s", text)
	}

	// Reach spelling delay
	for i := 0; i < 10; i++ {
		sb.Update()
	}

	if text := sb.Text(msg, 0); text != "" {
		t.Errorf("Expected empty text at delay, got %s", text)
	}

	// Update with default speed 4
	sb.Update() // 11
	sb.Update() // 12
	sb.Update() // 13
	sb.Update() // 14 -> spellingCount = 1

	if text := sb.Text(msg, 0); text != "H" {
		t.Errorf("Expected H, got %s", text)
	}

	// Test updating speed through Text method
	if text := sb.Text(msg, 2); text != "H" {
		t.Errorf("Expected H, got %s", text)
	}
	if sb.GetSpeed() != 2 {
		t.Errorf("Expected speed 2, got %d", sb.GetSpeed())
	}

	if sb.IsSpellingComplete() {
		t.Error("Expected spelling not to be complete")
	}

	sb.CompleteSpelling()
	if !sb.IsSpellingComplete() {
		t.Error("Expected spelling to be complete after CompleteSpelling()")
	}

	if text := sb.Text(msg, 0); text != "Hello" {
		t.Errorf("Expected Hello, got %s", text)
	}

	sb.ResetText()
	if sb.Count() != 0 {
		t.Error("Expected count to be 0 after ResetText()")
	}
	if sb.IsSpellingComplete() {
		t.Error("Expected spelling not to be complete after ResetText()")
	}
}

func TestSpeechBase_EmptyText(t *testing.T) {
	sb := NewSpeechBase(nil)
	if sb.IsSpellingComplete() {
		t.Error("Empty text should not be complete")
	}
}

func TestSpeechBase_PositionAndSpeed(t *testing.T) {
	sb := NewSpeechBase(nil)

	sb.SetPosition("top")
	if sb.GetPosition() != "top" {
		t.Errorf("Expected position top, got %s", sb.GetPosition())
	}

	sb.SetSpeed(10)
	if sb.GetSpeed() != 10 {
		t.Errorf("Expected speed 10, got %d", sb.GetSpeed())
	}
}
