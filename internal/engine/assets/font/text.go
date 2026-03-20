package font

import (
	"bytes"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type FontText struct {
	source *text.GoTextFaceSource
}

func NewFontText(fsys fs.FS, path string) (*FontText, error) {
	font, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, err
	}
	src, err := text.NewGoTextFaceSource(bytes.NewReader(font))
	if err != nil {
		return nil, err
	}
	return &FontText{source: src}, nil
}

// NewFace creates a new text face with the given size.
func (t *FontText) NewFace(size float64) *text.GoTextFace {
	return &text.GoTextFace{
		Source: t.source,
		Size:   size,
	}
}

func (t *FontText) Draw(screen *ebiten.Image, msg string, size float64, op *text.DrawOptions) {
	if t.source == nil {
		return
	}

	text.Draw(screen, msg, t.NewFace(size), op)
}
