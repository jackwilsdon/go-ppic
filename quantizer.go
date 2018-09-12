package ppic

import (
	"image"
	"image/color"
)

// bwQuantizer is a simple draw.Quantizer which only supports the colors black and white.
type bwQuantizer struct{}

func (bwQuantizer) Quantize(p color.Palette, m image.Image) color.Palette {
	return []color.Color{
		color.White,
		color.Black,
	}
}
