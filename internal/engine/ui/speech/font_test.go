package speech

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
)

func TestSpeechFont(t *testing.T) {
	// We can't easily create a real FontText without a font file,
	// but we can pass a nil source to test NewSpeechFont and basic Draw call.
	sf := NewSpeechFont(&font.FontText{}, 12.0, 1.5)

	if sf.size != 12.0 {
		t.Errorf("Expected size 12.0, got %f", sf.size)
	}

	if sf.LineSpacing != 1.5 {
		t.Errorf("Expected LineSpacing 1.5, got %f", sf.LineSpacing)
	}

	// Should not panic even if source is empty (FontText.Draw handles it)
	img := ebiten.NewImage(10, 10)
	sf.Draw(img, "test", &text.DrawOptions{})
}
