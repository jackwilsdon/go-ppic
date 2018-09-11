package main

import (
	"fmt"
	"github.com/jackwilsdon/go-ppic"
	"net/http"
	"os"
	"strconv"
)

func main() {
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

	addr := fmt.Sprintf("%s:%d", host, port)


	if err := http.ListenAndServe(addr, http.HandlerFunc(ppic.Handler)); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
	}
}
