package main

import (
	"ezproxy/rproxy"
	"net/http"
	"flag"
	"github.com/sirupsen/logrus"
	"os"
	"fmt"
)

func parseArgs() (string, int, int, *logrus.Logger) {
	/*
	* Parse cmd line args for proxy params
	*/

	var rhost, logfile string
	var rport, lport int
	var log = logrus.New()

	// Set flags. Note: I want to try colorizing this
	flag.StringVar(&rhost, "rhost", "", "The host to be proxied")
	flag.IntVar(&rport, "rport", 80, "The port of the host to be proxied")
	flag.IntVar(&lport, "lport", 8080, "The port the proxy will listen on")
	flag.StringVar(&logfile, "logging", "", "Logfile name")
	flag.Parse()

	if rhost == "" {
		// Require rhost, otherwise exit
		log.Fatal("Missing required argument --rhost")
	}

	if logfile != "" {
		// Create a new logger using logrus
		file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY, 0666)

		// Set the logger to write to file; if not specified, write to stderr
		if err == nil {
			log.Out = file
		} else {
			log.Info("No logfile specified, using default stderr")
		}
	}

	// Set logging format to json
	log.SetFormatter(&logrus.JSONFormatter{})

	return rhost, rport, lport, log
}

func main() {
	rhost, rport, lport, log := parseArgs()

	p, err := rproxy.NewProxy(rhost, rport, log)
	if err != nil {
		log.Fatal(err)
	}

	// Register the reverse proxy as the handler for all incoming requests
	http.HandleFunc("/", rproxy.ProxyRequestHandler(p.ReverseProxy))
	
	// Start the http server
	// Note: More control over the server's behavior is available by creating
	// a custom Server
	lportString := fmt.Sprintf(":%d", lport)
	log.Fatal(http.ListenAndServe(lportString, nil))
}