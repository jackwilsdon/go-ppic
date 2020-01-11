package ppic_test

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/jackwilsdon/go-ppic"
	"github.com/jackwilsdon/go-ppic/ppictest"
)

func BenchmarkGenerateSVG(b *testing.B) {
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
		ppic.GenerateSVG(grid, ppic.DefaultPalette)
	}
}

func TestGenerateSVG(t *testing.T) {
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
			palette: ppic.Palette{Foreground: color.RGBA{R: 0xFF, A: 0xFF}, Background: color.Black},
		},
	}

	for i, c := range cases {
		c := c

		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			grid := ppictest.Parse(c.grid)
			svg := ppic.GenerateSVG(grid, c.palette)

			if err := ppictest.CompareSVG(svg, c.palette, c.grid); err != nil {
				t.Error(err)
			}
		})
	}
}
