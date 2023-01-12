package reverseProxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"fmt"
)

func NewProxy(rhost string, rport int) (*httputil.ReverseProxy, error) {
	/*
	* Create a reverse proxy in front of the host
	*/

	// Obtain the host url from rhost and rport
	host := fmt.Sprintf("http://%s:%d", rhost, rport)
	url, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	// Create the proxy and hook requests / responses
	proxy := httputil.NewSingleHostReverseProxy(url)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		hookRequest(req)
	}
	proxy.ModifyResponse = hookResponse
	proxy.ErrorHandler = errorHandler

	return proxy, nil
}

func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	/*
	* Wrapper for the ServeHTTP method for forwarding requests and returning responses
	*/

	return func (rw http.ResponseWriter, req *http.Request) {
		proxy.ServeHTTP(rw, req)
	}
}

func hookRequest(req *http.Request) {
	/*
	* Hook incoming requests and modify them
	*/

	fmt.Printf("REQ: %v\n", req)
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