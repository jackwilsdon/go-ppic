package main

import (
	"fmt"
	"github.com/jackwilsdon/go-ppic"
	"image/png"
	"os"
	"path"
	"strconv"
)

func handleError(err error) {
	if err == nil {
		return
	}

	fmt.Fprintf(os.Stderr, "%s: %s\n", path.Base(os.Args[0]), err)
	os.Exit(1)
}

func main() {
	cmd := path.Base(os.Args[0])

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s text [size] > image.png\n", cmd)
		os.Exit(1)
	}

	txt := os.Args[1]
	size := 512

	if len(os.Args) > 2 {
		var err error

		size, err = strconv.Atoi(os.Args[2])

		handleError(err)
	}

	img, err := ppic.GenerateImage(txt, size, true, false)

	handleError(err)

	err = png.Encode(os.Stdout, img)

	handleError(err)
}
