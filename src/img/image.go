package img

import (
	"gards/src/gama"
	"gards/src/spotify"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"os"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

var (
	client_id     = os.Getenv("SPOTIFY_CLIENT_ID")
	client_secret = os.Getenv("SPOTIFY_CLIENT_SECRET")
)

func parseImage(img string) image.Image {
	res, err := http.Get(img)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	src, _, err := image.Decode(res.Body)
	if err != nil {
		panic(err)
	}

	return src
}

func GenerateCard(url string) {
	token, err := spotify.GetClientCredentialsToken(client_id, client_secret)
	if err != nil {
		log.Fatalf("Failed to get Spotify token: %v", err)
	}

	client := spotify.New(token)
	track, err := client.GetTrackFromURL(url)
	if err != nil {
		log.Fatalf("Failed to get track from URL: %v", err)
	}

	if len(track.Images) < 2 {
		log.Fatalf("Expected at least 2 images, got %d", len(track.Images))
	}

	cover := parseImage(track.Images[1].URL) // 300x300 (provided by the spotify api)

	const offset = 15
	width := 300 + (offset * 2)
	height := width + 100

	card := image.NewRGBA(image.Rect(0, 0, width, height))

	palette := gama.New(cover)
	colors, err := palette.Quantify(1)
	if err != nil || len(colors) == 0 {
		log.Fatalf("Failed to extract dominant color: %v", err)
	}

	avg := colors[0].(color.RGBA)

	var factor uint8 = 80

	bg := color.RGBA{
		R: avg.R + factor,
		G: avg.G + factor,
		B: avg.B + factor,
		A: 255,
	}

	draw.Draw(card, card.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)
	draw.Draw(card, image.Rect(offset, offset, offset+300, offset+300), cover, image.Point{}, draw.Over)

	textColor := color.Black
	RenderText(card, track.Name, offset, offset+300+20, textColor, 1.5)

	// TODO: remove this later.
	outFile, err := os.Create("result.png")
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer outFile.Close()

	if err := png.Encode(outFile, card); err != nil {
		log.Fatalf("Failed to encode image: %v", err)
	}
}
