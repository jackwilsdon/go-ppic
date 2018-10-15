// Copyright 2013 The Go Authors. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.

package httputils

import (
	"net/http"
	"strings"
)

type acceptSpec struct {
	Value string
	Q     float64
}

func parseAccept(header http.Header, key string) (specs []acceptSpec) {
loop:
	for _, s := range header[http.CanonicalHeaderKey(key)] {
		s = trimOWS(s)
		for {
			var spec acceptSpec
			spec.Value, s = expectTokenSlash(s)
			if spec.Value == "" {
				continue loop
			}
			spec.Q = 1.0
			s = trimOWS(s)
			if strings.HasPrefix(s, ";") {
				s = trimOWS(s[1:])
				if len(s) < 2 || !tokenEqual(s[:2], "q=") {
					continue loop
				}
				spec.Q, s = expectQuality(s[2:])
				if spec.Q < 0.0 {
					continue loop
				}
			}
			specs = append(specs, spec)
			s = trimOWS(s)
			if !strings.HasPrefix(s, ",") {
				continue loop
			}
			s = trimOWS(s[1:])
		}
	}
	return
}

func expectTokenSlash(s string) (token, rest string) {
	i := 0
	for ; i < len(s); i++ {
		b := s[i]
		if !isTokenRune(rune(b)) && b != '/' {
			break
		}
	}
	return s[:i], s[i:]
}

func expectQuality(s string) (q float64, rest string) {
	switch {
	case len(s) == 0:
		return -1, ""
	case s[0] == '0':
		q = 0
	case s[0] == '1':
		q = 1
	default:
		return -1, ""
	}
	s = s[1:]
	if !strings.HasPrefix(s, ".") {
		return q, s
	}
	s = s[1:]
	i := 0
	n := 0
	d := 1
	for ; i < len(s); i++ {
		b := s[i]
		if b < '0' || b > '9' {
			break
		}
		n = n*10 + int(b) - '0'
		d *= 10
	}
	// RFC 7231 section 5.3.1: Quality Values:
	//  A sender of qvalue MUST NOT generate more than three digits after the
	//  decimal point.  User configuration of these values ought to be
	//  limited in the same fashion.
	if i > 3 {
		return -1, ""
	}
	return q + float64(n)/float64(d), s[i:]
}
