package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	/* Adding this http endpoint to download live profiles */
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
)

// WORD OF WARNING: DO NOT USE IN PROD!
// This is expeirmental and error-handling is terrbile below.

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
	var cpuFile *os.File
	if *cpuprofile != "" {
		cpuFile, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		// Can increase the frequency of profiling as needed
		runtime.SetCPUProfileRate(500)
		if err := pprof.StartCPUProfile(cpuFile); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
	}

	// Below added to enable memory profiling
	var memFile *os.File
	if *memprofile != "" {
		memFile, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(memFile); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // subscribe to system signals
	onKill := func(c chan os.Signal, cpuFile *os.File, memFile *os.File) {
		select {
		case <-c:
			if *cpuprofile != "" {
				defer cpuFile.Close()
				defer pprof.StopCPUProfile()
			}
			if *memprofile != "" {
				defer memFile.Close()
			}
			defer os.Exit(0)
		}
	}

	mux := http.NewServeMux()
	mux.Handle("/", &helloHandler{})
	mux.Handle("/health", &healthHandler{})

	go onKill(c, cpuFile, memFile)

	fmt.Println("I've made it here")
	http.ListenAndServe(fmt.Sprintf(":%d", *addrFlag), mux)
}
