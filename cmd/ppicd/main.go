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
	portNum := flag.String("p", "3000", "Port number for the server to run on.")
	debug := flag.Bool("d", false, "Use to turn on debug option")
	gzip := flag.Bool("g", false, "Use to turn on gzip option")
	host := flag.String("h", "localhost", "Define the host for the server")

	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", ppic.Handler)

	// Enable profiling URLs if the debug option is set.
	if *debug {
		fmt.Println("Debug enabled")
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	var handler http.Handler = mux

	// Enable GZIP if it's not disabled.
	if *gzip {
		fmt.Println("gzip enabled")
		handler = gziphandler.Gzip(mux)
	}

	addr := fmt.Sprintf("%s:%s", *host, *portNum)
	fmt.Printf("Server at: %s\n", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
	}
}
