package ppic

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"math"
	"math/rand"
	"sync"
)

// ErrInvalidSize is an error caused by specifying a size which is not a multiple of 8.
var ErrInvalidSize = errors.New("size must be a multiple of 8")

// hashString hashes the provided string into an integer.
func hashString(s string) int64 {
	m := sha256.New()

	// Write our string to the MD5 hash calculator.
	fmt.Fprint(m, s)

	// Convert the first 8 bytes into a number.
	return int64(binary.BigEndian.Uint64(m.Sum(nil)))
}

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

// Generate returns an 8x8 grid of values based on the provided source text, optionally mirrored on the X or Y axis.
func Generate(k string, mX, mY bool) [8][8]bool {
	var d [8][8]bool

	w := len(d[0])
	h := len(d)

	// Work out the number of pixels we need.
	pc := w * h

	// If we're mirroring we need half the number of pixels.
	if mX {
		pc /= 2
	}

	// If we're mirroring we need half the number of pixels.
	if mY {
		pc /= 2
	}

	// Work out the number of bytes we need (8 pixels per byte).
	bc := pc / 8

	// Hash the string and create a random number source from it.
	hsh := hashString(k)
	src := rand.NewSource(hsh)
	rnd := rand.New(src)

	// The destination byte array for the random data.
	b := make([]byte, bc)

	// Fill in our buffer with some random data.
	r, err := rnd.Read(b)

	// The default source should never throw an error.
	if err != nil {
		panic(fmt.Sprintf("failed to get random data: %s", err))
	}

	// The default source should always return the requested number of bytes.
	if r != bc {
		panic(fmt.Sprintf("failed to get random data: expected %d bytes but got %d", bc, r))
	}

	// Populate the left half of the image.
	for i := 0; i < pc; i++ {
		// The byte we're interested in.
		bv := b[i/8]

		// The mask for the byte we're interested in.
		mask := byte(math.Pow(2, float64(i%8)))

		// The width we're dividing by.
		cw := w

		if mX {
			cw /= 2
		}

		// Work out the position of the pixel we're looking at.
		x := i % cw
		y := i / cw

		// Set the pixel based on whether or not the bit is set.
		d[y][x] = bv&mask > 0

		// Mirror on the X axis if we need to.
		if mX {
			d[y][w-x-1] = d[y][x]
		}

		// Mirror in the Y axis if we need to.
		if mY {
			d[h-y-1][x] = d[y][x]

			// If we're mirroring on both the X and Y axis then we need to set the bottom right value.
			if mX {
				d[h-y-1][w-x-1] = d[y][x]
			}
		}
	}

	return d
}

// GenerateImage returns an image for the specified source text, optionally mirrored on the X or Y axis.
func GenerateImage(k string, size int, mX, mY bool, p Palette) (image.Image, error) {
	if size%8 != 0 {
		return nil, ErrInvalidSize
	}

	// The size of each pixel in the image.
	pSize := size / 8

	// Create the image and image data.
	img := image.NewPaletted(image.Rect(0, 0, size, size), p.Palette())
	grid := Generate(k, mX, mY)

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
				rect(*img, x, y, pSize, uint8(c))
				wg.Done()
			}(x, y, c)
		}
	}

	// Wait for all of the goroutines to finish.
	wg.Wait()

	return img, nil
}
