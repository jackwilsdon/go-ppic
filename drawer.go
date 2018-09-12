package ppic

import (
	"image"
	"image/color"
	"image/draw"
)

// bwDrawer is a simple draw.Drawer which draws black and white images.
// This should be used alongside the bwQuantizer to ensure that the only
// colors being drawn are black or white.
// Note that this drawer will only accept an image.Gray and will panic
// in any other cases.
type bwDrawer struct{}

func (bwDrawer) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {
	img := src.(*image.Gray)
	col := color.RGBA{A: 0xFF}

	for y := 0; y != r.Dy(); y++ {
		for x := 0; x != r.Dx(); x++ {
			c := img.Pix[(y-img.Rect.Min.Y)*img.Stride+(x-img.Rect.Min.X)]

			col.R = c
			col.G = c
			col.B = c

			dst.Set(r.Min.X+x, r.Min.Y+y, &col)
		}
	}
}
