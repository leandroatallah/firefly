package main

import (
	"embed"
	"log"

	gamesetup "github.com/leandroatallah/firefly/internal/game/setup"
)

//go:embed assets/*
var embedFs embed.FS

func main() {
	err := gamesetup.Setup(embedFs)
	if err != nil {
		log.Fatal(err)
	}
}
