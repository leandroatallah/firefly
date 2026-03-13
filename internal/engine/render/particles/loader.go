package particles

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/leandroatallah/firefly/internal/engine/data/schemas"
)

// LoadConfig loads a particle configuration from a JSON file.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var pData schemas.ParticleData
	if err := json.Unmarshal(data, &pData); err != nil {
		return nil, err
	}

	// Resolve image path relative to the JSON file
	imagePath := filepath.Join(filepath.Dir(path), pData.Image)
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		// Try resolving relative to the working directory (project root)
		if _, err := os.Stat(pData.Image); err == nil {
			imagePath = pData.Image
		}
	}

	img, _, err := ebitenutil.NewImageFromFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load particle image %s: %w", imagePath, err)
	}

	frameCount := 1
	if pData.FrameWidth > 0 {
		frameCount = img.Bounds().Dx() / pData.FrameWidth
	}

	return &Config{
		Image:       img,
		FrameWidth:  pData.FrameWidth,
		FrameHeight: pData.FrameHeight,
		FrameCount:  frameCount,
		FrameRate:   pData.FrameRate,
	}, nil
}
