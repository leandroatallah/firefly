package main

import (
	"embed"
	"flag"
	"log"

	"github.com/boilerplate/ebiten-template/internal/engine/data/config"
	_ "github.com/boilerplate/ebiten-template/internal/engine/entity/actors" // Blank import to ensure init() is called
	gamesetup "github.com/boilerplate/ebiten-template/internal/game/app"
	_ "github.com/boilerplate/ebiten-template/internal/game/entity/actors/states" // Blank import to ensure init() is called
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
