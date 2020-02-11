package ppic_test

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/jackwilsdon/go-ppic"
	"github.com/jackwilsdon/go-ppic/ppictest"
)

func BenchmarkGenerateImage(b *testing.B) {
	grid := ppictest.Parse([8]string{
		"# #  # #",
		"# #### #",
		"        ",
		"# #  # #",
		"  #  #  ",
		"        ",
		"##    ##",
		"#      #",
	})

	for n := 0; n < b.N; n++ {
		if _, err := ppic.GenerateImage(grid, 512, ppic.DefaultPalette); err != nil {
			b.Errorf("error: %s", err)
		}
	}
}

func TestGenerateImage(t *testing.T) {
	cases := []struct {
		grid    [8]string
		size    int
		palette ppic.Palette
	}{
		{
			grid: [8]string{
				"# #  # #",
				"# #### #",
				"        ",
				"# #  # #",
				"  #  #  ",
				"        ",
				"##    ##",
				"#      #",
			},
			size:    512,
			palette: ppic.DefaultPalette,
		},
		{
			grid: [8]string{
				"# #  # #",
				"# #### #",
				"        ",
				"# #  # #",
				"  #  #  ",
				"        ",
				"##    ##",
				"#      #",
			},
			size:    512,
			palette: ppic.Palette{Foreground: color.RGBA{R: 0xFF, A: 0xFF}, Background: color.Black},
		},
	}

	for i, c := range cases {
		c := c

		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			grid := ppictest.Parse(c.grid)
			img, err := ppic.GenerateImage(grid, c.size, c.palette)

			if err != nil {
				t.Fatal(err)
			}

			err = ppictest.CompareImage(img, c.palette, c.grid)

			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestGenerateImageWithInvalidSize(t *testing.T) {
	grid := ppictest.Parse([8]string{
		"# #  # #",
		"# #### #",
		"        ",
		"# #  # #",
		"  #  #  ",
		"        ",
		"##    ##",
		"#      #",
	})

	_, err := ppic.GenerateImage(grid, 31, ppic.DefaultPalette)

	if err == nil || err != ppic.ErrInvalidSize {
		msg := "nil"

		if err != nil {
			msg = fmt.Sprintf("%q", msg)
		}

		t.Errorf("expected error to be %q but got %s", ppic.ErrInvalidSize, msg)
	}
}
