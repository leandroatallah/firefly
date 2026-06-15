// Package phaseskit defines genre constants for phase scenes.
// The engine declares the dumb Genre int type; this package names the values.
package phaseskit

import "github.com/boilerplate/ebiten-template/internal/engine/scene/phases"

const (
	// GenrePlatformer identifies a side-scrolling platformer phase.
	GenrePlatformer phases.Genre = iota + 1
	// GenreBeatemup identifies a beat-em-up phase with altitude-aware draw order.
	GenreBeatemup
	// GenreShepherd identifies a shepherd phase (platformer + rescue-sheep goal).
	GenreShepherd
)
