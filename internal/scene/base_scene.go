// - Add scene system (menu, playing, paused, game over)
// - Implement scene transitions and lifecycle management
package scene

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/leandroatallah/firefly/internal/audioplayer"
	"github.com/leandroatallah/firefly/internal/physics"
)

type Scene interface {
	Draw(screen *ebiten.Image)
	Update() error
	OnStart()
	OnFinish()
	SetManager(manager *SceneManager)
}

type BaseScene struct {
	boundaries        []physics.Body
	count             int
	Manager           *SceneManager
	audioPlayerStream map[string]*audio.Player
}

func NewScene() *BaseScene {
	return &BaseScene{}
}

func (s *BaseScene) Draw(screen *ebiten.Image) {
	panic("You should implement this method in derivated structs")
}

func (s *BaseScene) Update() error {
	panic("You should implement this method in derivated structs")
}

func (s *BaseScene) OnStart() {
	panic("You should implement this method in derivated structs")
}

func (s *BaseScene) OnFinish() {
	panic("You should implement this method in derivated structs")
}

func (s *BaseScene) Exit() {
	panic("You should implement this method in derivated structs")
}

func (s *BaseScene) AddBoundaries(boundaries ...physics.Body) {
	for _, o := range boundaries {
		s.boundaries = append(s.boundaries, o)
	}
}

func (s *BaseScene) SetManager(manager *SceneManager) {
	s.Manager = manager
}

// Audio methods
func (s *BaseScene) SetAudioStream(list []string) {
	s.audioPlayerStream = make(map[string]*audio.Player)
	for _, item := range list {
		player := s.NewAudioPlayer(item)
		s.audioPlayerStream[item] = player
	}
}

func (s *BaseScene) NewAudioPlayer(path string) *audio.Player {
	item := s.Manager.GetAudioData(path)
	player, err := audioplayer.NewAudioPlayer(s.Manager.AudioContext(), item)
	if err != nil {
		log.Fatal("Unable to create audio player")
	}
	return player
}

func (s BaseScene) PlayAudio(path string) {
	if player, exists := s.audioPlayerStream[path]; exists {
		player.Play()
	}
}

func (s *BaseScene) PauseAudio(path string) {
	if player, exists := s.audioPlayerStream[path]; exists {
		player.Pause()
	}
}
