package main

import (
	"gards/src/img"
	"log"
)

func init() {
	err := img.LoadFont("assets/Oswald-Bold.ttf") //wip
	if err != nil {
		log.Fatalf("Failed to load font: %v", err)
	}
}

func main() {
	img.GenerateCard("https://open.spotify.com/track/6XtNoK9E4N4LUtAsAgO5eA?si=ef7691fb65b240a2")
}
