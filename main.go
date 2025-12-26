package main

import (
	"embed"
	"flag"
	"log"

	gamesetup "github.com/leandroatallah/firefly/internal/game/app"
)

//go:embed assets/*
var embedFs embed.FS

func main() {
	flag.Parse()
	err := gamesetup.Setup(embedFs)
	if err != nil {
		log.Fatal(err)
	}
}
