package img

import (
	"gards/src/spotify"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"os"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

var (
	client_id     = os.Getenv("SPOTIFY_CLIENT_ID")
	client_secret = os.Getenv("SPOTIFY_CLIENT_SECRET")
)

func GenerateCard(url string) {
	track := fetchTrack(url)
	cover := parseImage(track.Images[0].URL)
	card := buildCard(track, cover)
	saveCard(card, "result.png")
}

func fetchTrack(url string) *spotify.Track {
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
		log.Fatalf("Expected at least 2 images, got %d", len(track.Images)) // future stuff
	}

	return track
}

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

func buildCard(track *spotify.Track, cover image.Image) image.Image {
	var (
		offset              = 25
		coverSize           = 500
		titleMaxLineLength  = 32
		artistMaxLineLength = 36
		titleFontSize       = 38.0
		artistFontSize      = 20.0
	)

	titleLineSpacing := int(titleFontSize * 1)
	artistLineSpacing := int(artistFontSize * 1.4)

	artists := strings.Join(track.Artists, ", ")

	nameLineCount := EstimateLineCount(track.Name, titleMaxLineLength)
	artistLineCount := EstimateLineCount(artists, artistMaxLineLength)
	totalTextHeight := nameLineCount*titleLineSpacing + artistLineCount*artistLineSpacing

	width := coverSize + offset*2
	height := coverSize + offset*5 + totalTextHeight

	card := image.NewRGBA(image.Rect(0, 0, width, height))

	colors := GetDominantColors(cover, 6)

	drawBackground(card, color.Black)
	drawCover(card, cover, offset, coverSize)

	drawColorBars(card, colors, offset, coverSize, titleLineSpacing)

	drawText(card, track.Name, artists, offset, coverSize, nameLineCount, titleLineSpacing, titleFontSize, artistFontSize)
	drawBottomLine(card, colors[0], height, width)

	return card
}

func drawColorBars(img *image.RGBA, colors []color.RGBA, offset, coverSize, titleSpacing int) {
	rectWidth := coverSize / 6
	rectHeight := 20
	gap := 0

	rectY := offset + coverSize + titleSpacing - rectHeight - 5

	for i, c := range colors {
		x := offset + i*(rectWidth+gap)
		rect := image.Rect(x, rectY, x+rectWidth, rectY+rectHeight)
		draw.Draw(img, rect, &image.Uniform{c}, image.Point{}, draw.Over)
	}
}

func drawBackground(img *image.RGBA, bg color.Color) {
	draw.Draw(img, img.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)
}

func drawCover(img *image.RGBA, cover image.Image, offset, size int) {
	draw.Draw(img, image.Rect(offset, offset, offset+size, offset+size), cover, image.Point{}, draw.Over)
}

func drawText(img *image.RGBA, title, artists string, offset, coverSize, titleLines, titleSpacing int, titleSize, artistSize float64) {
	textY := (offset * 3) + coverSize + titleSpacing
	if err := RenderText(img, title, offset, textY, color.White, titleSize); err != nil {
		log.Fatalf("Failed to render track name: %v", err)
	}

	textY += titleLines * titleSpacing
	if err := RenderText(img, artists, offset, textY, color.White, artistSize); err != nil {
		log.Fatalf("Failed to render artists: %v", err)
	}
}

func drawBottomLine(img *image.RGBA, col color.Color, height, width int) {
	rect := image.Rect(0, height-15, width, height)
	draw.Draw(img, rect, &image.Uniform{col}, image.Point{}, draw.Over)
}

func saveCard(card image.Image, path string) {
	outFile, err := os.Create(path)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer outFile.Close()

	if err := png.Encode(outFile, card); err != nil {
		log.Fatalf("Failed to encode image: %v", err)
	}
}
