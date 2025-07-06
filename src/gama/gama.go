// Thanks to n128 (https://github.com/nicolito128) for the original code.
// It can be found at https://github.com/nicolito128/gama.
// All of the code below has been extracted from the repo mentioned above.
// Credits to the original author.

package gama

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"runtime"
	"sync"
)

type Palette interface {
	Quantify(n int) ([]color.Color, error)
}

type paletteImpl struct {
	src      image.Image
	measures []color.Color
	wg       sync.WaitGroup
	mu       sync.Mutex
}

func New(source image.Image) Palette {
	plt := &paletteImpl{
		src: source,
	}
	return plt
}

func (pl *paletteImpl) Quantify(n int) ([]color.Color, error) {
	if n <= 0 {
		return nil, errors.New("palette length must be a non-zero positive value")
	}
	pl.measures = make([]color.Color, 0, n)

	maxParalellism := min(n, runtime.GOMAXPROCS(0), runtime.NumCPU())
	if maxParalellism <= 0 { // Just in case of something weird, lmao
		maxParalellism = 1
	}
	semaphore := make(chan struct{}, maxParalellism)

	width := pl.src.Bounds().Dx()
	height := pl.src.Bounds().Dy()
	numLines := height / n

	if n > height {
		return nil, errors.New("palette length must not exceeds the image height")
	}

	for i := range n {
		pl.wg.Add(1)
		semaphore <- struct{}{}

		go func(job int) {
			defer pl.wg.Done()
			defer func() {
				<-semaphore
			}()

			startY := job * numLines
			endY := startY + numLines
			if job == (n - 1) {
				endY = height
			}

			bucket := NewBucket()
			for y := startY; y < endY; y++ {
				for x := range width {
					bucket.Push(pl.src.At(x, y))
				}
			}
			m := bucket.Median()

			pl.mu.Lock()
			pl.measures = append(pl.measures, m)
			pl.mu.Unlock()
		}(i)
	}
	pl.wg.Wait()

	return pl.measures, nil
}

func ColorToHex(c color.Color, includeAlpha bool) string {
	r, g, b, a := c.RGBA()

	r8 := uint8(r >> 8)
	g8 := uint8(g >> 8)
	b8 := uint8(b >> 8)
	a8 := uint8(a >> 8)

	if includeAlpha {
		return fmt.Sprintf("#%02X%02X%02X%02X", r8, g8, b8, a8)
	}
	return fmt.Sprintf("#%02X%02X%02X", r8, g8, b8)
}
