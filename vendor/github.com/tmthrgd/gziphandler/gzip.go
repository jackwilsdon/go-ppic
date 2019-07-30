// Package gziphandler is a tiny Go package which wraps
// HTTP handlers to transparently gzip the response body,
// for clients which support it.
package gziphandler

import (
	"bufio"
	"compress/gzip"
	"io"
	"net"
	"net/http"
	"sync"

	"github.com/tmthrgd/httputils"
)

// 1500 bytes is the MTU size for the internet since that is
// the largest size allowed at the network layer. If you
// take a file that is 1300 bytes and compress it to 800
// bytes, it is still transmitted in that same 1500 byte
// packet regardless, so you’ve gained nothing. That being
// the case, you should restrict the gzip compression to
// files with a size greater than a single packet, 1400
// bytes is a safe value.
const defaultMinSize = 1400

var bufferPool = &sync.Pool{
	New: func() interface{} {
		buf := make([]byte, 0, defaultMinSize)
		return &buf
	},
}

var gzipWriterPools [gzip.BestCompression - gzip.HuffmanOnly + 1]sync.Pool

func gzipWriterPool(level int) *sync.Pool {
	return &gzipWriterPools[level-gzip.HuffmanOnly]
}

func gzipWriterGet(w io.Writer, level int) *gzip.Writer {
	if gw, ok := gzipWriterPool(level).Get().(*gzip.Writer); ok {
		gw.Reset(w)
		return gw
	}

	gw, _ := gzip.NewWriterLevel(w, level)
	return gw
}

func gzipWriterPut(gw *gzip.Writer, level int) {
	gzipWriterPool(level).Put(gw)
}

// These constants are copied from the gzip package, so
// that code that imports "github.com/tmthrgd/gziphandler"
// does not also have to import "compress/gzip".
const (
	NoCompression      = gzip.NoCompression
	BestSpeed          = gzip.BestSpeed
	BestCompression    = gzip.BestCompression
	DefaultCompression = gzip.DefaultCompression
	HuffmanOnly        = gzip.HuffmanOnly
)

// responseWriter provides an http.ResponseWriter interface,
// which gzips bytes before writing them to the underlying
// response. This doesn't close the writers, so don't forget
// to do that. It can be configured to skip response smaller
// than minSize.
type responseWriter struct {
	http.ResponseWriter

	h *handler

	gw *gzip.Writer

	// Holds the first part of the write before reaching
	// the minSize or the end of the write.
	buf *[]byte

	// Saves the WriteHeader value.
	code int
}

// WriteHeader just saves the response code until close or
// GZIP effective writes.
func (w *responseWriter) WriteHeader(code int) {
	if w.code == 0 {
		w.code = code
	}
}

// Write appends data to the gzip writer.
func (w *responseWriter) Write(b []byte) (int, error) {
	switch {
	case w.buf != nil && w.gw != nil:
		panic("gziphandler: both buf and gw are non nil in call to Write")
	// GZIP responseWriter is initialized. Use the GZIP
	// responseWriter.
	case w.gw != nil:
		return w.gw.Write(b)
	// We're operating in pass through mode.
	case w.buf == nil:
		return w.ResponseWriter.Write(b)
	}

	w.WriteHeader(http.StatusOK)

	// This may succeed if the Content-Type header was
	// explicitly set.
	if w.shouldPassThrough() {
		if err := w.startPassThrough(); err != nil {
			return 0, err
		}

		return w.ResponseWriter.Write(b)
	}

	if w.shouldBuffer(b) {
		// Save the write into a buffer for later.
		// This buffer will be flushed in either
		// startGzip or startPassThrough.
		*w.buf = append(*w.buf, b...)
		return len(b), nil
	}

	w.inferContentType(b)

	// Now that we've called inferContentType, we have
	// a Content-Type header.
	if w.shouldPassThrough() {
		if err := w.startPassThrough(); err != nil {
			return 0, err
		}

		return w.ResponseWriter.Write(b)
	}

	if err := w.startGzip(); err != nil {
		return 0, err
	}

	return w.gw.Write(b)
}

// startGzip initialize any GZIP specific informations.
func (w *responseWriter) startGzip() (err error) {
	h := w.Header()

	// Set the GZIP header.
	h.Set("Content-Encoding", "gzip")

	// if the Content-Length is already set, then calls
	// to Write on gzip will fail to set the
	// Content-Length header since its already set
	// See: https://github.com/golang/go/issues/14975.
	h.Del("Content-Length")

	// Write the header to gzip response.
	w.ResponseWriter.WriteHeader(w.code)

	// Bytes written during ServeHTTP are redirected to
	// this gzip writer before being written to the
	// underlying response.
	w.gw = gzipWriterGet(w.ResponseWriter, w.h.level)

	if buf := *w.buf; len(buf) != 0 {
		// Flush the buffer into the gzip response.
		_, err = w.gw.Write(buf)
	}

	w.releaseBuffer()
	return err
}

func (w *responseWriter) startPassThrough() (err error) {
	w.ResponseWriter.WriteHeader(w.code)

	if buf := *w.buf; len(buf) != 0 {
		_, err = w.ResponseWriter.Write(buf)
	}

	w.releaseBuffer()
	return err
}

func (w *responseWriter) releaseBuffer() {
	if w.buf == nil {
		panic("gziphandler: w.buf is nil in call to emptyBuffer")
	}

	*w.buf = (*w.buf)[:0]
	bufferPool.Put(w.buf)
	w.buf = nil
}

func (w *responseWriter) shouldBuffer(b []byte) bool {
	// If the all writes to date are bigger than the
	// minSize, we no longer need to buffer and we can
	// decide whether to enable compression or whether
	// to operate in pass through mode.
	return len(*w.buf)+len(b) < w.h.minSize
}

func (w *responseWriter) inferContentType(b []byte) {
	h := w.Header()

	// If content type is not set.
	if _, ok := h["Content-Type"]; ok {
		return
	}

	if buf := *w.buf; len(buf) != 0 {
		const sniffLen = 512
		if len(buf) >= sniffLen {
			b = buf
		} else if len(buf)+len(b) > sniffLen {
			b = append(buf, b[:sniffLen-len(buf)]...)
		} else {
			b = append(buf, b...)
		}
	}

	if len(b) == 0 {
		return
	}

	// It infer it from the uncompressed body.
	h.Set("Content-Type", http.DetectContentType(b))
}

func (w *responseWriter) shouldPassThrough() bool {
	if w.Header().Get("Content-Encoding") != "" {
		return true
	}

	return !w.handleContentType()
}

func (w *responseWriter) handleContentType() bool {
	// If contentTypes is empty, accept any content
	// type.
	if len(w.h.contentTypes) == 0 {
		return true
	}

	// If the Content-Type header is not set, return
	// as we haven't called inferContentType yet.
	ct, ok := w.Header()["Content-Type"]
	if !ok {
		return true
	}

	if len(ct) == 0 {
		return false
	}

	return httputils.MIMETypeMatches(ct[0], w.h.contentTypes)
}

// Close will close the gzip.Writer and will put it back in
// the gzipWriterPool.
func (w *responseWriter) Close() error {
	switch {
	case w.buf != nil && w.gw != nil:
		panic("gziphandler: both buf and gw are non nil in call to Close")
	// Buffer not nil means the regular response must
	// be returned.
	case w.buf != nil:
		return w.closeNonGzipped()
	// If the GZIP responseWriter is not set no need
	// to close it.
	case w.gw != nil:
		return w.closeGzipped()
	// Both buf and gw nil means we are operating in
	// pass through mode.
	default:
		return nil
	}
}

func (w *responseWriter) closeGzipped() error {
	err := w.gw.Close()

	gzipWriterPut(w.gw, w.h.level)
	w.gw = nil

	return err
}

func (w *responseWriter) closeNonGzipped() error {
	w.inferContentType(nil)

	w.WriteHeader(http.StatusOK)

	return w.startPassThrough()
}

// Flush flushes the underlying *gzip.Writer and then the
// underlying http.ResponseWriter if it is an http.Flusher.
// This makes responseWriter an http.Flusher.
func (w *responseWriter) Flush() {
	if w.gw == nil && w.buf != nil {
		// Fix for NYTimes/gziphandler#58:
		//  Only flush once startGzip or
		//  startPassThrough has been called.
		//
		// Flush is thus a no-op until the written
		// body exceeds minSize, or we've decided
		// not to compress.
		return
	}

	if w.gw != nil {
		w.gw.Flush()
	}

	if fw, ok := w.ResponseWriter.(http.Flusher); ok {
		fw.Flush()
	}
}

type handler struct {
	h http.Handler
	config
}

func (h *handler) shouldGzip(r *http.Request) bool {
	if h.config.shouldGzip != nil {
		switch h.config.shouldGzip(r) {
		case NegotiateGzip:
		case SkipGzip:
			return false
		case ForceGzip:
			return true
		}
	}

	match := httputils.Negotiate(r.Header, "Accept-Encoding", "gzip")
	return match == "gzip"
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Vary", "Accept-Encoding")

	if !h.shouldGzip(r) {
		h.h.ServeHTTP(w, r)
		return
	}

	gw := &responseWriter{
		ResponseWriter: w,

		h: h,

		buf: bufferPool.Get().(*[]byte),
	}
	defer func() {
		if err := gw.Close(); err != nil {
			httputils.RequestLogf(r, "gziphandler: error closing writer: %#v", err)
		}
	}()

	var rw http.ResponseWriter = gw

	_, cok := w.(http.CloseNotifier)
	_, hok := w.(http.Hijacker)
	_, pok := w.(http.Pusher)

	switch {
	case cok && hok:
		rw = closeNotifyHijackResponseWriter{gw}
	case cok && pok:
		rw = closeNotifyPusherResponseWriter{gw}
	case cok:
		rw = closeNotifyResponseWriter{gw}
	case hok:
		rw = hijackResponseWriter{gw}
	case pok:
		rw = pusherResponseWriter{gw}
	}

	h.h.ServeHTTP(rw, r)
}

// Gzip wraps an HTTP handler, to transparently gzip the
// response body if the client supports it (via the
// the Accept-Encoding header).
func Gzip(h http.Handler, opts ...Option) http.Handler {
	gzh := &handler{
		h: h,
		config: config{
			level:   DefaultCompression,
			minSize: defaultMinSize,
		},
	}

	for _, opt := range opts {
		opt(&gzh.config)
	}

	return gzh
}

// Wrapper returns a wrapper function (often known as
// middleware) which can be used to wrap an HTTP handler,
// to transparently gzip the response body if the client
// supports it (via the the Accept-Encoding header).
func Wrapper(opts ...Option) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return Gzip(h, opts...)
	}
}

type config struct {
	level        int
	minSize      int
	contentTypes []string
	shouldGzip   func(*http.Request) ShouldGzipType
}

// Option customizes the behaviour of the gzip handler.
type Option func(c *config)

// CompressionLevel is the gzip compression level to apply.
// See the level constants defined in this package.
//
// The default value adds gzip framing but performs no
// compression.
func CompressionLevel(level int) Option {
	if level < HuffmanOnly || level > BestCompression {
		panic("gziphandler: invalid compression level requested")
	}

	return func(c *config) {
		c.level = level
	}
}

// MinSize specifies the minimum size of a response before
// it will be compressed. Responses smaller than this value
// will not be compressed.
//
// If size is zero, all responses will be compressed.
//
// The default minimum size is 1400 bytes.
func MinSize(size int) Option {
	if size < 0 {
		panic("gziphandler: minimum size must not be negative")
	}

	return func(c *config) {
		c.minSize = size
	}
}

// ContentTypes specifies a list of MIME types to compare
// the Content-Type header to before compressing. If none
// match, the response will be returned as-is.
//
// MIME types are compared in a case-insensitive manner.
//
// A MIME type, but without any subtype, will match any
// more precise MIME type, i.e. image/* will match
// image/png, image/svg, image/gif and any other image
// types.
//
// Any directives that may be present in the Content-Type
// header will be skipped when comparing, i.e. text/html
// will match 'text/html; charset=utf-8'.
//
// By default, responses are gzipped regardless of
// Content-Type.
func ContentTypes(types []string) Option {
	types = append([]string(nil), types...)

	return func(c *config) {
		c.contentTypes = types
	}
}

// ShouldGzip provides control over when the handler should
// return a gzipped response. It allows handlers to implement
// logic that doesn't consult the request's Accept-Encoding
// header.
//
// By default, responses are only gzipped if the request's
// Accept-Encoding header indicates gzip support.
//
// Note: ShouldGzip does not affect MinSize or ContentTypes,
// it simply provides control over negotiating gzip support.
func ShouldGzip(fn func(*http.Request) ShouldGzipType) Option {
	return func(c *config) {
		c.shouldGzip = fn
	}
}

// ShouldGzipType controls how the handler determines gzip
// support.
type ShouldGzipType int

const (
	// NegotiateGzip defers gzip negotiation to the
	// request's Accept-Encoding header.
	NegotiateGzip ShouldGzipType = iota

	// SkipGzip skips gzipping the response.
	SkipGzip

	// ForceGzip ignores the request's Accept-Encoding
	// header and always gzips the response.
	// (See ShouldGzip note).
	ForceGzip
)

type (
	// Each of these structs is intentionally small (1 pointer wide) so
	// as to fit inside an interface{} without causing an allocaction.
	closeNotifyResponseWriter       struct{ *responseWriter }
	hijackResponseWriter            struct{ *responseWriter }
	pusherResponseWriter            struct{ *responseWriter }
	closeNotifyHijackResponseWriter struct{ *responseWriter }
	closeNotifyPusherResponseWriter struct{ *responseWriter }
)

var (
	_ http.CloseNotifier = closeNotifyResponseWriter{}
	_ http.CloseNotifier = closeNotifyHijackResponseWriter{}
	_ http.CloseNotifier = closeNotifyPusherResponseWriter{}
	_ http.Hijacker      = hijackResponseWriter{}
	_ http.Hijacker      = closeNotifyHijackResponseWriter{}
	_ http.Pusher        = pusherResponseWriter{}
	_ http.Pusher        = closeNotifyPusherResponseWriter{}
)

func (w closeNotifyResponseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (w closeNotifyHijackResponseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (w closeNotifyPusherResponseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (w hijackResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

func (w closeNotifyHijackResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

func (w pusherResponseWriter) Push(target string, opts *http.PushOptions) error {
	return w.ResponseWriter.(http.Pusher).Push(target, opts)
}

func (w closeNotifyPusherResponseWriter) Push(target string, opts *http.PushOptions) error {
	return w.ResponseWriter.(http.Pusher).Push(target, opts)
}
