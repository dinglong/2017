package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/emicklei/go-restful"
)

var (
	caFile, certFile, keyFile *string
	host                      *string
)

func init() {
	host = flag.String("host", "127.0.0.1", "master host")
	caFile = flag.String("ca", "root.ca.crt", "ca file")
	certFile = flag.String("cert", "proxy.crt", "cert file")
	keyFile = flag.String("key", "proxy.key", "private key file")
	flag.Parse()
}

func main() {
	http.DefaultServeMux.HandleFunc("/proxy/", proxy)

	restful.Filter(globalLogging)
	ws := new(restful.WebService)
	ws.Route(ws.GET("/hello").To(hello))
	restful.Add(ws)

	log.Fatal(http.ListenAndServeTLS(":8180", *certFile, *keyFile, http.DefaultServeMux))
}

func proxy(resp http.ResponseWriter, req *http.Request) {
	log.Printf("[in proxy (logger)] %s, %s\n", req.Method, req.URL)

	proxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "https", Host: *host})
	proxy.Transport = &http.Transport{
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       makeTLSConfig(),
	}

	path := strings.TrimPrefix(req.URL.Path, "/proxy")
	newReq, err := http.NewRequest(req.Method, path, req.Body)
	if err != nil {
		log.Fatalf("new request error, %v\n", err)
	}

	proxy.ServeHTTP(resp, newReq)

	log.Printf("[out proxy (logger)]\n")
}

func hello(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, "proxy")
}

func makeTLSConfig() *tls.Config {
	ca, err := ioutil.ReadFile(*caFile)
	if err != nil {
		log.Fatalf("read ca file error, %v\n", err)
	}

	cert, err := ioutil.ReadFile(*certFile)
	if err != nil {
		log.Fatalf("read cert file error, %v\n", err)
	}

	key, err := ioutil.ReadFile(*keyFile)
	if err != nil {
		log.Fatalf("read key file error, %v\n", err)
	}

	pool := x509.NewCertPool()
	if ok := pool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("append ca error, %v\n", err)
	}

	certificates, err := tls.X509KeyPair(cert, key)
	if ok := pool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("X509KeyPair error, %v\n", err)
	}

	return &tls.Config{
		RootCAs:      pool,
		Certificates: []tls.Certificate{certificates},
		// InsecureSkipVerify: true,
	}
}

// Global Filter
func globalLogging(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	log.Printf("[global-filter (logger)] %s, %s\n", req.Request.Method, req.Request.URL)
	chain.ProcessFilter(req, resp)
}
