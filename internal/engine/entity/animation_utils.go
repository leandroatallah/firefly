package entity

import (
	"image"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/render/sprites"
)

type Animatable interface {
	GetSpriteByState(state animation.SpriteState) *sprites.Sprite
	Position() image.Rectangle
	FrameRate() int
}

func IsAnimationFinished(tick int, entity Animatable, state animation.SpriteState) bool {
	if entity == nil {
		return true
	}

	sprite := entity.GetSpriteByState(state)
	if sprite == nil || sprite.Image == nil {
		return true
	}

	rect := entity.Position()
	if rect.Dx() == 0 {
		return true
	}

	elementWidth := sprite.Image.Bounds().Dx()
	frameCount := elementWidth / rect.Dx()

	frameRate := entity.FrameRate()
	if frameRate == 0 {
		frameRate = 1
	}

	duration := frameCount * frameRate
	return tick >= duration
}
