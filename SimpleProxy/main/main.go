package main

import (
	"SimpleProxy/reverseProxy"
	"net/http"
	"flag"
	"log"
	"fmt"
)

func parseArgs() (string, int, int) {
	/*
	* Parse cmd line args for proxy params
	*/

	var rhost string
	var rport, lport int

	// Set flags. Note: I want to try colorizing this
	flag.StringVar(&rhost, "rhost", "", "The host to be proxied")
	flag.IntVar(&rport, "rport", 80, "The port of the host to be proxied")
	flag.IntVar(&lport, "lport", 8080, "The port the proxy will listen on")
	flag.Parse()

	if rhost == "" {
		// Require rhost, otherwise exit
		log.Fatal("Missing required argument --rhost")
	}

	return rhost, rport, lport
}

func main() {
	rhost, rport, lport := parseArgs()

	proxy, err := reverseProxy.NewProxy(rhost, rport)
	if err != nil {
		log.Fatal(err)
	}

	// Register the reverse proxy as the handler for all incoming requests
	http.HandleFunc("/", reverseProxy.ProxyRequestHandler(proxy))
	
	// Start the http server
	// Note: More control over the server's behavior is available by creating
	// a custom Server
	lportString := fmt.Sprintf(":%d", lport)
	log.Fatal(http.ListenAndServe(lportString, nil))
}