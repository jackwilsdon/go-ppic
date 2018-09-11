package ppic_test

import (
	"github.com/jackwilsdon/go-ppic"
	"image/color"
	"reflect"
	"testing"
)

func TestGenerate(t *testing.T) {
	grid := ppic.Generate("jackwilsdon", false, false)

	if !reflect.DeepEqual(grid, [8][8]bool{
		{true, false, true, false, true, true, true, true},
		{true, false, true, true, false, true, true, true},
		{false, false, false, false, true, false, false, false},
		{true, false, true, false, false, false, true, true},
		{false, false, true, false, false, true, true, false},
		{false, false, false, false, false, true, false, true},
		{true, true, false, false, false, true, false, false},
		{true, false, false, false, false, true, true, true},
	}) {
		t.Error("generated grid does not match test data")
	}

	grid = ppic.Generate("jackwilsdon", true, false)

	if !reflect.DeepEqual(grid, [8][8]bool{
		{true, false, true, false, false, true, false, true},
		{true, true, true, true, true, true, true, true},
		{true, false, true, true, true, true, false, true},
		{false, true, true, true, true, true, true, false},
		{false, false, false, false, false, false, false, false},
		{true, false, false, false, false, false, false, true},
		{true, false, true, false, false, true, false, true},
		{false, false, true, true, true, true, false, false},
	}) {
		t.Error("generated grid does not match test data")
	}

	grid = ppic.Generate("jackwilsdon", false, true)

	if !reflect.DeepEqual(grid, [8][8]bool{
		{true, false, true, false, true, true, true, true},
		{true, false, true, true, false, true, true, true},
		{false, false, false, false, true, false, false, false},
		{true, false, true, false, false, false, true, true},
		{true, false, true, false, false, false, true, true},
		{false, false, false, false, true, false, false, false},
		{true, false, true, true, false, true, true, true},
		{true, false, true, false, true, true, true, true},
	}) {
		t.Error("generated grid does not match test data")
	}

	grid = ppic.Generate("jackwilsdon", true, true)

	if !reflect.DeepEqual(grid, [8][8]bool{
		{true, false, true, false, false, true, false, true},
		{true, true, true, true, true, true, true, true},
		{true, false, true, true, true, true, false, true},
		{false, true, true, true, true, true, true, false},
		{false, true, true, true, true, true, true, false},
		{true, false, true, true, true, true, false, true},
		{true, true, true, true, true, true, true, true},
		{true, false, true, false, false, true, false, true},
	}) {
		t.Error("generated grid does not match test data")
	}
}

func TestGenerateImage(t *testing.T) {
	img, err := ppic.GenerateImage("jackwilsdon", 512, true, false)

	if err != nil {
		t.Fatalf("error returned: %s", err)
	}

	if img == nil {
		t.Fatal("returned image is nil")
	}

	// Extract the image dimensions.
	b := img.Bounds()
	w := b.Dx()
	h := b.Dy()

	if w != 512 {
		t.Errorf("expected width to be 512 but got %d", w)
	}

	if h != 512 {
		t.Errorf("expected height to be 512 but got %d", h)
	}

	// We don't want to go any further if any of the dimensions are wrong.
	if t.Failed() {
		return
	}

	// A low resolution version of what we expect the image to look like.
	data := [8][8]color.Color{
		{color.Black, color.White, color.Black, color.White, color.White, color.Black, color.White, color.Black},
		{color.Black, color.Black, color.Black, color.Black, color.Black, color.Black, color.Black, color.Black},
		{color.Black, color.White, color.Black, color.Black, color.Black, color.Black, color.White, color.Black},
		{color.White, color.Black, color.Black, color.Black, color.Black, color.Black, color.Black, color.White},
		{color.White, color.White, color.White, color.White, color.White, color.White, color.White, color.White},
		{color.Black, color.White, color.White, color.White, color.White, color.White, color.White, color.Black},
		{color.Black, color.White, color.Black, color.White, color.White, color.Black, color.White, color.Black},
		{color.White, color.White, color.Black, color.Black, color.Black, color.Black, color.White, color.White},
	}

	// Pixel size is image size / 8 (the size of the actual grid).
	pSize := 512 / 8

	// Compare the image data to the low resolution version.
	for y, row := range data {
		for x, val := range row {
			// Get the color at the corresponding pixel.
			c := img.At(x*pSize, y*pSize)

			// Get the expected and actual RGBA values for the pixel.
			eR, eG, eB, eA := val.RGBA()
			r, g, b, a := c.RGBA()

			// Ensure that everything matches up.
			if eR != r || eG != g || eB != b || eA != a {
				t.Errorf("expected (%d, %d) to be [%d, %d, %d, %d] but got [%d, %d, %d, %d]", x*pSize, y*pSize, eR, eG, eB, eA, r, g, b, a)
			}
		}
	}
}

func TestGenerateImageWithInvalidSize(t *testing.T) {
	_, err := ppic.GenerateImage("jackwilsdon", 31, true, false)

	// Make sure we get the right error.
	if err == nil || err.Error() != "size must be a multiple of 8" {
		msg := "nil"

		if err != nil {
			msg = "\"" + err.Error() + "\""
		}

		t.Errorf("expected error to be \"size must be a multiple of 8\" but instead got %s", msg)
	}
}
