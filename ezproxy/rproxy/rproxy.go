package rproxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"github.com/sirupsen/logrus"
	"fmt"
)

type proxy struct {
	*httputil.ReverseProxy
	log 		     *logrus.Logger
	originalDirector func(*http.Request)
}

func NewProxy(rhost string, rport int, log *logrus.Logger) (*proxy, error) {
	/*
	* Create a reverse proxy in front of the host
	*/

	// Obtain the host url from rhost and rport
	host := fmt.Sprintf("http://%s:%d", rhost, rport)
	url, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	// Create the proxy
	p := &proxy{
		ReverseProxy: 	  httputil.NewSingleHostReverseProxy(url), 
		log: 			  log,
		originalDirector: nil,
		}

	// Hook requests/responses
	p.originalDirector = p.Director
	p.Director = p.hookRequest()
	p.ModifyResponse = hookResponse
	p.ErrorHandler = errorHandler

	return p, nil
}

func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	/*
	* Wrapper for the ServeHTTP method for forwarding requests and returning responses
	*/

	return func (rw http.ResponseWriter, req *http.Request) {
		proxy.ServeHTTP(rw, req)
	}
}

func (p *proxy) hookRequest() func(req *http.Request) {
	/*
	* Hook incoming requests, and return the function that modifies them
	*/

	fmt.Println("hooking requests")

	return func(req *http.Request) {
		// First, call the original director
		p.originalDirector(req)

		// Modify the request
		p.log.Info(req)

		return
	}
}

func hookResponse(resp *http.Response) error {
	/*
	* Hook outgoing responses and modify them
	*/

	fmt.Printf("RESP: %v\n", resp)
	return nil
}

func errorHandler(rw http.ResponseWriter, req *http.Request, err error) {
	 /*
	 * Log any errors from the reverse proxy
	 */

	 fmt.Printf("ERROR: %v\n", err)
}
