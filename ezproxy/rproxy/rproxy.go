package rproxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"github.com/sirupsen/logrus"
	//"io"
	"fmt"
	//"bytes"
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
	p.ModifyResponse = p.hookResponse()
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

	return func(req *http.Request) {
		// First, call the original director
		p.originalDirector(req)

		// Modify/log the request
		//msg := req.Header.Get("msg")
		p.log.Info("Request received")
		p.log.SetFormatter(&logrus.JSONFormatter{})
		p.log.Info(req)
		p.log.SetFormatter(&logrus.TextFormatter{})

		return
	}
}

func (p *proxy) hookResponse() func(*http.Response) error {
	/*
	* Hook outgoing responses and modify them
	*/

	return func(resp *http.Response) error {
		// Modify/log the response
		//buf, err := io.ReadAll(resp.Body)
        //if err != nil {
        //    return err
        //}
		//rdr1 := io.NopCloser(bytes.NewBuffer(buf))
		//rdr2 := io.NopCloser(bytes.NewBuffer(buf))
		//body := String(rdr1)
		//resp.Body = rdr2

		//msg := resp.Header.Get("User-Agent")
		p.log.Info("Response received")
		p.log.SetFormatter(&logrus.JSONFormatter{})
		p.log.Info(resp)
		p.log.SetFormatter(&logrus.TextFormatter{})

		return nil
	}
}

func errorHandler(rw http.ResponseWriter, req *http.Request, err error) {
	 /*
	 * Log any errors from the reverse proxy
	 */

	 fmt.Printf("ERROR: %v\n", err)
}
