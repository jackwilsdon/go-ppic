package ppic_test

import (
	"github.com/jackwilsdon/go-ppic"
	"image/color"
	"reflect"
	"testing"
)

func BenchmarkGenerate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ppic.Generate("jackwilsdon", false, false)
	}
}

func TestGenerate(t *testing.T) {
	cases := []struct {
		text     string
		mX       bool
		mY       bool
		expected [8][8]bool
	}{
		{
			text: "jackwilsdon",
			expected: [8][8]bool{
				{true, false, true, false, true, true, true, true},
				{true, false, true, true, false, true, true, true},
				{false, false, false, false, true, false, false, false},
				{true, false, true, false, false, false, true, true},
				{false, false, true, false, false, true, true, false},
				{false, false, false, false, false, true, false, true},
				{true, true, false, false, false, true, false, false},
				{true, false, false, false, false, true, true, true},
			},
		},
		{
			text: "jackwilsdon",
			mX:   true,
			expected: [8][8]bool{
				{true, false, true, false, false, true, false, true},
				{true, true, true, true, true, true, true, true},
				{true, false, true, true, true, true, false, true},
				{false, true, true, true, true, true, true, false},
				{false, false, false, false, false, false, false, false},
				{true, false, false, false, false, false, false, true},
				{true, false, true, false, false, true, false, true},
				{false, false, true, true, true, true, false, false},
			},
		},
		{
			text: "jackwilsdon",
			mY:   true,
			expected: [8][8]bool{
				{true, false, true, false, true, true, true, true},
				{true, false, true, true, false, true, true, true},
				{false, false, false, false, true, false, false, false},
				{true, false, true, false, false, false, true, true},
				{true, false, true, false, false, false, true, true},
				{false, false, false, false, true, false, false, false},
				{true, false, true, true, false, true, true, true},
				{true, false, true, false, true, true, true, true},
			},
		},
		{
			text: "jackwilsdon",
			mX:   true,
			mY:   true,
			expected: [8][8]bool{
				{true, false, true, false, false, true, false, true},
				{true, true, true, true, true, true, true, true},
				{true, false, true, true, true, true, false, true},
				{false, true, true, true, true, true, true, false},
				{false, true, true, true, true, true, true, false},
				{true, false, true, true, true, true, false, true},
				{true, true, true, true, true, true, true, true},
				{true, false, true, false, false, true, false, true},
			},
		},
	}

	for i, c := range cases {
		grid := ppic.Generate(c.text, c.mX, c.mY)

		if !reflect.DeepEqual(grid, c.expected) {
			t.Errorf("generated grid does not match test data for case %d", i)
		}
	}
}

func BenchmarkGenerateImage(b *testing.B) {
	for n := 0; n < b.N; n++ {
		if _, err := ppic.GenerateImage("jackwilsdon", 512, false, false, ppic.DefaultPalette); err != nil {
			b.Errorf("error: %s", err)
		}
	}
}

func TestGenerateImage(t *testing.T) {
	red := color.RGBA{R: 0xFF, A: 0xFF}

	cases := []struct {
		text    string
		size    int
		mX      bool
		mY      bool
		palette ppic.Palette
		image   [8][8]color.Color
	}{
		{
			text:    "jackwilsdon",
			size:    512,
			mX:      true,
			palette: ppic.DefaultPalette,
			image: [8][8]color.Color{
				{color.Black, color.White, color.Black, color.White, color.White, color.Black, color.White, color.Black},
				{color.Black, color.Black, color.Black, color.Black, color.Black, color.Black, color.Black, color.Black},
				{color.Black, color.White, color.Black, color.Black, color.Black, color.Black, color.White, color.Black},
				{color.White, color.Black, color.Black, color.Black, color.Black, color.Black, color.Black, color.White},
				{color.White, color.White, color.White, color.White, color.White, color.White, color.White, color.White},
				{color.Black, color.White, color.White, color.White, color.White, color.White, color.White, color.Black},
				{color.Black, color.White, color.Black, color.White, color.White, color.Black, color.White, color.Black},
				{color.White, color.White, color.Black, color.Black, color.Black, color.Black, color.White, color.White},
			},
		},
		{
			text:    "jackwilsdon",
			size:    512,
			mX:      true,
			palette: ppic.Palette{Foreground: red, Background: color.Black},
			image: [8][8]color.Color{
				{red, color.Black, red, color.Black, color.Black, red, color.Black, red},
				{red, red, red, red, red, red, red, red},
				{red, color.Black, red, red, red, red, color.Black, red},
				{color.Black, red, red, red, red, red, red, color.Black},
				{color.Black, color.Black, color.Black, color.Black, color.Black, color.Black, color.Black, color.Black},
				{red, color.Black, color.Black, color.Black, color.Black, color.Black, color.Black, red},
				{red, color.Black, red, color.Black, color.Black, red, color.Black, red},
				{color.Black, color.Black, red, red, red, red, color.Black, color.Black},
			},
		},
	}

	for i, c := range cases {
		img, err := ppic.GenerateImage(c.text, c.size, c.mX, c.mY, c.palette)

		if err != nil {
			t.Errorf("error returned for case %d: %s", i, err)
			continue
		}

		if img == nil {
			t.Errorf("returned image is nil for case %d", i)
			continue
		}

		// Extract the image dimensions.
		b := img.Bounds()
		w := b.Dx()
		h := b.Dy()

		if w != c.size {
			t.Errorf("expected width to be %d but got %d for case %d", c.size, w, i)
		}

		if h != c.size {
			t.Errorf("expected height to be %d but got %d for case %d", c.size, h, i)
		}

		// We don't want to go any further if any of the dimensions are wrong.
		if t.Failed() {
			continue
		}

		// Pixel size is image size / 8 (the size of the actual grid).
		pSize := c.size / 8

		// Compare the image data to the low resolution version.
		for y, row := range c.image {
			for x, val := range row {
				// Get the color at the corresponding pixel.
				c := img.At(x*pSize, y*pSize)

				// Get the expected and actual RGBA values for the pixel.
				eR, eG, eB, eA := val.RGBA()
				r, g, b, a := c.RGBA()

				// Ensure that everything matches up.
				if eR != r || eG != g || eB != b || eA != a {
					t.Errorf("expected (%d, %d) to be [%d, %d, %d, %d] but got [%d, %d, %d, %d] for case %d", x*pSize, y*pSize, eR, eG, eB, eA, r, g, b, a, i)
				}
			}
		}
	}
}

func TestGenerateImageWithInvalidSize(t *testing.T) {
	_, err := ppic.GenerateImage("jackwilsdon", 31, true, false, ppic.DefaultPalette)

	// Make sure we get the right error.
	if err == nil || err != ppic.ErrInvalidSize {
		msg := "nil"

		if err != nil {
			msg = "\"" + err.Error() + "\""
		}

		t.Errorf("expected error to be \"%s\" but got \"%s\"", ppic.ErrInvalidSize, msg)
	}
}
