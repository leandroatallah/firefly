package audioplayer

import (
	"bytes"
	"io"
	"os"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

const (
	sampleRate = 44100
	frequency  = 440
)

type AudioItem struct {
	name string
	data []byte
}

func (a *AudioItem) Name() string {
	return a.name
}
func (a *AudioItem) Data() []byte {
	return a.data
}

func LoadAudio(path string) (*AudioItem, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	bs := make([]byte, stat.Size())
	_, err = io.ReadFull(f, bs)
	if err != nil {
		return nil, err
	}

	return &AudioItem{path, bs}, nil
}

func NewContext() *audio.Context {
	return audio.NewContext(sampleRate)
}

func NewAudioPlayer(ctx *audio.Context, file []byte) (*audio.Player, error) {
	type audioStream interface {
		io.ReadSeeker
		Length() int64
	}

	// TODO: Extend to WAV and MP3 format
	s, err := vorbis.DecodeF32(bytes.NewReader(file))
	if err != nil {
		return nil, err
	}

	p, err := ctx.NewPlayerF32(s)
	if err != nil {
		return nil, err
	}

	return p, nil
}
