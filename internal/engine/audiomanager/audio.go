package audiomanager

import (
	"bytes"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const (
	sampleRate = 44100
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

type AudioManager struct {
	audioContext *audio.Context
	audioPlayers map[string]*audio.Player
	volume       float64
}

func NewAudioManager() *AudioManager {
	return &AudioManager{
		audioContext: audio.NewContext(sampleRate),
		audioPlayers: make(map[string]*audio.Player),
		volume:       1.0,
	}
}

func (am *AudioManager) Load(path string) (*AudioItem, error) {
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

func (am *AudioManager) LoadFromFS(fs fs.FS, path string) (*AudioItem, error) {
	f, err := fs.Open(path)
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

func (am *AudioManager) Add(name string, data []byte) {
	var s io.ReadSeeker
	var err error

	switch {
	case strings.HasSuffix(name, ".mp3"):
		s, err = mp3.DecodeWithSampleRate(sampleRate, bytes.NewReader(data))
		if err != nil {
			log.Printf("failed to decode mp3 file: %v", err)
			return
		}
	case strings.HasSuffix(name, ".ogg"):
		s, err = vorbis.DecodeWithSampleRate(sampleRate, bytes.NewReader(data))
		if err != nil {
			log.Printf("failed to decode ogg file: %v", err)
			return
		}
	case strings.HasSuffix(name, ".wav"):
		s, err = wav.DecodeWithSampleRate(sampleRate, bytes.NewReader(data))
		if err != nil {
			log.Printf("failed to decode wav file: %v", err)
			return
		}
	default:
		log.Printf("unsupported audio format: %s", name)
		return
	}

	p, err := am.audioContext.NewPlayer(s)
	if err != nil {
		log.Printf("failed to create audio player: %v", err)
		return
	}
	am.audioPlayers[name] = p
}

func (am *AudioManager) PlayMusic(name string) *audio.Player {
	player, ok := am.audioPlayers[name]
	if !ok {
		log.Printf("audio player not found: %s", name)
		return nil
	}
	player.SetVolume(am.volume)
	player.Play()
	return player
}

func (am *AudioManager) PauseMusic(name string) {
	player, ok := am.audioPlayers[name]
	if !ok {
		log.Printf("audio player not found: %s", name)
		return
	}
	player.Pause()
}

func (am *AudioManager) PlaySound(name string) *audio.Player {
	player, ok := am.audioPlayers[name]
	if !ok {
		log.Printf("audio player not found: %s", name)
		return nil
	}
	player.SetVolume(am.volume)
	player.Rewind()
	player.Play()
	return player
}

func (am *AudioManager) SetVolume(volume float64) {
	am.volume = volume
	for _, player := range am.audioPlayers {
		player.SetVolume(am.volume)
	}
}

func (am *AudioManager) Volume() float64 {
	return am.volume
}

func (am *AudioManager) PauseAll() {
	for _, player := range am.audioPlayers {
		player.Pause()
	}
}

func (am *AudioManager) FadeOutAll(duration time.Duration) {
	initialVolume := am.volume
	if initialVolume == 0 {
		return
	}

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		startTime := time.Now()

		for range ticker.C {
			elapsed := time.Since(startTime)
			if elapsed >= duration {
				am.SetVolume(0)
				am.PauseAll()
				return
			}

			progress := float64(elapsed) / float64(duration)
			newVolume := initialVolume * (1 - progress)
			if newVolume < 0 {
				newVolume = 0
			}
			am.SetVolume(newVolume)
		}
	}()
}

func (am *AudioManager) FadeOut(name string, duration time.Duration) {
	player, ok := am.audioPlayers[name]
	if !ok {
		log.Printf("audio player not found: %s", name)
		return
	}

	initialVolume := player.Volume()
	if initialVolume == 0 {
		return
	}

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		startTime := time.Now()

		for range ticker.C {
			elapsed := time.Since(startTime)
			if elapsed >= duration {
				player.SetVolume(0)
				player.Pause()
				return
			}

			progress := float64(elapsed) / float64(duration)
			newVolume := initialVolume * (1 - progress)
			if newVolume < 0 {
				newVolume = 0
			}
			player.SetVolume(newVolume)
		}
	}()
}

func (am *AudioManager) IsPlayingSomething() bool {
	for _, player := range am.audioPlayers {
		if player.IsPlaying() {
			return true
		}
	}
	return false
}

func (am *AudioManager) IsPlaying(name string) bool {
	audio, ok := am.audioPlayers[name]
	if !ok {
		return false
	}
	return audio.IsPlaying()
}
