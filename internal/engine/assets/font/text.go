package font

import (
	"bytes"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type FontText struct {
	source *text.GoTextFaceSource
}

func NewFontText(path string) (*FontText, error) {
	font, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	src, err := text.NewGoTextFaceSource(bytes.NewReader(font))
	if err != nil {
		return nil, err
	}
	return &FontText{source: src}, nil
}

func (t *FontText) Draw(screen *ebiten.Image, msg string, size float64, op *text.DrawOptions) {
	if t.source == nil {
		return
	}

	text.Draw(screen, msg, &text.GoTextFace{
		Source: t.source,
		Size:   size,
	}, op)
}
