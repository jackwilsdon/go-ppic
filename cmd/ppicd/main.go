package main

import (
	"fmt"
	"github.com/jackwilsdon/go-ppic"
	"github.com/tmthrgd/gziphandler"
	"net/http"
	"net/http/pprof"
	"os"
	"strconv"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", ppic.Handler)

	// Enable profiling URLs if the debug option is set.
	if os.Getenv("DEBUG") == "1" {
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	host := os.Getenv("HOST")
	port := 3000

	// Make sure we have a valid port.
	if portString, ok := os.LookupEnv("PORT"); ok {
		var err error

		port, err = strconv.Atoi(portString)

		if _, isNum := err.(*strconv.NumError); isNum {
			fmt.Fprintf(os.Stderr, "error: invalid port\n")
			os.Exit(1)
		}
	}

	var handler http.Handler = mux

	// Enable GZIP if it's not disabled.
	if os.Getenv("GZIP") != "0" {
		handler = gziphandler.Gzip(mux)
	}

	addr := fmt.Sprintf("%s:%d", host, port)

	if err := http.ListenAndServe(addr, handler); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
	}
}
