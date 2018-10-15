// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License license that can be found in
// the LICENSE file.

package httputils

import "net/http"

// Negotiate performs HTTP content negotiation via the
// Accept, Accept-Charset, Accept-Encoding, and
// Accept-Language HTTP headers.
//
// The highest ranking match is returned, or an empty
// string if none found.
//
// It does not resolve wildcard MIME subtypes when used
// with the Accept header. If this is important to you,
// include the relevant wildcards in offers.
//
// It does not resolve language subtypes, such as en-US to
// en. If this is important to you, include the relevant
// subtype in offers.
//
// It also does not handle the rejected identity
// (identity;q=0 or *;q=0 or */*;q=0), instead returning
// an empty string.
func Negotiate(header http.Header, name string, offers ...string) (match string) {
	specs := parseAccept(header, name)

	Q := 0.0
	for _, offer := range offers {
		for _, spec := range specs {
			if spec.Q > Q && tokenEqual(offer, spec.Value) {
				match, Q = offer, spec.Q
			}
		}
	}

	return match
}
