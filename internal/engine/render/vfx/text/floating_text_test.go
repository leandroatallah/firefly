package text_test

import (
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/assets/font"
	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	"github.com/boilerplate/ebiten-template/internal/engine/render/camera"
	"github.com/boilerplate/ebiten-template/internal/engine/render/vfx/text"
	"github.com/hajimehoshi/ebiten/v2"
)

// getModuleRoot returns the absolute path to the module root
func getModuleRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			panic("could not find go.mod")
		}
		dir = parent
	}
}

// mockFloatingText implements text.FloatingText for testing
type mockFloatingText struct {
	*text.SimpleFloatingText
	drawn bool
}

func (m *mockFloatingText) Draw(screen *ebiten.Image, cam *camera.Controller) {
	m.drawn = true
}

func newMockText(msg string, duration int) *mockFloatingText {
	return &mockFloatingText{
		SimpleFloatingText: text.NewFloatingText(msg, 0, 0, duration),
	}
}

func TestFloatingTextBase_Update_DecrementsDuration(t *testing.T) {
	ft := newMockText("Test", 10)

	ft.Update()

	if ft.Duration != 9 {
		t.Errorf("expected Duration=9, got %d", ft.Duration)
	}
}

func TestFloatingTextBase_Update_MarksRemoved(t *testing.T) {
	ft := newMockText("Test", 1)

	ft.Update()

	if !ft.IsComplete() {
		t.Error("expected text to be complete after duration reaches 0")
	}
}

func TestFloatingTextBase_IsComplete(t *testing.T) {
	ft := newMockText("Test", 5)

	if ft.IsComplete() {
		t.Error("expected text to not be complete initially")
	}

	// Update until complete
	for ft.Duration > 0 {
		ft.Update()
	}

	if !ft.IsComplete() {
		t.Error("expected text to be complete after duration expires")
	}
}

func TestFloatingTextBase_SetFont(t *testing.T) {
	ft := newMockText("Test", 10)
	ft.SetFont(nil)

	if ft.Font != nil {
		t.Error("expected font to be nil")
	}
}

func TestFloatingTextBase_SetColor(t *testing.T) {
	ft := newMockText("Test", 10)
	testColor := color.RGBA{255, 0, 0, 255}

	ft.SetColor(testColor)

	if ft.Color != testColor {
		t.Errorf("expected color %+v, got %+v", testColor, ft.Color)
	}
}

func TestNewFloatingTextWithVelocity(t *testing.T) {
	ft := text.NewFloatingTextWithVelocity("msg", 10, 20, 30, -1.5)
	if ft.VelocityY != -1.5 {
		t.Errorf("expected VelocityY -1.5, got %f", ft.VelocityY)
	}
}

func TestFloatingTextBase_DrawText_WithCamera(t *testing.T) {
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240})
	cam := camera.NewController(0, 0)
	cam.SetCenter(100, 100)

	ft := text.NewFloatingText("test", 100, 100, 10)

	// Load font from embedded FS
	moduleRoot := getModuleRoot()
	f, err := font.NewFontText(os.DirFS(moduleRoot), "assets/fonts/monogram.ttf")
	if err == nil {
		ft.SetFont(f)
	} else {
		// Fallback to testing without a font if assets are not reachable
		ft.SetFont(nil)
	}

	screen := ebiten.NewImage(320, 240)
	ft.DrawText(screen, 100, 100, cam)

	// Test without camera
	ft.DrawText(screen, 100, 100, nil)
}

func TestManager_NewManager(t *testing.T) {
	m := text.NewManager()
	if m == nil {
		t.Fatal("expected manager to be created")
	}
}

func TestManager_Add(t *testing.T) {
	m := text.NewManager()
	ft := newMockText("Test", 10)
	m.Add(ft)
}

func TestManager_Add_Nil(t *testing.T) {
	m := text.NewManager()
	m.Add(nil)
}

func TestManager_Add_AppliesDefaultFont(t *testing.T) {
	m := text.NewManager()
	var mockFont *font.FontText
	m.SetDefaultFont(mockFont)

	ft := newMockText("Test", 10)
	m.Add(ft)

	if ft.Font != mockFont {
		t.Error("expected default font to be applied")
	}
}

func TestManager_Update_RemovesCompleted(t *testing.T) {
	m := text.NewManager()

	ft1 := newMockText("Test1", 1)
	ft2 := newMockText("Test2", 10)

	m.Add(ft1)
	m.Add(ft2)
	m.Update()

	if !ft1.IsComplete() {
		t.Error("expected ft1 to be complete")
	}
	if ft2.IsComplete() {
		t.Error("expected ft2 to not be complete")
	}
}

func TestManager_Update_MultipleUpdates(t *testing.T) {
	m := text.NewManager()
	ft := newMockText("Test", 5)
	m.Add(ft)

	for i := 0; i < 4; i++ {
		m.Update()
	}

	if ft.IsComplete() {
		t.Error("expected text to still be active after 4 updates")
	}

	m.Update()

	if !ft.IsComplete() {
		t.Error("expected text to be complete after 5 updates")
	}
}

func TestManager_Draw(t *testing.T) {
	m := text.NewManager()
	ft := newMockText("Test", 10)
	ft.Color = color.White
	m.Add(ft)

	screen := ebiten.NewImage(100, 100)
	m.Draw(screen, nil)

	if !ft.drawn {
		t.Error("expected text to be drawn")
	}
}

func TestManager_Draw_OnlyActiveTexts(t *testing.T) {
	m := text.NewManager()

	ft1 := newMockText("Test1", 1)
	ft2 := newMockText("Test2", 100)

	m.Add(ft1)
	m.Add(ft2)
	m.Update()

	screen := ebiten.NewImage(100, 100)
	m.Draw(screen, nil)

	if !ft2.drawn {
		t.Error("expected ft2 to be drawn")
	}
	if ft1.drawn {
		t.Error("expected ft1 to not be drawn")
	}
}

func TestManager_Update_EmptyManager(t *testing.T) {
	m := text.NewManager()
	m.Update()
}

func TestManager_Draw_EmptyManager(t *testing.T) {
	m := text.NewManager()
	screen := ebiten.NewImage(100, 100)
	m.Draw(screen, nil)
}
