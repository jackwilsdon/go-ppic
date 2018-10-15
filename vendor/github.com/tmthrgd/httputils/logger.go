// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License license that can be found in
// the LICENSE file.

package httputils

import (
	"fmt"
	"log"
	"net/http"
)

// RequestLogf prints a message to the logger attached to
// the *http.Server that served the request.
// If nil, logging goes to os.Stderr via the log package's
// standard logger.
//
// Arguments are handled in the manner of fmt.Printf.
func RequestLogf(r *http.Request, format string, v ...interface{}) {
	srv, _ := r.Context().Value(http.ServerContextKey).(*http.Server)

	// Printf is defined as:
	// 	l.Output(2, fmt.Sprintf(format, v...))
	// by inlining, we skip our own frame.

	s := fmt.Sprintf(format, v...)

	if srv != nil && srv.ErrorLog != nil {
		srv.ErrorLog.Output(2, s)
	} else {
		log.Output(2, s)
	}
}
