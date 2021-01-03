package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var addrFlag string

func main() {

	/* Validate that the port provided is legit */
	flag.StringVar(&addrFlag, "addrFlag", "http://localhost:8080", "address on which to serve app")
	flag.Parse()

	resp, err := http.Get(addrFlag)
	if err != nil {
		log.Printf("cannot make request: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("cannot read body: %v", err)
		return
	}
	fmt.Println(resp.Status)
	fmt.Println(string(body))
}
