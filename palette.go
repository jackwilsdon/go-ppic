package ppic

import "image/color"

// Palette represents a pair of colors to use in image generation.
type Palette struct {
	Foreground color.Color
	Background color.Color
}

// Palette returns a color palette with the first color being the background and the second being the foreground.
func (p Palette) Palette() color.Palette {
	return color.Palette{p.Background, p.Foreground}
}

// DefaultPalette is the default black and white color palette.
var DefaultPalette = Palette{Foreground: color.Black, Background: color.White}
