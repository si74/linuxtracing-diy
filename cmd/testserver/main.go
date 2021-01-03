package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

type helloHandler struct {
}

func (*helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	io.WriteString(w, "hi there!")
}

type healthHandler struct {
}

func (*healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	io.WriteString(w, "okay!")
}

var addrFlag int

func main() {

	/* Validate that the port provided is legit */
	flag.IntVar(&addrFlag, "addrFlag", 8080, "address on which to serve app")
	flag.Parse()

	if addrFlag < 1 || addrFlag > 65535 {
		log.Fatal("provided flag must be between values of 1 and 65535")
	}

	mux := http.NewServeMux()
	mux.Handle("/", &helloHandler{})
	mux.Handle("/health", &healthHandler{})

	http.ListenAndServe(fmt.Sprintf(":%d", addrFlag), mux)
}
