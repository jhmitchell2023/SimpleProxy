package main

import (
	"net/http"
	"time"
	"flag"
	"fmt"
)

func parseArgs() int {
	/*
	* Parse cmd line args
	*/

	var lport int

	flag.IntVar(&lport, "lport", 8080, "Port to listen on")
	flag.Parse()

	return lport
}

func main() {
	/*
	* Create an http server listening on port 8080
	* Respond to incoming requests with "Hello!"
	*/

	http.HandleFunc("/", func (rw http.ResponseWriter, req *http.Request) {
		fmt.Printf("[SERVER] request receieved at: %s\n", time.Now())
		fmt.Fprint(rw, "Hello!\n")
	})

	lport := parseArgs()
	lportString := fmt.Sprintf(":%d", lport)
	http.ListenAndServe(lportString, nil)
}