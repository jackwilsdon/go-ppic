package ppic

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
)

// imageWriter represents a function which can write an image to a writer.
type imageWriter func(io.Writer, image.Image) error

// getImageSize extracts an image size from a set of URL values.
func getImageSize(q url.Values) (int, error) {
	ss := q.Get("size")

	if len(ss) == 0 {
		return 512, nil
	}

	s, err := strconv.Atoi(ss)

	if err != nil {
		return 0, err
	}

	return s, nil
}

// getMirroring extracts mirroring axes from a set of URL values.
// Axes are returned in the order (X, Y), with X being set by default if no
// mirroring parameters are specified.
func getMirroring(q url.Values) (bool, bool, error) {
	ms, ok := q["mirror"]

	if !ok || len(ms) == 0 {
		return true, false, nil
	}

	m := ms[0]

	mX := false
	mY := false

	for _, c := range m {
		if (c == 'x' && mX) || (c == 'y' && mY) {
			return false, false, fmt.Errorf("duplicate mirror axis: %c", c)
		}

		switch {
		case c == 'x':
			mX = true
		case c == 'y':
			mY = true
		default:
			return false, false, fmt.Errorf("unsupported mirror axis: %c", c)
		}
	}

	return mX, mY, nil
}

// getImageWriter returns an imageWriter for the specified path.
func getImageWriter(p string) imageWriter {
	ext := path.Ext(p)

	switch strings.ToLower(ext) {
	case ".gif":
		return func(w io.Writer, i image.Image) error {
			return gif.Encode(w, i, &gif.Options{NumColors: 2})
		}
	case ".jpg", ".jpeg":
		return func(w io.Writer, i image.Image) error {
			return jpeg.Encode(w, i, &jpeg.Options{Quality: 100})
		}
	case "", ".png":
		return func(w io.Writer, i image.Image) error {
			enc := png.Encoder{CompressionLevel: png.NoCompression}

			return enc.Encode(w, i)
		}
	default:
		return nil
	}
}

// Handler serves HTTP requests with generated images.
func Handler(res http.ResponseWriter, req *http.Request) {
	// We only support GETing images.
	if req.Method != http.MethodGet {
		res.Header().Set("Allow", http.MethodGet)
		res.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	writer := getImageWriter(req.URL.Path)

	// If we couldn't find a writer then we couldn't understand the extension.
	if writer == nil {
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(res, "error: unsupported file format")

		return
	}

	q := req.URL.Query()

	// Get the image size from the request.
	size, err := getImageSize(q)

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(res, "error: invalid size")

		return
	}

	// Get the mirroring axes from the request.
	mX, mY, err := getMirroring(q)

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(res, "error: %s", err)
		return
	}

	// Get the path without extension.
	txt := strings.TrimSuffix(req.URL.Path[1:], path.Ext(req.URL.Path))

	pal := DefaultPalette

	// Generate a palette based on the source text if we're not in monochrome mode.
	if _, mono := q["monochrome"]; !mono {
		pal = GeneratePalette(txt)
	}

	// Generate the grid.
	grid := Generate(txt, mX, mY)

	// Generate the image.
	img, err := GenerateImage(grid, size, pal)

	// Check if an invalid size was specified.
	if err == ErrInvalidSize {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(res, "error: %s", err)

		return
	}

	// Check if something else bad happened during generation.
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(res, "error: %s", err)

		return
	}

	// Write the image to the response.
	if err = writer(res, img); err != nil {
		fmt.Fprintf(res, "error: %s", err)
	}
}
