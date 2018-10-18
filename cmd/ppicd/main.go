package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/jackwilsdon/go-ppic"
	"github.com/tmthrgd/gziphandler"
)

func main() {
	// Build a list of the flags we support.
	host := flag.String("h", "", "host to run the server on")
	port := flag.Uint("p", 3000, "port to run the server on")
	debug := flag.Bool("d", false, "enable pprof debug routes")
	gzip := flag.Bool("z", false, "enable gzip compression")

	// Parse the command-line flags.
	flag.Parse()

	// Create a new server with our handler.
	mux := http.NewServeMux()
	mux.HandleFunc("/", ppic.Handler)

	// Enable pprof debug routes if the debug flag is set.
	if *debug {
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	var handler http.Handler = mux

	// Enable gzip compression if the gzip flag is set.
	if *gzip {
		handler = gziphandler.Gzip(handler)
	}

	// Build the address from the host and port.
	addr := fmt.Sprintf("%s:%d", *host, *port)

	if err := http.ListenAndServe(addr, handler); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
	}
}
