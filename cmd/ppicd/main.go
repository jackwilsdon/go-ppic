package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/felixge/httpsnoop"
	"github.com/jackwilsdon/go-ppic"
	"github.com/tmthrgd/gziphandler"
)

// addressToString converts an address to a string (who'd have thought!).
// This method will use "127.0.0.1" if the IP is unspecified.
func addressToString(addr net.Addr) string {
	tcp, isTCP := addr.(*net.TCPAddr)

	if !isTCP {
		return "unknown"
	}

	ip := tcp.IP.String()

	if tcp.IP.IsUnspecified() {
		ip = "127.0.0.1"
	}

	return fmt.Sprintf("%s:%d", ip, tcp.Port)
}

// withLogger wraps the specified handler in a logger and returns it.
func withLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Execute the handler and capture metrics.
		m := httpsnoop.CaptureMetrics(h, res, req)

		// Log the request.
		log.Printf("%s %s - %d (%s)", req.Method, req.URL.RequestURI(), m.Code, m.Duration)
	})
}

func main() {
	// Build a list of the flags we support.
	host := flag.String("h", "", "host to run the server on")
	port := flag.Uint("p", 3000, "port to run the server on")
	debug := flag.Bool("d", false, "enable pprof debug routes")
	gzip := flag.Bool("z", false, "enable gzip compression")
	verbose := flag.Bool("v", false, "enable verbose output")

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

	// Manually create the listener so we can work out what port it's on.
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	log.Printf("starting server on http://%s...\n", addressToString(listener.Addr()))

	// Wrap the handler in a logger if verbose mode is enabled.
	if *verbose {
		handler = withLogger(handler)
	}

	// Start serving on the listener.
	if err := http.Serve(listener, handler); err != nil {
		log.Fatalf("error: %s\n", err)
		os.Exit(1)
	}
}
