package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
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

func (g *Game) SetSceneManager(manager *scene.SceneManager) *Game {
	g.sceneManager = manager
	return g
}

func (g *Game) SetAudioContext(ctx *audio.Context) *Game {
	g.sceneManager.SetAudioContext(ctx)
	return g
}

func (g *Game) SetAudioAssets(items []*audioplayer.AudioItem) *Game {
	stream := make(map[string][]byte)
	for _, i := range items {
		stream[i.Name()] = i.Data()
	}
	g.sceneManager.SetAudioStream(stream)
	return g
}
