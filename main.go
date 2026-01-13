package main

import (
	"embed"
	"flag"
	"log"

	"github.com/leandroatallah/firefly/internal/engine/data/config"
	gamesetup "github.com/leandroatallah/firefly/internal/game/app"
)

//go:embed assets/*
var embedFs embed.FS

func main() {
	cfg := gamesetup.NewConfig()
	flag.Parse()
	config.Set(cfg)

	err := gamesetup.Setup(embedFs)
	if err != nil {
		log.Fatal(err)
	}
}
