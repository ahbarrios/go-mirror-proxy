package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"
	"sync"
	"time"
)

var trace = &httptrace.ClientTrace{
	GotConn: func(connInfo httptrace.GotConnInfo) {
		log.Printf("Got Conn: %+v\n", connInfo)
	},
	DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
		log.Printf("DNS Info: %+v\n", dnsInfo)
	},
}

var (
	addr    = new(string)
	proxy   = new(string)
	mirrors = new(string)
)

func init() {
	flag.StringVar(addr, "addr", "localhost:80", "The address to listen on")
	flag.StringVar(proxy, "proxy", "", "The main proxy address to forward requests to")
	flag.StringVar(mirrors, "mirrors", "", "This is a comma separated list of mirrors proxies to shadowing traffic")
}

// This is the main entrypoint for the gateway.
// It should be used to start a HTTP gateway to forwarding traffic to the specified proxy server
// and shadowing traffic to mirrors.
func main() {
	flag.Parse()
	if *proxy == "" {
		log.Fatal("The proxy address is required")
	}

	http.HandleFunc("/", gateway(*proxy, extractMirrors()...))

	log.Printf("Starting gateway on %s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

// extractMirrors it should be used to extract the mirrors from the `mirrors` command line arguments
func extractMirrors() (mrs []string) {
	mrs = strings.Split(*mirrors, ",")
	// clean up extra spaces
	for i, m := range mrs {
		mrs[i] = strings.TrimSpace(m)
	}
	return
}

// newTransportWithProxy it stablish safe defaults inherited from http.DefaultTransport and set a proxy
func newTransportWithProxy(proxy *url.URL) *http.Transport {
	return &http.Transport{
		Proxy:                 http.ProxyURL(proxy),
		DialContext:           (&net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

func gateway(proxy string, mirrors ...string) http.HandlerFunc {

	// Safely parse the proxies URLs first and panic if any of them is invalid
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		log.Fatalf("Error: %s\n", err.Error())
	}
	proxyTransport := newTransportWithProxy(proxyURL) // main proxy RoundTripper

	mirrorTransports := make([]*http.Transport, len(mirrors)) // mirrors proxies RoundTrippers
	for i, m := range mirrors {
		mURL, err := url.Parse(m)
		if err != nil {
			log.Fatalf("Error: %s\n", err.Error())
		}
		mirrorTransports[i] = newTransportWithProxy(mURL)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(httptrace.WithClientTrace(r.Context(), trace))

		// Forward the request to the main proxy
		res, err := proxyTransport.RoundTrip(r)
		if err != nil {
			log.Printf("Error: %s\n", err.Error())
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		defer res.Body.Close()

		// Shadow the request to the mirrors
		go shadowTraffic(r.Clone(context.Background()), mirrorTransports)

		// Copy the headers from the response to the writer
		for k, v := range res.Header {
			for _, v := range v {
				w.Header().Add(k, v)
			}
		}
		w.WriteHeader(res.StatusCode)

		// Deliver the response from main proxy to the client
		if _, err := io.Copy(w, res.Body); err != nil {
			log.Printf("Error: %s\n", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func shadowTraffic(r *http.Request, mirrors []*http.Transport) {
	var wg sync.WaitGroup
	wg.Add(len(mirrors))
	for _, t := range mirrors {
		go func(transport *http.Transport) {
			defer wg.Done()
			proxy, err := transport.Proxy(r)
			if err != nil {
				log.Printf("Couldn't extract proxy: %v - Error: %s\n", transport, err.Error())
				return
			}

			res, err := transport.RoundTrip(r)
			if err != nil {
				log.Printf("%s - Error: %v\n", proxy, err.Error())
				return
			}
			defer res.Body.Close()
			log.Printf("%s - Shadowed request for %s with response %v\n", proxy, r.URL.String(), res.Status)
		}(t)
	}
	wg.Wait()
}
