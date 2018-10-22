package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"

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

	// Manually create the listener so we can work out what port it's on.
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	log.Printf("Starting server on http://%s...\n", addressToString(listener.Addr()))

	// Start serving on the listener.
	if err := http.Serve(listener, handler); err != nil {
		log.Fatalf("error: %s\n", err)
		os.Exit(1)
	}
}
