package main

import (
	"flag"
	"io"
	"log"
	"net/http"
)

var addr = new(string)

func init() {
	flag.StringVar(addr, "addr", "localhost:8080", "The address to listen on")
}

func main() {
	flag.Parse()

	http.HandleFunc("/", simpleProxy)

	log.Printf("Starting server on %s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

// simpleProxy is a simple reverse proxy that forwards requests to the target
// using the simplest as dump implementation possible.
// This is too simple to support all the cases in the wild, such as protocol upgrades,
// streaming, SSE, Websocket, h2c, etc.
func simpleProxy(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request for %s\n", r.URL.String())
	res, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer res.Body.Close()

	// Copy the headers from the response to the writer
	for k, v := range res.Header {
		for _, v := range v {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(res.StatusCode)

	if _, err := io.Copy(w, res.Body); err != nil {
		log.Printf("Error: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
