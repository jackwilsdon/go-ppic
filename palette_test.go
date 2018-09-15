package ppic_test

import (
	"github.com/jackwilsdon/go-ppic"
	"image/color"
	"testing"
)

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
