package ppic

import (
	"bytes"
	"fmt"
	"image/color"

	svg "github.com/ajstarks/svgo"
)

func colorToHex(c color.Color) string {
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	return fmt.Sprintf("#%02x%02x%02x", rgba.R, rgba.G, rgba.B)
}

// GenerateSVG generates an SVG for the specified image.
func GenerateSVG(grid [8][8]bool, p Palette) (string, error) {
	buf := &bytes.Buffer{}
	s := svg.New(buf)

	// A custom version of s.Start which sets a viewBox without a width or height.
	buf.WriteString(`<?xml version="1.0"?>
<svg viewBox="0 0 8 8"
     shape-rendering="crispEdges"
     xmlns="http://www.w3.org/2000/svg"
     xmlns:xlink="http://www.w3.org/1999/xlink">
`)

	s.Rect(0, 0, 8, 8, fmt.Sprintf("fill: %s", colorToHex(p.Background)))

	for y, row := range grid {
		for x, val := range row {
			if !val {
				continue
			}

			s.Rect(x, y, 1, 1, fmt.Sprintf("fill: %s", colorToHex(p.Foreground)))
		}
	}

	s.End()

	return buf.String(), nil
}
