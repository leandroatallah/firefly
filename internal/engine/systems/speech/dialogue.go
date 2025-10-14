package speech

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Dialogue manages a sequence of dialogue lines.
type Dialogue struct {
	speech      Speech
	lines       []string
	currentLine int
}

func NewDialogue(speech Speech) *Dialogue {
	return &Dialogue{speech: speech}
}

func (d *Dialogue) Update() error {
	if err := d.speech.Update(); err != nil {
		return err
	}
	return nil
}

func (d *Dialogue) Draw(screen *ebiten.Image) {
	if len(d.lines) == 0 {
		return
	}

	line := d.GetCurrentLine()
	d.speech.Draw(screen, line)
}

func (d *Dialogue) GetCurrentLine() string {
	if d.currentLine >= len(d.lines) {
		return ""
	}
	return d.lines[d.currentLine]
}

func (d *Dialogue) AddLine(line string) {
	d.lines = append(d.lines, line)
}

func (d *Dialogue) SetLines(lines []string) {
	d.lines = lines
}

func (d *Dialogue) SetSpellingDelay(delay int) {
	d.speech.SetSpellingDelay(delay)
}

func (d *Dialogue) NextLine() bool {
	d.currentLine++

	if d.currentLine >= len(d.lines) {
		return false
	}

	d.speech.ResetText()

	return true
}

func (d *Dialogue) Show() {
	d.speech.Show()
}

func (d *Dialogue) Hide() {
	d.speech.Hide()
}

func (d *Dialogue) IsSpellingComplete() bool {
	return d.speech.IsSpellingComplete()
}

func (d *Dialogue) CompleteSpelling() {
	d.speech.CompleteSpelling()
}
