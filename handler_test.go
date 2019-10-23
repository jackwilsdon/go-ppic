package ppic_test

import (
	"bytes"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"
	"unicode"

	"github.com/jackwilsdon/go-ppic"
	"github.com/jackwilsdon/go-ppic/ppictest"
)

func isPrintable(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII || !unicode.IsPrint(r) {
			return false
		}
	}

	return true
}

func BenchmarkHandler(b *testing.B) {
	paths := []string{
		"/example.png",
		"/example.gif",
		"/example.jpg",
	}

	for _, p := range paths {
		p := p
		ext := path.Ext(p)

		b.Run(strings.ToUpper(ext[1:]), func(b *testing.B) {
			b.StopTimer()

			for n := 0; n < b.N; n++ {
				req, err := http.NewRequest(http.MethodGet, p, nil)

				if err != nil {
					b.Fatalf("http.NewRequest: %s", err)
				}

				rec := httptest.NewRecorder()

				b.StartTimer()

				ppic.Handler(rec, req)

				b.StopTimer()

				res := rec.Result()

				if res.StatusCode != http.StatusOK {
					b.Fatalf("expected status to be %d but got %d", http.StatusOK, res.StatusCode)
				}
			}
		})
	}
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
		method := method

		t.Run(method, func(t *testing.T) {
			req, err := http.NewRequest(method, "/example", nil)

			if err != nil {
				t.Fatalf("http.NewRequest: %s", err)
			}

			rec := httptest.NewRecorder()

			ppic.Handler(rec, req)

			res := rec.Result()

			if res.StatusCode != http.StatusMethodNotAllowed {
				t.Errorf("expected status to be %d but got %d", http.StatusMethodNotAllowed, res.StatusCode)
			}

			allow := res.Header.Get("Allow")

			if allow != http.MethodGet {
				t.Errorf("expected allow header to be %q but got %q", http.MethodGet, allow)
			}
		})
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
		c := c

		t.Run(c.path[1:], func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, c.path, nil)

			if err != nil {
				t.Fatalf("http.NewRequest: %s", err)
			}

			rec := httptest.NewRecorder()

			ppic.Handler(rec, req)

			res := rec.Result()

			if res.StatusCode != http.StatusOK {
				t.Errorf("expected status to be %d but got %d", http.StatusOK, res.StatusCode)
			}

			// Try and decode the image in the response.
			_, format, err := image.DecodeConfig(res.Body)

			if err != nil {
				t.Fatalf("failed to parse image: %s", err)
			}

			// If there is no error then we need to check the response format.
			if format != c.format {
				t.Errorf("expected format to be %q but got %q", c.format, format)
			}

			cType := res.Header.Get("Content-Type")

			if cType != c.contentType {
				t.Errorf("expected content type to be %q but got %q", c.contentType, cType)
			}
		})
	}
}

func TestHandlerSize(t *testing.T) {
	cases := []struct {
		path       string
		size       int
		statusCode int
		response   string
	}{
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
		c := c

		t.Run(c.path[1:], func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, c.path, nil)

			if err != nil {
				t.Fatalf("http.NewRequest: %s", err)
			}

			rec := httptest.NewRecorder()

			ppic.Handler(rec, req)

			res := rec.Result()

			if res.StatusCode != c.statusCode {
				t.Errorf("expected status %d but got %d", c.statusCode, res.StatusCode)
			}

			// If we're expecting a valid image then check its size.
			if c.statusCode == http.StatusOK {
				img, _, err := image.Decode(res.Body)

				if err != nil {
					t.Fatalf("failed to parse image: %s", err)
				}

				w := img.Bounds().Dx()
				h := img.Bounds().Dy()

				if w != c.size || h != c.size {
					t.Errorf("expected image to be %dx%d but got %dx%d", c.size, c.size, w, h)
				}
			} else {
				// Otherwise check the response text.
				buf := bytes.Buffer{}

				if _, err := buf.ReadFrom(res.Body); err != nil {
					t.Fatalf("failed to read from response buffer: %s", err)
				}

				txt := buf.String()

				if txt != c.response {
					// Ensure that we don't print out garbage.
					if !isPrintable(txt) {
						txt = "<binary data>"
					}

					t.Errorf("expected response body %q but got %q", c.response, txt)
				}
			}
		})
	}
}

func TestHandlerMirrorError(t *testing.T) {
	cases := []struct {
		path     string
		response string
	}{
		{"/example?mirror=a", "error: unsupported mirror axis: a"},
		{"/example?mirror=ax", "error: unsupported mirror axis: a"},
		{"/example?mirror=ay", "error: unsupported mirror axis: a"},
		{"/example?mirror=xa", "error: unsupported mirror axis: a"},
		{"/example?mirror=ya", "error: unsupported mirror axis: a"},
		{"/example?mirror=xx", "error: duplicate mirror axis: x"},
		{"/example?mirror=yy", "error: duplicate mirror axis: y"},
		{"/example?mirror=xay", "error: unsupported mirror axis: a"},
		{"/example?mirror=yax", "error: unsupported mirror axis: a"},
		{"/example?mirror=xya", "error: unsupported mirror axis: a"},
		{"/example?mirror=axy", "error: unsupported mirror axis: a"},
		{"/example?mirror=yxa", "error: unsupported mirror axis: a"},
		{"/example?mirror=ayx", "error: unsupported mirror axis: a"},
		{"/example?mirror=xxy", "error: duplicate mirror axis: x"},
		{"/example?mirror=xyy", "error: duplicate mirror axis: y"},
		{"/example?mirror=yxx", "error: duplicate mirror axis: x"},
		{"/example?mirror=yyx", "error: duplicate mirror axis: y"},
	}

	for _, c := range cases {
		c := c

		t.Run(c.path[1:], func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, c.path, nil)

			if err != nil {
				t.Fatalf("http.NewRequest: %s", err)
			}

			rec := httptest.NewRecorder()

			ppic.Handler(rec, req)

			res := rec.Result()

			if res.StatusCode != http.StatusBadRequest {
				t.Fatalf("expected status %d but got %d", http.StatusBadRequest, res.StatusCode)
			}

			buf := bytes.Buffer{}

			if _, err := buf.ReadFrom(res.Body); err != nil {
				t.Fatalf("failed to read from response buffer: %s", err)
			}

			txt := buf.String()

			if txt != c.response {
				// Ensure that we don't print out garbage.
				if !isPrintable(txt) {
					txt = "<binary data>"
				}

				t.Errorf("expected response body %q but got %q", c.response, txt)
			}
		})
	}
}

func TestHandler(t *testing.T) {
	cases := []struct {
		path    string
		palette ppic.Palette
		image   [8]string
	}{
		{
			path: "/jackwilsdon",
			palette: ppic.Palette{
				Foreground: color.RGBA{R: 0xEA, G: 0xE3, B: 0xA4, A: 0xFF},
				Background: color.White,
			},
			image: [8]string{
				"# #  # #",
				"# #### #",
				"        ",
				"# #  # #",
				"  #  #  ",
				"        ",
				"##    ##",
				"#      #",
			},
		},
		{
			path:    "/jackwilsdon?monochrome",
			palette: ppic.DefaultPalette,
			image: [8]string{
				"# #  # #",
				"# #### #",
				"        ",
				"# #  # #",
				"  #  #  ",
				"        ",
				"##    ##",
				"#      #",
			},
		},
		{
			path: "/jackwilsdon?mirror=xy",
			palette: ppic.Palette{
				Foreground: color.RGBA{R: 0xEA, G: 0xE3, B: 0xA4, A: 0xFF},
				Background: color.White,
			},
			image: [8]string{
				"# #  # #",
				"# #### #",
				"        ",
				"# #  # #",
				"# #  # #",
				"        ",
				"# #### #",
				"# #  # #",
			},
		},
		{
			path:    "/testing123?monochrome",
			palette: ppic.DefaultPalette,
			image: [8]string{
				"  ####  ",
				" ###### ",
				"#  ##  #",
				"        ",
				"        ",
				"##    ##",
				"  ####  ",
				" ##  ## ",
			},
		},
	}

	for _, c := range cases {
		c := c

		t.Run(c.path[1:], func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, c.path, nil)

			if err != nil {
				t.Fatalf("http.NewRequest: %s", err)
			}

			rec := httptest.NewRecorder()

			ppic.Handler(rec, req)

			res := rec.Result()

			// Try and decode the image in the response.
			img, _, err := image.Decode(res.Body)

			if err != nil {
				t.Fatalf("failed to parse image: %s", err)
			}

			err = ppictest.CompareImage(img, c.palette, c.image)

			if err != nil {
				t.Error(err)
			}
		})
	}
}
