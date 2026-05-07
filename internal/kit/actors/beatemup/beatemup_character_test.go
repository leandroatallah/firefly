package beatemup_test

import (
	"testing"

	"github.com/boilerplate/ebiten-template/internal/kit/actors/beatemup"
)

func TestNewBeatEmUpCharacter_NotNil(t *testing.T) {
	c := beatemup.NewBeatEmUpCharacter()
	if c == nil {
		t.Fatal("NewBeatEmUpCharacter returned nil")
	}
	if c.MeleeCharacter == nil {
		t.Fatal("BeatEmUpCharacter.MeleeCharacter is nil; expected embedded MeleeCharacter to be initialised")
	}
}
