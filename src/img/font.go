package img

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"

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
		Dot:  fixed.P(x, y),
	}
	drawer.DrawString(text)

	return nil
}
