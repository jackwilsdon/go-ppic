package ppic

import (
	"image"
	"sync"
)

// rect draws a rectangle on the provided surface.
func rect(img image.Paletted, x, y, size int, c uint8) {
	x *= size
	y *= size

	for cY := y; cY < y+size; cY++ {
		for cX := x; cX < x+size; cX++ {
			img.Pix[(cY-img.Rect.Min.Y)*img.Stride+(cX-img.Rect.Min.X)] = c
		}
	}
}

// GenerateImage returns an image for the specified grid.
func GenerateImage(grid [8][8]bool, size int, p Palette) (image.Image, error) {
	if size%8 != 0 {
		return nil, ErrInvalidSize
	}

	// The size of each pixel in the image.
	pSize := size / 8

	// Create the image and image data.
	img := image.NewPaletted(image.Rect(0, 0, size, size), p.Palette())

	// Create a wait group so we can wait for all of our goroutines to finish.
	wg := sync.WaitGroup{}

	// There are going to be 64 goroutines (it's an 8x8 image).
	wg.Add(64)

	// Draw the image data onto the image.
	for y, row := range grid {
		for x, val := range row {
			var c uint8

			// We're drawing the foreground if the value is set.
			if val {
				c = 1
			}

			// Draw the pixel.
			go func(x, y int, c uint8) {
				rect(*img, x, y, pSize, c)
				wg.Done()
			}(x, y, c)
		}
	}

	// Wait for all of the goroutines to finish.
	wg.Wait()

	return img, nil
}
