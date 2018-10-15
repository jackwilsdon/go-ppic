// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License license that can be found in
// the LICENSE file.

package httputils

import (
	"context"
	"net/http"

	"golang.org/x/sync/errgroup"
)

// Shutdown gracefully shutsdown several http.Server's.
func Shutdown(ctx context.Context, srvs ...*http.Server) error {
	var eg errgroup.Group

	for _, srv := range srvs {
		if srv == nil {
			continue
		}

		srv := srv
		eg.Go(func() error {
			return srv.Shutdown(ctx)
		})
	}

	return eg.Wait()
}
