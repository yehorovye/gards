package main

import (
	"gards/src/img"
	"log"
)

func init() {
	err := img.LoadFont("assets/inter.ttf") //wip
	if err != nil {
		log.Fatalf("Failed to load font: %v", err)
	}
}

func main() {
	img.GenerateCard("wip :D")
}
