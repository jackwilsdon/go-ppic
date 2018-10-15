// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httputils

import (
	"net"
	"time"
)

// TCPKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections.
//
// If period is zero, it defaults to three minutes.
//
// It is an unexported net/http type that is used by
// net/http.ListenAndServe and net/http.ListenAndServeTLS so dead
// TCP connections (e.g. closing laptop mid-download) eventually
// go away.
func TCPKeepAliveListener(ln *net.TCPListener, period time.Duration) net.Listener {
	if period == 0 {
		period = 3 * time.Minute
	}

	return &tcpKeepAliveListener{ln, period}
}

type tcpKeepAliveListener struct {
	ln *net.TCPListener

	period time.Duration
}

func (ln *tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.ln.AcceptTCP()
	if err != nil {
		return nil, err
	}

	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(ln.period)
	return tc, nil
}

func (ln *tcpKeepAliveListener) Close() error {
	return ln.ln.Close()
}

func (ln *tcpKeepAliveListener) Addr() net.Addr {
	return ln.ln.Addr()
}
