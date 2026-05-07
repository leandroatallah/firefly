package beatemup

import kitactors "github.com/boilerplate/ebiten-template/internal/kit/actors"

type BeatEmUpCharacter struct {
	*kitactors.MeleeCharacter
}

func NewBeatEmUpCharacter() *BeatEmUpCharacter {
	return &BeatEmUpCharacter{
		MeleeCharacter: kitactors.NewMeleeCharacter(),
	}
}
