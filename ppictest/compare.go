package ppictest

import (
	"fmt"
	"github.com/jackwilsdon/go-ppic"
	"image"
)

func Compare(grid [8][8]bool, expected [8]string) error {
	for y, row := range expected {
		if l := len(row); l != 8 {
			panic(fmt.Sprintf("len(expected[%d]) != 8 (got %d)", y, l))
		}

		for x, c := range row {
			if c != '#' && c != ' ' {
				panic(fmt.Sprintf("expected[%d][%d] is not '#' or ' ' (got %q)", y, x, c))
			}
		}
	}

	for y := range grid {
		for x := range grid[y] {
			exp := expected[y][x] == '#'

			if act := grid[y][x]; act != exp {
				return fmt.Errorf("expected grid[%d][%d] to be %t but got %t", y, x, exp, act)
			}
		}
	}

	return nil
}

func CompareImage(img image.Image, expectedPal ppic.Palette, expected [8]string) error {
	if expectedPal.Foreground == nil {
		panic("expectedPal.Foreground is nil")
	}

	if expectedPal.Background == nil {
		panic("expectedPal.Background is nil")
	}

	for y, row := range expected {
		if l := len(row); l != 8 {
			panic(fmt.Sprintf("len(expected[%d]) != 8 (got %d)", y, l))
		}

		for x, c := range row {
			if c != '#' && c != ' ' {
				panic(fmt.Sprintf("expected[%d][%d] is not '#' or ' ' (got %q)", y, x, c))
			}
		}
	}

	if img == nil {
		return fmt.Errorf("image is nil")
	}

	b := img.Bounds()
	w, h := b.Dx(), b.Dy()

	if w%8 != 0 {
		return fmt.Errorf("expected width to be divisible by 8 but got %d", w)
	}

	if h%8 != 0 {
		return fmt.Errorf("expected height to be divisible by 8 but got %d", h)
	}

	if w != h {
		return fmt.Errorf("expected width to be equal to height (got width %d, height %d)", w, h)
	}

	ps := w / 8

	// Loop through each pixel of the source image.
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// Find the character in the source text for this pixel.
			expChar := expected[y/ps][x/ps]
			expColor := expectedPal.Background

			// Use the foreground if this pixel is '#'.
			if expChar == '#' {
				expColor = expectedPal.Foreground
			}

			// Get the expected and actual RGBA values for the pixel.
			eR, eG, eB, eA := expColor.RGBA()
			r, g, b, a := img.At(x, y).RGBA()

			// Ensure that everything matches up.
			if eR != r || eG != g || eB != b || eA != a {
				return fmt.Errorf(
					"expected foreground at (%d, %d) to be #%02X%02X%02X%02X but got #%02X%02X%02X%02X",
					x,
					y,
					uint8(eR),
					uint8(eG),
					uint8(eB),
					uint8(eA),
					uint8(r),
					uint8(g),
					uint8(b),
					uint8(a),
				)
			}
		}
	}

	return nil
}
