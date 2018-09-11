package ppic

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
)

type surface interface {
	Set(x, y int, c color.Color)
}

// hashString hashes the provided string into an integer.
func hashString(s string) int64 {
	m := sha256.New()

	// Write our string to the MD5 hash calculator.
	fmt.Fprint(m, s)

	// Convert the first 8 bytes into a number.
	return int64(binary.BigEndian.Uint64(m.Sum(nil)))
}

// rect draws a rectangle on the provided surface.
func rect(s surface, r image.Rectangle, c color.Color) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			s.Set(x, y, c)
		}
	}
}

// Generate returns an 8x8 grid of values based on the provided source text, optionally mirrored on the X or Y axis.
func Generate(k string, mX, mY bool) (d [8][8]bool) {
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
	rnd.Read(b)

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

	return
}

// GenerateImage returns an image for the specified source text, optionally mirrored on the X or Y axis.
func GenerateImage(k string, size int, mX, mY bool) (image.Image, error) {
	if size % 8 != 0 {
		return nil, errors.New("size must be a multiple of 8")
	}

	// The size of each pixel in the image.
	pSize := size / 8

	// Create the image and image data.
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	grid := Generate(k, mX, mY)

	// Draw the image data onto the image.
	for y, row := range grid {
		for x, val := range row {
			c := color.White

			// We draw in black if it's set or white if it's unset.
			if val {
				c = color.Black
			}

			// Work out the pixel positions we're drawing at.
			pX := x * pSize
			pY := y * pSize

			// Draw the pixel.
			rect(img, image.Rect(pX, pY, pX+pSize, pY+pSize), c)
		}
	}

	return img, nil
}
