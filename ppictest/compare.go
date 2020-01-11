package ppictest

import (
	"encoding/xml"
	"fmt"
	"image"
	"image/color"

	"github.com/jackwilsdon/go-ppic"
)

// Validate the "expected" string for an 8x8 image.
func validateExpected(expected [8]string) {
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
}

// Compare an 8x8 grid to an expected image.
//
// Expected image must be 8 lines, each consisting of 8 characters.
func Compare(grid [8][8]bool, expected [8]string) error {
	validateExpected(expected)

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

// CompareImage compares an image with an expected image.
//
// Expected image must be 8 lines, each consisting of 8 characters.
func CompareImage(img image.Image, expectedPal ppic.Palette, expected [8]string) error {
	if expectedPal.Foreground == nil {
		panic("expectedPal.Foreground is nil")
	}

	if expectedPal.Background == nil {
		panic("expectedPal.Background is nil")
	}

	validateExpected(expected)

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
					x/ps,
					y/ps,
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

func flattenSVG(src string, p ppic.Palette) (grid [8][8]bool, err error) {
	svg := struct {
		XMLName xml.Name `xml:"svg"`
		ViewBox string   `xml:"viewBox,attr"`
		Rects   []struct {
			XMLName xml.Name `xml:"rect"`
			X       int      `xml:"x,attr"`
			Y       int      `xml:"y,attr"`
			Width   int      `xml:"width,attr"`
			Height  int      `xml:"height,attr"`
			Style   string   `xml:"style,attr"`
		} `xml:",any"`
	}{}

	if err = xml.Unmarshal([]byte(src), &svg); err != nil {
		err = fmt.Errorf("failed to parse SVG: %s", err)
		return
	}

	if svg.ViewBox != "0 0 8 8" {
		err = fmt.Errorf("expected viewBox to be \"0 0 8 8\" but got %q", svg.ViewBox)
		return
	}

	for i, rect := range svg.Rects {
		if rect.X < 0 || rect.X > 7 {
			err = fmt.Errorf("expected rect %d to have an x value in the range 0 to 7, got %d", i, rect.X)
			return
		}

		if rect.Y < 0 || rect.Y > 7 {
			err = fmt.Errorf("expected rect %d to have a y value in the range 0 to 7, got %d", i, rect.Y)
			return
		}

		if rect.Width < 0 || rect.X+rect.Width > 8 {
			err = fmt.Errorf("expected rect %d to have a width in the range 0 to %d, got %d", i, 8-rect.X, rect.Width)
			return
		}

		if rect.Height < 0 || rect.Y+rect.Height > 8 {
			err = fmt.Errorf("expected rect %d to have a height in the range 0 to %d, got %d", i, 8-rect.Y, rect.Height)
			return
		}

		fill := color.RGBA{A: 0xff}

		if _, err = fmt.Sscanf(rect.Style, "fill: #%02x%02x%02x", &fill.R, &fill.G, &fill.B); err != nil {
			err = fmt.Errorf("invalid style for rect %d: %s", i, rect.Style)
			return
		}

		fR, fG, fB, fA := fill.RGBA()
		fgR, fgG, fgB, fgA := p.Foreground.RGBA()
		bgR, bgG, bgB, bgA := p.Background.RGBA()

		if (fR != fgR || fG != fgG || fB != fgB || fA != fgA) && (fR != bgR || fG != bgG || fB != bgB || fA != bgA) {
			err = fmt.Errorf(
				"invalid color for rect %d: #%02X%02X%02X%02X (expected one of #%02X%02X%02X%02X, #%02X%02X%02X%02X)",
				i,
				fR,
				fG,
				fB,
				fA,
				fgR,
				fgG,
				fgB,
				fgA,
				bgR,
				bgG,
				bgB,
				bgA,
			)
			return
		}

		for y := rect.Y; y < rect.Y+rect.Height; y++ {
			for x := rect.X; x < rect.X+rect.Width; x++ {
				grid[y][x] = fill == p.Foreground
			}
		}
	}

	return
}

// CompareSVG compares an svg with an expected image.
//
// Expected image must be 8 lines, each consisting of 8 characters.
func CompareSVG(src string, expectedPal ppic.Palette, expected [8]string) error {
	if expectedPal.Foreground == nil {
		panic("expectedPal.Foreground is nil")
	}

	if expectedPal.Background == nil {
		panic("expectedPal.Background is nil")
	}

	validateExpected(expected)

	grid, err := flattenSVG(src, expectedPal)

	if err != nil {
		return fmt.Errorf("failed to flatten SVG: %s", err)
	}

	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			exp := expected[y][x] == '#'
			value := grid[y][x]

			if value != exp {
				expColor := expectedPal.Background
				color := expectedPal.Background

				if exp {
					expColor = expectedPal.Foreground
				}

				if value {
					color = expectedPal.Foreground
				}

				eR, eG, eB, eA := expColor.RGBA()
				r, g, b, a := color.RGBA()

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
