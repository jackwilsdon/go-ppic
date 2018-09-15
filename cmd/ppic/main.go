package main

import (
	"fmt"
	"github.com/jackwilsdon/go-ppic"
	"image/png"
	"os"
	"path"
	"strconv"
	"strings"
)

var cmd = path.Base(os.Args[0])

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s text [size] > image.png\n", cmd)
		os.Exit(1)
	}

	txt := os.Args[1]
	size := 512

	if len(os.Args) > 2 {
		var err error

		size, err = strconv.Atoi(os.Args[2])

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: invalid size \"%s\"\n", cmd, os.Args[2])
			os.Exit(1)
		}
	}

	// If we're trying to output to a terminal then prevent it.
	if isTerminal() {
		fmt.Fprintf(os.Stderr, "%s: refusing to output image to stdout (it looks like a terminal!)\n", cmd)

		args := strings.Join(os.Args[1:], " ")
		fmt.Fprintf(os.Stderr, "\ntry piping the output to a file:\n\t%s %s > image.png\n", cmd, args)

		os.Exit(1)
	}

	img, err := ppic.GenerateImage(txt, size, true, false, ppic.DefaultPalette)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: failed to generate image: %s\n", cmd, err)
		os.Exit(1)
	}

	err = png.Encode(os.Stdout, img)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: failed to generate PNG: %s\n", cmd, err)
		os.Exit(1)
	}
}
