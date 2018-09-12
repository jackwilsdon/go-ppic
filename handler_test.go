package ppic_test

import (
	"bytes"
	"github.com/jackwilsdon/go-ppic"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"net/http/httptest"
	"testing"
	"unicode"
)

func isPrintable(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII || !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

func runBenchmark(b *testing.B, path string) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()

		req, err := http.NewRequest(http.MethodGet, path, nil)

		if err != nil {
			b.Errorf("http.NewRequest: %s", err)
			continue
		}

		rec := httptest.NewRecorder()

		b.StartTimer()

		ppic.Handler(rec, req)

		b.StopTimer()

		res := rec.Result()

		if res.StatusCode != http.StatusOK {
			b.Errorf("expected status %d but got %d", http.StatusOK, res.StatusCode)
		}
	}
}

func BenchmarkHandlerPNG(b *testing.B) {
	runBenchmark(b, "/example.png")
}

func BenchmarkHandlerGIF(b *testing.B) {
	runBenchmark(b, "/example.gif")
}

func BenchmarkHandlerJPEG(b *testing.B) {
	runBenchmark(b, "/example.jpeg")
}

func TestHandlerMethod(t *testing.T) {
	methods := []string{
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}

	for _, method := range methods {
		req, err := http.NewRequest(method, "/example", nil)

		if err != nil {
			t.Errorf("http.NewRequest: %s", err)
			continue
		}

		rec := httptest.NewRecorder()

		ppic.Handler(rec, req)

		res := rec.Result()

		if res.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d for %s but got %d", http.StatusMethodNotAllowed, method, res.StatusCode)
		}

		allow := res.Header.Get("Allow")

		if allow != http.MethodGet {
			t.Errorf("expected allow header \"%s\" for %s but got \"%s\"", http.MethodGet, method, allow)
		}
	}
}

func TestHandlerType(t *testing.T) {
	cases := []struct {
		path        string
		contentType string
		format      string
	}{
		{"/example", "image/png", "png"},
		{"/example.png", "image/png", "png"},
		{"/example.gif", "image/gif", "gif"},
		{"/example.jpg", "image/jpeg", "jpeg"},
		{"/example.jpeg", "image/jpeg", "jpeg"},
	}

	for _, c := range cases {
		req, err := http.NewRequest(http.MethodGet, c.path, nil)

		if err != nil {
			t.Errorf("http.NewRequest: %s", err)
			return
		}

		rec := httptest.NewRecorder()

		ppic.Handler(rec, req)

		res := rec.Result()

		if res.StatusCode != http.StatusOK {
			t.Errorf("expected status %d for \"%s\" but got %d", http.StatusOK, c.path, res.StatusCode)
			continue
		}

		// Try and decode the image in the response.
		_, format, err := image.DecodeConfig(res.Body)

		if err != nil {
			t.Errorf("failed to parse image for \"%s\": %s", c.path, err)
		}

		// If there is no error then we need to check the response format.
		if err == nil && format != c.format {
			t.Errorf("expected format \"%s\" for \"%s\" but got \"%s\"", c.format, c.path, format)
		}

		cType := res.Header.Get("Content-Type")

		if cType != c.contentType {
			t.Errorf("expected content type \"%s\" for \"%s\" but got \"%s\"", c.contentType, c.path, cType)
		}
	}
}

func TestHandlerSize(t *testing.T) {
	cases := []struct {
		path       string
		size       int
		statusCode int
		response   string
	}{
		{"/example.", 512, http.StatusNotFound, "error: unsupported file format"},
		{"/example.?size=1024", 1024, http.StatusNotFound, "error: unsupported file format"},
		{"/example.?size=1023", 0, http.StatusNotFound, "error: unsupported file format"},
		{"/example.?size=foo", 0, http.StatusNotFound, "error: unsupported file format"},
		{"/example.foo", 512, http.StatusNotFound, "error: unsupported file format"},
		{"/example.foo?size=1024", 1024, http.StatusNotFound, "error: unsupported file format"},
		{"/example.foo?size=1023", 0, http.StatusNotFound, "error: unsupported file format"},
		{"/example.foo?size=foo", 0, http.StatusNotFound, "error: unsupported file format"},
		{"/example", 512, http.StatusOK, ""},
		{"/example?size=1024", 1024, http.StatusOK, ""},
		{"/example?size=1023", 0, http.StatusBadRequest, "error: size must be a multiple of 8"},
		{"/example?size=foo", 0, http.StatusBadRequest, "error: invalid size"},
		{"/example.png", 512, http.StatusOK, ""},
		{"/example.png?size=1024", 1024, http.StatusOK, ""},
		{"/example.png?size=1023", 0, http.StatusBadRequest, "error: size must be a multiple of 8"},
		{"/example.png?size=foo", 0, http.StatusBadRequest, "error: invalid size"},
		{"/example.gif", 512, http.StatusOK, ""},
		{"/example.gif?size=1024", 1024, http.StatusOK, ""},
		{"/example.gif?size=1023", 0, http.StatusBadRequest, "error: size must be a multiple of 8"},
		{"/example.gif?size=foo", 0, http.StatusBadRequest, "error: invalid size"},
		{"/example.jpg", 512, http.StatusOK, ""},
		{"/example.jpg?size=1024", 1024, http.StatusOK, ""},
		{"/example.jpg?size=1023", 0, http.StatusBadRequest, "error: size must be a multiple of 8"},
		{"/example.jpg?size=foo", 0, http.StatusBadRequest, "error: invalid size"},
		{"/example.jpeg", 512, http.StatusOK, ""},
		{"/example.jpeg?size=1024", 1024, http.StatusOK, ""},
		{"/example.jpeg?size=1023", 0, http.StatusBadRequest, "error: size must be a multiple of 8"},
		{"/example.jpeg?size=foo", 0, http.StatusBadRequest, "error: invalid size"},
	}

	for _, c := range cases {
		req, err := http.NewRequest(http.MethodGet, c.path, nil)

		if err != nil {
			t.Errorf("http.NewRequest: %s", err)
			return
		}

		rec := httptest.NewRecorder()

		ppic.Handler(rec, req)

		res := rec.Result()

		if res.StatusCode != c.statusCode {
			t.Errorf("expected status %d for \"%s\" but got %d", c.statusCode, c.path, res.StatusCode)
		}

		// We're only interested in checking the response if we aren't expecting an OK response.
		if c.statusCode != http.StatusOK {
			buf := bytes.Buffer{}

			if _, err := buf.ReadFrom(res.Body); err != nil {
				t.Errorf("failed to read from response buffer: %s", err)
				continue
			}

			txt := buf.String()

			if txt != c.response {
				// Ensure that we don't print out garbage.
				if !isPrintable(txt) {
					txt = "<binary data>"
				}

				t.Errorf("expected response body \"%s\" for \"%s\" but got \"%s\"", c.response, c.path, txt)
			}
		}
	}
}

func TestHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/jackwilsdon", nil)

	if err != nil {
		t.Errorf("http.NewRequest: %s", err)
		return
	}

	rec := httptest.NewRecorder()

	ppic.Handler(rec, req)

	res := rec.Result()

	// Try and decode the image in the response.
	img, _, err := image.Decode(res.Body)

	if err != nil {
		t.Fatalf("error returned: %s", err)
	}

	// Extract the image dimensions.
	b := img.Bounds()
	w := b.Dx()
	h := b.Dy()

	if w != 512 {
		t.Errorf("expected width to be 512 but got %d", w)
	}

	if h != 512 {
		t.Errorf("expected height to be 512 but got %d", h)
	}

	// We don't want to go any further if any of the dimensions are wrong.
	if t.Failed() {
		return
	}

	// A low resolution version of what we expect the image to look like.
	data := [8][8]color.Color{
		{color.Black, color.White, color.Black, color.White, color.White, color.Black, color.White, color.Black},
		{color.Black, color.Black, color.Black, color.Black, color.Black, color.Black, color.Black, color.Black},
		{color.Black, color.White, color.Black, color.Black, color.Black, color.Black, color.White, color.Black},
		{color.White, color.Black, color.Black, color.Black, color.Black, color.Black, color.Black, color.White},
		{color.White, color.White, color.White, color.White, color.White, color.White, color.White, color.White},
		{color.Black, color.White, color.White, color.White, color.White, color.White, color.White, color.Black},
		{color.Black, color.White, color.Black, color.White, color.White, color.Black, color.White, color.Black},
		{color.White, color.White, color.Black, color.Black, color.Black, color.Black, color.White, color.White},
	}

	// Pixel size is image size / 8 (the size of the actual grid).
	pSize := 512 / 8

	// Compare the image data to the low resolution version.
	for y, row := range data {
		for x, val := range row {
			// Get the color at the corresponding pixel.
			c := img.At(x*pSize, y*pSize)

			// Get the expected and actual RGBA values for the pixel.
			eR, eG, eB, eA := val.RGBA()
			r, g, b, a := c.RGBA()

			// Ensure that everything matches up.
			if eR != r || eG != g || eB != b || eA != a {
				t.Errorf("expected (%d, %d) to be [%d, %d, %d, %d] but got [%d, %d, %d, %d]", x*pSize, y*pSize, eR, eG, eB, eA, r, g, b, a)
			}
		}
	}
}
