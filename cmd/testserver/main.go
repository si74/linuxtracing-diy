package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	/* Adding this http endpoint to download live profiles */
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
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

var addrFlag = flag.Int("addrFlag", 8080, "address on which to serve app")

// Note: These flags were added to consider cases of cpu or memory profiling
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {

	/* Validate that the port provided is legit */
	flag.Parse()

	if *addrFlag < 1 || *addrFlag > 65535 {
		log.Fatal("provided flag must be between values of 1 and 65535")
	}

	// Below added to enable cpu profiling
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	mux := http.NewServeMux()
	mux.Handle("/", &helloHandler{})
	mux.Handle("/health", &healthHandler{})

	http.ListenAndServe(fmt.Sprintf(":%d", *addrFlag), mux)

	// Below added to enable memory profiling
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}
