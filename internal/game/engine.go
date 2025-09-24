package game

import (
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/audioplayer"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/scene"
)

type Game struct {
	sceneManager *scene.SceneManager
	state        GameState
}

func NewGame() *Game {
	return &Game{}
}

func (g *Game) Update() error {
	g.sceneManager.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.sceneManager.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}

func (g *Game) ChangeState(state GameState) *Game {
	state.SetContext(g)
	g.state = state
	if g.state != nil {
		g.state.OnStart()
	}
	return g
}

func (g *Game) SetSceneManager() *Game {
	sceneManager := scene.NewSceneManager()
	sceneFactory := scene.NewDefaultSceneFactory()

	sceneManager.SetFactory(sceneFactory)
	sceneFactory.SetManager(sceneManager)

	g.sceneManager = sceneManager
	return g
}

func (g *Game) SetAudioContext() *Game {
	ctx := audioplayer.NewContext()
	g.sceneManager.SetAudioContext(ctx)
	return g
}

func (g *Game) LoadAudioAssets() []*audioplayer.AudioItem {
	audioAssets := []*audioplayer.AudioItem{}
	files, err := os.ReadDir("assets")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		// TODO: Extend to WAV and MP3 format
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".ogg") {
			audioItem, err := audioplayer.LoadAudio("assets/" + file.Name())
			if err != nil {
				log.Fatal(err)
			}
			audioAssets = append(audioAssets, audioItem)
		}
	}
	return audioAssets
}

func (g *Game) SetAudioAssets() *Game {
	items := g.LoadAudioAssets()
	stream := make(map[string][]byte)
	for _, i := range items {
		stream[i.Name()] = i.Data()
	}
	g.sceneManager.SetAudioStream(stream)
	return g
}
