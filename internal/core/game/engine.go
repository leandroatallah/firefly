package game

import (
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/core/scene"
	"github.com/leandroatallah/firefly/internal/systems/audiomanager"
)

type Game struct {
	sceneManager *scene.SceneManager
	state        GameState
	audioManager *audiomanager.AudioManager
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

func (g *Game) SetState(state GameState) *Game {
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

func (g *Game) SetAudioManager() *Game {
	g.audioManager = audiomanager.NewAudioManager()
	g.loadAudioAssets()
	g.sceneManager.SetAudioManager(g.audioManager)
	return g
}

func (g *Game) loadAudioAssets() {
	files, err := os.ReadDir("assets")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if !file.IsDir() && (strings.HasSuffix(file.Name(), ".ogg") || strings.HasSuffix(file.Name(), ".wav")) {
			audioItem, err := g.audioManager.Load("assets/" + file.Name())
			if err != nil {
				log.Fatal(err)
			}
			g.audioManager.Add(audioItem.Name(), audioItem.Data())
		}
	}
}
