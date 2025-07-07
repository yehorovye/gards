package img

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

var loadedFont *opentype.Font

func LoadFont(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read font file: %w", err)
	}
	loadedFont, err = opentype.Parse(data)
	if err != nil {
		return fmt.Errorf("failed to parse font: %w", err)
	}
	return nil
}

func EstimateLineCount(text string, maxCharsPerLine int) int {
	words := strings.Fields(text)
	if len(words) == 0 {
		return 0
	}

	lines := 0
	current := ""

	for _, word := range words {
		if len(current)+len(word)+1 <= maxCharsPerLine {
			if current == "" {
				current = word
			} else {
				current += " " + word
			}
		} else {
			lines++
			current = word
		}
	}
	if current != "" {
		lines++
	}
	return lines
}

func RenderText(dst draw.Image, text string, x, y int, col color.Color, size float64) error {
	if loadedFont == nil {
		return fmt.Errorf("font not loaded; call font.LoadFont() first")
	}

	face, err := opentype.NewFace(loadedFont, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return fmt.Errorf("failed to create font face: %w", err)
	}
	defer face.Close()

	drawer := &font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(col),
		Face: face,
	}

	const maxLineLength = 20
	lineHeight := int(size * 1.4)

	words := strings.Fields(text)
	var lines []string
	var current string

	for _, word := range words {
		if len(current)+len(word)+1 <= maxLineLength {
			if current == "" {
				current = word
			} else {
				current += " " + word
			}
		} else {
			lines = append(lines, current)
			current = word
		}
	}
	if current != "" {
		lines = append(lines, current)
	}

	for i, line := range lines {
		drawer.Dot = fixed.P(x, y+i*lineHeight)
		drawer.DrawString(line)
	}

	return nil
}
