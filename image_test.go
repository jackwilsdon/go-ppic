package ppic_test

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/jackwilsdon/go-ppic"
	"github.com/jackwilsdon/go-ppic/ppictest"
)

func BenchmarkGenerateImage(b *testing.B) {
	for n := 0; n < b.N; n++ {
		if _, err := ppic.GenerateImage("jackwilsdon", 512, false, false, ppic.DefaultPalette); err != nil {
			b.Errorf("error: %s", err)
		}
	}
}

func TestGenerateImage(t *testing.T) {
	cases := []struct {
		text    string
		size    int
		mX      bool
		mY      bool
		palette ppic.Palette
		image   [8]string
	}{
		{
			text:    "jackwilsdon",
			size:    512,
			mX:      true,
			palette: ppic.DefaultPalette,
			image: [8]string{
				"# #  # #",
				"########",
				"# #### #",
				" ###### ",
				"        ",
				"#      #",
				"# #  # #",
				"  ####  ",
			},
		},
		{
			text:    "jackwilsdon",
			size:    512,
			mX:      true,
			palette: ppic.Palette{Foreground: color.RGBA{R: 0xFF, A: 0xFF}, Background: color.Black},
			image: [8]string{
				"# #  # #",
				"########",
				"# #### #",
				" ###### ",
				"        ",
				"#      #",
				"# #  # #",
				"  ####  ",
			},
		},
	}

	for i, c := range cases {
		c := c

		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			img, err := ppic.GenerateImage(c.text, c.size, c.mX, c.mY, c.palette)

			if err != nil {
				t.Fatal(err)
			}

			err = ppictest.CompareImage(img, c.palette, c.image)

			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestGenerateImageWithInvalidSize(t *testing.T) {
	_, err := ppic.GenerateImage("jackwilsdon", 31, true, false, ppic.DefaultPalette)

	if err == nil || err != ppic.ErrInvalidSize {
		msg := "nil"

		if err != nil {
			msg = fmt.Sprintf("%q", msg)
		}

		t.Errorf("expected error to be %q but got %s", ppic.ErrInvalidSize, msg)
	}
}
