package main

import (
	"embed"

	gamesetup "github.com/leandroatallah/firefly/internal/game/setup"
)

//go:embed assets/*
var embedFs embed.FS

func main() {
	gamesetup.Setup(embedFs)
}
