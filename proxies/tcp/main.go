package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

var addr = new(string)

func init() {
	flag.StringVar(addr, "addr", "localhost:8080", "The address to listen on")
}

func main() {
	flag.Parse()

	log.Printf("Starting server on %s\n", *addr)
	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error: %s\n", err.Error())
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	req, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		return
	}

	port := ":80"
	if p := req.URL.Port(); p != "" {
		port = ":" + p
	}

	rconn, err := net.DialTimeout("tcp", req.Host+port, 30*time.Second)
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		conn.Write([]byte("HTTP/1.1 502 Bad Gateway\r\n\r\n"))
		return
	}
	defer rconn.Close()

	// request
	if err := req.Write(rconn); err != nil {
		log.Printf("Error: %s\n", err.Error())
		return
	}
	// response
	if _, err := io.Copy(conn, rconn); err != nil {
		log.Printf("Error: %s\n", err.Error())
		return
	}
}
