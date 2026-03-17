package audio

import (
	"bytes"
	"context"
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
	"github.com/leandroatallah/firefly/internal/engine/data/config"
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
	noSound      bool
	fadeCancel   map[string]context.CancelFunc
	paused       map[string]bool
}

func NewAudioManager() *AudioManager {
	initialVolume := 1.0
	noSound := config.Get().NoSound
	if noSound {
		initialVolume = 0.0
	}
	return &AudioManager{
		audioContext: audio.NewContext(sampleRate),
		audioPlayers: make(map[string]*audio.Player),
		volume:       initialVolume,
		noSound:      noSound,
		fadeCancel:   make(map[string]context.CancelFunc),
		paused:       make(map[string]bool),
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

func (am *AudioManager) PlayMusic(name string, loop bool) *audio.Player {
	if am.noSound {
		return nil
	}

	if cancel, ok := am.fadeCancel[name]; ok {
		cancel()
		delete(am.fadeCancel, name)
	}
	if cancel, ok := am.fadeCancel["_all"]; ok {
		cancel()
		delete(am.fadeCancel, "_all")
	}

	player, ok := am.audioPlayers[name]
	if !ok {
		log.Printf("audio player not found: %s", name)
		return nil
	}

	delete(am.paused, name)

	player.SetVolume(am.volume)
	player.Rewind()
	player.Play()

	if loop {
		go func() {
			for {
				for player.IsPlaying() {
					time.Sleep(100 * time.Millisecond)
					if _, exists := am.fadeCancel[name]; exists {
						return
					}
					if _, exists := am.fadeCancel["_all"]; exists {
						return
					}
				}
				if am.paused[name] {
					for am.paused[name] {
						time.Sleep(100 * time.Millisecond)
						if _, exists := am.fadeCancel[name]; exists {
							delete(am.paused, name)
							return
						}
					}
					player.Play()
					continue
				}
				player.Rewind()
				player.Play()
			}
		}()
	}

	return player
}

func (am *AudioManager) PauseMusic(name string) {
	player, ok := am.audioPlayers[name]
	if !ok {
		log.Printf("audio player not found: %s", name)
		return
	}
	am.paused[name] = true
	player.Pause()
}

func (am *AudioManager) ResumeMusic(name string) {
	if am.noSound {
		return
	}
	player, ok := am.audioPlayers[name]
	if !ok {
		log.Printf("audio player not found: %s", name)
		return
	}
	delete(am.paused, name)
	player.SetVolume(am.volume)
	player.Play()
}

func (am *AudioManager) PlaySound(name string) *audio.Player {
	if am.noSound {
		return nil
	}
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
	if am.noSound {
		return
	}
	am.volume = volume
	for _, player := range am.audioPlayers {
		player.SetVolume(am.volume)
	}
}

func (am *AudioManager) Volume() float64 {
	return am.volume
}

func (am *AudioManager) PauseAll() {
	for _, p := range am.audioPlayers {
		p.Pause()
	}
}

func (am *AudioManager) FadeOutAll(duration time.Duration) {
	if am.noSound {
		return
	}
	currentVolume := am.volume
	if currentVolume == 0 {
		return
	}

	if am.fadeCancel["_all"] != nil {
		am.fadeCancel["_all"]()
	}

	ctx, cancel := context.WithCancel(context.Background())
	am.fadeCancel["_all"] = cancel

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		startTime := time.Now()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				elapsed := time.Since(startTime)
				if elapsed >= duration {
					for _, p := range am.audioPlayers {
						p.SetVolume(currentVolume)
						p.Rewind()
						p.Pause()
					}
					delete(am.fadeCancel, "_all")
					return
				}

				progress := float64(elapsed) / float64(duration)
				newVolume := currentVolume * (1 - progress)
				if newVolume < 0 {
					newVolume = 0
				}
				for _, p := range am.audioPlayers {
					p.SetVolume(newVolume)
				}
			}
		}
	}()
}

func (am *AudioManager) FadeOut(name string, duration time.Duration) {
	if am.noSound {
		return
	}
	player, ok := am.audioPlayers[name]
	if !ok {
		log.Printf("audio player not found: %s", name)
		return
	}

	initialVolume := player.Volume()
	if initialVolume == 0 {
		return
	}

	if cancel, ok := am.fadeCancel[name]; ok {
		cancel()
		delete(am.fadeCancel, name)
	}
	if cancel, ok := am.fadeCancel["_all"]; ok {
		cancel()
		delete(am.fadeCancel, "_all")
	}

	ctx, cancel := context.WithCancel(context.Background())
	am.fadeCancel[name] = cancel

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		startTime := time.Now()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				elapsed := time.Since(startTime)
				if elapsed >= duration {
					player.SetVolume(am.volume)
					player.Rewind()
					player.Pause()
					delete(am.fadeCancel, name)
					return
				}

				progress := float64(elapsed) / float64(duration)
				newVolume := initialVolume * (1 - progress)
				if newVolume < 0 {
					newVolume = 0
				}
				player.SetVolume(newVolume)
			}
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
