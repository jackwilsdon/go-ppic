package ppic_test

import (
	"image/color"
	"testing"

	"github.com/jackwilsdon/go-ppic"
)

func colorsEqual(a, b color.Color) bool {
	aR, aG, aB, aA := a.RGBA()
	bR, bG, bB, bA := b.RGBA()

	return aR == bR && aG == bG && aB == bB && aA == bA
}

func TestPalette(t *testing.T) {
	p := ppic.Palette{Foreground: color.Black, Background: color.White}
	pp := p.Palette()

	// Make sure that the first item is the background color.
	if p.Background != pp[0] {
		t.Errorf("expected pp[0] to be %#v but got %#v", p.Background, pp[0])
	}

	// Make sure that the second item is the foreground color.
	if p.Foreground != pp[1] {
		t.Errorf("expected pp[1] to be %#v but got %#v", p.Foreground, pp[1])
	}
}

func TestGeneratePalette(t *testing.T) {
	cases := []struct {
		text    string
		palette ppic.Palette
	}{
		{
			text: "",
			palette: ppic.Palette{
				Foreground: color.RGBA{R: 0xFC, G: 0x1C, B: 0x14, A: 0xFF},
				Background: color.White,
			},
		},
		{
			text: "jackwilsdon",
			palette: ppic.Palette{
				Foreground: color.RGBA{R: 0xEA, G: 0xE3, B: 0xA4, A: 0xFF},
				Background: color.White,
			},
		},
		{
			text: "testing, 123",
			palette: ppic.Palette{
				Foreground: color.RGBA{R: 0xBD, G: 0x3A, B: 0x3B, A: 0xFF},
				Background: color.White,
			},
		},
	}

	for _, c := range cases {
		c := c
		name := c.text

		if len(name) == 0 {
			name = "[empty]"
		}

		t.Run(name, func(t *testing.T) {
			p := ppic.GeneratePalette(c.text)

			// Check the foreground color.
			if !colorsEqual(c.palette.Foreground, p.Foreground) {
				eR, eG, eB, eA := c.palette.Foreground.RGBA()
				aR, aG, aB, aA := p.Foreground.RGBA()

				t.Errorf(
					"expected foreground to be %02X%02X%02X%02X but got %02X%02X%02X%02X",
					uint8(eR),
					uint8(eG),
					uint8(eB),
					uint8(eA),
					uint8(aR),
					uint8(aG),
					uint8(aB),
					uint8(aA),
				)
			}

			// Check the background color.
			if !colorsEqual(c.palette.Background, p.Background) {
				eR, eG, eB, eA := c.palette.Background.RGBA()
				aR, aG, aB, aA := p.Background.RGBA()

				t.Errorf(
					"expected background to be %02X%02X%02X%02X but got %02X%02X%02X%02X",
					uint8(eR),
					uint8(eG),
					uint8(eB),
					uint8(eA),
					uint8(aR),
					uint8(aG),
					uint8(aB),
					uint8(aA),
				)
			}
		})
	}
}
