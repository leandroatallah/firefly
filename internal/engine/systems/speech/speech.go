package speech

import "github.com/hajimehoshi/ebiten/v2"

type Speech interface {
	ID() string
	Show()
	Hide()
	Visible() bool
	Text(msg string) string
	ResetText()
	SetSpellingDelay(d int)
	IsSpellingComplete() bool
	CompleteSpelling()
	Count() int
	Update() error
	Draw(screen *ebiten.Image, text string)
}
