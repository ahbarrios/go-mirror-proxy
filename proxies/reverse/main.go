package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
)

var addr = new(string)

func init() {
	flag.StringVar(addr, "addr", "localhost:8080", "The address to listen on")
}

func main() {
	flag.Parse()

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			log.Printf("Received request for %s\n", r.In.URL.String())
			r.SetURL(r.In.URL)
			r.SetXForwarded()
			r.Out.Host = r.In.Host // if desired
		},
	}

	log.Printf("Starting server on %s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, proxy))
}
