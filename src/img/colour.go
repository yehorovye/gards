package img

import (
	"image"
	"image/color"
	"math"
	"math/rand"
)

type rgb struct {
	r, g, b float64
}

func toRGB(c color.Color) rgb {
	r, g, b, _ := c.RGBA()
	return rgb{
		r: float64(r >> 8),
		g: float64(g >> 8),
		b: float64(b >> 8),
	}
}

func distance(a, b rgb) float64 {
	dr := a.r - b.r
	dg := a.g - b.g
	db := a.b - b.b
	return dr*dr + dg*dg + db*db
}

func meanColor(colors []rgb) rgb {
	var sum rgb
	for _, c := range colors {
		sum.r += c.r
		sum.g += c.g
		sum.b += c.b
	}
	n := float64(len(colors))
	return rgb{sum.r / n, sum.g / n, sum.b / n}
}

func quantizeWithRand(pixels []rgb, k int, rnd *rand.Rand) []rgb {
	centroids := make([]rgb, k)
	for i := range centroids {
		centroids[i] = pixels[rnd.Intn(len(pixels))]
	}

	assignments := make([]int, len(pixels))
	for range 10 {
		for i, p := range pixels {
			minDist := math.MaxFloat64
			for j, c := range centroids {
				d := distance(p, c)
				if d < minDist {
					minDist = d
					assignments[i] = j
				}
			}
		}

		clusters := make([][]rgb, k)
		for i, a := range assignments {
			clusters[a] = append(clusters[a], pixels[i])
		}

		for i := range centroids {
			if len(clusters[i]) > 0 {
				centroids[i] = meanColor(clusters[i])
			}
		}
	}

	return centroids
}

func GetDominantColors(img image.Image, k int) []color.RGBA {
	rnd := rand.New(rand.NewSource(42))

	bounds := img.Bounds()
	var pixels []rgb
	for y := bounds.Min.Y; y < bounds.Max.Y; y += 2 {
		for x := bounds.Min.X; x < bounds.Max.X; x += 2 {
			pixels = append(pixels, toRGB(img.At(x, y)))
		}
	}

	quantized := quantizeWithRand(pixels, k, rnd)

	colors := make([]color.RGBA, len(quantized))
	for i, c := range quantized {
		colors[i] = color.RGBA{
			R: uint8(c.r),
			G: uint8(c.g),
			B: uint8(c.b),
			A: 255,
		}
	}
	return colors
}
