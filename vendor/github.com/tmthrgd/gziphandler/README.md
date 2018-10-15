Gzip Handler
============

[![GoDoc](https://godoc.org/github.com/tmthrgd/gziphandler?status.svg)](https://godoc.org/github.com/tmthrgd/gziphandler)
[![Build Status](https://travis-ci.org/tmthrgd/gziphandler.svg?branch=master)](https://travis-ci.org/tmthrgd/gziphandler)
[![Go Report Card](https://goreportcard.com/badge/github.com/tmthrgd/gziphandler)](https://goreportcard.com/report/github.com/tmthrgd/gziphandler)
[![Coverage Status](https://coveralls.io/repos/github/tmthrgd/gziphandler/badge.svg?branch=master)](https://coveralls.io/github/tmthrgd/gziphandler?branch=master)

This is a tiny Go package which wraps HTTP handlers to transparently gzip the
response body, for clients which support it. Although it's usually simpler to
leave that to a reverse proxy (like nginx or Varnish), this package is useful
when that's undesirable.

## Usage

Call `Gzip` with any handler (an object which implements the `http.Handler`
interface), and it'll return a new handler which gzips the response. For
example:

```go
package main

import (
	"io"
	"net/http"

	"github.com/tmthrgd/gziphandler"
)

func main() {
	withoutGz := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, "Hello, World")
	})

	withGz := gziphandler.Gzip(withoutGz)

	http.Handle("/", withGz)
	http.ListenAndServe("0.0.0.0:8000", nil)
}
```

## License

[Apache 2.0][license].

[license]:  https://github.com/tmthrgd/gziphandler/blob/master/LICENSE.md
