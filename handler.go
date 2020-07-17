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

// getImageWriter returns an imageWriter for the specified extension.
func getImageWriter(ext string) imageWriter {
	switch ext {
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

	q := req.URL.Query()

	// Get the path without extension.
	ext := path.Ext(req.URL.Path)
	txt := strings.TrimSuffix(req.URL.Path[1:], ext)

	pal := DefaultPalette

	// Generate a palette based on the source text if we're not in monochrome mode.
	if _, mono := q["monochrome"]; !mono {
		pal = GeneratePalette(txt)
	}

	// Generate the grid.
	grid := Generate(txt, true, false)

	// Lower-case the extension to make comparisons case-insensitive.
	ext = strings.ToLower(ext)

	if ext == ".svg" {
		img := GenerateSVG(grid, pal)

		res.Write([]byte(img))

		return
	}

	writer := getImageWriter(ext)

	if writer != nil {
		// Get the image size from the request.
		size, err := getImageSize(q)

		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(res, "error: invalid size")

			return
		}

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

		return
	}

	res.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(res, "error: unsupported file format")
}
