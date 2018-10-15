// +build go1.9
// +build ignore

package gziphandler

import "net/http"

const (
	// DefaultQValue is the default qvalue to assign to an encoding if no explicit qvalue is set.
	// This is actually kind of ambiguous in RFC 2616, so hopefully it's correct.
	// The examples seem to indicate that it is.
	//
	// Deprecated: maintained for compatibility, no longer used.
	DefaultQValue = 1.0

	// DefaultMinSize defines the minimum size to reach to enable compression.
	// It's 1400 bytes.
	//
	// Deprecated: maintained for compatibility.
	DefaultMinSize = defaultMinSize
)

// DeprecatedHandler is a type alias to to group deprecated
// functions together for godoc.
//
// This will be removed once godoc automatically groups
// deprecated symbols, see golang/go#17056.
type DeprecatedHandler = http.Handler

// DeprecatedMiddleware is a type alias to to group
// deprecated functions together for godoc.
//
// This will be removed once godoc automatically groups
// deprecated symbols, see golang/go#17056.
type DeprecatedMiddleware = func(http.Handler) http.Handler

// GzipHandler wraps an HTTP handler, to transparently gzip
// the response body if the client supports it (via the
// Accept-Encoding header). This will compress at the
// default compression level.
//
// Deprecated: maintained for compatibility, use Gzip.
func GzipHandler(h http.Handler) DeprecatedHandler {
	return Gzip(h)
}

// GzipHandlerWithOpts ...
//
// Deprecated: maintained for compatibility, use Wrapper.
func GzipHandlerWithOpts(opts ...option) (DeprecatedMiddleware, error) {
	return Wrapper(opts), nil
}

// MustNewGzipLevelHandler behaves just like
// NewGzipLevelHandler except that in an error case it
// panics rather than returning an error.
//
// Deprecated: maintained for compatibility, use Wrapper.
func MustNewGzipLevelHandler(level int) DeprecatedMiddleware {
	return Wrapper(CompressionLevel(level))
}

// NewGzipLevelAndMinSize behave as NewGzipLevelHandler
// except it let the caller specify the minimum size before
// compression.
//
// Deprecated: maintained for compatibility, use Wrapper.
func NewGzipLevelAndMinSize(level, minSize int) (DeprecatedMiddleware, error) {
	return Wrapper(CompressionLevel(level), MinSize(minSize)), nil
}

// NewGzipLevelHandler returns a wrapper function (often
// known as middleware) which can be used to wrap an HTTP
// handler to transparently gzip the response body if the
// client supports it (via the Accept-Encoding header).
// Responses will be encoded at the given gzip compression
// level. An error will be returned only if an invalid gzip
// compression level is given, so if one can ensure the
// level is valid, the returned error can be safely ignored.
//
// Deprecated: maintained for compatibility, use Wrapper.
func NewGzipLevelHandler(level int) (DeprecatedMiddleware, error) {
	return Wrapper(CompressionLevel(level)), nil
}

/*type GzipResponseWriter
  func (w *GzipResponseWriter) Close() error
  func (w *GzipResponseWriter) Flush()
  func (w *GzipResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error)
  func (w *GzipResponseWriter) Push(target string, opts *http.PushOptions) error
  func (w *GzipResponseWriter) Write(b []byte) (int, error)
  func (w *GzipResponseWriter) WriteHeader(code int)
*/
