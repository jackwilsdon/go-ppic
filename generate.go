package ppic

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
)

// ErrInvalidSize is an error caused by specifying a size which is not a multiple of 8.
var ErrInvalidSize = errors.New("size must be a multiple of 8")

// Generate returns an 8x8 grid of values based on the provided source text, optionally mirrored along the X or Y axis.
func Generate(k string, mX, mY bool) (img [8][8]bool) {
	// Hash the string and create a random number source from it.
	hsh := hashString(k)
	src := rand.NewSource(hsh)
	rnd := rand.New(src)

	// Create a buffer and fill it with random data.
	buf := make([]byte, 8)
	n, err := rnd.Read(buf)

	// The default source should never throw an error.
	if err != nil {
		panic(fmt.Sprintf("failed to get random data: %s", err))
	}

	// The default source should always return the requested number of bytes.
	if n != len(buf) {
		panic(fmt.Sprintf("failed to get random data: expected %d bytes but got %d", len(buf), n))
	}

	for i := 0; i < 64; i++ {
		// Work out which bit of the current byte we're interested in.
		mask := byte(math.Pow(2, float64(i%8)))

		// Work out the position of the pixel we're looking at.
		x := i % 8
		y := i / 8

		// Set the pixel based on whether or not the bit is set.
		img[y][x] = buf[i/8]&mask > 0

		// If we're mirroring along the Y axis and past the center line then draw the mirrored pixels.
		if mX && x >= 4 {
			img[y][x] = img[y][7-x]
		}

		// If we're mirroring along the X axis and past the center line then draw the mirrored pixels.
		if mY && y >= 4 {
			img[y][x] = img[7-y][x]
		}
	}

	return img
}
