package main

import (
	"flag"
	"io"
	"log"
	"net/http"

	"github.com/emicklei/go-restful"
)

var (
	certFile, keyFile *string
)

func init() {
	certFile = flag.String("cert", "wrapper.crt", "cert file")
	keyFile = flag.String("key", "wrapper.key", "private key file")
	flag.Parse()
}

func main() {
	restful.Filter(globalLogging)
	ws := new(restful.WebService)
	ws.Route(ws.GET("/hello").To(hello))
	restful.Add(ws)
	log.Fatal(http.ListenAndServeTLS(":8080", *certFile, *keyFile, http.DefaultServeMux))
}

func hello(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, "wrapper")
}

// Global Filter
func globalLogging(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	log.Printf("[global-filter (logger)] %s, %s\n", req.Request.Method, req.Request.URL)
	chain.ProcessFilter(req, resp)
}
