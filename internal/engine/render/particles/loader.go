package particles

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/boilerplate/ebiten-template/internal/engine/data/schemas"
)

// LoadConfig loads a particle configuration from a JSON file.
func LoadConfig(fsys fs.FS, path string) (*Config, error) {
	data, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, err
	}

	var pData schemas.ParticleData
	if err := json.Unmarshal(data, &pData); err != nil {
		return nil, err
	}

	// Resolve image path relative to the JSON file
	imagePath := pData.Image
	if !filepath.IsAbs(imagePath) {
		// Try loading from the same directory as the JSON file
		imagePath = filepath.Join(filepath.Dir(path), pData.Image)
	}

	imgData, err := fs.ReadFile(fsys, imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load particle image %s: %w", imagePath, err)
	}

	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(imgData))
	if err != nil {
		return nil, fmt.Errorf("failed to parse particle image %s: %w", imagePath, err)
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
